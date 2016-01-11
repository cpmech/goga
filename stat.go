// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

type Stat struct {

	// stat
	Nfeval   int           // number of function evaluations
	XfltBest [][]float64   // best results after RunMany
	XintBest [][]int       // best results after RunMany
	SysTime  time.Duration // system (real/CPU) time

	// check multi-objective
	F1F0func   func(f0 float64) float64 // f1(f0) function
	F1F0err    []float64                // max(error(f1))
	F1F0spread []float64                // spreading on (f0,f1) space
}

// RunMany runs many trials in order to produce statistical data
func (o *Optimiser) RunMany(dirout, fnkey string) {

	// benchmark
	t0 := time.Now()
	defer func() {
		o.SysTime = time.Now().Sub(t0)
	}()

	// disable verbose flag temporarily
	if o.Verbose {
		defer func() {
			o.Verbose = true
		}()
		o.Verbose = false
	}

	// remove previous results
	if fnkey != "" {
		io.RemoveAll(dirout + "/" + fnkey + "-*.res")
	}

	// perform trials
	for itrial := 0; itrial < o.Ntrials; itrial++ {

		// re-generate solutions
		o.Nfeval = 0
		if itrial > 0 {
			o.generate_solutions(itrial)
		}

		// save initial solutions
		if fnkey != "" {
			WriteAllValues(dirout, io.Sf("%s-%04d_ini", fnkey, itrial), o)
		}

		// solve
		o.Solve()

		// sort
		if o.Nova < 2 {
			SortByOva(o.Solutions, 0)
		} else {
			SortByFrontThenOva(o.Solutions, 0)
		}

		// find best
		if o.Solutions[0].Feasible() {
			xf, xi := o.Solutions[0].GetCopyResults()
			if o.Nflt > 0 {
				o.XfltBest = append(o.XfltBest, xf)
			}
			if o.Nflt > 0 {
				o.XintBest = append(o.XintBest, xi)
			}
		}

		// check multi-objective results
		if o.F1F0func != nil {
			var f1_err_max float64
			var found bool
			for _, sol := range o.Solutions {
				if sol.FrontId == 0 {
					f0, f1 := sol.Ova[0], sol.Ova[1]
					f1_cor := o.F1F0func(f0)
					f1_err := math.Abs(f1 - f1_cor)
					f1_err_max = utl.Max(f1_err_max, f1_err)
					found = true
				}
			}
			if found {
				o.F1F0err = append(o.F1F0err, f1_err_max)
			}
		}

		// spreading on multi-objective problems
		if o.Nova > 1 {
			if o.Solutions[0].FrontId == 0 && o.Solutions[1].FrontId == 0 {
				dist := 0.0
				for i := 1; i < o.Nsol; i++ {
					if o.Solutions[i].FrontId == 0 && o.Solutions[i].DistNeigh > 1e-7 {
						F0, F1 := o.Solutions[i-1].Ova[0], o.Solutions[i-1].Ova[1]
						f0, f1 := o.Solutions[i].Ova[0], o.Solutions[i].Ova[1]
						dist += math.Sqrt(math.Pow(f0-F0, 2.0) + math.Pow(f1-F1, 2.0))
					}
				}
				if false {
					io.Pforan("dist = %v\n", dist)
				}
				o.F1F0spread = append(o.F1F0spread, dist)
			}
		}

		// save final solutions
		if fnkey != "" {
			f0min := o.Solutions[0].Ova[0]
			for _, sol := range o.Solutions {
				f0min = utl.Min(f0min, sol.Ova[0])
			}
			WriteAllValues(dirout, io.Sf("%s-%04d_f0min=%g", fnkey, itrial, f0min), o)
		}

		// debug
		if false {
			plt.Reset()
			PlotOvaOvaPareto(io.Sf("fig_flt05_%d", itrial), o, nil, 0, 1, func() {
				np := 101
				F0 := utl.LinSpace(0, 1, np)
				F1 := make([]float64, np)
				for i := 0; i < np; i++ {
					F1[i] = o.F1F0func(F0[i])
				}
				plt.Plot(F0, F1, "'b-'")
			}, nil, false)
		}

	}
}

// StatMinProb prints statistical analysis when using MinProb
func (o *Optimiser) StatMinProb(idxF, hlen int, Fref float64, verbose bool) (fmin, fave, fmax, fdev float64, F []float64) {
	if o.MinProb == nil {
		io.Pfred("_warning_ MinProb is <nil>\n")
		return
	}
	nfb := len(o.XfltBest)
	nib := len(o.XintBest)
	if nfb+nib == 0 {
		fmin, fave, fmax, fdev = 1e30, 1e30, 1e30, 1e30
		io.Pfred("_warning_ XfltBest and XintBest are not available. Call RunMany first.\n")
		return
	}
	nbest := utl.Imax(nfb, nib)
	var xf []float64
	var xi []int
	F = make([]float64, nbest)
	cpu := 0
	for i := 0; i < nbest; i++ {
		if nfb > 0 {
			xf = o.XfltBest[i]
		}
		if nib > 0 {
			xi = o.XintBest[i]
		}
		o.MinProb(o.F[cpu], o.G[cpu], o.H[cpu], xf, xi, cpu)
		F[i] = o.F[cpu][idxF]
	}
	if nbest < 2 {
		fmin, fave, fmax = F[0], F[0], F[0]
		return
	}
	fmin, fave, fmax, fdev = rnd.StatBasic(F, true)
	if verbose {
		io.Pf("fmin = %v\n", fmin)
		io.PfYel("fave = %v (%v)\n", fave, Fref)
		io.Pf("fmax = %v\n", fmax)
		io.Pf("fdev = %v\n\n", fdev)
		io.Pf(rnd.BuildTextHist(nice_num(fmin-0.05, 2), nice_num(fmax+0.05, 2), 11, F, "%.2f", hlen))
	}
	return
}

// StatMultiObj prints statistical analysis for multi-objective problems
//  emin, eave, emax, edev -- errors on f1(f0)
func (o *Optimiser) StatMultiObj(hlen int, SpreadRef float64, verbose bool) (emin, eave, emax, edev, smin, save, smax, sdev float64) {
	if len(o.F1F0err) > 2 {
		emin, eave, emax, edev = rnd.StatBasic(o.F1F0err, true)
		if verbose {
			io.Pf("\nerror on Pareto front\n")
			io.Pf("emin = %v\n", emin)
			io.PfYel("eave = %v\n", eave)
			io.Pf("emax = %v\n", emax)
			io.Pf("edev = %v\n\n", edev)
			io.Pf(rnd.BuildTextHist(nice_num(emin-0.05, 2), nice_num(emax+0.05, 2), 11, o.F1F0err, "%.2f", hlen))
		}
	}
	if len(o.F1F0spread) > 2 {
		S := make([]float64, len(o.F1F0spread))
		for i, s := range o.F1F0spread {
			S[i] = s / SpreadRef
		}
		smin, save, smax, sdev = rnd.StatBasic(S, true)
		if verbose {
			io.Pf("\nspreading on Pareto front (ref = %g)\n", SpreadRef)
			io.Pf("smin = %v\n", smin)
			io.PfYel("save = %v\n", save)
			io.Pf("smax = %v\n", smax)
			io.Pf("sdev = %v\n\n", sdev)
			io.Pf(rnd.BuildTextHist(nice_num(smin-0.05, 2), nice_num(smax+0.05, 2), 11, o.F1F0spread, "%.2f", hlen))
		}
	}
	return
}

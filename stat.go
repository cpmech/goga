// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

type Stat struct {

	// stat
	Nfeval  int           // number of function evaluations
	SysTime time.Duration // system (real/CPU) time

	// formatting data for reports
	RptName    string    // problem name
	RptFref    []float64 // reference OVAs
	RptXref    []float64 // reference flts
	RptFmin    []float64 // min OVAs for reports/graphs
	RptFmax    []float64 // max OVAs for reports/graphs
	RptFmtF    string    // format for fmin, fave and fmax
	RptFmtFdev string    // format for fdev
	RptFmtE    string    // format for emin, eave and emax
	RptFmtEdev string    // format for edev
	RptFmtL    string    // format for lmin, lave and lmax
	RptFmtLdev string    // format for ldev
	RptFmtX    string    // format for x values
	HistFmt    string    // format in histogram
	HistNdig   int       // number of digits in histogram
	HistNsta   int       // number of stations in histogram
	HistLen    int       // number of characters (bar length) in histogram

	// RunMany: best solutions
	BestOvas      [][]float64 // best OVAs [nova][nsamples]
	BestFlts      [][]float64 // best flts [nflt][nsamples]
	BestInts      [][]int     // best ints [nint][nsamples]
	BestOfBestOva []float64   // [nova]
	BestOfBestFlt []float64   // [nflt]
	BestOfBestInt []int       // [nint]

	// RunMany: checking multi-obj problems
	F1F0_func      func(f0 float64) float64  // f1(f0) function
	F1F0_err       []float64                 // max(error(f1))
	F1F0_arcLen    []float64                 // arc-length: spreading on (f0,f1) space
	F1F0_arcLenRef float64                   // reference arc-length along f1(f0) curve
	F1F0_f0ranges  [][]float64               // ranges of f0 values to compute arc-length
	Multi_fcnErr   func(f []float64) float64 // computes Pareto-optimal front error with many OVAs
	Multi_err      []float64                 // max(error(f[i]))
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

	// allocate variables
	o.BestOvas = make([][]float64, o.Nova)
	o.BestFlts = make([][]float64, o.Nflt)
	o.BestInts = make([][]int, o.Nint)
	o.BestOfBestOva = make([]float64, o.Nova)
	o.BestOfBestFlt = make([]float64, o.Nflt)
	o.BestOfBestInt = make([]int, o.Nint)

	// perform trials
	for itrial := 0; itrial < o.Nsamples; itrial++ {

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
		if o.Nova > 1 { // multi-objective
			SortByFrontThenOva(o.Solutions, 0)
		} else { // single-objective
			SortByOva(o.Solutions, 0)
		}

		// feasible solution
		if o.Solutions[0].Feasible() {

			// best solution
			best := o.Solutions[0]
			for i := 0; i < o.Nova; i++ {
				o.BestOvas[i] = append(o.BestOvas[i], best.Ova[i])
			}
			for i := 0; i < o.Nflt; i++ {
				o.BestFlts[i] = append(o.BestFlts[i], best.Flt[i])
			}
			for i := 0; i < o.Nint; i++ {
				o.BestInts[i] = append(o.BestInts[i], best.Int[i])
			}

			// best of all trials
			first_best := len(o.BestOvas[0]) == 1
			if first_best {
				copy(o.BestOfBestOva, best.Ova)
				copy(o.BestOfBestFlt, best.Flt)
				copy(o.BestOfBestInt, best.Int)
			} else {
				if best.Ova[0] < o.BestOfBestOva[0] {
					copy(o.BestOfBestOva, best.Ova)
					copy(o.BestOfBestFlt, best.Flt)
					copy(o.BestOfBestInt, best.Int)
				}
			}

			// check multi-objective results
			if o.F1F0_func != nil {
				var rms_err float64
				var nfeasible int
				for _, sol := range o.Solutions {
					if sol.Feasible() && sol.FrontId == 0 {
						f0, f1 := sol.Ova[0], sol.Ova[1]
						f1_cor := o.F1F0_func(f0)
						rms_err += math.Pow(f1-f1_cor, 2.0)
						nfeasible++
					}
				}
				if nfeasible > 0 {
					rms_err = math.Sqrt(rms_err / float64(nfeasible))
					o.F1F0_err = append(o.F1F0_err, rms_err)
				}
			}

			// arc-length along Pareto front
			if o.Nova == 2 {
				if best.Feasible() && best.FrontId == 0 && o.Solutions[1].FrontId == 0 {
					dist := 0.0
					for i := 1; i < o.Nsol; i++ {
						if o.Solutions[i].FrontId == 0 {
							F0, F1 := o.Solutions[i-1].Ova[0], o.Solutions[i-1].Ova[1]
							f0, f1 := o.Solutions[i].Ova[0], o.Solutions[i].Ova[1]
							if o.F1F0_f0ranges != nil {
								a := o.find_f0_spot(F0)
								b := o.find_f0_spot(f0)
								if a == -1 || b == -1 {
									continue
								}
								if a != b {
									//io.Pforan("\nF0=%g is in [%g,%g]\n", F0, o.F1F0_f0ranges[a][0], o.F1F0_f0ranges[a][1])
									//io.Pfpink("f0=%g is in [%g,%g]\n", f0, o.F1F0_f0ranges[b][0], o.F1F0_f0ranges[b][1])
									continue
								}
							}
							dist += math.Sqrt(math.Pow(f0-F0, 2.0) + math.Pow(f1-F1, 2.0))
						}
					}
					o.F1F0_arcLen = append(o.F1F0_arcLen, dist)
				}
			}

			// multiple OVAs
			if o.Nova > 1 && o.Multi_fcnErr != nil {
				var rms_err float64
				var nfeasible int
				for _, sol := range o.Solutions {
					if sol.Feasible() && sol.FrontId == 0 {
						f_err := o.Multi_fcnErr(sol.Ova)
						rms_err += f_err * f_err
						nfeasible++
					}
				}
				if nfeasible > 0 {
					rms_err = math.Sqrt(rms_err / float64(nfeasible))
					o.Multi_err = append(o.Multi_err, rms_err)
				}
			}

			// save final solutions
			if fnkey != "" {
				f0min := best.Ova[0]
				for _, sol := range o.Solutions {
					f0min = utl.Min(f0min, sol.Ova[0])
				}
				WriteAllValues(dirout, io.Sf("%s-%04d_f0min=%g", fnkey, itrial, f0min), o)
			}
		}
	}
}

// StatF computes statistical information corresponding to objective function idxF
func StatF(o *Optimiser, idxF int, verbose bool) (fmin, fave, fmax, fdev float64, F []float64) {
	nsamples := len(o.BestOvas[idxF])
	if nsamples == 0 {
		if verbose {
			io.Pfred("there are no samples for statistical analysis\n")
		}
		return
	}
	F = make([]float64, nsamples)
	if nsamples == 1 {
		F[0] = o.BestOvas[idxF][0]
		fmin, fave, fmax = F[0], F[0], F[0]
		return
	}
	for i, f := range o.BestOvas[idxF] {
		F[i] = f
	}
	fmin, fave, fmax, fdev = rnd.StatBasic(F, true)
	if verbose {
		str := "\n"
		if len(o.RptFref) == o.Nova {
			str = io.Sf(" (%g)\n", o.RptFref[idxF])
		}
		io.Pf("fmin = %g\n", fmin)
		io.Pf("fave = %g"+str, fave)
		io.Pf("fmax = %g\n", fmax)
		io.Pf("fdev = %g\n", fdev)
		o.fix_formatting_data()
		io.Pf(rnd.BuildTextHist(nice(fmin-0.05, o.HistNdig), nice(fmax+0.05, o.HistNdig),
			o.HistNsta, F, o.HistFmt, o.HistLen))
	}
	return
}

// StatF1F0 prints statistical analysis for two-objective problems
//  emin, eave, emax, edev -- errors on f1(f0)
//  lmin, lave, lmax, ldev -- arc-lengths along f1(f0) curve
func StatF1F0(o *Optimiser, verbose bool) (emin, eave, emax, edev float64, E []float64, lmin, lave, lmax, ldev float64, L []float64) {
	if len(o.F1F0_err) == 0 && len(o.F1F0_arcLen) == 0 {
		io.Pfred("there are no samples for statistical analysis\n")
		return
	}
	o.fix_formatting_data()
	if len(o.F1F0_err) > 2 {
		E = make([]float64, len(o.F1F0_err))
		copy(E, o.F1F0_err)
		emin, eave, emax, edev = rnd.StatBasic(E, true)
		if verbose {
			io.Pf("\nerror on Pareto front\n")
			io.Pf("emin = %g\n", emin)
			io.Pf("eave = %g\n", eave)
			io.Pf("emax = %g\n", emax)
			io.Pf("edev = %g\n", edev)
			io.Pf(rnd.BuildTextHist(nice(emin-0.05, o.HistNdig), nice(emax+0.05, o.HistNdig),
				o.HistNsta, E, o.HistFmt, o.HistLen))
		}
	}
	if len(o.F1F0_arcLen) > 2 {
		den := 1.0
		if o.F1F0_arcLenRef > 0 {
			den = o.F1F0_arcLenRef
		}
		L := make([]float64, len(o.F1F0_arcLen))
		for i, l := range o.F1F0_arcLen {
			L[i] = l / den
		}
		lmin, lave, lmax, ldev = rnd.StatBasic(L, true)
		if verbose {
			io.Pf("\nnormalised arc length along Pareto front (ref = %g)\n", o.F1F0_arcLenRef)
			io.Pf("lmin = %g\n", lmin)
			io.Pf("lave = %g\n", lave)
			io.Pf("lmax = %g\n", lmax)
			io.Pf("ldev = %g\n", ldev)
			io.Pf(rnd.BuildTextHist(nice(lmin-0.05, o.HistNdig), nice(lmax+0.05, o.HistNdig),
				o.HistNsta, L, o.HistFmt, o.HistLen))
		}
	}
	return
}

// StatMulti prints statistical analysis for multi-objective problems
//  emin, eave, emax, edev -- errors on f1(f0)
//  lmin, lave, lmax, ldev -- arc-lengths along f1(f0) curve
func StatMulti(o *Optimiser, verbose bool) (emin, eave, emax, edev float64, E []float64) {
	if len(o.Multi_err) < 2 {
		io.Pfred("there are no samples for statistical analysis\n")
		return
	}
	o.fix_formatting_data()
	E = make([]float64, len(o.Multi_err))
	copy(E, o.Multi_err)
	emin, eave, emax, edev = rnd.StatBasic(E, true)
	if verbose {
		io.Pf("\nerror on Pareto front (multi)\n")
		io.Pf("emin = %g\n", emin)
		io.Pf("eave = %g\n", eave)
		io.Pf("emax = %g\n", emax)
		io.Pf("edev = %g\n", edev)
		io.Pf(rnd.BuildTextHist(nice(emin-0.05, o.HistNdig), nice(emax+0.05, o.HistNdig),
			o.HistNsta, E, o.HistFmt, o.HistLen))
	}
	return
}

// fix_formatting_data fixes formatting data and data for histograms
func (o *Stat) fix_formatting_data() {
	if o.RptFmtF == "" {
		o.RptFmtF = "%g"
	}
	if o.RptFmtFdev == "" {
		o.RptFmtFdev = "%g"
	}
	if o.RptFmtE == "" {
		o.RptFmtE = "%.8e"
	}
	if o.RptFmtEdev == "" {
		o.RptFmtEdev = "%.8e"
	}
	if o.RptFmtL == "" {
		o.RptFmtL = "%g"
	}
	if o.RptFmtLdev == "" {
		o.RptFmtLdev = "%.8e"
	}
	if o.RptFmtX == "" {
		o.RptFmtX = "%g"
	}
	if o.HistFmt == "" {
		o.HistFmt = "%.2f"
	}
	if o.HistNdig == 0 {
		o.HistNdig = 3
	}
	if o.HistNsta == 0 {
		o.HistNsta = 8
	}
	if o.HistLen == 0 {
		o.HistLen = 20
	}
}

// find_f0_spot finds where f0 falls in
func (o *Stat) find_f0_spot(f0 float64) (idx int) {
	for i, f0vals := range o.F1F0_f0ranges {
		if f0 >= f0vals[0] && f0 <= f0vals[1] {
			return i
		}
	}
	l := len(o.F1F0_f0ranges) - 1
	if f0 > o.F1F0_f0ranges[l][0] {
		return l
	}
	if f0 < o.F1F0_f0ranges[0][1] {
		return 0
	}
	return -1
}

// nice returns a truncated float
func nice(x float64, ndigits int) float64 {
	s := io.Sf("%."+io.Sf("%d", ndigits)+"f", x)
	return io.Atof(s)
}

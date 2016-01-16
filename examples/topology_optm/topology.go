// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

// main function
func main() {

	// GA parameters
	var opt goga.Optimiser
	opt.Read("ga-data.json")

	// FEM
	data := make([]*FemData, opt.Ncpu)
	for i := 0; i < opt.Ncpu; i++ {
		data[i] = NewData(i)
	}

	// set integers
	if data[0].Opt.BinInt {
		opt.CxInt = goga.CxInt
		opt.MtInt = goga.MtIntBin
		opt.BinInt = data[0].Ncells
	}

	// set floats
	opt.FltMin = make([]float64, data[0].Nareas)
	opt.FltMax = make([]float64, data[0].Nareas)
	for i := 0; i < data[0].Nareas; i++ {
		opt.FltMin[i] = data[0].Opt.Amin
		opt.FltMax[i] = data[0].Opt.Amax
	}

	// initialise optimiser
	opt.Nova = 2 // weight and deflection
	opt.Noor = 4 // mobility, feasibility, maxdeflection, stress
	opt.Init(goga.GenTrialSolutions, func(sol *goga.Solution, cpu int) {
		mob, fail, weight, umax, _, errU, errS := data[cpu].RunFEM(sol.Int, sol.Flt, false)
		sol.Ova[0] = weight
		sol.Ova[1] = umax
		sol.Oor[0] = mob
		sol.Oor[1] = fail
		sol.Oor[2] = errU
		sol.Oor[3] = errS
	}, nil, 0, 0, 0)

	// initial solutions
	//io.Pforan("%s\n", PrintSolutions(data[0], opt.Solutions))
	io.Pforan("DtExc = %v\n", opt.DtExc)

	// initial solutions
	var sols0 []*goga.Solution
	if false {
		sols0 = opt.GetSolutionsCopy()
	}

	// solve
	opt.Verbose = true
	opt.Solve()

	// check
	var nfailed, nsuccess int
	for _, sol := range opt.Solutions {
		var failed bool
		for _, oor := range sol.Oor {
			if oor > 0 {
				failed = true
				break
			}
		}
		if failed {
			nfailed++
		} else {
			nsuccess++
		}
	}
	if nfailed > 0 {
		io.PfRed("N failed = %d out of %d\n", nfailed, opt.Nsol)
	} else {
		io.PfGreen("N success = %d out of %d\n", nsuccess, opt.Nsol)
	}

	// save results
	goga.SortByOva(opt.Solutions, 0)
	fnkey := data[0].Analysis.Sim.Key
	var log, res bytes.Buffer
	io.Ff(&log, opt.LogParams())
	io.Ff(&res, PrintSolutions(data[0], opt.Solutions))
	io.Ff(&res, io.Sf("\n\nnfailed = %d\n", nfailed))
	io.WriteFileVD("/tmp/goga", fnkey+".log", &log)
	io.WriteFileVD("/tmp/goga", fnkey+".res", &res)

	// plot
	feasibleOnly := true
	plt.SetForEps(0.8, 355)
	if strings.HasPrefix(fnkey, "truss10bar") {
		_, ref, _ := io.ReadTable("p460_fig300.dat")
		plt.Plot(ref["w"], ref["u"], "'b-'")
	}
	fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
	fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
	goga.PlotOvaOvaPareto(&opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
	plt.Gll("weight ($f_0$)", "deflection ($f_1)$", "leg_out=1, leg_ncol=4, leg_hlen=1.5")
	plt.SaveD("/tmp/goga", fnkey+".eps")
}

type FltFormatter []float64

func (o FltFormatter) String() (l string) {
	for _, val := range o {
		if val < 1e-9 {
			l += "       "
		} else {
			l += io.Sf("%7.2f", val)
		}
	}
	return l
}

func PrintSolutions(fed *FemData, sols []*goga.Solution) (l string) {
	goga.SortByOva(sols, 0)
	l = io.Sf("%8s%6s%6s |%s\n", "weight", "umax", "smax", "areas")
	for _, sol := range sols {
		mob, fail, weight, umax, smax, errU, errS := fed.RunFEM(sol.Int, sol.Flt, false)
		if mob > 0 || fail > 0 || errU > 0 || errS > 0 {
			l += io.Sf("%20s |%s\n", "unfeasible    ", FltFormatter(sol.Flt))
			continue
		}
		l += io.Sf("%8.1f%6.2f%6.2f |%s\n", weight, umax, smax, FltFormatter(sol.Flt))
	}
	return
}

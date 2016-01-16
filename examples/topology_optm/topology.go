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

	// input filename
	fn, fnkey := io.ArgToFilename(0, "ground10", ".sim", true)

	// GA parameters
	var opt goga.Optimiser
	opt.Read("ga-" + fnkey + ".json")

	// FEM
	data := make([]*FemData, opt.Ncpu)
	for i := 0; i < opt.Ncpu; i++ {
		data[i] = NewData(fn, fnkey, i)
	}
	io.Pforan("MaxWeight = %v\n", data[0].MaxWeight)

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
		mob, fail, weight, umax, _, errU, errS := data[cpu].RunFEM(sol.Int, sol.Flt, 0, false)
		sol.Ova[0] = weight
		sol.Ova[1] = umax
		sol.Oor[0] = mob
		sol.Oor[1] = fail
		sol.Oor[2] = errU
		sol.Oor[3] = errS
	}, nil, 0, 0, 0)

	// initial solutions
	var sols0 []*goga.Solution
	if false {
		sols0 = opt.GetSolutionsCopy()
	}

	// solve
	opt.Verbose = true
	opt.Solve()
	goga.SortByOva(opt.Solutions, 0)

	// check
	front0 := make([]*goga.Solution, 0)
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
			if sol.FrontId == 0 {
				front0 = append(front0, sol)
			}
		}
	}
	if nfailed > 0 {
		io.PfRed("N failed = %d out of %d\n", nfailed, opt.Nsol)
	} else {
		io.PfGreen("N success = %d out of %d\n", nsuccess, opt.Nsol)
		io.PfGreen("N front 0 = %d\n", len(front0))
	}

	// save results
	var log, res bytes.Buffer
	io.Ff(&log, opt.LogParams())
	io.Ff(&res, PrintSolutions(data[0], opt.Solutions))
	io.Ff(&res, io.Sf("\n\nnfailed = %d\n", nfailed))
	io.WriteFileVD("/tmp/goga", fnkey+".log", &log)
	io.WriteFileVD("/tmp/goga", fnkey+".res", &res)

	// plot Pareto-optimal front
	feasibleOnly := true
	plt.SetForEps(0.8, 355)
	if strings.HasPrefix(fnkey, "ground10") {
		_, ref, _ := io.ReadTable("p460_fig300.dat")
		plt.Plot(ref["w"], ref["u"], "'b-'")
	}
	fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
	fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
	goga.PlotOvaOvaPareto(&opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
	plt.Gll("weight ($f_0$)", "deflection ($f_1)$", "leg_out=1, leg_ncol=4, leg_hlen=1.5")
	if strings.HasPrefix(fnkey, "ground10") {
		plt.AxisRange(1800, 14000, 1, 6)
	}

	// plot selected results
	nfront0 := len(front0)
	if nfront0 > 2 {
		m := nfront0 / 2
		l := nfront0 - 1
		io.Pforan("nfront0=%d m=%d l=%v\n", nfront0, m, l)
		_, _, weight, umax, _, _, _ := data[0].RunFEM(front0[0].Int, front0[0].Flt, 0, false)
		plt.Text(weight, umax, "1", "size=7")
		plt.PlotOne(weight, umax, "'g*', zorder=100")
		_, _, weight, umax, _, _, _ = data[0].RunFEM(front0[m].Int, front0[m].Flt, 0, false)
		plt.Text(weight, umax, "2", "size=7")
		plt.PlotOne(weight, umax, "'g*', zorder=100")
		_, _, weight, umax, _, _, _ = data[0].RunFEM(front0[l].Int, front0[l].Flt, 0, false)
		plt.Text(weight, umax, "3", "size=7")
		plt.PlotOne(weight, umax, "'g*', zorder=100")
		plt.PyCmds(`
from pylab import axes, setp
a = axes([0.2, 0.75, 0.20, 0.10], axisbg='#dcdcdc')
setp(a, xticks=[0,720], yticks=[0,360])
axis('equal')
axis('off')
`)
		data[0].RunFEM(front0[0].Int, front0[0].Flt, 1, false)
		plt.PyCmds(`
a = axes([0.40, 0.28, 0.20, 0.10], axisbg='#dcdcdc')
setp(a, xticks=[0,720], yticks=[0,360])
axis('equal')
axis('off')
`)
		data[0].RunFEM(front0[m].Int, front0[m].Flt, 2, false)
		plt.PyCmds(`
a = axes([0.7, 0.18, 0.20, 0.10], axisbg='#dcdcdc')
setp(a, xticks=[0,720], yticks=[0,360])
axis('equal')
axis('off')
`)
		data[0].RunFEM(front0[l].Int, front0[l].Flt, 3, false)
	}

	// save
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
		mob, fail, weight, umax, smax, errU, errS := fed.RunFEM(sol.Int, sol.Flt, 0, false)
		if mob > 0 || fail > 0 || errU > 0 || errS > 0 {
			l += io.Sf("%20s |%s\n", "unfeasible    ", FltFormatter(sol.Flt))
			continue
		}
		l += io.Sf("%8.1f%6.2f%6.2f |%s\n", weight, umax, smax, FltFormatter(sol.Flt))
	}
	return
}

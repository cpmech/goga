// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"strings"
	"time"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

// main function
func main() {

	// flags
	benchmark := false
	ncpuMax := 16

	// benchmarking
	if benchmark {
		var nsol, tmax int
		var et time.Duration
		X := make([]float64, ncpuMax)
		T := make([]float64, ncpuMax)
		S := make([]float64, ncpuMax) // speedup
		S[0] = 1
		for i := 0; i < ncpuMax; i++ {
			io.Pf("\n\n")
			nsol, tmax, et = run(i + 1)
			io.PfYel("elaspsedTime = %v\n", et)
			X[i] = float64(i + 1)
			T[i] = et.Seconds()
			if i > 0 {
				S[i] = T[0] / T[i] // Told / Tnew
			}
		}

		plt.Reset(true, nil)
		plt.Plot(X, S, &plt.A{C: "b", M: ".", L: io.Sf("speedup: $N_{sol}=%d,\\,t_{max}=%d$", nsol, tmax)})
		plt.Plot([]float64{1, 16}, []float64{1, 16}, &plt.A{C: "k"})
		plt.Gll("$N_{cpu}:\\;$ number of groups", "speedup", &plt.A{LegOut: true})
		plt.DoubleYscale("$T_{sys}:\\;$ system time [s]")
		plt.Plot(X, T, &plt.A{C: "gray"})
		plt.Save("/tmp/goga", "topology-speedup")
		return
	}

	// normal run
	run(-1)
}

// run runs optimiser
func run(ncpu int) (nsol, tmax int, elaspsedTime time.Duration) {

	// options
	doPlot := true
	doReport := true

	// input filename
	fn, fnkey := io.ArgToFilename(0, "ground10", ".sim", true)

	// parameters
	var opt goga.Optimiser
	opt.Read("ga-" + fnkey + ".json")
	opt.GenType = "rnd"
	nsol, tmax = opt.Nsol, opt.Tmax
	postproc := true
	if ncpu > 0 {
		opt.Ncpu = ncpu
		postproc = false
	}

	// FEM
	data := make([]*FemData, opt.Ncpu)
	for i := 0; i < opt.Ncpu; i++ {
		data[i] = NewData(fn, fnkey, i)
	}
	io.Pf("MaxWeight = %v\n", data[0].MaxWeight)

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

	// benchmark
	initialTime := time.Now()
	defer func() {
		elaspsedTime = time.Now().Sub(initialTime)
	}()

	// solve
	opt.Verbose = true
	opt.Solve()
	goga.SortSolutions(opt.Solutions, 0)

	// post processing
	if !postproc {
		return
	}

	// check
	nfailed, front0 := goga.CheckFront0(&opt, true)

	// save results
	var log, res bytes.Buffer
	io.Ff(&log, opt.LogParams())
	io.Ff(&res, PrintSolutions(data[0], opt.Solutions))
	io.Ff(&res, io.Sf("\n\nnfailed = %d\n", nfailed))
	io.WriteFileVD("/tmp/goga", fnkey+".log", &log)
	io.WriteFileVD("/tmp/goga", fnkey+".res", &res)

	// selected results
	ia, ib, ic, id, ie := 0, 0, 0, 0, 0
	nfront0 := len(front0)
	if nfront0 > 4 {
		ib = nfront0 / 10
		ic = nfront0 / 5
		id = nfront0 / 2
		ie = nfront0 - 1
	}
	A := front0[ia]
	B := front0[ib]
	C := front0[ic]
	D := front0[id]
	E := front0[ie]

	// plot Pareto-optimal front
	if doPlot {
		var ref map[string][]float64
		if strings.HasPrefix(fnkey, "ground10") {
			_, ref, _ = io.ReadTable("p460_fig300.dat")
		}
		plt.Reset(true, &plt.A{Eps: true})
		pp := goga.NewPlotParams(false)
		pp.FnKey = fnkey
		pp.FeasibleOnly = true
		pp.ArgsFront = &plt.A{C: "r", M: ".", Ls: "none", Ms: 6, Mec: "black", Mew: 0.3, L: "Goga best front"}
		pp.Xlabel = "weight ($f_0$)"
		pp.Ylabel = "deflection ($f_1)$"
		pp.Extra = func() {
			plt.Plot(ref["w"], ref["u"], &plt.A{C: "k", L: "reference"})
			wid, hei := 0.20, 0.10
			drawTruss(data[0], "A", A, 0.17, 0.75, wid, hei)
			drawTruss(data[0], "B", B, 0.20, 0.55, wid, hei)
			drawTruss(data[0], "C", C, 0.28, 0.33, wid, hei)
			drawTruss(data[0], "D", D, 0.47, 0.22, wid, hei)
			drawTruss(data[0], "E", E, 0.70, 0.18, wid, hei)
			if strings.HasPrefix(fnkey, "ground10") {
				plt.AxisRange(1800, 14000, 1, 6)
			}
			plt.HideBorders(&plt.A{HideT: true, HideR: true})
			io.Pf("truss drawn\n")
		}
		opt.PlotOvaOvaPareto(sols0, 0, 1, pp)
	}

	// Report
	if doReport {
		title := "Shape and topology optimisation. Results."
		label := "topoFront"
		document := true
		compact := true
		texResults("/tmp/goga", "tmp_"+fnkey, title, label, data[0], A, B, C, D, E, document, compact)
		document = false
		texResults("/tmp/goga", fnkey, title, label, data[0], A, B, C, D, E, document, compact)
	}
	return
}

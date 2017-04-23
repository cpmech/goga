// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// options
	doPlot := true
	doReport := true
	constantSeed := false

	// goga parameters
	opt := new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 200
	opt.Ncpu = 4
	opt.Tmax = 500
	//opt.Tmax = 100
	opt.Verbose = true
	opt.VerbTime = true
	opt.VerbStat = true
	opt.GenType = "latin"
	opt.Nsamples = 1
	opt.EpsH = 1e-3

	// flags
	problem := 4
	checkOnly := false
	lossless := problem < 4

	// generators
	var sys System
	sys.Init(2.834, lossless, checkOnly)
	ngs := len(sys.G) // number of generators == number of variables (Nx)
	opt.FltMin = make([]float64, ngs)
	opt.FltMax = make([]float64, ngs)
	for i := 0; i < ngs; i++ {
		opt.FltMin[i], opt.FltMax[i] = sys.G[i].Pmin, sys.G[i].Pmax
	}

	// problem variables
	var nf, ng, nh int // number of objective functions and constraints
	var fcn goga.MinProb_t

	// problems
	switch problem {

	// problem # 1: lossless and unsecured: cost only
	case 1:
		opt.RptXref = []float64{0.10954, 0.29967, 0.52447, 1.01601, 0.52469, 0.35963}
		opt.RptFref = []float64{600.114}
		nf, nh = 1, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = sys.FuelCost(x)
			h[0] = sys.Balance(x)
		}

	// lossless and unsecured: emission only
	case 2:
		opt.RptXref = []float64{0.40584, 0.45915, 0.53797, 0.38300, 0.53791, 0.51012}
		opt.RptFref = []float64{0.19420}
		nf, nh = 1, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = sys.Emission(x)
			h[0] = sys.Balance(x)
		}

	// lossless and unsecured: cost and emission
	case 3:
		nf, nh = 2, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = sys.FuelCost(x)
			f[1] = sys.Emission(x)
			h[0] = sys.Balance(x)
		}
		checkOnly = false

	// with loss but unsecured: cost and emission
	default:
		nf, nh = 2, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = sys.FuelCost(x)
			f[1] = sys.Emission(x)
			h[0] = sys.Balance(x)
		}
		checkOnly = false
	}

	// check best known solution
	if checkOnly {
		check(fcn, ng, nh, opt.RptXref, opt.RptFref[0], 1e-6)
		return
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// initial solutions
	var sols0 []*goga.Solution
	if false {
		sols0 = opt.GetSolutionsCopy()
	}

	// solve
	opt.RunMany("", "", constantSeed)
	goga.SortSolutions(opt.Solutions, 0)

	// post-processing
	fnkey := io.Sf("eed%d", problem)
	switch problem {
	case 1:
		// selected results
		idxsel := 120
		m, l := idxsel, opt.Nsol-1
		A, B, C := opt.Solutions[0], opt.Solutions[m], opt.Solutions[l]

		// results
		PrintResults(&sys, []*goga.Solution{A, B, C}, opt.RptXref, opt.RptFref[0], 0.22214)

		// stat
		io.Pf("\n")
		opt.PrintStatF(0)

	case 2:
		// selected results
		idxsel := 120
		m, l := idxsel, opt.Nsol-1
		A, B, C := opt.Solutions[0], opt.Solutions[m], opt.Solutions[l]

		// results
		PrintResults(&sys, []*goga.Solution{A, B, C}, opt.RptXref, 638.260, opt.RptFref[0])

		// stat
		io.Pf("\n")
		opt.PrintStatF(0)

	default:

		// check
		nfailed, front0 := goga.CheckFront0(opt, true)

		// selected results
		ia, ib, ic, id, ie, iF := 0, 0, 0, 0, 0, 0
		nfront0 := len(front0)
		if nfront0 > 4 {
			ib = int(float64(nfront0) * 0.3)
			ic = int(float64(nfront0) * 0.5)
			id = int(float64(nfront0) * 0.7)
			ie = int(float64(nfront0) * 0.9)
			iF = nfront0 - 1
		}
		A := front0[ia]
		B := front0[ib]
		C := front0[ic]
		D := front0[id]
		E := front0[ie]
		F := front0[iF]

		// print results
		selected := []*goga.Solution{A, B, C, D, E, F}
		strRes := PrintResults(&sys, selected, nil, -1, -1)

		// save results
		var log, res bytes.Buffer
		io.Ff(&log, opt.LogParams())
		io.Ff(&res, strRes)
		io.Ff(&res, io.Sf("\n\nnfailed = %d\n", nfailed))
		io.WriteFileVD("/tmp/goga", fnkey+".log", &log)
		io.WriteFileVD("/tmp/goga", fnkey+".res", &res)

		// load reference data
		title := "Economic emission dispatch. Results: Lossy case."
		label := "eedLossy"
		var dat map[string][]float64
		if problem == 3 {
			title = "Economic emission dispatch. Results: Lossless case."
			label = "eedLossless"
			_, dat, _ = io.ReadTable("abido2006-fig6.dat")
		} else {
			_, dat, _ = io.ReadTable("abido2006-fig8.dat")
		}

		// plot
		if doPlot {
			argsSel := &plt.A{C: "b", M: "*"}
			argsTxt := &plt.A{C: "b"}
			addSelected := func(key string, P *goga.Solution) {
				plt.PlotOne(P.Ova[0], P.Ova[1], argsSel)
				plt.Text(P.Ova[0]+1, P.Ova[1], key, argsTxt)
			}
			plt.Reset(true, &plt.A{Eps: true, Prop: 0.75, WidthPt: 300})
			pp := goga.NewPlotParams(false)
			pp.FnKey = fnkey
			pp.Xlabel = "$f_0:\\;$ cost~[\\$/h]"
			pp.Ylabel = "$f_1:\\quad$ emission~[ton/h]"
			pp.ArgsLeg = &plt.A{LegOut: false}
			pp.Extra = func() {
				plt.Plot(dat["cost"], dat["emission"], &plt.A{C: "#2c7248", Ls: "-", L: "reference"})
				for i, selected := range selected {
					addSelected(selectedKeys[i], selected)
				}
				plt.HideBorders(&plt.A{HideR: true, HideT: true})
				if problem == 3 {
					plt.AxisXmin(599)
				}
			}
			opt.PlotOvaOvaPareto(sols0, 0, 1, pp)
		}

		// tex file
		if doReport {
			document := true
			compact := true
			texResults("/tmp/goga", "tmp_"+fnkey, title, label, &sys, selected, document, compact)
			document = false
			texResults("/tmp/goga", fnkey, title, label, &sys, selected, document, compact)
		}
	}
}

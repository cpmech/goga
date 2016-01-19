// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// goga parameters
	opt := new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 200
	opt.Ncpu = 4
	opt.Tf = 500
	opt.Verbose = false
	opt.Ntrials = 1
	opt.GenType = "latin"
	opt.EpsH = 1e-3

	// flags
	problem := 4
	checkOnly := false

	// generators
	Pdemand := 2.834
	Ploss := 0.0
	gs, B00, B0, B := NewGenerators(false)
	ngs := len(gs) // number of generators == number of variables (Nx)
	opt.FltMin = make([]float64, ngs)
	opt.FltMax = make([]float64, ngs)
	for i := 0; i < ngs; i++ {
		opt.FltMin[i], opt.FltMax[i] = gs[i].Pmin, gs[i].Pmax
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
			sumP := la.VecAccum(x)
			f[0] = gs.FuelCost(x)
			h[0] = sumP - Pdemand - Ploss
		}

	// lossless and unsecured: emission only
	case 2:
		opt.RptXref = []float64{0.40584, 0.45915, 0.53797, 0.38300, 0.53791, 0.51012}
		opt.RptFref = []float64{0.19420}
		nf, nh = 1, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			sumP := la.VecAccum(x)
			f[0] = gs.Emission(x)
			h[0] = sumP - Pdemand - Ploss
		}

	// lossless and unsecured: cost and emission
	case 3:
		nf, nh = 2, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			sumP := la.VecAccum(x)
			f[0] = gs.FuelCost(x)
			f[1] = gs.Emission(x)
			h[0] = sumP - Pdemand - Ploss
		}

	// with loss but unsecured: cost and emission
	default:
		nf, nh = 2, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			sumP := la.VecAccum(x)
			f[0] = gs.FuelCost(x)
			f[1] = gs.Emission(x)
			Ploss = B00
			for i := 0; i < ngs; i++ {
				Ploss += B0[i] * x[i]
				for j := 0; j < ngs; j++ {
					Ploss += x[i] * B[i][j] * x[j]
				}
			}
			h[0] = sumP - Pdemand - Ploss
		}
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
	opt.RunMany("", "")

	// post-processing
	switch problem {
	case 1:

		// results
		print_single(opt, gs, opt.RptXref, opt.RptFref[0], 0.22214, Pdemand, Ploss)

		// stat
		io.Pf("\n")
		goga.StatF(opt, 0, true)

	case 2:

		// results
		print_single(opt, gs, opt.RptXref, 638.260, opt.RptFref[0], Pdemand, Ploss)

		// stat
		io.Pf("\n")
		goga.StatF(opt, 0, true)

	default:

		// stat
		goga.StatF1F0(opt, true)

		// plot
		feasibleOnly := true
		plt.SetForEps(0.8, 300)
		fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
		fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
		goga.PlotOvaOvaPareto(opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
		var dat map[string][]float64
		if problem == 3 {
			_, dat, _ = io.ReadTable("abido2006-fig6.dat")
		} else {
			_, dat, _ = io.ReadTable("abido2006-fig8.dat")
		}
		plt.Plot(dat["cost"], dat["emission"], "'b-', label='reference'")
		plt.Gll("$f_0:\\;$ cost~[\\$/h]", "$f_1:\\quad$ emission~[ton/h]", "")
		plt.SaveD("/tmp/goga", io.Sf("ecoemission_prob%d.eps", problem))

		// report
		goga.TexF1F0Report("/tmp/goga", "tmp_ecoemission", "ecoemission", 10, true, []*goga.Optimiser{opt})
		goga.TexF1F0Report("/tmp/goga", "ecoemission", "ecoemission", 10, false, []*goga.Optimiser{opt})
	}
}

func print_single(opt *goga.Optimiser, gs Generators, Pref []float64, costRef, emisRef, Pdemand, Ploss float64) {
	n := 9*2 + 8*6 + 12
	P := Pref
	sumP := la.VecAccum(P)
	bal_err := math.Abs(sumP - Pdemand - Ploss)
	io.Pf("%s", io.StrThickLine(n))
	io.Pf("%9s%9s%8s%8s%8s%8s%8s%8s%12s\n", "cost", "emis", "P1", "P2", "P3", "P4", "P5", "P6", "bal.err")
	io.Pf("%s", io.StrThinLine(n))
	io.Pfyel("%9.4f%9.5f%8.5f%8.5f%8.5f%8.5f%8.5f%8.5f%12.4e\n", costRef, emisRef, P[0], P[1], P[2], P[3], P[4], P[5], bal_err)
	goga.SortByOva(opt.Solutions, 0)
	for i, sol := range opt.Solutions {
		P = sol.Flt
		sumP = la.VecAccum(P)
		cost := gs.FuelCost(P)
		emis := gs.Emission(P)
		bal_err = math.Abs(sumP - Pdemand - Ploss)
		io.Pf("%9.4f%9.5f%8.5f%8.5f%8.5f%8.5f%8.5f%8.5f%12.4e\n", cost, emis, P[0], P[1], P[2], P[3], P[4], P[5], bal_err)
		if i > 20 {
			break
		}
	}
	io.Pf("%s", io.StrThickLine(n))
}

func check(fcn goga.MinProb_t, ng, nh int, xs []float64, fs, ϵ float64) {
	f := make([]float64, 1)
	g := make([]float64, ng)
	h := make([]float64, nh)
	cpu := 0
	fcn(f, g, h, xs, nil, cpu)
	io.Pfblue2("xs = %v\n", xs)
	io.Pfblue2("f(x)=%g  (%g)  diff=%g\n", f[0], fs, math.Abs(fs-f[0]))
	for i, v := range g {
		unfeasible := false
		if v < 0 {
			unfeasible = true
		}
		if unfeasible {
			io.Pfred("g%d(x) = %g\n", i, v)
		} else {
			io.Pfgreen("g%d(x) = %g\n", i, v)
		}
	}
	for i, v := range h {
		unfeasible := false
		if math.Abs(v) > ϵ {
			unfeasible = true
		}
		if unfeasible {
			io.Pfred("h%d(x) = %g\n", i, v)
		} else {
			io.Pfgreen("h%d(x) = %g\n", i, v)
		}
	}
}

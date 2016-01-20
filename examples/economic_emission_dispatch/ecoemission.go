// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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
	problem := 3
	idxsel := 120
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
	opt.RunMany("", "")

	// selected solutions
	goga.SortByOva(opt.Solutions, 0)
	m, l := idxsel, opt.Nsol-1
	A, B, C := opt.Solutions[0], opt.Solutions[m], opt.Solutions[l]

	// post-processing
	switch problem {
	case 1:

		// results
		print_results(&sys, A, B, C, opt.RptXref, opt.RptFref[0], 0.22214)

		// stat
		io.Pf("\n")
		goga.StatF(opt, 0, true)

	case 2:

		// results
		print_results(&sys, A, B, C, opt.RptXref, 638.260, opt.RptFref[0])

		// stat
		io.Pf("\n")
		goga.StatF(opt, 0, true)

	default:

		// results
		print_results(&sys, A, B, C, nil, -1, -1)

		// stat
		goga.StatF1F0(opt, true)

		// plot
		feasibleOnly := true
		plt.SetForEps(0.8, 300)
		fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
		fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
		goga.PlotOvaOvaPareto(opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
		plt.PlotOne(A.Ova[0], A.Ova[1], "'g*', zorder=1000, clip_on=0")
		plt.PlotOne(B.Ova[0], B.Ova[1], "'g*', zorder=1000, clip_on=0")
		plt.PlotOne(C.Ova[0], C.Ova[1], "'g*', zorder=1000, clip_on=0")
		plt.Text(A.Ova[0]+1, A.Ova[1], "A", "")
		plt.Text(B.Ova[0]+1, B.Ova[1], "B", "")
		plt.Text(C.Ova[0]+1, C.Ova[1], "C", "")
		var dat map[string][]float64
		if problem == 3 {
			_, dat, _ = io.ReadTable("abido2006-fig6.dat")
		} else {
			_, dat, _ = io.ReadTable("abido2006-fig8.dat")
		}
		plt.Plot(dat["cost"], dat["emission"], "'b-', label='reference'")
		plt.Gll("$f_0:\\;$ cost~[\\$/h]", "$f_1:\\quad$ emission~[ton/h]", "")
		if problem == 3 {
			plt.AxisXmin(599)
		}
		fnkey := io.Sf("ecoemission_prob%d", problem)
		plt.SaveD("/tmp/goga", fnkey+".eps")

		// report
		goga.TexF1F0Report("/tmp/goga", "tmp_"+fnkey, fnkey, 10, true, []*goga.Optimiser{opt})
		goga.TexF1F0Report("/tmp/goga", fnkey, fnkey, 10, false, []*goga.Optimiser{opt})
		tex_results("/tmp/goga", "res_"+fnkey, &sys, A, B, C, true)
	}
}

func print_results(sys *System, A, B, C *goga.Solution, Pref []float64, costRef, emisRef float64) {
	n := 9*2 + 8*6 + 12 + 3
	io.Pf("%s", io.StrThickLine(n))
	io.Pf("%3s%9s%9s%8s%8s%8s%8s%8s%8s%12s\n", "pt", "cost", "emis", "P1", "P2", "P3", "P4", "P5", "P6", "bal.err")
	io.Pf("%s", io.StrThinLine(n))
	if Pref != nil {
		io.Pfyel("%3s%9.4f%9.5f%8.5f%8.5f%8.5f%8.5f%8.5f%8.5f%12.4e\n", "ref", costRef, emisRef, Pref[0], Pref[1], Pref[2], Pref[3], Pref[4], Pref[5], sys.Balance(Pref))
	}
	writeline := func(pt string, P []float64) {
		io.Pf("%3s%9.4f%9.5f%8.5f%8.5f%8.5f%8.5f%8.5f%8.5f%12.4e\n", pt, sys.FuelCost(P), sys.Emission(P), P[0], P[1], P[2], P[3], P[4], P[5], sys.Balance(P))
	}
	writeline("A", A.Flt)
	writeline("B", B.Flt)
	writeline("C", C.Flt)
	io.Pf("%s", io.StrThickLine(n))
}

func tex_results(dirout, fnkey string, sys *System, A, B, C *goga.Solution, dorun bool) {
	buf := new(bytes.Buffer)
	io.Ff(buf, `\documentclass[a4paper]{article}

\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{booktabs}
\usepackage[margin=1.5cm,footskip=0.5cm]{geometry}

\title{GOGA Report}
\author{Dorival Pedroso}

\begin{document}

\begin{table} \centering
\caption{goga: Parameters}
\begin{tabular}[c]{cccccccccc} \toprule
point & cost & emission & $P_0$ & $P_1$ & $P_2$ & $P_3$ & $P_4$ & $P_5$ & $h_0$ \\ \hline
`)

	writeline := func(pt string, P []float64) {
		io.Ff(buf, "%s & $%.4f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.4e$ \\\\\n", pt, sys.FuelCost(P), sys.Emission(P), P[0], P[1], P[2], P[3], P[4], P[5], sys.Balance(P))
	}

	writeline("A", A.Flt)
	writeline("B", B.Flt)
	writeline("C", C.Flt)

	io.Ff(buf, `
\bottomrule
\end{tabular}
\label{tab:ecoemission}
\end{table}
\end{document}`)

	tex := fnkey + ".tex"
	io.WriteFileVD(dirout, tex, buf)
	if dorun {
		_, err := io.RunCmd(true, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory=/tmp/goga/", tex)
		if err != nil {
			chk.Panic("%v", err)
		}
		io.PfBlue("file <%s/%s.pdf> generated\n", dirout, fnkey)
	}
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

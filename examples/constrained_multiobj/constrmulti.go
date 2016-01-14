// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	PI = math.Pi
)

func CTPgenerator(θ, a, b, c, d, e float64) goga.MinProb_t {
	return func(f, g, h, x []float64, ξ []int, cpu int) {
		sin, cos := math.Sin(θ), math.Cos(θ)
		c0 := 1.0
		for i := 1; i < len(x); i++ {
			c0 += x[i]
		}
		f[0] = x[0]
		f[1] = c0 * (1.0 - f[0]/c0)
		if true {
			c1 := cos*(f[1]-e) - sin*f[0]
			c2 := sin*(f[1]-e) + cos*f[0]
			c3 := math.Sin(b * PI * math.Pow(c2, c))
			g[0] = c1 - a*math.Pow(math.Abs(c3), d)
		}
	}
}

func CTPplotter(θ, a, b, c, d, e, f1max float64) func() {
	return func() {
		np := 201
		X, Y := utl.MeshGrid2D(0, 1, 0, f1max, np, np)
		Z1 := utl.DblsAlloc(np, np)
		Z2 := utl.DblsAlloc(np, np)
		sin, cos := math.Sin(θ), math.Cos(θ)
		for j := 0; j < np; j++ {
			for i := 0; i < np; i++ {
				c1 := cos*(Y[i][j]-e) - sin*X[i][j]
				c2 := sin*(Y[i][j]-e) + cos*X[i][j]
				c3 := math.Sin(b * PI * math.Pow(c2, c))
				Z1[i][j] = c1
				Z2[i][j] = c1 - a*math.Pow(math.Abs(c3), d)
			}
		}
		plt.ContourSimple(X, Y, Z1, false, 7, "linestyles=['--'], linewidths=[0.7], colors=['g'], levels=[0]")
		plt.ContourSimple(X, Y, Z2, false, 7, "linestyles=['-'],  linewidths=[0.7], colors=['b'], levels=[0]")
	}
}

func solve_problem(problem, ntrials int, doplot bool) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Ncpu = 1
	opt.Tf = 500
	opt.Verbose = false
	opt.Ntrials = ntrials
	opt.GenType = "latin"

	// options for report
	opt.RptFmtE = "%.5e"
	opt.RptFmtL = "%.6f"
	opt.RptFmtEdev = "%.4e"
	opt.RptFmtLdev = "%.4e"

	// problem variables
	nx := 10
	opt.RptName = io.Sf("CTP%d", problem)
	opt.Nsol = 10 * nx
	opt.Ncpu = 1
	opt.FltMin = make([]float64, nx)
	opt.FltMax = make([]float64, nx)
	for i := 0; i < nx; i++ {
		opt.FltMin[i] = 0
		opt.FltMax[i] = 1
	}
	nf, ng, nh := 2, 1, 0

	// extra problem variables
	var f1max float64
	var fcn goga.MinProb_t
	var extraplot func()

	// problems
	switch problem {

	// problem # 1 -- CTP1, Deb 2001, p367, fig 225
	case 1:
		ng = 2
		f1max = 1.0
		a0, b0 := 0.858, 0.541
		a1, b1 := 0.728, 0.295
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c0 := 1.0
			for i := 1; i < len(x); i++ {
				c0 += x[i]
			}
			f[0] = x[0]
			f[1] = c0 * math.Exp(-x[0]/c0)
			if true {
				g[0] = f[1] - a0*math.Exp(-b0*f[0])
				g[1] = f[1] - a1*math.Exp(-b1*f[0])
			}
		}
		extraplot = func() {
			np := 201
			X, Y := utl.MeshGrid2D(0, 1, 0, 1, np, np)
			Z1 := utl.DblsAlloc(np, np)
			Z2 := utl.DblsAlloc(np, np)
			for j := 0; j < np; j++ {
				for i := 0; i < np; i++ {
					f0, f1 := X[i][j], Y[i][j]
					Z1[i][j] = f1 - a0*math.Exp(-b0*f0)
					Z2[i][j] = f1 - a1*math.Exp(-b1*f0)
				}
			}
			plt.ContourSimple(X, Y, Z1, false, 7, "linestyles=['--'], linewidths=[0.7], colors=['g'], levels=[0]")
			plt.ContourSimple(X, Y, Z2, false, 7, "linestyles=['-'],  linewidths=[0.7], colors=['b'], levels=[0]")
		}

	// problem # 2 -- CTP2, Deb 2001, p368/369, fig 226
	case 2:
		f1max = 1.2
		θ, a, b := -0.2*PI, 0.2, 10.0
		c, d, e := 1.0, 6.0, 1.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)

	// problem # 3 -- CTP3, Deb 2001, p368/370, fig 227
	case 3:
		f1max = 1.2
		θ, a, b := -0.2*PI, 0.1, 10.0
		c, d, e := 1.0, 0.5, 1.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)

	// problem # 4 -- CTP4, Deb 2001, p368/370, fig 228
	case 4:
		f1max = 2.0
		θ, a, b := -0.2*PI, 0.75, 10.0
		c, d, e := 1.0, 0.5, 1.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)

	// problem # 5 -- CTP5, Deb 2001, p368/371, fig 229
	case 5:
		f1max = 1.2
		θ, a, b := -0.2*PI, 0.1, 10.0
		c, d, e := 2.0, 0.5, 1.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)

	// problem # 6 -- CTP6, Deb 2001, p368/372, fig 230
	case 6:
		f1max = 5.0
		θ, a, b := 0.1*PI, 40.0, 0.5
		c, d, e := 1.0, 2.0, -2.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)

	// problem # 7 -- CTP7, Deb 2001, p368/373, fig 231
	case 7:
		f1max = 1.2
		θ, a, b := -0.05*PI, 40.0, 5.0
		c, d, e := 1.0, 6.0, 0.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)

	// problem # 8 -- CTP8, Deb 2001, p368/373, fig 232
	case 8:
		ng = 2
		f1max = 5.0
		θ1, a, b := 0.1*PI, 40.0, 0.5
		c, d, e := 1.0, 2.0, -2.0
		θ2, A, B := -0.05*PI, 40.0, 2.0
		C, D, E := 1.0, 6.0, 0.0
		sin1, cos1 := math.Sin(θ1), math.Cos(θ1)
		sin2, cos2 := math.Sin(θ2), math.Cos(θ2)
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c0 := 1.0
			for i := 1; i < len(x); i++ {
				c0 += x[i]
			}
			f[0] = x[0]
			f[1] = c0 * (1.0 - f[0]/c0)
			if true {
				c1 := cos1*(f[1]-e) - sin1*f[0]
				c2 := sin1*(f[1]-e) + cos1*f[0]
				c3 := math.Sin(b * PI * math.Pow(c2, c))
				g[0] = c1 - a*math.Pow(math.Abs(c3), d)
				d1 := cos2*(f[1]-E) - sin2*f[0]
				d2 := sin2*(f[1]-E) + cos2*f[0]
				d3 := math.Sin(B * PI * math.Pow(d2, C))
				g[1] = d1 - A*math.Pow(math.Abs(d3), D)
			}
		}
		extraplot = func() {
			np := 201
			X, Y := utl.MeshGrid2D(0, 1, 0, f1max, np, np)
			Z1 := utl.DblsAlloc(np, np)
			Z2 := utl.DblsAlloc(np, np)
			Z3 := utl.DblsAlloc(np, np)
			Z4 := utl.DblsAlloc(np, np)
			for j := 0; j < np; j++ {
				for i := 0; i < np; i++ {
					c1 := cos1*(Y[i][j]-e) - sin1*X[i][j]
					c2 := sin1*(Y[i][j]-e) + cos1*X[i][j]
					c3 := math.Sin(b * PI * math.Pow(c2, c))
					d1 := cos2*(Y[i][j]-E) - sin2*X[i][j]
					d2 := sin2*(Y[i][j]-E) + cos2*X[i][j]
					d3 := math.Sin(B * PI * math.Pow(d2, C))
					Z1[i][j] = c1
					Z2[i][j] = c1 - a*math.Pow(math.Abs(c3), d)
					Z3[i][j] = d1
					Z4[i][j] = d1 - A*math.Pow(math.Abs(d3), D)
				}
			}
			plt.ContourSimple(X, Y, Z1, false, 7, "linestyles=['--'], linewidths=[0.7], colors=['g'], levels=[0]")
			plt.ContourSimple(X, Y, Z2, false, 7, "linestyles=['-'],  linewidths=[0.7], colors=['b'], levels=[0]")
			plt.ContourSimple(X, Y, Z3, false, 7, "linestyles=['--'], linewidths=[0.7], colors=['g'], levels=[0]")
			plt.ContourSimple(X, Y, Z4, false, 7, "linestyles=['-'],  linewidths=[0.7], colors=['b'], levels=[0]")
		}

	default:
		chk.Panic("problem %d is not available", problem)
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

	// check
	var failed bool
	for _, sol := range opt.Solutions {
		for _, oor := range sol.Oor {
			if oor > 0 {
				failed = true
			}
		}
	}
	if failed {
		io.PfRed("failed\n")
	} else {
		io.PfGreen("OK\n")
	}

	// plot
	if doplot {
		feasibleOnly := false
		plt.SetForEps(0.8, 300)
		fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
		fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
		goga.PlotOvaOvaPareto(opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
		extraplot()
		if problem < 8 {
			plt.Text(0.05, 0.05, "unfeasible", "color='gray', ha='left',va='bottom'")
			plt.Text(0.95, f1max-0.05, "feasible", "color='gray', ha='right',va='top'")
		}
		plt.AxisYrange(0, f1max)
		plt.SaveD("/tmp/goga", io.Sf("constrmulti_%s.eps", opt.RptName))
	}
	return
}

func main() {
	ntrials := 1
	P := utl.IntRange2(1, 9)
	//P := []int{1}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem, ntrials, true)
	}
	io.Pf("\n-------------------------- generating report --------------------------\nn")
	nRowPerTab := 6
	goga.TexF1F0Report("/tmp/goga", "tmp_constrmulti", "constrmulti", nRowPerTab, true, opts)
	goga.TexF1F0Report("/tmp/goga", "constrmulti", "constrmulti", nRowPerTab, false, opts)
	//io.Pf("\n%v\n", opts[0].LogParams())
}

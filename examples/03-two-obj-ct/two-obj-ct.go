// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

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

// main function
func main() {

	// problem numbers
	P := utl.IntRange2(0, 9)
	//P := []int{0}

	// allocate and run each problem
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = twoObjCt(problem)
	}

	// report
	io.Pf("\n----------------------------------- generating report -----------------------------------\n\n")
	rpt := goga.NewTexReport(opts)
	rpt.ShowDescription = false
	rpt.ShowLmin = false
	rpt.ShowLave = false
	rpt.ShowLmax = false
	rpt.ShowLdev = false
	rpt.Title = "Constrained two-objective problems"
	rpt.Generate("/tmp/goga", "two-obj-ct")
}

// twoObjCt runs two-objective problem with constraints
func twoObjCt(problem int) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// options
	doPlot := false
	withSols0 := false
	constantSeed := false

	// parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Ncpu = 3
	opt.Tmax = 500
	opt.Verbose = false
	opt.VerbStat = false
	opt.GenType = "latin"
	opt.DEC = 0.1
	opt.Nsamples = 3 //////////////////////// increase this number

	// options for report
	opt.RptFmtE = "%.2e"
	opt.RptFmtEdev = "%.2e"

	// problem variables
	nx := 10
	opt.RptName = io.Sf("CTP%d", problem)
	opt.Nsol = 120
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

	// problem # 0 -- TNK
	case 0:
		ng = 2
		f1max = 1.21
		opt.RptName = "TNK"
		opt.FltMin = []float64{0, 0}
		opt.FltMax = []float64{PI, PI}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = x[0]
			f[1] = x[1]
			g[0] = x[0]*x[0] + x[1]*x[1] - 1.0 - 0.1*math.Cos(16.0*math.Atan2(x[0], x[1]))
			g[1] = 0.5 - math.Pow(x[0]-0.5, 2.0) - math.Pow(x[1]-0.5, 2.0)
		}
		extraplot = func() {
			np := 301
			X, Y := utl.MeshGrid2d(0, 1.3, 0, 1.3, np, np)
			Z1, Z2, Z3 := utl.Alloc(np, np), utl.Alloc(np, np), utl.Alloc(np, np)
			for j := 0; j < np; j++ {
				for i := 0; i < np; i++ {
					g1 := 0.5 - math.Pow(X[i][j]-0.5, 2.0) - math.Pow(Y[i][j]-0.5, 2.0)
					if g1 >= 0 {
						Z1[i][j] = X[i][j]*X[i][j] + Y[i][j]*Y[i][j] - 1.0 - 0.1*math.Cos(16.0*math.Atan2(Y[i][j], X[i][j]))
					} else {
						Z1[i][j] = -1
					}
					Z2[i][j] = X[i][j]*X[i][j] + Y[i][j]*Y[i][j] - 1.0 - 0.1*math.Cos(16.0*math.Atan2(Y[i][j], X[i][j]))
					Z3[i][j] = g1
				}
			}
			plt.ContourF(X, Y, Z1, &plt.A{Levels: []float64{0, 2}, NoCbar: true, Lw: 0.5, CmapIdx: 6, Fsz: 5})
			plt.Text(0.3, 0.95, "0.000", &plt.A{Fsz: 5, Rot: 10})
			plt.ContourL(X, Y, Z2, &plt.A{Colors: []string{"k"}, Levels: []float64{0}, Lw: 0.7})
			plt.ContourL(X, Y, Z3, &plt.A{Colors: []string{"k"}, Levels: []float64{0}, Lw: 1.0})
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] - 1.0 - 0.1*math.Cos(16.0*math.Atan2(f[0], f[1]))
		}

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
		f0a := math.Log(a0) / (b0 - 1.0)
		f1a := math.Exp(-f0a)
		f0b := math.Log(a0/a1) / (b0 - b1)
		f1b := a0 * math.Exp(-b0*f0b)
		opt.Multi_fcnErr = func(f []float64) float64 {
			if f[0] < f0a {
				return f[1] - math.Exp(-f[0])
			}
			if f[0] < f0b {
				return f[1] - a0*math.Exp(-b0*f[0])
			}
			return f[1] - a1*math.Exp(-b1*f[0])
		}
		extraplot = func() {
			np := 201
			X, Y := utl.MeshGrid2d(0, 1, 0, 1, np, np)
			Z := utl.Alloc(np, np)
			for j := 0; j < np; j++ {
				for i := 0; i < np; i++ {
					Z[i][j] = opt.Multi_fcnErr([]float64{X[i][j], Y[i][j]})
				}
			}
			plt.ContourF(X, Y, Z, &plt.A{Levels: []float64{0, 0.6}, NoCbar: true, Lw: 0.5, CmapIdx: 6})
			F0 := utl.LinSpace(0, 1, 21)
			F1r := make([]float64, len(F0))
			F1s := make([]float64, len(F0))
			F1t := make([]float64, len(F0))
			for i, f0 := range F0 {
				F1r[i] = math.Exp(-f0)
				F1s[i] = a0 * math.Exp(-b0*f0)
				F1t[i] = a1 * math.Exp(-b1*f0)
			}
			plt.Plot(F0, F1r, &plt.A{C: "blue", Ls: "--"})
			plt.Plot(F0, F1s, &plt.A{C: "green", Ls: "--"})
			plt.Plot(F0, F1t, &plt.A{C: "gray", Ls: "--"})
			plt.PlotOne(f0a, f1a, &plt.A{C: "k", M: "|", Ms: 20})
			plt.PlotOne(f0b, f1b, &plt.A{C: "k", M: "|", Ms: 20})
		}

	// problem # 2 -- CTP2, Deb 2001, p368/369, fig 226
	case 2:
		f1max = 1.2
		θ, a, b := -0.2*PI, 0.2, 10.0
		c, d, e := 1.0, 6.0, 1.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)
		opt.Multi_fcnErr = CTPerror1(θ, a, b, c, d, e)

	// problem # 3 -- CTP3, Deb 2001, p368/370, fig 227
	case 3:
		f1max = 1.2
		θ, a, b := -0.2*PI, 0.1, 10.0
		c, d, e := 1.0, 0.5, 1.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)
		opt.Multi_fcnErr = CTPerror1(θ, a, b, c, d, e)

	// problem # 4 -- CTP4, Deb 2001, p368/370, fig 228
	case 4:
		f1max = 2.0
		θ, a, b := -0.2*PI, 0.75, 10.0
		c, d, e := 1.0, 0.5, 1.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)
		opt.Multi_fcnErr = CTPerror1(θ, a, b, c, d, e)

	// problem # 5 -- CTP5, Deb 2001, p368/371, fig 229
	case 5:
		f1max = 1.2
		θ, a, b := -0.2*PI, 0.1, 10.0
		c, d, e := 2.0, 0.5, 1.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = CTPplotter(θ, a, b, c, d, e, f1max)
		opt.Multi_fcnErr = CTPerror1(θ, a, b, c, d, e)

	// problem # 6 -- CTP6, Deb 2001, p368/372, fig 230
	case 6:
		f1max = 5.0
		θ, a, b := 0.1*PI, 40.0, 0.5
		c, d, e := 1.0, 2.0, -2.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		extraplot = func() {
			np := 201
			X, Y := utl.MeshGrid2d(0, 1, 0, 20, np, np)
			Z := utl.Alloc(np, np)
			for j := 0; j < np; j++ {
				for i := 0; i < np; i++ {
					Z[i][j] = CTPconstraint(θ, a, b, c, d, e, X[i][j], Y[i][j])
				}
			}
			plt.ContourF(X, Y, Z, &plt.A{Levels: []float64{-30, -15, 0, 15, 30}, NoCbar: true, Lw: 0.5, CmapIdx: 6, Fsz: 5})
		}
		opt.Multi_fcnErr = CTPerror1(θ, a, b, c, d, e)

	// problem # 7 -- CTP7, Deb 2001, p368/373, fig 231
	case 7:
		f1max = 1.2
		θ, a, b := -0.05*PI, 40.0, 5.0
		c, d, e := 1.0, 6.0, 0.0
		fcn = CTPgenerator(θ, a, b, c, d, e)
		opt.Multi_fcnErr = func(f []float64) float64 { return f[1] - (1.0 - f[0]) }
		extraplot = func() {
			np := 201
			X, Y := utl.MeshGrid2d(0, 1, 0, f1max, np, np)
			Z1 := utl.Alloc(np, np)
			Z2 := utl.Alloc(np, np)
			for j := 0; j < np; j++ {
				for i := 0; i < np; i++ {
					Z1[i][j] = opt.Multi_fcnErr([]float64{X[i][j], Y[i][j]})
					Z2[i][j] = CTPconstraint(θ, a, b, c, d, e, X[i][j], Y[i][j])
				}
			}
			plt.ContourF(X, Y, Z2, &plt.A{Levels: []float64{0, 3}, NoCbar: true, Lw: 0.5, CmapIdx: 6, Fsz: 5})
			plt.ContourL(X, Y, Z1, &plt.A{Levels: []float64{0}, Colors: []string{"b"}, Ls: "--", Lw: 0.7})
		}

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
			np := 401
			X, Y := utl.MeshGrid2d(0, 1, 0, 20, np, np)
			Z1 := utl.Alloc(np, np)
			Z2 := utl.Alloc(np, np)
			Z3 := utl.Alloc(np, np)
			for j := 0; j < np; j++ {
				for i := 0; i < np; i++ {
					c1 := cos1*(Y[i][j]-e) - sin1*X[i][j]
					c2 := sin1*(Y[i][j]-e) + cos1*X[i][j]
					c3 := math.Sin(b * PI * math.Pow(c2, c))
					d1 := cos2*(Y[i][j]-E) - sin2*X[i][j]
					d2 := sin2*(Y[i][j]-E) + cos2*X[i][j]
					d3 := math.Sin(B * PI * math.Pow(d2, C))
					Z1[i][j] = c1 - a*math.Pow(math.Abs(c3), d)
					Z2[i][j] = d1 - A*math.Pow(math.Abs(d3), D)
					if Z1[i][j] >= 0 && Z2[i][j] >= 0 {
						Z3[i][j] = 1
					} else {
						Z3[i][j] = -1
					}
				}
			}
			plt.ContourF(X, Y, Z3, &plt.A{Colors: []string{"white", "grey"}, NoLabels: false, NoCbar: true, Lw: 0.5, Fsz: 5})
			plt.ContourL(X, Y, Z1, &plt.A{Levels: []float64{0}, Colors: []string{"gray"}, Ls: "--", Lw: 0.7})
			plt.ContourL(X, Y, Z2, &plt.A{Levels: []float64{0}, Colors: []string{"gray"}, Ls: "--", Lw: 0.7})
		}
		opt.Multi_fcnErr = CTPerror1(θ1, a, b, c, d, e)

	default:
		chk.Panic("problem %d is not available", problem)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// initial solutions
	var sols0 []*goga.Solution
	if withSols0 {
		sols0 = opt.GetSolutionsCopy()
	}

	// solve
	opt.RunMany("", "", constantSeed)
	opt.PrintStatMultiE()

	// check
	goga.CheckFront0(opt, true)

	// plot
	if doPlot {
		plt.Reset(false, nil)
		pp := goga.NewPlotParams(false)
		pp.FnKey = opt.RptName
		pp.Extra = func() {
			extraplot()
			if problem > 0 && problem < 6 {
				plt.Text(0.05, 0.05, "unfeasible", &plt.A{C: "gray", Ha: "left", Va: "bottom"})
				plt.Text(0.95, f1max-0.05, "feasible", &plt.A{C: "white", Ha: "right", Va: "top"})
			}
			glb := &plt.A{Rot: -7, C: "gray", Ha: "left", Va: "bottom"}
			wcb := &plt.A{Rot: -7, C: "white", Ha: "center", Va: "bottom"}
			if opt.RptName == "CTP6" {
				plt.Text(0.02, 0.15, "unfeasible", glb)
				plt.Text(0.02, 6.50, "unfeasible", glb)
				plt.Text(0.02, 13.0, "unfeasible", glb)
				plt.Text(0.50, 2.40, "feasible", wcb)
				plt.Text(0.50, 8.80, "feasible", wcb)
				plt.Text(0.50, 15.30, "feasible", wcb)
			}
			if opt.RptName == "TNK" {
				plt.Text(0.05, 0.05, "unfeasible", glb)
				plt.Text(0.80, 0.85, "feasible", &plt.A{C: "white", Ha: "left", Va: "top"})
				plt.Equal()
				plt.AxisRange(0, 1.22, 0, 1.22)
			}
		}
		opt.PlotOvaOvaPareto(sols0, 0, 1, pp)
	}
	return
}

// CTPconstraint implements the constraint function in CTP problem
func CTPconstraint(θ, a, b, c, d, e float64, f0, f1 float64) (g0 float64) {
	sθ, cθ := math.Sin(θ), math.Cos(θ)
	c1 := cθ*(f1-e) - sθ*f0
	c2 := sθ*(f1-e) + cθ*f0
	c3 := math.Sin(b * PI * math.Pow(c2, c))
	return c1 - a*math.Pow(math.Abs(c3), d)
}

// CTPgenerator generates CTP problem
func CTPgenerator(θ, a, b, c, d, e float64) goga.MinProb_t {
	return func(f, g, h, x []float64, ξ []int, cpu int) {
		c0 := 1.0
		for i := 1; i < len(x); i++ {
			c0 += x[i]
		}
		f[0] = x[0]
		f[1] = c0 * (1.0 - f[0]/c0)
		g[0] = CTPconstraint(θ, a, b, c, d, e, f[0], f[1])
	}
}

// CTPplotter plots CTP problem
func CTPplotter(θ, a, b, c, d, e, f1max float64) func() {
	return func() {
		np := 401
		X, Y := utl.MeshGrid2d(0, 1, 0, f1max, np, np)
		Z1 := utl.Alloc(np, np)
		Z2 := utl.Alloc(np, np)
		sθ, cθ := math.Sin(θ), math.Cos(θ)
		for j := 0; j < np; j++ {
			for i := 0; i < np; i++ {
				f0, f1 := X[i][j], Y[i][j]
				Z1[i][j] = cθ*(f1-e) - sθ*f0
				Z2[i][j] = CTPconstraint(θ, a, b, c, d, e, X[i][j], Y[i][j])
			}
		}
		//plt.Contour(X, Y, Z2, "levels=[0,2],cbar=0,lwd=0.5,fsz=5,cmapidx=6")
		//plt.ContourSimple(X, Y, Z1, false, 7, "linestyles=['--'], linewidths=[0.7], colors=['b'], levels=[0]")
	}
}

// CTPerror1 implements the error function in CTP problem
func CTPerror1(θ, a, b, c, d, e float64) func(f []float64) float64 {
	return func(f []float64) float64 {
		return CTPconstraint(θ, a, b, c, d, e, f[0], f[1])
	}
}

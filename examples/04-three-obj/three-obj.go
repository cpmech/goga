// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
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

func solve_problem(problem int) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 200
	opt.Ncpu = 5
	opt.Tf = 500
	opt.Nsamples = 2
	opt.DEC = 0.01

	// options for report
	opt.HistNsta = 6
	opt.HistLen = 13
	opt.RptFmtE = "%.4e"
	opt.RptFmtL = "%.4e"
	opt.RptFmtEdev = "%.3e"
	opt.RptFmtLdev = "%.3e"
	opt.RptFmin = make([]float64, 3)
	opt.RptFmax = make([]float64, 3)
	for i := 0; i < 3; i++ {
		opt.RptFmax[i] = 1
	}

	// problem variables
	var αcone float64        // cone half-opening angle
	var nf, ng, nh int       // number of functions
	var fcn goga.MinProb_t   // functions
	var plot_solution func() // plot solution in 3D

	// problems
	switch problem {

	// DTLZ1
	case 1:
		opt.RptName = "DTLZ1"
		opt.FltMin = make([]float64, 7)
		opt.FltMax = make([]float64, 7)
		for i := 0; i < 7; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c := 5.0
			for i := 2; i < 7; i++ {
				c += math.Pow((x[i]-0.5), 2.0) - math.Cos(20.0*PI*(x[i]-0.5))
			}
			c *= 100.0
			f[0] = 0.5 * x[0] * x[1] * (1.0 + c)
			f[1] = 0.5 * x[0] * (1.0 - x[1]) * (1.0 + c)
			f[2] = 0.5 * (1.0 - x[0]) * (1.0 + c)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0] + f[1] + f[2] - 0.5
		}
		plot_solution = func() { plot_plane(false) }
		opt.RptFmax = []float64{0.5, 0.5, 0.5}

	// DTLZ2
	case 2:
		opt.RptName = "DTLZ2"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() { plot_sphere(false) }

	// DTLZ3
	case 3:
		opt.RptName = "DTLZ3"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c := 10.0
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i]-0.5), 2.0) - math.Cos(20.0*PI*(x[i]-0.5))
			}
			c *= 100.0
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() { plot_sphere(false) }

	// DTLZ4
	case 4:
		opt.RptName = "DTLZ4"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			a := 100.0
			f[0] = (1.0 + c) * math.Cos(math.Pow(x[0], a)*PI/2.0) * math.Cos(math.Pow(x[1], a)*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(math.Pow(x[0], a)*PI/2.0) * math.Sin(math.Pow(x[1], a)*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(math.Pow(x[0], a)*PI/2.0)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() { plot_sphere(false) }

	// DTLZ2x (convex)
	case 5:
		opt.RptName = "DTLZ2x"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
			f[0] = math.Pow(f[0], 4.0)
			f[1] = math.Pow(f[1], 4.0)
			f[2] = math.Pow(f[2], 2.0)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return math.Pow(math.Abs(f[0]), 0.5) + math.Pow(math.Abs(f[1]), 0.5) + f[2] - 1.0
		}
		plot_solution = func() { plot_convex(1.0, false) }

	// DTLZ2c (constraint)
	case 6:
		opt.RptName = "DTLZ2c"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 1, 0
		//αcone = math.Atan(1.0 / SQ2) // <<< touches lower plane
		//αcone = PI/2.0 - αcone // <<< touches upper plane
		αcone = 15.0 * PI / 180.0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
			g[0] = math.Tan(αcone) - cone_angle(f)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() {
			plot_sphere(false)
			plot_cone(αcone, true)
		}

	// Superquadric 1
	case 7:
		opt.RptName = "SUQ1"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		a, b, c := 0.5, 0.5, 0.5
		A, B, C := 2.0/a, 2.0/b, 2.0/c
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * cosX(x[0]*PI/2.0, A) * cosX(x[1]*PI/2.0, A)
			f[1] = (1.0 + c) * cosX(x[0]*PI/2.0, B) * sinX(x[1]*PI/2.0, B)
			f[2] = (1.0 + c) * sinX(x[0]*PI/2.0, C)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return math.Pow(math.Abs(f[0]), a) + math.Pow(math.Abs(f[1]), b) + math.Pow(math.Abs(f[2]), c) - 1.0
		}
		plot_solution = func() { plot_superquadric(a, b, c, false) }

	// Superquadric 2
	case 8:
		opt.RptName = "SUQ2"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		a, b, c := 2.0, 1.0, 0.5
		A, B, C := 2.0/a, 2.0/b, 2.0/c
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * cosX(x[0]*PI/2.0, A) * cosX(x[1]*PI/2.0, A)
			f[1] = (1.0 + c) * cosX(x[0]*PI/2.0, B) * sinX(x[1]*PI/2.0, B)
			f[2] = (1.0 + c) * sinX(x[0]*PI/2.0, C)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return math.Pow(math.Abs(f[0]), a) + math.Pow(math.Abs(f[1]), b) + math.Pow(math.Abs(f[2]), c) - 1.0
		}
		plot_solution = func() { plot_superquadric(a, b, c, false) }

	default:
		chk.Panic("problem %d is not available", problem)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.RunMany("", "")
	goga.StatMulti(opt, true)

	// check
	goga.CheckFront0(opt, true)

	// print results
	if false {
		goga.SortByOva(opt.Solutions, 0)
		m, l := opt.Nsol/2, opt.Nsol-1
		A, B, C := opt.Solutions[0], opt.Solutions[m], opt.Solutions[l]
		io.Pforan("A = %v\n", A.Flt)
		io.Pforan("B = %v\n", B.Flt)
		io.Pforan("C = %v\n", C.Flt)
	}

	// plot results
	if false {
		py_plot3(0, 1, nf-1, opt, plot_solution, true, true)
	}

	// vtk
	if false {
		ptRad := 0.015
		if opt.RptName == "DTLZ1" {
			ptRad = 0.01
		}
		vtk_plot3(opt, αcone, ptRad, true, true)
	}

	// star plot
	if false {
		plt.SetForEps(1, 300)
		goga.PlotStar(opt)
		plt.SaveD("/tmp/goga", io.Sf("starplot_%s.eps", opt.RptName))
	}

	// write all results
	if false {
		goga.WriteAllValues("/tmp/goga", "res_three-obj", opt)
	}
	return
}

func main() {
	P := utl.IntRange2(1, 9)
	//P := []int{2}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem)
	}
	io.Pf("\n-------------------------- generating report --------------------------\nn")
	rpt := goga.NewTexReport(opts)
	rpt.NRowPerTab = 10
	rpt.Type = 1
	rpt.Title = "Unconstrained and constrained three objective problems."
	rpt.Fnkey = "three-obj"
	rpt.Generate()
}

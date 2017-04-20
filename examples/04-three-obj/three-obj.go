// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// constants
const (
	PI        = math.Pi
	Color     = "#1e5ec9"
	Linewidth = 0.5
	Eps       = true
)

// main function
func main() {

	// problem numbers
	//P := utl.IntRange2(1, 9)
	//P := []int{7, 8}
	P := []int{8}

	// allocate and run each problem
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = threeObj(problem)
	}

	return

	// report
	io.Pf("\n----------------------------------- generating report -----------------------------------\n\n")
	rpt := goga.NewTexReport(opts)
	rpt.ShowDescription = false
	rpt.ShowLmin = false
	rpt.ShowLave = false
	rpt.ShowLmax = false
	rpt.ShowLdev = false
	rpt.Title = "Constrained and unconstrained three-objective problems"
	rpt.Generate("/tmp/goga", "three-obj")
}

// threeObj runs three-obj problem
func threeObj(problem int) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// options
	plotVTK := false
	plotPy := true
	plotStar := false
	writeRes := false
	printResults := false
	constantSeed := false

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 200
	opt.Ncpu = 5
	opt.Tmax = 500
	opt.DEC = 0.01
	opt.Nsamples = 2 /////////////////////// increase this number

	// options for report
	opt.RptFmtE = "%.2e"
	opt.RptFmtEdev = "%.2e"
	opt.RptFmin = make([]float64, 3)
	opt.RptFmax = make([]float64, 3)
	for i := 0; i < 3; i++ {
		opt.RptFmax[i] = 1
	}

	// plot arguments
	pltArgs := &plt.A{C: Color, Lw: Linewidth}

	// problem variables
	var αcone float64      // cone half-opening angle
	var nf, ng, nh int     // number of functions
	var fcn goga.MinProb_t // functions
	var surface func()     // plot solution in 3D

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
		surface = func() {
			p := []float64{0.5, 0, 0}
			n := []float64{1, 1, 1}
			d := -n[0]*p[0] - n[1]*p[1] - n[2]*p[2]
			xmin, xmax := 0.0, 0.5
			ymin, ymax := 0.0, 0.5
			nu, nv := 21, 21
			X, Y, Z := utl.MeshGrid2dF(xmin, xmax, ymin, ymax, nu, nv, func(x, y float64) float64 {
				z := (-d - n[0]*x - n[1]*y) / n[2]
				if z < -0.01 {
					return math.NaN()
				}
				return z
			})
			plt.Wireframe(X, Y, Z, pltArgs)
		}
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
		surface = func() {
			plt.Hemisphere(nil, 1, 0, 90, 21, 21, false, pltArgs)
		}

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
		surface = func() {
			plt.Hemisphere(nil, 1, 0, 90, 21, 21, false, pltArgs)
		}

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
		surface = func() {
			plt.Hemisphere(nil, 1, 0, 90, 21, 21, false, pltArgs)
		}

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
		surface = func() {
			level := 1.0
			X, Y, Z := utl.MeshGrid2dF(0, 1, 0, 1, 21, 21, func(x, y float64) float64 {
				z := level - math.Sqrt(x) - math.Sqrt(y)
				if z < -0.01 {
					z = math.NaN()
				}
				return z
			})
			plt.Wireframe(X, Y, Z, pltArgs)
		}

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
		alpDeg := 15.0
		αcone = alpDeg * PI / 180.0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
			g[0] = math.Tan(αcone) - plt.CalcDiagAngle(f)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		surface = func() {
			plt.Hemisphere(nil, 1, 0, 90, 21, 21, false, pltArgs)
			plt.ConeDiag([]float64{0, 0, 0}, alpDeg, 1.3, 7, 31, &plt.A{C: "k", Lw: Linewidth})
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
			var r float64
			for i := 2; i < 12; i++ {
				r += math.Pow((x[i] - 0.5), 2.0)
			}
			r = 1.0 + r
			f[0] = r * fun.SuqCos(x[0]*PI/2.0, A) * fun.SuqCos(x[1]*PI/2.0, A)
			f[1] = r * fun.SuqCos(x[0]*PI/2.0, B) * fun.SuqSin(x[1]*PI/2.0, B)
			f[2] = r * fun.SuqSin(x[0]*PI/2.0, C)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return math.Pow(math.Abs(f[0]), a) + math.Pow(math.Abs(f[1]), b) + math.Pow(math.Abs(f[2]), c) - 1.0
		}
		surface = func() {
			p := []float64{0, 0, 0}
			r := []float64{1, 1, 1}
			e := []float64{a, b, c}
			alpmin, alpmax := 0.0, 90.0
			etamin, etamax := 0.0, 90.0
			nalp, neta := 21, 21
			plt.Superquadric(p, r, e, alpmin, alpmax, etamin, etamax, nalp, neta, pltArgs)
		}

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
			var r float64
			for i := 2; i < 12; i++ {
				r += math.Pow((x[i] - 0.5), 2.0)
			}
			r = 1.0 + r
			f[0] = r * fun.SuqCos(x[0]*PI/2.0, A) * fun.SuqCos(x[1]*PI/2.0, A)
			f[1] = r * fun.SuqCos(x[0]*PI/2.0, B) * fun.SuqSin(x[1]*PI/2.0, B)
			f[2] = r * fun.SuqSin(x[0]*PI/2.0, C)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return math.Pow(math.Abs(f[0]), a) + math.Pow(math.Abs(f[1]), b) + math.Pow(math.Abs(f[2]), c) - 1.0
		}
		surface = func() {
			p := []float64{0, 0, 0}
			r := []float64{1, 1, 1}
			e := []float64{a, b, c}
			alpmin, alpmax := 0.0, 90.0
			etamin, etamax := 0.0, 90.0
			nalp, neta := 21, 21
			plt.Superquadric(p, r, e, alpmin, alpmax, etamin, etamax, nalp, neta, pltArgs)
		}

	default:
		chk.Panic("problem %d is not available", problem)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.RunMany("", "", constantSeed)
	opt.PrintStatMultiE()

	// check
	goga.CheckFront0(opt, true)

	// print results
	if printResults {
		goga.SortSolutions(opt.Solutions, 0)
		m, l := opt.Nsol/2, opt.Nsol-1
		A, B, C := opt.Solutions[0], opt.Solutions[m], opt.Solutions[l]
		io.Pforan("A = %v\n", A.Flt)
		io.Pforan("B = %v\n", B.Flt)
		io.Pforan("C = %v\n", C.Flt)
	}

	// plot results
	if plotPy {
		PyPlot3(0, 1, nf-1, opt, surface, true)
	}

	// vtk
	if plotVTK {
		ptRad := 0.015
		if opt.RptName == "DTLZ1" {
			ptRad = 0.01
		}
		twice := false
		VtkPlot3(opt, αcone, ptRad, true, twice)
	}

	// star plot
	if plotStar {
		plt.Reset(false, nil)
		opt.PlotStar()
		plt.Save("/tmp/goga", io.Sf("starplot_%s", opt.RptName))
	}

	// write all results
	if writeRes {
		goga.WriteAllValues("/tmp/goga", "res_three-obj", opt)
	}
	return
}

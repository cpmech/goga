// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_flt01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt01. sin⁶(5 π x) multimodal")

	// parameters
	var opt Optimiser
	opt.Default()
	opt.Nsol = 20
	opt.Ncpu = 1
	opt.FltMin = []float64{0}
	opt.FltMax = []float64{1}
	nf, ng, nh := 1, 0, 0

	// initialise optimiser
	yfcn := func(x float64) float64 { return math.Pow(math.Sin(5.0*math.Pi*x), 6.0) }
	opt.Init(GenTrialSolutions, nil, func(f, g, h, x []float64, y []int, cpu int) {
		f[0] = -yfcn(x[0])
	}, nf, ng, nh)

	// initial solutions
	sols0 := opt.GetSolutionsCopy()

	// solve
	opt.Solve()

	// plot
	if chk.Verbose {
		pp := NewPlotParams(false)
		pp.FnKey = "fig-flt01"
		pp.YfuncX = yfcn
		opt.PlotFltOva(sols0, 0, 0, -1, pp)
	}
}

func Test_flt02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt02. quadratic function with inequalities")

	// parameters
	var opt Optimiser
	opt.Default()
	opt.Nsol = 20
	opt.Ncpu = 1
	opt.FltMin = []float64{-2, -2}
	opt.FltMax = []float64{2, 2}
	nf, ng, nh := 1, 5, 0

	// initialise optimiser
	opt.Init(GenTrialSolutions, nil, func(f, g, h, x []float64, y []int, cpu int) {
		f[0] = x[0]*x[0]/2.0 + x[1]*x[1] - x[0]*x[1] - 2.0*x[0] - 6.0*x[1]
		g[0] = 2.0 - x[0] - x[1]     // ≥ 0
		g[1] = 2.0 + x[0] - 2.0*x[1] // ≥ 0
		g[2] = 3.0 - 2.0*x[0] - x[1] // ≥ 0
		g[3] = x[0]                  // ≥ 0
		g[4] = x[1]                  // ≥ 0
	}, nf, ng, nh)

	// initial solutions
	sols0 := opt.GetSolutionsCopy()

	// solve
	opt.Solve()

	// log
	io.Pforan("%v\n", opt.LogParams())

	// plot
	if chk.Verbose {
		pp := NewPlotParams(false)
		pp.FnKey = "fig-flt02"
		opt.PlotFltFltContour(sols0, 0, 1, 0, pp)
	}
}

func Test_flt03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt03. circle with equality constraint")

	// geometry
	xe := 1.0                      // centre of circle
	le := -0.4                     // selected level of f(x)
	ys := xe - (1.0+le)/math.Sqrt2 // coordinates of minimum point with level=le
	y0 := 2.0*ys + xe              // vertical axis intersect of straight line defined by c(x)
	xc := []float64{xe, xe}        // centre

	// parameters
	var opt Optimiser
	opt.Default()
	opt.Nsol = 20
	opt.Ncpu = 1
	opt.Verbose = false
	opt.FltMin = []float64{-1, -1}
	opt.FltMax = []float64{3, 3}
	nf, ng, nh := 1, 0, 1

	// initialise optimiser
	opt.Init(GenTrialSolutions, nil, func(f, g, h, x []float64, y []int, cpu int) {
		res := 0.0
		for i := 0; i < len(x); i++ {
			res += (x[i] - xc[i]) * (x[i] - xc[i])
		}
		f[0] = math.Sqrt(res) - 1.0
		h[0] = x[0] + x[1] + xe - y0
	}, nf, ng, nh)

	// initial solutions
	sols0 := opt.GetSolutionsCopy()

	// solve
	opt.Solve()

	// plot
	if chk.Verbose {
		pp := NewPlotParams(false)
		pp.FnKey = "fig-flt03"
		pp.AxEqual = true
		plt.Reset(false, nil)
		opt.PlotFltFltContour(sols0, 0, 1, 0, pp)
	}
}

func Test_flt04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt04. two-bar truss. Pareto-optimal")

	// data. from Coelho (2007) page 19
	ρ := 0.283 // lb/in³
	H := 100.0 // in
	P := 1e4   // lb
	E := 3e7   // lb/in²
	σ0 := 2e4  // lb/in²

	// parameters
	var opt Optimiser
	opt.Default()
	opt.Nsol = 30
	opt.Ncpu = 1
	opt.Tmax = 100
	opt.LatinDup = 5
	opt.FltMin = []float64{0.1, 0.5}
	opt.FltMax = []float64{2.25, 2.5}
	nf, ng, nh := 2, 2, 0

	// initialise optimiser
	TSQ2 := 2.0 * math.Sqrt2
	opt.Init(GenTrialSolutions, nil, func(f, g, h, x []float64, y []int, cpu int) {
		f[0] = 2.0 * ρ * H * x[1] * math.Sqrt(1.0+x[0]*x[0])
		f[1] = P * H * math.Pow(1.0+x[0]*x[0], 1.5) * math.Sqrt(1.0+math.Pow(x[0], 4.0)) / (TSQ2 * E * x[0] * x[0] * x[1])
		g[0] = σ0 - P*(1.0+x[0])*math.Sqrt(1.0+x[0]*x[0])/(TSQ2*x[0]*x[1])
		g[1] = σ0 - P*(1.0-x[0])*math.Sqrt(1.0+x[0]*x[0])/(TSQ2*x[0]*x[1])
	}, nf, ng, nh)

	// initial solutions
	sols0 := opt.GetSolutionsCopy()

	// solve
	opt.Solve()

	// plot
	if chk.Verbose {
		_, dat := io.ReadTable("data/coelho-fig1.6.dat")
		pp := NewPlotParams(false)
		pp.FnKey = "fig-flt04"
		pp.WithAll = true
		pp.Extra = func() {
			plt.Plot(dat["f1"], dat["f2"], &plt.A{C: "b", Ms: 3, Mec: "b"})
			plt.AxisRange(0, 250, 0, 0.15)
		}
		opt.PlotOvaOvaPareto(sols0, 0, 1, pp)
	}
}

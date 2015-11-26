// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func Test_flt01(tst *testing.T) {

	verbose()
	chk.PrintTitle("flt01. sin⁶(5 π x) multimodal")

	// parameters
	var opt Optimiser
	opt.Default()
	opt.Ngrp = 1
	opt.Nsol = 20
	opt.FltMin = []float64{0}
	opt.FltMax = []float64{1}
	nf, ng, nh := 1, 0, 0

	// initialise optimiser
	yfcn := func(x float64) float64 { return math.Pow(math.Sin(5.0*math.Pi*x), 6.0) }
	opt.Init(GenTrialSolutions, nil, func(f, g, h, x []float64, ξ []int, grp int) {
		f[0] = -yfcn(x[0])
	}, nf, ng, nh)

	// initial solutions
	sols0 := opt.GetSolutionsCopy()

	// solve
	opt.Solve()

	// plot
	PlotFltOva("fig_flt03", &opt, sols0, 0, 0, 201, -1, yfcn, nil, false)
}

func Test_flt02(tst *testing.T) {

	verbose()
	chk.PrintTitle("flt02. quadratic function with inequalities")

	// parameters
	var opt Optimiser
	opt.Default()
	opt.Ngrp = 1
	opt.Nsol = 20
	opt.FltMin = []float64{-2, -2}
	opt.FltMax = []float64{2, 2}
	nf, ng, nh := 1, 5, 0

	// initialise optimiser
	opt.Init(GenTrialSolutions, nil, func(f, g, h, x []float64, ξ []int, grp int) {
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

	// plot
	PlotFltFltContour("fig_flt01", &opt, sols0, 0, 1, 0, nil, false)
}

func Test_flt03(tst *testing.T) {

	verbose()
	chk.PrintTitle("flt03. circle with equality constraint")

	// geometry
	xe := 1.0                      // centre of circle
	le := -0.4                     // selected level of f(x)
	ys := xe - (1.0+le)/math.Sqrt2 // coordinates of minimum point with level=le
	y0 := 2.0*ys + xe              // vertical axis intersect of straight line defined by c(x)
	xc := []float64{xe, xe}        // centre

	// parameters
	var opt Optimiser
	//opt.Default()
	opt.Read("ga-data.json")
	opt.Verbose = false
	opt.FltMin = []float64{-1, -1}
	opt.FltMax = []float64{3, 3}
	nf, ng, nh := 1, 0, 1

	// initialise optimiser
	opt.Init(GenTrialSolutions, nil, func(f, g, h, x []float64, ξ []int, grp int) {
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
	PlotFltFltContour("fig_flt02", &opt, sols0, 0, 1, 0, nil, false)
}

func Test_flt04(tst *testing.T) {

	verbose()
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
	opt.Ngrp = 1
	opt.Nsol = 30
	opt.FltMin = []float64{0.1, 0.5}
	opt.FltMax = []float64{2.25, 2.5}
	nf, ng, nh := 2, 2, 0

	// initialise optimiser
	TSQ2 := 2.0 * math.Sqrt2
	opt.Init(GenTrialSolutions, nil, func(f, g, h, x []float64, ξ []int, grp int) {
		f[0] = 2.0 * ρ * H * x[1] * math.Sqrt(1.0+x[0]*x[0])
		f[1] = P * H * math.Pow(1.0+x[0]*x[0], 1.5) * math.Sqrt(1.0+math.Pow(x[0], 4.0)) / (TSQ2 * E * x[0] * x[0] * x[1])
		g[0] = σ0 - P*(1.0+x[0])*math.Sqrt(1.0+x[0]*x[0])/(TSQ2*x[0]*x[1])
		g[1] = σ0 - P*(1.0-x[0])*math.Sqrt(1.0+x[0]*x[0])/(TSQ2*x[0]*x[1])
	}, nf, ng, nh)

	// solve
	opt.Solve()

	// reference data
	_, dat, _ := io.ReadTable("data/coelho-fig1.6.dat")

	// plot
	PlotOvaOvaPareto("fig_flt04", &opt, nil, 0, 1, func() {
		plt.Plot(dat["f1"], dat["f2"], "'b*',ms=3,markeredgecolor='b'")
	}, nil, false)
}

func Test_flt05(tst *testing.T) {

	verbose()
	chk.PrintTitle("flt05. ZDT problems")

	// parameters
	var opt Optimiser
	opt.Default()
	opt.Ngrp = 1
	opt.Nsol = 30
	opt.Problem = 1

	// problem variables
	var pname string
	var fmin, fmax []float64
	var f1f0 func(f0 float64) float64
	var nf, ng, nh int
	var fcn MinProb_t

	// problems
	switch opt.Problem {

	// ZDT1, Deb 2001, p356
	case 1:
		pname = "ZDT1"
		n := 30
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, grp int) {
			f[0] = x[0]
			sum := 0.0
			for i := 1; i < n; i++ {
				sum += x[i]
			}
			c0 := 1.0 + 9.0*sum/float64(n-1)
			c1 := 1.0 - math.Sqrt(f[0]/c0)
			f[1] = c0 * c1
		}
		fmin = []float64{0, 0}
		fmax = []float64{1, 1}
		f1f0 = func(f0 float64) float64 {
			return 1.0 - math.Sqrt(f0)
		}

	// ZDT2, Deb 2001, p356
	case 2:
		pname = "ZDT2"
		n := 30
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, grp int) {
			f[0] = x[0]
			sum := 0.0
			for i := 1; i < n; i++ {
				sum += x[i]
			}
			c0 := 1.0 + 9.0*sum/float64(n-1)
			c1 := 1.0 - math.Pow(f[0]/c0, 2.0)
			f[1] = c0 * c1
		}
		fmin = []float64{0, 0}
		fmax = []float64{1, 1}
		f1f0 = func(f0 float64) float64 {
			return 1.0 - math.Pow(f0, 2.0)
		}

	// ZDT3, Deb 2001, p356
	case 3:
		pname = "ZDT3"
		n := 30
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, grp int) {
			f[0] = x[0]
			sum := 0.0
			for i := 1; i < n; i++ {
				sum += x[i]
			}
			c0 := 1.0 + 9.0*sum/float64(n-1)
			c1 := 1.0 - math.Sqrt(f[0]/c0) - math.Sin(10.0*math.Pi*f[0])*f[0]/c0
			f[1] = c0 * c1
		}
		fmin = []float64{0, -1}
		fmax = []float64{1, 1}
		f1f0 = func(f0 float64) float64 {
			return 1.0 - math.Sqrt(f0) - math.Sin(10.0*math.Pi*f0)*f0
		}

	// ZDT4, Deb 2001, p358
	case 4:
		pname = "ZDT4"
		n := 10
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = -5
			opt.FltMax[i] = 5
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, grp int) {
			f[0] = x[0]
			sum := 0.0
			w := 4.0 * math.Pi
			for i := 1; i < n; i++ {
				sum += x[i]*x[i] - 10.0*math.Cos(w*x[i])
			}
			c0 := 1.0 + 10.0*float64(n-1) + sum
			c1 := 1.0 - math.Sqrt(f[0]/c0)
			f[1] = c0 * c1
		}
		fmin = []float64{0, 0}
		fmax = []float64{1, 1}
		f1f0 = func(f0 float64) float64 {
			return 1.0 - math.Sqrt(f0)
		}

	// FON (Fonseca and Fleming 1995), Deb 2001, p339
	case 5:
		pname = "FON"
		n := 10
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		coef := 1.0 / math.Sqrt(float64(n))
		fcn = func(f, g, h, x []float64, ξ []int, grp int) {
			sum0, sum1 := 0.0, 0.0
			for i := 0; i < n; i++ {
				sum0 += math.Pow(x[i]-coef, 2.0)
				sum1 += math.Pow(x[i]+coef, 2.0)
			}
			f[0] = 1.0 - math.Exp(-sum0)
			f[1] = 1.0 - math.Exp(-sum1)
		}
		fmin = []float64{0, 0}
		fmax = []float64{0.98, 1}
		f1f0 = func(f0 float64) float64 {
			return 1.0 - math.Exp(-math.Pow(2.0-math.Sqrt(-math.Log(1.0-f0)), 2.0))
		}

	// ZDT6, Deb 2001, p360
	case 6:
		pname = "ZDT6"
		n := 10
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, grp int) {
			w := 6.0 * math.Pi
			f[0] = 1.0 - math.Exp(-4.0*x[0])*math.Pow(math.Sin(w*x[0]), 6.0)
			sum := 0.0
			for i := 1; i < n; i++ {
				sum += x[i]
			}
			w = float64(n - 1)
			c0 := 1.0 + w*math.Pow(sum/w, 0.25)
			c1 := 1.0 - math.Pow(f[0]/c0, 2.0)
			f[1] = c0 * c1
		}
		fmin = []float64{0, 0}
		fmax = []float64{1, 1}
		f1f0 = func(f0 float64) float64 {
			return 1.0 - math.Pow(f0, 2.0)
		}

	default:
		chk.Panic("problem %d is not available", opt.Problem)
	}

	// initialise optimiser
	opt.Init(GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.Solve()

	// plot
	PlotOvaOvaPareto(io.Sf("fig_flt05_%s", pname), &opt, nil, 0, 1, func() {
		np := 101
		F0 := utl.LinSpace(fmin[0], fmax[0], np)
		F1 := make([]float64, np)
		for i := 0; i < np; i++ {
			F1[i] = f1f0(F0[i])
		}
		plt.Plot(F0, F1, io.Sf("'b-', label='%s'", pname))
	}, nil, false)
}

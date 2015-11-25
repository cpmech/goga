// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_flt01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt01. quadratic function with inequalities")

	// parameters
	var opt Optimiser
	opt.Default()
	opt.Ngrp = 1
	opt.Nsol = 12
	opt.DEpc = 1
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

	// plot
	test_contour_plot_run_plot("fig_flt01", &opt)
}

func Test_flt02(tst *testing.T) {

	verbose()
	chk.PrintTitle("flt02. circle with equality constraint")

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

	// plot
	test_contour_plot_run_plot("fig_flt02", &opt)
}

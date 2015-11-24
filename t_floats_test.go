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
)

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
	opt.Read("ga-data.json")
	//opt.Tf = 2
	//opt.DtMig = 1
	opt.Verbose = false
	//opt.Seed = 1234
	opt.FltMin = []float64{-1, -1}
	opt.FltMax = []float64{3, 3}
	//opt.IntMin = []int{1, 1}
	//opt.IntMax = []int{10, 10}
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
	io.Pf("%+v\n", opt)

	// plot initial solution
	if chk.Verbose {
		plt.SetForEps(0.8, 455)
		opt.PlotContour(0, 0, 1, ContourParams{})
		opt.PlotSolutions(0, 1, plt.Fmt{M: "o", C: "k", Ls: "none", Ms: 3}, false)
	}

	// solve
	opt.Solve()

	// plot final solution
	if chk.Verbose {
		opt.PlotSolutions(0, 1, plt.Fmt{M: "o", C: "k", Ls: "none", Ms: 6}, true)
		plt.Equal()
		plt.Gll("$x_0$", "$x_1$", "")
		plt.SaveD("/tmp/goga", "fig_flt02.eps")
	}
}

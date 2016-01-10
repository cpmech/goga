// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func solve_problem(problem int) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 20
	opt.Ncpu = 1
	opt.Tf = 40
	opt.Verbose = false
	opt.Problem = problem
	opt.EpsH = 1e-3

	// problem variables
	var ng, nh int
	var fcn goga.MinProb_t
	var cprms goga.ContourParams
	cprms.Npts = 201
	eps_prop := 0.8
	eps_size := 400.0

	// problems
	switch opt.Problem {

	// problem # 1: quadratic function with inequalities
	case 1:
		opt.FltMin = []float64{-2, -2}
		opt.FltMax = []float64{2, 2}
		ng, nh = 5, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = x[0]*x[0]/2.0 + x[1]*x[1] - x[0]*x[1] - 2.0*x[0] - 6.0*x[1]
			g[0] = 2.0 - x[0] - x[1]     // ≥ 0
			g[1] = 2.0 + x[0] - 2.0*x[1] // ≥ 0
			g[2] = 3.0 - 2.0*x[0] - x[1] // ≥ 0
			g[3] = x[0]                  // ≥ 0
			g[4] = x[1]                  // ≥ 0
		}

	// problem # 2: circle with equality constraint
	case 2:
		xe := 1.0                      // centre of circle
		le := -0.4                     // selected level of f(x)
		ys := xe - (1.0+le)/math.Sqrt2 // coordinates of minimum point with level=le
		y0 := 2.0*ys + xe              // vertical axis intersect of straight line defined by c(x)
		xc := []float64{xe, xe}        // centre
		opt.FltMin = []float64{-1, -1}
		opt.FltMax = []float64{3, 3}
		ng, nh = 0, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			res := 0.0
			for i := 0; i < len(x); i++ {
				res += (x[i] - xc[i]) * (x[i] - xc[i])
			}
			f[0] = math.Sqrt(res) - 1.0
			h[0] = x[0] + x[1] + xe - y0
		}

	// problem # 3: Deb (2000) narrow crescent-shaped region
	case 3:
		opt.FltMin = []float64{0, 0}
		opt.FltMax = []float64{6, 6}
		ng, nh = 2, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = math.Pow(x[0]*x[0]+x[1]-11.0, 2.0) + math.Pow(x[0]+x[1]*x[1]-7.0, 2.0)
			g[0] = 4.84 - math.Pow(x[0]-0.05, 2.0) - math.Pow(x[1]-2.5, 2.0) // ≥ 0
			g[1] = x[0]*x[0] + math.Pow(x[1]-2.5, 2.0) - 4.84                // ≥ 0
		}
		cprms.Args = "levels=[1, 10, 25, 50, 100, 200, 400, 600, 1000, 1500]"
		cprms.Csimple = true
		cprms.Lwg = 0.6
		eps_prop = 1
	}

	// initialise optimiser
	nf := 1
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// initial solutions
	sols0 := opt.GetSolutionsCopy()

	// solve
	opt.Solve()

	// results
	goga.SortByOva(opt.Solutions, 0)
	best := opt.Solutions[0]
	io.Pforan("X_best = %v\n", best.Flt)
	io.Pforan("F_best = %v\n", best.Ova[0])
	io.Pforan("Oor    = %v\n", best.Oor)

	// text box
	extra := func() {
		str := io.Sf("$f(%.6f,%.6f)=%.6f$", best.Flt[0], best.Flt[1], best.Ova[0])
		plt.Text(0.5, 0.03, str, "size=12, transform=gca().transAxes, ha='center', zorder=1000, bbox=dict(boxstyle='round', facecolor='wheat')")
	}

	// plot
	plt.SetForEps(eps_prop, eps_size)
	goga.PlotFltFltContour(io.Sf("simpleoptm%d", opt.Problem), opt, sols0, 0, 1, 0, cprms, extra, true)
	return opt
}

func main() {
	P := []int{1, 2, 3}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem)
	}
	goga.TexPrmsReport("/tmp/goga", "simpleoptm", opts, 5)
	io.Pf("\n%v\n", opts[0].LogParams())
}

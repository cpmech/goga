// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// input filename
	_, fnkey := io.ArgToFilename(0, "frame2d", ".sim", true)

	// simple problems
	var opts []*goga.Optimiser
	if fnkey == "simple" {
		io.Pf("\n\n\n")
		//P := []int{1}
		P := utl.IntRange2(1, 19)
		opts = make([]*goga.Optimiser, len(P))
		for i, problem := range P {
			opts[i] = solve_problem(fnkey, problem)
		}
	} else {
		opts = []*goga.Optimiser{solve_problem(fnkey, 0)}
	}
	if opts[0].PlotSet1 {
		return
	}

	if opts[0].Nsamples > 1 {
		io.Pf("\n")
		rpt := goga.NewTexReport(opts)
		rpt.ShowDEC = false
		rpt.Type = 4
		rpt.TextSize = ""
		rpt.Title = "FORM Reliability: " + fnkey
		rpt.Fnkey = "rel-" + fnkey
		rpt.Generate()
	}
}

func solve_problem(fnkey string, problem int) (opt *goga.Optimiser) {

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()

	// options for report
	opt.RptFmtF = "%.4f"
	opt.RptFmtX = "%.3f"
	opt.RptFmtFdev = "%.1e"
	opt.RptWordF = "\\beta"
	opt.HistFmt = "%.2f"
	opt.HistNdig = 3
	opt.HistDelFmin = 0.005
	opt.HistDelFmax = 0.005

	// FORM data
	var lsft LSF_T
	var vars rnd.Variables

	// simple problem or FEM sim
	if fnkey == "simple" {
		opt.Read("ga-simple.json")
		opt.ProbNum = problem
		lsft, vars = get_simple_data(opt)
		fnkey += io.Sf("-%d", opt.ProbNum)
		io.Pf("\n----------------------------------- simple problem %d --------------------------------\n", opt.ProbNum)
	} else {
		opt.Read("ga-" + fnkey + ".json")
		lsft, vars = get_femsim_data(opt, fnkey)
		io.Pf("\n----------------------------------- femsim %s --------------------------------\n", fnkey)
	}

	// set limits
	nx := len(vars)
	opt.FltMin = make([]float64, nx)
	opt.FltMax = make([]float64, nx)
	for i, dat := range vars {
		opt.FltMin[i] = dat.Min
		opt.FltMax[i] = dat.Max
	}

	// log input
	var buf bytes.Buffer
	io.Ff(&buf, "%s", opt.LogParams())
	io.WriteFileVD("/tmp/gosl", fnkey+".log", &buf)

	// initialise distributions
	err := vars.Init()
	if err != nil {
		chk.Panic("cannot initialise distributions:\n%v", err)
	}

	// plot distributions
	if opt.PlotSet1 {
		io.Pf(". . . . . . . .  plot distributions  . . . . . . . .\n")
		np := 201
		for i, dat := range vars {
			plt.SetForEps(0.75, 250)
			dat.PlotPdf(np, "'b-',lw=2,zorder=1000")
			//plt.AxisXrange(dat.Min, dat.Max)
			plt.SetXnticks(15)
			plt.SaveD("/tmp/sims", io.Sf("distr-%s-%d.eps", fnkey, i))
		}
		return
	}

	// objective function
	nf := 1
	var ng, nh int
	var fcn goga.MinProb_t
	var obj goga.ObjFunc_t
	switch opt.Strategy {

	// argmin_x{ β(y(x)) | lsf(x) ≤ 0 }
	//  f ← sqrt(y dot y)
	//  g ← -lsf(x) ≥ 0
	//  h ← out-of-range in case Transform fails
	case 0:
		ng, nh = 1, 1
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {

			// original and normalised variables
			h[0] = 0
			y, invalid := vars.Transform(x)
			if invalid {
				h[0] = 1
				return
			}

			// objective value
			f[0] = math.Sqrt(la.VecDot(y, y)) // β

			// inequality constraint
			lsf, failed := lsft(x, cpu)
			g[0] = -lsf
			h[0] = failed
		}

	// argmin_x{ β(y(x)) | lsf(x) = 0 }
	//  f  ← sqrt(y dot y)
	//  h0 ← lsf(x)
	//  h1 ← out-of-range in case Transform fails
	case 1:
		ng, nh = 0, 2
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {

			// original and normalised variables
			h[0], h[1] = 0, 0
			y, invalid := vars.Transform(x)
			if invalid {
				h[0], h[1] = 1, 1
				return
			}

			// objective value
			f[0] = math.Sqrt(la.VecDot(y, y)) // β

			// equality constraint
			lsf, failed := lsft(x, cpu)
			h[0] = lsf
			h[1] = failed

			// induce minmisation of h0
			//f[0] += math.Abs(lsf)
		}

	case 2:
		opt.Nova = 1
		opt.Noor = 2
		obj = func(sol *goga.Solution, cpu int) {

			// clear out-of-range values
			sol.Oor[0] = 0 // invalid transformation or FEM failed
			sol.Oor[1] = 0 // g(x) ≤ 0 was violated

			// original and normalised variables
			x := sol.Flt
			y, invalid := vars.Transform(x)
			if invalid {
				sol.Oor[0] = goga.INF
				sol.Oor[1] = goga.INF
				return
			}

			// objective value
			sol.Ova[0] = math.Sqrt(la.VecDot(y, y)) // β

			// inequality constraint
			lsf, failed := lsft(x, cpu)
			sol.Oor[0] = failed
			sol.Oor[1] = fun.Ramp(lsf)
		}

	default:
		chk.Panic("strategy %d is not available", opt.Strategy)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, obj, fcn, nf, ng, nh)

	// solve
	io.Pf(". . . . . . . .  running  . . . . . . . .\n")
	opt.RunMany("", "")
	goga.StatF(opt, 0, true)
	io.Pfblue2("Tsys = %v\n", opt.SysTimeAve)

	// check
	goga.CheckFront0(opt, true)

	// results
	sols := goga.GetFeasible(opt.Solutions)
	if len(sols) > 0 {
		goga.SortByOva(sols, 0)
		best := sols[0]
		io.Pforan("x    = %.6f\n", best.Flt)
		io.Pforan("xref = %.6f\n", opt.RptXref)
		io.Pforan("β = %v  (%v)\n", best.Ova[0], opt.RptFref[0])
	}
	return
}

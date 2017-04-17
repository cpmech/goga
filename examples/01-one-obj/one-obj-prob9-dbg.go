// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// options
	doCheck := false
	plotTime := false
	plotContour := true
	printResults := true

	// GA parameters
	opt := new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 50
	opt.Ncpu = 1
	opt.Tmax = 7000
	opt.EpsH = 1e-3
	opt.Verbose = true
	opt.GenType = "latin"
	//opt.GenType = "halton"
	//opt.GenType = "rnd"
	opt.NormFlt = false
	//opt.UseMesh = true
	//opt.Nbry = 3

	// define problem
	cf := 1.0
	opt.RptName = "9"
	opt.RptFref = []float64{0.0539498478}
	opt.RptXref = []float64{-1.717143, 1.595709, 1.827247, -0.7636413, -0.7636450}
	opt.FltMin = []float64{-2.3, -2.3, -3.2, -3.2, -3.2}
	opt.FltMax = []float64{+2.3, +2.3, +3.2, +3.2, +3.2}
	ng, nh := 0, 3
	fcn := func(f, g, h, x []float64, y []int, cpu int) {
		//f[0] = math.Exp(x[0] * x[1] * x[2] * x[3] * x[4])
		//f[0] = x[0] * x[1] * x[2] * x[3] * x[4]
		//f[0] = math.Exp(x[0] * x[1] * x[2] * x[3] * x[4] * cf)
		f[0] = math.Exp(x[0]*x[1]*x[2]*x[3]*x[4]) * cf
		h[0] = x[0]*x[0] + x[1]*x[1] + x[2]*x[2] + x[3]*x[3] + x[4]*x[4] - 10.0
		h[1] = x[1]*x[2] - 5.0*x[3]*x[4]
		h[2] = math.Pow(x[0], 3.0) + math.Pow(x[1], 3.0) + 1.0
	}

	// check
	if doCheck {
		f := make([]float64, 1)
		h := make([]float64, 3)
		fcn(f, nil, h, opt.RptXref, nil, 0)
		io.Pforan("f(xref)  = %g  (%g)\n", f[0], opt.RptFref[0])
		io.Pforan("h0(xref) = %g\n", h[0])
		io.Pforan("h1(xref) = %g\n", h[1])
		io.Pforan("h2(xref) = %g\n", h[2])
	}

	// initialise optimiser
	nf := 1
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// output function
	T := make([]float64, opt.Tmax+1)                    // [nT]
	X := utl.Deep3alloc(opt.Nflt, opt.Nsol, opt.Tmax+1) // [nx][nsol][nT]
	F := utl.Deep3alloc(opt.Nova, opt.Nsol, opt.Tmax+1) // [nf][nsol][nT]
	U := utl.Deep3alloc(opt.Noor, opt.Nsol, opt.Tmax+1) // [nu][nsol][nT]
	opt.Output = func(time int, sols []*goga.Solution) {
		T[time] = float64(time)
		for j, s := range sols {
			for i := 0; i < opt.Nflt; i++ {
				X[i][j][time] = s.Flt[i]
			}
			for i := 0; i < opt.Nova; i++ {
				F[i][j][time] = s.Ova[i]
			}
			for i := 0; i < opt.Noor; i++ {
				U[i][j][time] = s.Oor[i]
			}
		}
	}

	// initial population
	fnk := "one-obj-prob9-dbg"
	sols0 := opt.GetSolutionsCopy()
	goga.WriteAllValues("/tmp/goga", fnk, opt)

	// solve
	opt.Solve()

	// best
	goga.SortSolutions(opt.Solutions, 0)
	best := goga.NewSolution(0, 1, &opt.Parameters)
	opt.Solutions[0].CopyInto(best)

	// print
	if printResults {
		//io.Pf("%13s%13s%13s%13s%10s\n", "f0", "u0", "u1", "u2", "feasible")
		io.Pf("%11s%11s%11s%11s%11s%11s%6s\n", "f0", "x0", "x1", "x2", "x3", "x4", "feas.")
		for i := 0; i < 10; i++ {
			s := opt.Solutions[i]
			//io.Pf("%13.5e%13.5e%13.5e%13.5e%10v\n", s.Ova[0], s.Oor[0], s.Oor[1], s.Oor[2], s.Feasible())
			io.Pf("%11.6f%11.6f%11.6f%11.6f%11.6f%11.6f%6v\n", s.Ova[0], s.Flt[0], s.Flt[1], s.Flt[2], s.Flt[3], s.Flt[4], s.Feasible())
		}
	}

	// plot: time series
	if plotTime {
		a, b := 100, len(T)
		//a, b := 300, 400
		plt.SetForEps(2.0, 400)
		nrow := opt.Nflt + opt.Nova + opt.Noor
		for j := 0; j < opt.Nsol; j++ {
			for i := 0; i < opt.Nflt; i++ {
				plt.Subplot(nrow, 1, 1+i)
				plt.Plot(T[a:b], X[i][j][a:b], "")
				plt.Gll("$t$", io.Sf("$x_%d$", i), "")
			}
		}
		for j := 0; j < opt.Nsol; j++ {
			for i := 0; i < opt.Nova; i++ {
				plt.Subplot(nrow, 1, 1+opt.Nflt+i)
				plt.Plot(T[a:b], F[i][j][a:b], "")
				plt.Gll("$t$", io.Sf("$f_%d$", i), "")
			}
		}
		for j := 0; j < opt.Nsol; j++ {
			for i := 0; i < opt.Noor; i++ {
				plt.Subplot(nrow, 1, 1+opt.Nflt+opt.Nova+i)
				plt.Plot(T[a:b], U[i][j][a:b], "")
				plt.Gll("$t$", io.Sf("$u_%d$", i), "")
			}
		}
		plt.SaveD("/tmp/goga", fnk+"-time.eps")
	}

	// plot: x-relationships
	if plotContour {
		pp := goga.NewPlotParams(false)
		pp.DirOut = "/tmp/goga"
		pp.FnKey = "contours-one-obj"
		pp.ContourAt = "best"
		pp.Npts = 101
		pp.Extra = func() {
			s := best
			plt.Subplot(4, 4, 14)
			plt.Text(0, 0, io.Sf("$f=%.8f$\n$x_0=%+.8f$\n$x_1=%+.8f$\n$x_2=%+.8f$\n$x_3=%+.8f$\n$x_4=%+.8f$",
				s.Ova[0], s.Flt[0], s.Flt[1], s.Flt[2], s.Flt[3], s.Flt[4]),
				"fontsize=14, ha='center', va='center'")
			plt.AxisRange(-1, 1, -1, 1)
			plt.AxisOff()
		}
		plt.SetForEps(1, 800)
		opt.PlotFltFltContour(sols0, -1, -1, 0, pp)
	}
}

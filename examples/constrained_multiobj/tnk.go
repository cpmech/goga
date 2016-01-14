// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	PI = math.Pi
)

func main() {

	// parameters
	opt := new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 50
	opt.Ncpu = 1
	opt.Tf = 500
	opt.Verbose = false
	opt.GenType = "latin"

	// options for report
	opt.RptFmtE = "%.5e"
	opt.RptFmtL = "%.6f"
	opt.RptFmtEdev = "%.4e"
	opt.RptFmtLdev = "%.4e"

	// problem: TNK
	opt.FltMin = []float64{0, 0}
	opt.FltMax = []float64{PI, PI}
	nf, ng, nh := 2, 2, 0
	fcn := func(f, g, h, x []float64, ξ []int, cpu int) {
		f[0] = x[0]
		f[1] = x[1]
		g[0] = x[0]*x[0] + x[1]*x[1] - 1.0 - 0.1*math.Cos(16.0*math.Atan2(x[0], x[1]))
		g[1] = 0.5 - math.Pow(x[0]-0.5, 2.0) - math.Pow(x[1]-0.5, 2.0)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// initial solutions
	var sols0 []*goga.Solution
	if false {
		sols0 = opt.GetSolutionsCopy()
	}

	// solve
	opt.RunMany("", "")

	// check
	var failed bool
	for _, sol := range opt.Solutions {
		for _, oor := range sol.Oor {
			if oor > 0 {
				failed = true
			}
		}
	}
	if failed {
		io.PfRed("failed\n")
	} else {
		io.PfGreen("OK\n")
	}

	// plot
	if true {
		feasibleOnly := false
		plt.SetForEps(0.8, 300)
		fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
		fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
		goga.PlotOvaOvaPareto(opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
		np := 201
		X, Y := utl.MeshGrid2D(0, 1.3, 0, 1.3, np, np)
		Z1 := utl.DblsAlloc(np, np)
		Z2 := utl.DblsAlloc(np, np)
		for j := 0; j < np; j++ {
			for i := 0; i < np; i++ {
				Z1[i][j] = X[i][j]*X[i][j] + Y[i][j]*Y[i][j] - 1.0 - 0.1*math.Cos(16.0*math.Atan2(Y[i][j], X[i][j]))
				Z2[i][j] = 0.5 - math.Pow(X[i][j]-0.5, 2.0) - math.Pow(Y[i][j]-0.5, 2.0)
			}
		}
		plt.ContourSimple(X, Y, Z1, false, 7, "linestyles=['--'], linewidths=[0.7], colors=['g'], levels=[0], clip_on=0")
		plt.ContourSimple(X, Y, Z2, false, 7, "linestyles=['-'],  linewidths=[0.7], colors=['b'], levels=[0], clip_on=0")
		plt.Equal()
		plt.AxisRange(0, 1.2, 0, 1.21)
		plt.SaveD("/tmp/goga", "tnk.eps")
	}
}

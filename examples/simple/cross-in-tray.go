// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	. "math"
	"time"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

// objective function
func fcn(f, g, h, x []float64, y []int, cpu int) {
	f[0] = -0.0001 * Pow(Abs(Sin(x[0])*Sin(x[1])*Exp(Abs(100-Sqrt(Pow(x[0], 2)+Pow(x[1], 2))/Pi)))+1, 0.1)
}

// main function
func main() {

	// problem definition
	nf := 1 // number of objective functions
	ng := 0 // number of inequality constraints
	nh := 0 // number of equality constraints

	// the solver (optimiser)
	var opt goga.Optimiser
	opt.Default()                    // must call this to set default constants
	opt.FltMin = []float64{-10, -10} // must set minimum
	opt.FltMax = []float64{+10, +10} // must set maximum
	opt.Nsol = 80

	// initialise the solver
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve problem
	tstart := time.Now()
	opt.Solve()
	cputime := time.Now().Sub(tstart)

	// print results
	fvec := []float64{0} // temporary vector to use with fcn
	xbest := opt.Solutions[0].Flt
	fcn(fvec, nil, nil, xbest, nil, 0)
	io.Pf("xBest    = %v\n", xbest)
	io.Pf("f(xBest) = %v\n", fvec)

	// plotting
	pp := goga.NewPlotParams(false)
	pp.Npts = 101
	pp.ArgsF.NoLines = true
	pp.ArgsF.CmapIdx = 4
	plt.Reset(true, nil)
	plt.Title(io.Sf("Nsol(pop.size)=%d  Tmax(generations)=%d CpuTime=%v", opt.Nsol, opt.Tmax, io.RoundDuration(cputime, 1e3)), &plt.A{Fsz: 8})
	opt.PlotContour(0, 1, 0, pp)
	plt.PlotOne(xbest[0], xbest[1], &plt.A{C: "r", Mec: "r", M: "*"})
	plt.Gll("$x_0$", "$x_1$", nil)
	plt.Save("/tmp/goga", "cross-in-tray")
}

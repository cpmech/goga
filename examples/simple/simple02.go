// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	. "math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// objective function
func fcn(f, g, h, x []float64, y []int, cpu int) {
	f[0] = Cos(8*x[0]*Pi) * Exp(Cos(x[0]*Pi)-1)
}

// main function
func main() {

	// problem definition
	nf := 1 // number of objective functions
	ng := 0 // number of inequality constraints
	nh := 0 // number of equality constraints

	// the solver (optimiser)
	var opt goga.Optimiser
	opt.Default()             // must call this to set default constants
	opt.FltMin = []float64{0} // must set minimum
	opt.FltMax = []float64{2} // must set maximum

	// initialise the solver
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve problem
	opt.Solve()

	// auxiliary
	fvec := []float64{0} // temporary vector to use with fcn
	xvec := []float64{0} // temporary vector to use with fcn

	// print results
	xBest := opt.Solutions[0].Flt[0]
	xvec[0] = xBest
	fcn(fvec, nil, nil, xvec, nil, 0)
	fBest := fvec[0]
	io.Pf("xBest    = %v\n", xBest)
	io.Pf("f(xBest) = %v\n", fBest)

	// generate f(x) curve
	X := utl.LinSpace(opt.FltMin[0], opt.FltMax[0], 1001)
	F := utl.GetMapped(X, func(x float64) float64 {
		xvec[0] = x
		fcn(fvec, nil, nil, xvec, nil, 0)
		return fvec[0]
	})

	// plotting
	plt.Reset(true, nil)
	plt.PlotOne(xBest, fBest, &plt.A{L: "best", C: "#bf2e64", M: ".", Ms: 15, NoClip: true})
	opt.PlotAddFltOva(0, 0, opt.Solutions, 1, &plt.A{L: "all", C: "g", Ls: "none", M: ".", NoClip: true})
	plt.Plot(X, F, &plt.A{L: "f(x)", C: "#0077d2"})
	plt.Gll("$x$", "$f$", nil)
	plt.Save("/tmp/goga", "simple02")
}

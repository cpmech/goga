// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// goga can use two types of objective functions:
//
//   (A) the higher-level one: MinProb_t which takes the vector of random variables x and returns
//       the objectives values in f. It may also return the inequality constraints in g and the
//       equality constraints in h. It also accepts integer random variables in y
//
//   (B) the lower-level one: ObjFunc_t which takes the pointer to a candidate solution object
//       (Solution) and fills the Ova array in this object with the objective values. In this method
//       the vector of random variables x is stored as Flt
//
// Both functions take the cpu number as input if that's necessary (rarely)
//
// The functions definitions of each case are shown below
//
// ObjFunc_t defines the objective fuction
//   type ObjFunc_t func(sol *Solution, cpu int)
//
// MinProb_t defines objective functon for specialised minimisation problem
//   type MinProb_t func(f, g, h, x []float64, y []int, cpu int)

// case A: finding the minimum of 2.0 + (1+x)²
func fcnA(f, g, h, x []float64, y []int, cpu int) {
	f[0] = 2.0 + (1.0+x[0])*(1.0+x[0])
}

// case B: finding the minimum of
func fcnB(sol *goga.Solution, cpu int) {
	x := sol.Flt
	sol.Ova[0] = 2.0 + (1.0+x[0])*(1.0+x[0])
}

// main function
func main() {

	// problem definition
	nf := 1 // number of objective functions
	ng := 0 // number of inequality constraints
	nh := 0 // number of equality constraints

	// the solver (optimiser)
	var opt goga.Optimiser
	opt.Default()              // must call this to set default constants
	opt.FltMin = []float64{-2} // must set minimum
	opt.FltMax = []float64{2}  // must set maximum

	// initialise the solver
	useMethodA := false
	if useMethodA {
		opt.Init(goga.GenTrialSolutions, nil, fcnA, nf, ng, nh)
	} else {
		opt.Init(goga.GenTrialSolutions, fcnB, nil, nf, ng, nh)
	}

	// solve problem
	opt.Solve()

	// print results
	xBest := opt.Solutions[0].Flt[0]
	fBest := 2.0 + (1.0+xBest)*(1.0+xBest)
	io.Pf("xBest    = %v\n", xBest)
	io.Pf("f(xBest) = %v\n", fBest)

	// plotting
	fvec := []float64{0} // temporary vector to use with fcnA
	xvec := []float64{0} // temporary vector to use with fcnA
	X := utl.LinSpace(-2, 2, 101)
	F := utl.GetMapped(X, func(x float64) float64 {
		xvec[0] = x
		fcnA(fvec, nil, nil, xvec, nil, 0)
		return fvec[0]
	})
	plt.Reset(true, nil)
	plt.PlotOne(xBest, fBest, &plt.A{C: "r", M: "o", Ms: 20, NoClip: true})
	plt.Plot(X, F, nil)
	plt.Gll("$x$", "$f$", nil)
	plt.Save("/tmp/goga", "simple01")
}

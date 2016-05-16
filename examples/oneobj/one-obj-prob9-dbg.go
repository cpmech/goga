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

	// GA parameters
	opt := new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 6
	opt.Ncpu = 1
	opt.Tf = 10
	opt.EpsH = 1e-3
	opt.Verbose = true
	opt.GenType = "latin"
	//opt.GenType = "halton"
	//opt.GenType = "rnd"
	opt.NormFlt = false
	opt.UseMesh = true
	opt.Nbry = 3

	// define problem
	opt.RptName = "9"
	opt.RptFref = []float64{0.0539498478}
	opt.RptXref = []float64{-1.717143, 1.595709, 1.827247, -0.7636413, -0.7636450}
	opt.FltMin = []float64{-2.3, -2.3, -3.2, -3.2, -3.2}
	opt.FltMax = []float64{+2.3, +2.3, +3.2, +3.2, +3.2}
	ng, nh := 0, 3
	fcn := func(f, g, h, x []float64, ξ []int, cpu int) {
		f[0] = math.Exp(x[0] * x[1] * x[2] * x[3] * x[4])
		h[0] = x[0]*x[0] + x[1]*x[1] + x[2]*x[2] + x[3]*x[3] + x[4]*x[4] - 10.0
		h[1] = x[1]*x[2] - 5.0*x[3]*x[4]
		h[2] = math.Pow(x[0], 3.0) + math.Pow(x[1], 3.0) + 1.0
	}

	// check
	if false {
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
	T := make([]float64, opt.Tf+1)                    // [nT]
	X := utl.Deep3alloc(opt.Nflt, opt.Nsol, opt.Tf+1) // [nx][nsol][nT]
	F := utl.Deep3alloc(opt.Nova, opt.Nsol, opt.Tf+1) // [nf][nsol][nT]
	U := utl.Deep3alloc(opt.Noor, opt.Nsol, opt.Tf+1) // [nu][nsol][nT]
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
	//S0 := opt.GetSolutionsCopy()
	goga.WriteAllValues("/tmp/goga", fnk, opt)

	// solve
	opt.Solve()

	// print
	if false {
		io.Pf("%13s%13s%13s%13s%10s\n", "f0", "u0", "u1", "u2", "feasible")
		for _, s := range opt.Solutions {
			io.Pf("%13.5e%13.5e%13.5e%13.5e%10v\n", s.Ova[0], s.Oor[0], s.Oor[1], s.Oor[2], s.Feasible())
		}
	}

	// plot: time series
	//a, b := 100, len(T)
	a, b := 0, 1 //len(T)
	if false {
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
	if true {
		plt.SetForEps(1, 700)
		ncol := opt.Nflt - 1
		for i := 0; i < opt.Nflt-1; i++ {
			for j := i + 1; j < opt.Nflt; j++ {
				plt.Subplot(ncol, ncol, i*ncol+j)
				if opt.UseMesh {
					opt.Meshes[i][j].CalcDerived(0)
					opt.Meshes[i][j].Draw2d(false, false, nil, 0)
				}
				for k := 0; k < opt.Nsol; k++ {
					plt.Plot(X[i][k][a:b], X[j][k][a:b], "ls='none', marker='.'")
				}
				plt.Gll(io.Sf("$x_%d$", i), io.Sf("$x_%d$", j), "")
			}
		}
		plt.SaveD("/tmp/goga", fnk+"-x.eps")
	}
}

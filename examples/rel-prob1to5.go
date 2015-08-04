// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"math"
	"time"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// Problems where all variables are standard variables => μ=0 and σ=1 => y = x
	// References
	//  [1] Santos SR, Matioli LC and Beck AT. New optimization algorithms for structural
	//      reliability analysis. Computer Modeling in Engineering & Sciences, 83(1):23-56; 2012
	//      doi:10.3970/cmes.2012.083.023
	//  [2] Borri A and Speranzini E. Structural reliability analysis using a standard deterministic
	//      finite element code. Structural Safety, 19(4):361-382; 1997
	//      doi:10.1016/S0167-4730(97)00017-9
	//  [3] Grooteman F.  Adaptive radial-based importance sampling method or structural
	//      reliability. Structural safety, 30:533-542; 2008
	//      doi:10.1016/j.strusafe.2007.10.002
	//  [4] Wang L and Grandhi RV. Higher-order failure probability calculation using nonlinear
	//      approximations. Computer Methods in Applied Mechanics and Engineering, 168(1-4):185-206;
	//      1999 doi:10.1016/S0045-7825(98)00140-6

	// catch errors
	defer func() {
		if err := recover(); err != nil {
			io.PfRed("ERROR: %v\n", err)
		}
	}()

	// read parameters
	fn := "rel-prob1to6"
	fn, _ = io.ArgToFilename(0, fn, ".json", true)
	C := goga.ReadConfParams(fn)
	io.Pf("\n%s\nproblem # %v\n", utl.PrintThickLine(80), C.Problem)

	// initialise random numbers generator
	rnd.Init(C.Seed)

	// problems's data: limit state function
	npts := 41
	var g func(x []float64) float64
	var βref float64
	var xref, xmin, xmax []float64
	switch C.Problem {

	// problem # 1 of [1] and Eq. (A.5) of [2]
	case 1:
		g = func(x []float64) float64 {
			return 0.1*math.Pow(x[0]-x[1], 2.0) - (x[0]+x[1])/math.Sqrt2 + 2.5
		}
		βref = 2.5 // from [1]
		xref = []float64{1.7677, 1.7677}
		xmin, xmax = []float64{-5, -5}, []float64{5, 5}

	// problem # 2 of [1] and Eq. (A.6) of [2]
	case 2:
		g = func(x []float64) float64 {
			return -0.5*math.Pow(x[0]-x[1], 2.0) - (x[0]+x[1])/math.Sqrt2 + 3.0
		}
		βref = 1.658 // from [2]
		xref = []float64{-0.7583, 1.4752}
		xmin, xmax = []float64{-5, -5}, []float64{5, 5}

	// problem # 3 from [1] and # 6 from [3]
	case 3:
		g = func(x []float64) float64 {
			return 2.0 - x[1] - 0.1*math.Pow(x[0], 2) + 0.06*math.Pow(x[0], 3)
		}
		βref = 2.0 // from [1]
		xref = []float64{0, 2}
		xmin, xmax = []float64{-5, -5}, []float64{5, 5}

	// problem # 4 from [1] and # 8 from [3]
	case 4:
		g = func(x []float64) float64 {
			return 3.0 - x[1] + 256.0*math.Pow(x[0], 4.0)
		}
		npts = 101
		βref = 3.0 // from [1]
		xref = []float64{0, 3}
		xmin, xmax = []float64{-5, -5}, []float64{5, 5}

	// problem # 5 from [1] and # 1 from [4] (modified)
	case 5:
		shift := 0.1
		g = func(x []float64) float64 {
			return 1.0 + math.Pow(x[0]+x[1]+shift, 2.0)/4.0 - 4.0*math.Pow(x[0]-x[1]+shift, 2.0)
		}
		βref = 0.3536 // from [1]
		xref = []float64{-βref * math.Sqrt2 / 2.0, βref * math.Sqrt2 / 2.0}
		xmin, xmax = []float64{-1, -1}, []float64{1, 1}

	default:
		chk.Panic("problem number %d is invalid", C.Problem)
	}

	// objective value function
	ovfunc := func(ind *goga.Individual, idIsland, t int, report *bytes.Buffer) (ova, oor float64) {
		x := []float64{ind.GetFloat(0), ind.GetFloat(1)} // must be inside ovfunc to avoid data race problems
		if C.Strategy == 1 {
			ova = la.VecDot(x, x)
			oor = utl.GtePenalty(0, g(x), 1)
		} else {
			fp := utl.GtePenalty(1e-2, math.Abs(g(x)), 1)
			ova = la.VecDot(x, x) + fp
			oor = fp
		}
		return
	}

	// evolver
	evo := goga.NewEvolverFloatChromo(C, xmin, xmax, ovfunc, goga.NewBingoFloats(xmin, xmax))

	// benchmarking
	cpu0 := time.Now()

	// for a number of trials
	betas := make([]float64, C.Ntrials)
	for i := 0; i < C.Ntrials; i++ {

		// reset population
		if i > 0 {
			for _, isl := range evo.Islands {
				isl.Pop.GenFloatRandom(C, xmin, xmax)
			}
		}
		pop0 := evo.Islands[0].Pop.GetCopy()

		// run
		check := i == C.Ntrials-1
		verbose := check
		doreport := check
		evo.Run(verbose, doreport)
		β := calc_beta(evo.Best, βref, xref, verbose)
		betas[i] = β

		// plot contour
		if check {
			if C.DoPlot {
				pop1 := evo.Islands[0].Pop
				goga.PlotTwoVarsContour("/tmp/goga", io.Sf("rel-prob%d", C.Problem), pop0, pop1, evo.Best,
					xmin, xmax, npts, true, func() { plt.SetXnticks(11); plt.SetYnticks(11) }, g, g)
			}
		}
	}

	// benchmarking
	io.Pfcyan("\nelapsed time = %v\n", time.Now().Sub(cpu0))

	// analysis
	βmin, βave, βmax, βdev := rnd.StatBasic(betas, true)
	io.Pf("\nβmin = %v\n", βmin)
	io.PfYel("βave = %v\n", βave)
	io.Pf("βmax = %v\n", βmax)
	io.Pf("βdev = %v\n\n", βdev)
	io.Pf(rnd.BuildTextHist(nice_num(βmin-0.005), nice_num(βmax+0.005), 11, betas, "%.3f", 60))
}

func nice_num(x float64) float64 {
	s := io.Sf("%.2f", x)
	return io.Atof(s)
}

func calc_beta(best *goga.Individual, βref float64, xref []float64, verbose bool) (β float64) {
	xs := make([]float64, best.Nfltgenes)
	for i := 0; i < best.Nfltgenes; i++ {
		xs[i] = best.GetFloat(i)
	}
	β = math.Sqrt(la.VecDot(xs, xs))
	if verbose {
		io.Pf("\nova  = %g  oor = %g\n", best.Ova, best.Oor)
		io.Pf("x    = %v\n", xs)
		io.Pf("xref = %v\n", xref)
		io.PfYel("β    = %g", β)
		io.Pf(" (%g)\n", βref)
	}
	return
}

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
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// catch errors
	defer func() {
		if err := recover(); err != nil {
			io.PfRed("ERROR: %v\n", err)
		}
	}()

	// Problem # 1:
	//  All variables are standard variables => μ=0 and σ=1 => y = x

	// read parameters
	fn := "rel-prob1"
	fn, fnkey := io.ArgToFilename(0, fn, ".json", true)
	C := goga.ReadConfParams(fn)

	// initialise random numbers generator
	rnd.Init(C.Seed)

	// limit state function
	g := func(x []float64) float64 {
		return 0.1*math.Pow(x[0]-x[1], 2.0) - (x[0]+x[1])/math.Sqrt2 + 2.5
	}

	// objective value function
	x := make([]float64, 2)
	ovfunc := func(ind *goga.Individual, idIsland, time int, report *bytes.Buffer) (ova, oor float64) {
		x[0], x[1] = ind.GetFloat(0), ind.GetFloat(1)
		fp := utl.GtePenalty(1e-2, math.Abs(g(x)), 1)
		ova = la.VecDot(x, x) + fp
		oor = fp
		return
	}

	// bingo
	ndim := 2
	vmin, vmax := -2.0, 2.0
	xmin, xmax := utl.DblVals(ndim, vmin), utl.DblVals(ndim, vmax)
	bingo := goga.NewBingoFloats(xmin, xmax)

	// evolver
	βref := 2.5
	evo := goga.NewEvolverFloatChromo(C, xmin, xmax, ovfunc, bingo)

	// benchmarking
	cpu0 := time.Now()

	// for a number of trials
	ntrials := 100
	betas := make([]float64, ntrials)
	for i := 0; i < ntrials; i++ {

		// reset population
		if i > 0 {
			for _, isl := range evo.Islands {
				isl.Pop.GenFloatRandom(C, xmin, xmax)
			}
		}
		pop0 := evo.Islands[0].Pop

		// run
		check := i == ntrials-1
		verbose := check
		doreport := check
		evo.Run(verbose, doreport)
		β := calc_beta(evo.Best, βref, verbose)
		betas[i] = β

		// plot contour
		if check {
			if C.DoPlot {
				xmin := []float64{-1, -1}
				xmax := []float64{5, 5}
				goga.PlotTwoVarsContour("/tmp/goga", fnkey, pop0, evo.Islands[0].Pop, evo.Best,
					xmin, xmax, 41, true, nil, g, g)
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

func calc_beta(best *goga.Individual, βref float64, verbose bool) (β float64) {
	xs := make([]float64, best.Nfltgenes)
	for i := 0; i < best.Nfltgenes; i++ {
		xs[i] = best.GetFloat(i)
	}
	β = math.Sqrt(la.VecDot(xs, xs))
	if verbose {
		io.Pf("\nova = %g  oor = %g\n", best.Ova, best.Oor)
		io.Pf("x   = %v\n", xs)
		io.PfYel("β   = %g", β)
		io.Pf(" (%g)\n", βref)
	}
	return
}

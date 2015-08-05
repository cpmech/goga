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

	// read parameters
	fn := "rel-nd-simple"
	fn, _ = io.ArgToFilename(0, fn, ".json", true)
	C := goga.ReadConfParams(fn)

	// initialise random numbers generator
	rnd.Init(C.Seed)

	// problem # 6 of [1] and case # 3 of [3]
	g := func(x []float64) float64 {
		sum := 0.0
		for i := 0; i < 9; i++ {
			sum += x[i] * x[i]
		}
		return 2.0 - 0.015*sum - x[9]
	}
	βref := 2.0 // from [1]
	xmin := utl.DblVals(10, -5)
	xmax := utl.DblVals(10, 5)

	// objective value function
	ovfunc := func(ind *goga.Individual, idIsland, t int, report *bytes.Buffer) (ova, oor float64) {
		x := ind.GetFloats()
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

		// run
		check := i == C.Ntrials-1
		verbose := check
		doreport := check
		evo.Run(verbose, doreport)
		β := calc_beta(evo.Best, βref, verbose)
		betas[i] = β
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
	xs := best.GetFloats()
	β = math.Sqrt(la.VecDot(xs, xs))
	if verbose {
		io.Pf("\nova  = %g  oor = %g\n", best.Ova, best.Oor)
		io.PfYel("β    = %g", β)
		io.Pf(" (%g)\n", βref)
	}
	return
}

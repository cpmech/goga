// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"math/rand"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_evo01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo01. organise sequence of ints")
	io.Pf("\n")

	// initialise random numbers generator
	rnd.Init(0) // 0 => use current time as seed

	// mutation function
	mtfunc := func(A []int, nchanges int, pm float64, extra interface{}) {
		size := len(A)
		if !rnd.FlipCoin(pm) || size < 1 {
			return
		}
		pos := rnd.IntGetUniqueN(0, size, nchanges)
		for _, i := range pos {
			if A[i] == 1 {
				A[i] = 0
			}
			if A[i] == 0 {
				A[i] = 1
			}
		}
	}

	// objective function
	ovfunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) (ova, oor float64) {
		score := 0.0
		count := 0
		for _, val := range ind.Ints {
			if val == 0 && count%2 == 0 {
				score += 1.0
			}
			if val == 1 && count%2 != 0 {
				score += 1.0
			}
			count++
		}
		ova = 1.0 / (1.0 + score)
		return
	}

	// reference individual
	nvals := 20
	ref := NewIndividual(1, utl.IntVals(nvals, 1))
	for i := 0; i < nvals; i++ {
		ref.Ints[i] = rand.Intn(2)
	}

	// bingo
	bingo := NewBingoInts(utl.IntVals(nvals, 0), utl.IntVals(nvals, 1))

	// parameters
	C := NewConfParams()
	C.Nisl = 2
	C.Ninds = 20
	C.FnKey = "" //"test_evo01"
	C.MtIntFunc = mtfunc
	C.RegTol = 0
	C.CalcDerived()

	// evolver
	evo := NewEvolver(C, ref, ovfunc, bingo)

	// run
	evo.Run(true, false)

	// results
	ideal := 1.0 / (1.0 + float64(nvals))
	io.PfGreen("\nBest = %v\nBestOV = %v  (ideal=%v)\n", evo.Best.Ints, evo.Best.Ova, ideal)
}

func Test_evo02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo02")

	// initialise random numbers generator
	//rnd.Init(0) // 0 => use current time as seed
	rnd.Init(1111) // 0 => use current time as seed

	f := func(x []float64) float64 { return x[0]*x[0]/2.0 + x[1]*x[1] - x[0]*x[1] - 2.0*x[0] - 6.0*x[1] }
	c1 := func(x []float64) float64 { return x[0] + x[1] - 2.0 }      // ≤ 0
	c2 := func(x []float64) float64 { return -x[0] + 2.0*x[1] - 2.0 } // ≤ 0
	c3 := func(x []float64) float64 { return 2.0*x[0] + x[1] - 3.0 }  // ≤ 0
	c4 := func(x []float64) float64 { return -x[0] }                  // ≤ 0
	c5 := func(x []float64) float64 { return -x[1] }                  // ≤ 0

	// objective function
	p := 1.0
	x := make([]float64, 2)
	ovfunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) (ova, oor float64) {
		x[0], x[1] = ind.GetFloat(0), ind.GetFloat(1)
		ova = f(x)
		oor += utl.GtePenalty(0, c1(x), p)
		oor += utl.GtePenalty(0, c2(x), p)
		oor += utl.GtePenalty(0, c3(x), p)
		oor += utl.GtePenalty(0, c4(x), p)
		oor += utl.GtePenalty(0, c5(x), p)
		return
	}

	// parameters
	C := NewConfParams()
	C.Pll = true
	C.Nisl = 4
	C.Ninds = 20
	if chk.Verbose {
		C.FnKey = "test_evo02"
		C.DoPlot = true
	}
	C.CalcDerived()

	// bingo
	ndim := 2
	vmin, vmax := -2.0, 2.0
	xmin, xmax := utl.DblVals(ndim, vmin), utl.DblVals(ndim, vmax)
	bingo := NewBingoFloats(xmin, xmax)

	// evolver
	evo := NewEvolverFloatChromo(C, xmin, xmax, ovfunc, bingo)
	verbose := true
	doreport := true
	pop0 := evo.Islands[0].Pop.GetCopy()
	evo.Run(verbose, doreport)

	// results
	io.PfGreen("\nx=%g (%g)\n", evo.Best.GetFloat(0), 2.0/3.0)
	io.PfGreen("y=%g (%g)\n", evo.Best.GetFloat(1), 4.0/3.0)
	io.PfGreen("BestOV=%g (%g)\n", evo.Best.Ova, f([]float64{2.0 / 3.0, 4.0 / 3.0}))

	// plot contour
	if C.DoPlot {
		PlotTwoVarsContour("/tmp/goga", "contour_evo02", pop0, evo.Islands[0].Pop, evo.Best,
			xmin, xmax, 41, true, nil, f, c1, c2, c3, c4, c5)
	}
}

func Test_evo03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo03")

	//rnd.Init(0)

	// geometry
	xe := 1.0                      // centre of circle
	le := -0.4                     // selected level of f(x)
	ys := xe - (1.0+le)/math.Sqrt2 // coordinates of minimum point with level=le
	y0 := 2.0*ys + xe              // vertical axis intersect of straight line defined by c(x)
	xc := []float64{xe, xe}        // centre
	nx := len(xc)
	f := func(x []float64) (res float64) {
		for i := 0; i < nx; i++ {
			res += (x[i] - xc[i]) * (x[i] - xc[i])
		}
		return math.Sqrt(res) - 1
	}
	c := func(x []float64) (res float64) {
		return x[0] + x[1] + xe - y0
	}

	// objective function
	p := 1.0
	x := make([]float64, nx)
	ovfunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) (ova, oor float64) {
		x[0], x[1] = ind.GetFloat(0), ind.GetFloat(1)
		fp := utl.GtePenalty(1e-2, math.Abs(c(x)), p)
		ova = f(x) + fp
		oor = fp
		return
	}

	// parameters
	C := NewConfParams()
	C.Pll = true
	C.Nisl = 4
	C.Ninds = 20
	if chk.Verbose {
		C.FnKey = "test_evo03"
		C.DoPlot = chk.Verbose
	}
	C.CalcDerived()

	// bingo
	ndim := 2
	vmin, vmax := -1.0, 3.0
	xmin, xmax := utl.DblVals(ndim, vmin), utl.DblVals(ndim, vmax)
	bingo := NewBingoFloats(xmin, xmax)

	// evolver
	evo := NewEvolverFloatChromo(C, xmin, xmax, ovfunc, bingo)
	pop0 := evo.Islands[0].Pop.GetCopy()
	verbose := true
	doreport := true
	evo.Run(verbose, doreport)

	// results
	xbest := []float64{evo.Best.GetFloat(0), evo.Best.GetFloat(1)}
	io.PfGreen("\nx=%g (%g)\n", xbest[0], ys)
	io.PfGreen("y=%g (%g)\n", xbest[1], ys)
	io.PfGreen("BestOV=%g (%g)\n\n", evo.Best.Ova, le)

	// plot contour
	if C.DoPlot {
		extra := func() {
			plt.PlotOne(ys, ys, "'o', markeredgecolor='yellow', markerfacecolor='none', markersize=10")
		}
		PlotTwoVarsContour("/tmp/goga", "contour_evo03", pop0, evo.Islands[0].Pop, evo.Best,
			xmin, xmax, 41, true, extra, f, c)
	}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_evo01(tst *testing.T) {

	verbose()
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
	ovfunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) (ov, oor float64) {
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
		ov = 1.0 / (1.0 + score)
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
	bingo.UseIntRnd = true

	// parameters
	C := NewConfParams()
	C.Nisl = 3
	C.Ninds = 6
	C.FnKey = "test_evo01"
	C.MtIntFunc = mtfunc

	// evolver
	evo := NewEvolver(C, ref, ovfunc, bingo)

	// run
	evo.Run(true)

	// results
	ideal := 1.0 / (1.0 + float64(nvals))
	io.PfGreen("\nBest = %v\nBestOV = %v  (ideal=%v)\n", evo.Best.Ints, evo.Best.Ova, ideal)

	// plot
	if C.DoPlot {
		for i, isl := range evo.Islands {
			first := i == 0
			last := i == C.Nisl-1
			isl.PlotOvs(".png", "", 0, C.Tf, true, "%.6f", first, last)
		}
	}
}

func Test_evo02(tst *testing.T) {

	verbose()
	chk.PrintTitle("evo02")

	// initialise random numbers generator
	rnd.Init(0) // 0 => use current time as seed

	f := func(x, y float64) float64 { return x*x/2.0 + y*y - x*y - 2.0*x - 6.0*y }
	c1 := func(x, y float64) float64 { return x + y - 2.0 }      // ≤ 0
	c2 := func(x, y float64) float64 { return -x + 2.0*y - 2.0 } // ≤ 0
	c3 := func(x, y float64) float64 { return 2.0*x + y - 3.0 }  // ≤ 0
	c4 := func(x, y float64) float64 { return -x }               // ≤ 0
	c5 := func(x, y float64) float64 { return -y }               // ≤ 0

	// objective function
	p := 1000.0
	ovfunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) (ov, oor float64) {
		x := ind.GetFloat(0)
		y := ind.GetFloat(1)
		ov = f(x, y)
		oor += utl.GtePenalty(0, c1(x, y), p)
		oor += utl.GtePenalty(0, c2(x, y), p)
		oor += utl.GtePenalty(0, c3(x, y), p)
		oor += utl.GtePenalty(0, c4(x, y), p)
		oor += utl.GtePenalty(0, c5(x, y), p)
		return
	}

	// parameters
	C := NewConfParams()
	C.Nisl = 1
	C.Ninds = 10
	C.FnKey = "test_evo02"
	C.DoPlot = true
	C.Noise = 0

	// bingo
	ndim := 2
	vmin, vmax := -2.0, 2.0
	xmin, xmax := utl.DblVals(ndim, vmin), utl.DblVals(ndim, vmax)
	bingo := NewBingoFloats(xmin, xmax)
	bingo.UseFltRnd = false

	// populations
	pops := make([]Population, C.Nisl)
	for i := 0; i < C.Nisl; i++ {
		pops[i] = NewPopFloatRandom(C, xmin, xmax)
	}

	// evolver
	evo := NewEvolverPop(C, pops, ovfunc, bingo)

	// plot contour
	if C.DoPlot {
		plt.SetForEps(0.8, 300)
		n, nn := 41, 7
		X, Y := utl.MeshGrid2D(vmin, vmax, vmin, vmax, n, n)
		Z := la.MatAlloc(n, n)
		xx, yy := utl.MeshGrid2D(vmin, vmax, vmin, vmax, nn, nn)
		z1 := la.MatAlloc(nn, nn)
		z2 := la.MatAlloc(nn, nn)
		z3 := la.MatAlloc(nn, nn)
		z4 := la.MatAlloc(nn, nn)
		z5 := la.MatAlloc(nn, nn)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				Z[i][j] = f(X[i][j], Y[i][j])
			}
		}
		for i := 0; i < nn; i++ {
			for j := 0; j < nn; j++ {
				z1[i][j] = c1(xx[i][j], yy[i][j])
				z2[i][j] = c2(xx[i][j], yy[i][j])
				z3[i][j] = c3(xx[i][j], yy[i][j])
				z4[i][j] = c4(xx[i][j], yy[i][j])
				z5[i][j] = c5(xx[i][j], yy[i][j])
			}
		}
		plt.Contour(X, Y, Z, "")
		plt.ContourSimple(xx, yy, z1, "levels=[0], colors=['yellow']")
		plt.ContourSimple(xx, yy, z2, "levels=[0], colors=['yellow']")
		plt.ContourSimple(xx, yy, z3, "levels=[0], colors=['yellow']")
		plt.ContourSimple(xx, yy, z4, "levels=[0], colors=['yellow'], linestyles=['--']")
		plt.ContourSimple(xx, yy, z5, "levels=[0], colors=['yellow'], linestyles=['--']")
		for _, ind := range evo.Islands[0].Pop {
			x := ind.GetFloat(0)
			y := ind.GetFloat(1)
			plt.PlotOne(x, y, "'k.', clip_on=0")
		}
	}

	plt.Equal()
	plt.SaveD("/tmp/goga", "test_evo02_contour.eps")
	return

	// run
	evo.Run(true)
	io.PfGreen("\nx=%g (%g)\n", evo.Best.GetFloat(0), 2.0/3.0)
	io.PfGreen("y=%g (%g)\n", evo.Best.GetFloat(1), 4.0/3.0)
	io.PfGreen("BestOV=%g (%g)\n", evo.Best.Ova, f(2.0/3.0, 4.0/3.0))

	// plot population on contour
	if C.DoPlot {
		for _, ind := range evo.Islands[0].Pop {
			x := ind.GetFloat(0)
			y := ind.GetFloat(1)
			plt.PlotOne(x, y, "'g*'")
		}
		x := evo.Best.GetFloat(0)
		y := evo.Best.GetFloat(1)
		plt.PlotOne(x, y, "'y*', ms=8")
		plt.Equal()
		plt.SaveD("/tmp/goga", "test_evo02_contour.eps")
	}

	// plot
	if C.DoPlot {
		evo.Islands[0].PlotOvs(".png", "", 10, C.Tf, true, "%.6f", true, true)
	}
}

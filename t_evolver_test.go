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
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_evo01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo01")

	// initialise random numbers generator
	rnd.Init(0) // 0 => use current time as seed

	// objective function
	ovfunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		ind.ObjValue = 1.0 / (1.0 + (ind.GetFloat(0)+ind.GetFloat(1)+ind.GetFloat(2))/3.0)
	}

	// reference population
	nbases := 8
	pop := NewPopFloatChromo(nbases, [][]float64{
		{11, 21, 31},
		{12, 22, 32},
		{13, 23, 33},
		{14, 24, 34},
		{15, 25, 35},
		{16, 26, 36},
	})

	// evolver
	evo := NewEvolverPop([]Population{pop}, ovfunc)
	evo.FnKey = "test_evo01"

	// set island
	evo.Islands[0].Roulette = true
	evo.Islands[0].ShowBases = true

	// run
	tf := 100
	dtout := 10
	dtmig := 20
	dtreg := 30
	nreg := 0
	io.Pf("\n")
	evo.Run(tf, dtout, dtmig, dtreg, nreg, true)

	// plot
	//if true {
	if false {
		evo.Islands[0].PlotOvs("/tmp", "fig_evo01", "", tf, true, "%.6f", true, true)
	}
}

func Test_evo02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo02. organise sequence of ints")
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
	ovfunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
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
		ind.ObjValue = 1.0 / (1.0 + score)
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

	// evolver
	nislands := 3
	ninds := 6
	evo := NewEvolver(nislands, ninds, ref, bingo, ovfunc)
	for _, isl := range evo.Islands {
		isl.MtProbs = make(map[string]float64)
		isl.MtProbs["int"] = 0.01
		isl.MtIntFunc = mtfunc
	}

	// saving files
	evo.FnKey = "evo02"

	// run
	tf := 100
	dtout := 20
	dtmig := 40
	dtreg := 50
	nreg := -1
	evo.Run(tf, dtout, dtmig, dtreg, nreg, true)

	// results
	ideal := 1.0 / (1.0 + float64(nvals))
	io.PfGreen("\nBest = %v\nBestOV = %v  (ideal=%v)\n", evo.Best.Ints, evo.Best.ObjValue, ideal)

	// plot
	//if true {
	if false {
		for i, isl := range evo.Islands {
			first := i == 0
			last := i == nislands-1
			isl.PlotOvs("/tmp", "fig_evo02", "", tf, true, "%.6f", first, last)
		}
	}
}

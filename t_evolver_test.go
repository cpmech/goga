// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
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
	ovfunc := func(ind *Individual, time int, best *Individual) {
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

	// run
	tf := 100
	dtout := 10
	dtmig := 20
	io.Pf("\n")
	evo.Run(tf, dtout, dtmig)
}

func Test_evo02(tst *testing.T) {

	verbose()
	chk.PrintTitle("evo02. organise sequence of ints")

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
	ovfunc := func(ind *Individual, time int, best *Individual) {
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
	bingo := NewBingoInts([]int{0}, []int{1})

	// evolver
	nislands := 2
	ninds := 10
	evo := NewEvolver(nislands, ninds, ref, bingo, ovfunc)
	for _, isl := range evo.Islands {
		isl.MtProbs = make(map[string]float64)
		isl.MtProbs["int"] = 0.01
		isl.MtIntFunc = mtfunc
		io.Pforan("\n%v\n", isl.Pop.Output(nil))
	}

	// run
	tf := 100
	dtout := 20
	dtmig := 1000
	evo.Run(tf, dtout, dtmig)

	// results
	io.Pf("\n")
	for _, isl := range evo.Islands {
		isl.MtProbs = make(map[string]float64)
		isl.MtProbs["int"] = 0.01
		isl.MtIntFunc = mtfunc
		io.Pfgreen("%v\n", isl.Pop.Output(nil))
	}
	io.PfGreen("\nBest = %v\nBestOV = %v\n", evo.Best.Ints, evo.Best.ObjValue)
}

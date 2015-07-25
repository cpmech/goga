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

	nbases := 8
	pop := NewPopFloatChromo(nbases, [][]float64{
		{11, 21, 31},
		{12, 22, 32},
		{13, 23, 33},
		{14, 24, 34},
		{15, 25, 35},
		{16, 26, 36},
	})

	ofunc := func(ind *Individual, time int, best *Individual) {
		ind.ObjValue = 1.0 / (1.0 + (ind.GetFloat(0)+ind.GetFloat(1)+ind.GetFloat(2))/3.0)
	}

	isl := NewIsland(pop, ofunc)
	isl.MtProbs = make(map[string]float64)
	isl.MtProbs["flt"] = 0.0
	io.Pforan("%v\n", isl.Pop.Output(nil))
	io.Pforan("best = %v\n", isl.Pop[0].Output(nil))
	io.Pf("\n")

	tf := 100
	dtout := 10
	dtmig := 20
	evo := Evolver{[]*Island{isl}}
	evo.Run(tf, dtout, dtmig)
}

func Test_evo02(tst *testing.T) {

	verbose()
	chk.PrintTitle("evo02. organise sequence of ints")

	rnd.Init(0)

	mtfunc := func(A []int, nchanges int, pm float64, extra interface{}) {
		size := len(A)
		if !rnd.FlipCoin(pm) || size < 1 {
			return
		}
		pos := rnd.IntGetUniqueN(0, size, nchanges)
		io.Pforan("here = %v\n", 1)
		for _, i := range pos {
			if A[i] == 1 {
				A[i] = 0
			}
			if A[i] == 0 {
				A[i] = 1
			}
		}
	}

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

	// template individual
	nvals := 20
	ind := NewIndividual(1, utl.IntVals(nvals, 1))
	for i := 0; i < nvals; i++ {
		ind.Ints[i] = rand.Intn(2)
	}

	// bingo and population
	ninds := 10
	bingo := NewBingoInts(nvals, 0, 1)
	bingo.UseIntRnd = true
	pop := NewPopRandom(ninds, ind, bingo)

	// island and evolver
	isl := NewIsland(pop, ovfunc)
	isl.MtProbs = make(map[string]float64)
	isl.MtProbs["int"] = 0.01
	isl.MtIntFunc = mtfunc
	io.Pforan("%v\n", isl.Pop.Output(nil))

	// run
	tf := 200
	dtout := 10
	dtmig := 1000
	evo := Evolver{[]*Island{isl}}
	evo.Run(tf, dtout, dtmig)
	io.Pfgreen("%v\n", isl.Pop.Output(nil))
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func Test_island01(tst *testing.T) {

	verbose()
	chk.PrintTitle("island01")

	genes := [][]float64{
		{11, 21, 31},
		{12, 22, 32},
		{13, 23, 33},
		{14, 24, 34},
		{15, 25, 35},
		{16, 26, 36},
	}

	// parameters
	C := NewConfParams()
	C.Ninds = len(genes)
	C.Nbases = 1
	C.Rnk = false
	C.RangeFlt = [][]float64{
		{10, 20},
		{20, 30},
		{30, 40},
	}

	// generator
	nova := 2
	noor := 3
	C.PopFltGen = func(pop Population, ninds, nbases int, noise float64, args interface{}, frange [][]float64) Population {
		o := make([]*Individual, ninds)
		for i := 0; i < ninds; i++ {
			o[i] = NewIndividual(nova, noor, nbases, genes[i])
		}
		return o
	}

	// objective function
	// the best will have the largest genes (x,y,z);
	// but with the first gene (x) smaller than or equal to 13
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		x, y, z := ind.GetFloat(0), ind.GetFloat(1), ind.GetFloat(2)
		ind.Ovas[0] = -(x + y + z)
		ind.Ovas[1] = z
		ind.Oors[0] = utl.GtePenalty(13, x, 1)
		ind.Oors[1] = utl.GtePenalty(24, y, 1)
		ind.Oors[2] = utl.GtePenalty(z, 35, 1)
		return
	}

	sorted_genes := [][]float64{
		{13, 23, 33},
		{12, 22, 32},
		{11, 21, 31},
		{14, 24, 34},
		{15, 25, 35},
		{16, 26, 36},
	}

	// island
	isl := NewIsland(0, C)
	io.Pf("%v", isl.Pop.Output(nil, false))
	io.Pf("best = %v\n", isl.Pop[0].Output(nil, false))
	chk.Vector(tst, "best", 1e-17, isl.Pop[0].Floats, []float64{13, 23, 33})
	for i, ind := range isl.Pop {
		chk.Vector(tst, io.Sf("ind%d", i), 1e-17, ind.Floats, sorted_genes[i])
	}
	/*
		io.Pforan("sovas0 = %.4f\n", isl.sovas[0])
		io.Pforan("sovas1 = %.4f\n", isl.sovas[1])
		io.Pfcyan("soors0 = %.4f\n", isl.soors[0])
		io.Pfcyan("soors1 = %.4f\n", isl.soors[1])
		io.Pfcyan("soors2 = %.4f\n", isl.soors[2])
	*/

	return

	isl.SelectReprodAndRegen(0, false, false, false)
	io.Pfcyan("%v\n", isl.Pop.Output(nil, false))

	isl.SelectReprodAndRegen(1, false, false, false)
	io.Pforan("%v\n", isl.Pop.Output(nil, false))

	isl.SelectReprodAndRegen(2, false, false, false)
	io.Pfcyan("%v\n", isl.Pop.Output(nil, false))

	// TODO: more tests required
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_island01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("island01")

	genes := [][]float64{
		{11, 21, 31},
		{13, 23, 33},
		{15, 25, 35},
		{12, 22, 32},
		{16, 26, 36},
		{14, 24, 34},
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
	C.PopFltGen = func(pop Population, ninds, nbases int, noise float64, args interface{}, frange [][]float64) Population {
		o := make([]*Individual, ninds)
		for i := 0; i < ninds; i++ {
			o[i] = NewIndividual(nbases, genes[i])
		}
		return o
	}

	// objective function
	// the best will have the largest genes (x,y,z);
	// but with the first gene (x) smaller than or equal to 13
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) (ov, oor float64) {
		x, y, z := ind.GetFloat(0), ind.GetFloat(1), ind.GetFloat(2)
		ov = 1.0 / (1.0 + (x+y+z)/3.0)
		if ind.GetFloat(0) > 13 {
			oor = x - 10
		}
		return
	}

	// island
	isl := NewIsland(0, C)
	io.Pforan("%v\n", isl.Pop.Output(nil, false))
	io.Pforan("best = %v\n", isl.Pop[0].Output(nil, false))
	chk.Vector(tst, "best", 1e-17, isl.Pop[0].Floats, []float64{13, 23, 33})

	isl.SelectReprodAndRegen(0, false, false, false)
	io.Pfcyan("%v\n", isl.Pop.Output(nil, false))

	isl.SelectReprodAndRegen(1, false, false, false)
	io.Pforan("%v\n", isl.Pop.Output(nil, false))

	isl.SelectReprodAndRegen(2, false, false, false)
	io.Pfcyan("%v\n", isl.Pop.Output(nil, false))

	// TODO: more tests required
}

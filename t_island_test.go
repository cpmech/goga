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

	nbases := 1
	pop := NewPopFloatChromo(nbases, [][]float64{
		{11, 21, 31},
		{13, 23, 33},
		{15, 25, 35},
		{12, 22, 32},
		{16, 26, 36},
		{14, 24, 34},
	})

	// the best will have the largest genes (x,y,z);
	// but with the first gene (x) smaller than or equal to 13
	ofunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) (ov, oor float64) {
		x, y, z := ind.GetFloat(0), ind.GetFloat(1), ind.GetFloat(2)
		ov = 1.0 / (1.0 + (x+y+z)/3.0)
		if ind.GetFloat(0) > 13 {
			oor = x - 10
		}
		return
	}

	// parameters
	C := NewConfParams()
	C.Ninds = len(pop)
	C.Rnk = false

	// bingo
	bingo := NewBingoFloats([]float64{-100, -200, -300}, []float64{100, 200, 300})

	// island
	isl := NewIsland(0, C, pop, ofunc, bingo)
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

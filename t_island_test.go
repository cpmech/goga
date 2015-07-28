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
		{12, 22, 32},
		{13, 23, 33},
		{14, 24, 34},
		{15, 25, 35},
		{16, 26, 36},
	})

	ofunc := func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		ind.ObjValue = 1.0 / (1.0 + (ind.GetFloat(0)+ind.GetFloat(1)+ind.GetFloat(2))/3.0)
	}

	isl := NewIsland(0, pop, ofunc)
	io.Pforan("%v\n", isl.Pop.Output(nil, false))
	io.Pforan("best = %v\n", isl.Pop[0].Output(nil, false))
	chk.Vector(tst, "best", 1e-17, isl.Pop[0].Floats, []float64{16, 26, 36})
	chk.Scalar(tst, "ind2", 1e-17, isl.Pop[2].ObjValue, 0.04) // 1/25

	isl.SelectAndReprod(0)
	io.Pfcyan("%v\n", isl.Pop.Output(nil, false))

	isl.SelectAndReprod(1)
	io.Pforan("%v\n", isl.Pop.Output(nil, false))

	isl.SelectAndReprod(2)
	io.Pfcyan("%v\n", isl.Pop.Output(nil, false))

	// TODO: more tests required
}

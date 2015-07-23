// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_island01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("island01")

	pop := NewPopFloatChromo(1, [][]float64{
		{11, 21, 31},
		{12, 22, 32},
		{13, 23, 33},
		{14, 24, 34},
		{15, 25, 35},
		{16, 26, 36},
	})

	ofunc := func(ind *Individual, time int, best *Individual) {
		ind.ObjValue = 1.0 / (1.0 + ind.GetFloat(0) + ind.GetFloat(1) + ind.GetFloat(2))
	}

	isl := NewIsland(pop, ofunc)
	best := isl.BestInd
	io.Pforan("%v\n", isl.Pop.Output(nil))
	io.Pforan("%v\n", best.Output(nil))
	chk.Vector(tst, "best", 1e-17, best.Floats, []float64{16, 26, 36})
}

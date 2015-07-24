// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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

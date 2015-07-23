// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func Test_island01(tst *testing.T) {

	verbose()
	chk.PrintTitle("island01")

	pop := NewPopFloatChromo(1, [][]float64{
		{0},
		{1},
		{2},
		{3},
		{4},
		{5},
		{6},
		{7},
		{8},
		{9},
		{10},
	})

	ofunc := func(ind *Individual, time int, best *Individual) {
		ind.ObjValue = ind.Floats[0]
	}

	isl := NewIsland(pop, ofunc)
	best := isl.oldbest
	io.Pforan("%v\n", isl.Pop.Output(nil))
	io.Pforan("%v\n", best.Output(nil))

	isl.C.UseRanking = true
	isl.C.RnkPressure = 2
	isl.Reproduction(0)
	io.Pf("f = %v\n", isl.fitness)
	io.Pfyel("p = %.5f\n", isl.prob)
	io.Pfcyan("m = %.5f\n", isl.cumprob)
	chk.Scalar(tst, "sum(f)", 1e-15, la.VecAccum(isl.fitness), 11)
	chk.Vector(tst, "f", 1e-15, isl.fitness, []float64{2, 1.8, 1.6, 1.4, 1.2, 1, 0.8, 0.6, 0.4, 0.2, 0})
	chk.Vector(tst, "p", 1e-15, isl.prob, []float64{2 / 11.0, 1.8 / 11.0, 1.6 / 11.0, 1.4 / 11.0, 1.2 / 11.0, 1 / 11.0, 0.8 / 11.0, 0.6 / 11.0, 0.4 / 11.0, 0.2 / 11.0, 0})
	chk.Vector(tst, "m", 1e-15, isl.cumprob, []float64{2 / 11.0, 3.8 / 11.0, 5.4 / 11.0, 6.8 / 11.0, 8 / 11.0, 9 / 11.0, 9.8 / 11.0, 10.4 / 11.0, 10.8 / 11.0, 11.0 / 11.0, 1})
	//sol := []float64{2 / 11.0, 3.8 / 11.0, 5.4 / 11.0, 6.8 / 11.0, 8 / 11.0, 9 / 11.0, 9.8 / 11.0, 10.4 / 11.0, 10.8 / 11.0, 11.0 / 11.0, 1}
	//for i := 0; i < len(isl.cumprob); i++ {
	//chk.PrintAnaNum("", 1e-15, isl.cumprob[i], sol[i], true)
	//}
}

func Test_island02(tst *testing.T) {

	verbose()
	chk.PrintTitle("island02")

	pop := NewPopFloatChromo(1, [][]float64{
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
	best := isl.oldbest
	io.Pforan("%v\n", isl.Pop.Output(nil))
	io.Pforan("%v\n", best.Output(nil))
	chk.Vector(tst, "best", 1e-17, best.Floats, []float64{16, 26, 36})
	chk.Scalar(tst, "ind2", 1e-17, isl.Pop[2].ObjValue, 0.04) // 1/25
}

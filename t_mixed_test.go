// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
)

func Test_mix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mix01. mixed: float and ints")

	// parameters
	C := NewConfParams()
	C.Pll = false
	C.Nisl = 1
	C.Ninds = 24
	C.NumInts = 2
	C.RangeFlt = [][]float64{
		{-2, 2}, // gene # 0: min and max
		{-2, 2}, // gene # 1: min and max
	}
	C.PopFltGen = PopFltAndBinGen
	C.Ops.MtInt = IntBinMutation
	C.Ops.FltCxName = "de"
	C.Nova = 1
	C.Noor = 0
	C.CalcDerived()
	rnd.Init(C.Seed)

	// objective function
	C.OvaOor = func(isl int, ind *Individual) {
		x := ind.Floats
		ind.Ovas[0] = math.Abs(math.Sqrt(x[0]*x[0]+x[1]*x[1]) - 1.0)
		if ind.Ints[0]+ind.Ints[1] != 1 {
			ind.Ovas[0] = 1
		}
		return
	}

	// run optimisation
	evo := NewEvolver(C)
	evo.Run()

	// plot
	if io.Verbose {
		inds := evo.GetFeasible()
		plt.SetForEps(1, 400)
		for _, ind := range inds {
			x := ind.Floats
			args := "'ro'"
			if ind.Ints[0] == 0 && ind.Ints[1] == 1 {
				args = "'go'"
			}
			if ind.Ints[0] == 1 && ind.Ints[1] == 0 {
				args = "'bo'"
			}
			plt.PlotOne(x[0], x[1], args+",clip_on=0")
		}
		plt.Circle(0, 0, 1, "ec='b', clip_on=0")
		plt.Equal()
		plt.Gll("x", "y", "")
		plt.SaveD("/tmp/goga", "fig_mix01.eps")
	}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_sort01(tst *testing.T) {

	verbose()
	chk.PrintTitle("sort01. non-dominated sort")

	// genes
	genes := [][]float64{
		{1.0, 0.3},
		{0.4, 0.7},
		{0.5, 1.0},
		{0.3, 0.2},
		{0.8, 0.3},
		{0.2, 0.6},
		{0.6, 0.0},
		{0.7, 0.6},
		{0.8, 0.7},
		{0.5, 0.0},
		{0.0, 0.5},
		{0.4, 0.4},
		{0.3, 0.9},
		{0.9, 0.2},
	}

	// parameters
	C := NewConfParams()
	C.Pll = false
	C.Nisl = 1
	C.Ninds = len(genes)
	C.Nova = 2
	C.RangeFlt = [][]float64{
		{0, 1},
		{0, 1},
	}
	C.CalcDerived()

	// generator
	C.PopFltGen = func(isl int, cc *ConfParams) Population {
		o := make([]*Individual, cc.Ninds)
		for i := 0; i < cc.Ninds; i++ {
			o[i] = NewIndividual(cc.Nova, cc.Noor, cc.Nbases, genes[i])
			o[i].Id = i
		}
		return o
	}

	// objective function
	C.OvaOor = func(ind *Individual, isl, time int, report *bytes.Buffer) {
		ind.Ovas[0] = ind.GetFloat(0)
		ind.Ovas[1] = ind.GetFloat(1)
	}

	// run
	isl := NewIsland(0, C)
	isl.NonDomSort(isl.Pop)
	isl.CalcMinMaxOva(isl.Pop)
	isl.CalcCrowdDist(isl.Pop)

	// check fronts
	chk.IntAssert(isl.nfronts, 4)
	for r := 0; r < isl.nfronts; r++ {
		ids := make([]int, isl.fsizes[r])
		for s := 0; s < isl.fsizes[r]; s++ {
			i := isl.fronts[r][s]
			ind := isl.Pop[i]
			ids[s] = ind.Id
		}
		io.Pforan("ids = %v\n", ids)
		switch r {
		case 0:
			chk.Ints(tst, "front 0", ids, []int{9, 3, 10})
		case 1:
			chk.Ints(tst, "front 1", ids, []int{6, 11, 5})
		case 2:
			chk.Ints(tst, "front 2", ids, []int{4, 13, 1, 12, 7})
		case 3:
			chk.Ints(tst, "front 3", ids, []int{0, 2, 8})
		}
	}

	// check crowd distances
	io.Pf("\n")
	cdist := [][]float64{
		{INF, 1.0, INF},           // front 0: 9, 3, 10
		{INF, 1.0, INF},           // front 1: 6, 11, 5
		{0.6, INF, 0.7, INF, 0.8}, // front 2: 4, 13, 1, 12, 7
		{INF, INF, 1.2},           // front 3: 0, 2, 8
	}
	for r := 0; r < isl.nfronts; r++ {
		for s := 0; s < isl.fsizes[r]; s++ {
			i := isl.fronts[r][s]
			ind := isl.Pop[i]
			chk.Scalar(tst, io.Sf("D%2d%3d", r, ind.Id), 1e-14, ind.Cdist, cdist[r][s])
		}
		io.Pf("\n")
	}

	// plot
	if io.Verbose {
		plt.SetForEps(1, 455)
		for _, ind := range isl.Pop {
			x, y := ind.Ovas[0], ind.Ovas[1]
			plt.PlotOne(x, y, "'r.', clip_on=0")
			plt.Text(x, y, io.Sf("%d", ind.Id), "size=7, clip_on=0")
		}
		plt.Equal()
		plt.SetXnticks(13)
		plt.SetYnticks(13)
		plt.AxisRange(-0.1, 1.1, -0.1, 1.1)
		plt.Gll("$f_0$", "$f_1$", "")
		plt.SaveD("/tmp/goga", "fig_sort01.eps")
	}
}

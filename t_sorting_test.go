// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"sort"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_sort01(tst *testing.T) {

	//verbose()
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
		{0.6, 0.0},
		{0.2, 0.6},
	}

	// parameters
	C := NewConfParams()
	C.Pll = false
	C.Nisl = 1
	C.Ninds = len(genes)
	C.NparGrp = 2
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
			o[i] = NewIndividual(cc.Nova, cc.Noor, genes[i], nil)
			o[i].Id = i
		}
		return o
	}

	// objective function
	C.OvaOor = func(isl int, ind *Individual) {
		ind.Ovas[0] = ind.Floats[0]
		ind.Ovas[1] = ind.Floats[1]
	}

	// run
	isl := NewIsland(0, C)
	pop := isl.Pop
	pop.SortById()
	for i, ind := range pop {
		io.Pf("%2d => %2d %v\n", i, ind.Id, ind.Floats)
	}
	ninds := len(pop)
	ovamin, ovamax := make([]float64, 2), make([]float64, 2)
	fsizes := make([]int, ninds)
	fronts := make([][]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		fronts[i] = make([]*Individual, ninds)
	}
	nfronts := Metrics(ovamin, ovamax, fsizes, fronts, pop)

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

	// check number of wins/losses
	io.Pf("\n")
	NW := []int{0, 2, 0, 9, 2, 5, 5, 1, 0, 8, 7, 4, 1, 1, 5, 5}
	WO := [][]int{
		{},     // 0
		{2, 8}, // 1
		{},     // 2
		{0, 1, 2, 4, 7, 8, 11, 12, 13}, // 3
		{0, 8},           // 4
		{1, 2, 7, 8, 12}, // 5
		{0, 4, 7, 8, 13}, // 6
		{8},              // 7
		{},               // 8
		{0, 2, 4, 6, 7, 8, 13, 14}, // 9
		{1, 2, 5, 7, 8, 12, 15},    // 10
		{1, 2, 7, 8},               // 11
		{2},                        // 12
		{0},                        // 13
		{0, 4, 7, 8, 13},           // 14
		{1, 2, 7, 8, 12},           // 15
	}
	for i, A := range pop {
		if A.Nwins != NW[i] {
			chk.Panic("%2d : number of wins is incorrect %d ! %d\n", A.Id, A.Nwins, NW[i])
		}
		winover := make([]int, len(A.WinOver))
		for j := 0; j < A.Nwins; j++ {
			winover[j] = A.WinOver[j].Id
		}
		wo := winover[:A.Nwins]
		sort.Ints(wo)
		chk.Ints(tst, io.Sf("%2d: wo = %v\n", A.Id, wo), wo, WO[i])
	}

	// check fronts
	io.Pf("\n")
	FR := [][]int{
		{3, 9, 10},
		{5, 6, 11, 14, 15},
		{1, 4, 7, 12, 13},
		{0, 2, 8},
	}
	io.Pforan("nfronts = %v\n", nfronts)
	chk.IntAssert(nfronts, 4)
	for i := 0; i < nfronts; i++ {
		io.Pforan("front # %d : front size = %v\n", i, fsizes[i])
		ids := make([]int, fsizes[i])
		for j := 0; j < fsizes[i]; j++ {
			ids[j] = fronts[i][j].Id
		}
		sort.Ints(ids)
		chk.Ints(tst, io.Sf("front %d : %v\n", i, ids), ids, FR[i])
	}

	// check ranks (front ids)
	io.Pf("\n")
	frontids := make([]int, ninds)
	for i, ind := range pop {
		frontids[i] = ind.FrontId
	}
	chk.Ints(tst, "front ids", frontids, []int{3, 2, 3, 0, 2, 1, 1, 2, 3, 0, 0, 1, 2, 2, 1, 1})

	// check limits
	io.Pf("\n")
	chk.Vector(tst, "ovamin", 1e-15, ovamin, []float64{0, 0})
	chk.Vector(tst, "ovamax", 1e-15, ovamax, []float64{1, 1})

	// check neighbour distances
	io.Pf("\n")
	DD := [][]float64{
		{0.1, 0.1}, //  0
		{0.1, 0.2}, //  1
		{0.1, 0.2}, //  2
		{0.1, 0.2}, //  3
		{0.1, 0.1}, //  4
		{0.0, 0.0}, //  5
		{0.0, 0.0}, //  6
		{0.1, 0.1}, //  7
		{0.1, 0.1}, //  8
		{0.1, 0.0}, //  9
		{0.1, 0.2}, // 10
		{0.1, 0.2}, // 11
		{0.1, 0.2}, // 12
		{0.1, 0.1}, // 13
		{0.0, 0.0}, // 14
		{0.0, 0.0}, // 15
	}
	for i, ind := range pop {
		chk.Scalar(tst, io.Sf("%2d : dNeigh = %g", ind.Id, ind.DistNeigh), 1e-15, ind.DistNeigh, math.Sqrt(math.Pow(DD[i][0], 2.0)+math.Pow(DD[i][1], 2.0)))
	}

	// check crowd distances
	io.Pf("\n")
	DC := [][]float64{
		{INF, 1.0, 1.0, 1.0}, // 0: dx,dx, dy,dy
		{0.1, 0.3, 0.2, 0.1}, // 1
		{INF, 1.0, 1.0, 1.0}, // 2
		{0.3, 0.2, 0.3, 0.2}, // 3
		{0.1, 0.1, 0.3, 0.1}, // 4
		{INF, 1.0, 1.0, 1.0}, // 5
		{INF, 1.0, 1.0, 1.0}, // 6
		{0.3, 0.1, 0.1, 0.3}, // 7
		{0.3, 0.2, 0.3, 0.4}, // 8
		{INF, 1.0, 1.0, 1.0}, // 9
		{INF, 1.0, 1.0, 1.0}, // 10
		{0.2, 0.2, 0.2, 0.4}, // 11
		{INF, 1.0, 1.0, 1.0}, // 12
		{INF, 1.0, 1.0, 1.0}, // 13
		{INF, 1.0, 1.0, 1.0}, // 14
		{INF, 1.0, 1.0, 1.0}, // 15
	}
	for i, ind := range pop {
		chk.Scalar(tst, io.Sf("%2d : dCrowd = %g", ind.Id, ind.DistCrowd), 1e-15, ind.DistCrowd, DC[i][0]*DC[i][1]+DC[i][2]*DC[i][3])
	}
}

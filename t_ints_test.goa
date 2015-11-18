// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"sort"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_int01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("int01. organise sequence of ints")
	io.Pf("\n")

	// initialise random numbers generator
	rnd.Init(0) // 0 => use current time as seed

	// parameters
	C := NewConfParams()
	C.Nova = 1
	C.Noor = 0
	C.Nisl = 1
	C.Ninds = 20
	C.RegTol = 0
	C.NumInts = 20
	//C.GAtype = "crowd"
	C.NparGrp = 2
	C.Tf = 50
	C.Verbose = chk.Verbose
	C.CalcDerived()

	// mutation function
	C.Ops.MtInt = IntBinMutation

	// generation function
	C.PopIntGen = PopBinGen

	// objective function
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		score := 0.0
		count := 0
		for _, val := range ind.Ints {
			if val == 0 && count%2 == 0 {
				score += 1.0
			}
			if val == 1 && count%2 != 0 {
				score += 1.0
			}
			count++
		}
		ind.Ovas[0] = 1.0 / (1.0 + score)
		return
	}

	// run optimisation
	evo := NewEvolver(C)
	evo.Run()

	// results
	ideal := 1.0 / (1.0 + float64(C.NumInts))
	io.PfGreen("\nBest = %v\nBestOV = %v  (ideal=%v)\n", evo.Best.Ints, evo.Best.Ovas[0], ideal)
}

func Test_int02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("int02. TSP")

	// location / coordinates of stations
	locations := [][]float64{
		{60, 200}, {180, 200}, {80, 180}, {140, 180}, {20, 160}, {100, 160}, {200, 160},
		{140, 140}, {40, 120}, {100, 120}, {180, 100}, {60, 80}, {120, 80}, {180, 60},
		{20, 40}, {100, 40}, {200, 40}, {20, 20}, {60, 20}, {160, 20},
	}
	nstations := len(locations)

	// parameters
	C := NewConfParams()
	C.Nova = 1
	C.Noor = 0
	C.Nisl = 4
	C.Ninds = 24
	C.RegTol = 0.3
	C.RegPct = 0.2
	//C.Dtmig = 30
	C.GAtype = "crowd"
	C.ParetoPhi = 0.1
	C.Elite = false
	C.DoPlot = false //chk.Verbose
	//C.Rws = true
	C.SetIntOrd(nstations)
	C.CalcDerived()

	// initialise random numbers generator
	rnd.Init(0)

	// objective value function
	C.OvaOor = func(ind *Individual, idIsland, t int, report *bytes.Buffer) {
		L := locations
		ids := ind.Ints
		//io.Pforan("ids = %v\n", ids)
		dist := 0.0
		for i := 1; i < nstations; i++ {
			a, b := ids[i-1], ids[i]
			dist += math.Sqrt(math.Pow(L[b][0]-L[a][0], 2.0) + math.Pow(L[b][1]-L[a][1], 2.0))
		}
		a, b := ids[nstations-1], ids[0]
		dist += math.Sqrt(math.Pow(L[b][0]-L[a][0], 2.0) + math.Pow(L[b][1]-L[a][1], 2.0))
		ind.Ovas[0] = dist
		return
	}

	// evolver
	evo := NewEvolver(C)

	// print initial population
	pop := evo.Islands[0].Pop
	//io.Pf("\n%v\n", pop.Output(nil, false))

	// 0,4,8,11,14,17,18,15,12,19,13,16,10,6,1,3,7,9,5,2 894.363
	if false {
		for i, x := range []int{0, 4, 8, 11, 14, 17, 18, 15, 12, 19, 13, 16, 10, 6, 1, 3, 7, 9, 5, 2} {
			pop[0].Ints[i] = x
		}
		evo.Islands[0].CalcOvs(pop, 0)
		evo.Islands[0].CalcDemeritsCdistAndSort(pop)
	}

	// check initial population
	ints := make([]int, nstations)
	if false {
		for i := 0; i < C.Ninds; i++ {
			for j := 0; j < nstations; j++ {
				ints[j] = pop[i].Ints[j]
			}
			sort.Ints(ints)
			chk.Ints(tst, "ints", ints, utl.IntRange(nstations))
		}
	}

	// run
	evo.Run()
	//io.Pf("%v\n", pop.Output(nil, false))
	io.Pfgreen("best = %v\n", evo.Best.Ints)
	io.Pfgreen("best OVA = %v  (871.117353844847)\n\n", evo.Best.Ovas[0])

	// best = [18 17 14 11 8 4 0 2 5 9 12 7 6 1 3 10 16 13 19 15]
	// best OVA = 953.4643474956656

	// best = [8 11 14 17 18 15 12 19 16 13 10 6 1 3 7 9 5 2 0 4]
	// best OVA = 871.117353844847

	// best = [5 2 0 4 8 11 14 17 18 15 12 19 16 13 10 6 1 3 7 9]
	// best OVA = 871.1173538448469

	// best = [6 10 13 16 19 15 18 17 14 11 8 4 0 2 5 9 12 7 3 1]
	// best OVA = 880.7760751923065

	// check final population
	if false {
		for i := 0; i < C.Ninds; i++ {
			for j := 0; j < nstations; j++ {
				ints[j] = pop[i].Ints[j]
			}
			sort.Ints(ints)
			chk.Ints(tst, "ints", ints, utl.IntRange(nstations))
		}
	}

	// plot travelling salesman path
	if C.DoPlot {
		plt.SetForEps(1, 300)
		X, Y := make([]float64, nstations), make([]float64, nstations)
		for k, id := range evo.Best.Ints {
			X[k], Y[k] = locations[id][0], locations[id][1]
			plt.PlotOne(X[k], Y[k], "'r.', ms=5, clip_on=0, zorder=20")
			plt.Text(X[k], Y[k], io.Sf("%d", id), "fontsize=7, clip_on=0, zorder=30")
		}
		plt.Plot(X, Y, "'b-', clip_on=0, zorder=10")
		plt.Plot([]float64{X[0], X[nstations-1]}, []float64{Y[0], Y[nstations-1]}, "'b-', clip_on=0, zorder=10")
		plt.Equal()
		plt.AxisRange(10, 210, 10, 210)
		plt.Gll("$x$", "$y$", "")
		plt.SaveD("/tmp/goga", "test_evo04.eps")
	}
}

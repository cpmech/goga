// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
)

func main() {

	// GA parameters
	C := goga.ReadConfParams("tsp-simple.json")
	rnd.Init(C.Seed)

	// location / coordinates of stations
	locations := [][]float64{
		{60, 200}, {180, 200}, {80, 180}, {140, 180}, {20, 160}, {100, 160}, {200, 160},
		{140, 140}, {40, 120}, {100, 120}, {180, 100}, {60, 80}, {120, 80}, {180, 60},
		{20, 40}, {100, 40}, {200, 40}, {20, 20}, {60, 20}, {160, 20},
	}
	nstations := len(locations)
	C.SetIntOrd(nstations)
	C.CalcDerived()

	// objective value function
	C.OvaOor = func(ind *goga.Individual, idIsland, time int, report *bytes.Buffer) {
		L := locations
		ids := ind.Ints
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
	nova, noor := 1, 0
	evo := goga.NewEvolver(nova, noor, C)
	evo.Run()

	// results
	io.Pfgreen("best = %v\n", evo.Best.Ints)
	io.Pfgreen("best OVA = %v  (871.117353844847)\n\n", evo.Best.Ovas[0])

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

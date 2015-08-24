// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func Test_island01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("island01")

	genes := [][]float64{
		{11, 21, 31},
		{12, 22, 32},
		{13, 23, 33},
		{14, 24, 34},
		{15, 25, 35},
		{16, 26, 36},
	}

	// parameters
	C := NewConfParams()
	C.Ninds = len(genes)
	C.Nbases = 1
	C.Rnk = false
	C.RangeFlt = [][]float64{
		{10, 20},
		{20, 30},
		{30, 40},
	}

	// generator
	C.PopFltGen = func(ninds, nova, noor, nbases int, noise float64, args interface{}, frange [][]float64) Population {
		o := make([]*Individual, ninds)
		for i := 0; i < ninds; i++ {
			o[i] = NewIndividual(nova, noor, nbases, genes[i])
		}
		return o
	}

	// objective function
	// the best will have the largest genes (x,y,z);
	// but with the first gene (x) smaller than or equal to 13
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		x, y, z := ind.GetFloat(0), ind.GetFloat(1), ind.GetFloat(2)
		ind.Ovas[0] = -(x + y + z)
		ind.Ovas[1] = z
		ind.Oors[0] = utl.GtePenalty(13, x, 1)
		ind.Oors[1] = utl.GtePenalty(24, y, 1)
		ind.Oors[2] = utl.GtePenalty(z, 35, 1)
		return
	}

	sorted_genes := [][]float64{
		{13, 23, 33},
		{14, 24, 34},
		{12, 22, 32},
		{11, 21, 31},
		{15, 25, 35},
		{16, 26, 36},
	}

	// island
	id := 0
	nova := 2
	noor := 3
	isl := NewIsland(id, nova, noor, C)
	io.Pf("\n%v", isl.Pop.Output(nil, true, false, -1))
	io.Pf("best = %v\n", isl.Pop[0].Output(nil, false))
	vals := make([]float64, 3)
	for i := 0; i < 3; i++ {
		vals[i] = isl.Pop[0].GetFloat(i)
	}
	chk.Vector(tst, "best", 1e-14, vals, []float64{13, 23, 33})
	for i, ind := range isl.Pop {
		for j := 0; j < 3; j++ {
			vals[j] = ind.GetFloat(j)
		}
		chk.Vector(tst, io.Sf("ind%d", i), 1e-14, vals, sorted_genes[i])
	}
	/*
		io.Pforan("sovas0 = %.4f\n", isl.sovas[0])
		io.Pforan("sovas1 = %.4f\n", isl.sovas[1])
		io.Pfcyan("soors0 = %.4f\n", isl.soors[0])
		io.Pfcyan("soors1 = %.4f\n", isl.soors[1])
		io.Pfcyan("soors2 = %.4f\n", isl.soors[2])
	*/

	// stat
	minrho, averho, maxrho, devrho := isl.FltStat()
	io.Pforan("\nallbases[0] = %.6f\n", isl.allbases[0])
	io.Pforan("allbases[1] = %.6f\n", isl.allbases[1])
	io.Pforan("allbases[2] = %.6f\n", isl.allbases[2])
	if C.Nbases > 1 {
		io.Pforan("allbases[3] = %.6f\n", isl.allbases[3])
		io.Pforan("allbases[4] = %.6f\n", isl.allbases[4])
		io.Pforan("allbases[5] = %.6f\n", isl.allbases[5])
	}
	io.Pfcyan("devbases[0] = %.6f\n", isl.devbases[0])
	io.Pfcyan("devbases[1] = %.6f\n", isl.devbases[1])
	io.Pfcyan("devbases[2] = %.6f\n", isl.devbases[2])
	io.Pforan("minrho = %v\n", minrho)
	io.Pforan("averho = %v\n", averho)
	io.Pforan("maxrho = %v\n", maxrho)
	io.Pforan("devrho = %v\n", devrho)

	isl.Run(0, false, false)
	io.Pf("\n%v\n", isl.Pop.Output(nil, true, false, -1))

	isl.Run(1, false, false)
	io.Pforan("%v\n", isl.Pop.Output(nil, true, false, -1))

	isl.Run(2, false, false)
	io.Pf("%v\n", isl.Pop.Output(nil, true, false, -1))

	// TODO: more tests required
}

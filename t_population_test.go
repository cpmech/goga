// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_pop01(tst *testing.T) {

	verbose()
	chk.PrintTitle("pop01")

	// chromosomes
	chromos := [][]float64{
		{0, 1, 2, 3, 4, 5, -1, -1},
		{0, 1, 2, 1, 2, 1, -2, -1},
		{0, 5, 4, 3, 2, 1, -3, -1},
		{0, 1, 1, 1, 2, 2, -4, -1},
		{0, 2, 2, 2, 1, 1, -2, -1},
		{0, 3, 3, 3, 1, 1, -3, -1},
		{0, 4, 4, 4, 3, 1, -4, -1},
		{0, 4, 3, 3, 3, 3, -2, -1},
		{0, 1, 1, 2, 2, 0, -3, -1},
		{0, 0, 0, 0, 0, 0, -4, -1},
	}

	// objective values and fitness values
	ovs := []float64{11, 21, 10, 12, 13, 31, 41, 11, 51, 11}
	fits := []float64{0.1, 0.2, 1.2, 5.1, 0.1, 0.3, 2.1, 8.1, 0.5, 3.1}

	// init population
	nbases := 2
	ninds := len(chromos)
	pop := make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = new(Individual)
		err := pop[i].InitFloatChromo(nbases, chromos[i])
		if err != nil {
			tst.Errorf("InitFloatChromo failed:%v\n", err)
			return
		}
		pop[i].ObjValue = ovs[i]
		pop[i].Fitness = fits[i]
	}

	/*
		pop := Population{
			&Individual{ngenes, nbases, C[0], , , nil}, // 0
			&Individual{ngenes, nbases, C[0], , , nil}, // 1
			&Individual{ngenes, nbases, C[0], , , nil}, // 2
			&Individual{ngenes, nbases, C[0], , , nil}, // 3
			&Individual{ngenes, nbases, C[0], , , nil}, // 4
			&Individual{ngenes, nbases, C[0], , , nil}, // 5
			&Individual{ngenes, nbases, C[0], , , nil}, // 6
			&Individual{ngenes, nbases, C[0], , , nil}, // 7
			&Individual{ngenes, nbases, C[0], , , nil}, // 8
			&Individual{ngenes, nbases, C[0], , , nil}, // 9
		}
	*/

	/*
		for _, ind := range pop {
			ind.CalcGenes()
		}

		io.Pforan("\npopulation before sorting\n")
		io.Pf("%v\n", pop.GenTable(nil, nil, false))

		chk.Vector(tst, "genes0", 1e-17, pop[0].Genes, []float64{1, 5, 9, -2}) // 0
		chk.Vector(tst, "genes1", 1e-17, pop[1].Genes, []float64{1, 3, 3, -3}) // 1
		chk.Vector(tst, "genes2", 1e-17, pop[2].Genes, []float64{5, 7, 3, -4}) // 2
		chk.Vector(tst, "genes3", 1e-17, pop[3].Genes, []float64{1, 2, 4, -5}) // 3
		chk.Vector(tst, "genes4", 1e-17, pop[4].Genes, []float64{2, 4, 2, -3}) // 4
		chk.Vector(tst, "genes5", 1e-17, pop[5].Genes, []float64{3, 6, 2, -4}) // 5
		chk.Vector(tst, "genes6", 1e-17, pop[6].Genes, []float64{4, 8, 4, -5}) // 6
		chk.Vector(tst, "genes7", 1e-17, pop[7].Genes, []float64{4, 6, 6, -3}) // 7
		chk.Vector(tst, "genes8", 1e-17, pop[8].Genes, []float64{1, 3, 2, -4}) // 8
		chk.Vector(tst, "genes9", 1e-17, pop[9].Genes, []float64{0, 0, 0, -5}) // 9

		pop.Sort()

		io.Pforan("\npopulation after sorting\n")
		io.Pf("%v\n", pop.GenTable(nil, nil, false))

		chk.Vector(tst, "genes0", 1e-17, pop[0].Genes, []float64{4, 6, 6, -3}) // 7
		chk.Vector(tst, "genes1", 1e-17, pop[1].Genes, []float64{1, 2, 4, -5}) // 3
		chk.Vector(tst, "genes2", 1e-17, pop[2].Genes, []float64{0, 0, 0, -5}) // 9
		chk.Vector(tst, "genes3", 1e-17, pop[3].Genes, []float64{4, 8, 4, -5}) // 6
		chk.Vector(tst, "genes4", 1e-17, pop[4].Genes, []float64{5, 7, 3, -4}) // 2
		chk.Vector(tst, "genes5", 1e-17, pop[5].Genes, []float64{1, 3, 2, -4}) // 8
		chk.Vector(tst, "genes6", 1e-17, pop[6].Genes, []float64{3, 6, 2, -4}) // 5
		chk.Vector(tst, "genes7", 1e-17, pop[7].Genes, []float64{1, 3, 3, -3}) // 1
		chk.Vector(tst, "genes8", 1e-17, pop[8].Genes, []float64{2, 4, 2, -3}) // 4
		chk.Vector(tst, "genes9", 1e-17, pop[9].Genes, []float64{1, 5, 9, -2}) // 0
	*/
}

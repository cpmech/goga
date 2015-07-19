// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_pop01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("pop01")

	ngenes := 4
	nbases := 2

	pop := Population{
		&Individual{ngenes, nbases, []float64{0, 1, 2, 3, 4, 5, -1, -1}, 11, 0.1, nil}, // 0
		&Individual{ngenes, nbases, []float64{0, 1, 2, 1, 2, 1, -2, -1}, 21, 0.2, nil}, // 1
		&Individual{ngenes, nbases, []float64{0, 5, 4, 3, 2, 1, -3, -1}, 10, 1.2, nil}, // 2
		&Individual{ngenes, nbases, []float64{0, 1, 1, 1, 2, 2, -4, -1}, 12, 5.1, nil}, // 3
		&Individual{ngenes, nbases, []float64{0, 2, 2, 2, 1, 1, -2, -1}, 13, 0.1, nil}, // 4
		&Individual{ngenes, nbases, []float64{0, 3, 3, 3, 1, 1, -3, -1}, 31, 0.3, nil}, // 5
		&Individual{ngenes, nbases, []float64{0, 4, 4, 4, 3, 1, -4, -1}, 41, 2.1, nil}, // 6
		&Individual{ngenes, nbases, []float64{0, 4, 3, 3, 3, 3, -2, -1}, 11, 8.1, nil}, // 7
		&Individual{ngenes, nbases, []float64{0, 1, 1, 2, 2, 0, -3, -1}, 51, 0.5, nil}, // 8
		&Individual{ngenes, nbases, []float64{0, 0, 0, 0, 0, 0, -4, -1}, 11, 3.1, nil}, // 9
	}

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
}

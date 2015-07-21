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

	// genes
	genes := [][]float64{
		{1, 5, -2},
		{1, 3, -3},
		{5, 7, -4},
		{1, 2, -5},
		{2, 4, -3},
	}

	// objective values and fitness values
	ovs := []float64{11, 21, 10, 12, 13}
	fits := []float64{0.1, 0.2, 1.2, 5.1, 0.1}

	// init population
	nbases := 2
	pop := NewPopFloatChromo(nbases, genes)
	for i, ind := range pop {
		ind.ObjValue = ovs[i]
		ind.Fitness = fits[i]
	}
	io.Pforan("%v\n", pop.Output(nil))

	// check floats and subfloats
	for i, ind := range pop {
		for j, g := range ind.Chromo {
			chk.Scalar(tst, io.Sf("i%dg%d", i, j), 1e-17, g.GetFloat(), genes[i][j])
			chk.Scalar(tst, io.Sf("i%dg%d bases", i, j), 1e-15, g.Fbases[0]+g.Fbases[1], genes[i][j])
		}
	}

	// print bases
	io.Pf("\nbases (before)\n")
	io.Pf("%s\n", pop.OutFloatBases("%7.4f"))

	// change subfloats
	bases := [][]float64{
		{10, 1, 14, 5, -10, -1},
		{10, 1, 12, 1, -20, -1},
		{10, 5, 12, 1, -30, -1},
		{10, 1, 12, 2, -40, -1},
		{10, 2, 11, 1, -20, -1},
	}
	ngenes := len(genes[0])
	for i, b := range bases {
		for j := 0; j < ngenes; j++ {
			pop[i].Chromo[j].SetFbases(0, b[j*nbases:(j+1)*nbases])
		}
	}

	// print bases
	io.Pfyel("bases (after)\n")
	io.Pfyel("%s\n", pop.OutFloatBases("%g"))

	// checkf floats
	chk.Scalar(tst, "i0g0", 1e-17, pop[0].Chromo[0].GetFloat(), 11)
	chk.Scalar(tst, "i0g1", 1e-17, pop[0].Chromo[1].GetFloat(), 19)
	chk.Scalar(tst, "i0g2", 1e-17, pop[0].Chromo[2].GetFloat(), -11)

	chk.Scalar(tst, "i1g0", 1e-17, pop[1].Chromo[0].GetFloat(), 11)
	chk.Scalar(tst, "i1g1", 1e-17, pop[1].Chromo[1].GetFloat(), 13)
	chk.Scalar(tst, "i1g2", 1e-17, pop[1].Chromo[2].GetFloat(), -21)

	chk.Scalar(tst, "i2g0", 1e-17, pop[2].Chromo[0].GetFloat(), 15)
	chk.Scalar(tst, "i2g1", 1e-17, pop[2].Chromo[1].GetFloat(), 13)
	chk.Scalar(tst, "i2g2", 1e-17, pop[2].Chromo[2].GetFloat(), -31)

	chk.Scalar(tst, "i3g0", 1e-17, pop[3].Chromo[0].GetFloat(), 11)
	chk.Scalar(tst, "i3g1", 1e-17, pop[3].Chromo[1].GetFloat(), 14)
	chk.Scalar(tst, "i3g2", 1e-17, pop[3].Chromo[2].GetFloat(), -41)

	chk.Scalar(tst, "i4g0", 1e-17, pop[4].Chromo[0].GetFloat(), 12)
	chk.Scalar(tst, "i4g1", 1e-17, pop[4].Chromo[1].GetFloat(), 12)
	chk.Scalar(tst, "i4g2", 1e-17, pop[4].Chromo[2].GetFloat(), -21)
}

func Test_pop02(tst *testing.T) {

	verbose()
	chk.PrintTitle("pop02")

	genes := [][]float64{
		{1, 5}, // 0
		{1, 3}, // 1
		{5, 7}, // 2
		{1, 2}, // 3
		{2, 4}, // 4
		{3, 6}, // 5
		{4, 8}, // 6
		{4, 6}, // 7
		{1, 3}, // 8
		{0, 0}, // 9
	}

	// objective values and fitness values
	ovs := []float64{11, 21, 10, 12, 13, 31, 41, 11, 51, 11}
	fits := []float64{0.1, 0.2, 1.2, 5.1, 0.1, 0.3, 2.1, 8.1, 0.5, 3.1}

	// init population
	nbases := 2
	pop := NewPopFloatChromo(nbases, genes)
	for i, ind := range pop {
		ind.ObjValue = ovs[i]
		ind.Fitness = fits[i]
	}
	io.Pforan("%v\n", pop.Output(nil))

	pop.Sort()

	io.Pfyel("%v\n", pop.Output(nil))

	genes_sorted := [][]float64{
		{4, 6}, // 7
		{1, 2}, // 3
		{0, 0}, // 9
		{4, 8}, // 6
		{5, 7}, // 2
		{1, 3}, // 8
		{3, 6}, // 5
		{1, 3}, // 1
		{2, 4}, // 4
		{1, 5}, // 0
	}

	for i, ind := range pop {
		for j, g := range ind.Chromo {
			chk.Scalar(tst, io.Sf("i%dg%d", i, j), 1e-17, g.GetFloat(), genes_sorted[i][j])
		}
	}
}

func Test_pop03(tst *testing.T) {

	verbose()
	chk.PrintTitle("pop03")

	var ind Individual
	nbases := 3
	ind.InitChromo(nbases,
		[]int{1, 20, 300},
		[]float64{4.4, 5.5, 666},
		[]string{"abc", "b", "c"},
		[]byte("SGA"),
		[][]byte{[]byte("ABC"), []byte("DEF"), []byte("GHI")},
		[]Func_t{
			func(g *Gene) string { return "f0" },
			func(g *Gene) string { return "f1" },
			func(g *Gene) string { return "f2" },
		},
	)

	ninds := 5
	pop := NewPopReference(ninds, &ind)
	io.Pf("\n%v\n", pop.Output(nil))

	bingo := NewExampleBingo()
	for i, ind := range pop {
		for j, g := range ind.Chromo {
			idx := i
			if j > 0 {
				idx = -1
			}
			s := bingo.Draw(idx, ninds)
			g.SetInt(s.Int)
			g.SetFloat(s.Flt)
			g.SetString(s.String)
			g.SetByte(s.Byte)
			g.SetBytes(s.Bytes)
			g.SetFunc(s.Func)
		}
	}

	io.Pfyel("\n%v\n", pop.Output(nil))
}

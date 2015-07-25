// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Bingo collects values to be drawn in random operations
type Bingo struct {
	UseIntRnd bool        // generate random integers instead of selecting from grid
	UseFltRnd bool        // generate random float point numbers instead of selecting from grid
	IntRange  [][]int     // [ngene][nsamples] min and max integers
	FltRange  [][]float64 // [ngene][nsamples] min and max float point numbers
	PoolWords [][]string  // [ngene][nsamples] pool of words to be used in Gene.String
	PoolBytes [][]byte    // [ngene][nsamples] pool of bytes to be used in Gene.Byte
	PoolBtxt  [][]string  // [ngene][nsamples] pool of byte-words to be used in Gene.Bytes
	PoolFuncs [][]Func_t  // [ngene][nsamples] pool of functions
}

// NewExampleBingo returns a new Bingo with example values
func NewExampleBingo() *Bingo {
	return &Bingo{
		false, false,
		[][]int{{-10, 10}, {-20, 20}, {-30, 30}, {-40, 40}},
		[][]float64{{-123.0, 321.0}, {-1, 1}, {0, 1}},
		[][]string{
			{"circle", "square", "pentagon", "b-spline", "line", "point"},
			{"a", "b", "c", "d", "e", "f", "g"},
			{"int", "float64", "string", "byte"},
		},
		[][]byte{
			[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		},
		[][]string{
			{"apple", "banana", "mango", "orange", "peach", "kiwi"},
			{"red", "green", "blue", "cyan", "magenta", "black", "white"},
			{"x", "y", "z"},
		},
		[][]Func_t{
			{
				func(i *Individual) string { return "I" },
				func(i *Individual) string { return "F" },
				func(i *Individual) string { return "S" },
				func(i *Individual) string { return "B" },
				func(i *Individual) string { return "T" },
				func(i *Individual) string { return "F" },
			},
			{
				func(i *Individual) string { return "f0" },
				func(i *Individual) string { return "f1" },
				func(i *Individual) string { return "f2" },
				func(i *Individual) string { return "f3" },
			},
			{
				func(i *Individual) string { return "g0" },
				func(i *Individual) string { return "g1" },
				func(i *Individual) string { return "g2" },
				func(i *Individual) string { return "g3" },
			},
		},
	}
}

// NewBingoInts creates a bingo to generate int numbers between xmin and xmax
//  Input:
//   xmin -- min values of genes. len(xmin) == ngenes
//   xmax -- max values of genes
func NewBingoInts(xmin, xmax []int) (o *Bingo) {
	chk.IntAssert(len(xmin), len(xmax))
	o = new(Bingo)
	ngenes := len(xmin)
	o.IntRange = utl.IntsAlloc(ngenes, 2)
	for i := 0; i < ngenes; i++ {
		o.IntRange[i][0] = xmin[i]
		o.IntRange[i][1] = xmax[i]
	}
	return
}

// NewBingoFloats creates a bingo to generate float point numbers between xmin and xmax
//  Input:
//   nbases -- number of subdivisions of each gene
//   xmin   -- min values of genes (to be subdivided). len(xmin) == ngenes
//   xmax   -- max values of genes (to be subdivided)
func NewBingoFloats(nbases, xmin, xmax []float64) (o *Bingo) {
	chk.IntAssert(len(xmin), len(xmax))
	o = new(Bingo)
	ngenes := len(xmin)
	o.FltRange = utl.DblsAlloc(ngenes, 2)
	for i := 0; i < ngenes; i++ {
		o.FltRange[i][0] = xmin[i]
		o.FltRange[i][1] = xmax[i]
	}
	return
}

// DrawInt randomly selects an int from data pool
//  Input:
//   iInd  -- index of individual used to  compute value in Range: val = min + idx * Δ
//            use iInd = -1 to randomly choose between min and max
//   iGene -- index of gene
//   nInd  -- number of individuals
func (o Bingo) DrawInt(iInd, iGene, nInd int) int {
	if iGene < len(o.IntRange) {
		chk.IntAssert(len(o.IntRange[iGene]), 2)
		xmin := o.IntRange[iGene][0]
		xmax := o.IntRange[iGene][1]
		if iInd < 0 || nInd < 2 || o.UseIntRnd {
			return rnd.Int(xmin, xmax)
		}
		return xmin + iInd*(xmax-xmin)/(nInd-1)
	}
	return 0
}

// DrawFloat randomly selects a float point number from data pool
//  Input:
//   iInd  -- index of individual used to  compute value in Range: val = min + idx * Δ
//            use iInd = -1 to randomly choose between min and max
//   iGene -- index of gene
//   nInd  -- number of individuals
func (o Bingo) DrawFloat(iInd, iGene, nInd int) float64 {
	if iGene < len(o.FltRange) {
		chk.IntAssert(len(o.FltRange[iGene]), 2)
		xmin := o.FltRange[iGene][0]
		xmax := o.FltRange[iGene][1]
		if iInd < 0 || nInd < 2 || o.UseFltRnd {
			return rnd.Float64(xmin, xmax)
		}
		return xmin + float64(iInd)*(xmax-xmin)/float64(nInd-1)
	}
	return 0
}

// DrawString randomly selects a string from data pool
//  Input:
//   iGene -- index of gene
func (o Bingo) DrawString(iGene int) string {
	if iGene < len(o.PoolWords) {
		nw := len(o.PoolWords[iGene])
		return o.PoolWords[iGene][rnd.Int(0, nw-1)]
	}
	return ""
}

// DrawKey randomly selects a byte from data pool
//  Input:
//   iGene -- index of gene
func (o Bingo) DrawKey(iGene int) byte {
	if iGene < len(o.PoolBytes) {
		nb := len(o.PoolBytes[iGene])
		return o.PoolBytes[iGene][rnd.Int(0, nb-1)]
	}
	return 0
}

// DrawBytes randomly selects a []byte from data pool
//  Input:
//   iGene -- index of gene
func (o Bingo) DrawBytes(iGene int) []byte {
	if iGene < len(o.PoolBtxt) {
		nt := len(o.PoolBtxt[iGene])
		return []byte(o.PoolBtxt[iGene][rnd.Int(0, nt-1)])
	}
	return nil
}

// DrawFunc randomly selects a function from data pool
//  Input:
//   iGene -- index of gene
func (o Bingo) DrawFunc(iGene int) Func_t {
	if iGene < len(o.PoolFuncs) {
		nf := len(o.PoolFuncs[iGene])
		return o.PoolFuncs[iGene][rnd.Int(0, nf-1)]
	}
	return nil
}

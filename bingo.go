// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

// BingoSample holds the reslts of one Bingo.Draw
type BingoSample struct {
	Int    int     // an integer
	Flt    float64 // a float point number
	String string  // a string
	Byte   byte    // a byte
	Bytes  []byte  // a set of bytes
	Func   Func_t  // a function
}

// Bingo collects values to be drawn in random operations
type Bingo struct {
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
		[][]int{{-10, 10}, {-20, 20}, {-30, 30}},
		[][]float64{{-123.0, 321.0}, {-1, 1}, {0, 1}},
		[][]string{
			{"circle", "square", "pentagon", "b-spline", "line", "point"},
			{"a", "b", "c", "d"},
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
				func(g *Gene) string { return io.Sf("I=%v", g.GetInt()) },
				func(g *Gene) string { return io.Sf("F=%v", g.GetFloat()) },
				func(g *Gene) string { return io.Sf("S=%v", g.GetString()) },
				func(g *Gene) string { return io.Sf("B=%x", g.GetByte()) },
				func(g *Gene) string { return io.Sf("T=%s", string(g.GetBytes())) },
				func(g *Gene) string { return io.Sf("F=%v", g.GetFunc()) },
			},
			{
				func(g *Gene) string { return "f0" },
				func(g *Gene) string { return "f1" },
				func(g *Gene) string { return "f2" },
				func(g *Gene) string { return "f3" },
			},
			{
				func(g *Gene) string { return "g0" },
				func(g *Gene) string { return "g1" },
				func(g *Gene) string { return "g2" },
				func(g *Gene) string { return "g3" },
			},
		},
	}
}

// Draw randomly selects
//  Input:
//   iInd  -- index of individual used to  compute value in Range: val = min + idx * Δ
//            use iInd = -1 to randomly choose between min and max
//   iGene -- index of gene
//   nInd  -- number of individuals
func (o Bingo) Draw(iInd, iGene, nInd int) (sample BingoSample) {

	// integer
	if iGene < len(o.IntRange) {
		chk.IntAssert(len(o.IntRange[iGene]), 2)
		xmin := o.IntRange[iGene][0]
		xmax := o.IntRange[iGene][1]
		if iInd < 0 || nInd < 2 {
			sample.Int = rnd.Int(xmin, xmax)
		} else {
			sample.Int = xmin + iInd*(xmax-xmin)/(nInd-1)
		}
	}

	// float point number
	if iGene < len(o.FltRange) {
		chk.IntAssert(len(o.FltRange[iGene]), 2)
		xmin := o.FltRange[iGene][0]
		xmax := o.FltRange[iGene][1]
		if iInd < 0 || nInd < 2 {
			sample.Flt = rnd.Float64(xmin, xmax)
		} else {
			sample.Flt = xmin + float64(iInd)*(xmax-xmin)/float64(nInd-1)
		}
	}

	// string
	if iGene < len(o.PoolWords) {
		nw := len(o.PoolWords[iGene])
		sample.String = o.PoolWords[iGene][rnd.Int(0, nw-1)]
	}

	// byte
	if iGene < len(o.PoolBytes) {
		nb := len(o.PoolBytes[iGene])
		sample.Byte = o.PoolBytes[iGene][rnd.Int(0, nb-1)]
	}

	// bytes
	if iGene < len(o.PoolBtxt) {
		nt := len(o.PoolBtxt[iGene])
		sample.Bytes = []byte(o.PoolBtxt[iGene][rnd.Int(0, nt-1)])
	}

	// function
	if iGene < len(o.PoolFuncs) {
		nf := len(o.PoolFuncs[iGene])
		sample.Func = o.PoolFuncs[iGene][rnd.Int(0, nf-1)]
	}
	return
}

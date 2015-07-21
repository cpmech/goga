// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
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
	IntRange  []int     // min and max integers
	FltRange  []float64 // min and max float point numbers
	PoolWords []string  // pool of words to be used in Gene.String
	PoolBytes []byte    // pool of bytes to be used in Gene.Byte
	PoolBtxt  []string  // pool of byte-words to be used in Gene.Bytes
	PoolFuncs []Func_t  // pool of functions
}

// NewExampleBingo returns a new Bingo with example values
func NewExampleBingo() *Bingo {
	return &Bingo{
		[]int{-10, 10},
		[]float64{-123.0, 321.0},
		[]string{"circle", "square", "pentagon", "b-spline", "line", "point"},
		[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		[]string{"apple", "banana", "mango", "orange", "peach", "kiwi"},
		[]Func_t{
			func(g *Gene) string { return io.Sf("I=%v", g.GetInt()) },
			func(g *Gene) string { return io.Sf("F=%v", g.GetFloat()) },
			func(g *Gene) string { return io.Sf("S=%v", g.GetString()) },
			func(g *Gene) string { return io.Sf("B=%x", g.GetByte()) },
			func(g *Gene) string { return io.Sf("T=%s", string(g.GetBytes())) },
			func(g *Gene) string { return io.Sf("F=%v", g.GetFunc()) },
		},
	}
}

// Draw randomly selects
//  Input:
//   idx -- index to compute value in Range: val = min + idx * Δ
//          Note: use idx = -1 to randomly chose between min and max
//          e.g. idx == index of individual
//   num -- number of values from min to max
//          e.g. num == number of individuals
func (o Bingo) Draw(idx, num int) (sample BingoSample) {

	// integer
	if len(o.IntRange) == 2 {
		if idx < 0 || num < 2 {
			sample.Int = rnd.Int(o.IntRange[0], o.IntRange[1])
			sample.Flt = rnd.Float64(o.FltRange[0], o.FltRange[1])
		} else {
			sample.Int = o.IntRange[0] + idx*(o.IntRange[1]-o.IntRange[0])/(num-1)
			sample.Flt = o.FltRange[0] + float64(idx)*(o.FltRange[1]-o.FltRange[0])/float64(num-1)
		}
	}

	// float point number
	if len(o.FltRange) == 2 {
		if idx < 0 || num < 2 {
			sample.Int = rnd.Int(o.IntRange[0], o.IntRange[1])
			sample.Flt = rnd.Float64(o.FltRange[0], o.FltRange[1])
		} else {
			sample.Int = o.IntRange[0] + idx*(o.IntRange[1]-o.IntRange[0])/(num-1)
			sample.Flt = o.FltRange[0] + float64(idx)*(o.FltRange[1]-o.FltRange[0])/float64(num-1)
		}
	}

	// sizes of slices
	nw := len(o.PoolWords)
	nb := len(o.PoolBytes)
	nt := len(o.PoolBtxt)
	nf := len(o.PoolFuncs)

	// string
	if nw > 0 {
		sample.String = o.PoolWords[rnd.Int(0, nw-1)]
	}

	// byte
	if nb > 0 {
		sample.Byte = o.PoolBytes[rnd.Int(0, nb-1)]
	}

	// bytes
	if nt > 0 {
		sample.Bytes = []byte(o.PoolBtxt[rnd.Int(0, nt-1)])
	}

	// function
	if nf > 0 {
		sample.Func = o.PoolFuncs[rnd.Int(0, nf-1)]
	}
	return
}

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
	Float  float64 // a float point number
	String string  // a string
	Byte   byte    // a byte
	Bytes  []byte  // a set of bytes
	Func   Func_t  // a function
}

// Bingo collects values to be drawn in random operations
type Bingo struct {
	IntRange   []int     // min and max integers
	FloatRange []float64 // min and max float point numbers
	PoolBytes  []byte    // pool of bytes to be used in Gene.Byte
	PoolWords  []string  // pool of words to be used in Gene.String
	PoolBwords []string  // pool of byte-words to be used in Gene.Bytes
	PoolFuncs  []Func_t  // pool of functions
}

// Init initialises Bingo with template values
func (o *Bingo) Init() {
	o.IntRange = []int{-10, 10}
	o.FloatRange = []float64{-123.0, 321.0}
	o.PoolBytes = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	o.PoolBwords = []string{"apple", "banana", "mango", "orange", "peach", "kiwi"}
	o.PoolFuncs = []Func_t{
		func(g *Gene) string { return io.Sf("myInt=%d", g.GetInt()) },
		func(g *Gene) string { return io.Sf("myFlt=%d", g.GetFloat()) },
		func(g *Gene) string { return io.Sf("myStr=%d", g.GetString()) },
		func(g *Gene) string { return io.Sf("myByt=%d", g.GetByte()) },
		func(g *Gene) string { return io.Sf("myBys=%d", g.GetBytes()) },
		func(g *Gene) string { return io.Sf("myFcs=%d", g.GetFunc()) },
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
	if idx < 0 || num < 2 {
		sample.Int = rnd.Int(o.IntRange[0], o.IntRange[1])
		sample.Float = rnd.Float64(o.FloatRange[0], o.FloatRange[1])
	} else {
		sample.Int = o.IntRange[0] + idx*(o.IntRange[1]-o.IntRange[0])/(num-1)
		sample.Float = o.FloatRange[0] + float64(idx)*(o.FloatRange[1]-o.FloatRange[0])/float64(num-1)
	}

	return
}

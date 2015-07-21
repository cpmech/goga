// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/io"

// Func_tt defines a type for a generic function to be used as a gene value
type Func_tt func(ind *Individual) string

// CxFunc_t defines function to perform the crossover operation
//type CxFunc_t func(a, b, A, B interface{}, cuts []int, prob float64)

// Individual implements one individual in a population
type Individual struct {

	// data
	ObjValue float64 // objective value
	Fitness  float64 // fitness
	Nfloats  int     // number of floats
	Nbases   int     // number of bases to split Floats

	// chromosome
	Ints    []int     // integers
	Floats  []float64 // floats [nfloats * nbases]
	Strings []string  // strings
	Keys    []byte    // 1D bytes
	Bytes   [][]byte  // 2D bytes
	Funcs   []Func_tt // functions
}

// NewIndividual allocates a new individual
//  Input:
//   nbases -- used to split genes of floats into smaller parts
//   slices -- slices of ints, floats, strings, bytes, []bytes, and/or Func_t
//  Notes:
//   1) the slices in 'genes' can all be combined to define genes with mixed data;
//   2) the slices can also be nil, except for one of them.
func NewIndividual(nbases int, slices ...interface{}) (o *Individual) {
	o = new(Individual)
	for _, slice := range slices {
		switch s := slice.(type) {
		case []int:
			o.Ints = make([]int, len(s))
			copy(o.Ints, s)

		case []float64:
			o.Nfloats = len(s)
			o.Nbases = nbases
			if o.Nbases > 1 {
				o.Floats = SimpleChromo(s, nbases)
			} else {
				o.Floats = make([]float64, o.Nfloats*o.Nbases)
				copy(o.Floats, s)
			}

		case []string:
			o.Strings = make([]string, len(s))
			copy(o.Strings, s)

		case []byte:
			o.Keys = make([]byte, len(s))
			copy(o.Keys, s)

		case [][]byte:
			o.Bytes = make([][]byte, len(s))
			for i, x := range s {
				o.Bytes[i] = make([]byte, len(x))
				copy(o.Bytes[i], x)
			}

		case []Func_tt:
			o.Funcs = make([]Func_tt, len(s))
			copy(o.Funcs, s)
		}
	}
	return
}

// GetCopy returns a copy of this individual
func (o Individual) GetCopy() (x *Individual) {

	x = new(Individual)
	x.ObjValue = o.ObjValue
	x.Fitness = o.Fitness
	x.Nfloats = o.Nfloats
	x.Nbases = o.Nbases

	if o.Ints != nil {
		x.Ints = make([]int, len(o.Ints))
		copy(x.Ints, o.Ints)
	}

	if o.Floats != nil {
		x.Floats = make([]float64, len(o.Floats))
		copy(x.Floats, o.Floats)
	}

	if o.Strings != nil {
		x.Strings = make([]string, len(o.Strings))
		copy(x.Strings, o.Strings)
	}

	if o.Keys != nil {
		x.Keys = make([]byte, len(o.Keys))
		copy(x.Keys, o.Keys)
	}

	if o.Bytes != nil {
		x.Bytes = make([][]byte, len(o.Bytes))
		for i, b := range o.Bytes {
			x.Bytes[i] = make([]byte, len(b))
			copy(x.Bytes[i], b)
		}
	}

	if o.Funcs != nil {
		x.Funcs = make([]Func_tt, len(o.Funcs))
		copy(x.Funcs, o.Funcs)
	}
	return
}

// genetic algorithm routines //////////////////////////////////////////////////////////////////////

type IntCxFunc_t func(a, b, A, B []int, cuts []int, pc float64)
type FltCxFunc_t func(a, b, A, B []float64, cuts []int, pc float64)
type StrCxFunc_t func(a, b, A, B []string, cuts []int, pc float64)
type KeyCxFunc_t func(a, b, A, B []byte, cuts []int, pc float64)
type BytCxFunc_t func(a, b, A, B [][]byte, cuts []int, pc float64)
type FunCxFunc_t func(a, b, A, B []Func_tt, cuts []int, pc float64)

// Crossover performs the crossover between chromosomes of two individuals A and B
// resulting in the chromosomes of other two individuals a and b
//  Input:
//   A and B -- parents
//   cuts    -- [6][...] positions for cuts in the augmented/whole chromosome
//              len(cuts) == 6: {int, flt, string, byte, bytes, func}
//              len(cuts[i]) defines the number of cuts. use -1 for random generation of cuts
//              cuts == nil => use random numbers
//   probs   -- [6] probabilities. use nil for default value
//   cxfucns -- [6] crossover functions. use nil for default ones
//  Output:
//   a and b -- offspring
//func Crossover(a, b, A, B *Individual, cuts map[string][]int, probs map[string]float64, intcxf IntCxFunc_t, fltcxf FltCxFunc_t, strcxf StrCxFunc_t, keycxf KeyCxFunc_t, bytcxf BytCxFunc_t, funcxf FunCxFunc_t) {
func Crossover(a, b, A, B *Individual, cuts [][]int, probs []float64, intcxf IntCxFunc_t, fltcxf FltCxFunc_t, strcxf StrCxFunc_t, keycxf KeyCxFunc_t, bytcxf BytCxFunc_t, funcxf FunCxFunc_t) {
	if cuts == nil {
		cuts = [][]int{{-1}, {-1}, {-1}, {-1}, {-1}, {-1}}
	}
	if probs == nil {
		probs = []float64{0.8, 0.8, 0.8, 0.8, 0.8, 0.8}
	}
	if intcxf == nil {
		intcxf = IntCrossover
	}
	if fltcxf == nil {
		fltcxf = FltCrossover
	}
	if strcxf == nil {
		strcxf = StrCrossover
	}
	if keycxf == nil {
		keycxf = KeyCrossover
	}
	if bytcxf == nil {
		bytcxf = BytCrossover
	}
	if funcxf == nil {
		funcxf = FunCrossover
	}
	if A.Ints != nil {
		intcxf(a.Ints, b.Ints, A.Ints, B.Ints, cuts[0], probs[0])
	}
	if A.Floats != nil {
		fltcxf(a.Floats, b.Floats, A.Floats, B.Floats, cuts[1], probs[1])
	}
	if A.Strings != nil {
		strcxf(a.Strings, b.Strings, A.Strings, B.Strings, cuts[2], probs[2])
	}
	if A.Keys != nil {
		keycxf(a.Keys, b.Keys, A.Keys, B.Keys, cuts[3], probs[3])
	}
	if A.Bytes != nil {
		bytcxf(a.Bytes, b.Bytes, A.Bytes, B.Bytes, cuts[4], probs[4])
	}
	if A.Funcs != nil {
		funcxf(a.Funcs, b.Funcs, A.Funcs, B.Funcs, cuts[5], probs[5])
	}
}

// output //////////////////////////////////////////////////////////////////////////////////////////

// GetStringSizes returns the sizes of strings representing each gene type
//  sizes -- [6][...] sizes of strings for {int, flt, string, byte, bytes, func}
func (o Individual) GetStringSizes() (sizes [][]int) {

	sizes = make([][]int, 6)
	if o.Ints != nil {
		sizes[0] = make([]int, len(o.Ints))
		for i, x := range o.Ints {
			sizes[0][i] = imax(sizes[0][i], len(io.Sf("%v", x)))
		}
	}

	if o.Floats != nil {
		sizes[1] = make([]int, len(o.Floats))
		for i, x := range o.Floats {
			if o.Nbases > 1 {
				for j := 0; j < o.Nbases; j++ {
					x += o.Floats[i*o.Nbases+j]
				}
			}
			sizes[1][i] = imax(sizes[1][i], len(io.Sf("%v", x)))
		}
	}

	if o.Strings != nil {
		sizes[2] = make([]int, len(o.Strings))
		for i, x := range o.Strings {
			sizes[2][i] = imax(sizes[2][i], len(io.Sf("%v", x)))
		}
	}

	if o.Keys != nil {
		sizes[3] = make([]int, len(o.Keys))
		for i, x := range o.Keys {
			sizes[3][i] = imax(sizes[3][i], len(io.Sf("%v", x)))
		}
	}

	if o.Bytes != nil {
		sizes[4] = make([]int, len(o.Bytes))
		for i, x := range o.Bytes {
			sizes[4][i] = imax(sizes[4][i], len(io.Sf("%v", string(x))))
		}
	}

	if o.Funcs != nil {
		sizes[5] = make([]int, len(o.Funcs))
		for i, x := range o.Funcs {
			sizes[5][i] = imax(sizes[5][i], len(io.Sf("%v", x)))
		}
	}
	return
}

// Output returns a string representation of this individual
//  fmts -- [6][...] formats of strings for {int, flt, string, byte, bytes, func}
//          use fmts == nil to choose default ones
func (o *Individual) Output(fmts [][]string) (l string) {

	if fmts == nil {
		fmts = [][]string{{" %d"}, {" %g"}, {" %q"}, {" %x"}, {" %q"}, {" %q"}}
	}

	fmt := func(itype, idx int) (s string) {
		s = fmts[itype][0]
		if idx < len(fmts[itype]) {
			s = fmts[itype][idx]
		}
		return
	}

	for i, x := range o.Ints {
		l += io.Sf(fmt(0, i), x)
	}

	for i := 0; i < o.Nfloats; i++ {
		x := o.Floats[i]
		if o.Nbases > 1 {
			x = 0
			for j := 0; j < o.Nbases; j++ {
				x += o.Floats[i*o.Nbases+j]
			}
		}
		l += io.Sf(fmt(1, i), x)
	}

	for i, x := range o.Strings {
		l += io.Sf(fmt(2, i), x)
	}

	for i, x := range o.Keys {
		l += io.Sf(fmt(3, i), x)
	}

	for i, x := range o.Bytes {
		l += io.Sf(fmt(4, i), string(x))
	}

	for i, x := range o.Funcs {
		l += io.Sf(fmt(5, i), x(o))
	}

	return
}

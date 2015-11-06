// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Individual implements one individual in a population
type Individual struct {

	// data
	Ovas      []float64 // objective values
	Oors      []float64 // out-of-range values: sum of positive distances from constraints
	Demerit   float64   // quantity for comparing individuals. 0=good 1=bad 2=worse(oor) 3=worst(oor)
	Nfltgenes int       // number of floats == number of float64 genes
	Nbases    int       // number of bases to split Floats

	// chromosome
	Ints    []int     // integers
	Floats  []float64 // floats [nFLTgenes * nbases]
	Strings []string  // strings
	Keys    []byte    // 1D bytes
	Bytes   [][]byte  // 2D bytes
	Funcs   []Func_t  // functions
}

// NewIndividual allocates a new individual
//  Input:
//   nbases -- used to split genes of floats into smaller parts
//   slices -- slices of ints, floats, strings, bytes, []bytes, and/or Func_t
//  Notes:
//   1) the slices in 'genes' can all be combined to define genes with mixed data;
//   2) the slices can also be nil, except for one of them.
func NewIndividual(nova, noor, nbases int, slices ...interface{}) (o *Individual) {
	o = new(Individual)
	o.Ovas = make([]float64, nova)
	o.Oors = make([]float64, noor)
	for _, slice := range slices {
		switch s := slice.(type) {
		case []int:
			o.Ints = make([]int, len(s))
			copy(o.Ints, s)

		case []float64:
			o.Nfltgenes = len(s)
			o.Nbases = nbases
			if o.Nbases > 1 {
				o.Floats = SimpleChromo(s, nbases)
			} else {
				o.Floats = make([]float64, o.Nfltgenes*o.Nbases)
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

		case []Func_t:
			o.Funcs = make([]Func_t, len(s))
			copy(o.Funcs, s)
		}
	}
	return
}

// GetCopy returns a copy of this individual
func (o Individual) GetCopy() (x *Individual) {

	x = new(Individual)
	x.Ovas = make([]float64, len(o.Ovas))
	x.Oors = make([]float64, len(o.Oors))
	copy(x.Ovas, o.Ovas)
	copy(x.Oors, o.Oors)
	x.Demerit = o.Demerit
	x.Nfltgenes = o.Nfltgenes
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
		x.Funcs = make([]Func_t, len(o.Funcs))
		copy(x.Funcs, o.Funcs)
	}
	return
}

// CopyInto copies this individual's data into another individual
func (o Individual) CopyInto(x *Individual) {

	copy(x.Ovas, o.Ovas)
	copy(x.Oors, o.Oors)
	x.Demerit = o.Demerit
	x.Nfltgenes = o.Nfltgenes
	x.Nbases = o.Nbases

	if o.Ints != nil {
		copy(x.Ints, o.Ints)
	}

	if o.Floats != nil {
		copy(x.Floats, o.Floats)
	}

	if o.Strings != nil {
		copy(x.Strings, o.Strings)
	}

	if o.Keys != nil {
		copy(x.Keys, o.Keys)
	}

	if o.Bytes != nil {
		for i, b := range o.Bytes {
			copy(x.Bytes[i], b)
		}
	}

	if o.Funcs != nil {
		copy(x.Funcs, o.Funcs)
	}
	return
}

// Feasible returns whether this individual is feasible or not
func (o Individual) Feasible() bool {
	for _, oor := range o.Oors {
		if oor > 0 {
			return false
		}
	}
	return true
}

// IndCompareDet compares individual 'A' with another one 'B'. Deterministic method
func IndCompareDet(A, B *Individual) (A_dominates, B_dominates bool) {
	var A_nviolations, B_nviolations int
	for i := 0; i < len(A.Oors); i++ {
		if A.Oors[i] > 0 {
			A_nviolations++
		}
		if B.Oors[i] > 0 {
			B_nviolations++
		}
	}
	if A_nviolations > 0 {
		if B_nviolations > 0 {
			if A_nviolations < B_nviolations {
				A_dominates = true
				return
			}
			if B_nviolations < A_nviolations {
				B_dominates = true
				return
			}
			A_dominates, B_dominates = utl.DblsParetoMin(A.Oors, B.Oors)
			if !A_dominates && !B_dominates {
				A_dominates, B_dominates = utl.DblsParetoMin(A.Ovas, B.Ovas)
			}
			return
		}
		B_dominates = true
		return
	}
	if B_nviolations > 0 {
		A_dominates = true
		return
	}
	A_dominates, B_dominates = utl.DblsParetoMin(A.Ovas, B.Ovas)
	return
}

// IndCompareProb compares individual 'A' with another one 'B' using probabilistic Pareto method
func IndCompareProb(A, B *Individual, φ float64) (A_dominates bool) {
	var A_nviolations, B_nviolations int
	for i := 0; i < len(A.Oors); i++ {
		if A.Oors[i] > 0 {
			A_nviolations++
		}
		if B.Oors[i] > 0 {
			B_nviolations++
		}
	}
	if A_nviolations > 0 {
		if B_nviolations > 0 {
			if A_nviolations < B_nviolations {
				A_dominates = true
				return
			}
			if B_nviolations < A_nviolations {
				A_dominates = false
				return
			}
			var B_dominates bool
			A_dominates, B_dominates = utl.DblsParetoMin(A.Oors, B.Oors)
			if !A_dominates && !B_dominates {
				A_dominates = utl.DblsParetoMinProb(A.Ovas, B.Ovas, φ)
			}
			return
		}
		A_dominates = false
		return
	}
	if B_nviolations > 0 {
		A_dominates = true
		return
	}
	A_dominates = utl.DblsParetoMinProb(A.Ovas, B.Ovas, φ)
	return
}

// IndDistance computes a distance measure from individual 'A' to another individual 'B'
func IndDistance(A, B *Individual, imin, imax []int, fmin, fmax []float64, ovspace bool) (dist float64) {
	if ovspace {
		for i := 0; i < len(A.Ovas); i++ {
			dist += math.Pow(A.Ovas[i]-B.Ovas[i], 2.0)
		}
		dist = math.Sqrt(dist)
		return
	}
	nints := len(A.Ints)
	dints := 0.0
	for i := 0; i < nints; i++ {
		dints += math.Pow(float64(A.Ints[i]-B.Ints[i])/(1e-15+float64(imax[i]-imin[i])), 2.0)
	}
	if nints > 0 {
		dints = math.Sqrt(dints / float64(nints))
	}
	nflts := len(A.Floats)
	dflts := 0.0
	for i := 0; i < nflts; i++ {
		dflts += math.Pow((A.Floats[i]-B.Floats[i])/(1e-15+fmax[i]-fmin[i]), 2.0)
	}
	if nflts > 0 {
		dflts = math.Sqrt(dflts / float64(nflts))
	}
	return dints + dflts
}

// genetic algorithm routines //////////////////////////////////////////////////////////////////////

// IndCrossover performs the crossover between chromosomes of two individuals A and B
// resulting in the chromosomes of other two individuals a and b
func IndCrossover(a, b, A, B, C, D *Individual, time int, ops *OpsData) {
	if A.Ints != nil {
		if D == nil {
			ops.CxInt(a.Ints, b.Ints, A.Ints, B.Ints, nil, nil, time, ops)
		} else {
			ops.CxInt(a.Ints, b.Ints, A.Ints, B.Ints, C.Ints, D.Ints, time, ops)
		}
	}
	if A.Floats != nil {
		if D == nil {
			ops.CxFlt(a.Floats, b.Floats, A.Floats, B.Floats, nil, nil, time, ops)
		} else {
			ops.CxFlt(a.Floats, b.Floats, A.Floats, B.Floats, C.Floats, D.Floats, time, ops)
		}
	}
	if A.Strings != nil {
		if D == nil {
			ops.CxStr(a.Strings, b.Strings, A.Strings, B.Strings, nil, nil, time, ops)
		} else {
			ops.CxStr(a.Strings, b.Strings, A.Strings, B.Strings, C.Strings, D.Strings, time, ops)
		}
	}
	if A.Keys != nil {
		if D == nil {
			ops.CxKey(a.Keys, b.Keys, A.Keys, B.Keys, nil, nil, time, ops)
		} else {
			ops.CxKey(a.Keys, b.Keys, A.Keys, B.Keys, C.Keys, D.Keys, time, ops)
		}
	}
	if A.Bytes != nil {
		if D == nil {
			ops.CxByt(a.Bytes, b.Bytes, A.Bytes, B.Bytes, nil, nil, time, ops)
		} else {
			ops.CxByt(a.Bytes, b.Bytes, A.Bytes, B.Bytes, C.Bytes, D.Bytes, time, ops)
		}
	}
	if A.Funcs != nil {
		if D == nil {
			ops.CxFun(a.Funcs, b.Funcs, A.Funcs, B.Funcs, nil, nil, time, ops)
		} else {
			ops.CxFun(a.Funcs, b.Funcs, A.Funcs, B.Funcs, C.Funcs, D.Funcs, time, ops)
		}
	}
}

// IndMutation performs the mutation operation in the chromosomes of an individual
func IndMutation(A *Individual, time int, ops *OpsData) {
	if A.Ints != nil {
		ops.MtInt(A.Ints, time, ops)
	}
	if A.Floats != nil {
		ops.MtFlt(A.Floats, time, ops)
	}
	if A.Strings != nil {
		ops.MtStr(A.Strings, time, ops)
	}
	if A.Keys != nil {
		ops.MtKey(A.Keys, time, ops)
	}
	if A.Bytes != nil {
		ops.MtByt(A.Bytes, time, ops)
	}
	if A.Funcs != nil {
		ops.MtFun(A.Funcs, time, ops)
	}
}

// handle bases ////////////////////////////////////////////////////////////////////////////////////

// SetFloat returns the float corresponding to gene 'i'
//  igene -- is the index of gene/float in [0, Nfltgenes]
func (o *Individual) SetFloat(igene int, x float64) {
	if o.Nbases > 1 {
		values := make([]float64, o.Nbases)
		rnd.Float64s(values, 0, 1)
		sum := la.VecAccum(values)
		for j := 0; j < o.Nbases; j++ {
			o.Floats[igene*o.Nbases+j] = x * values[j] / sum
		}
		return
	}
	o.Floats[igene] = x
}

// GetFloat returns the float corresponding to gene 'i'
//  igene -- is the index of gene/float in [0, Nfltgenes]
func (o Individual) GetFloat(igene int) (x float64) {
	if o.Nbases > 1 {
		for j := 0; j < o.Nbases; j++ {
			x += o.Floats[igene*o.Nbases+j]
		}
		return
	}
	return o.Floats[igene]
}

// GetFloats returns all float genes
func (o Individual) GetFloats() (x []float64) {
	x = make([]float64, o.Nfltgenes)
	for i := 0; i < o.Nfltgenes; i++ {
		x[i] = o.GetFloat(i)
	}
	return
}

// output //////////////////////////////////////////////////////////////////////////////////////////

// GetStringSizes returns the sizes of strings representing each gene type
//  sizes -- [6][...] sizes of strings for {int, flt, string, byte, bytes, func}
func (o *Individual) GetStringSizes() (sizes [][]int) {

	sizes = make([][]int, 6)
	if o.Ints != nil {
		sizes[0] = make([]int, len(o.Ints))
		for i, x := range o.Ints {
			sizes[0][i] = utl.Imax(sizes[0][i], len(io.Sf("%v", x)))
		}
	}

	if o.Floats != nil {
		sizes[1] = make([]int, o.Nfltgenes)
		for i := 0; i < o.Nfltgenes; i++ {
			x := o.Floats[i]
			if o.Nbases > 1 {
				x = 0
				for j := 0; j < o.Nbases; j++ {
					x += o.Floats[i*o.Nbases+j]
				}
			}
			sizes[1][i] = utl.Imax(sizes[1][i], len(io.Sf("%v", x)))
		}
	}

	if o.Strings != nil {
		sizes[2] = make([]int, len(o.Strings))
		for i, x := range o.Strings {
			sizes[2][i] = utl.Imax(sizes[2][i], len(io.Sf("%v", x)))
		}
	}

	if o.Keys != nil {
		sizes[3] = make([]int, len(o.Keys))
		for i, x := range o.Keys {
			sizes[3][i] = utl.Imax(sizes[3][i], len(io.Sf("%v", x)))
		}
	}

	if o.Bytes != nil {
		sizes[4] = make([]int, len(o.Bytes))
		for i, x := range o.Bytes {
			sizes[4][i] = utl.Imax(sizes[4][i], len(io.Sf("%v", string(x))))
		}
	}

	if o.Funcs != nil {
		sizes[5] = make([]int, len(o.Funcs))
		for i, x := range o.Funcs {
			sizes[5][i] = utl.Imax(sizes[5][i], len(io.Sf("%v", x(o))))
		}
	}
	return
}

// Output returns a string representation of this individual
//  fmts      -- ["int","flt","str","key","byt","fun"][ngenes] print formats for each gene
//               use fmts == nil to choose default ones
//  showBases -- show bases, if any
func (o *Individual) Output(fmts map[string][]string, showBases bool) (l string) {

	if fmts == nil {
		fmts = map[string][]string{"int": {" %d"}, "flt": {" %g"}, "str": {" %q"}, "key": {" %x"}, "byt": {" %q"}, "fun": {" %q"}}
	}

	fmt := func(name string, idx int) (s string) {
		s = fmts[name][0]
		if idx < len(fmts[name]) {
			s = fmts[name][idx]
		}
		return
	}

	for i, x := range o.Ints {
		l += io.Sf(fmt("int", i), x)
	}

	for i := 0; i < o.Nfltgenes; i++ {
		x := o.Floats[i]
		if o.Nbases > 1 {
			x = 0
			for j := 0; j < o.Nbases; j++ {
				x += o.Floats[i*o.Nbases+j]
			}
		}
		l += io.Sf(fmt("flt", i), x)
	}

	for i, x := range o.Strings {
		l += io.Sf(fmt("str", i), x)
	}

	for i, x := range o.Keys {
		l += io.Sf(fmt("key", i), x)
	}

	for i, x := range o.Bytes {
		l += io.Sf(fmt("byt", i), string(x))
	}

	for i, x := range o.Funcs {
		l += io.Sf(fmt("fun", i), x(o))
	}

	if showBases && len(o.Floats) > 0 {
		for i, x := range o.Floats {
			if i%o.Nbases == 0 {
				if i == 0 {
					l += " ||"
				} else {
					l += " |"
				}
			}
			l += io.Sf("%11.3e", x)
		}
	}
	return
}

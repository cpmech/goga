// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
)

// OpsData holds data for crossover and mutation operators
type OpsData struct {

	// input
	IntPc     float64 // probability of crossover for ints
	IntPm     float64 // probability of mutation for ints
	Ncuts     int     // number of cuts during crossover
	Nchanges  int     // number of changes during mutation
	Mmax      float64 // multiplier for mutation
	Cuts      []int   // specified cuts for crossover. can be <nil>
	OrdSti    []int   // {start, end, insertPoint}. can be <nil>
	DebEtac   float64 // Deb's SBX crossover parameter
	DebEtam   float64 // Deb's parameter-based mutation parameter
	DEpc      float64 // differential-evolution crossover probability
	DEmult    float64 // differential-evolution multiplier
	FltCxName string  // crossover function name. ""=default; "mw"=BLX-α evolution; "db"=Deb's SBX; "de"=differential; "cl"=classic
	FltMtName string  // mutation function name. ""=default; "mw"=Michaelewicz; "db"=Deb's parameter-based

	// derived
	Use4inds bool        // crossover needs 4 individuals (A,B,C,D); e.g. with differential evolution (de)
	Tmax     float64     // max number of generations
	Xrange   [][]float64 // [ngenes][2] genes minimum and maximum values

	// crossover functions
	CxInt CxIntFunc_t // int crossover function
	CxFlt CxFltFunc_t // flt crossover function

	// mutation functions
	MtInt MtIntFunc_t // int mutation function
	MtFlt MtFltFunc_t // flt mutation function
}

// SetDefault sets default values
func (o *OpsData) SetDefault() {

	// input
	o.IntPc = 0.8
	o.IntPm = 0.01
	o.Ncuts = 2
	o.Nchanges = 1
	o.Mmax = 2
	o.DebEtac = 1
	o.DebEtam = 1
	o.DEpc = 0.1
	o.DEmult = 0.5

	// crossover functions
	o.CxInt = IntCrossover
	o.CxFlt = FltCrossoverDB

	// mutation functions
	o.MtInt = IntMutation
	o.MtFlt = FltMutationDB
}

// CalcDerived sets derived quantities
func (o *OpsData) CalcDerived(Tf int, xrange [][]float64) {
	o.Tmax = float64(Tf)
	o.Xrange = xrange
	switch o.FltCxName {
	case "db":
		o.CxFlt = FltCrossoverDB
	case "de":
		o.CxFlt = FltCrossoverDE
		o.Use4inds = true
	}
	switch o.FltMtName {
	case "db":
		o.MtFlt = FltMutationDB
	}
}

// EnforceRange makes sure x is within given range
func (o *OpsData) EnforceRange(igene int, x float64) float64 {
	if x < o.Xrange[igene][0] {
		return o.Xrange[igene][0]
	}
	if x > o.Xrange[igene][1] {
		return o.Xrange[igene][1]
	}
	return x
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// GenerateCxEnds randomly computes the end positions of cuts in chromosomes
//  Input:
//   size  -- size of chromosome
//   ncuts -- number of cuts to be used, unless cuts != nil
//   cuts  -- cut positions. can be nil => use ncuts instead
//  Output:
//   ends -- end positions where the last one equals size
//  Example:
//        0 1 2 3 4 5 6 7
//    A = a b c d e f g h    size = 8
//         ↑       ↑     ↑   cuts = [1, 5]
//         1       5     8   ends = [1, 5, 8]
func GenerateCxEnds(size, ncuts int, cuts []int) (ends []int) {

	// handle small slices
	if size < 2 {
		return
	}
	if size == 2 {
		return []int{1, size}
	}

	// cuts slice is given
	if len(cuts) > 0 {
		ncuts = len(cuts)
		ends = make([]int, ncuts+1)
		ends[ncuts] = size
		for i, cut := range cuts {
			if cut < 1 || cut >= size {
				chk.Panic("cut=%d is outside the allowed range: 1 ≤ cut ≤ size-1", cut)
			}
			if i > 0 {
				if cut == cuts[i-1] {
					chk.Panic("repeated cut values are not allowed: cuts=%v", cuts)
				}
			}
			ends[i] = cut
		}
		sort.Ints(ends)
		return
	}

	// randomly generate cuts
	if ncuts < 1 {
		ncuts = 1
	}
	if ncuts >= size {
		ncuts = size - 1
	}
	ends = make([]int, ncuts+1)
	ends[ncuts] = size

	// pool of values for selections
	pool := rnd.IntGetUniqueN(1, size, ncuts)
	sort.Ints(pool)
	for i := 0; i < ncuts; i++ {
		ends[i] = pool[i]
	}
	return
}

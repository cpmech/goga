// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
)

// OpsData holds data for crossover and mutation operators
type OpsData struct {

	// input
	Pc        float64 // probability of crossover
	Pm        float64 // probability of mutation
	Ncuts     int     // number of cuts during crossover
	Nchanges  int     // number of changes during mutation
	MwiczB    float64 // Michalewicz' power coefficient
	BlxAlp    float64 // BLX-α coefficient
	Mmax      float64 // multiplier for mutation
	Cuts      []int   // specified cuts for crossover. can be <nil>
	OrdSti    []int   // {start, end, insertPoint}. can be <nil>
	EnfRange  bool    // do enforce range
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
	CxStr CxStrFunc_t // str crossover function
	CxKey CxKeyFunc_t // key crossover function
	CxByt CxBytFunc_t // byt crossover function
	CxFun CxFunFunc_t // fun crossover function

	// mutation functions
	MtInt MtIntFunc_t // int mutation function
	MtFlt MtFltFunc_t // flt mutation function
	MtStr MtStrFunc_t // str mutation function
	MtKey MtKeyFunc_t // key mutation function
	MtByt MtBytFunc_t // byt mutation function
	MtFun MtFunFunc_t // fun mutation function
}

// SetDefault sets default values
func (o *OpsData) SetDefault() {

	// input
	o.Pc = 0.8
	o.Pm = 0.01
	o.Ncuts = 2
	o.Nchanges = 1
	o.MwiczB = 2.0
	o.BlxAlp = 0.5
	o.Mmax = 2
	o.EnfRange = true
	o.DebEtac = 1
	o.DebEtam = 1
	o.DEpc = 0.1
	o.DEmult = 0.5

	// crossover functions
	o.CxInt = IntCrossover
	o.CxFlt = FltCrossoverDB
	o.CxStr = StrCrossover
	o.CxKey = KeyCrossover
	o.CxByt = BytCrossover
	o.CxFun = FunCrossover

	// mutation functions
	o.MtInt = IntMutation
	o.MtFlt = FltMutationDB
	o.MtStr = StrMutation
	o.MtKey = KeyMutation
	o.MtByt = BytMutation
	o.MtFun = FunMutation
}

// CalcDerived sets derived quantities
func (o *OpsData) CalcDerived(Tf int, xrange [][]float64) {
	o.Tmax = float64(Tf)
	o.Xrange = xrange
	switch o.FltCxName {
	case "mw":
		o.CxFlt = FltCrossoverMW
	case "db":
		o.CxFlt = FltCrossoverDB
	case "de":
		o.CxFlt = FltCrossoverDE
		o.Use4inds = true
	case "mix":
		o.CxFlt = FltCrossoverMix
		o.Use4inds = true
	case "cl":
		o.CxFlt = FltCrossover
	}
	switch o.FltMtName {
	case "mw":
		o.MtFlt = FltMutationMW
	case "db":
		o.MtFlt = FltMutationDB
	case "cl":
		o.MtFlt = FltMutation
	}
}

// MwiczDelta computes Michalewicz' Δ function
func (o *OpsData) MwiczDelta(t, x float64) float64 {
	r := rand.Float64()
	return (1.0 - math.Pow(r, math.Pow(1.0-t/o.Tmax, o.MwiczB))) * x
}

// EnforceRange makes sure x is within given range
func (o *OpsData) EnforceRange(igene int, x float64) float64 {
	if !o.EnfRange {
		return x
	}
	if x < o.Xrange[igene][0] {
		return o.Xrange[igene][0]
	}
	if x > o.Xrange[igene][1] {
		return o.Xrange[igene][1]
	}
	return x
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// Report generates report
func (o *OpsData) Report(buf *bytes.Buffer) {
	io.Ff(buf, `
# input
Pc        = %v # probability of crossover
Pm        = %v # probability of mutation
Ncuts     = %v # number of cuts during crossover
Nchanges  = %v # number of changes during mutation
MwiczB    = %v # Michalewicz' power coefficient
BlxAlp    = %v # BLX-α coefficient
Mmax      = %v # multiplier for mutation
Cuts      = %v # specified cuts for crossover. can be <nil>
OrdSti    = %v # {start, end, insertPoint}. can be <nil>
EnfRange  = %v # do enforce range
DebEtac   = %v # Deb's SBX crossover parameter
DebEtam   = %v # Deb's parameter-based mutation parameter
DEpc      = %v # differential-evolution crossover probability
DEmult    = %v # differential-evolution multiplier
FltCxName = %q # crossover function name. ""=default; "mw"=BLX-α evolution; "db"=Deb's SBX; "de"=differential; "cl"=classic
FltMtName = %q # mutation function name. ""=default; "mw"=Michaelewicz; "db"=Deb's parameter-based

# derived
Use4inds = %v # crossover needs 4 individuals (A,B,C,D); e.g. with differential evolution (de)
Tmax     = %v # max number of generations
Xrange   = %v # [ngenes][2] genes minimum and maximum values

# crossover functions
CxInt = %v # int crossover function
CxFlt = %v # flt crossover function
CxStr = %v # str crossover function
CxKey = %v # key crossover function
CxByt = %v # byt crossover function
CxFun = %v # fun crossover function

# mutation functions
MtInt = %v # int mutation function
MtFlt = %v # flt mutation function
MtStr = %v # str mutation function
MtKey = %v # key mutation function
MtByt = %v # byt mutation function
MtFun = %v # fun mutation function
`, o.Pc, o.Pm, o.Ncuts, o.Nchanges, o.MwiczB, o.BlxAlp, o.Mmax, o.Cuts, o.OrdSti, o.EnfRange,
		o.DebEtac, o.DebEtam, o.DEpc, o.DEmult, o.FltCxName, o.FltMtName,
		o.Use4inds, o.Tmax, o.Xrange,
		chk.GetFunctionName(o.CxInt), chk.GetFunctionName(o.CxFlt), chk.GetFunctionName(o.CxStr),
		chk.GetFunctionName(o.CxKey), chk.GetFunctionName(o.CxByt), chk.GetFunctionName(o.CxFun),
		chk.GetFunctionName(o.MtInt), chk.GetFunctionName(o.MtFlt), chk.GetFunctionName(o.MtStr),
		chk.GetFunctionName(o.MtKey), chk.GetFunctionName(o.MtByt), chk.GetFunctionName(o.MtFun))
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// SimpleChromo splits 'genes' into 'nbases' unequal parts
//  Input:
//    genes  -- a slice whose size equals to the number of genes
//    nbases -- number of bases used to split 'genes'
//  Output:
//    chromo -- the chromosome
//
//  Example:
//
//    genes = [0, 1, 2, ... nbases-1,  0, 1, 2, ... nbases-1]
//             \___________________/   \___________________/
//                    gene # 0               gene # 1
//
func SimpleChromo(genes []float64, nbases int) (chromo []float64) {
	ngenes := len(genes)
	chromo = make([]float64, ngenes*nbases)
	values := make([]float64, nbases)
	var sumv float64
	for i, g := range genes {
		rnd.Float64s(values, 0, 1)
		sumv = la.VecAccum(values)
		for j := 0; j < nbases; j++ {
			chromo[i*nbases+j] = g * values[j] / sumv
		}
	}
	return
}

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

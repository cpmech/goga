// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"math/rand"
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// OpsData holds data for crossover and mutation operators
type OpsData struct {

	// constants
	Pc       float64     // probability of crossover
	Pm       float64     // probability of mutation
	Ncuts    int         // number of cuts during crossover
	Nchanges int         // number of changes during mutation
	Tmax     float64     // max number of generations
	MwiczB   float64     // Michalewicz' power coefficient
	BlxAlp   float64     // BLX-α coefficient
	Mmax     float64     // multiplier for mutation
	Cuts     []int       // specified cuts for crossover. can be <nil>
	OrdSti   []int       // {start, end, insertPoint}. can be <nil>
	Xrange   [][]float64 // [ngenes][2] genes minimum and maximum values
	EnfRange bool        // do enforce range
	DebEtac  float64     // Deb's SBX crossover parameter
	DebEtam  float64     // Deb's parameter-based mutation parameter

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

	// constants
	o.Pc = 0.8
	o.Pm = 0.01
	o.Ncuts = 2
	o.Nchanges = 1
	o.MwiczB = 2.0
	o.BlxAlp = 0.5
	o.Mmax = 2
	o.EnfRange = true
	o.DebEtac = 1
	o.DebEtam = 100

	// crossover functions
	o.CxInt = IntCrossover
	o.CxFlt = FltCrossoverDeb
	o.CxStr = StrCrossover
	o.CxKey = KeyCrossover
	o.CxByt = BytCrossover
	o.CxFun = FunCrossover

	// mutation functions
	o.MtInt = IntMutation
	o.MtFlt = FltMutationDeb
	o.MtStr = StrMutation
	o.MtKey = KeyMutation
	o.MtByt = BytMutation
	o.MtFun = FunMutation
}

// CalcDerived sets derived quantities
func (o *OpsData) CalcDerived(Tf int, xrange [][]float64) {
	o.Tmax = float64(Tf)
	o.Xrange = xrange
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

// crossover ///////////////////////////////////////////////////////////////////////////////////////

// IntCrossover performs the crossover of genetic data from A and B
//  Output:
//   a and b -- offspring
//  Example:
//         0 1 2 3 4 5 6 7
//     A = a b c d e f g h    size = 8
//     B = * . . . . * * *    cuts = [1, 5]
//          ↑       ↑     ↑   ends = [1, 5, 8]
//          1       5     8
//     a = a . . . . f g h
//     b = * b c d e * * *
func IntCrossover(a, b, A, B []int, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ops.Ncuts, ops.Cuts)
	swap := false
	start := 0
	for _, end := range ends {
		if swap {
			for j := start; j < end; j++ {
				b[j], a[j] = A[j], B[j]
			}
		} else {
			for j := start; j < end; j++ {
				a[j], b[j] = A[j], B[j]
			}
		}
		start = end
		swap = !swap
	}
	return
}

// IntOrdCrossover performs the crossover in a pair of individuals with integer numbers
// that correspond to a ordered sequence, e.g. for traveling salesman problem
//  Output:
//    a and b -- offspring chromosomes
//  Note: using OX1 method explained in [1] (proposed in [2])
//  References:
//   [1] Larrañaga P, Kuijpers CMH, Murga RH, Inza I and Dizdarevic S. Genetic Algorithms for the
//       Travelling Salesman Problem: A Review of Representations and Operators. Artificial
//       Intelligence Review, 13:129-170; 1999. doi:10.1023/A:1006529012972
//   [2] Davis L. Applying Adaptive Algorithms to Epistatic Domains. Proceedings of International
//       Joint Conference on Artificial Intelligence, 162-164; 1985.
//  Example:
//   data:
//         0 1   2 3 4   5 6 7
//     A = a b | c d e | f g h        size = 8
//     B = b d | f h g | e c a        cuts = [2, 5]
//             ↑       ↑       ↑      ends = [2, 5, 8]
//             2       5       8
//   first step: copy subtours
//     a = . . | f h g | . . .
//     b = . . | c d e | . . .
//   second step: copy unique from subtour's end, position 5
//               start adding here
//                       ↓                           5 6 7   0 1   2 3 4
//     a = d e | f h g | a b c         get from A: | f̶ g̶ h̶ | a b | c d e
//     b = h g | c d e | a b f         get from B: | e̶ c̶ a | b d̶ | f h g
func IntOrdCrossover(a, b, A, B []int, time int, ops *OpsData) (notused []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 3 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	var s, t int
	if len(ops.Cuts) == 2 {
		s, t = ops.Cuts[0], ops.Cuts[1]
	} else {
		s = rnd.Int(1, size-2)
		t = rnd.Int(s+1, size-1)
	}
	chk.IntAssertLessThan(s, t)
	acore := B[s:t]
	bcore := A[s:t]
	ncore := t - s
	acorehas := make(map[int]bool) // TODO: check if map can be replaced => improve efficiency
	bcorehas := make(map[int]bool)
	for i := 0; i < ncore; i++ {
		a[s+i] = acore[i]
		b[s+i] = bcore[i]
		acorehas[acore[i]] = true
		bcorehas[bcore[i]] = true
	}
	ja, jb := t, t
	for i := 0; i < size; i++ {
		k := (i + t) % size
		if !acorehas[A[k]] {
			a[ja] = A[k]
			ja++
			if ja == size {
				ja = 0
			}
		}
		if !bcorehas[B[k]] {
			b[jb] = B[k]
			jb++
			if jb == size {
				jb = 0
			}
		}
	}
	return
}

// FltCrossover performs the crossover of genetic data from A and B
//  Output:
//   a and b -- offspring
//  Example:
//         0 1 2 3 4 5 6 7
//     A = a b c d e f g h    size = 8
//     B = * . . . . * * *    cuts = [1, 5]
//          ↑       ↑     ↑   ends = [1, 5, 8]
//          1       5     8
//     a = a . . . . f g h
//     b = * b c d e * * *
func FltCrossover(a, b, A, B []float64, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ops.Ncuts, ops.Cuts)
	swap := false
	start := 0
	for _, end := range ends {
		if swap {
			for j := start; j < end; j++ {
				b[j], a[j] = A[j], B[j]
			}
		} else {
			for j := start; j < end; j++ {
				a[j], b[j] = A[j], B[j]
			}
		}
		start = end
		swap = !swap
	}
	return
}

// FltCrossoverBlx implements the BLS-α crossover by Eshelman et al. (1993); see also Herrera (1998)
//  Output:
//   a and b -- offspring
func FltCrossoverBlx(a, b, A, B []float64, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) {
		for i := 0; i < size; i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	α := ops.BlxAlp
	chk.IntAssert(len(ops.Xrange), len(A))
	var cmin, cmax, δ float64
	for i := 0; i < size; i++ {
		cmin = utl.Min(A[i], B[i])
		cmax = utl.Max(A[i], B[i])
		δ = cmax - cmin
		a[i] = rnd.Float64(cmin-α*δ, cmax+α*δ)
		b[i] = rnd.Float64(cmin-α*δ, cmax+α*δ)
		a[i] = ops.EnforceRange(i, a[i])
		b[i] = ops.EnforceRange(i, b[i])
	}
	return
}

// FltCrossoverDeb implements Deb's simulated binary crossover (SBX)
func FltCrossoverDeb(a, b, A, B []float64, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) {
		for i := 0; i < size; i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	cc := 1.0 / (ops.DebEtac + 1.0)
	var u, α, β, βb, x1, x2, xl, xu float64
	for i := 0; i < size; i++ {
		x1, x2 = A[i], B[i]
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		u = rnd.Float64(0, 1)
		if ops.EnfRange {
			xl, xu = ops.Xrange[i][0], ops.Xrange[i][1]
			β = 1.0 + (2.0/(1e-15+x2-x1))*utl.Min(x1-xl, xu-x2)
			α = 2.0 - math.Pow(β, -(ops.DebEtac+1.0))
			if u <= 1.0/α {
				βb = math.Pow(α*u, cc)
			} else {
				βb = math.Pow(1.0/(2.0-α*u), cc)
			}
		} else {
			if u <= 0.5 {
				βb = math.Pow(2.0*u, cc)
			} else {
				βb = math.Pow(0.5/(1.0-u), cc)
			}
		}
		a[i] = ops.EnforceRange(i, 0.5*(x1+x2-βb*(x2-x1)))
		b[i] = ops.EnforceRange(i, 0.5*(x1+x2+βb*(x2-x1)))
	}
	return
}

// StrCrossover performs the crossover of genetic data from A and B
//  Output:
//   a and b -- offspring
//  Example:
//         0 1 2 3 4 5 6 7
//     A = a b c d e f g h    size = 8
//     B = * . . . . * * *    cuts = [1, 5]
//          ↑       ↑     ↑   ends = [1, 5, 8]
//          1       5     8
//     a = a . . . . f g h
//     b = * b c d e * * *
func StrCrossover(a, b, A, B []string, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ops.Ncuts, ops.Cuts)
	swap := false
	start := 0
	for _, end := range ends {
		if swap {
			for j := start; j < end; j++ {
				b[j], a[j] = A[j], B[j]
			}
		} else {
			for j := start; j < end; j++ {
				a[j], b[j] = A[j], B[j]
			}
		}
		start = end
		swap = !swap
	}
	return
}

// KeyCrossover performs the crossover of genetic data from A and B
//  Output:
//   a and b -- offspring
//  Example:
//         0 1 2 3 4 5 6 7
//     A = a b c d e f g h    size = 8
//     B = * . . . . * * *    cuts = [1, 5]
//          ↑       ↑     ↑   ends = [1, 5, 8]
//          1       5     8
//     a = a . . . . f g h
//     b = * b c d e * * *
func KeyCrossover(a, b, A, B []byte, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ops.Ncuts, ops.Cuts)
	swap := false
	start := 0
	for _, end := range ends {
		if swap {
			for j := start; j < end; j++ {
				b[j], a[j] = A[j], B[j]
			}
		} else {
			for j := start; j < end; j++ {
				a[j], b[j] = A[j], B[j]
			}
		}
		start = end
		swap = !swap
	}
	return
}

// BytCrossover performs the crossover of genetic data from A and B
//  Output:
//   a and b -- offspring
//  Example:
//         0 1 2 3 4 5 6 7
//     A = a b c d e f g h    size = 8
//     B = * . . . . * * *    cuts = [1, 5]
//          ↑       ↑     ↑   ends = [1, 5, 8]
//          1       5     8
//     a = a . . . . f g h
//     b = * b c d e * * *
func BytCrossover(a, b, A, B [][]byte, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			copy(a[i], A[i])
			copy(b[i], B[i])
		}
		return
	}
	ends = GenerateCxEnds(size, ops.Ncuts, ops.Cuts)
	swap := false
	start := 0
	for _, end := range ends {
		if swap {
			for j := start; j < end; j++ {
				copy(b[j], A[j])
				copy(a[j], B[j])
			}
		} else {
			for j := start; j < end; j++ {
				copy(a[j], A[j])
				copy(b[j], B[j])
			}
		}
		start = end
		swap = !swap
	}
	return
}

// FunCrossover performs the crossover of genetic data from A and B
//  Output:
//   a and b -- offspring
//  Example:
//         0 1 2 3 4 5 6 7
//     A = a b c d e f g h    size = 8
//     B = * . . . . * * *    cuts = [1, 5]
//          ↑       ↑     ↑   ends = [1, 5, 8]
//          1       5     8
//     a = a . . . . f g h
//     b = * b c d e * * *
func FunCrossover(a, b, A, B []Func_t, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ops.Ncuts, ops.Cuts)
	swap := false
	start := 0
	for _, end := range ends {
		if swap {
			for j := start; j < end; j++ {
				b[j], a[j] = A[j], B[j]
			}
		} else {
			for j := start; j < end; j++ {
				a[j], b[j] = A[j], B[j]
			}
		}
		start = end
		swap = !swap
	}
	return
}

// mutation ////////////////////////////////////////////////////////////////////////////////////////

// IntMutation performs the mutation of genetic data from A
//  Output: modified individual 'A'
func IntMutation(A []int, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		m := rnd.Int(1, int(ops.Mmax))
		if rnd.FlipCoin(0.5) {
			A[i] += m * A[i]
		} else {
			A[i] -= m * A[i]
		}
	}
}

// IntBinMutation performs the mutation of a binary chromosome
//  Output: modified individual 'A'
func IntBinMutation(A []int, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		if A[i] == 0 {
			A[i] = 1
		} else {
			A[i] = 0
		}
	}
}

// IntOrdMutation performs the mutation of genetic data from a ordered list of integers A
//  Output: modified individual 'A'
//  Note: using DM method as explained in [1] (citing [2])
//  References:
//   [1] Larrañaga P, Kuijpers CMH, Murga RH, Inza I and Dizdarevic S. Genetic Algorithms for the
//       Travelling Salesman Problem: A Review of Representations and Operators. Artificial
//       Intelligence Review, 13:129-170; 1999. doi:10.1023/A:1006529012972
//   [2] Michalewicz Z. Genetic Algorithms + Data Structures = Evolution Programs. Berlin
//       Heidelberg: Springer Verlag; 1992
//       Joint Conference on Artificial Intelligence, 162-164; 1985.
//
//  DM displacement mutation method:
//   Ex:
//           0 1 2 3 4 5 6 7
//       A = a b c d e f g h   s = 2
//              ↑     ↑        t = 5
//              2     5
//
//       core = c d e  (subtour)  ncore = t - s = 5 - 2 = 3
//
//                0 1 2 3 4
//       remain = a b f g h  (remaining)  nrem = size - ncore = 8 - 3 = 5
//                       ↑
//                       4 = ins
func IntOrdMutation(A []int, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 3 {
		if size == 2 {
			A[0], A[1] = A[1], A[0]
		}
		return
	}
	var s, t, ncore, nrem, ins int
	if ops.OrdSti != nil {
		s, t, ins = ops.OrdSti[0], ops.OrdSti[1], ops.OrdSti[2]
		ncore = t - s
		nrem = size - ncore
	} else {
		s = rnd.Int(1, size-2)
		t = rnd.Int(s+1, size-1)
		ncore = t - s
		nrem = size - ncore
		ins = rnd.Int(1, nrem)
	}
	core := make([]int, ncore)
	remain := make([]int, nrem)
	var jc, jr int
	for i := 0; i < size; i++ {
		if i >= s && i < t {
			core[jc] = A[i]
			jc++
		} else {
			remain[jr] = A[i]
			jr++
		}
	}
	jc, jr = 0, 0
	for i := 0; i < size; i++ {
		if i < ins {
			A[i] = remain[jr]
			jr++
		} else {
			if jc < ncore {
				A[i] = core[jc]
				jc++
			} else {
				A[i] = remain[jr]
				jr++
			}
		}
	}
}

// FltMutation performs the mutation of genetic data from A
//  Output: modified individual 'A'
func FltMutation(A []float64, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		m := rnd.Float64(1, ops.Mmax)
		if rnd.FlipCoin(0.5) {
			A[i] += m * A[i]
		} else {
			A[i] -= m * A[i]
		}
	}
}

// FltMutationMwicz implements the non-uniform mutation (Michaelewicz, 1992; Herrera, 1998)
// See also Michalewicz (1996) page 103
func FltMutationMwicz(A []float64, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	t := float64(time)
	chk.IntAssert(len(ops.Xrange), len(A))
	for i := 0; i < size; i++ {
		xmin := ops.Xrange[i][0]
		xmax := ops.Xrange[i][1]
		if rnd.FlipCoin(0.5) {
			A[i] += ops.MwiczDelta(t, xmax-A[i])
		} else {
			A[i] -= ops.MwiczDelta(t, A[i]-xmin)
		}
		A[i] = ops.EnforceRange(i, A[i])
	}
}

//  FltMutationDeb implements Deb's parameter-based mutation operator
func FltMutationDeb(A []float64, time int, ops *OpsData) {
	size := len(A)
	chk.IntAssert(len(ops.Xrange), size)
	t := float64(time)
	r := 1.0 / float64(size)
	pm := r + (t/ops.Tmax)*(1.0-r)
	if !rnd.FlipCoin(pm) || size < 1 {
		return
	}
	ηm := 100.0 + t
	cm := 1.0 / (ηm + 1.0)
	var u, Δx, φ, δ, δb, xl, xu float64
	for i := 0; i < size; i++ {
		u = rnd.Float64(0, 1)
		xl, xu = ops.Xrange[i][0], ops.Xrange[i][1]
		Δx = xu - xl
		if ops.EnfRange {
			δ = utl.Min(A[i]-xl, xu-A[i]) / Δx
			φ = math.Pow(1.0-δ, ηm+1.0)
			if u <= 0.5 {
				δb = math.Pow(2.0*u+(1.0-2.0*u)*φ, cm) - 1.0
			} else {
				δb = 1.0 - math.Pow(2.0-2.0*u+(2.0*u-1.0)*φ, cm)
			}
		} else {
			if u <= 0.5 {
				δb = math.Pow(2.0*u, cm) - 1.0
			} else {
				δb = 1.0 - math.Pow(2.0-2.0*u, cm)
			}
		}
		A[i] = ops.EnforceRange(i, A[i]+δb*Δx)
	}
}

// StrMutation performs the mutation of genetic data from A
//  Output: modified individual 'A'
func StrMutation(A []string, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		A[i] = "TODO" // TODO: improve this
	}
}

// KeyMutation performs the mutation of genetic data from A
//  Output: modified individual 'A'
func KeyMutation(A []byte, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		v := rnd.Int(0, 100)
		A[i] = byte(v) // TODO: improve this
	}
}

// BytMutation performs the mutation of genetic data from A
//  Output: modified individual 'A'
func BytMutation(A [][]byte, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		v := rnd.Int(0, 100)
		A[i][0] = byte(v) // TODO: improve this
	}
}

// FunMutation performs the mutation of genetic data from A
//  Output: modified individual 'A'
func FunMutation(A []Func_t, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		// TODO: improve this
		A[i] = func(ind *Individual) string { return "mutated" }
	}
}

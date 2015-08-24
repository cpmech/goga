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

// CalcFitness maps objective values into [0, 1]; thus returning the fitness function values
//  Input:
//    ovs -- objective values
//  Output:
//    f -- fitness function values
func CalcFitness(f, ovs []float64) {
	chk.IntAssert(len(f), len(ovs))
	ymin, ymax := la.VecMinMax(ovs)
	if math.Abs(ymax-ymin) < 1e-14 {
		la.VecFill(f, 1)
		return
	}
	for i := 0; i < len(ovs); i++ {
		f[i] = (ymax - ovs[i]) / (ymax - ymin)
	}
}

// CalcFitnessRanking computes fitness corresponding to a linear ranking
//  Input:
//    ninds -- number of individuals
//    sp    -- selective pressure; must be inside [1, 2]
//  Output:
//    f -- ranked fitnesses
func CalcFitnessRanking(ninds int, sp float64) (f []float64) {
	if sp < 1.0 || sp > 2.0 {
		sp = 1.2
	}
	f = make([]float64, ninds)
	for i := 0; i < ninds; i++ {
		f[i] = 2.0 - sp + 2.0*(sp-1.0)*float64(ninds-i-1)/float64(ninds-1)
	}
	return
}

// selection ///////////////////////////////////////////////////////////////////////////////////////

// RouletteSelect selects n individuals
//  Input:
//    cumprob -- cumulated probabilities (from sorted population)
//    sample  -- a list of random numbers; can be nil
//  Output:
//    selinds -- selected individuals (indices). len(selinds) == nsel
func RouletteSelect(selinds []int, cumprob []float64, sample []float64) {
	nsel := len(selinds)
	chk.IntAssertLessThanOrEqualTo(nsel, len(cumprob))
	if sample == nil {
		var s float64
		for i := 0; i < nsel; i++ {
			s = rand.Float64()
			for j, m := range cumprob {
				if m > s {
					selinds[i] = j
					break
				}
			}
		}
		return
	}
	chk.IntAssert(len(sample), nsel)
	for i, s := range sample {
		for j, m := range cumprob {
			if m > s {
				selinds[i] = j
				break
			}
		}
	}
}

// SUSselect performs the Stochastic-Universal-Sampling selection
//  Input:
//    cumprob -- cumulated probabilities (from sorted population)
//    pb      -- one random number corresponding to the first probability (pointer/position)
//               use pb = -1 to generate a random value here
//  Output:
//    selinds -- selected individuals (indices)
func SUSselect(selinds []int, cumprob []float64, pb float64) {
	nsel := len(selinds)
	chk.IntAssertLessThanOrEqualTo(nsel, len(cumprob))
	dp := 1.0 / float64(nsel)
	if pb < 0 {
		pb = rnd.Float64(0, dp)
	}
	var j int
	for i := 0; i < nsel; i++ {
		j = 0
		for pb > cumprob[j] {
			j += 1
		}
		pb += dp
		selinds[i] = j
	}
}

// FilterPairs generates 2 lists with ninds/2 items each corresponding to selected pairs
// for reprodoction. Repeated indices in pairs are avoided.
//  Input:
//   selinds -- list of selected individuals len(selinds) == ninds
//  Output:
//   A and B -- [ninds/2] lists with pairs
func FilterPairs(A, B []int, selinds []int) {
	ninds := len(selinds)
	chk.IntAssert(len(A), ninds/2)
	chk.IntAssert(len(B), ninds/2)
	var a, b int
	var aux []int
	for i := 0; i < ninds/2; i++ {
		a, b = selinds[2*i], selinds[2*i+1]
		if a == b {
			if len(aux) == 0 {
				aux = rnd.IntGetShuffled(selinds)
			} else {
				rnd.IntShuffle(aux)
			}
			for _, s := range aux {
				if s != a {
					b = s
					break
				}
			}
		}
		A[i], B[i] = a, b
	}
}

// crossover ///////////////////////////////////////////////////////////////////////////////////////

// IntCrossover performs the crossover of genetic data from A and B
//  Input:
//   A and B -- parents
//   ncuts   -- number of cuts to be used, unless cuts != nil
//   cuts    -- cut positions. can be nil => use ncuts instead
//   pc      -- probability of crossover
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
func IntCrossover(a, b, A, B []int, time, ncuts int, cuts []int, pc float64, extra interface{}) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ncuts, cuts)
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
//  Input:
//    A and B -- parents' chromosomes
//    dum     -- not used
//    cuts    -- 2 cut positions. if len(cuts) != 2, 2 cuts are randomly generated
//    pc      -- probability of crossover
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
func IntOrdCrossover(a, b, A, B []int, time, dum int, cuts []int, pc float64, extra interface{}) (notused []int) {
	size := len(A)
	if !rnd.FlipCoin(pc) || size < 3 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	var s, t int
	if len(cuts) == 2 {
		s, t = cuts[0], cuts[1]
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
//  Input:
//   A and B -- parents
//   ncuts   -- number of cuts to be used, unless cuts != nil
//   cuts    -- cut positions. can be nil => use ncuts instead
//   pc      -- probability of crossover
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
func FltCrossover(a, b, A, B []float64, time, ncuts int, cuts []int, pc float64, extra interface{}) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ncuts, cuts)
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
func FltCrossoverBlx(a, b, A, B []float64, time, ncuts int, cuts []int, pc float64, extra interface{}) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(pc) {
		for i := 0; i < size; i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	α := extra.(float64)
	var cmin, cmax, δ float64
	for i := 0; i < size; i++ {
		cmin = utl.Min(A[i], B[i])
		cmax = utl.Max(A[i], B[i])
		δ = cmax - cmin
		a[i] = rnd.Float64(cmin-α*δ, cmax+α*δ)
		b[i] = rnd.Float64(cmin-α*δ, cmax+α*δ)
	}
	return
}

// StrCrossover performs the crossover of genetic data from A and B
//  Input:
//   A and B -- parents
//   ncuts   -- number of cuts to be used, unless cuts != nil
//   cuts    -- cut positions. can be nil => use ncuts instead
//   pc      -- probability of crossover
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
func StrCrossover(a, b, A, B []string, time, ncuts int, cuts []int, pc float64, extra interface{}) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ncuts, cuts)
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
//  Input:
//   A and B -- parents
//   ncuts   -- number of cuts to be used, unless cuts != nil
//   cuts    -- cut positions. can be nil => use ncuts instead
//   pc      -- probability of crossover
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
func KeyCrossover(a, b, A, B []byte, time, ncuts int, cuts []int, pc float64, extra interface{}) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ncuts, cuts)
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
//  Input:
//   A and B -- parents
//   ncuts   -- number of cuts to be used, unless cuts != nil
//   cuts    -- cut positions. can be nil => use ncuts instead
//   pc      -- probability of crossover
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
func BytCrossover(a, b, A, B [][]byte, time, ncuts int, cuts []int, pc float64, extra interface{}) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			copy(a[i], A[i])
			copy(b[i], B[i])
		}
		return
	}
	ends = GenerateCxEnds(size, ncuts, cuts)
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
//  Input:
//   A and B -- parents
//   ncuts   -- number of cuts to be used, unless cuts != nil
//   cuts    -- cut positions. can be nil => use ncuts instead
//   pc      -- probability of crossover
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
func FunCrossover(a, b, A, B []Func_t, time, ncuts int, cuts []int, pc float64, extra interface{}) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ncuts, cuts)
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

// mutation ////////////////////////////////////////////////////////////////////////////////////////

// IntMutation performs the mutation of genetic data from A
//  Input:
//   A        -- individual
//   nchanges -- number of changes of genes
//   pm       -- probability of mutation
//   extra    -- an integer corresponding to the max value for multiplier 'm'
//  Output: modified individual 'A'
func IntMutation(A []int, time, nchanges int, pm float64, extra interface{}) {
	size := len(A)
	if !rnd.FlipCoin(pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, nchanges)
	mmax := 10
	if val, ok := extra.(int); ok {
		mmax = val
	}
	for _, i := range pos {
		m := rnd.Int(1, mmax)
		if rnd.FlipCoin(0.5) {
			A[i] += m * A[i]
		} else {
			A[i] -= m * A[i]
		}
	}
}

// IntOrdMutation performs the mutation of genetic data from a ordered list of integers A
//  Input:
//   A     -- individual
//   dum1  -- not used
//   pm    -- probability of mutation
//   sti   -- if []int{start, end, insertPoint} != nil, use it; otherwise, use random
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
func IntOrdMutation(A []int, time, dum1 int, pm float64, sti interface{}) {
	size := len(A)
	if !rnd.FlipCoin(pm) || size < 3 {
		if size == 2 {
			A[0], A[1] = A[1], A[0]
		}
		return
	}
	var s, t, ncore, nrem, ins int
	if sti != nil {
		res := sti.([]int)
		s, t, ins = res[0], res[1], res[2]
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
//  Input:
//   A        -- individual
//   nchanges -- number of changes of genes
//   pm       -- probability of mutation
//   extra    -- an integer corresponding to the max value for multiplier 'm'
//  Output: modified individual 'A'
func FltMutation(A []float64, time, nchanges int, pm float64, extra interface{}) {
	size := len(A)
	if !rnd.FlipCoin(pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, nchanges)
	mmax := 10.0
	if val, ok := extra.(float64); ok {
		mmax = val
	}
	for _, i := range pos {
		m := rnd.Float64(1, mmax)
		if rnd.FlipCoin(0.5) {
			A[i] += m * A[i]
		} else {
			A[i] -= m * A[i]
		}
	}
}

// FltMutationNonUni implements the non-uniform mutation (Michaelewicz, 1992; Herrera, 1998)
//func FltMutationNonUni(A []float64, nchanges int, pm float64, extra interface{}) {
//}

// StrMutation performs the mutation of genetic data from A
//  Input:
//   A        -- individual
//   nchanges -- number of changes of genes
//   pm       -- probability of mutation
//   extra    -- an integer corresponding to the max value for multiplier 'm'
//  Output: modified individual 'A'
func StrMutation(A []string, time, nchanges int, pm float64, extra interface{}) {
	size := len(A)
	if !rnd.FlipCoin(pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, nchanges)
	for _, i := range pos {
		A[i] = "TODO" // TODO: improve this
	}
}

// KeyMutation performs the mutation of genetic data from A
//  Input:
//   A        -- individual
//   nchanges -- number of changes of genes
//   pm       -- probability of mutation
//   extra    -- an integer corresponding to the max value for multiplier 'm'
//  Output: modified individual 'A'
func KeyMutation(A []byte, time, nchanges int, pm float64, extra interface{}) {
	size := len(A)
	if !rnd.FlipCoin(pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, nchanges)
	for _, i := range pos {
		v := rnd.Int(0, 100)
		A[i] = byte(v) // TODO: improve this
	}
}

// BytMutation performs the mutation of genetic data from A
//  Input:
//   A        -- individual
//   nchanges -- number of changes of genes
//   pm       -- probability of mutation
//   extra    -- an integer corresponding to the max value for multiplier 'm'
//  Output: modified individual 'A'
func BytMutation(A [][]byte, time, nchanges int, pm float64, extra interface{}) {
	size := len(A)
	if !rnd.FlipCoin(pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, nchanges)
	for _, i := range pos {
		v := rnd.Int(0, 100)
		A[i][0] = byte(v) // TODO: improve this
	}
}

// FunMutation performs the mutation of genetic data from A
//  Input:
//   A        -- individual
//   nchanges -- number of changes of genes
//   pm       -- probability of mutation
//   extra    -- an integer corresponding to the max value for multiplier 'm'
//  Output: modified individual 'A'
func FunMutation(A []Func_t, time, nchanges int, pm float64, extra interface{}) {
	size := len(A)
	if !rnd.FlipCoin(pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, nchanges)
	for _, i := range pos {
		// TODO: improve this
		A[i] = func(ind *Individual) string { return "mutated" }
	}
}

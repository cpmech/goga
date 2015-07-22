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

// Fitness maps objective values into [0, 1]; thus returning the fitness function values
//  Input:
//    ovs -- objective values
//  Output:
//    f -- fitness function values
func Fitness(f, ovs []float64) {
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

// Ranking computes fitness corresponding to a linear ranking
//  Input:
//    ninds -- number of individuals
//    sp    -- selective pressure; must be inside [1, 2]
//  Output:
//    f -- ranked fitnesses
func Ranking(ninds int, sp float64) (f []float64) {
	if sp < 1.0 || sp > 2.0 {
		sp = 1.2
	}
	f = make([]float64, ninds)
	for i := 0; i < ninds; i++ {
		f[i] = 2.0 - sp + 2.0*(sp-1.0)*float64(ninds-i-1)/float64(ninds-1)
	}
	return
}

// CumSum returns the cumulative sum of the elements in p
//  Input:
//   p -- values
//  Output:
//   cs -- cumulated sum
func CumSum(cs, p []float64) {
	chk.IntAssert(len(cs), len(p))
	if len(p) < 1 {
		return
	}
	cs[0] = p[0]
	for i := 1; i < len(p); i++ {
		cs[i] = cs[i-1] + p[i]
	}
}

// RouletteSelect selects n individuals
//  Input:
//    cumprob -- cumulated probabilities (from sorted population)
//    sample  -- a list of random numbers; can be nil
//  Output:
//    selinds -- selected individuals (indices). len(selinds) == nsel
func RouletteSelect(selinds []int, cumprob []float64, sample []float64) {
	nsel := len(selinds)
	chk.IntAssertLessThan(nsel, len(cumprob))
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
	chk.IntAssertLessThan(nsel, len(cumprob))
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

//     0 1 2 3 4 5 6 7
// A = a b c d e f g h    size = 8
// B = * . . . . * * *    ends = [1, 5, 8]
//      ↑       ↑     ↑
//      1       5     8
// a = a . . . . f g h
// b = * b c d e * * *
//
func IntCrossover(a, b, A, B []int, cuts []int, pc float64) {
	size := len(A)
	if !rnd.FlipCoin(pc) || size < 2 {
		copy(a, A)
		copy(b, B)
		return
	}
	ends := GenerateCxEnds(size, cuts)
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
}

func FltCrossover(a, b, A, B []float64, cuts []int, pc float64) {
}

func StrCrossover(a, b, A, B []string, cuts []int, pc float64) {
}

func KeyCrossover(a, b, A, B []byte, cuts []int, pc float64) {
}

func BytCrossover(a, b, A, B [][]byte, cuts []int, pc float64) {
}

func FunCrossover(a, b, A, B []Func_tt, cuts []int, pc float64) {
}

// GenerateCxEnds randomly computes the end positions of cuts in chromosomes
//  Input:
//   size -- size of chromosome
//   cuts -- cut positions. use -1 or greater than size-1 for random value
//  Output:
//   ends -- end positions where the last one equals size
//  Example:
//        0 1 2 3 4 5 6 7
//    A = a b c d e f g h    size = 8
//         ↑       ↑     ↑   cuts = [1, 5]
//         1       5     8   ends = [1, 5, 8]
func GenerateCxEnds(size int, cuts []int) (ends []int) {
	if size < 2 {
		return
	}
	if size == 2 {
		return []int{1, size}
	}
	if len(cuts) == 0 {
		return []int{rnd.Int(1, size-1), size}
	}
	ncuts := len(cuts)
	if ncuts >= size-1 {
		ncuts = size - 1
	}
	ends = make([]int, ncuts+1)
	ends[ncuts] = size
	pool := utl.IntRange2(1, size)
	rnd.IntShuffle(pool)
	for i := 0; i < ncuts; i++ {
		if cuts[i] < 0 || cuts[i] > size-1 {
			ends[i] = pool[i]
		} else {
			if i > 0 {
				if cuts[i] == cuts[i-1] { // avoid repeated cuts
					ends[i] = pool[i]
					continue
				}
			}
			ends[i] = cuts[i]
		}
	}
	sort.Ints(ends)
	return ends
}

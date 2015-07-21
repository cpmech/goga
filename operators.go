// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"math/rand"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
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

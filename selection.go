// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math/rand"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
)

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

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

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

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/rnd"

// GenTrialSolutions generates (initial) trial solutions
func GenTrialSolutions(sols []*Solution, prms *Parameters) {

	// floats
	n := len(sols)
	if prms.Nflt > 0 {
		K := rnd.LatinIHS(prms.Nflt, n, prms.LatinDup)
		for i := 0; i < n; i++ {
			for j := 0; j < prms.Nflt; j++ {
				sols[i].Flt[j] = prms.FltMin[j] + float64(K[j][i]-1)*prms.DelFlt[j]/float64(n-1)
			}
		}
	}

	// skip if there are no ints
	if prms.Nint < 2 {
		return
	}

	// binary numbers
	if prms.BinInt > 0 {
		for i := 0; i < n; i++ {
			for j := 0; j < prms.Nint; j++ {
				if rnd.FlipCoin(0.5) {
					sols[i].Int[j] = 1
				} else {
					sols[i].Int[j] = 0
				}
			}
		}
		return
	}

	// general integers
	L := rnd.LatinIHS(prms.Nint, n, prms.LatinDup)
	for i := 0; i < n; i++ {
		for j := 0; j < prms.Nint; j++ {
			sols[i].Int[j] = prms.IntMin[j] + (L[j][i]-1)*prms.DelInt[j]/(n-1)
		}
	}
}

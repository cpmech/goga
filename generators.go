// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
)

// GenTrialSolutions generates (initial) trial solutions
func GenTrialSolutions(sols []*Solution, prms *Parameters, reset bool) {

	// reset solutions
	if reset {
		for id, sol := range sols {
			sol.Reset(id)
		}
	}

	// floats
	n := len(sols) // cannot use Nsol here because subsets of Solutions may be provided; e.g. parallel code
	if prms.Nflt > 0 {

		// interior points
		switch prms.GenType {
		case "latin":
			K := rnd.LatinIHS(prms.Nflt, n, prms.LatinDup)
			for i := 0; i < n; i++ {
				for j := 0; j < prms.Nflt; j++ {
					sols[i].Flt[j] = prms.FltMin[j] + float64(K[j][i]-1)*prms.DelFlt[j]/float64(n-1)
				}
			}
		case "halton":
			H := rnd.HaltonPoints(prms.Nflt, n)
			for i := 0; i < n; i++ {
				for j := 0; j < prms.Nflt; j++ {
					sols[i].Flt[j] = prms.FltMin[j] + H[j][i]*prms.DelFlt[j]
				}
			}
		default:
			for i := 0; i < n; i++ {
				for j := 0; j < prms.Nflt; j++ {
					sols[i].Flt[j] = rnd.Float64(prms.FltMin[j], prms.FltMax[j])
				}
			}
		}

		// extra points
		if prms.UseMesh {
			initX := func(isol int) {
				for k := 0; k < prms.Nflt; k++ {
					sols[isol].Flt[k] = (prms.FltMin[k] + prms.FltMax[k]) / 2.0
				}
			}
			isol := prms.Nsol - prms.NumExtraSols
			for i := 0; i < prms.Nflt-1; i++ {
				for j := i + 1; j < prms.Nflt; j++ {
					// (min,min) corner
					initX(isol)
					sols[isol].Flt[i] = prms.FltMin[i]
					sols[isol].Flt[j] = prms.FltMin[j]
					sols[isol].Fixed = true
					isol++
					// (min,max) corner
					initX(isol)
					sols[isol].Flt[i] = prms.FltMin[i]
					sols[isol].Flt[j] = prms.FltMax[j]
					sols[isol].Fixed = true
					isol++
					// (max,max) corner
					initX(isol)
					sols[isol].Flt[i] = prms.FltMax[i]
					sols[isol].Flt[j] = prms.FltMax[j]
					sols[isol].Fixed = true
					isol++
					// (max,min) corner
					initX(isol)
					sols[isol].Flt[i] = prms.FltMax[i]
					sols[isol].Flt[j] = prms.FltMin[j]
					sols[isol].Fixed = true
					isol++
					// Xi-min middle points
					ndelta := float64(prms.Nbry - 1)
					for m := 0; m < prms.Nbry-2; m++ {
						initX(isol)
						sols[isol].Flt[i] = prms.FltMin[i]
						sols[isol].Flt[j] = prms.FltMin[j] + float64(m+1)*prms.DelFlt[j]/ndelta
						sols[isol].Fixed = true
						isol++
					}
					// Xi-max middle points
					for m := 0; m < prms.Nbry-2; m++ {
						initX(isol)
						sols[isol].Flt[i] = prms.FltMax[i]
						sols[isol].Flt[j] = prms.FltMin[j] + float64(m+1)*prms.DelFlt[j]/ndelta
						sols[isol].Fixed = true
						isol++
					}
					// Xj-min middle points
					for m := 0; m < prms.Nbry-2; m++ {
						initX(isol)
						sols[isol].Flt[i] = prms.FltMin[i] + float64(m+1)*prms.DelFlt[i]/ndelta
						sols[isol].Flt[j] = prms.FltMin[j]
						sols[isol].Fixed = true
						isol++
					}
					// Xj-max middle points
					for m := 0; m < prms.Nbry-2; m++ {
						initX(isol)
						sols[isol].Flt[i] = prms.FltMin[i] + float64(m+1)*prms.DelFlt[i]/ndelta
						sols[isol].Flt[j] = prms.FltMax[j]
						sols[isol].Fixed = true
						isol++
					}
				}
			}
			chk.IntAssert(isol, prms.Nsol)
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

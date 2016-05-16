// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
)

// GenTrialSolutions generates (initial) trial solutions
func GenTrialSolutions(sols []*Solution, prms *Parameters) {

	// floats
	n := len(sols) // cannot use Nsol here because subsets of Solutions may be provided; e.g. parallel code
	if prms.Nx > 0 {

		// interior points
		switch prms.GenType {
		case "latin":
			K := rnd.LatinIHS(prms.Nx, n, prms.LatinDup)
			for i := 0; i < n; i++ {
				for j := 0; j < prms.Nx; j++ {
					sols[i].Flt[j] = float64(K[j][i]-1) / float64(n-1)
				}
			}
		case "halton":
			H := rnd.HaltonPoints(prms.Nx, n)
			for i := 0; i < n; i++ {
				for j := 0; j < prms.Nx; j++ {
					sols[i].Flt[j] = H[j][i]
				}
			}
		default:
			for i := 0; i < n; i++ {
				for j := 0; j < prms.Nx; j++ {
					sols[i].Flt[j] = rnd.Float64(0, 1)
				}
			}
		}

		// extra points
		if prms.UseMesh {
			initX := func(isol int) {
				for k := 0; k < prms.Nx; k++ {
					sols[isol].Flt[k] = 0.5
				}
			}
			isol := prms.Nsol - prms.NumExtraSols
			for i := 0; i < prms.Nx-1; i++ {
				for j := i + 1; j < prms.Nx; j++ {
					// (min,min) corner
					initX(isol)
					sols[isol].Flt[i] = 0.0
					sols[isol].Flt[j] = 0.0
					sols[isol].Fixed = true
					isol++
					// (min,max) corner
					initX(isol)
					sols[isol].Flt[i] = 0.0
					sols[isol].Flt[j] = 1.0
					sols[isol].Fixed = true
					isol++
					// (max,max) corner
					initX(isol)
					sols[isol].Flt[i] = 1.0
					sols[isol].Flt[j] = 1.0
					sols[isol].Fixed = true
					isol++
					// (max,min) corner
					initX(isol)
					sols[isol].Flt[i] = 1.0
					sols[isol].Flt[j] = 0.0
					sols[isol].Fixed = true
					isol++
					// Xi-min middle points
					ndelta := float64(prms.Nbry - 1)
					for m := 0; m < prms.Nbry-2; m++ {
						initX(isol)
						sols[isol].Flt[i] = 0.0
						sols[isol].Flt[j] = float64(m+1) / ndelta
						sols[isol].Fixed = true
						isol++
					}
					// Xi-max middle points
					for m := 0; m < prms.Nbry-2; m++ {
						initX(isol)
						sols[isol].Flt[i] = 1.0
						sols[isol].Flt[j] = float64(m+1) / ndelta
						sols[isol].Fixed = true
						isol++
					}
					// Xj-min middle points
					for m := 0; m < prms.Nbry-2; m++ {
						initX(isol)
						sols[isol].Flt[i] = float64(m+1) / ndelta
						sols[isol].Flt[j] = 0.0
						sols[isol].Fixed = true
						isol++
					}
					// Xj-max middle points
					for m := 0; m < prms.Nbry-2; m++ {
						initX(isol)
						sols[isol].Flt[i] = float64(m+1) / ndelta
						sols[isol].Flt[j] = 1.0
						sols[isol].Fixed = true
						isol++
					}
				}
			}
			chk.IntAssert(isol, prms.Nsol)
		}
	}

	// skip if there are no ints
	if prms.Nk < 2 {
		return
	}

	// binary numbers
	if prms.BinInt > 0 {
		for i := 0; i < n; i++ {
			for j := 0; j < prms.Nk; j++ {
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
	L := rnd.LatinIHS(prms.Nk, n, prms.LatinDup)
	for i := 0; i < n; i++ {
		for j := 0; j < prms.Nk; j++ {
			sols[i].Int[j] = prms.Kmin[j] + (L[j][i]-1)*prms.Dk[j]/(n-1)
		}
	}
}

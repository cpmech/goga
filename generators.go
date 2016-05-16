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
	nsol := len(sols) - prms.NumExtraSols // may be smaller than Nsol when using multiple CPUs
	if prms.Nx > 0 {

		// interior points
		switch prms.GenType {
		case "latin":
			K := rnd.LatinIHS(prms.Nx, nsol, prms.LatinDup)
			for i := 0; i < nsol; i++ {
				for j := 0; j < prms.Nx; j++ {
					sols[i].Flt[j] = prms.GapX + (1.0-2.0*prms.GapX)*float64(K[j][i]-1)/float64(nsol-1)
				}
			}
		case "halton":
			H := rnd.HaltonPoints(prms.Nx, nsol)
			for i := 0; i < nsol; i++ {
				for j := 0; j < prms.Nx; j++ {
					sols[i].Flt[j] = prms.GapX + (1.0-2.0*prms.GapX)*H[j][i]
				}
			}
		default:
			for i := 0; i < nsol; i++ {
				for j := 0; j < prms.Nx; j++ {
					sols[i].Flt[j] = rnd.Float64(prms.GapX, 1.0-prms.GapX)
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
			isol := nsol
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
			chk.IntAssert(isol, len(sols))
		}
	}

	// skip if there are no ints
	if prms.Nk < 2 {
		return
	}

	// binary numbers
	if prms.BinInt > 0 {
		for i := 0; i < nsol; i++ {
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
	L := rnd.LatinIHS(prms.Nk, nsol, prms.LatinDup)
	for i := 0; i < nsol; i++ {
		for j := 0; j < prms.Nk; j++ {
			sols[i].Int[j] = prms.Kmin[j] + (L[j][i]-1)*prms.Dk[j]/(nsol-1)
		}
	}
}

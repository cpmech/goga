// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"math/rand"

	"github.com/cpmech/gosl/rnd"
)

func de_operator(u, x, x0, x1, x2 []float64, prms *Parameters) {
	C := 1.0
	if prms.DEuseC {
		C = 1.0 + rand.NormFloat64()*0.25
	}
	F := prms.DEmult * C
	K := 0.5 * (F + 1.0)
	n := len(x)
	I := rnd.Int(0, n-1)
	for i := 0; i < n; i++ {
		if rnd.FlipCoin(prms.DEpc) || i == I {
			if rnd.FlipCoin(prms.DEpm) {
				u[i] = x0[i] + F*(x1[i]-x2[i])
			} else {
				u[i] = x0[i] + K*(x1[i]+x2[i]-2.0*x0[i])
			}
		} else {
			u[i] = x[i]
		}
		u[i] = prms.EnforceRange(i, u[i])
	}
}

// CxFltDE implements the differential-evolution crossover
func CxFltDE(a, b, A, B, x0, x1, x2, y0, y1, y2 []float64, prms *Parameters) {
	de_operator(a, A, x0, x1, x2, prms)
	de_operator(b, B, y0, y1, y2, prms)
}

// MtFltDeb implements Deb's parameter-based mutation operator
//  [1] Deb K and Tiwari S (2008) Omni-optimizer: A generic evolutionary algorithm for single
//      and multi-objective optimization. European Journal of Operational Research, 185:1062-1087.
func MtFltDeb(A []float64, prms *Parameters) {

	// skip mutation
	if !rnd.FlipCoin(prms.PmFlt) {
		return
	}

	// for each gene
	size := len(A)
	pm := 1.0 / float64(size)
	ηm := prms.DebEtam
	cm := 1.0 / (ηm + 1.0)
	var u, Δx, φ1, φ2, δ1, δ2, δb, xl, xu float64
	for i := 0; i < size; i++ {

		// leave basis unmodified
		if !rnd.FlipCoin(pm) {
			continue
		}

		// range
		xl, xu = prms.FltMin[i], prms.FltMax[i]
		Δx = xu - xl

		// mutation
		u = rnd.Float64(0, 1)
		δ1 = (A[i] - xl) / Δx
		δ2 = (xu - A[i]) / Δx
		if u <= 0.5 {
			φ1 = math.Pow(1.0-δ1, ηm+1.0)
			δb = math.Pow(2.0*u+(1.0-2.0*u)*φ1, cm) - 1.0
		} else {
			φ2 = math.Pow(1.0-δ2, ηm+1.0)
			δb = 1.0 - math.Pow(2.0-2.0*u+(2.0*u-1.0)*φ2, cm)
		}
		A[i] = prms.EnforceRange(i, A[i]+δb*Δx)
	}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/rnd"
)

// CxFltDE implements the differential-evolution crossover
func CxFltDE(a, b, A, B, C, D, E, F []float64, prms *Parameters) {
	n := len(A)
	ia := rnd.Int(0, n-1)
	ib := rnd.Int(0, n-1)
	var x float64
	for i := 0; i < n; i++ {

		// a
		if rnd.FlipCoin(prms.DEpc) || i == ia {
			x = B[i] + prms.DEmult*(C[i]-D[i]) //+ prms.DEmult*(D[i]-E[i]) + prms.DEmult*(E[i]-F[i])
		} else {
			x = A[i]
		}
		a[i] = prms.EnforceRange(i, x)

		// b
		if rnd.FlipCoin(prms.DEpc) || i == ib {
			x = A[i] + prms.DEmult*(D[i]-C[i]) //+ prms.DEmult*(E[i]-D[i]) + prms.DEmult*(F[i]-E[i])
		} else {
			x = B[i]
		}
		b[i] = prms.EnforceRange(i, x)
	}
	return
}

// CxFltDeb implements Deb's simulated binary crossover (SBX)
func CxFltDeb(a, b, A, B, C, D, E, F []float64, prms *Parameters) {

	// for each gene
	ϵ := 1e-10
	cc := 1.0 / (prms.DebEtac + 1.0)
	size := len(A)
	var u, α, β, βb, x1, x2, δx, xl, xu float64
	for i := 0; i < size; i++ {

		// parents' basis values
		x1, x2 = A[i], B[i]
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		δx = x2 - x1

		// copy only
		if rnd.FlipCoin(0.5) || δx < ϵ {
			a[i], b[i] = A[i], B[i]
			continue
		}

		// random number
		u = rnd.Float64(0, 1)

		// range
		xl, xu = prms.FltMin[i], prms.FltMax[i]

		// first offspring
		β = 1.0 + 2.0*(x1-xl)/δx
		α = 2.0 - math.Pow(β, -(prms.DebEtac+1.0))
		if u <= 1.0/α {
			βb = math.Pow(α*u, cc)
		} else {
			βb = math.Pow(1.0/(2.0-α*u), cc)
		}
		a[i] = prms.EnforceRange(i, 0.5*(x1+x2-βb*δx))

		// second offspring
		β = 1.0 + 2.0*(xu-x2)/δx
		α = 2.0 - math.Pow(β, -(prms.DebEtac+1.0))
		if u <= (1.0 / α) {
			βb = math.Pow(α*u, cc)
		} else {
			βb = math.Pow(1.0/(2.0-α*u), cc)
		}
		b[i] = prms.EnforceRange(i, 0.5*(x1+x2+βb*δx))
	}
	return
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

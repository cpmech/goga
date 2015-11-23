// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
)

// Deb and Tiwari /////////////////////////////////////////////////////////////////////////////////

// FltCrossoverDB implements Deb's simulated binary crossover (SBX)
func FltCrossoverDB(a, b, A, B, unusedC, unusedD []float64, time int, ops *OpsData) (ends []int) {

	// check
	chk.IntAssert(len(ops.Xrange), len(A))

	// for each gene
	ϵ := 1e-10
	cc := 1.0 / (ops.DebEtac + 1.0)
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
		xl, xu = ops.Xrange[i][0], ops.Xrange[i][1]

		// first offspring
		β = 1.0 + 2.0*(x1-xl)/δx
		α = 2.0 - math.Pow(β, -(ops.DebEtac+1.0))
		if u <= 1.0/α {
			βb = math.Pow(α*u, cc)
		} else {
			βb = math.Pow(1.0/(2.0-α*u), cc)
		}
		a[i] = ops.EnforceRange(i, 0.5*(x1+x2-βb*δx))

		// second offspring
		β = 1.0 + 2.0*(xu-x2)/δx
		α = 2.0 - math.Pow(β, -(ops.DebEtac+1.0))
		if u <= (1.0 / α) {
			βb = math.Pow(α*u, cc)
		} else {
			βb = math.Pow(1.0/(2.0-α*u), cc)
		}
		b[i] = ops.EnforceRange(i, 0.5*(x1+x2+βb*δx))
	}
	return
}

//  FltMutationDB implements Deb's parameter-based mutation operator
//  References:
//   [1] Deb K and Tiwari S (2008) Omni-optimizer: A generic evolutionary algorithm for single
//       and multi-objective optimization. European Journal of Operational Research, 185:1062-1087.
func FltMutationDB(A []float64, time int, ops *OpsData) {

	// check
	size := len(A)
	chk.IntAssert(len(ops.Xrange), size)

	// skip mutation
	if !rnd.FlipCoin(ops.FltPm) {
		return
	}

	// for each gene
	pm := 1.0 / float64(size)
	ηm := ops.DebEtam
	cm := 1.0 / (ηm + 1.0)
	var u, Δx, φ1, φ2, δ1, δ2, δb, xl, xu float64
	for i := 0; i < size; i++ {

		// leave basis unmodified
		if !rnd.FlipCoin(pm) {
			continue
		}

		// range
		xl, xu = ops.Xrange[i][0], ops.Xrange[i][1]
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
		A[i] = ops.EnforceRange(i, A[i]+δb*Δx)
	}
}

// Differential-Evolution /////////////////////////////////////////////////////////////////////////

// FltCrossoverDE implements the differential-evolution crossover
func FltCrossoverDE(a, b, A, B, C, D []float64, time int, ops *OpsData) (ends []int) {
	n := len(A)
	sa := rnd.Int(0, n-1)
	sb := rnd.Int(0, n-1)
	var x float64
	for s := 0; s < n; s++ {

		// a
		if rnd.FlipCoin(ops.DEpc) || s == sa {
			x = B[s] + ops.DEmult*(C[s]-D[s])
		} else {
			x = A[s]
		}
		a[s] = ops.EnforceRange(s, x)

		// b
		if rnd.FlipCoin(ops.DEpc) || s == sb {
			x = A[s] + ops.DEmult*(D[s]-C[s])
		} else {
			x = B[s]
		}
		b[s] = ops.EnforceRange(s, x)
	}
	return
}

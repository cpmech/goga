// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// general ////////////////////////////////////////////////////////////////////////////////////////

// FltCrossover performs the crossover of genetic data from A and B
//  Output:
//   a and b -- offspring
//  Example:
//         0 1 2 3 4 5 6 7
//     A = a b c d e f g h    size = 8
//     B = * . . . . * * *    cuts = [1, 5]
//          ↑       ↑     ↑   ends = [1, 5, 8]
//          1       5     8
//     a = a . . . . f g h
//     b = * b c d e * * *
func FltCrossover(a, b, A, B []float64, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ops.Ncuts, ops.Cuts)
	swap := false
	start := 0
	for _, end := range ends {
		if swap {
			for j := start; j < end; j++ {
				b[j], a[j] = A[j], B[j]
			}
		} else {
			for j := start; j < end; j++ {
				a[j], b[j] = A[j], B[j]
			}
		}
		start = end
		swap = !swap
	}
	return
}

// FltMutation performs the mutation of genetic data from A
//  Output: modified individual 'A'
func FltMutation(A []float64, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		m := rnd.Float64(1, ops.Mmax)
		if rnd.FlipCoin(0.5) {
			A[i] += m * A[i]
		} else {
			A[i] -= m * A[i]
		}
	}
}

// BLX-α and Michaelicz ///////////////////////////////////////////////////////////////////////////

// FltCrossoverBlx implements the BLS-α crossover by Eshelman et al. (1993); see also Herrera (1998)
//  Output:
//   a and b -- offspring
func FltCrossoverBlx(a, b, A, B []float64, time int, ops *OpsData) (ends []int) {
	chk.IntAssert(len(ops.Xrange), len(A))
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) {
		for i := 0; i < size; i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	α := ops.BlxAlp
	var cmin, cmax, δ float64
	for i := 0; i < size; i++ {
		cmin = utl.Min(A[i], B[i])
		cmax = utl.Max(A[i], B[i])
		δ = cmax - cmin
		a[i] = rnd.Float64(cmin-α*δ, cmax+α*δ)
		b[i] = rnd.Float64(cmin-α*δ, cmax+α*δ)
		a[i] = ops.EnforceRange(i, a[i])
		b[i] = ops.EnforceRange(i, b[i])
	}
	return
}

// FltMutationMwicz implements the non-uniform mutation (Michaelewicz, 1992; Herrera, 1998)
// See also Michalewicz (1996) page 103
func FltMutationMwicz(A []float64, time int, ops *OpsData) {
	chk.IntAssert(len(ops.Xrange), len(A))
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	t := float64(time)
	for i := 0; i < size; i++ {
		xmin := ops.Xrange[i][0]
		xmax := ops.Xrange[i][1]
		if rnd.FlipCoin(0.5) {
			A[i] += ops.MwiczDelta(t, xmax-A[i])
		} else {
			A[i] -= ops.MwiczDelta(t, A[i]-xmin)
		}
		A[i] = ops.EnforceRange(i, A[i])
	}
}

// Deb and Tiwari /////////////////////////////////////////////////////////////////////////////////

// FltCrossoverDeb implements Deb's simulated binary crossover (SBX)
func FltCrossoverDeb(a, b, A, B []float64, time int, ops *OpsData) (ends []int) {

	// check
	chk.IntAssert(len(ops.Xrange), len(A))

	// copy only
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) {
		for i := 0; i < size; i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}

	// for each gene
	ϵ := 1e-10
	cc := 1.0 / (ops.DebEtac + 1.0)
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

		// crossover
		u = rnd.Float64(0, 1)
		if ops.EnfRange {

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

		} else {
			if u <= 0.5 {
				βb = math.Pow(2.0*u, cc)
			} else {
				βb = math.Pow(0.5/(1.0-u), cc)
			}
			a[i] = 0.5 * (x1 + x2 - βb*δx)
			b[i] = 0.5 * (x1 + x2 + βb*δx)
		}
	}
	return
}

//  FltMutationDeb implements Deb's parameter-based mutation operator
//  References:
//   [1] Deb K and Tiwari S (2008) Omni-optimizer: A generic evolutionary algorithm for single
//       and multi-objective optimization. European Journal of Operational Research, 185:1062-1087.
func FltMutationDeb(A []float64, time int, ops *OpsData) {

	// check
	size := len(A)
	chk.IntAssert(len(ops.Xrange), size)

	// no mutation
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
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
		if ops.EnfRange {
			δ1 = (A[i] - xl) / Δx
			δ2 = (xu - A[i]) / Δx
			if u <= 0.5 {
				φ1 = math.Pow(1.0-δ1, ηm+1.0)
				δb = math.Pow(2.0*u+(1.0-2.0*u)*φ1, cm) - 1.0
			} else {
				φ2 = math.Pow(1.0-δ2, ηm+1.0)
				δb = 1.0 - math.Pow(2.0-2.0*u+(2.0*u-1.0)*φ2, cm)
			}
		} else {
			if u <= 0.5 {
				δb = math.Pow(2.0*u, cm) - 1.0
			} else {
				δb = 1.0 - math.Pow(2.0-2.0*u, cm)
			}
		}
		A[i] = ops.EnforceRange(i, A[i]+δb*Δx)
	}
}

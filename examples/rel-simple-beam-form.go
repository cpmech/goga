// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

// deflection computes the max deflection at mid-span of simply supported beam
// with uniform distributed load q and concentrated load at midspan
//  Input:
//   L  -- span
//   EI -- Young's modulus times cross-sectional moment of inertia
//   p  -- concentrated load at mid-span
//   q  -- distributed load
//  Output:
//   max deflection (vertical displacement) at mid-span
func deflection(L, EI, p, q float64) float64 {
	L3 := math.Pow(L, 3.0)
	return p*L3/EI/48.0 + 5.0*q*L3*L/EI/384.0
}
func deflection_derivs(L, EI, p, q float64) (dp, dq float64) {
	L3 := math.Pow(L, 3.0)
	dp = L3 / EI / 48.0
	dq = 5.0 * L3 * L / EI / 384.0
	return
}

func main() {

	// Example from
	// Achintya Haldar, Reliability-Based Structura Design, 2005

	δlim := 0.0381 // [m] max allowed deflection
	L := 9.144     // [m]
	EI := 182262.0 // [kN m²]
	μp := 111.2    // [kN] mean of p (lognormal random variable)
	σp := 11.12    // [kN] deviation of p
	μq := 35.03    // [kN/m] mean value of q (normal random variable)
	σq := 5.25     // [kN/m] standard deviation of q

	_ = L
	_ = EI
	_ = σq

	// limit state function and derivative
	gfcn := func(p, q float64) float64 {
		return δlim - deflection(L, EI, p, q)
	}
	dgfcn := func(p, q float64) (dgdp, dgdq float64) {
		dp, dq := deflection_derivs(L, EI, p, q)
		dgdp, dgdq = -dp, -dq
		return
	}

	_ = gfcn

	var ln rnd.LogNormal
	ln.Sig = σp / μp
	ln.Mu = math.Log(μp) - ln.Sig*ln.Sig/2.0 // TODO: check this equation
	ln.CalcDerived()

	β := 3.0
	ps := μp // p-star: trial intersection point
	qs := μq // q-star: trial intersection point

	ps = 123.841

	fp := ln.Pdf(ps)
	Φinvp := (math.Log(ps) - ln.Mu) / ln.Sig
	Φinvp = 1.13
	φp := math.Exp(-Φinvp*Φinvp/2.0) / math.Sqrt2 / math.SqrtPi
	σpN := φp / fp
	μpN := ps - Φinvp*σpN

	dp, dq := dgfcn(ps, qs)
	den := math.Sqrt(math.Pow(dp*σpN, 2.0) + math.Pow(dq*σq, 2.0))
	αp := dp * σpN / den
	αq := dq * σq / den

	io.Pfpink("αp = %v\n", αp)
	io.Pfpink("αq = %v\n", αq)

	_ = β
	_ = μpN
}

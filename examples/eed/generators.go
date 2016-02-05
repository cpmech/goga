// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// Generator holds data for one generator
type Generator struct {
	a, b, c       float64 // cost coefficients
	α, β, γ, ζ, λ float64 // emission coefficients
	Pmin, Pmax    float64 // capacity constraints
}

// System holds generators and other data
type System struct {
	Pdemand  float64
	Lossless bool
	B00      float64
	B0       []float64
	B        [][]float64
	G        []Generator
}

// Init initialises system
func (o *System) Init(Pdemand float64, lossless, check bool) {

	o.Pdemand = Pdemand
	o.Lossless = lossless

	// units:
	//  a [$ / (hMW²)]
	//  b [$ / (hMW)]
	//  c [$ / h]
	//  α [tons / (hMW²)]
	//  β [tons / (hMW)]
	//  γ [tons / h]
	//  ζ [tons / h]
	//  λ [MW⁻¹]
	//  Pmin [MW / 100]
	//  Pmax [MW / 100]

	o.G = []Generator{
		{a: 10, b: 200, c: 100, α: 4.091e-2, β: -5.554e-2, γ: 6.490e-2, ζ: 2.0e-4, λ: 2.857, Pmin: 0.05, Pmax: 0.5},
		{a: 10, b: 150, c: 120, α: 2.543e-2, β: -6.047e-2, γ: 5.638e-2, ζ: 5.0e-4, λ: 3.333, Pmin: 0.05, Pmax: 0.6},
		{a: 20, b: 180, c: 40., α: 4.258e-2, β: -5.094e-2, γ: 4.586e-2, ζ: 1.0e-6, λ: 8.000, Pmin: 0.05, Pmax: 1.0},
		{a: 10, b: 100, c: 60., α: 5.326e-2, β: -3.550e-2, γ: 3.380e-2, ζ: 2.0e-3, λ: 2.000, Pmin: 0.05, Pmax: 1.2},
		{a: 20, b: 180, c: 40., α: 4.258e-2, β: -5.094e-2, γ: 4.586e-2, ζ: 1.0e-6, λ: 8.000, Pmin: 0.05, Pmax: 1.0},
		{a: 10, b: 150, c: 100, α: 6.131e-2, β: -5.555e-2, γ: 5.151e-2, ζ: 1.0e-5, λ: 6.667, Pmin: 0.05, Pmax: 0.6},
	}

	o.B00 = 0.00098573
	o.B0 = []float64{-0.0107, +0.0060, -0.0017, +0.0009, +0.0002, +0.0030}
	o.B = [][]float64{
		{+0.1382, -0.0299, +0.0044, -0.0022, -0.0010, -0.0008},
		{-0.0299, +0.0487, -0.0025, +0.0004, +0.0016, +0.0041},
		{+0.0044, -0.0025, +0.0182, -0.0070, -0.0066, -0.0066},
		{-0.0022, +0.0004, -0.0070, +0.0137, +0.0050, +0.0033},
		{-0.0010, +0.0016, -0.0066, +0.0050, +0.0109, +0.0005},
		{-0.0008, +0.0041, -0.0066, +0.0033, +0.0005, +0.0244},
	}

	if check {

		// lossless and unsecured: cost only
		P_best_cost := []float64{0.10954, 0.29967, 0.52447, 1.01601, 0.52469, 0.35963}
		c := o.FuelCost(P_best_cost)
		e := o.Emission(P_best_cost)
		io.Pf("lossless and unsecured: cost only\n")
		io.Pforan("c = %.3f (600.114)\n", c)
		io.Pforan("e = %.5f (0.22214)\n", e)
		P_best_cost = []float64{0.1265, 0.2843, 0.5643, 1.0468, 0.5278, 0.2801}
		c = o.FuelCost(P_best_cost)
		io.Pfgreen("c = %.3f\n", c)
		Pdemand := 2.834
		o.PrintConstraints(P_best_cost, Pdemand, true)

		// lossless and unsecured: emission only
		P_best_emission := []float64{0.40584, 0.45915, 0.53797, 0.38300, 0.53791, 0.51012}
		c = o.FuelCost(P_best_emission)
		e = o.Emission(P_best_emission)
		io.Pf("\nlossless and unsecured: emission only\n")
		io.Pforan("c = %.3f (638.260)\n", c)
		io.Pforan("e = %.5f (0.19420)\n", e)

		P_best_cost = []float64{0.1500, 0.3000, 0.5500, 1.0500, 0.4600, 0.3500}
		c = o.FuelCost(P_best_cost)
		e = o.Emission(P_best_cost)
		io.Pforan("\nc = %.3f (606.314)\n", c)
		io.Pforan("e = %.5f (0.22330)\n", e)

		P_best_emission = []float64{0.4000, 0.4500, 0.5500, 0.4000, 0.5500, 0.5000}
		c = o.FuelCost(P_best_emission)
		e = o.Emission(P_best_emission)
		io.Pforan("\nc = %.3f (639.600)\n", c)
		io.Pforan("e = %.5f (0.19424)\n", e)
	}
	return
}

// FuelCost computes the total $/h fuel cost for given P[i] (power outputs)
func (o System) FuelCost(P []float64) (C float64) {
	for i, g := range o.G {
		C += g.a + g.b*P[i] + g.c*P[i]*P[i]
	}
	return
}

// Emission computes the total ton/h emmision of atmospheric pollutants
func (o System) Emission(P []float64) (E float64) {
	for i, g := range o.G {
		E += g.α + g.β*P[i] + g.γ*P[i]*P[i] + g.ζ*math.Exp(g.λ*P[i])
	}
	return
}

// Ploss computes P_{loss}
func (o System) Ploss(P []float64) (res float64) {
	if o.Lossless {
		return
	}
	res = o.B00
	ngs := len(o.G)
	for i := 0; i < ngs; i++ {
		res += o.B0[i] * P[i]
		for j := 0; j < ngs; j++ {
			res += P[i] * o.B[i][j] * P[j]
		}
	}
	return
}

// Balance computes load balance
func (o System) Balance(P []float64) float64 {
	return math.Abs(la.VecAccum(P) - o.Pdemand - o.Ploss(P))
}

// PrintConstraints prints violated or not constraints
func (o System) PrintConstraints(P []float64, Pdemand float64, full bool) {
	sumP := 0.0
	for i, g := range o.G {
		if full {
			io.Pfyel("P%d range error = %v\n", i, utl.GtePenalty(P[i], g.Pmin, 1)+utl.GtePenalty(g.Pmax, P[i], 1))
		}
		sumP += P[i]
	}
	Ploss := 0.0
	io.Pf("balance error = %v\n", math.Abs(sumP-Pdemand-Ploss))
}

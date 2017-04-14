// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "github.com/cpmech/gosl/io"

func main() {

	G := NewGenerators(false)

	B := [][]float64{
		{+0.1382, -0.0299, +0.0044, -0.0022, -0.0010, -0.0008},
		{-0.0299, +0.0487, -0.0025, +0.0004, +0.0016, +0.0041},
		{+0.0044, -0.0025, +0.0182, -0.0070, -0.0066, -0.0066},
		{-0.0022, +0.0004, -0.0070, +0.0137, +0.0050, +0.0033},
		{-0.0010, +0.0016, -0.0066, +0.0050, +0.0109, +0.0005},
		{-0.0008, +0.0041, -0.0066, +0.0033, +0.0005, +0.0244},
	}
	B0 := []float64{-0.0107, +0.0060, -0.0017, +0.0009, +0.0002, +0.0030}
	B00 := 0.00098573
	Pload := 283.4 // MW

	P := []float64{12.6481 / 100.0, 28.4796 / 100.0, 58.2603 / 100.0, 99.1763 / 100.0, 52.2454 / 100.0, 35.1321 / 100.0}

	// reference
	Cref := 605.9865           // total fuel cost
	Eref := 0.2204             // emmision
	PlossRef := 2.5417 / 100.0 // loss

	_ = Pload

	N := 6
	Ploss := B00
	for i := 0; i < N; i++ {
		Ploss += B0[i] * P[i]
		for j := 0; j < N; j++ {
			Ploss += P[i] * B[i][j] * P[j]
		}
	}

	io.Pforan("Cost     = %v (%g)\n", G.FuelCost(P), Cref)
	io.Pforan("Emission = %v (%g)\n", G.Emission(P), Eref)
	io.Pforan("Ploss    = %v (%g)\n", Ploss, PlossRef)
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/rnd"

// DiffEvol performs the differential-evolution operation
func DiffEvol(xnew, x, x0, x1, x2 []float64, prms *Parameters) {
	C, F := prms.DiffEvolC, prms.DiffEvolF
	if prms.DiffEvolUseCmult {
		C *= rnd.Float64(0, 0.5)
	}
	if prms.DiffEvolUseFmult {
		F *= rnd.Float64(0, 0.5)
	}
	K := 0.5 * (F + 1.0)
	G := rnd.Float64(0, 0.5)
	n := len(x)
	I := rnd.Int(0, n-1)
	mutation := rnd.Int(1, 3)
	for i := 0; i < n; i++ {
		if rnd.FlipCoin(C) || i == I {
			switch mutation {
			case 1:
				xnew[i] = x0[i] + F*x1[i] - G*x2[i]
			case 2:
				xnew[i] = x0[i] + K*(x1[i]+x2[i]-2.0*x0[i])
			default:
				xnew[i] = x0[i] + F*(x1[i]-x2[i])
			}
		} else {
			xnew[i] = x[i]
		}
		xnew[i] = prms.EnforceRange(i, xnew[i])
	}
}

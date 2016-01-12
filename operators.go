// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/rnd"

// DiffEvol performs the differential-evolution operation
func DiffEvol(xnew, x, x0, x1, x2 []float64, prms *Parameters) {
	C := rnd.Float64(0.0, 1.0)
	F := rnd.Float64(0.0, 1.0)
	K := 0.5 * (F + 1.0)
	n := len(x)
	I := rnd.Int(0, n-1)
	mutation := rnd.FlipCoin(0.5)
	for i := 0; i < n; i++ {
		if rnd.FlipCoin(C) || i == I {
			if mutation {
				xnew[i] = x0[i] + K*(x1[i]+x2[i]-2.0*x0[i])
			} else {
				xnew[i] = x0[i] + F*(x1[i]-x2[i])
			}
		} else {
			xnew[i] = x[i]
		}
		if xnew[i] < prms.FltMin[i] || xnew[i] > prms.FltMax[i] {
			xnew[i] = x[i]
		}
	}
}

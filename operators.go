// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/rnd"

// DiffEvol performs the differential-evolution operation
func DiffEvol(xnew, x, x0, x1, x2 []float64, prms *Parameters) {
	n := len(xnew)
	F := rnd.Float64(0.0, 1.0)
	I := rnd.Int(0, n-1)
	for i := 0; i < n; i++ {
		if rnd.FlipCoin(prms.DEC) || i == I {
			xnew[i] = x0[i] + F*(x1[i]-x2[i])
			if xnew[i] < 0 {
				xnew[i] = 0
			}
			if xnew[i] > 1 {
				xnew[i] = 1
			}
		} else {
			xnew[i] = x[i]
		}
	}
}

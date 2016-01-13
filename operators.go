// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/rnd"
)

// DiffEvol performs the differential-evolution operation
func DiffEvol(xnew, x, x0, x1, x2 []float64, prms *Parameters) {
	n := len(x)
	C := rnd.Float64(0.0, 1.0/math.Pow(float64(n), 0.1))
	F := rnd.Float64(0.0, 1.0)
	I := rnd.Int(0, n-1)
	for i := 0; i < n; i++ {
		if rnd.FlipCoin(C) || i == I {
			xnew[i] = x0[i] + F*(x1[i]-x2[i])
			if xnew[i] < prms.FltMin[i] {
				xnew[i] = prms.FltMin[i]
			}
			if xnew[i] > prms.FltMax[i] {
				xnew[i] = prms.FltMax[i]
			}
		} else {
			xnew[i] = x[i]
		}
	}
}

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
	var C float64
	n := len(x)
	if prms.Nova == 1 {
		C = 0.8 // this makes one-obj:prob-9 to work perfectly
	} else {
		//C = 0.8 / math.Pow(float64(n+prms.Nova), 0.5)
		//C = 0.1

		//xa, ya := 6.0, 0.8
		//xb, yb := 18.0, 0.01
		//C = yb + (yb-ya)*(float64(n)-xb)/(xb-xa)

		nn := float64(n)
		ca, cb, nb := 1.0, 0.01, 20.0
		if nn > nb {
			C = cb
		} else {
			C = cb + (ca-cb)*(math.Cos(math.Pi*nn/nb)+1.0)/2.0
		}
	}
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

// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/rnd"

// DiffEvol performs the differential-evolution operation
func DiffEvol(xnew, x, x0, x1, x2 []float64, prms *Parameters) {

	// normalise variables
	r, r0, r1, r2 := prms.Normalise4(x, x0, x1, x2)

	// perform DE
	n := len(xnew)
	F := rnd.Float64(0.0, 1.0)
	I := rnd.Int(0, n-1)
	for i := 0; i < n; i++ {
		if rnd.FlipCoin(prms.DEC) || i == I {
			xnew[i] = r0[i] + F*(r1[i]-r2[i])
			if prms.NormFlt {
				if xnew[i] < 0 {
					xnew[i] = 0
				}
				if xnew[i] > 1 {
					xnew[i] = 1
				}
			} else {
				if xnew[i] < prms.FltMin[i] {
					xnew[i] = prms.FltMin[i]
				}
				if xnew[i] > prms.FltMax[i] {
					xnew[i] = prms.FltMax[i]
				}
			}
		} else {
			xnew[i] = r[i]
		}
	}

	// de-normalise result
	prms.DeNormalise1(xnew)
}

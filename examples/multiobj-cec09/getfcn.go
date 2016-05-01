// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore
package main

/*
#cgo CFLAGS: -O3
#cgo LDFLAGS: -lm
#include "cec09.h"
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
)

func getfcn(problem int) (opt *goga.Optimiser) {

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.EpsH = 0.0001

	// dims
	nx := []int{30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
	nf := []int{2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3}
	ng := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 2, 2, 1, 1, 1}
	nh := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// get fcn
	idx := problem - 1
	var fcn goga.MinProb_t // functions
	switch problem {
	case 1:
		opt.RptName = "UF1"
		opt.Nsol = 100
		opt.FltMin = make([]float64, nx[idx])
		opt.FltMax = make([]float64, nx[idx])
		for i := 0; i < nx[idx]; i++ {
			opt.FltMin[i] = -1
			opt.FltMax[i] = 1
		}
		opt.FltMin[0] = 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.UF1(
				(*C.double)(unsafe.Pointer(&x[0])),
				(*C.double)(unsafe.Pointer(&f[0])),
				(C.int)(nx[idx]),
			)
		}
	default:
		chk.Panic("problem %d is not available", problem)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf[idx], ng[idx], nh[idx])
	return
}

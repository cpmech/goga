// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore
package main

import (
	"github.com/cpmech/goga"
	"github.com/cpmech/goga/examples/multiobj-cec09/cec09"
	"github.com/cpmech/gosl/chk"
)

func getfcn(problem string) (opt *goga.Optimiser) {

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.RptName = problem
	opt.EpsH = 0.0001
	opt.Nsol = 100

	// problem data
	opt.FltMin = cec09.Xmin[problem]
	opt.FltMax = cec09.Xmax[problem]
	nx := cec09.Nx[problem]
	nf := cec09.Nf[problem]
	ng := cec09.Nf[problem]
	nh := cec09.Nf[problem]
	chk.IntAssert(nx, len(opt.FltMin))
	chk.IntAssert(nx, len(opt.FltMax))

	// function
	var fcn goga.MinProb_t
	switch problem {
	case "UF1":
		fcn = cec09.UF1
	default:
		chk.Panic("problem %d is not available", problem)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)
	return
}

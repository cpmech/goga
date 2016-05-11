// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/goga"
	"github.com/cpmech/goga/examples/mulobj-cec09/cec09"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func solve_problem(problem string) (opt *goga.Optimiser) {

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.RptName = problem
	opt.EpsH = 0.0001
	opt.Nsamples = 1
	opt.Tf = 5000

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
	case "UF2":
		fcn = cec09.UF2
	case "UF3":
		fcn = cec09.UF3
	default:
		chk.Panic("problem %d is not available", problem)
	}

	nx = 3

	// load reference values
	opt.Multi_fStar = cec09.PFdata(problem)

	// number of trial solutions
	opt.Nsol = 600
	opt.Ncpu = 16
	if nf == 3 {
		opt.Nsol = 150
		opt.Ncpu = 3
	}
	if nf > 3 {
		opt.Nsol = 800
		opt.Ncpu = 8
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.RunMany("", "")
	goga.StatMulti(opt, true)

	// check
	goga.CheckFront0(opt, true)

	// plot
	if nf == 2 {
		plot2(opt, true)
	}
	plot3x(opt, false, 0, 1, 2, 0.02)
	return
}

func main() {
	P := []string{"UF3"}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem)
	}
	io.Pf("\n-------------------------- generating report --------------------------\nn")
	nRowPerTab := 10
	title := "CEC09 functions"
	goga.TexReport("/tmp/goga", "tmp_cec09", title, "cec09", 3, nRowPerTab, true, opts)
	goga.TexReport("/tmp/goga", "cec09", title, "cec09", 3, nRowPerTab, false, opts)
}

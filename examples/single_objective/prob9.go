// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
)

func main() {

	// GA parameters
	opt := new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 50
	opt.Ncpu = 1
	opt.Tf = 7000
	opt.Ntrials = 1000
	opt.EpsH = 1e-3
	opt.Verbose = false
	opt.GenType = "latin"

	// options for report
	opt.RptFmtF = "%.7f"
	opt.RptFmtFdev = "%.7f"
	opt.RptFmtX = "%.5f"

	opt.RptName = "9"
	opt.RptFref = []float64{0.0539498478}
	opt.RptXref = []float64{-1.717143, 1.595709, 1.827247, -0.7636413, -0.7636450}
	opt.FltMin = []float64{-2.3, -2.3, -3.2, -3.2, -3.2}
	opt.FltMax = []float64{+2.3, +2.3, +3.2, +3.2, +3.2}
	ng, nh := 0, 3
	fcn := func(f, g, h, x []float64, ξ []int, cpu int) {
		f[0] = math.Exp(x[0] * x[1] * x[2] * x[3] * x[4])
		h[0] = x[0]*x[0] + x[1]*x[1] + x[2]*x[2] + x[3]*x[3] + x[4]*x[4] - 10.0
		h[1] = x[1]*x[2] - 5.0*x[3]*x[4]
		h[2] = math.Pow(x[0], 3.0) + math.Pow(x[1], 3.0) + 1.0
	}

	// check
	if false {
		f := make([]float64, 1)
		h := make([]float64, 3)
		fcn(f, nil, h, opt.RptXref, nil, 0)
		io.Pforan("f(xref)  = %g  (%g)\n", f[0], opt.RptFref[0])
		io.Pforan("h0(xref) = %g\n", h[0])
		io.Pforan("h1(xref) = %g\n", h[1])
		io.Pforan("h2(xref) = %g\n", h[2])
		return
	}

	// initialise optimiser
	nf := 1
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	//opt.RunMany("/tmp/goga", "functions")
	opt.RunMany("", "")
	opt.HistNsta = 11
	goga.StatF(opt, 0, true)

	goga.TexSingleObjReport("/tmp/goga", "tmp_prob9", "probNine", 10, true, []*goga.Optimiser{opt})
	goga.TexSingleObjReport("/tmp/goga", "prob9", "probNine", 10, false, []*goga.Optimiser{opt})
}

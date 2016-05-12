// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "github.com/cpmech/goga"

func main() {

	// GA parameters
	opt := new(goga.Optimiser)
	opt.Default()
	opt.Tf = 500
	opt.Nsamples = 1000
	opt.Verbose = false
	opt.GenType = "latin"

	// enlarge box; add more constraint equations
	strategy2 := true

	// options for report
	opt.HistNsta = 6
	opt.HistLen = 13
	opt.RptFmtF = "%.5f"
	opt.RptFmtFdev = "%.2e"
	opt.RptFmtX = "%.3f"

	opt.Ncpu = 4
	opt.RptName = "2"
	opt.RptFref = []float64{-15.0}
	opt.RptXref = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 3, 3, 1}
	opt.FltMin = make([]float64, 13)
	opt.FltMax = make([]float64, 13)
	xmin, xmax := 0.0, 1.0
	if strategy2 {
		xmin, xmax = -0.5, 1.5
	}
	for i := 0; i < 9; i++ {
		opt.FltMin[i], opt.FltMax[i] = xmin, xmax
	}
	opt.FltMin[12], opt.FltMax[12] = xmin, xmax
	xmin, xmax = 0, 100
	if strategy2 {
		xmin, xmax = -1, 101
	}
	for i := 9; i < 12; i++ {
		opt.FltMin[i], opt.FltMax[i] = xmin, xmax
	}
	ng := 9
	if strategy2 {
		ng += 9 + 9 + 3 + 3 + 2
	}
	fcn := func(f, g, h, x []float64, ξ []int, cpu int) {
		s1, s2, s3 := 0.0, 0.0, 0.0
		for i := 0; i < 4; i++ {
			s1 += x[i]
			s2 += x[i] * x[i]
		}
		for i := 4; i < 13; i++ {
			s3 += x[i]
		}
		f[0] = 5.0*(s1-s2) - s3
		g[0] = 10.0 - 2.0*x[0] - 2.0*x[1] - x[9] - x[10]
		g[1] = 10.0 - 2.0*x[0] - 2.0*x[2] - x[9] - x[11]
		g[2] = 10.0 - 2.0*x[1] - 2.0*x[2] - x[10] - x[11]
		g[3] = 8.0*x[0] - x[9]
		g[4] = 8.0*x[1] - x[10]
		g[5] = 8.0*x[2] - x[11]
		g[6] = 2.0*x[3] + x[4] - x[9]
		g[7] = 2.0*x[5] + x[6] - x[10]
		g[8] = 2.0*x[7] + x[8] - x[11]
		if strategy2 {
			for i := 0; i < 9; i++ {
				g[9+i] = x[i]
				g[18+i] = 1.0 - x[i]
			}
			for i := 0; i < 3; i++ {
				g[27+i] = x[9+i]
				g[30+i] = 100.0 - x[9+i]
			}
			g[33] = x[12]
			g[34] = 1.0 - x[12]
		}
	}

	// number of trial solutions
	opt.Nsol = len(opt.FltMin) * 10

	// initialise optimiser
	nf, nh := 1, 0
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.RunMany("", "")
	goga.StatF(opt, 0, true)
	opts := []*goga.Optimiser{opt}
	textSize := `\scriptsize  \setlength{\tabcolsep}{0.5em}`
	miniPageSz, histTextSize := "4.1cm", `\fontsize{5pt}{6pt}`
	nRowPerTab := 9
	title := "Constrained single objective problem 2"
	goga.TexReport("/tmp/goga", "tmp_one-obj-prob2", title, "one-ob-prob2", 1, nRowPerTab, true, false, textSize, miniPageSz, histTextSize, opts)
	goga.TexReport("/tmp/goga", "one-obj-prob2", title, "one-obj-prob2", 1, nRowPerTab, false, false, textSize, miniPageSz, histTextSize, opts)
}

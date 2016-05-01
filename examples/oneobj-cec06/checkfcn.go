// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/gosl/io"
)

func check(problem int, tolf, tolg float64) {
	io.Pf("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> problem %d <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", problem)
	opt := getfcn(problem)
	x := opt.RptXref
	f := make([]float64, opt.Nf)
	g := make([]float64, opt.Ng)
	h := make([]float64, opt.Nh)
	opt.MinProb(f, g, h, x, nil, 0)
	err := math.Abs(f[0] - opt.RptFref[0])
	io.Pforan("nx=%d nf=%d ng=%d nh=%d\n", opt.Nflt, opt.Nf, opt.Ng, opt.Nh)
	io.Pforan("x = %v\n", x)
	io.Pforan("f = %v  err = %v\n", f, err)
	io.Pforan("g = %v\n", g)
	io.Pforan("h = %v\n", h)
	for i := 0; i < opt.Ng; i++ {
		if g[i] < 0 {
			io.PfRed("unfeasible on g\n")
		}
	}
	for i := 0; i < opt.Nh; i++ {
		if math.Abs(h[i]) > opt.EpsH {
			io.PfRed("unfeasible on h\n")
		}
	}
	if err > tolf {
		io.PfRed("err is too big\n")
	}
}

func main() {
	tolf, tolg := 1e-10, 1e-10
	for i := 1; i <= 24; i++ {
		check(i, tolf, tolg)
	}
}

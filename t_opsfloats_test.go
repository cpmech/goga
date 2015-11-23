// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_cxdeb01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cxdeb01. Deb's crossover")

	var ops OpsData
	ops.SetDefault()
	ops.Xrange = [][]float64{{-3, 3}, {-4, 4}}

	rnd.Init(0)

	A := []float64{-1, 1}
	B := []float64{1, 2}
	a := make([]float64, len(A))
	b := make([]float64, len(A))
	FltCrossoverDB(a, b, A, B, nil, nil, 0, &ops)
	io.Pforan("A = %v\n", A)
	io.Pforan("B = %v\n", B)
	io.Pfcyan("a = %.6f\n", a)
	io.Pfcyan("b = %.6f\n", b)

	nsamples := 1000
	a0s, a1s := make([]float64, nsamples), make([]float64, nsamples)
	b0s, b1s := make([]float64, nsamples), make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		FltCrossoverDB(a, b, B, A, nil, nil, 0, &ops)
		a0s[i], a1s[i] = a[0], a[1]
		b0s[i], b1s[i] = b[0], b[1]
	}
	ha0 := rnd.Histogram{Stations: []float64{-4, -3.5, -3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1}}
	hb0 := rnd.Histogram{Stations: []float64{0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 5, 5.5, 6}}
	ha1 := rnd.Histogram{Stations: utl.LinSpace(-4, 4, 11)}
	hb1 := rnd.Histogram{Stations: utl.LinSpace(-4, 4, 11)}
	ha0.Count(a0s, true)
	hb0.Count(b0s, true)
	ha1.Count(a1s, true)
	hb1.Count(b1s, true)

	io.Pforan("\na0s\n")
	io.Pf("%s", rnd.TextHist(ha0.GenLabels("%.1f"), ha0.Counts, 60))
	io.Pforan("b0s\n")
	io.Pf("%s", rnd.TextHist(hb0.GenLabels("%.1f"), hb0.Counts, 60))

	io.Pforan("\na1s\n")
	io.Pf("%s", rnd.TextHist(ha1.GenLabels("%.1f"), ha1.Counts, 60))
	io.Pforan("b1s\n")
	io.Pf("%s", rnd.TextHist(hb1.GenLabels("%.1f"), hb1.Counts, 60))
}

func Test_mtdeb01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mtdeb01. Deb's mutation")

	var ops OpsData
	ops.SetDefault()
	ops.Xrange = [][]float64{{-3, 3}, {-4, 4}}

	rnd.Init(0)

	A := []float64{-1, 1}
	io.Pforan("before: A = %v\n", A)
	FltMutationDB(A, 10, &ops)
	io.Pforan("after:  A = %v\n", A)

	ha0 := rnd.Histogram{Stations: utl.LinSpace(-3, 3, 11)}

	nsamples := 1000
	aa := make([]float64, len(A))
	a0s := make([]float64, nsamples)
	for _, t := range []int{0, 50, 100} {
		for i := 0; i < nsamples; i++ {
			copy(aa, A)
			FltMutationDB(aa, t, &ops)
			a0s[i] = aa[0]
		}
		ha0.Count(a0s, true)
		io.Pf("\ntime = %d\n", t)
		io.Pf("%s", rnd.TextHist(ha0.GenLabels("%.1f"), ha0.Counts, 60))
	}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	PI := math.Pi

	yf := func(x float64) float64 {
		return 1.0 - math.Sqrt(x) - math.Sin(10.0*math.Pi*x)*x
	}

	dydx := func(x float64) float64 {
		return -math.Sin(10.0*PI*x) - 10.0*PI*x*math.Cos(10.0*PI*x) - 1.0/(2.0*math.Sqrt(x))
	}

	var nlsDY num.NlSolver
	nlsDY.Init(1, func(fx, x []float64) error {
		fx[0] = dydx(x[0])
		return nil
	}, nil, nil, false, true, nil)
	defer nlsDY.Clean()

	X := []float64{0.09, 0.25, 0.45, 0.65, 0.85}
	Y := make([]float64, len(X))
	for i, x := range X {

		// find min
		xx := []float64{x}
		err := nlsDY.Solve(xx, true)
		if err != nil {
			io.PfRed("dydx nls failed:\n%v", err)
			return
		}
		X[i] = xx[0]
		Y[i] = yf(X[i])
	}

	// find next point along horizontal line
	Xnext := []float64{0.2, 0.4, 0.6, 0.8}
	Ynext := make([]float64, len(Xnext))
	for i, xnext := range Xnext {
		var nlsX num.NlSolver
		nlsX.Init(1, func(fx, x []float64) error {
			fx[0] = Y[i] - yf(x[0])
			return nil
		}, nil, nil, false, true, nil)
		defer nlsX.Clean()
		xx := []float64{xnext}
		err := nlsX.Solve(xx, true)
		if err != nil {
			io.PfRed("dydx nls failed:\n%v", err)
			return
		}
		Xnext[i] = xx[0]
		Ynext[i] = yf(Xnext[i])
	}

	// auxiliary points
	XX := []float64{
		0, X[0],
		Xnext[0], X[1],
		Xnext[1], X[2],
		Xnext[2], X[3],
		Xnext[3], X[4],
	}
	YY := []float64{
		1, Y[0],
		Ynext[0], Y[1],
		Ynext[1], Y[2],
		Ynext[2], Y[3],
		Ynext[3], Y[4],
	}
	io.Pforan("XX = %.3f\n", XX)
	io.Pforan("YY = %.3f\n", YY)

	// find arc-length
	arclen := 0.0
	for i := 0; i < len(XX); i += 2 {
		a := XX[i]
		b := XX[i+1]
		if i == 0 {
			a += 1e-7
		}
		var quad num.Simp
		quad.Init(func(x float64) float64 {
			return math.Sqrt(1.0 + math.Pow(dydx(x), 2.0))
		}, a, b, 1e-4)
		res, err := quad.Integrate()
		if err != nil {
			io.PfRed("quad failed:\n%v", err)
			return
		}
		arclen += res
		io.Pf("int(...) from %.15f to %.15f = %g\n", a, b, res)
	}
	io.Pforan("arclen = %v\n", arclen)

	np := 201
	xx := utl.LinSpace(0, 1, np)
	yy := make([]float64, np)
	for i := 0; i < np; i++ {
		yy[i] = yf(xx[i])
	}
	plt.Plot(xx, yy, "'b-', clip_on=0")
	for i, x := range X {
		plt.PlotOne(x, Y[i], "'r|', mew=2, ms=30, clip_on=0")
	}
	for i, x := range Xnext {
		plt.PlotOne(x, Ynext[i], "'r_', mew=2, ms=30, clip_on=0")
	}
	for i := 0; i < len(XX); i += 2 {
		x0, y0 := XX[i], YY[i]
		x1, y1 := XX[i+1], YY[i+1]
		plt.Arrow(x0, y0, x1, y1, "")
	}
	plt.SetXnticks(11)
	plt.Gll("x", "y", "")
	plt.SaveD("/tmp/goga", "calcZDT3pts.eps")
}

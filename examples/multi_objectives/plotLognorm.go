// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func main() {

	μ := 0.5
	σ := 0.25
	ori := false
	xmin := 0.0
	xmax := 3.0

	//μ = -0.8
	//σ = 0.4
	//ori = true

	var dist rnd.DistLogNormal
	dist.Init(&rnd.VarData{M: μ, S: σ, Pori: ori})
	io.Pforan("μ=%g σ=%v\n", dist.M, dist.S)

	n := 101
	X := utl.LinSpace(xmin, xmax, n)
	Y := make([]float64, n)
	for i := 0; i < n; i++ {
		Y[i] = dist.Pdf(X[i])
	}

	nsamples := 1000
	Xs := make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		Xs[i] = rnd.Lognormal(μ, σ, ori)
	}

	var hist rnd.Histogram
	hist.Stations = utl.LinSpace(xmin, xmax, 21)
	hist.Count(Xs, true)

	plt.SetForEps(0.8, 280)
	plt.Plot(X, Y, io.Sf("clip_on=0,zorder=10,label=r'$\\mu=%.4f,\\;\\sigma=%.4f$'", dist.M, dist.S))
	hist.PlotDensity(nil, "")
	plt.Gll("$x$", "$f(x)$", "leg_out=1, leg_ncol=2")
	plt.SaveD("/tmp/goga", "figLognorm.eps")
}

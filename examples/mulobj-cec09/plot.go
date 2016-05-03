// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {
	opt := getfcn("UF1")

	feasibleOnly := true
	plt.SetForEps(0.8, 300)
	fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
	fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
	goga.PlotOvaOvaPareto(opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
	np := 201
	F0 := utl.LinSpace(fmin[0], fmax[0], np)
	F1 := make([]float64, np)
	for i := 0; i < np; i++ {
		F1[i] = opt.F1F0_func(F0[i])
	}
	plt.Plot(F0, F1, io.Sf("'b-', label='%s'", opt.RptName))
	for _, f0vals := range opt.F1F0_f0ranges {
		f0A, f0B := f0vals[0], f0vals[1]
		f1A, f1B := opt.F1F0_func(f0A), opt.F1F0_func(f0B)
		plt.PlotOne(f0A, f1A, "'g_', mew=1.5, ms=10, clip_on=0")
		plt.PlotOne(f0B, f1B, "'g|', mew=1.5, ms=10, clip_on=0")
	}
	plt.AxisRange(fmin[0], fmax[0], fmin[1], fmax[1])
	plt.Gll("$f_0$", "$f_1$", "")
	plt.SaveD("/tmp/goga", io.Sf("%s.eps", opt.RptName))

}

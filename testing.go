// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func test_ova_plot_run_plot(fnkey string, opt *Optimiser, np int, ovaMult float64, fcn func(x float64) float64) {

	// plot initial solution
	if chk.Verbose {
		X := utl.LinSpace(opt.FltMin[0], opt.FltMax[0], np)
		Y := make([]float64, np)
		for i := 0; i < np; i++ {
			Y[i] = fcn(X[i])
		}
		plt.SetForEps(0.8, 455)
		plt.Plot(X, Y, "'b-'")
		opt.PlotSolOvas(0, 0, ovaMult, plt.Fmt{L: "initial", M: "o", C: "k", Ls: "none", Ms: 3}, false)
	}

	// solve
	opt.Solve()

	// sort
	SortByOva(opt.Solutions, 0)
	best := opt.Solutions[0]

	// plot final solution
	if chk.Verbose {
		opt.PlotSolOvas(0, 0, ovaMult, plt.Fmt{L: "final", M: "o", C: "r", Ls: "none", Ms: 6}, true)
		plt.PlotOne(best.Flt[0], best.Ova[0]*ovaMult, "'g*', markeredgecolor='g', label='best', clip_on=0, zorder=20")
		plt.Equal()
		plt.Gll("$x$", "$y$", "leg_out=1, leg_ncol=4, leg_hlen=1.5")
		plt.SaveD("/tmp/goga", fnkey+".eps")
	}
}

func test_contour_plot_run_plot(fnkey string, opt *Optimiser) {

	// plot initial solution
	if chk.Verbose {
		plt.SetForEps(0.8, 455)
		opt.PlotContour(0, 0, 1, ContourParams{})
		opt.PlotSolutions(0, 1, plt.Fmt{L: "initial", M: "o", C: "k", Ls: "none", Ms: 3}, false)
	}

	// solve
	opt.Solve()

	// sort
	SortByOva(opt.Solutions, 0)
	best := opt.Solutions[0]

	// plot final solution
	if chk.Verbose {
		opt.PlotSolutions(0, 1, plt.Fmt{L: "final", M: "o", C: "k", Ls: "none", Ms: 6}, true)
		plt.PlotOne(best.Flt[0], best.Flt[1], "'g*', markeredgecolor='g', label='best', clip_on=0, zorder=20")
		plt.Equal()
		plt.Gll("$x_0$", "$x_1$", "leg_out=1, leg_ncol=4, leg_hlen=1.5")
		plt.SaveD("/tmp/goga", fnkey+".eps")
	}
}

func test_pareto_plot_run_plot(fnkey string, opt *Optimiser, initial bool) {

	// plot initial solution
	if chk.Verbose {
		plt.SetForEps(0.8, 455)
		if initial {
			opt.PlotOvas(0, 1, plt.Fmt{L: "initial", M: "o", C: "k", Ls: "none", Ms: 3}, false)
		}
	}

	// solve
	opt.Solve()

	// plot final solution
	if chk.Verbose {
		opt.PlotOvas(0, 1, plt.Fmt{L: "final", M: "o", C: "r", Ls: "none", Ms: 4}, false)
		plt.Gll("$f_0$", "$f_1$", "leg_out=1, leg_ncol=4, leg_hlen=1.5")
		plt.SaveD("/tmp/goga", fnkey+".eps")
	}
}

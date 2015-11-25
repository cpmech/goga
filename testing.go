// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
)

func test_contour_plot_run_plot(fnkey string, opt *Optimiser) {

	// plot initial solution
	if chk.Verbose {
		plt.SetForEps(0.8, 455)
		opt.PlotContour(0, 0, 1, ContourParams{})
		opt.PlotSolutions(0, 1, plt.Fmt{M: "o", C: "k", Ls: "none", Ms: 3}, false)
	}

	// solve
	opt.Solve()

	// plot final solution
	if chk.Verbose {
		opt.PlotSolutions(0, 1, plt.Fmt{M: "o", C: "k", Ls: "none", Ms: 6}, true)
		plt.Equal()
		plt.Gll("$x_0$", "$x_1$", "")
		plt.SaveD("/tmp/goga", fnkey+".eps")
	}
}

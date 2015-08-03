// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

type TwoVarsObjFunc_t func(x []float64) float64    // len(x) == 2
type TwoVarsConstraint_t func(x []float64) float64 // len(x) == 2

// PlotTwoVarsContour plots contour for two variables problem. len(x) == 2
// c   -- constraint. can be nil
// isl -- island. can be nil
func PlotTwoVarsContour(xmin, xmax []float64, g TwoVarsObjFunc_t, c TwoVarsConstraint_t, isl *Island, extra func()) {
	chk.IntAssert(len(xmin), 2)
	chk.IntAssert(len(xmax), 2)
	plt.SetForEps(0.8, 350)
	np := 41
	X, Y := utl.MeshGrid2D(xmin[0], xmax[0], xmin[1], xmax[1], np, np)
	Z := la.MatAlloc(np, np)
	var C [][]float64
	if c != nil {
		C = la.MatAlloc(np, np)
	}
	x := make([]float64, 2)
	for i := 0; i < np; i++ {
		for j := 0; j < np; j++ {
			x[0], x[1] = X[i][j], Y[i][j]
			Z[i][j] = g(x)
			if c != nil {
				C[i][j] = c(x)
			}
		}
	}
	plt.Contour(X, Y, Z, "")
	if c != nil {
		plt.ContourSimple(X, Y, C, "levels=[0], colors=['yellow'], linewidths=[2]")
	}
	if isl != nil {
		for _, ind := range isl.Pop {
			x := ind.GetFloat(0)
			y := ind.GetFloat(1)
			plt.PlotOne(x, y, "'k.'")
		}
	}
	if extra != nil {
		extra()
	}
}

// PlotTwoVarsSave plots final population and save figure
//  xmin, xmax -- range. use nil for automatic
func PlotTwoVarsSave(dirout, fname string, xmin, xmax []float64, isl *Island, best *Individual) {
	if fname == "" {
		return
	}
	if dirout == "" {
		dirout = "/tmp/goga"
	}
	if isl != nil {
		for _, ind := range isl.Pop {
			x := ind.GetFloat(0)
			y := ind.GetFloat(1)
			plt.PlotOne(x, y, "'m*'")
		}
	}
	if best != nil {
		x := best.GetFloat(0)
		y := best.GetFloat(1)
		plt.PlotOne(x, y, "'g*', ms=8")
	}
	plt.Equal()
	if xmin != nil {
		plt.AxisRange(xmin[0], xmax[0], xmin[1], xmax[1])
	}
	plt.SaveD(dirout, fname)
}

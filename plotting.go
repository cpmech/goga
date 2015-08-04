// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// TwoVarsFunc_t defines a function to plot contours (len(x)==2)
type TwoVarsFunc_t func(x []float64) float64

// PlotTwoVarsContour plots contour for two variables problem. len(x) == 2
// dirout  -- directory to save files
// fnkey   -- file name key for eps figure
// pop0    -- initial population. can be <nil> if individuals are not to be plotted
// pop1    -- final population. can be <nil> if individuals are not to be plotted
// best    -- best individual. can be <nil>
// xmin    -- min x[0] and x[1]. use <nil>for automatic values
// xmax    -- max x[0] and x[1]. use <nil>for automatic values
// np      -- number of points for contour
// axrange -- axes range: if true, use xmin and xmax
// extra   -- called just before saving figure
// f       -- function to plot filled contour. can be <nil>
// gfs     -- functions to plot contour @ level 0. can be <nil>
func PlotTwoVarsContour(dirout, fnkey string, pop0, pop1 Population, best *Individual,
	xmin, xmax []float64, np int, axrange bool, extra func(),
	f TwoVarsFunc_t, gfs ...TwoVarsFunc_t) {
	if fnkey == "" {
		return
	}
	chk.IntAssert(len(xmin), 2)
	chk.IntAssert(len(xmax), 2)
	X, Y := utl.MeshGrid2D(xmin[0], xmax[0], xmin[1], xmax[1], np, np)
	var Zf [][]float64
	var Zg [][][]float64
	if f != nil {
		Zf = la.MatAlloc(np, np)
	}
	if len(gfs) > 0 {
		Zg = utl.Deep3alloc(len(gfs), np, np)
	}
	x := make([]float64, 2)
	for i := 0; i < np; i++ {
		for j := 0; j < np; j++ {
			x[0], x[1] = X[i][j], Y[i][j]
			if f != nil {
				Zf[i][j] = f(x)
			}
			for k, g := range gfs {
				Zg[k][i][j] = g(x)
			}
		}
	}
	plt.Reset()
	plt.SetForEps(0.8, 400)
	if f != nil {
		plt.Contour(X, Y, Zf, "")
	}
	for k, _ := range gfs {
		plt.ContourSimple(X, Y, Zg[k], "levels=[0], colors=['yellow'], linewidths=[2], clip_on=0")
	}
	if pop0 != nil {
		for _, ind := range pop0 {
			x := ind.GetFloat(0)
			y := ind.GetFloat(1)
			plt.PlotOne(x, y, "'k.', zorder=20, clip_on=0")
		}
	}
	if pop1 != nil {
		for _, ind := range pop1 {
			x := ind.GetFloat(0)
			y := ind.GetFloat(1)
			plt.PlotOne(x, y, "'m*', zorder=30, clip_on=0")
		}
	}
	if extra != nil {
		extra()
	}
	if best != nil {
		x := best.GetFloat(0)
		y := best.GetFloat(1)
		plt.PlotOne(x, y, "'g*', ms=8, zorder=50, clip_on=0")
	}
	if dirout == "" {
		dirout = "/tmp/goga"
	}
	plt.Cross("clr='grey'")
	plt.SetXnticks(11)
	plt.SetYnticks(11)
	plt.Equal()
	if axrange {
		plt.AxisRange(xmin[0], xmax[0], xmin[1], xmax[1])
	}
	plt.SaveD(dirout, fnkey+".eps")
}

// PlotOvs plots objective values versus time
func PlotOvs(isl *Island, ext, args string, t0, tf int, withtxt bool, numfmt string, first, last bool) {
	if isl.C.DoPlot == false || isl.C.FnKey == "" {
		return
	}
	if first {
		plt.SetForEps(0.75, 250)
	}
	var y []float64
	if tf == -1 {
		y = isl.OVA[t0:]
		tf = len(isl.OVA)
	} else {
		y = isl.OVA[t0:tf]
	}
	n := len(y)
	T := utl.LinSpace(float64(t0), float64(tf), n)
	me := (tf-t0)/10 + isl.Id
	if me < 1 {
		me = 1
	}
	if len(args) > 0 {
		args += ","
	}
	plt.Plot(T, y, io.Sf("%s marker='%s', markersize=%d, markevery=%d, zorder=10, clip_on=0", args, get_marker(isl.Id), get_mrksz(isl.Id), me))
	if withtxt {
		plt.Text(T[0], y[0], io.Sf(numfmt, y[0]), "ha='left'")
		plt.Text(T[n-1], y[n-1], io.Sf(numfmt, y[n-1]), "ha='right'")
	}
	if last {
		plt.Gll("time", "objective value", "")
		plt.SaveD(isl.C.DirOut, isl.C.FnKey+"_ova"+ext)
	}
}

// PlotOor plots out-of-range values versus time
func PlotOor(isl *Island, ext, args string, t0, tf int, withtxt bool, numfmt string, first, last bool) {
	if isl.C.DoPlot == false || isl.C.FnKey == "" {
		return
	}
	if first {
		plt.SetForEps(0.75, 250)
	}
	var y []float64
	if tf == -1 {
		y = isl.OOR[t0:]
		tf = len(isl.OOR)
	} else {
		y = isl.OOR[t0:tf]
	}
	n := len(y)
	T := utl.LinSpace(float64(t0), float64(tf), n)
	me := (tf-t0)/10 + isl.Id
	if me < 1 {
		me = 1
	}
	if len(args) > 0 {
		args += ","
	}
	plt.Plot(T, y, io.Sf("%s marker='%s', markersize=%d, markevery=%d, zorder=10, clip_on=0", args, get_marker(isl.Id), get_mrksz(isl.Id), me))
	if withtxt {
		plt.Text(T[0], y[0], io.Sf(numfmt, y[0]), "ha='left'")
		plt.Text(T[n-1], y[n-1], io.Sf(numfmt, y[n-1]), "ha='right'")
	}
	if last {
		plt.Gll("time", "out-of-range value", "")
		plt.SaveD(isl.C.DirOut, isl.C.FnKey+"_oor"+ext)
	}
}

// get_marker returns a marker for graphs
func get_marker(i int) string {
	pool := []string{"", "+", ".", "x", "s", "o", "*"}
	return pool[i%len(pool)]
}

// get_mrksz returns a marker size for graphs
func get_mrksz(i int) int {
	pool := []int{6, 6, 6, 3, 6, 6, 6}
	return pool[i%len(pool)]
}

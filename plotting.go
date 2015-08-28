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

// TwoVarsTrans_t defines a tranformation x → y (len(x)==len(y)==2)
type TwoVarsTrans_t func(x []float64) (y []float64, invalid bool)

// PlotTwoVarsContour plots contour for two variables problem. len(x) == 2
//  Input:
//   dirout  -- directory to save files
//   fnkey   -- file name key for eps figure
//   pop0    -- initial population. can be <nil> if individuals are not to be plotted
//   pop1    -- final population. can be <nil> if individuals are not to be plotted
//   best    -- best individual. can be <nil>
//   np      -- number of points for contour
//   lw_g    -- linewidth for g functions
//   cargs   -- arguments for contour command
//   extra   -- called just before saving figure
//   csimple -- use simple contour for f function
//   axequal -- axis.equal
//   vrange  -- [2][2] range of x and y values; e.g.: [][]float64{{xmin,xmax},{ymin,ymax}}
//   vmax    -- max 1 values
//   istrans -- vrange and individuals are transformed y-values; otherwise they are x-values
//   tplot   -- plot transformed plot; needs T and Ti.
//   T       -- transformation: x → y
//   Ti      -- transformation: y → x
//   f       -- function to plot filled contour. can be <nil>
//   gs      -- functions to plot contour @ level 0. can be <nil>
//  Note: g(x) operates on original x values
func PlotTwoVarsContour(dirout, fnkey string, pop0, pop1 Population, best *Individual, np int, lw_g float64, cargs string, extra func(), csimple, axequal bool,
	vrange [][]float64, istrans, tplot bool, T, Ti TwoVarsTrans_t, f TwoVarsFunc_t, gs ...TwoVarsFunc_t) {
	if fnkey == "" {
		return
	}
	chk.IntAssert(len(vrange), 2)
	V0, V1 := utl.MeshGrid2D(vrange[0][0], vrange[0][1], vrange[1][0], vrange[1][1], np, np)
	var Zf [][]float64
	var Zg [][][]float64
	if f != nil {
		Zf = la.MatAlloc(np, np)
	}
	if len(gs) > 0 {
		Zg = utl.Deep3alloc(len(gs), np, np)
	}
	dotrans := !istrans && tplot // do transform
	untrans := istrans && !tplot // un-transform
	x := make([]float64, 2)
	for i := 0; i < np; i++ {
		for j := 0; j < np; j++ {
			if istrans {
				x, invalid := Ti([]float64{V0[i][j], V1[i][j]}) // x ← T⁻¹(y)
				if invalid {
					chk.Panic("cannot plot contour due to invalid transformation")
				}
				if !tplot {
					V0[i][j], V1[i][j] = x[0], x[1]
				}
			} else {
				x[0], x[1] = V0[i][j], V1[i][j]
				if tplot {
					y, invalid := T(x) // v ← y = T(x)
					if invalid {
						chk.Panic("cannot plot contour due to invalid transformation")
					}
					V0[i][j], V1[i][j] = y[0], y[1]
				}
			}
			if f != nil {
				Zf[i][j] = f(x)
			}
			for k, g := range gs {
				Zg[k][i][j] = g(x)
			}
		}
	}
	plt.Reset()
	plt.SetForEps(0.8, 350)
	if f != nil {
		cmapidx := 0
		if tplot {
			cmapidx = 4
		}
		if cargs != "" {
			cargs = "," + cargs
		}
		if csimple {
			plt.ContourSimple(V0, V1, Zf, true, 7, "colos=['k'], fsz=7"+cargs)
		} else {
			plt.Contour(V0, V1, Zf, io.Sf("fsz=7, cmapidx=%d"+cargs, cmapidx))
		}
	}
	for k, _ := range gs {
		plt.ContourSimple(V0, V1, Zg[k], false, 7, io.Sf("zorder=5, levels=[0], colors=['yellow'], linewidths=[%g], clip_on=0", lw_g))
	}
	get_v := func(ind *Individual) (v []float64) {
		v = ind.GetFloats()
		if dotrans {
			y, invalid := T(v)
			if invalid {
				chk.Panic("cannot plot contour due to invalid transformation")
			}
			v[0], v[1] = y[0], y[1]
		}
		if untrans {
			x, invalid := Ti(v)
			if invalid {
				chk.Panic("cannot plot contour due to invalid transformation")
			}
			v[0], v[1] = x[0], x[1]
		}
		return
	}
	if pop0 != nil {
		for i, ind := range pop0 {
			l := ""
			if i == 0 {
				l = "initial population"
			}
			v := get_v(ind)
			plt.PlotOne(v[0], v[1], io.Sf("'k.', zorder=20, clip_on=0, label='%s'", l))
		}
	}
	if pop1 != nil {
		for i, ind := range pop1 {
			l := ""
			if i == 0 {
				l = "final population"
			}
			v := get_v(ind)
			plt.PlotOne(v[0], v[1], io.Sf("'ko', ms=6, zorder=30, clip_on=0, label='%s', markerfacecolor='none'", l))
		}
	}
	if extra != nil {
		extra()
	}
	if best != nil {
		v := get_v(best)
		plt.PlotOne(v[0], v[1], "'m*', zorder=50, clip_on=0, label='best', markeredgecolor='m'")
	}
	if dirout == "" {
		dirout = "."
	}
	plt.Cross("clr='grey'")
	if axequal {
		plt.Equal()
	}
	urange := vrange
	if istrans && !tplot {
		vmin := []float64{vrange[0][0], vrange[1][0]}
		xmin, invalid := Ti(vmin)
		if invalid {
			chk.Panic("cannot plot contour due to invalid transformation")
		}
		vmax := []float64{vrange[0][1], vrange[1][1]}
		xmax, invalid := Ti(vmax)
		if invalid {
			chk.Panic("cannot plot contour due to invalid transformation")
		}
		urange = [][]float64{{xmin[0], xmax[0]}, {xmin[1], xmax[1]}}
	}
	if !istrans && tplot {
		vmin := []float64{vrange[0][0], vrange[1][0]}
		ymin, invalid := T(vmin)
		if invalid {
			chk.Panic("cannot plot contour due to invalid transformation")
		}
		vmax := []float64{vrange[0][1], vrange[1][1]}
		ymax, invalid := T(vmax)
		if invalid {
			chk.Panic("cannot plot contour due to invalid transformation")
		}
		urange = [][]float64{{ymin[0], ymax[0]}, {ymin[1], ymax[1]}}
	}
	plt.AxisRange(urange[0][0], urange[0][1], urange[1][0], urange[1][1])
	args := "leg_out=1, leg_ncol=4, leg_hlen=1.5"
	if tplot {
		plt.Gll("$y_0$", "$y_1$", args)
	} else {
		plt.Gll("$x_0$", "$x_1$", args)
	}
	plt.SaveD(dirout, fnkey+".eps")
}

// PlotOvs plots objective values versus time
func PlotOvs(isl *Island, ext, args string, t0, tf int, first, last bool) {
	if isl.C.DoPlot == false || isl.C.FnKey == "" {
		return
	}
	if first {
		plt.SetForEps(0.75, 250)
	}
	me := (tf-t0)/20 + isl.Id
	if me < 1 {
		me = 1
	}
	if len(args) > 0 {
		args += ","
	}
	nova := len(isl.Pop[0].Ovas)
	for i := 0; i < nova; i++ {
		plt.Plot(isl.OutTimes[t0:], isl.OutOvas[i][t0:], io.Sf("%s marker='%s', markersize=%d, markevery=%d, zorder=10, clip_on=0", args, get_marker(isl.Id), get_mrksz(isl.Id), me))
	}
	if last {
		plt.Gll("time", "objective value", "")
		plt.SaveD(isl.C.DirOut, isl.C.FnKey+"_ova"+ext)
	}
}

// PlotOor plots out-of-range values versus time
func PlotOor(isl *Island, ext, args string, t0, tf int, first, last bool) {
	if isl.C.DoPlot == false || isl.C.FnKey == "" {
		return
	}
	if first {
		plt.SetForEps(0.75, 250)
	}
	me := (tf-t0)/20 + isl.Id
	if me < 1 {
		me = 1
	}
	if len(args) > 0 {
		args += ","
	}
	noor := len(isl.Pop[0].Oors)
	for i := 0; i < noor; i++ {
		plt.Plot(isl.OutTimes[t0:], isl.OutOors[i][t0:], io.Sf("%s marker='%s', markersize=%d, markevery=%d, zorder=10, clip_on=0", args, get_marker(isl.Id), get_mrksz(isl.Id), me))
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

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// ContourParams holds parameters to plot contours
type ContourParams struct {
	Npts    int       // number of points for contour
	CmapIdx int       // colormap index
	Csimple bool      // simple contour
	AxEqual bool      // axes-equal
	Lwg     float64   // linewidth for g functions
	Lwh     float64   // linewidth for h functions
	Args    string    // extra arguments for plot
	Extra   func()    // extra function
	Xrange  []float64 // to override x-range
	Yrange  []float64 // to override y-range
}

// PlotContour plots contour
func (o *Optimiser) PlotContour(iFlt, jFlt, iOva int, prms ContourParams) {

	// fix parameters
	if prms.Npts < 3 {
		prms.Npts = 41
	}
	if prms.Lwg < 0.1 {
		prms.Lwg = 1.5
	}
	if prms.Lwh < 0.1 {
		prms.Lwh = 1.5
	}

	// limits and meshgrid
	xmin, xmax := o.Xmin[iFlt], o.Xmax[iFlt]
	ymin, ymax := o.Xmin[jFlt], o.Xmax[jFlt]
	if prms.Xrange != nil {
		xmin, xmax = prms.Xrange[0], prms.Xrange[1]
	}
	if prms.Yrange != nil {
		ymin, ymax = prms.Yrange[0], prms.Yrange[1]
	}

	// auxiliary variables
	X, Y := utl.MeshGrid2D(xmin, xmax, ymin, ymax, prms.Npts, prms.Npts)
	Zf := utl.DblsAlloc(prms.Npts, prms.Npts)
	var Zg [][][]float64
	var Zh [][][]float64
	if o.Ng > 0 {
		Zg = utl.Deep3alloc(o.Ng, prms.Npts, prms.Npts)
	}
	if o.Nh > 0 {
		Zh = utl.Deep3alloc(o.Nh, prms.Npts, prms.Npts)
	}

	// compute values
	x := make([]float64, 2)
	grp := 0
	for i := 0; i < prms.Npts; i++ {
		for j := 0; j < prms.Npts; j++ {
			x[0], x[1] = X[i][j], Y[i][j]
			o.MinProb(o.F[grp], o.G[grp], o.H[grp], x, nil, grp)
			Zf[i][j] = o.F[grp][iOva]
			for k, g := range o.G[grp] {
				Zg[k][i][j] = g
			}
			for k, h := range o.H[grp] {
				Zh[k][i][j] = h
			}
		}
	}

	// plot f
	if prms.Csimple {
		plt.ContourSimple(X, Y, Zf, true, 7, "colors=['k'], fsz=7, "+prms.Args)
	} else {
		plt.Contour(X, Y, Zf, io.Sf("fsz=7, cmapidx=%d, "+prms.Args, prms.CmapIdx))
	}

	// plot g
	clr := "yellow"
	if prms.Csimple {
		clr = "blue"
	}
	for _, g := range Zg {
		plt.ContourSimple(X, Y, g, false, 7, io.Sf("zorder=5, levels=[0], colors=['%s'], linewidths=[%g], clip_on=0", clr, prms.Lwg))
	}

	// plot h
	clr = "yellow"
	if prms.Csimple {
		clr = "blue"
	}
	for _, h := range Zh {
		plt.ContourSimple(X, Y, h, false, 7, io.Sf("zorder=5, levels=[0], colors=['%s'], linewidths=[%g], clip_on=0", clr, prms.Lwh))
	}
}

// PlotAddFltFlt adds flt-flt points to existent plot
func (o *Optimiser) PlotAddFltFlt(iFlt, jFlt int, sols []*Solution, denormalise bool, fmt plt.Fmt, emptyMarker bool) {
	nsol := len(sols)
	x, y := make([]float64, nsol), make([]float64, nsol)
	if denormalise {
		for i, sol := range sols {
			x[i] = o.Xmin[iFlt] + sol.Flt[iFlt]*o.Dx[iFlt]
			y[i] = o.Xmin[jFlt] + sol.Flt[jFlt]*o.Dx[jFlt]
		}
	} else {
		for i, sol := range sols {
			x[i], y[i] = sol.Flt[iFlt], sol.Flt[jFlt]
		}
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	if emptyMarker {
		args += io.Sf(",markeredgecolor='%s',markerfacecolor='none'", fmt.C)
	}
	plt.Plot(x, y, args)
}

// PlotAddFltOva adds flt-ova points to existent plot
func (o *Optimiser) PlotAddFltOva(iFlt, iOva int, sols []*Solution, ovaMult float64, fmt plt.Fmt, emptyMarker bool) {
	nsol := len(sols)
	x, y := make([]float64, nsol), make([]float64, nsol)
	for i, sol := range sols {
		x[i], y[i] = sol.Flt[iFlt], sol.Ova[iOva]*ovaMult
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	if emptyMarker {
		args += io.Sf(",markeredgecolor='%s',markerfacecolor='none'", fmt.C)
	}
	plt.Plot(x, y, args)
}

// PlotAddOvaOva adds ova-ova points to existent plot
func (o *Optimiser) PlotAddOvaOva(iOva, jOva int, sols []*Solution, feasibleOnly bool, fmt *plt.Fmt) {
	var x, y []float64
	for _, sol := range sols {
		if sol.Feasible() || !feasibleOnly {
			x = append(x, sol.Ova[iOva])
			y = append(y, sol.Ova[jOva])
		}
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	plt.Plot(x, y, args)
}

// PlotAddParetoFront highlights Pareto front
func (o *Optimiser) PlotAddParetoFront(iOva, jOva int, sols []*Solution, feasibleOnly bool, fmt *plt.Fmt) {
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	x, y, _ := GetParetoFront(iOva, jOva, sols, feasibleOnly)
	plt.Plot(x, y, args)
}

// PlotFltOva plots flt-ova points
func PlotFltOva(fnkey string, opt *Optimiser, sols0 []*Solution, iFlt, iOva, np int, ovaMult float64, fcn func(x float64) float64, extra func(), equalAxes bool) {
	if fcn != nil {
		X := utl.LinSpace(0, 1, np)
		Y := make([]float64, np)
		for i := 0; i < np; i++ {
			Y[i] = fcn(X[i])
		}
		plt.Plot(X, Y, "'b-'")
	}
	if sols0 != nil {
		opt.PlotAddFltOva(iFlt, iOva, sols0, ovaMult, plt.Fmt{L: "initial", M: "o", C: "k", Ls: "none", Ms: 3}, false)
	}
	SortByOva(opt.Solutions, iOva)
	best := opt.Solutions[0]
	opt.PlotAddFltOva(iFlt, iOva, opt.Solutions, ovaMult, plt.Fmt{L: "final", M: "o", C: "r", Ls: "none", Ms: 6}, true)
	plt.PlotOne(best.Flt[iFlt], best.Ova[iOva]*ovaMult, "'g*', markeredgecolor='g', label='best', clip_on=0, zorder=20")
	if extra != nil {
		extra()
	}
	if equalAxes {
		plt.Equal()
	}
	plt.Gll(io.Sf("$x_%d$", iFlt), io.Sf("$f_%d$", iOva), "leg_out=1, leg_ncol=4, leg_hlen=1.5")
	plt.SaveD("/tmp/goga", fnkey+".eps")
}

// PlotFltFlt plots flt-flt contour
func PlotFltFltContour(fnkey string, opt *Optimiser, sols0 []*Solution, iFlt, jFlt, iOva int, denormalise bool, cprms ContourParams, extra func(), equalAxes bool) {
	clr1 := "green"
	clr2 := "magenta"
	if cprms.Csimple {
		clr1 = "red"
		clr2 = "green"
	}
	opt.PlotContour(iFlt, jFlt, iOva, cprms)
	if sols0 != nil {
		opt.PlotAddFltFlt(iFlt, jFlt, sols0, denormalise, plt.Fmt{L: "initial", M: "o", C: "k", Ls: "none", Ms: 3}, false)
	}
	SortByOva(opt.Solutions, iOva)
	best := opt.Solutions[0]
	opt.PlotAddFltFlt(iFlt, jFlt, opt.Solutions, denormalise, plt.Fmt{L: "final", M: "o", C: clr2, Ls: "none", Ms: 7}, true)
	xs, ys := best.Flt[iFlt], best.Flt[jFlt]
	if denormalise {
		xs = opt.Xmin[iFlt] + best.Flt[iFlt]*opt.Dx[iFlt]
		ys = opt.Xmin[jFlt] + best.Flt[jFlt]*opt.Dx[jFlt]
	}
	plt.PlotOne(xs, ys, io.Sf("'k*', markersize=6, color='%s', markeredgecolor='%s', label='best', clip_on=0, zorder=20", clr1, clr1))
	if extra != nil {
		extra()
	}
	if equalAxes {
		plt.Equal()
	}
	plt.Gll(io.Sf("$x_%d$", iFlt), io.Sf("$x_%d$", jFlt), "leg_out=1, leg_ncol=4, leg_hlen=1.5")
	plt.SaveD("/tmp/goga", fnkey+".eps")
}

// PlotOvaOvaPareto plots ova-ova Pareto values
//  fmtAll   -- format for all points; use nil if not requested
//  fmtFront -- format for Pareto front; use nil if not requested
func PlotOvaOvaPareto(opt *Optimiser, sols0 []*Solution, iOva, jOva int, feasibleOnly bool, fmtAll, fmtFront *plt.Fmt) {
	if sols0 != nil {
		opt.PlotAddOvaOva(iOva, jOva, sols0, feasibleOnly, &plt.Fmt{L: "initial", M: "+", C: "g", Ls: "none", Ms: 4})
	}
	if fmtAll != nil {
		opt.PlotAddOvaOva(iOva, jOva, opt.Solutions, feasibleOnly, fmtAll)
	}
	if fmtFront != nil {
		opt.PlotAddParetoFront(iOva, jOva, opt.Solutions, feasibleOnly, fmtFront)
	}
	plt.Gll(io.Sf("$f_%d$", iOva), io.Sf("$f_%d$", jOva), "leg_out=1, leg_ncol=4, leg_hlen=1.5")
}

// PlotStar plots star with normalised OVAs
func PlotStar(opt *Optimiser) {
	nf := opt.Nf
	dθ := 2.0 * math.Pi / float64(nf)
	θ0 := 0.0
	if nf == 3 {
		θ0 = -math.Pi / 6.0
	}
	for _, ρ := range []float64{0.25, 0.5, 0.75, 1.0} {
		plt.Circle(0, 0, ρ, "ec='gray',lw=0.5,zorder=5")
	}
	arrowM, textM := 1.1, 1.15
	for i := 0; i < nf; i++ {
		θ := θ0 + float64(i)*dθ
		xi, yi := 0.0, 0.0
		xf, yf := arrowM*math.Cos(θ), arrowM*math.Sin(θ)
		plt.Arrow(xi, yi, xf, yf, "sc=10,st='->',lw=0.7,zorder=10,clip_on=0")
		plt.PlotOne(xf, yf, "'k+', ms=0")
		xf, yf = textM*math.Cos(θ), textM*math.Sin(θ)
		plt.Text(xf, yf, io.Sf("%d", i), "size=6,zorder=10,clip_on=0")
	}
	X, Y := make([]float64, nf+1), make([]float64, nf+1)
	clr := false
	neg := false
	step := 1
	count := 0
	colors := []string{"m", "orange", "g", "r", "b", "k"}
	var ρ float64
	for i, sol := range opt.Solutions {
		if sol.Feasible() && sol.FrontId == 0 && i%step == 0 {
			for j := 0; j < nf; j++ {
				if neg {
					ρ = 1.0 - sol.Ova[j]/(opt.RptFmax[j]-opt.RptFmin[j])
				} else {
					ρ = sol.Ova[j] / (opt.RptFmax[j] - opt.RptFmin[j])
				}
				θ := θ0 + float64(j)*dθ
				X[j], Y[j] = ρ*math.Cos(θ), ρ*math.Sin(θ)
			}
			X[nf], Y[nf] = X[0], Y[0]
			if clr {
				j := count % len(colors)
				plt.Plot(X, Y, io.Sf("'k-',color='%s',markersize=3,clip_on=0", colors[j]))
			} else {
				plt.Plot(X, Y, "'r-',marker='.',markersize=3,clip_on=0")
			}
			count++
		}
	}
	plt.Equal()
	plt.AxisOff()
}

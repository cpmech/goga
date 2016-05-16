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

// PlotFunctionsContour plots contour of f, g, and h functions
func (o *Optimiser) PlotFcnContour(iFlt, jFlt, iOva int, prms *ContourParams) {

	// fix parameters
	if prms == nil {
		prms = new(ContourParams)
	}
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

// PlotAddXvsX adds X versus X points to existent plot
func (o *Optimiser) PlotAddXvsX(iFlt, jFlt int, sols []*Solution, denormalise bool, fmt plt.Fmt, emptyMarker bool) {
	nsol := len(sols)
	X, Y := make([]float64, nsol), make([]float64, nsol)
	if denormalise {
		for i, sol := range sols {
			X[i] = o.Xmin[iFlt] + sol.Flt[iFlt]*o.Dx[iFlt]
			Y[i] = o.Xmin[jFlt] + sol.Flt[jFlt]*o.Dx[jFlt]
		}
	} else {
		for i, sol := range sols {
			X[i] = sol.Flt[iFlt]
			Y[i] = sol.Flt[jFlt]
		}
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	if emptyMarker {
		args += io.Sf(",markeredgecolor='%s',markerfacecolor='none'", fmt.C)
	}
	plt.Plot(X, Y, args)
}

// PlotAddXvsF adds X versus F points to existent plot
func (o *Optimiser) PlotAddXvsF(iFlt, iOva int, sols []*Solution, denormalise bool, ovaMult float64, fmt plt.Fmt, emptyMarker bool) {
	nsol := len(sols)
	X, Y := make([]float64, nsol), make([]float64, nsol)
	if denormalise {
		for i, sol := range sols {
			X[i] = o.Xmin[iFlt] + sol.Flt[iFlt]*o.Dx[iFlt]
			Y[i] = sol.Ova[iOva] * ovaMult
		}
	} else {
		for i, sol := range sols {
			X[i] = sol.Flt[iFlt]
			Y[i] = sol.Ova[iOva] * ovaMult
		}
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	if emptyMarker {
		args += io.Sf(",markeredgecolor='%s',markerfacecolor='none'", fmt.C)
	}
	plt.Plot(X, Y, args)
}

// PlotAddFvsF adds F versus F points to existent plot
func (o *Optimiser) PlotAddFvsF(iOva, jOva int, sols []*Solution, feasibleOnly bool, fmt *plt.Fmt) {
	var X, Y []float64
	for _, sol := range sols {
		if sol.Feasible() || !feasibleOnly {
			X = append(X, sol.Ova[iOva])
			Y = append(Y, sol.Ova[jOva])
		}
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	plt.Plot(X, Y, args)
}

// PlotAddFront add points on Pareto front
func (o *Optimiser) PlotAddFront(iOva, jOva int, sols []*Solution, feasibleOnly bool, fmt *plt.Fmt) {
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	X, Y, _ := GetParetoFront(iOva, jOva, sols, feasibleOnly)
	plt.Plot(X, Y, args)
}

// PlotXvsF plots X versus F points
//  If fcn is not nil, a curve with nptsFcn is added
func (o *Optimiser) PlotXvsF(sols0 []*Solution, iFlt, iOva int, denormalise bool, ovaMult float64, fcn func(x float64) float64, nptsFcn int) {
	if fcn != nil {
		var X []float64
		if denormalise {
			X = utl.LinSpace(o.Xmin[iFlt], o.Xmax[iFlt], nptsFcn)
		} else {
			X = utl.LinSpace(0, 1, nptsFcn)
		}
		Y := make([]float64, nptsFcn)
		for i := 0; i < nptsFcn; i++ {
			Y[i] = fcn(X[i])
		}
		plt.Plot(X, Y, "'b-'")
	}
	if sols0 != nil {
		o.PlotAddXvsF(iFlt, iOva, sols0, denormalise, ovaMult, plt.Fmt{L: "initial", M: "o", C: "k", Ls: "none", Ms: 3}, false)
	}
	o.PlotAddXvsF(iFlt, iOva, o.Solutions, denormalise, ovaMult, plt.Fmt{L: "final", M: "o", C: "r", Ls: "none", Ms: 6}, true)
	SortByOva(o.Solutions, iOva)
	best := o.Solutions[0]
	xs := best.Flt[iFlt]
	ys := best.Ova[iOva] * ovaMult
	xl := io.Sf("$\\bar{x}_%d$", iFlt)
	if denormalise {
		xl = io.Sf("$x_%d$", iFlt)
		xs = o.Xmin[iFlt] + best.Flt[iFlt]*o.Dx[iFlt]
	}
	plt.PlotOne(xs, ys, "'g*', markeredgecolor='g', label='best', clip_on=0, zorder=20")
	plt.Gll(xl, io.Sf("$f_%d$", iOva), "leg_out=1, leg_ncol=4, leg_hlen=1.5")
}

// PlotContour plots contour
func (o *Optimiser) PlotContour(sols0 []*Solution, iFlt, jFlt, iOva int, denormalise bool, cprms *ContourParams) {
	clr1 := "green"
	clr2 := "magenta"
	if cprms == nil {
		cprms = new(ContourParams)
	}
	if cprms.Csimple {
		clr1 = "red"
		clr2 = "green"
	}
	o.PlotFcnContour(iFlt, jFlt, iOva, cprms)
	if sols0 != nil {
		o.PlotAddXvsX(iFlt, jFlt, sols0, denormalise, plt.Fmt{L: "initial", M: "o", C: "k", Ls: "none", Ms: 3}, false)
	}
	o.PlotAddXvsX(iFlt, jFlt, o.Solutions, denormalise, plt.Fmt{L: "final", M: "o", C: clr2, Ls: "none", Ms: 7}, true)
	SortByOva(o.Solutions, iOva)
	best := o.Solutions[0]
	xs := best.Flt[iFlt]
	ys := best.Flt[jFlt]
	xl := io.Sf("$\\bar{x}_%d$", iFlt)
	yl := io.Sf("$\\bar{x}_%d$", jFlt)
	if denormalise {
		xs = o.Xmin[iFlt] + best.Flt[iFlt]*o.Dx[iFlt]
		ys = o.Xmin[jFlt] + best.Flt[jFlt]*o.Dx[jFlt]
		xl = io.Sf("$x_%d$", iFlt)
		yl = io.Sf("$x_%d$", jFlt)
	}
	plt.PlotOne(xs, ys, io.Sf("'k*', markersize=6, color='%s', markeredgecolor='%s', label='best', clip_on=0, zorder=20", clr1, clr1))
	plt.Gll(xl, yl, "leg_out=1, leg_ncol=4, leg_hlen=1.5")
}

// PlotFront plots Pareto front (F versus F values)
//  fmtAll   -- format for all points; use nil if not requested
//  fmtFront -- format for Pareto front; use nil if not requested
func (o *Optimiser) PlotFront(sols0 []*Solution, iOva, jOva int, feasibleOnly bool, fmtAll, fmtFront *plt.Fmt) {
	if sols0 != nil {
		o.PlotAddFvsF(iOva, jOva, sols0, feasibleOnly, &plt.Fmt{L: "initial", M: "+", C: "g", Ls: "none", Ms: 4})
	}
	if fmtAll != nil {
		o.PlotAddFvsF(iOva, jOva, o.Solutions, feasibleOnly, fmtAll)
	}
	if fmtFront != nil {
		o.PlotAddFront(iOva, jOva, o.Solutions, feasibleOnly, fmtFront)
	}
	plt.Gll(io.Sf("$f_%d$", iOva), io.Sf("$f_%d$", jOva), "leg_out=1, leg_ncol=4, leg_hlen=1.5")
}

// PlotStar plots star with normalised OVAs
func (o *Optimiser) PlotStar() {
	nf := o.Nf
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
	for i, sol := range o.Solutions {
		if sol.Feasible() && sol.FrontId == 0 && i%step == 0 {
			for j := 0; j < nf; j++ {
				if neg {
					ρ = 1.0 - sol.Ova[j]/(o.RptFmax[j]-o.RptFmin[j])
				} else {
					ρ = sol.Ova[j] / (o.RptFmax[j] - o.RptFmin[j])
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

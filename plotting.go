// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

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

// Plot plots contour
func (o *Optimiser) PlotContour(idxF, iFlt, jFlt int, prms ContourParams) {

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
	xmin, xmax := o.FltMin[iFlt], o.FltMax[iFlt]
	ymin, ymax := o.FltMin[jFlt], o.FltMax[jFlt]
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
			Zf[i][j] = o.F[grp][idxF]
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

func (o *Optimiser) PlotSolutions(iFlt, jFlt int, fmt plt.Fmt, emptyMarker bool) {
	x, y := make([]float64, o.NsolTot), make([]float64, o.NsolTot)
	for i, sol := range o.Solutions {
		x[i], y[i] = sol.Flt[iFlt], sol.Flt[jFlt]
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	if emptyMarker {
		args += io.Sf(",markeredgecolor='%s',markerfacecolor='none'", fmt.C)
	}
	plt.Plot(x, y, args)
}

func (o *Optimiser) PlotSolOvas(iFlt, jOva int, ovaMult float64, fmt plt.Fmt, emptyMarker bool) {
	x, y := make([]float64, o.NsolTot), make([]float64, o.NsolTot)
	for i, sol := range o.Solutions {
		x[i], y[i] = sol.Flt[iFlt], sol.Ova[jOva]*ovaMult
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	if emptyMarker {
		args += io.Sf(",markeredgecolor='%s',markerfacecolor='none'", fmt.C)
	}
	plt.Plot(x, y, args)
}

func (o *Optimiser) PlotOvas(iOva, jOva int, fmt plt.Fmt, emptyMarker bool) {
	x, y := make([]float64, o.NsolTot), make([]float64, o.NsolTot)
	for i, sol := range o.Solutions {
		x[i], y[i] = sol.Ova[iOva], sol.Ova[jOva]
	}
	args := fmt.GetArgs("") + ",clip_on=0,zorder=10"
	if emptyMarker {
		args += io.Sf(",markeredgecolor='%s',markerfacecolor='none'", fmt.C)
	}
	plt.Plot(x, y, args)
}

func (o *Optimiser) PlotParetoFront(iOva, jOva int, fmt plt.Fmt, emptyMarker bool) {
}

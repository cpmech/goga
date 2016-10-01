// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// PlotParams holds parameters to customize plots
type PlotParams struct {
	DirOut       string    // output directory; default = "/tmp/goga"
	FnKey        string    // filename key
	FnExt        string    // filename extension; default = ".eps" IMPORTANT: "." is required
	FmtSols0     plt.Fmt   // format for points indicating initial solutions
	FmtSols      plt.Fmt   // format for points indicating final solutions
	FmtBest      plt.Fmt   // format for points indicating best solution
	FmtFront     plt.Fmt   // format for points on Pareto front
	YfuncX       YfuncX_t  // y(x) function to plot from FltMin[iFlt] to FltMax[iFlt]
	FmtYfX       plt.Fmt   // format for y(x) function
	NptsYfX      int       // number of points for y(x) function
	Extra        func()    // extra plotting commands
	AxEqual      bool      // make axes equal
	Xlabel       string    // xlabel. "" means use default
	Ylabel       string    // xlabel. "" means use default
	LegPrms      string    // legend parameters
	FeasibleOnly bool      // plot feasible solutions only
	WithAll      bool      // with all points
	NoFront      bool      // do not show Pareto front
	Npts         int       // number of points for contour
	CmapIdx      int       // colormap index
	Cbar         bool      // with color bar
	Simple       bool      // simple contour
	FmtF         plt.Fmt   // format for f(x) function
	FmtG         plt.Fmt   // format for g(x) function
	FmtH         plt.Fmt   // format for h(x) function
	FmtA         plt.Fmt   // format for auxiliary field
	Xrange       []float64 // to override x-range
	Yrange       []float64 // to override y-range
	IdxH         int       // index of h function to plot. -1 means all
	Refx         []float64 // reference x vector with the other values in case Nflt > 2
	NoF          bool      // without f(x)
	NoG          bool      // without g(x)
	NoH          bool      // without h(x)
	WithAux      bool      // plot Solution.Aux (with the same colors as g(x))
}

// NewPlotParams allocates and sets default PlotParams
func NewPlotParams(simple bool) (o *PlotParams) {

	o = new(PlotParams)
	o.DirOut = "/tmp/goga"
	o.FnKey = "plt-goga"
	o.FnExt = ".eps"

	o.FmtSols0.L = "initial"
	o.FmtSols0.M = "o"
	o.FmtSols0.C = "k"
	o.FmtSols0.Ls = "none"
	o.FmtSols0.Ms = 3

	o.FmtSols.L = "final"
	o.FmtSols.M = "o"
	o.FmtSols.C = "magenta"
	o.FmtSols.Ls = "none"
	o.FmtSols.Ms = 7
	o.FmtSols.Void = true

	o.FmtBest.L = "best"
	o.FmtBest.M = "*"
	o.FmtBest.C = "#00b30d"
	o.FmtBest.Mec = "white"
	o.FmtBest.Ms = 6
	o.FmtBest.Mew = 0.3
	o.FmtBest.Z = 20

	o.FmtFront.L = "front"
	o.FmtFront.M = "*"
	o.FmtFront.Ls = "none"
	o.FmtFront.C = "red"
	o.FmtFront.Mec = "black"
	o.FmtFront.Ms = 6
	o.FmtFront.Mew = 0.3
	o.FmtFront.Z = 20

	o.NptsYfX = 41
	o.FmtYfX.L = "y(x)"
	o.FmtYfX.C = "blue"
	o.FmtYfX.Ls = "--"

	if simple {
		o.FmtSols.C = "#00b30d"
		o.FmtBest.C = "red"
		o.FmtBest.Mec = "black"
	}

	o.LegPrms = "leg_out=1, leg_ncol=4, leg_hlen=1.5"

	o.Npts = 41
	o.FmtF.C = "black"
	o.FmtF.Lw = 1
	o.FmtG.C = "yellow"
	o.FmtG.Lw = 1.5
	o.FmtH.C = "yellow"
	o.FmtH.Lw = 1.5
	o.FmtA.C = "yellow"
	o.FmtA.Lw = 1.5
	o.Simple = simple
	return
}

// PlotContour plots contour
func (o *Optimiser) PlotContour(iFlt, jFlt, iOva int, pp *PlotParams) {

	// check
	var x []float64
	if pp.Refx == nil {
		if iFlt > 1 || jFlt > 1 {
			chk.Panic("Refx vector must be given to PlotContour when iFlt or jFlt > 1")
		}
		x = make([]float64, 2)
	} else {
		x = make([]float64, len(pp.Refx))
		copy(x, pp.Refx)
	}

	// limits and meshgrid
	xmin, xmax := o.FltMin[iFlt], o.FltMax[iFlt]
	ymin, ymax := o.FltMin[jFlt], o.FltMax[jFlt]
	if pp.Xrange != nil {
		xmin, xmax = pp.Xrange[0], pp.Xrange[1]
	}
	if pp.Yrange != nil {
		ymin, ymax = pp.Yrange[0], pp.Yrange[1]
	}

	// check objective function
	var sol *Solution // copy of solution for objective function
	if o.MinProb == nil {
		pp.NoG = true
		pp.NoH = true
		sol = NewSolution(o.Nsol, 0, &o.Parameters)
		o.Solutions[0].CopyInto(sol)
		if pp.Refx != nil {
			copy(sol.Flt, pp.Refx)
		}
	}

	// auxiliary variables
	X, Y := utl.MeshGrid2D(xmin, xmax, ymin, ymax, pp.Npts, pp.Npts)
	var Zf [][]float64
	var Zg [][][]float64
	var Zh [][][]float64
	var Za [][]float64
	if !pp.NoF {
		Zf = utl.DblsAlloc(pp.Npts, pp.Npts)
	}
	if o.Ng > 0 && !pp.NoG {
		Zg = utl.Deep3alloc(o.Ng, pp.Npts, pp.Npts)
	}
	if o.Nh > 0 && !pp.NoH {
		Zh = utl.Deep3alloc(o.Nh, pp.Npts, pp.Npts)
	}
	if pp.WithAux {
		Za = utl.DblsAlloc(pp.Npts, pp.Npts)
	}

	// compute values
	grp := 0
	for i := 0; i < pp.Npts; i++ {
		for j := 0; j < pp.Npts; j++ {
			x[iFlt], x[jFlt] = X[i][j], Y[i][j]
			if o.MinProb == nil {
				copy(sol.Flt, x)
				o.ObjFunc(sol, grp)
				if !pp.NoF {
					Zf[i][j] = sol.Ova[iOva]
				}
				if pp.WithAux {
					Za[i][j] = sol.Aux
				}
			} else {
				o.MinProb(o.F[grp], o.G[grp], o.H[grp], x, nil, grp)
				if !pp.NoF {
					Zf[i][j] = o.F[grp][iOva]
				}
				if !pp.NoG {
					for k, g := range o.G[grp] {
						Zg[k][i][j] = g
					}
				}
				if !pp.NoH {
					for k, h := range o.H[grp] {
						Zh[k][i][j] = h
					}
				}
			}
		}
	}

	// plot f
	if !pp.NoF {
		txt := "cbar=0"
		if pp.Cbar {
			txt = ""
		}
		if pp.Simple {
			plt.ContourSimple(X, Y, Zf, true, 7, io.Sf("colors=['%s'], fsz=7, %s", pp.FmtF.C, txt))
		} else {
			plt.Contour(X, Y, Zf, io.Sf("fsz=7, cmapidx=%d, %s", pp.CmapIdx, txt))
		}
	}

	// plot g
	if !pp.NoG {
		for _, g := range Zg {
			plt.ContourSimple(X, Y, g, false, 7, io.Sf("zorder=5, levels=[0], colors=['%s'], linewidths=[%g], clip_on=0", pp.FmtG.C, pp.FmtG.Lw))
		}
	}

	// plot h
	if !pp.NoH {
		for i, h := range Zh {
			if i == pp.IdxH || pp.IdxH < 0 {
				plt.ContourSimple(X, Y, h, false, 7, io.Sf("zorder=5, levels=[0], colors=['%s'], linewidths=[%g], clip_on=0", pp.FmtH.C, pp.FmtH.Lw))
			}
		}
	}

	// plot aux
	if pp.WithAux {
		plt.ContourSimple(X, Y, Za, false, 7, io.Sf("zorder=5, levels=[0], colors=['%s'], linewidths=[%g], clip_on=0", pp.FmtA.C, pp.FmtA.Lw))
	}
}

// PlotAddFltFlt adds flt-flt points to existent plot
func (o *Optimiser) PlotAddFltFlt(iFlt, jFlt int, sols []*Solution, fmt *plt.Fmt) {
	nsol := len(sols)
	x, y := make([]float64, nsol), make([]float64, nsol)
	for i, sol := range sols {
		x[i], y[i] = sol.Flt[iFlt], sol.Flt[jFlt]
	}
	plt.Plot(x, y, fmt.GetArgs(""))
}

// PlotAddFltOva adds flt-ova points to existent plot
func (o *Optimiser) PlotAddFltOva(iFlt, iOva int, sols []*Solution, ovaMult float64, fmt *plt.Fmt) {
	nsol := len(sols)
	x, y := make([]float64, nsol), make([]float64, nsol)
	for i, sol := range sols {
		x[i], y[i] = sol.Flt[iFlt], sol.Ova[iOva]*ovaMult
	}
	plt.Plot(x, y, fmt.GetArgs(""))
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
	plt.Plot(x, y, fmt.GetArgs(""))
}

// PlotAddParetoFront highlights Pareto front
func (o *Optimiser) PlotAddParetoFront(iOva, jOva int, sols []*Solution, feasibleOnly bool, fmt *plt.Fmt) {
	x, y, _ := GetParetoFront(iOva, jOva, sols, feasibleOnly)
	plt.Plot(x, y, fmt.GetArgs(""))
}

// PlotFltOva plots flt-ova points
func (o *Optimiser) PlotFltOva(sols0 []*Solution, iFlt, iOva int, ovaMult float64, pp *PlotParams) {
	if pp.YfuncX != nil {
		X := utl.LinSpace(o.FltMin[iFlt], o.FltMax[iFlt], pp.NptsYfX)
		Y := make([]float64, pp.NptsYfX)
		for i := 0; i < pp.NptsYfX; i++ {
			Y[i] = pp.YfuncX(X[i])
		}
		plt.Plot(X, Y, pp.FmtYfX.GetArgs(""))
	}
	if sols0 != nil {
		o.PlotAddFltOva(iFlt, iOva, sols0, ovaMult, &pp.FmtSols0)
	}
	o.PlotAddFltOva(iFlt, iOva, o.Solutions, ovaMult, &pp.FmtSols)
	best, _ := GetBestFeasible(o, iOva)
	if best != nil {
		plt.PlotOne(best.Flt[iFlt], best.Ova[iOva]*ovaMult, pp.FmtBest.GetArgs(""))
	}
	if pp.Extra != nil {
		pp.Extra()
	}
	if pp.AxEqual {
		plt.Equal()
	}
	plt.Gll(io.Sf("$x_{%d}$", iFlt), io.Sf("$f_{%d}$", iOva), "leg_out=1, leg_ncol=4, leg_hlen=1.5")
	plt.SaveD(pp.DirOut, pp.FnKey+pp.FnExt)
}

// PlotFltFlt plots flt-flt contour
// use iFlt==-1 || jFlt==-1 to plot all combinations
func (o *Optimiser) PlotFltFltContour(sols0 []*Solution, iFlt, jFlt, iOva int, pp *PlotParams) {
	best, _ := GetBestFeasible(o, iOva)
	plotAll := iFlt < 0 || jFlt < 0
	plotCommands := func(i, j int) {
		o.PlotContour(i, j, iOva, pp)
		if sols0 != nil {
			o.PlotAddFltFlt(i, j, sols0, &pp.FmtSols0)
		}
		o.PlotAddFltFlt(i, j, o.Solutions, &pp.FmtSols)
		if best != nil {
			plt.PlotOne(best.Flt[i], best.Flt[j], pp.FmtBest.GetArgs(""))
		}
		if pp.Extra != nil {
			pp.Extra()
		}
		if pp.AxEqual {
			plt.Equal()
		}
	}
	if plotAll {
		idx := 1
		ncol := o.Nflt - 1
		for row := 0; row < o.Nflt; row++ {
			idx += row
			for col := row + 1; col < o.Nflt; col++ {
				plt.Subplot(ncol, ncol, idx)
				plt.SplotGap(0.0, 0.0)
				plotCommands(col, row)
				if col > row+1 {
					plt.SetXnticks(0)
					plt.SetYnticks(0)
				} else {
					plt.Gll(io.Sf("$x_{%d}$", col), io.Sf("$x_{%d}$", row), "leg=0")
				}
				idx++
			}
		}
		idx = ncol*(ncol-1) + 1
		plt.Subplot(ncol, ncol, idx)
		plt.AxisOff()
		// TODO: fix formatting of open marker, add star to legend
		plt.DrawLegend([]plt.Fmt{pp.FmtSols0, pp.FmtSols, pp.FmtBest}, 8, "center", false, "")
	} else {
		plotCommands(iFlt, jFlt)
		if pp.Xlabel == "" {
			plt.Gll(io.Sf("$x_{%d}$", iFlt), io.Sf("$x_{%d}$", jFlt), pp.LegPrms)
		} else {
			plt.Gll(pp.Xlabel, pp.Ylabel, pp.LegPrms)
		}
	}
	plt.SaveD(pp.DirOut, pp.FnKey+pp.FnExt)
}

// PlotOvaOvaPareto plots ova-ova Pareto values
func (o *Optimiser) PlotOvaOvaPareto(sols0 []*Solution, iOva, jOva int, pp *PlotParams) {
	if sols0 != nil {
		o.PlotAddOvaOva(iOva, jOva, sols0, pp.FeasibleOnly, &pp.FmtSols0)
	}
	if pp.WithAll {
		o.PlotAddOvaOva(iOva, jOva, o.Solutions, pp.FeasibleOnly, &pp.FmtSols)
	}
	if !pp.NoFront {
		o.PlotAddParetoFront(iOva, jOva, o.Solutions, pp.FeasibleOnly, &pp.FmtFront)
	}
	if pp.Extra != nil {
		pp.Extra()
	}
	plt.Gll(io.Sf("$f_{%d}$", iOva), io.Sf("$f_{%d}$", jOva), pp.LegPrms)
	plt.SaveD(pp.DirOut, pp.FnKey+pp.FnExt)
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

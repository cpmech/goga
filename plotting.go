// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// PlotParams holds parameters to customize plots
type PlotParams struct {

	// output directory and filename
	DirOut string // output directory; default = "/tmp/goga"
	FnKey  string // filename key

	// auxiliary
	YfuncX  YfuncX_t // y(x) function to plot from FltMin[iFlt] to FltMax[iFlt]
	NptsYfX int      // number of points for y(x) function
	Extra   func()   // extra plotting commands

	// options
	NoF          bool      // without f(x)
	NoG          bool      // without g(x)
	NoH          bool      // without h(x)
	WithAux      bool      // plot Solution.Aux (with the same colors as g(x))
	OnlyAux      bool      // plot only Solution.Aux
	Limits       bool      // plot limits of variables in contour
	FeasibleOnly bool      // plot feasible solutions only
	WithAll      bool      // with all points
	NoFront      bool      // do not show Pareto front
	Simple       bool      // simple contour
	Npts         int       // number of points for contour
	IdxH         int       // index of h function to plot. -1 means all
	RefX         []float64 // reference x vector with the other values in case Nflt > 2
	Xrange       []float64 // to override x-range
	Yrange       []float64 // to override y-range
	ContourAt    string    // if RefX==nil, options: "minimum", "maximum", "middle", "best" (default)
	Xlabel       string    // x-axis label
	Ylabel       string    // y-axis label

	// plot arguments
	AxEqual    bool   // plot: equal scales
	ArgsLeg    *plt.A // plot: arguments for legend
	ArgsF      *plt.A // plot: f(x) function
	ArgsG      *plt.A // plot: g(x) function
	ArgsH      *plt.A // plot: h(x) function
	ArgsAux    *plt.A // plot: format for auxiliary field
	ArgsSols0  *plt.A // plot: points indicating initial solutions
	ArgsSols   *plt.A // plot: points indicating final solutions
	ArgsBest   *plt.A // plot: points indicating best solution
	ArgsFront  *plt.A // plot: points on Pareto front
	ArgsYfX    *plt.A // plot: y(x) function
	ArgsSimple *plt.A // plot: simple contours
}

// NewPlotParams allocates and sets default PlotParams
func NewPlotParams(simple bool) (o *PlotParams) {

	// output directory and filename
	o = new(PlotParams)
	o.DirOut = "/tmp/goga"
	o.FnKey = "plt-goga"

	// auxiliary
	o.NptsYfX = 101
	o.Npts = 21

	// plot arguments
	o.ArgsLeg = &plt.A{CmapIdx: 0, LegOut: true, LegNcol: 4, LegHlen: 1.5}
	o.ArgsF = &plt.A{}
	o.ArgsG = &plt.A{Colors: []string{"y"}, Levels: []float64{0}, Lw: 2}
	o.ArgsH = &plt.A{Colors: []string{"y"}, Levels: []float64{0}, Lw: 2}
	o.ArgsAux = &plt.A{Colors: []string{"y"}, Levels: []float64{0}, Lw: 2}
	o.ArgsSols0 = &plt.A{C: "k", M: "o", Ms: 3, Ls: "none", L: "initial"}
	o.ArgsSols = &plt.A{C: "m", M: "o", Ms: 7, Ls: "none", L: "final", Void: true}
	o.ArgsBest = &plt.A{C: "#00b30d", M: "*", Ms: 6, Ls: "none", Mec: "white", Mew: 0.3, L: "best"}
	o.ArgsFront = &plt.A{C: "r", M: ".", Ms: 6, Ls: "none", Mec: "black", Mew: 0.3, L: "front"}
	o.ArgsYfX = &plt.A{C: "b", Ls: "--", L: "y(x)"}
	o.ArgsSimple = &plt.A{Colors: []string{"y"}, Levels: []float64{0}, Lw: 2}

	// options
	o.Simple = simple
	if simple {
		o.ArgsSols.C = "#00b30d"
		o.ArgsBest.C = "red"
		o.ArgsBest.Mec = "black"
	}
	return
}

// PlotContour plots contour with other components @ x=RefX
//  If RefX==nil, x can be either @ minimum, maximum, middle or best
//  For x @ best, Solutions will be sorted
//  Input:
//    pp -- plotting parameters [may be nil]
func (o *Optimiser) PlotContour(iFlt, jFlt, iOva int, pp *PlotParams) {

	// check
	var x []float64
	if pp == nil {
		pp = NewPlotParams(false)
	}
	if pp.RefX == nil {
		x = make([]float64, o.Nflt)
		switch pp.ContourAt {
		case "minimum":
			copy(x, o.FltMin)
		case "maximum":
			copy(x, o.FltMax)
		case "middle":
			for k := 0; k < o.Nflt; k++ {
				x[k] = (o.FltMin[k] + o.FltMax[k]) / 2.0
			}
		default:
			SortSolutions(o.Solutions, 0)
			copy(x, o.Solutions[0].Flt)
		}
	} else {
		x = make([]float64, len(pp.RefX))
		copy(x, pp.RefX)
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
		if pp.RefX != nil {
			copy(sol.Flt, pp.RefX)
		}
	}

	// auxiliary variables
	X, Y := utl.MeshGrid2d(xmin, xmax, ymin, ymax, pp.Npts, pp.Npts)
	var Zf [][]float64
	var Zg [][][]float64
	var Zh [][][]float64
	var Za [][]float64
	if !pp.NoF {
		Zf = utl.Alloc(pp.Npts, pp.Npts)
	}
	if o.Ng > 0 && !pp.NoG {
		Zg = utl.Deep3alloc(o.Ng, pp.Npts, pp.Npts)
	}
	if o.Nh > 0 && !pp.NoH {
		Zh = utl.Deep3alloc(o.Nh, pp.Npts, pp.Npts)
	}
	if pp.WithAux {
		Za = utl.Alloc(pp.Npts, pp.Npts)
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
	if !pp.NoF && !pp.OnlyAux {
		if pp.Simple {
			plt.ContourL(X, Y, Zf, pp.ArgsSimple)
		} else {
			plt.ContourF(X, Y, Zf, pp.ArgsF)
		}
	}

	// plot g
	if !pp.NoG && !pp.OnlyAux {
		for _, g := range Zg {
			plt.ContourL(X, Y, g, pp.ArgsG)
		}
	}

	// plot h
	if !pp.NoH && !pp.OnlyAux {
		for i, h := range Zh {
			if i == pp.IdxH || pp.IdxH < 0 {
				plt.ContourL(X, Y, h, pp.ArgsH)
			}
		}
	}

	// plot aux
	if pp.WithAux {
		plt.ContourL(X, Y, Za, pp.ArgsAux)
	}

	// limits
	if pp.Limits {
		plt.Plot(
			[]float64{o.FltMin[iFlt], o.FltMax[iFlt], o.FltMax[iFlt], o.FltMin[iFlt], o.FltMin[iFlt]},
			[]float64{o.FltMin[jFlt], o.FltMin[jFlt], o.FltMax[jFlt], o.FltMax[jFlt], o.FltMin[jFlt]},
			nil,
		)
	}
}

// PlotAddFltFlt adds flt-flt points to existent plot
func (o *Optimiser) PlotAddFltFlt(iFlt, jFlt int, sols []*Solution, args *plt.A) {
	nsol := len(sols)
	x, y := make([]float64, nsol), make([]float64, nsol)
	for i, sol := range sols {
		x[i], y[i] = sol.Flt[iFlt], sol.Flt[jFlt]
	}
	plt.Plot(x, y, args)
}

// PlotAddFltOva adds flt-ova points to existent plot
func (o *Optimiser) PlotAddFltOva(iFlt, iOva int, sols []*Solution, ovaMult float64, args *plt.A) {
	nsol := len(sols)
	x, y := make([]float64, nsol), make([]float64, nsol)
	for i, sol := range sols {
		x[i], y[i] = sol.Flt[iFlt], sol.Ova[iOva]*ovaMult
	}
	plt.Plot(x, y, args)
}

// PlotAddOvaOva adds ova-ova points to existent plot
func (o *Optimiser) PlotAddOvaOva(iOva, jOva int, sols []*Solution, feasibleOnly bool, args *plt.A) {
	var x, y []float64
	for _, sol := range sols {
		if sol.Feasible() || !feasibleOnly {
			x = append(x, sol.Ova[iOva])
			y = append(y, sol.Ova[jOva])
		}
	}
	plt.Plot(x, y, args)
}

// PlotAddParetoFront highlights Pareto front
func (o *Optimiser) PlotAddParetoFront(iOva, jOva int, sols []*Solution, feasibleOnly bool, args *plt.A) {
	x, y, _ := GetParetoFront(iOva, jOva, sols, feasibleOnly)
	plt.Plot(x, y, args)
}

// PlotFltOva plots flt-ova points
func (o *Optimiser) PlotFltOva(sols0 []*Solution, iFlt, iOva int, ovaMult float64, pp *PlotParams) {
	if pp.YfuncX != nil {
		X := utl.LinSpace(o.FltMin[iFlt], o.FltMax[iFlt], pp.NptsYfX)
		Y := make([]float64, pp.NptsYfX)
		for i := 0; i < pp.NptsYfX; i++ {
			Y[i] = pp.YfuncX(X[i])
		}
		plt.Plot(X, Y, pp.ArgsYfX)
	}
	if sols0 != nil {
		o.PlotAddFltOva(iFlt, iOva, sols0, ovaMult, pp.ArgsSols0)
	}
	o.PlotAddFltOva(iFlt, iOva, o.Solutions, ovaMult, pp.ArgsSols)
	best, _ := GetBestFeasible(o, iOva)
	if best != nil {
		plt.PlotOne(best.Flt[iFlt], best.Ova[iOva]*ovaMult, pp.ArgsBest)
	}
	if pp.Extra != nil {
		pp.Extra()
	}
	if pp.AxEqual {
		plt.Equal()
	}
	plt.Gll(io.Sf("$x_{%d}$", iFlt), io.Sf("$f_{%d}$", iOva), &plt.A{LegOut: true, LegNcol: 4, LegHlen: 1.5})
	plt.Save(pp.DirOut, pp.FnKey)
}

// PlotFltFlt plots flt-flt contour
// use iFlt==-1 || jFlt==-1 to plot all combinations
func (o *Optimiser) PlotFltFltContour(sols0 []*Solution, iFlt, jFlt, iOva int, pp *PlotParams) {
	best, _ := GetBestFeasible(o, iOva)
	plotAll := iFlt < 0 || jFlt < 0
	plotCommands := func(i, j int) {
		o.PlotContour(i, j, iOva, pp)
		if sols0 != nil {
			o.PlotAddFltFlt(i, j, sols0, pp.ArgsSols0)
		}
		o.PlotAddFltFlt(i, j, o.Solutions, pp.ArgsSols)
		if best != nil {
			plt.PlotOne(best.Flt[i], best.Flt[j], pp.ArgsBest)
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
					plt.Gll(io.Sf("$x_{%d}$", col), io.Sf("$x_{%d}$", row), nil)
				}
				idx++
			}
		}
		idx = ncol*(ncol-1) + 1
		plt.Subplot(ncol, ncol, idx)
		plt.AxisOff()
		// TODO: fix formatting of open marker, add star to legend
		plt.LegendX([]*plt.A{pp.ArgsSols0, pp.ArgsSols, pp.ArgsBest}, pp.ArgsLeg)
	} else {
		plotCommands(iFlt, jFlt)
		plt.Gll(io.Sf("$x_{%d}$", iFlt), io.Sf("$x_{%d}$", jFlt), pp.ArgsLeg)
	}
	plt.Save(pp.DirOut, pp.FnKey)
}

// PlotOvaOvaPareto plots ova-ova Pareto values
func (o *Optimiser) PlotOvaOvaPareto(sols0 []*Solution, iOva, jOva int, pp *PlotParams) {
	if sols0 != nil {
		o.PlotAddOvaOva(iOva, jOva, sols0, pp.FeasibleOnly, pp.ArgsSols0)
	}
	if pp.WithAll {
		o.PlotAddOvaOva(iOva, jOva, o.Solutions, pp.FeasibleOnly, pp.ArgsSols)
	}
	if !pp.NoFront {
		o.PlotAddParetoFront(iOva, jOva, o.Solutions, pp.FeasibleOnly, pp.ArgsFront)
	}
	if pp.Extra != nil {
		pp.Extra()
	}
	xl, yl := io.Sf("$f_{%d}$", iOva), io.Sf("$f_{%d}$", jOva)
	if pp.Xlabel != "" {
		xl = pp.Xlabel
	}
	if pp.Ylabel != "" {
		yl = pp.Ylabel
	}
	plt.Gll(xl, yl, pp.ArgsLeg)
	plt.Save(pp.DirOut, pp.FnKey)
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
		plt.Circle(0, 0, ρ, &plt.A{Ec: "grey", Lw: 0.5})
	}
	arrowM, textM := 1.1, 1.15
	for i := 0; i < nf; i++ {
		θ := θ0 + float64(i)*dθ
		xi, yi := 0.0, 0.0
		xf, yf := arrowM*math.Cos(θ), arrowM*math.Sin(θ)
		plt.Arrow(xi, yi, xf, yf, &plt.A{Scale: 10, Style: "->", Lw: 0.7})
		plt.PlotOne(xf, yf, &plt.A{C: "k", M: "+"})
		xf, yf = textM*math.Cos(θ), textM*math.Sin(θ)
		plt.Text(xf, yf, io.Sf("%d", i), nil)
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
				plt.Plot(X, Y, &plt.A{C: colors[j], Ms: 3})
			} else {
				plt.Plot(X, Y, &plt.A{C: "r", M: ".", Ms: 3})
			}
			count++
		}
	}
	plt.Equal()
	plt.AxisOff()
}

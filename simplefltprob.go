// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// SimpleFltFcn_t simple float problem function type
type SimpleFltFcn_t func(f, g, h []float64, x []float64)

// SimpleFltProb implements optimisation problems defined by:
//  min  {f0(x), f1(x), f2(x), ...}  nf functions
//       g0(x) ≥ 0
//       g1(x) ≥ 0
//  s.t. ....  ≥ 0  ng inequalities
//       h0(x) = 0
//       h1(x) = 0
//       ....  = 0  nh equalities
type SimpleFltProb struct {

	// data
	Fcn SimpleFltFcn_t // function evaluator
	F   []float64      // functions
	G   []float64      // inequalities
	H   []float64      // equalities
	C   *ConfParams    // configuration parameters
	Evo *Evolver       // evolver

	// auxiliary
	NumfmtX string // number format for x
	NumfmtF string // number format for f(x)

	// results and stat
	Xbest     [][]float64 // [nfeasible][nx] (max=ntrials) the best feasible floats
	Nfeasible int         // counter for feasible results

	// plotting
	PopsIni    []Population // [nisl] initial populations in all islands for the best trial
	PopsBest   []Population // [nisl] the best populations for the best trial
	PltDirout  string       // directory to save files
	PltIdxF    int          // index of which f[i] to plot
	PltNpts    int          // number of points for contour
	PltCmapIdx int          // colormap index
	PltCsimple bool         // simple contour
	PltAxEqual bool         // axes-equal
	PltLwg     float64      // linewidth for g functions
	PltLwh     float64      // linewidth for h functions
	PltArgs    string       // extra arguments for plot
	PltExtra   func()       // extra function
}

// Init initialises simple flot problem structure
func NewSimpleFltProb(fcn SimpleFltFcn_t, nf, ng, nh int, C *ConfParams) (o *SimpleFltProb) {

	// data
	o = new(SimpleFltProb)
	o.Fcn = fcn
	o.F = make([]float64, nf)
	o.G = make([]float64, ng)
	o.H = make([]float64, nh)
	o.C = C
	o.C.Nova = nf
	o.C.Noor = ng + nh

	// objective function
	o.C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		x := ind.GetFloats()
		o.Fcn(o.F, o.G, o.H, x)
		for i, f := range o.F {
			ind.Ovas[i] = f
		}
		for i, g := range o.G {
			ind.Oors[i] = utl.GtePenalty(g, 0.0, 1) // g[i] ≥ 0
		}
		ng := len(o.G)
		for i, h := range o.H {
			h = math.Abs(h)
			ind.Ovas[0] += h
			ind.Oors[ng+i] = utl.GtePenalty(o.C.Eps1, h, 1) // ϵ ≥ |h[i]|
		}
	}

	// evolver
	o.Evo = NewEvolver(o.C)

	// auxiliary
	o.NumfmtX = "%8.5f"
	o.NumfmtF = "%8.5f"

	// results and stat
	nx := len(o.C.RangeFlt)
	o.Xbest = utl.DblsAlloc(o.C.Ntrials, nx)

	// plotting
	if nx != 2 {
		o.C.DoPlot = false
	}
	if o.C.DoPlot {
		o.PopsIni = o.Evo.GetPopulations()
		o.PltDirout = "/tmp/goga"
		o.PltNpts = 41
		o.PltLwg = 1.5
		o.PltLwh = 1.5
	}
	return
}

// Run runs optimisations
func (o *SimpleFltProb) Run(verbose bool) {

	// run all trials
	for itrial := 0; itrial < o.C.Ntrials; itrial++ {

		// reset populations
		if itrial > 0 {
			for id, isl := range o.Evo.Islands {
				isl.Pop = o.C.PopFltGen(id, o.C)
				isl.CalcOvs(isl.Pop, 0)
				isl.CalcDemeritsAndSort(isl.Pop)
			}
		}

		// run evolution
		o.Evo.Run()

		// results
		xbest := o.Evo.Best.GetFloats()
		o.Fcn(o.F, o.G, o.H, xbest)

		// check if best is unfeasible
		unfeasible := false
		for _, g := range o.G {
			if g < 0 {
				unfeasible = true
				break
			}
		}
		for _, h := range o.H {
			if math.Abs(h) > o.C.Eps1 {
				unfeasible = true
				break
			}
		}

		// feasible results
		if !unfeasible {
			for i, x := range xbest {
				o.Xbest[o.Nfeasible][i] = x
			}
			o.Nfeasible++
		}

		// message
		if verbose {
			io.Pfyel("x*="+o.NumfmtX+" f="+o.NumfmtF, xbest, o.F)
			if unfeasible {
				io.Pfred(" unfeasible\n")
			} else {
				io.Pfgreen(" ok\n")
			}
		}

		// best populations
		if o.C.DoPlot {
			if o.Nfeasible == 1 {
				o.PopsBest = o.Evo.GetPopulations()
			} else {
				fcur := utl.DblCopy(o.F)
				o.Fcn(o.F, o.G, o.H, o.Xbest[o.Nfeasible-1])
				cur_dom, _ := utl.DblsParetoMin(fcur, o.F)
				if cur_dom {
					o.PopsBest = o.Evo.GetPopulations()
				}
			}
		}
	}
}

// Plot plots contour
func (o *SimpleFltProb) Plot(fnkey string) {

	// check
	if !o.C.DoPlot {
		return
	}

	// limits and meshgrid
	xmin, xmax := o.C.RangeFlt[0][0], o.C.RangeFlt[0][1]
	ymin, ymax := o.C.RangeFlt[1][0], o.C.RangeFlt[1][1]

	// auxiliary variables
	X, Y := utl.MeshGrid2D(xmin, xmax, ymin, ymax, o.PltNpts, o.PltNpts)
	Zf := utl.DblsAlloc(o.PltNpts, o.PltNpts)
	var Zg [][][]float64
	var Zh [][][]float64
	ng := len(o.G)
	nh := len(o.H)
	if ng > 0 {
		Zg = utl.Deep3alloc(ng, o.PltNpts, o.PltNpts)
	}
	if nh > 0 {
		Zh = utl.Deep3alloc(nh, o.PltNpts, o.PltNpts)
	}

	// compute values
	x := make([]float64, 2)
	for i := 0; i < o.PltNpts; i++ {
		for j := 0; j < o.PltNpts; j++ {
			x[0], x[1] = X[i][j], Y[i][j]
			o.Fcn(o.F, o.G, o.H, x)
			Zf[i][j] = o.F[o.PltIdxF]
			for k, g := range o.G {
				Zg[k][i][j] = g
			}
			for k, h := range o.H {
				Zh[k][i][j] = h
			}
		}
	}

	// prepare plot area
	plt.Reset()
	plt.SetForEps(0.8, 350)

	// plot f
	if o.PltArgs != "" {
		o.PltArgs = "," + o.PltArgs
	}
	if o.PltCsimple {
		plt.ContourSimple(X, Y, Zf, true, 7, "colors=['k'], fsz=7"+o.PltArgs)
	} else {
		plt.Contour(X, Y, Zf, io.Sf("fsz=7, cmapidx=%d"+o.PltArgs, o.PltCmapIdx))
	}

	// plot g
	clr := "yellow"
	if o.PltCsimple {
		clr = "blue"
	}
	for _, g := range Zg {
		plt.ContourSimple(X, Y, g, false, 7, io.Sf("zorder=5, levels=[0], colors=['%s'], linewidths=[%g], clip_on=0", clr, o.PltLwg))
	}

	// plot h
	clr = "yellow"
	if o.PltCsimple {
		clr = "blue"
	}
	for _, h := range Zh {
		plt.ContourSimple(X, Y, h, false, 7, io.Sf("zorder=5, levels=[0], colors=['%s'], linewidths=[%g], clip_on=0", clr, o.PltLwh))
	}

	// initial populations
	l := "initial population"
	for _, pop := range o.PopsIni {
		for _, ind := range pop {
			x := ind.GetFloats()
			plt.PlotOne(x[0], x[1], io.Sf("'k.', zorder=20, clip_on=0, label='%s'", l))
			l = ""
		}
	}

	// final populations
	l = "final population"
	for _, pop := range o.PopsBest {
		for _, ind := range pop {
			x := ind.GetFloats()
			plt.PlotOne(x[0], x[1], io.Sf("'ko', ms=6, zorder=30, clip_on=0, label='%s', markerfacecolor='none'", l))
			l = ""
		}
	}

	// extra
	if o.PltExtra != nil {
		o.PltExtra()
	}

	// best result
	if o.Nfeasible > 0 {
		x, _, _, _ := o.find_best()
		plt.PlotOne(x[0], x[1], "'m*', zorder=50, clip_on=0, label='best', markeredgecolor='m'")
	}

	// save figure
	plt.Cross("clr='grey'")
	if o.PltAxEqual {
		plt.Equal()
	}
	plt.AxisRange(xmin, xmax, ymin, ymax)
	plt.Gll("$x_0$", "$x_1$", "leg_out=1, leg_ncol=4, leg_hlen=1.5")
	plt.SaveD(o.PltDirout, fnkey+".eps")
}

func (o *SimpleFltProb) find_best() (x, f, g, h []float64) {
	chk.IntAssertLessThan(0, o.Nfeasible) // 0 < nfeasible
	nx := len(o.C.RangeFlt)
	nf := len(o.F)
	ng := len(o.G)
	nh := len(o.H)
	x = make([]float64, nx)
	f = make([]float64, nf)
	g = make([]float64, ng)
	h = make([]float64, nh)
	copy(x, o.Xbest[0])
	o.Fcn(f, g, h, x)
	for i := 1; i < o.Nfeasible; i++ {
		o.Fcn(o.F, o.G, o.H, o.Xbest[i])
		_, other_dom := utl.DblsParetoMin(f, o.F)
		if other_dom {
			copy(x, o.Xbest[i])
			copy(f, o.F)
			copy(g, o.G)
			copy(h, o.H)
		}
	}
	return
}

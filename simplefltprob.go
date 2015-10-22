// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// ParetoF1F0_t Pareto front solution
type ParetoF1F0_t func(f0 float64) float64

// SimpleFltFcn_t simple float problem function type
type SimpleFltFcn_t func(f, g, h, x []float64, isl int)

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
	C   *ConfParams    // configuration parameters

	// sandbox
	nf int         // number of functions
	ng int         // number of inequalities
	nh int         // number of equalities
	ff [][]float64 // [nisl][nf] functions
	gg [][]float64 // [nisl][ng] inequalities
	hh [][]float64 // [nisl][nh] equalities

	// evolver
	Evo *Evolver // evolver

	// auxiliary
	NumfmtX string // number format for x
	NumfmtF string // number format for f(x)
	ShowCts bool   // show g(x) and/or h(x) values if verbose

	// results and stat
	Xbest     [][]float64 // [nfeasible][nx] (max=ntrials) the best feasible floats
	Nfeasible int         // counter for feasible results

	// stat about Pareto front
	ParNdiv   int          // number of divisions of bins
	ParF1F0   ParetoF1F0_t // known solution
	ParFmin   []float64    // known solution: ova min
	ParFmax   []float64    // known solution: ova max
	ParRadM   []float64    // radius multiplier to find near bins
	ParNray   int          // number of rays to find near bins
	ParBins   gm.Bins      // bins close to Pareto front
	ParSelB   map[int]bool // selected bins
	ParDisErr []float64    // dist errors
	ParSpread []float64    // spreads

	// plotting
	PopsIni         []Population // [nisl] initial populations in all islands for the best trial
	PopsBest        []Population // [nisl] the best populations for the best trial
	PltDirout       string       // directory to save files
	PltIdxF         int          // index of which f[i] to plot
	PltNpts         int          // number of points for contour
	PltCmapIdx      int          // colormap index
	PltCsimple      bool         // simple contour
	PltAxEqual      bool         // axes-equal
	PltLwg          float64      // linewidth for g functions
	PltLwh          float64      // linewidth for h functions
	PltArgs         string       // extra arguments for plot
	PltExtra        func()       // extra function
	PltXrange       []float64    // to override x-range
	PltYrange       []float64    // to override y-range
	PltShowIniPop   bool         // show ini population
	PltShowFinalPop bool         // show final population
}

// Init initialises simple flot problem structure
func NewSimpleFltProb(fcn SimpleFltFcn_t, nf, ng, nh int, C *ConfParams) (o *SimpleFltProb) {

	// data
	o = new(SimpleFltProb)
	o.Fcn = fcn
	o.C = C
	o.C.Nova = nf
	o.C.Noor = ng + nh

	// sandbox
	o.nf, o.ng, o.nh = nf, ng, nh
	o.ff = utl.DblsAlloc(o.C.Nisl, o.nf)
	o.gg = utl.DblsAlloc(o.C.Nisl, o.ng)
	o.hh = utl.DblsAlloc(o.C.Nisl, o.nh)

	// objective function
	o.C.OvaOor = func(ind *Individual, isl, time int, report *bytes.Buffer) {
		x := ind.GetFloats()
		o.Fcn(o.ff[isl], o.gg[isl], o.hh[isl], x, isl)
		for i, f := range o.ff[isl] {
			ind.Ovas[i] = f
		}
		for i, g := range o.gg[isl] {
			ind.Oors[i] = utl.GtePenalty(g, 0.0, 1) // g[i] ≥ 0
		}
		for i, h := range o.hh[isl] {
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

	// Pareto front
	o.ParNdiv = 20
	o.ParRadM = []float64{0.02, 0.04}
	o.ParNray = 8

	// plotting
	if o.C.DoPlot {
		o.PopsIni = o.Evo.GetPopulations()
		o.PltDirout = "/tmp/goga"
		o.PltNpts = 41
		o.PltLwg = 1.5
		o.PltLwh = 1.5
		o.PltShowIniPop = true
		o.PltShowFinalPop = true
	}
	return
}

// Run runs optimisations
func (o *SimpleFltProb) Run(verbose bool) (nfeval int) {

	// benchmark
	if verbose {
		time0 := time.Now()
		defer func() {
			io.Pfblue2("\ncpu time = %v\n", time.Now().Sub(time0))
		}()
	}

	// Pareto front stat
	pareto_stat := false
	if o.C.Nova > 1 && o.ParF1F0 != nil && len(o.ParFmin) > 1 && len(o.ParFmax) > 1 {
		pareto_stat = true
		o.pareto_bins(0, 1)
		o.ParDisErr = make([]float64, o.C.Ntrials)
		o.ParSpread = make([]float64, o.C.Ntrials)
	}

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

		// number of function evaluations
		if itrial == 0 {
			nfeval = o.Evo.GetNfeval()
		}

		// results
		isl := 0
		xbest := o.Evo.Best.GetFloats()
		o.Fcn(o.ff[isl], o.gg[isl], o.hh[isl], xbest, isl)

		// check if best is unfeasible
		unfeasible := false
		for _, g := range o.gg[0] {
			if g < 0 {
				unfeasible = true
				break
			}
		}
		for _, h := range o.hh[0] {
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
			io.Pfyel("%3d x*="+o.NumfmtX+" f="+o.NumfmtF, itrial, xbest, o.ff[0])
			if o.ShowCts {
				if o.ng > 0 {
					io.Pfcyan(" g="+o.NumfmtF, o.gg[0])
				}
				if o.nh > 0 {
					io.Pfcyan(" h="+o.NumfmtF, o.hh[0])
				}
			}
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
				fcur := utl.DblCopy(o.ff[0])
				o.Fcn(o.ff[isl], o.gg[isl], o.hh[isl], o.Xbest[o.Nfeasible-1], isl)
				cur_dom, _ := utl.DblsParetoMin(fcur, o.ff[0])
				if cur_dom {
					o.PopsBest = o.Evo.GetPopulations()
				}
			}
		}

		// Pareto front
		if pareto_stat {
			o.ParDisErr[itrial], o.ParSpread[itrial] = o.pareto_front(0, 1)
		}
	}
	return
}

// Stat prints statistical analysis
func (o *SimpleFltProb) Stat(idxF, hlen int, Fref float64) (fmin, fave, fmax, fdev float64) {
	fmin, fave, fmax, fdev = 1e30, 1e30, 1e30, 1e30
	if o.Nfeasible < 1 {
		return
	}
	F := make([]float64, o.Nfeasible)
	isl := 0
	for i := 0; i < o.Nfeasible; i++ {
		o.Fcn(o.ff[isl], o.gg[isl], o.hh[isl], o.Xbest[i], isl)
		F[i] = o.ff[isl][idxF]
	}
	if o.Nfeasible < 2 {
		fmin, fave, fmax = F[0], F[0], F[0]
		return
	}
	fmin, fave, fmax, fdev = rnd.StatBasic(F, true)
	io.Pf("fmin = %v\n", fmin)
	io.PfYel("fave = %v (%v)\n", fave, Fref)
	io.Pf("fmax = %v\n", fmax)
	io.Pf("fdev = %v\n\n", fdev)
	io.Pf(rnd.BuildTextHist(nice_num(fmin-0.05), nice_num(fmax+0.05), 11, F, "%.2f", hlen))
	return
}

// StatPareto print stat about Pareto front
func (o *SimpleFltProb) StatPareto() {

	// distance error
	var emin, eave, emax, edev float64
	if len(o.ParDisErr) > 1 {
		emin, eave, emax, edev = rnd.StatBasic(o.ParDisErr, true)
	} else {
		emin, eave, emax = o.ParDisErr[0], o.ParDisErr[0], o.ParDisErr[0]
	}
	io.Pfcyan("\ndist_err_min = %g\n", emin)
	io.PfCyan("dist_err_ave = %g\n", eave)
	io.Pfcyan("dist_err_max = %g\n", emax)
	io.Pfcyan("dist_err_dev = %g\n", edev)

	// spread
	var smin, save, smax, sdev float64
	if len(o.ParSpread) > 1 {
		smin, save, smax, sdev = rnd.StatBasic(o.ParSpread, true)
	} else {
		smin, save, smax = o.ParSpread[0], o.ParSpread[0], o.ParSpread[0]
	}
	io.Pfgreen("\nspread_min = %g\n", smin)
	io.PfGreen("spread_ave = %g\n", save)
	io.Pfgreen("spread_max = %g\n", smax)
	io.Pfgreen("spread_dev = %g\n", sdev)
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
	if o.PltXrange != nil {
		xmin, xmax = o.PltXrange[0], o.PltXrange[1]
	}
	if o.PltYrange != nil {
		ymin, ymax = o.PltYrange[0], o.PltYrange[1]
	}

	// auxiliary variables
	X, Y := utl.MeshGrid2D(xmin, xmax, ymin, ymax, o.PltNpts, o.PltNpts)
	Zf := utl.DblsAlloc(o.PltNpts, o.PltNpts)
	var Zg [][][]float64
	var Zh [][][]float64
	if o.ng > 0 {
		Zg = utl.Deep3alloc(o.ng, o.PltNpts, o.PltNpts)
	}
	if o.nh > 0 {
		Zh = utl.Deep3alloc(o.nh, o.PltNpts, o.PltNpts)
	}

	// compute values
	x := make([]float64, 2)
	isl := 0
	for i := 0; i < o.PltNpts; i++ {
		for j := 0; j < o.PltNpts; j++ {
			x[0], x[1] = X[i][j], Y[i][j]
			o.Fcn(o.ff[isl], o.gg[isl], o.hh[isl], x, isl)
			Zf[i][j] = o.ff[0][o.PltIdxF]
			for k, g := range o.gg[0] {
				Zg[k][i][j] = g
			}
			for k, h := range o.hh[0] {
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
	if o.PltShowIniPop {
		for _, pop := range o.PopsIni {
			for _, ind := range pop {
				x := ind.GetFloats()
				plt.PlotOne(x[0], x[1], io.Sf("'k.', zorder=20, clip_on=0, label='%s'", l))
				l = ""
			}
		}
	}

	// final populations
	l = "final population"
	if o.PltShowFinalPop {
		for _, pop := range o.PopsBest {
			for _, ind := range pop {
				x := ind.GetFloats()
				plt.PlotOne(x[0], x[1], io.Sf("'ko', ms=6, zorder=30, clip_on=0, label='%s', markerfacecolor='none'", l))
				l = ""
			}
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
	x = make([]float64, nx)
	f = make([]float64, o.nf)
	g = make([]float64, o.ng)
	h = make([]float64, o.nh)
	copy(x, o.Xbest[0])
	isl := 0
	o.Fcn(f, g, h, x, isl)
	for i := 1; i < o.Nfeasible; i++ {
		o.Fcn(o.ff[isl], o.gg[isl], o.hh[isl], o.Xbest[i], isl)
		_, other_dom := utl.DblsParetoMin(f, o.ff[0])
		if other_dom {
			copy(x, o.Xbest[i])
			copy(f, o.ff[0])
			copy(g, o.gg[0])
			copy(h, o.hh[0])
		}
	}
	return
}

// pareto_bins sets bins touching solution front
func (o *SimpleFltProb) pareto_bins(I, J int) {
	o.ParBins.Init(o.ParFmin, []float64{o.ParFmax[I] * 1.1, o.ParFmax[J] * 1.1}, o.ParNdiv)
	o.ParSelB = make(map[int]bool)
	select_bin := func(pt []float64) {
		idx := o.ParBins.CalcIdx(pt)
		if idx >= 0 {
			o.ParSelB[idx] = true
		}
	}
	diag := math.Sqrt(math.Pow(o.ParFmax[I]-o.ParFmin[I], 2.0) + math.Pow(o.ParFmax[J]-o.ParFmin[J], 2.0))
	tmp := utl.LinSpace(o.ParFmin[I], o.ParFmax[I], 3*o.ParNdiv)
	pt := make([]float64, 2)
	for i := 0; i < len(tmp); i++ {
		f0, f1 := tmp[i], o.ParF1F0(tmp[i])
		pt[0], pt[1] = f0, f1
		select_bin(pt)
		for j := 0; j < o.ParNray; j++ {
			α := float64(j) * 2.0 * math.Pi / float64(o.ParNray)
			for k := 0; k < len(o.ParRadM); k++ {
				ρ := o.ParRadM[k] * diag
				δf0 := ρ * math.Cos(α)
				δf1 := ρ * math.Sin(α)
				pt[0], pt[1] = f0+δf0, f1+δf1
				select_bin(pt)
			}
		}
	}
}

// pareto_front computes stat about Pareto front
func (o *SimpleFltProb) pareto_front(idxf0, idxf1 int) (disterr, spread float64) {

	// Pareto-front
	feasible := o.Evo.GetFeasible()
	ovas, _ := o.Evo.GetResults(feasible)
	ovafront, _ := o.Evo.GetParetoFront(feasible, ovas, nil)
	f0front, f1front := o.Evo.GetFrontOvas(idxf0, idxf1, ovafront)
	f0fin := utl.DblsGetColumn(idxf0, ovas)
	f1fin := utl.DblsGetColumn(idxf1, ovas)

	// solution-quality: distance
	for i, f0 := range f0front {
		dist := math.Abs(f1front[i] - o.ParF1F0(f0))
		disterr = utl.Max(disterr, dist)
	}

	// solution-quality: spread
	pt := make([]float64, 2)
	for i := 0; i < len(f0fin); i++ {
		pt[0], pt[1] = f0fin[i], f1fin[i]
		if pt[0] < o.ParFmax[0]*1.1 && pt[1] < o.ParFmax[1]*1.1 {
			err := o.ParBins.Append(pt, i)
			if err != nil {
				chk.Panic("cannot append item:\n%v", err)
			}
		}
	}
	for idx, _ := range o.ParSelB {
		bin := o.ParBins.All[idx]
		if bin != nil {
			if len(bin.Entries) > 0 {
				spread += 1
			}
		}
	}
	spread = spread / float64(len(o.ParSelB))
	return
}

// nice_num returns a truncated float
func nice_num(x float64) float64 {
	s := io.Sf("%.2f", x)
	return io.Atof(s)
}

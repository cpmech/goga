// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_flt01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt01. quadratic with inequalities")

	// parameters
	C := NewConfParams()
	C.Pll = false
	C.Nisl = 1
	C.Ninds = 12
	C.GAtype = "crowd"
	C.Ops.FltCxName = "de"
	C.CrowdSize = 3
	C.RangeFlt = [][]float64{
		{-2, 2}, // gene # 0: min and max
		{-2, 2}, // gene # 1: min and max
	}
	C.PopFltGen = PopFltGen
	if chk.Verbose {
		C.DoPlot = chk.Verbose
	}
	C.CalcDerived()
	rnd.Init(C.Seed)

	// functions
	fcn := func(f, g, h []float64, x []float64, isl int) {
		f[0] = x[0]*x[0]/2.0 + x[1]*x[1] - x[0]*x[1] - 2.0*x[0] - 6.0*x[1]
		g[0] = 2.0 - x[0] - x[1]     // ≥ 0
		g[1] = 2.0 + x[0] - 2.0*x[1] // ≥ 0
		g[2] = 3.0 - 2.0*x[0] - x[1] // ≥ 0
		g[3] = x[0]                  // ≥ 0
		g[4] = x[1]                  // ≥ 0
	}

	// simple problem
	sim := NewSimpleFltProb(fcn, 1, 5, 0, C)
	sim.Run(chk.Verbose)

	// plot
	sim.Plot("test_flt01")
	C.Report("/tmp/goga", "tst_flt01")
}

func Test_flt02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt02. circle with equality constraint")

	// parameters
	C := NewConfParams()
	C.Eps1 = 1e-3
	C.Pll = false
	C.Nisl = 4
	C.Ninds = 12
	C.Ntrials = 1
	if chk.Verbose {
		C.Ntrials = 40
	}
	C.Verbose = false
	C.Dtmig = 50
	C.CrowdSize = 2
	C.CompProb = false
	C.GAtype = "crowd"
	C.Ops.FltCxName = "de"
	C.RangeFlt = [][]float64{
		{-1, 3}, // gene # 0: min and max
		{-1, 3}, // gene # 1: min and max
	}
	C.Latin = true
	C.PopFltGen = PopFltGen
	if chk.Verbose {
		C.FnKey = ""
		if C.Ntrials == 1 {
			C.DoPlot = true
		}
	}
	C.Ops.EnfRange = true
	C.NumFmts = map[string][]string{"flt": {"%8.4f", "%8.4f"}}
	C.ShowDem = true
	C.RegTol = 0.01
	C.CalcDerived()
	rnd.Init(C.Seed)

	// geometry
	xe := 1.0                      // centre of circle
	le := -0.4                     // selected level of f(x)
	ys := xe - (1.0+le)/math.Sqrt2 // coordinates of minimum point with level=le
	y0 := 2.0*ys + xe              // vertical axis intersect of straight line defined by c(x)
	xc := []float64{xe, xe}        // centre
	nx := len(xc)

	// functions
	fcn := func(f, g, h []float64, x []float64, isl int) {
		res := 0.0
		for i := 0; i < nx; i++ {
			res += (x[i] - xc[i]) * (x[i] - xc[i])
		}
		f[0] = math.Sqrt(res) - 1
		h[0] = x[0] + x[1] + xe - y0
	}

	// simple problem
	sim := NewSimpleFltProb(fcn, 1, 0, 1, C)
	sim.Run(chk.Verbose)

	// stat
	io.Pf("\n")
	sim.Stat(0, 60, -0.4)

	// plot
	sim.PltExtra = func() {
		plt.PlotOne(ys, ys, "'o', markeredgecolor='yellow', markerfacecolor='none', markersize=10")
	}
	sim.Plot("test_flt02")
}

func Test_flt03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt03. sin⁶(5 π x) multimodal")

	// configuration
	C := NewConfParams()
	C.Nova = 1
	C.Noor = 2
	C.Nisl = 4
	C.Ninds = 24
	C.GAtype = "crowd"
	C.Ops.FltCxName = "de"
	C.CrowdSize = 3
	C.ParetoPhi = 0.01
	C.CompProb = true
	C.Tf = 100
	C.Dtmig = 60
	C.RangeFlt = [][]float64{{0, 0.9999999999999}}
	C.PopFltGen = PopFltGen
	C.CalcDerived()
	rnd.Init(C.Seed)

	// post-processing function
	values := utl.Deep3alloc(C.Tf/10, C.Nisl, C.Ninds)
	C.PostProc = func(idIsland, time int, pop Population) {
		if time%10 == 0 {
			k := time / 10
			for i, ind := range pop {
				values[k][idIsland][i] = ind.GetFloat(0)
			}
		}
	}

	// functions
	yfcn := func(x float64) float64 { return math.Pow(math.Sin(5.0*math.Pi*x), 6.0) }
	fcn := func(f, g, h []float64, x []float64, isl int) {
		f[0] = -yfcn(x[0])
	}

	// simple problem
	sim := NewSimpleFltProb(fcn, 1, 0, 0, C)
	sim.Run(chk.Verbose)

	// write histograms and plot
	if chk.Verbose {

		// write histograms
		var buf bytes.Buffer
		hist := rnd.Histogram{Stations: utl.LinSpace(0, 1, 13)}
		for k := 0; k < C.Tf/10; k++ {
			for i := 0; i < C.Nisl; i++ {
				clear := false
				if i == 0 {
					clear = true
				}
				hist.Count(values[k][i], clear)
			}
			io.Ff(&buf, "\ntime=%d\n%v", k*10, rnd.TextHist(hist.GenLabels("%4.2f"), hist.Counts, 60))
		}
		io.WriteFileVD("/tmp/goga", "test_flt03_hist.txt", &buf)

		// plot
		plt.SetForEps(0.8, 300)
		xmin := sim.Evo.Islands[0].Pop[0].GetFloat(0)
		xmax := xmin
		for k := 0; k < C.Nisl; k++ {
			for _, ind := range sim.Evo.Islands[k].Pop {
				x := ind.GetFloat(0)
				y := yfcn(x)
				xmin = utl.Min(xmin, x)
				xmax = utl.Max(xmax, x)
				plt.PlotOne(x, y, "'r.',clip_on=0,zorder=20")
			}
		}
		np := 401
		X := utl.LinSpace(0, 1, np)
		Y := make([]float64, np)
		for i := 0; i < np; i++ {
			Y[i] = yfcn(X[i])
		}
		plt.Plot(X, Y, "'b-',clip_on=0,zorder=10")
		plt.Gll("$x$", "$y$", "")
		plt.SaveD("/tmp/goga", "test_flt03_func.eps")
	}
}

func Test_flt04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt04. two-bar truss. Pareto-optimal")

	// configuration
	C := NewConfParams()
	C.Nisl = 4
	C.Ninds = 24
	C.GAtype = "crowd"
	C.Ops.FltCxName = "de"
	C.Ops.DEpc = 0.1
	C.Ops.DEmult = 0.5
	C.CrowdSize = 3
	C.ParetoPhi = 0.05
	C.Tf = 100
	C.Dtmig = 10
	C.RangeFlt = [][]float64{{0.1, 2.25}, {0.5, 2.5}}
	C.PopFltGen = PopFltGen
	C.CalcDerived()
	rnd.Init(C.Seed)

	// data
	// from Coelho (2007) page 19
	ρ := 0.283 // lb/in³
	H := 100.0 // in
	P := 1e4   // lb
	E := 3e7   // lb/in²
	σ0 := 2e4  // lb/in²

	// functions
	TSQ2 := 2.0 * math.Sqrt2
	fcn := func(f, g, h []float64, x []float64, isl int) {
		f[0] = 2.0 * ρ * H * x[1] * math.Sqrt(1.0+x[0]*x[0])
		f[1] = P * H * math.Pow(1.0+x[0]*x[0], 1.5) * math.Sqrt(1.0+math.Pow(x[0], 4.0)) / (TSQ2 * E * x[0] * x[0] * x[1])
		g[0] = σ0 - P*(1.0+x[0])*math.Sqrt(1.0+x[0]*x[0])/(TSQ2*x[0]*x[1])
		g[1] = σ0 - P*(1.0-x[0])*math.Sqrt(1.0+x[0]*x[0])/(TSQ2*x[0]*x[1])
	}

	// objective value function
	C.OvaOor = func(ind *Individual, isl, t int, report *bytes.Buffer) {
		x := ind.GetFloats()
		f := make([]float64, 2)
		g := make([]float64, 2)
		fcn(f, g, nil, x, isl)
		ind.Ovas[0] = f[0]
		ind.Ovas[1] = f[1]
		ind.Oors[0] = utl.GtePenalty(g[0], 0, 1)
		ind.Oors[1] = utl.GtePenalty(g[1], 0, 1)
	}

	// simple problem
	sim := NewSimpleFltProb(fcn, 2, 2, 0, C)
	sim.Run(chk.Verbose)

	// results
	if chk.Verbose {

		// reference data
		_, dat, _ := io.ReadTable("data/coelho-fig1.6.dat")

		// Pareto-front
		feasible := sim.Evo.GetFeasible()
		ovas, _ := sim.Evo.GetResults(feasible)
		ovafront, _ := sim.Evo.GetParetoFront(feasible, ovas, nil)
		xova, yova := sim.Evo.GetFrontOvas(0, 1, ovafront)

		// plot
		plt.SetForEps(0.75, 355)
		plt.Plot(dat["f1"], dat["f2"], "'b-',ms=3")
		x := utl.DblsGetColumn(0, ovas)
		y := utl.DblsGetColumn(1, ovas)
		plt.Plot(x, y, "'r.'")
		plt.Plot(xova, yova, "'ko',markerfacecolor='none',ms=6")
		plt.Gll("$f_1$", "$f_2$", "")
		plt.SaveD("/tmp/goga", "test_flt04.eps")
	}
}

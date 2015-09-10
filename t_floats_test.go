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

	// initialise random numbers generator
	rnd.Init(0) // 0 => use current time as seed

	// parameters
	C := NewConfParams()
	C.Pll = false
	C.Nisl = 1
	C.Ninds = 12
	C.GAtype = "crowd"
	C.CrowdSize = 3
	C.DiffEvol = true
	C.RangeFlt = [][]float64{
		{-2, 2}, // gene # 0: min and max
		{-2, 2}, // gene # 1: min and max
	}
	C.PopFltGen = PopFltGen
	if chk.Verbose {
		C.DoPlot = chk.Verbose
	}
	C.CalcDerived()

	// functions
	fcn := func(f, g, h []float64, x []float64) {
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
}

func Test_flt02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt02. circle with equality constraint")

	// initialise random numbers generator
	rnd.Init(0)

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
	C.CrowdSize = 3
	C.CompProb = false
	C.GAtype = "crowd"
	C.DiffEvol = true
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

	// geometry
	xe := 1.0                      // centre of circle
	le := -0.4                     // selected level of f(x)
	ys := xe - (1.0+le)/math.Sqrt2 // coordinates of minimum point with level=le
	y0 := 2.0*ys + xe              // vertical axis intersect of straight line defined by c(x)
	xc := []float64{xe, xe}        // centre
	nx := len(xc)

	// functions
	fcn := func(f, g, h []float64, x []float64) {
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
	C.Ninds = 12
	C.GAtype = "crowd"
	C.CrowdSize = 2
	C.ParetoPhi = 0.01
	C.CompProb = true
	C.Elite = false
	C.Noise = 0.05
	C.DoPlot = false
	C.RegTol = 0
	C.Tf = 100
	C.Dtmig = 60
	C.RangeFlt = [][]float64{{0, 1}}
	C.PopFltGen = PopFltGen
	C.SetNbasesFixOp(8)
	C.CalcDerived()

	// initialise random numbers generator
	rnd.Init(0)

	// function
	yfcn := func(x float64) float64 {
		return math.Pow(math.Sin(5.0*math.Pi*x), 6.0)
	}

	// objective value function
	C.OvaOor = func(ind *Individual, idIsland, t int, report *bytes.Buffer) {
		x := ind.GetFloat(0)
		ind.Ovas[0] = -yfcn(x)
		ind.Oors[0] = utl.GtePenalty(x, 0, 1)
		ind.Oors[1] = utl.GtePenalty(1, x, 1)
	}

	// post-processing function
	values := utl.Deep3alloc(C.Tf/10, C.Nisl, C.Ninds)
	C.PostProc = func(idIsland, time int, pop Population) {
		if time%10 == 0 && false {
			k := time / 10
			for i, ind := range pop {
				values[k][idIsland][i] = ind.GetFloat(0)
			}
		}
	}

	// run
	evo := NewEvolver(C)
	evo.Run()

	// print population
	for _, isl := range evo.Islands {
		io.Pf("%v", isl.Pop.Output(C))
	}

	// write histograms and plot
	if chk.Verbose {

		// write histograms
		var buf bytes.Buffer
		hist := rnd.Histogram{Stations: utl.LinSpace(0, 1, 26)}
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
		xmin := evo.Islands[0].Pop[0].GetFloat(0)
		xmax := xmin
		for k := 0; k < C.Nisl; k++ {
			for _, ind := range evo.Islands[k].Pop {
				x := ind.GetFloat(0)
				y := yfcn(x)
				xmin = utl.Min(xmin, x)
				xmax = utl.Max(xmax, x)
				plt.PlotOne(x, y, "'r.',clip_on=0,zorder=20")
			}
		}
		np := 401
		//X := utl.LinSpace(xmin, xmax, np)
		X := utl.LinSpace(0, 1, np)
		Y := make([]float64, np)
		for i := 0; i < np; i++ {
			Y[i] = yfcn(X[i])
		}
		plt.Plot(X, Y, "'b-',clip_on=0,zorder=10")
		plt.Gll("$x$", "$y$", "")
		//plt.AxisXrange(0, 1)
		plt.SaveD("/tmp/goga", "test_flt03_func.eps")
	}
}

func Test_flt04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flt04. two-bar truss. Pareto-optimal")

	// configuration
	C := NewConfParams()
	C.Nova = 2
	C.Noor = 2
	C.Nisl = 4
	C.Ninds = 24
	//C.GAtype = "std"
	C.GAtype = "crowd"
	//C.GAtype = "sharing"
	C.DiffEvol = true
	C.Elite = false
	C.CrowdSize = 3
	C.ParetoPhi = 0.05
	C.ShAlp = 0.5
	C.ShSig = 0.001
	C.ShPhen = false
	C.Noise = 0.05
	C.DoPlot = false
	C.RegTol = 0
	C.Tf = 100
	C.Dtmig = 25
	C.RangeFlt = [][]float64{{0.1, 2.25}, {0.5, 2.5}}
	C.PopFltGen = PopFltGen
	C.Latin = true
	C.CalcDerived()

	// initialise random numbers generator
	rnd.Init(0)

	// data
	// from Coelho (2007) page 19
	ρ := 0.283 // lb/in³
	h := 100.0 // in
	P := 1e4   // lb
	E := 3e7   // lb/in²
	σ0 := 2e4  // lb/in²

	// functions
	twosq2 := 2.0 * math.Sqrt2
	f1 := func(x []float64) float64 {
		return 2.0 * ρ * h * x[1] * math.Sqrt(1.0+x[0]*x[0])
	}
	f2 := func(x []float64) float64 {
		return P * h * math.Pow(1.0+x[0]*x[0], 1.5) * math.Sqrt(1.0+math.Pow(x[0], 4.0)) / (twosq2 * E * x[0] * x[0] * x[1])
	}
	g1 := func(x []float64) float64 {
		return P*(1.0+x[0])*math.Sqrt(1.0+x[0]*x[0])/(twosq2*x[0]*x[1]) - σ0
	}
	g2 := func(x []float64) float64 {
		return P*(1.0-x[0])*math.Sqrt(1.0+x[0]*x[0])/(twosq2*x[0]*x[1]) - σ0
	}

	// objective value function
	C.OvaOor = func(ind *Individual, idIsland, t int, report *bytes.Buffer) {
		x := ind.GetFloats()
		ind.Ovas[0] = f1(x)
		ind.Ovas[1] = f2(x)
		ind.Oors[0] = utl.GtePenalty(0, g1(x), 1)
		ind.Oors[1] = utl.GtePenalty(0, g2(x), 1)
		//ind.Oors[2] = utl.GtePenalty(x[0], 0, 1)
		//ind.Oors[3] = utl.GtePenalty(x[1], 0, 1)
	}

	// run
	evo := NewEvolver(C)
	evo.Run()

	// results
	if chk.Verbose {
		_, dat, _ := io.ReadTable("data/coelho-fig1.6.dat")
		feasible := evo.GetFeasible()
		ovas, _ := evo.GetResults(feasible)
		ovafront, _ := evo.GetParetoFront(feasible, ovas, nil)
		xova, yova := evo.GetFrontOvas(0, 1, ovafront)
		//for _, ind := range feasible {
		//x := ind.GetFloats()
		//io.Pforan("f1=%8.4f f2=%8.4f g1=%12.4f g2=%12.4f\n", f1(x), f2(x), g1(x), g2(x))
		//io.Pfyel("ovas = %v\n", ind.Ovas)
		//io.Pfpink("oors = %v\n", ind.Oors)
		//}
		plt.SetForEps(0.75, 355)
		plt.Plot(dat["f1"], dat["f2"], "'k+',ms=3")
		x := utl.DblsGetColumn(0, ovas)
		y := utl.DblsGetColumn(1, ovas)
		plt.Plot(x, y, "'r.'")
		plt.Plot(xova, yova, "'ko',markerfacecolor='none',ms=6")
		plt.Gll("$f_1$", "$f_2$", "")
		plt.SaveD("/tmp/goga", "test_flt04.eps")
	}
}

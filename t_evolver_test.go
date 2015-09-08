// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_evo01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo01. organise sequence of ints")
	io.Pf("\n")

	// initialise random numbers generator
	rnd.Init(0) // 0 => use current time as seed

	// parameters
	C := NewConfParams()
	C.Nova = 1
	C.Noor = 0
	C.Nisl = 1
	C.Ninds = 20
	C.RegTol = 0
	C.NumInts = 20
	//C.GAtype = "crowd"
	C.CrowdSize = 2
	C.Tf = 50
	C.Verbose = chk.Verbose
	C.CalcDerived()

	// mutation function
	C.Ops.MtInt = func(A []int, time int, ops *OpsData) {
		size := len(A)
		if !rnd.FlipCoin(ops.Pm) || size < 1 {
			return
		}
		pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
		for _, i := range pos {
			if A[i] == 1 {
				A[i] = 0
			}
			if A[i] == 0 {
				A[i] = 1
			}
		}
	}

	// generation function
	C.PopIntGen = func(id int, cc *ConfParams) Population {
		o := make([]*Individual, cc.Ninds)
		genes := make([]int, cc.NumInts)
		for i := 0; i < cc.Ninds; i++ {
			for j := 0; j < cc.NumInts; j++ {
				genes[j] = rand.Intn(2)
			}
			o[i] = NewIndividual(cc.Nova, cc.Noor, cc.Nbases, genes)
		}
		return o
	}

	// objective function
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		score := 0.0
		count := 0
		for _, val := range ind.Ints {
			if val == 0 && count%2 == 0 {
				score += 1.0
			}
			if val == 1 && count%2 != 0 {
				score += 1.0
			}
			count++
		}
		ind.Ovas[0] = 1.0 / (1.0 + score)
		return
	}

	// run optimisation
	evo := NewEvolver(C)
	evo.Run()

	// results
	ideal := 1.0 / (1.0 + float64(C.NumInts))
	io.PfGreen("\nBest = %v\nBestOV = %v  (ideal=%v)\n", evo.Best.Ints, evo.Best.Ovas[0], ideal)
}

func Test_evo02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo02")

	// initialise random numbers generator
	rnd.Init(0) // 0 => use current time as seed
	//rnd.Init(1111) // 0 => use current time as seed

	// parameters
	C := NewConfParams()
	C.Nova = 1
	C.Noor = 5
	C.Pll = false
	C.Nisl = 1
	C.Ninds = 20
	C.RegTol = 0
	//C.GAtype = "std"
	C.GAtype = "crowd"
	C.CrowdSize = 2
	C.ParetoPhi = 0.01
	//C.Elite = true
	C.Verbose = false
	C.RangeFlt = [][]float64{
		{-2, 2}, // gene # 0: min and max
		{-2, 2}, // gene # 1: min and max
	}
	C.PopFltGen = PopFltGen
	if chk.Verbose {
		C.FnKey = "test_evo02"
		C.DoPlot = false
	}
	//C.SetNbasesFixOp(8)
	C.CalcDerived()

	f := func(x []float64) float64 { return x[0]*x[0]/2.0 + x[1]*x[1] - x[0]*x[1] - 2.0*x[0] - 6.0*x[1] }
	c1 := func(x []float64) float64 { return x[0] + x[1] - 2.0 }      // ≤ 0
	c2 := func(x []float64) float64 { return -x[0] + 2.0*x[1] - 2.0 } // ≤ 0
	c3 := func(x []float64) float64 { return 2.0*x[0] + x[1] - 3.0 }  // ≤ 0
	c4 := func(x []float64) float64 { return -x[0] }                  // ≤ 0
	c5 := func(x []float64) float64 { return -x[1] }                  // ≤ 0

	// objective function
	p := 1.0
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		x := ind.GetFloats()
		ind.Ovas[0] = f(x)
		ind.Oors[0] = utl.GtePenalty(0, c1(x), p)
		ind.Oors[1] = utl.GtePenalty(0, c2(x), p)
		ind.Oors[2] = utl.GtePenalty(0, c3(x), p)
		ind.Oors[3] = utl.GtePenalty(0, c4(x), p)
		ind.Oors[4] = utl.GtePenalty(0, c5(x), p)
		return
	}

	// evolver
	evo := NewEvolver(C)
	pop0 := evo.Islands[0].Pop.GetCopy()
	evo.Run()

	// results
	io.PfGreen("\nx=%g (%g)\n", evo.Best.GetFloat(0), 2.0/3.0)
	io.PfGreen("y=%g (%g)\n", evo.Best.GetFloat(1), 4.0/3.0)
	io.PfGreen("BestOV=%g (%g)\n", evo.Best.Ovas[0], f([]float64{2.0 / 3.0, 4.0 / 3.0}))

	// plot contour
	if C.DoPlot {
		PlotTwoVarsContour("/tmp/goga", "contour_evo02", pop0, evo.Islands[0].Pop, evo.Best, 41, 2, "", nil, false, true,
			C.RangeFlt, false, false, nil, nil, f, c1, c2, c3, c4, c5)
	}
}

func Test_evo03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo03")

	rnd.Init(0)

	// parameters
	C := NewConfParams()
	C.Nova = 1
	C.Noor = 1
	C.Pll = false
	C.Nisl = 4
	C.Ninds = 12
	C.Ntrials = 1
	if chk.Verbose {
		C.Ntrials = 40
	}
	C.Verbose = false
	C.Dtmig = 50
	C.Ops.Pm = 0.1
	C.CrowdSize = 2
	C.ParetoPhi = 0
	//C.GAtype = "std"
	C.GAtype = "crowd"
	//C.GAtype = "sharing"
	C.Elite = false
	C.RangeFlt = [][]float64{
		{-1, 3}, // gene # 0: min and max
		{-1, 3}, // gene # 1: min and max
	}
	C.Latin = true
	C.PopFltGen = PopFltGen
	if chk.Verbose {
		C.FnKey = "" //"test_evo03"
		C.DoPlot = false
	}
	//C.SetBlxMwicz()
	//C.SetNbasesFixOp(8)
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
	f := func(x []float64) (res float64) {
		for i := 0; i < nx; i++ {
			res += (x[i] - xc[i]) * (x[i] - xc[i])
		}
		return math.Sqrt(res) - 1
	}
	c := func(x []float64) (res float64) {
		return x[0] + x[1] + xe - y0
	}

	// objective function
	p := 1.0
	C.OvaOor = func(ind *Individual, idIsland, time int, report *bytes.Buffer) {
		x := ind.GetFloats()
		fp := utl.GtePenalty(1e-2, math.Abs(c(x)), p)
		ind.Ovas[0] = f(x) + fp
		ind.Oors[0] = fp
		return
	}

	// evolver
	evo := NewEvolver(C)

	// run ntrials times
	pops0 := make([]Population, C.Nisl)
	for i := 0; i < C.Ntrials; i++ {

		// reset populations
		if i > 0 {
			evo.ResetAllPop()
		}

		// initial populations
		for k, isl := range evo.Islands {
			pops0[k] = isl.Pop.GetCopy()
		}

		// run
		evo.Run()

		// results
		if false {
			xbest := []float64{evo.Best.GetFloat(0), evo.Best.GetFloat(1)}
			io.PfGreen("\nx=%g (%g)\n", xbest[0], ys)
			io.PfGreen("y=%g (%g)\n", xbest[1], ys)
		}
		ova := evo.Best.Ovas[0]
		if ova > 0 {
			io.PfRed("BestOV=%g (%g)\n", ova, le)
		} else if math.Abs(ova)-0.25 < 0.1 {
			io.Pforan("BestOV=%g (%g)\n", ova, le)
		} else {
			io.PfGreen("BestOV=%g (%g)\n", ova, le)
		}
		//io.Pf("%v\n", evo.Islands[0].Pop.Output(C))

		// plot contour
		if C.DoPlot {
			extra := func() {
				plt.PlotOne(ys, ys, "'o', markeredgecolor='yellow', markerfacecolor='none', markersize=10")
				for k := 1; k < C.Nisl; k++ {
					for _, ind := range pops0[k] {
						v := ind.GetFloats()
						plt.PlotOne(v[0], v[1], "'k.', zorder=20, clip_on=0")
					}
					for _, ind := range evo.Islands[k].Pop {
						v := ind.GetFloats()
						plt.PlotOne(v[0], v[1], "'ko', ms=6, zorder=30, clip_on=0, markerfacecolor='none'")
					}
				}
			}
			PlotTwoVarsContour("/tmp/goga", io.Sf("contour_evo03_%02d", i), pops0[0], evo.Islands[0].Pop, evo.Best, 41, 2, "", extra, false, true,
				C.RangeFlt, false, false, nil, nil, f, c)
		}
	}
}

func Test_evo04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo04. TSP")

	// location / coordinates of stations
	locations := [][]float64{
		{60, 200}, {180, 200}, {80, 180}, {140, 180}, {20, 160}, {100, 160}, {200, 160},
		{140, 140}, {40, 120}, {100, 120}, {180, 100}, {60, 80}, {120, 80}, {180, 60},
		{20, 40}, {100, 40}, {200, 40}, {20, 20}, {60, 20}, {160, 20},
	}
	nstations := len(locations)

	// parameters
	C := NewConfParams()
	C.Nova = 1
	C.Noor = 0
	C.Nisl = 1
	C.Ninds = 20
	C.RegTol = 0.3
	C.RegPct = 0.2
	//C.Dtmig = 30
	C.GAtype = "crowd"
	C.ParetoPhi = 0.1
	C.Elite = false
	C.DoPlot = false //chk.Verbose
	//C.Rws = true
	C.SetIntOrd(nstations)
	C.CalcDerived()

	// initialise random numbers generator
	rnd.Init(0)

	// objective value function
	C.OvaOor = func(ind *Individual, idIsland, t int, report *bytes.Buffer) {
		L := locations
		ids := ind.Ints
		//io.Pforan("ids = %v\n", ids)
		dist := 0.0
		for i := 1; i < nstations; i++ {
			a, b := ids[i-1], ids[i]
			dist += math.Sqrt(math.Pow(L[b][0]-L[a][0], 2.0) + math.Pow(L[b][1]-L[a][1], 2.0))
		}
		a, b := ids[nstations-1], ids[0]
		dist += math.Sqrt(math.Pow(L[b][0]-L[a][0], 2.0) + math.Pow(L[b][1]-L[a][1], 2.0))
		ind.Ovas[0] = dist
		return
	}

	// evolver
	evo := NewEvolver(C)

	// print initial population
	pop := evo.Islands[0].Pop
	//io.Pf("\n%v\n", pop.Output(nil, false))

	// 0,4,8,11,14,17,18,15,12,19,13,16,10,6,1,3,7,9,5,2 894.363
	if false {
		for i, x := range []int{0, 4, 8, 11, 14, 17, 18, 15, 12, 19, 13, 16, 10, 6, 1, 3, 7, 9, 5, 2} {
			pop[0].Ints[i] = x
		}
		evo.Islands[0].CalcOvs(pop, 0)
		evo.Islands[0].CalcDemeritsAndSort(pop)
	}

	// check initial population
	ints := make([]int, nstations)
	if false {
		for i := 0; i < C.Ninds; i++ {
			for j := 0; j < nstations; j++ {
				ints[j] = pop[i].Ints[j]
			}
			sort.Ints(ints)
			chk.Ints(tst, "ints", ints, utl.IntRange(nstations))
		}
	}

	// run
	evo.Run()
	//io.Pf("%v\n", pop.Output(nil, false))
	io.Pfgreen("best = %v\n", evo.Best.Ints)
	io.Pfgreen("best OVA = %v  (871.117353844847)\n\n", evo.Best.Ovas[0])

	// best = [18 17 14 11 8 4 0 2 5 9 12 7 6 1 3 10 16 13 19 15]
	// best OVA = 953.4643474956656

	// best = [8 11 14 17 18 15 12 19 16 13 10 6 1 3 7 9 5 2 0 4]
	// best OVA = 871.117353844847

	// best = [5 2 0 4 8 11 14 17 18 15 12 19 16 13 10 6 1 3 7 9]
	// best OVA = 871.1173538448469

	// best = [6 10 13 16 19 15 18 17 14 11 8 4 0 2 5 9 12 7 3 1]
	// best OVA = 880.7760751923065

	// check final population
	if false {
		for i := 0; i < C.Ninds; i++ {
			for j := 0; j < nstations; j++ {
				ints[j] = pop[i].Ints[j]
			}
			sort.Ints(ints)
			chk.Ints(tst, "ints", ints, utl.IntRange(nstations))
		}
	}

	// plot travelling salesman path
	if C.DoPlot {
		plt.SetForEps(1, 300)
		X, Y := make([]float64, nstations), make([]float64, nstations)
		for k, id := range evo.Best.Ints {
			X[k], Y[k] = locations[id][0], locations[id][1]
			plt.PlotOne(X[k], Y[k], "'r.', ms=5, clip_on=0, zorder=20")
			plt.Text(X[k], Y[k], io.Sf("%d", id), "fontsize=7, clip_on=0, zorder=30")
		}
		plt.Plot(X, Y, "'b-', clip_on=0, zorder=10")
		plt.Plot([]float64{X[0], X[nstations-1]}, []float64{Y[0], Y[nstations-1]}, "'b-', clip_on=0, zorder=10")
		plt.Equal()
		plt.AxisRange(10, 210, 10, 210)
		plt.Gll("$x$", "$y$", "")
		plt.SaveD("/tmp/goga", "test_evo04.eps")
	}
}

func Test_evo05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo05. sin⁶(5 π x)")

	// configuration
	C := NewConfParams()
	C.Nova = 1
	C.Noor = 2
	C.Nisl = 4
	C.Ninds = 12
	C.GAtype = "crowd"
	//C.GAtype = "sharing"
	C.CrowdSize = 2
	C.ParetoPhi = 0.01
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
		io.WriteFileVD("/tmp/goga", "test_evo05_hist.txt", &buf)

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
		plt.SaveD("/tmp/goga", "test_evo05_func.eps")
	}
}

func Test_evo06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("evo06. two-bar truss. Pareto-optimal")

	// configuration
	C := NewConfParams()
	C.Nova = 2
	C.Noor = 2
	C.Nisl = 4
	C.Ninds = 24
	//C.GAtype = "std"
	C.GAtype = "crowd"
	//C.GAtype = "sharing"
	C.Elite = true
	C.CrowdSize = 2
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
		plt.SaveD("/tmp/goga", "test_evo06.eps")
	}
}

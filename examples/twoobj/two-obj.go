// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func solve_problem(problem int) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Ncpu = 1
	opt.Tf = 500
	opt.Verbose = false
	opt.Nsamples = 1000
	opt.GenType = "latin"
	opt.DEC = 0.1

	// options for report
	opt.HistNsta = 6
	opt.HistLen = 13
	opt.RptFmtE = "%.4e"
	opt.RptFmtL = "%.4e"
	opt.RptFmtEdev = "%.3e"
	opt.RptFmtLdev = "%.3e"

	// problem variables
	var fmin, fmax []float64
	var nf, ng, nh int
	var fcn goga.MinProb_t

	// problems
	switch problem {

	// problem # 1 -- ZDT1, Deb 2001, p356
	case 1:
		opt.Ncpu = 6
		opt.RptName = "ZDT1"
		n := 30
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = x[0]
			sum := 0.0
			for i := 1; i < n; i++ {
				sum += x[i]
			}
			c0 := 1.0 + 9.0*sum/float64(n-1)
			c1 := 1.0 - math.Sqrt(f[0]/c0)
			f[1] = c0 * c1
		}
		fmin = []float64{0, 0}
		fmax = []float64{1, 1}
		opt.F1F0_func = func(f0 float64) float64 {
			return 1.0 - math.Sqrt(f0)
		}
		// arc length = sqrt(5)/2 + log(sqrt(5)+2)/4 ≈ 1.4789428575445975
		opt.F1F0_arcLenRef = math.Sqrt(5.0)/2.0 + math.Log(math.Sqrt(5.0)+2.0)/4.0

	// problem # 2 -- ZDT2, Deb 2001, p356
	case 2:
		opt.Ncpu = 6
		opt.RptName = "ZDT2"
		n := 30
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = x[0]
			sum := 0.0
			for i := 1; i < n; i++ {
				sum += x[i]
			}
			c0 := 1.0 + 9.0*sum/float64(n-1)
			c1 := 1.0 - math.Pow(f[0]/c0, 2.0)
			f[1] = c0 * c1
		}
		fmin = []float64{0, 0}
		fmax = []float64{1, 1}
		opt.F1F0_func = func(f0 float64) float64 {
			return 1.0 - math.Pow(f0, 2.0)
		}
		// arc length = sqrt(5)/2 + log(sqrt(5)+2)/4 ≈ 1.4789428575445975
		opt.F1F0_arcLenRef = math.Sqrt(5.0)/2.0 + math.Log(math.Sqrt(5.0)+2.0)/4.0

	// problem # 3 -- ZDT3, Deb 2001, p356
	case 3:
		opt.Ncpu = 6
		opt.RptName = "ZDT3"
		n := 30
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = x[0]
			sum := 0.0
			for i := 1; i < n; i++ {
				sum += x[i]
			}
			c0 := 1.0 + 9.0*sum/float64(n-1)
			c1 := 1.0 - math.Sqrt(f[0]/c0) - math.Sin(10.0*math.Pi*f[0])*f[0]/c0
			f[1] = c0 * c1
		}
		fmin = []float64{0, -1}
		fmax = []float64{1, 1}
		opt.F1F0_func = func(f0 float64) float64 {
			return 1.0 - math.Sqrt(f0) - math.Sin(10.0*math.Pi*f0)*f0
		}
		opt.F1F0_f0ranges = [][]float64{
			{0.000000100000000, 0.083001534925223},
			{0.182228728029413, 0.257762363387862},
			{0.409313674808657, 0.453882104088830},
			{0.618396794416602, 0.652511703804663},
			{0.823331798326633, 0.851832865436414},
		}
		opt.F1F0_arcLenRef = 1.811

	// problem # 4 -- ZDT4, Deb 2001, p358
	case 4:
		opt.Ncpu = 2
		opt.RptName = "ZDT4"
		n := 10
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		opt.FltMin[0] = 0
		opt.FltMax[0] = 1
		for i := 1; i < n; i++ {
			opt.FltMin[i] = -5
			opt.FltMax[i] = 5
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = x[0]
			sum := 0.0
			w := 4.0 * math.Pi
			for i := 1; i < n; i++ {
				sum += x[i]*x[i] - 10.0*math.Cos(w*x[i])
			}
			c0 := 1.0 + 10.0*float64(n-1) + sum
			c1 := 1.0 - math.Sqrt(f[0]/c0)
			f[1] = c0 * c1
		}
		fmin = []float64{0, 0}
		fmax = []float64{1, 1}
		opt.F1F0_func = func(f0 float64) float64 {
			return 1.0 - math.Sqrt(f0)
		}
		// arc length = sqrt(5)/2 + log(sqrt(5)+2)/4 ≈ 1.4789428575445975
		opt.F1F0_arcLenRef = math.Sqrt(5.0)/2.0 + math.Log(math.Sqrt(5.0)+2.0)/4.0

	// problem # 5 -- FON (Fonseca and Fleming 1995), Deb 2001, p339
	case 5:
		opt.DEC = 0.8
		opt.Ncpu = 2
		opt.RptName = "FON"
		n := 10
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = -4
			opt.FltMax[i] = 4
		}
		nf, ng, nh = 2, 0, 0
		coef := 1.0 / math.Sqrt(float64(n))
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			sum0, sum1 := 0.0, 0.0
			for i := 0; i < n; i++ {
				sum0 += math.Pow(x[i]-coef, 2.0)
				sum1 += math.Pow(x[i]+coef, 2.0)
			}
			f[0] = 1.0 - math.Exp(-sum0)
			f[1] = 1.0 - math.Exp(-sum1)
		}
		fmin = []float64{0, 0}
		fmax = []float64{0.98, 1}
		opt.F1F0_func = func(f0 float64) float64 {
			return 1.0 - math.Exp(-math.Pow(2.0-math.Sqrt(-math.Log(1.0-f0)), 2.0))
		}
		opt.F1F0_arcLenRef = 1.45831385

	// problem # 6 -- ZDT6, Deb 2001, p360
	case 6:
		opt.Ncpu = 2
		opt.RptName = "ZDT6"
		n := 10
		opt.FltMin = make([]float64, n)
		opt.FltMax = make([]float64, n)
		for i := 0; i < n; i++ {
			opt.FltMin[i] = 0
			opt.FltMax[i] = 1
		}
		nf, ng, nh = 2, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			w := 6.0 * math.Pi
			f[0] = 1.0 - math.Exp(-4.0*x[0])*math.Pow(math.Sin(w*x[0]), 6.0)
			sum := 0.0
			for i := 1; i < n; i++ {
				sum += x[i]
			}
			w = float64(n - 1)
			c0 := 1.0 + 9.0*math.Pow(sum/w, 0.25)
			c1 := 1.0 - math.Pow(f[0]/c0, 2.0)
			f[1] = c0 * c1
		}
		opt.F1F0_func = func(f0 float64) float64 {
			return 1.0 - math.Pow(f0, 2.0)
		}
		xs := math.Atan(9.0*math.Pi) / (6.0 * math.Pi)
		f0min := 1.0 - math.Exp(-4.0*xs)*math.Pow(math.Sin(6.0*math.Pi*xs), 6.0)
		f1max := opt.F1F0_func(f0min)
		io.Pforan("xs=%v f0min=%v f1max=%v\n", xs, f0min, f1max)
		// xs=0.08145779687998356 f0min=0.2807753188153699 f1max=0.9211652203441274
		fmin = []float64{f0min, 0}
		fmax = []float64{1, 1}
		opt.F1F0_arcLenRef = 1.184

	default:
		chk.Panic("problem %d is not available", problem)
	}

	// number of trial solutions and number of groups
	opt.Nsol = len(opt.FltMin) * 10

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// initial solutions
	var sols0 []*goga.Solution
	if false {
		sols0 = opt.GetSolutionsCopy()
	}

	// solve
	opt.RunMany("", "")
	goga.StatF1F0(opt, true)

	// check
	goga.CheckFront0(opt, true)

	// plot
	if true {
		feasibleOnly := true
		plt.SetForEps(0.8, 300)
		fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
		fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
		goga.PlotOvaOvaPareto(opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
		np := 201
		F0 := utl.LinSpace(fmin[0], fmax[0], np)
		F1 := make([]float64, np)
		for i := 0; i < np; i++ {
			F1[i] = opt.F1F0_func(F0[i])
		}
		plt.Plot(F0, F1, io.Sf("'b-', label='%s'", opt.RptName))
		for _, f0vals := range opt.F1F0_f0ranges {
			f0A, f0B := f0vals[0], f0vals[1]
			f1A, f1B := opt.F1F0_func(f0A), opt.F1F0_func(f0B)
			plt.PlotOne(f0A, f1A, "'g_', mew=1.5, ms=10, clip_on=0")
			plt.PlotOne(f0B, f1B, "'g|', mew=1.5, ms=10, clip_on=0")
		}
		plt.AxisRange(fmin[0], fmax[0], fmin[1], fmax[1])
		plt.Gll("$f_0$", "$f_1$", "")
		plt.SaveD("/tmp/goga", io.Sf("%s.eps", opt.RptName))
	}
	return
}

func main() {
	textSize := `\scriptsize  \setlength{\tabcolsep}{0.5em}`
	miniPageSz, histTextSize := "4.1cm", `\fontsize{5pt}{6pt}`
	P := utl.IntRange2(1, 7)
	//P := []int{1, 2, 4, 6}
	//P := []int{4, 5, 6}
	//P := []int{4}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem)
	}
	io.Pf("\n-------------------------- generating report --------------------------\nn")
	nRowPerTab := 9
	title := "Unconstrained two objective problems"
	goga.TexReport("/tmp/goga", "tmp_two-obj", title, "two-obj", 2, nRowPerTab, true, false, textSize, miniPageSz, histTextSize, opts)
	goga.TexReport("/tmp/goga", "two-obj", title, "two-obj", 2, nRowPerTab, false, false, textSize, miniPageSz, histTextSize, opts)
}

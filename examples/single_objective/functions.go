// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func solve_problem(problem, ntrials int, checkOnly bool) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Ncpu = 1
	opt.Tf = 500
	opt.EpsH = 0.01
	opt.Verbose = false
	opt.Ntrials = ntrials
	opt.GenType = "latin"

	// options for report
	opt.RptFmtF = "%.7f"
	opt.RptFmtFdev = "%.7f"
	opt.RptFmtX = "%.5f"

	// problem variables
	var ng, nh int         // number of functions
	var fcn goga.MinProb_t // functions

	// problems. Examples from Deb (2000) An efficient constraint handling method for genetic algorithms
	switch problem {

	// problem # 1 -- Deb's problem 1
	case 1:
		opt.RptName = "1"
		opt.RptFref = []float64{13.59085}
		opt.RptXref = []float64{2.246826, 2.381865}
		opt.FltMin = []float64{0, 0}
		opt.FltMax = []float64{6, 6}
		ng = 2
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = math.Pow(x[0]*x[0]+x[1]-11.0, 2.0) + math.Pow(x[0]+x[1]*x[1]-7.0, 2.0)
			g[0] = 4.84 - math.Pow(x[0]-0.05, 2.0) - math.Pow(x[1]-2.5, 2.0)
			g[1] = x[0]*x[0] + math.Pow(x[1]-2.5, 2.0) - 4.84
		}
		opt.RptFmtFdev = "%.4e"

	// problem # 2 -- Deb's problem 3 -- Z. Michalewicz 1995
	case 2:
		opt.Ncpu = 4
		opt.RptName = "2"
		opt.RptFref = []float64{-15.0}
		opt.RptXref = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 3, 3, 1}
		opt.FltMin = make([]float64, 13)
		opt.FltMax = make([]float64, 13)
		for i := 0; i < 9; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		for i := 9; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 100
		}
		opt.FltMin[12], opt.FltMax[12] = 0, 1
		ng = 9
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			s1, s2, s3 := 0.0, 0.0, 0.0
			for i := 0; i < 4; i++ {
				s1 += x[i]
				s2 += x[i] * x[i]
			}
			for i := 4; i < 13; i++ {
				s3 += x[i]
			}
			f[0] = 5.0*(s1-s2) - s3
			g[0] = 10.0 - 2.0*x[0] - 2.0*x[1] - x[9] - x[10]
			g[1] = 10.0 - 2.0*x[0] - 2.0*x[2] - x[9] - x[11]
			g[2] = 10.0 - 2.0*x[1] - 2.0*x[2] - x[10] - x[11]
			g[3] = 8.0*x[0] - x[9]
			g[4] = 8.0*x[1] - x[10]
			g[5] = 8.0*x[2] - x[11]
			g[6] = 2.0*x[3] + x[4] - x[9]
			g[7] = 2.0*x[5] + x[6] - x[10]
			g[8] = 2.0*x[7] + x[8] - x[11]
		}
		opt.RptFmtFdev = "%.4e"
		opt.RptFmtX = "%.3f"

	// problem # 3 -- Deb's problem 6 -- Z. Michalewicz 1996 -- D.M. Himmelblau 1972
	case 3:
		opt.RptName = "3"
		opt.RptFref = []float64{-30665.5}
		opt.RptXref = []float64{78.0, 33.0, 29.995, 45.0, 36.776}
		opt.FltMin = []float64{78, 33, 27, 27, 27}
		opt.FltMax = []float64{102, 45, 45, 45, 45}
		ng = 6
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c0 := 85.334407 + 0.0056858*x[1]*x[4] + 0.0006262*x[0]*x[3] - 0.0022053*x[2]*x[4]
			c1 := 80.51249 + 0.0071317*x[1]*x[4] + 0.0029955*x[0]*x[1] + 0.0021813*x[2]*x[2]
			c2 := 9.300961 + 0.0047026*x[2]*x[4] + 0.0012547*x[0]*x[2] + 0.0019085*x[2]*x[3]
			f[0] = 5.3578547*x[2]*x[2] + 0.8356891*x[0]*x[4] + 37.293239*x[0] - 40792.141
			g[0] = c0
			g[1] = 92.0 - c0
			g[2] = c1 - 90.0
			g[3] = 110.0 - c1
			g[4] = c2 - 20.0
			g[5] = 25.0 - c2
		}
		opt.RptFmtF = "%.3f"
		opt.RptFmtFdev = "%.3f"

	// problem # 4 -- Deb's problem 2 -- D.M. Himmelblau 1972 -- W. Hock, K. Schittkowski 1981
	case 4:
		opt.RptName = "4"
		opt.RptFref = []float64{-1.90513375}
		opt.RptXref = []float64{705.1803, 68.60005, 102.90001, 282.324999, 37.5850413}
		opt.FltMin = []float64{704.4148, 68.6, 0.0, 193.0, 25.0}
		opt.FltMax = []float64{906.3855, 288.88, 134.75, 287.0966, 84.1988}
		acf := []float64{0, 17.505, 11.275, 214.228, 7.458, 0.961, 1.612, 0.146, 107.99, 922.693, 926.832, 18.766, 1072.163, 8961.448, 0.063, 71084.33, 2802713.0}
		bcf := []float64{0, 1053.6667, 35.03, 665.585, 584.463, 265.916, 7.046, 0.222, 273.366, 1286.105, 1444.046, 537.141, 3247.039, 26844.086, 0.386, 140000.0, 12146108.0}
		ng = 38
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c, y := make([]float64, 18), make([]float64, 18)
			y[0] = x[1] + x[2] + 41.6
			c[0] = 0.024*x[3] - 4.62
			y[1] = 12.5/c[0] + 12.0
			c[1] = 0.0003535*x[0]*x[0] + 0.5311*x[0] + 0.08705*y[1]*x[0]
			c[2] = 0.052*x[0] + 78.0 + 0.002377*y[1]*x[0]
			y[2] = c[1] / c[2]
			y[3] = 19.0 * y[2]
			c[3] = 0.04782*(x[0]-y[2]) + 0.1956*(x[0]-y[2])*(x[0]-y[2])/x[1] + 0.6376*y[3] + 1.594*y[2]
			c[4] = 100.0 * x[1]
			c[5] = x[0] - y[2] - y[3]
			c[6] = 0.95 - c[3]/c[4]
			y[4] = c[5] * c[6]
			y[5] = x[0] - y[4] - y[3] - y[2]
			c[7] = 0.995 * (y[3] + y[4])
			y[6] = c[7] / y[0]
			y[7] = c[7] / 3798.0
			c[8] = y[6] - 0.0663*y[6]/y[7] - 0.3153
			y[8] = 96.82/c[8] + 0.321*y[0]
			y[9] = 1.29*y[4] + 1.258*y[3] + 2.29*y[2] + 1.71*y[5]
			y[10] = 1.71*x[0] - 0.452*y[3] + 0.58*y[2]
			c[9] = 12.3 / 752.3
			c[10] = 1.75 * y[1] * 0.995 * x[0]
			c[11] = 0.995*y[9] + 1998.0
			y[11] = c[9]*x[0] + c[10]/c[11]
			y[12] = c[11] - 1.75*y[1]
			y[13] = 3623.0 + 64.4*x[1] + 58.4*x[2] + 146312.0/(y[8]+x[4])
			c[12] = 0.995*y[9] + 60.8*x[1] + 48.0*x[3] - 0.1121*y[13] - 5095.0
			y[14] = y[12] / c[12]
			y[15] = 148000.0 - 331000.0*y[14] + 40*y[12] - 61.0*y[14]*y[12]
			c[13] = 2324.0*y[9] - 28740000.0*y[1]
			y[16] = 14130000.0 - 1328.0*y[9] - 531.0*y[10] + c[13]/c[11]
			c[14] = y[12]/y[14] - y[12]/0.52
			c[15] = 1.104 - 0.72*y[14]
			c[16] = y[8] + x[4]
			f[0] = -5.843e-7*y[16] + 1.17e-4*y[13] + 2.358e-5*y[12] + 1.502e-6*y[15] + 0.0321*y[11] + 0.004324*y[4] + 1.0e-4*c[14]/c[15] + 37.48*y[1]/c[11] + 0.1365
			g[0] = 1.5*x[1] - x[2]
			g[1] = y[0] - 213.1
			g[2] = 405.23 - y[0]
			for i := 1; i < 17; i++ {
				g[i+2] = y[i] - acf[i]
			}
			for i := 1; i < 17; i++ {
				g[i+18] = bcf[i] - y[i]
			}
			g[35] = y[3] - 0.28/0.72*y[4]
			g[36] = 21.0 - 3496.0*y[1]/c[11]
			g[37] = 62212.0/c[16] - 110.6 - y[0]
		}

	// problem # 5 -- G.V. Reklaitis, A. Ravindran, K.M. Ragsdell 1983 -- Coello 2002
	case 5:
		opt.RptName = "5"
		opt.RptFref = []float64{2.38116}
		opt.RptXref = []float64{0.2444, 6.2187, 8.2915, 0.2444}
		opt.FltMin = []float64{0.125, 0.1, 0.1, 0.1}
		opt.FltMax = []float64{10.0, 10.0, 10.0, 10.0}
		ng = 5
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			E := 30e+6
			G := 12e+6
			P := 6000.0
			L := 14.0
			τmax := 13600.0
			σmax := 30000.0
			δmax := 0.25
			SQ2 := math.Sqrt2
			R := math.Sqrt(x[1]*x[1]/4.0 + math.Pow((x[0]+x[2])/2.0, 2.0))
			J := 2.0 * SQ2 * x[0] * x[1] * (x[1]*x[1]/12.0 + math.Pow((x[0]+x[2])/2.0, 2.0))
			M := P * (L + x[1]/2.0)
			τd := P / (SQ2 * x[0] * x[1])
			τdd := M * R / J
			τ := math.Sqrt(τd*τd + τd*τdd*x[1]/R + τdd*τdd)
			LL := L * L
			x2two := x[2] * x[2]
			x3six := math.Pow(x[3], 6.0)
			σ := 6.0 * P * L / (x2two * x[3])
			δ := 4.0 * P * L * LL / (E * x2two * x[2] * x[3])
			Pc := 4.013 * E * math.Sqrt(x2two*x3six/36.0) * (1.0 - 0.5*x[2]*math.Sqrt(E/(4.0*G))/L) / LL
			f[0] = 1.10471*x[0]*x[0]*x[1] + 0.04811*x[2]*x[3]*(14.0+x[1])
			g[0] = x[3] - x[0]
			g[1] = Pc - P
			g[2] = τmax - τ
			g[3] = σmax - σ
			g[4] = δmax - δ
		}

	// problem # 6 -- Deb's problem 5 -- Z. Michalewicz 1995
	case 6:
		opt.Ncpu = 2
		opt.RptName = "6"
		opt.RptFref = []float64{680.6300573}
		opt.RptXref = []float64{2.330499, 1.951372, -0.4775414, 4.365726, -0.6244870, 1.038131, 1.594227}
		opt.FltMin = make([]float64, 7)
		opt.FltMax = make([]float64, 7)
		for i := 0; i < 7; i++ {
			opt.FltMin[i], opt.FltMax[i] = -10, 10
		}
		ng = 4
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = math.Pow(x[0]-10.0, 2.0) + 5.0*math.Pow(x[1]-12.0, 2.0) + math.Pow(x[2], 4.0) + 3.0*math.Pow(x[3]-11.0, 2.0) + 10.0*math.Pow(x[4], 6.0) + 7.0*math.Pow(x[5], 2.0) + math.Pow(x[6], 4.0) - 4.0*x[5]*x[6] - 10.0*x[5] - 8.0*x[6]
			g[0] = 127.0 - 2.0*x[0]*x[0] - 3.0*math.Pow(x[1], 4.0) - x[2] - 4.0*x[3]*x[3] - 5.0*x[4]
			g[1] = 282.0 - 7.0*x[0] - 3.0*x[1] - 10.0*x[2]*x[2] - x[3] + x[4]
			g[2] = 196.0 - 23.0*x[0] - x[1]*x[1] - 6.0*x[5]*x[5] + 8.0*x[6]
			g[3] = -4.0*x[0]*x[0] - x[1]*x[1] + 3.0*x[0]*x[1] - 2.0*x[2]*x[2] - 5.0*x[5] + 11.0*x[6]
		}

	// problem # 7 -- Deb's problem 8 -- Z. Michalewicz 1995
	case 7:
		opt.Ncpu = 2
		opt.RptName = "7"
		opt.RptFref = []float64{24.3062091}
		opt.RptXref = []float64{2.171996, 2.363683, 8.773926, 5.095984, 0.9906548, 1.430574, 1.321644, 9.828726, 8.280092, 8.375927}
		opt.FltMin = make([]float64, 10)
		opt.FltMax = make([]float64, 10)
		for i := 0; i < 10; i++ {
			opt.FltMin[i], opt.FltMax[i] = -10, 10
		}
		ng = 8
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = x[0]*x[0] + x[1]*x[1] + x[0]*x[1] - 14.0*x[0] - 16.0*x[1] + math.Pow(x[2]-10.0, 2.0) + 4.0*math.Pow(x[3]-5.0, 2.0) + math.Pow(x[4]-3.0, 2.0) + 2.0*math.Pow(x[5]-1.0, 2.0) + 5.0*x[6]*x[6] + 7.0*math.Pow(x[7]-11.0, 2.0) + 2.0*math.Pow(x[8]-10.0, 2.0) + math.Pow(x[9]-7.0, 2.0) + 45.0
			g[0] = 105.0 - 4.0*x[0] - 5.0*x[1] + 3.0*x[6] - 9.0*x[7]
			g[1] = -10.0*x[0] + 8.0*x[1] + 17.0*x[6] - 2.0*x[7]
			g[2] = 8.0*x[0] - 2.0*x[1] - 5.0*x[8] + 2.0*x[9] + 12.0
			g[3] = -3.0*math.Pow(x[0]-2, 2.0) - 4.0*math.Pow(x[1]-3.0, 2.0) - 2.0*x[2]*x[2] + 7.0*x[3] + 120.0
			g[4] = -5.0*x[0]*x[0] - 8.0*x[1] - math.Pow(x[2]-6.0, 2.0) + 2.0*x[3] + 40.0
			g[5] = -x[0]*x[0] - 2.0*math.Pow(x[1]-2.0, 2.0) + 2.0*x[0]*x[1] - 14.0*x[4] + 6.0*x[5]
			g[6] = -0.5*math.Pow(x[0]-8.0, 2.0) - 2.0*math.Pow(x[1]-4.0, 2.0) - 3.0*x[4]*x[4] + x[5] + 30.0
			g[7] = 3.0*x[0] - 6.0*x[1] - 12.0*math.Pow(x[8]-8.0, 2.0) + 7.0*x[9]
		}
		opt.RptFmtX = "%.4f"

	// problem # 8 -- Deb's problem 4 -- Z. Michalewicz 1995
	case 8:
		opt.Ncpu = 2
		opt.RptName = "8"
		opt.RptFref = []float64{7049.330923}
		opt.RptXref = []float64{579.3167, 1359.943, 5110.071, 182.0174, 295.5985, 217.9799, 286.4162, 395.5979}
		opt.FltMin = make([]float64, 8)
		opt.FltMax = make([]float64, 8)
		opt.FltMin[0], opt.FltMax[0] = 100, 10000
		opt.FltMin[1], opt.FltMax[1] = 1000, 10000
		opt.FltMin[2], opt.FltMax[2] = 1000, 10000
		for i := 3; i < 8; i++ {
			opt.FltMin[i], opt.FltMax[i] = 10, 1000
		}
		ng = 6
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = x[0] + x[1] + x[2]
			g[0] = 1.0 - 0.0025*(x[3]+x[5])
			g[1] = 1.0 - 0.0025*(x[4]+x[6]-x[3])
			g[2] = 1.0 - 0.01*(x[7]-x[4])
			g[3] = x[0]*x[5] - 833.33252*x[3] - 100.0*x[0] + 83333.333
			g[4] = x[1]*x[6] - 1250.0*x[4] - x[1]*x[3] + 1250.0*x[3]
			g[5] = x[2]*x[7] - x[2]*x[4] + 2500.0*x[4] - 1250000
		}
		opt.RptFmtX = "%.3f"

	// problem # 9 -- Deb's problem 7 -- Z. Michalewicz 1995
	case 9:
		opt.RptName = "9"
		opt.RptFref = []float64{0.0539498478}
		opt.RptXref = []float64{-1.717143, 1.595709, 1.827247, -0.7636413, -0.7636450}
		opt.FltMin = []float64{-2.3, -2.3, -3.2, -3.2, -3.2}
		opt.FltMax = []float64{+2.3, +2.3, +3.2, +3.2, +3.2}
		nh = 3
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = math.Exp(x[0] * x[1] * x[2] * x[3] * x[4])
			h[0] = x[0]*x[0] + x[1]*x[1] + x[2]*x[2] + x[3]*x[3] + x[4]*x[4] - 10.0
			h[1] = x[1]*x[2] - 5.0*x[3]*x[4]
			h[2] = math.Pow(x[0], 3.0) + math.Pow(x[1], 3.0) + 1.0
		}

	default:
		chk.Panic("problem %d is not available", problem)
	}

	// check best known solution
	if checkOnly {
		check(fcn, ng, nh, opt.RptXref, opt.RptFref[0], 1e-6)
		return
	}

	// number of trial solutions
	opt.Nsol = len(opt.FltMin) * 10

	// initialise optimiser
	nf := 1
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	//opt.RunMany("/tmp/goga", "functions")
	opt.RunMany("", "")
	goga.StatF(opt, 0, true)
	return
}

func main() {
	ntrials := 100
	checkOnly := false
	P := utl.IntRange2(1, 10)
	//P := []int{5}
	//P := []int{4}
	//P := []int{7}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem, ntrials, checkOnly)
	}
	if !checkOnly {
		io.Pf("\n-------------------------- generating report --------------------------\nn")
		nRowPerTab := 5
		goga.TexSingleObjReport("/tmp/goga", "functions", nRowPerTab, opts)
		io.Pf("\n%v\n", opts[0].LogParams())
	}
}

func check(fcn goga.MinProb_t, ng, nh int, xs []float64, fs, ϵ float64) {
	f := make([]float64, 1)
	g := make([]float64, ng)
	h := make([]float64, nh)
	cpu := 0
	fcn(f, g, h, xs, nil, cpu)
	io.Pfblue2("xs = %v\n", xs)
	io.Pfblue2("f(x)=%g  (%g)  diff=%g\n", f[0], fs, math.Abs(fs-f[0]))
	for i, v := range g {
		unfeasible := false
		if v < 0 {
			unfeasible = true
		}
		if unfeasible {
			io.Pfred("g%d(x) = %g\n", i, v)
		} else {
			io.Pfgreen("g%d(x) = %g\n", i, v)
		}
	}
	for i, v := range h {
		unfeasible := false
		if math.Abs(v) > ϵ {
			unfeasible = true
		}
		if unfeasible {
			io.Pfred("h%d(x) = %g\n", i, v)
		} else {
			io.Pfgreen("h%d(x) = %g\n", i, v)
		}
	}
}

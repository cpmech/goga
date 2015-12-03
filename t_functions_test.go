// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_functions(tst *testing.T) {

	verbose()
	chk.PrintTitle("simple functions")

	// GA parameters
	var opt Optimiser
	opt.Default()
	opt.Ncpu = 1
	opt.Verbose = false
	opt.Problem = 2
	opt.Ntrials = 100
	opt.DEpc = 0.5
	opt.DEmult = 0.5
	//opt.PmFlt = 1
	//opt.DebEtam = 1

	// problem variables
	var ng, nh int    // number of functions
	var fs float64    // best known solution
	var xs []float64  // best known solution
	var fcn MinProb_t // functions

	// problems
	switch opt.Problem {

	// problem # 1
	case 1:
		opt.Nsol = 20
		opt.Tf = 100
		opt.FltMin = []float64{0, 0}
		opt.FltMax = []float64{6, 6}
		xs = []float64{2.246826, 2.381865}
		fs = 13.59085
		ng = 2
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = math.Pow(x[0]*x[0]+x[1]-11.0, 2.0) + math.Pow(x[0]+x[1]*x[1]-7.0, 2.0)
			g[0] = 4.84 - math.Pow(x[0]-0.05, 2.0) - math.Pow(x[1]-2.5, 2.0)
			g[1] = x[0]*x[0] + math.Pow(x[1]-2.5, 2.0) - 4.84
		}

	// problem # 2
	case 2:
		opt.Ncpu = 1
		opt.Nsol = 50
		opt.Tf = 200
		opt.DtExc = opt.Tf / 20
		opt.FltMin = []float64{704.4148, 68.6, 0.0, 193.0, 25.0}
		opt.FltMax = []float64{906.3855, 288.88, 134.75, 287.0966, 84.1988}
		xs = []float64{705.1803, 68.60005, 102.90001, 282.324999, 37.5850413}
		fs = -1.90513375
		ng = 38
		acf := []float64{0, 17.505, 11.275, 214.228, 7.458, 0.961, 1.612, 0.146, 107.99, 922.693, 926.832, 18.766, 1072.163, 8961.448, 0.063, 71084.33, 2802713.0}
		bcf := []float64{0, 1053.6667, 35.03, 665.585, 584.463, 265.916, 7.046, 0.222, 273.366, 1286.105, 1444.046, 537.141, 3247.039, 26844.086, 0.386, 140000.0, 12146108.0}
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

	// problem # 3
	case 3:
		opt.Nsol = 24
		opt.Tf = 500
		opt.FltMin = make([]float64, 13)
		opt.FltMax = make([]float64, 13)
		for i := 0; i < 9; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		for i := 9; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 100
		}
		opt.FltMin[12], opt.FltMax[12] = 0, 1
		xs = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 3, 3, 1}
		fs = -15.0
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

	// problem # 4
	case 4:
		opt.Nsol = 24
		opt.Tf = 1000
		opt.FltMin = make([]float64, 8)
		opt.FltMax = make([]float64, 8)
		opt.FltMin[0], opt.FltMax[0] = 100, 10000
		opt.FltMin[1], opt.FltMax[1] = 1000, 10000
		opt.FltMin[2], opt.FltMax[2] = 1000, 10000
		for i := 3; i < 8; i++ {
			opt.FltMin[i], opt.FltMax[i] = 10, 1000
		}
		xs = []float64{579.3167, 1359.943, 5110.071, 182.0174, 295.5985, 217.9799, 286.4162, 395.5979}
		fs = 7049.330923
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

	// problem # 5
	case 5:
		opt.Nsol = 24
		opt.FltMin = make([]float64, 7)
		opt.FltMax = make([]float64, 7)
		for i := 0; i < 7; i++ {
			opt.FltMin[i], opt.FltMax[i] = -10, 10
		}
		xs = []float64{2.330499, 1.951372, -0.4775414, 4.365726, -0.6244870, 1.038131, 1.594227}
		fs = 680.6300573
		ng = 4
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = math.Pow(x[0]-10.0, 2.0) + 5.0*math.Pow(x[1]-12.0, 2.0) + math.Pow(x[2], 4.0) + 3.0*math.Pow(x[3]-11.0, 2.0) +
				10.0*math.Pow(x[4], 6.0) + 7.0*math.Pow(x[5], 2.0) + math.Pow(x[6], 4.0) - 4.0*x[5]*x[6] - 10.0*x[5] - 8.0*x[6]
			g[0] = 127.0 - 2.0*x[0]*x[0] - 3.0*math.Pow(x[1], 4.0) - x[2] - 4.0*x[3]*x[3] - 5.0*x[4]
			g[1] = 282.0 - 7.0*x[0] - 3.0*x[1] - 10.0*x[2]*x[2] - x[3] + x[4]
			g[2] = 196.0 - 23.0*x[0] - x[1]*x[1] - 6.0*x[5]*x[5] + 8.0*x[6]
			g[3] = -4.0*x[0]*x[0] - x[1]*x[1] + 3.0*x[0]*x[1] - 2.0*x[2]*x[2] - 5.0*x[5] + 11.0*x[6]
		}

	// problem # 6. also Coelho² 2002: Himmelblau's problem
	case 6:
		opt.Nsol = 12
		opt.FltMin = []float64{78, 33, 27, 27, 27}
		opt.FltMax = []float64{102, 45, 45, 45, 45}
		xs = []float64{78.0, 33.0, 29.995, 45.0, 36.776}
		fs = -30665.5
		ng = 6
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = 5.3578547*x[2]*x[2] + 0.8356891*x[0]*x[4] + 37.293239*x[0] - 40792.141
			g[0] = 85.334407 + 0.0056858*x[1]*x[4] + 0.0006262*x[0]*x[3] - 0.0022053*x[2]*x[4]
			g[1] = 92.0 - (85.334407 + 0.0056858*x[1]*x[4] + 0.0006262*x[0]*x[3] - 0.0022053*x[2]*x[4])
			g[2] = 80.51249 + 0.0071317*x[1]*x[4] + 0.0029955*x[0]*x[1] + 0.0021813*x[2]*x[2] - 90.0
			g[3] = 110.0 - (80.51249 + 0.0071317*x[1]*x[4] + 0.0029955*x[0]*x[1] + 0.0021813*x[2]*x[2])
			g[4] = 9.300961 + 0.0047026*x[2]*x[4] + 0.0012547*x[0]*x[2] + 0.0019085*x[2]*x[3] - 20.0
			g[5] = 25.0 - (9.300961 + 0.0047026*x[2]*x[4] + 0.0012547*x[0]*x[2] + 0.0019085*x[2]*x[3])
		}

	// problem # 7: Michaelwicz (1996) page 146
	case 7:
		opt.Nsol = 24
		opt.Tf = 2000
		opt.FltMin = []float64{-2.3, -2.3, -3.2, -3.2, -3.2}
		opt.FltMax = []float64{+2.3, +2.3, +3.2, +3.2, +3.2}
		xs = []float64{-1.717143, 1.595709, 1.827247, -0.7636413, -0.7636450}
		fs = 0.0539498478
		nh = 3
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			f[0] = math.Exp(x[0] * x[1] * x[2] * x[3] * x[4])
			h[0] = x[0]*x[0] + x[1]*x[1] + x[2]*x[2] + x[3]*x[3] + x[4]*x[4] - 10.0
			h[1] = x[1]*x[2] - 5.0*x[3]*x[4]
			h[2] = math.Pow(x[0], 3.0) + math.Pow(x[1], 3.0) + 1.0
		}

	// problem # 8
	case 8:
		opt.Nsol = 24
		opt.FltMin = make([]float64, 10)
		opt.FltMax = make([]float64, 10)
		for i := 0; i < 10; i++ {
			opt.FltMin[i], opt.FltMax[i] = -10, 10
		}
		xs = []float64{2.171996, 2.363683, 8.773926, 5.095984, 0.9906548, 1.430574, 1.321644, 9.828726, 8.280092, 8.375927}
		fs = 24.3062091
		ng = 8
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			x1, x2, x3, x4, x5, x6, x7, x8, x9, x10 := x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7], x[8], x[9]
			f[0] = x1*x1 + x2*x2 + x1*x2 - 14.0*x1 - 16.0*x2 + math.Pow(x3-10.0, 2.0) + 4.0*math.Pow(x4-5.0, 2.0) +
				math.Pow(x5-3.0, 2.0) + 2.0*math.Pow(x6-1.0, 2.0) + 5.0*x7*x7 + 7.0*math.Pow(x8-11.0, 2.0) +
				2.0*math.Pow(x9-10.0, 2.0) + math.Pow(x10-7.0, 2.0) + 45.0
			g[0] = 105.0 - 4.0*x1 - 5.0*x2 + 3.0*x7 - 9.0*x8
			g[1] = -10.0*x1 + 8.0*x2 + 17.0*x7 - 2.0*x8
			g[2] = 8.0*x1 - 2.0*x2 - 5.0*x9 + 2.0*x10 + 12.0
			g[3] = -3.0*math.Pow(x1-2, 2.0) - 4.0*math.Pow(x2-3.0, 2.0) - 2.0*x3*x3 + 7.0*x4 + 120.0
			g[4] = -5.0*x1*x1 - 8.0*x2 - math.Pow(x3-6.0, 2.0) + 2.0*x4 + 40.0
			g[5] = -x1*x1 - 2.0*math.Pow(x2-2.0, 2.0) + 2.0*x1*x2 - 14.0*x5 + 6.0*x6
			g[6] = -0.5*math.Pow(x1-8.0, 2.0) - 2.0*math.Pow(x2-4.0, 2.0) - 3.0*x5*x5 + x6 + 30.0
			g[7] = 3.0*x1 - 6.0*x2 - 12.0*math.Pow(x9-8.0, 2.0) + 7.0*x10
		}

	// problem # 9: welded beam design. Coello² (2002)
	case 9:
		opt.Nsol = 12
		opt.FltMin = []float64{0.125, 0.1, 0.1, 0.1}
		opt.FltMax = []float64{10.0, 10.0, 10.0, 10.0}
		xs = []float64{0.2444, 6.2187, 8.2915, 0.2444}
		fs = 2.38116
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

	default:
		chk.Panic("problem %d is not available", opt.Problem)
	}

	// check best known solution
	if false {
		check(fcn, ng, nh, xs, fs, 1e-6)
	}

	// initialise optimiser
	nf := 1
	opt.Init(GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.RunMany()
	opt.StatMinProb(0, 60, fs, true)

	// results
	SortByOva(opt.Solutions, 0)
	best := opt.Solutions[0]
	io.Pfyel("     xs = %.6f\n", xs)
	io.Pfyel("     fs = %.6f\n", fs)
	io.Pf("best x  = %.6f\n", best.Flt)
	io.Pf("best f  = %.6f\n", best.Ova[0])
}

func check(fcn MinProb_t, ng, nh int, xs []float64, fs, ϵ float64) {
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

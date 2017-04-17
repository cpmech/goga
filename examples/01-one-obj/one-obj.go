// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// main function
func main() {

	// problem numbers
	P := utl.IntRange2(1, 10)
	//P := []int{1}

	// check problem constraints
	checkOnly := false

	// allocate and run each problem
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = oneObj(problem, checkOnly)
	}

	// skip report if just checking constraints
	if checkOnly {
		return
	}

	// report
	io.Pf("\n----------------------------------- generating report -----------------------------------\n\n")
	rpt := goga.NewTexReport(opts)
	rpt.UseGeom = true
	rpt.Landscape = false
	rpt.DescHeader = "ref"
	rpt.SetColumnsSingleObj(false, false)
	rpt.Title = "Constrained single-objective problems"
	rpt.Generate("/tmp/goga", "one-obj")
}

// oneObj runs one-obj problem
func oneObj(problem int, checkOnly bool) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// options
	constantSeed := false

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Ncpu = 1
	opt.Tmax = 500
	opt.Verbose = false
	opt.VerbStat = false
	opt.GenType = "latin"
	opt.Nsamples = 3 /////////////////////// increase this number

	// options for report
	opt.RptFmtF = "%.5f"
	opt.RptFmtFdev = "%.5f"
	opt.RptFmtX = "%.5f"

	// problem variables
	var ng, nh int         // number of functions
	var fcn goga.MinProb_t // functions

	// problems. Examples from Deb (2000) An efficient constraint handling method for genetic algorithms
	switch problem {

	// problem # 1 -- Deb's problem 1
	//   is a simple problem presented in \citep{deb:00} with 2 variables. It has 2 inequalities.
	case 1:
		opt.RptName = "1"
		opt.RptFref = []float64{13.59085}
		opt.RptXref = []float64{2.246826, 2.381865}
		opt.FltMin = []float64{0, 0}
		opt.FltMax = []float64{6, 6}
		ng = 2
		fcn = func(f, g, h, x []float64, y []int, cpu int) {
			f[0] = math.Pow(x[0]*x[0]+x[1]-11.0, 2.0) + math.Pow(x[0]+x[1]*x[1]-7.0, 2.0)
			g[0] = 4.84 - math.Pow(x[0]-0.05, 2.0) - math.Pow(x[1]-2.5, 2.0)
			g[1] = x[0]*x[0] + math.Pow(x[1]-2.5, 2.0) - 4.84
		}
		opt.RptFmtFdev = "%.2e"
		opt.RptDesc = `\cite{deb:00}-T1`

	// problem # 2 -- Deb's problem 3 -- Z. Michalewicz 1995
	//   has 13 variables, 9 inequalities and is relatively easy because it involves only quadratic
	//   terms in the objective function and has linear constraints. It corresponds to
	//   \emph{case 1} in \citep{mich:95} and
	//   \emph{test 3} in \citep{deb:00}
	case 2:
		opt.Ncpu = 4
		opt.RptName = "2"
		opt.RptFref = []float64{-15.0}
		opt.RptXref = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 3, 3, 1}
		opt.FltMin = make([]float64, 13)
		opt.FltMax = make([]float64, 13)
		strategy2 := true // enlarge box; add more constraint equations
		xmin, xmax := 0.0, 1.0
		if strategy2 {
			xmin, xmax = -0.5, 1.5
		}
		for i := 0; i < 9; i++ {
			opt.FltMin[i], opt.FltMax[i] = xmin, xmax
		}
		opt.FltMin[12], opt.FltMax[12] = xmin, xmax
		xmin, xmax = 0, 100
		if strategy2 {
			xmin, xmax = -1, 101
		}
		for i := 9; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = xmin, xmax
		}
		ng = 9
		if strategy2 {
			ng += 9 + 9 + 3 + 3 + 2
		}
		fcn = func(f, g, h, x []float64, y []int, cpu int) {
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
			if strategy2 {
				for i := 0; i < 9; i++ {
					g[9+i] = x[i]
					g[18+i] = 1.0 - x[i]
				}
				for i := 0; i < 3; i++ {
					g[27+i] = x[9+i]
					g[30+i] = 100.0 - x[9+i]
				}
				g[33] = x[12]
				g[34] = 1.0 - x[12]
			}
		}
		opt.RptFmtFdev = "%.2e"
		opt.RptFmtX = "%.3f"
		opt.RptDesc = `\cite{deb:00}-T3`

	// problem # 3 -- Deb's problem 6 -- Z. Michalewicz 1996 -- D.M. Himmelblau 1972
	//   has 5 variables, 6 inequalities and has normal difficulty since it involves a quadratic
	//   objective function and has quadratic constraint functions. It corresponds to
	//   \emph{problem 83 (page 102)} in \citep{hock:81} and
	//   \emph{test 6} in \citep{deb:00}
	case 3:
		opt.RptName = "3"
		opt.RptFref = []float64{-30665.5}
		opt.RptXref = []float64{78.0, 33.0, 29.995, 45.0, 36.776}
		opt.FltMin = []float64{78, 33, 27, 27, 27}
		opt.FltMax = []float64{102, 45, 45, 45, 45}
		ng = 6
		fcn = func(f, g, h, x []float64, y []int, cpu int) {
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
		opt.RptFmtF = "%.2f"
		opt.RptFmtFdev = "%.2e"
		opt.RptDesc = `\cite{deb:00}-T6`

	// problem # 4 -- Deb's problem 2 -- D.M. Himmelblau 1972 -- W. Hock, K. Schittkowski 1981
	//   has 5 variables, 38 inequalities and is of moderate difficulty because it involves
	//   a large number of nonlinear constraints up to the second order. It corresponds to
	//   \emph{{problem~85} (page 104)} in \citep{hock:81} and
	//   \emph{test 2} in \citep{deb:00}.
	case 4:
		opt.RptName = "4"
		opt.RptFref = []float64{-1.90513375}
		opt.RptXref = []float64{705.1803, 68.60005, 102.90001, 282.324999, 37.5850413}
		opt.FltMin = []float64{704.4148, 68.6, 0.0, 193.0, 25.0}
		opt.FltMax = []float64{906.3855, 288.88, 134.75, 287.0966, 84.1988}
		acf := []float64{0, 17.505, 11.275, 214.228, 7.458, 0.961, 1.612, 0.146, 107.99, 922.693, 926.832, 18.766, 1072.163, 8961.448, 0.063, 71084.33, 2802713.0}
		bcf := []float64{0, 1053.6667, 35.03, 665.585, 584.463, 265.916, 7.046, 0.222, 273.366, 1286.105, 1444.046, 537.141, 3247.039, 26844.086, 0.386, 140000.0, 12146108.0}
		ng = 38
		fcn = func(f, g, h, x []float64, ints []int, cpu int) {
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
		opt.RptFmtFdev = "%.2e"
		opt.RptDesc = `\cite{deb:00}-T2`

	// problem # 5 -- A. Ravindran, K. M. Ragsdell, and G. V. Reklaitis (2007)
	//   has 4 variables, 5 inequalities and is of moderate difficulty because of the cubic
	//   expressions in the objective and constraint functions. It corresponds to
	//   the design of a welded beam firstly presented in \citep{rags:76} (see also \citep{ravi:07})
	//   where the objective function is the system cost. It corresponds to
	//   the last problem in \cite{deb:00}.
	//   As shown in \citep{ravi:07} (page 592), after substitutions, the problem is defined by
	case 5:
		opt.RptName = "5"
		opt.RptFref = []float64{2.34027}
		opt.RptXref = []float64{0.2536, 7.141, 7.1044, 0.2536}
		//opt.RptFref = []float64{28.0671}
		//opt.RptXref = []float64{2, 4, 4, 3}
		opt.FltMin = []float64{0.125, 0.0, 0.0, 0.125}
		opt.FltMax = []float64{10.0, 10.0, 10.0, 10.0}
		ng = 5
		fcn = func(f, g, h, x []float64, y []int, cpu int) {
			F, τd := 6000.0, 13600.0
			c0 := x[0] * x[0] * x[1] * x[1]
			c1 := x[1]*x[1] + 3.0*math.Pow(x[0]+x[2], 2.0)
			f[0] = 1.10471*x[0]*x[0]*x[1] + 0.04811*x[2]*x[3]*(14.0+x[1])
			g[0] = (τd / F) - math.Sqrt(1.0/(2.0*c0)+3.0*(28.0+x[1])/(x[0]*x[0]*x[1]*c1)+4.5*math.Pow(28.0+x[1], 2.0)*(x[1]*x[1]+math.Pow(x[0]+x[2], 2.0))/(c0*c1*c1))
			g[1] = x[2]*x[2]*x[3] - 12.8
			g[2] = x[3] - x[0]
			g[3] = x[2]*math.Pow(x[3], 3.0)*(1.0-0.02823*x[2]) - 0.09267
			g[4] = math.Pow(x[2], 3.0)*x[3] - 8.7808
		}
		opt.RptFmtFdev = "%.2e"
		opt.RptDesc = `\cite{deb:00}-T9`

	// problem # 6 -- Deb's problem 5 -- Z. Michalewicz 1995
	//   has 7 variables, 4 inequalities and is of moderate difficulty with a fourth order term
	//   in the objective function and nonlinear constraints. It corresponds to
	//   \emph{problem 100 (page 111)} in \citep{hock:81},
	//   \emph{case 3} in \citep{mich:95} and
	//   \emph{test 5} in \citep{deb:00}.
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
		fcn = func(f, g, h, x []float64, y []int, cpu int) {
			f[0] = math.Pow(x[0]-10.0, 2.0) + 5.0*math.Pow(x[1]-12.0, 2.0) + math.Pow(x[2], 4.0) + 3.0*math.Pow(x[3]-11.0, 2.0) + 10.0*math.Pow(x[4], 6.0) + 7.0*math.Pow(x[5], 2.0) + math.Pow(x[6], 4.0) - 4.0*x[5]*x[6] - 10.0*x[5] - 8.0*x[6]
			g[0] = 127.0 - 2.0*x[0]*x[0] - 3.0*math.Pow(x[1], 4.0) - x[2] - 4.0*x[3]*x[3] - 5.0*x[4]
			g[1] = 282.0 - 7.0*x[0] - 3.0*x[1] - 10.0*x[2]*x[2] - x[3] + x[4]
			g[2] = 196.0 - 23.0*x[0] - x[1]*x[1] - 6.0*x[5]*x[5] + 8.0*x[6]
			g[3] = -4.0*x[0]*x[0] - x[1]*x[1] + 3.0*x[0]*x[1] - 2.0*x[2]*x[2] - 5.0*x[5] + 11.0*x[6]
		}
		opt.RptFmtFdev = "%.3e"
		opt.RptDesc = `\cite{deb:00}-T5`

	// problem # 7 -- Deb's problem 8 -- Z. Michalewicz 1995
	//   has 10 variables, 8 inequalities and is also of moderate difficulty with quadratic terms
	//   and nonlinear constraints. It corresponds to
	//   \emph{problem 113 (page 122)} in \citep{hock:81},
	//   \emph{case 5} in \citep{mich:95} and
	//   \emph{test 8} in \citep{deb:00}.
	case 7:
		opt.Tmax = 1000
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
		fcn = func(f, g, h, x []float64, y []int, cpu int) {
			f[0] = x[0]*x[0] + x[1]*x[1] + x[0]*x[1] - 14.0*x[0] - 16.0*x[1] + math.Pow(x[2]-10.0, 2.0) + 4.0*math.Pow(x[3]-5.0, 2.0) + math.Pow(x[4]-3.0, 2.0) + 2.0*math.Pow(x[5]-1.0, 2.0) + 5.0*x[6]*x[6] + 7.0*math.Pow(x[7]-11.0, 2.0) + 2.0*math.Pow(x[8]-10.0, 2.0) + math.Pow(x[9]-7.0, 2.0) + 45.0
			g[0] = 105.0 - 4.0*x[0] - 5.0*x[1] + 3.0*x[6] - 9.0*x[7]
			g[1] = -10.0*x[0] + 8.0*x[1] + 17.0*x[6] - 2.0*x[7]
			g[2] = 8.0*x[0] - 2.0*x[1] - 5.0*x[8] + 2.0*x[9] + 12.0
			g[3] = -3.0*math.Pow(x[0]-2, 2.0) - 4.0*math.Pow(x[1]-3.0, 2.0) - 2.0*x[2]*x[2] + 7.0*x[3] + 120.0
			g[4] = -5.0*x[0]*x[0] - 8.0*x[1] - math.Pow(x[2]-6.0, 2.0) + 2.0*x[3] + 40.0
			g[5] = -0.5*math.Pow(x[0]-8.0, 2.0) - 2.0*math.Pow(x[1]-4.0, 2.0) - 3.0*x[4]*x[4] + x[5] + 30.0
			g[6] = -x[0]*x[0] - 2.0*math.Pow(x[1]-2.0, 2.0) + 2.0*x[0]*x[1] - 14.0*x[4] + 6.0*x[5]
			g[7] = 3.0*x[0] - 6.0*x[1] - 12.0*math.Pow(x[8]-8.0, 2.0) + 7.0*x[9]
		}
		opt.RptFmtFdev = "%.2e"
		opt.RptDesc = `\cite{deb:00}-T8`

	// problem # 8 -- Deb's problem 4 -- Z. Michalewicz 1995
	//   has 8 variables, 6 inequalities and is a difficult problem as observed in
	//   \citep{mich:96,deb:00}. It has a linear objective function and nonlinear constraints and
	//   corresponds to
	//   \emph{problem 106 (page 115; heat exchanger design)} in \citep{hock:81},
	//   \emph{case 2} in \citep{mich:95} and
	//   \emph{test 4} in \citep{deb:00}.
	case 8:
		opt.Tmax = 5000
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
		fcn = func(f, g, h, x []float64, y []int, cpu int) {
			f[0] = x[0] + x[1] + x[2]
			g[0] = 1.0 - 0.0025*(x[3]+x[5])
			g[1] = 1.0 - 0.0025*(x[4]+x[6]-x[3])
			g[2] = 1.0 - 0.01*(x[7]-x[4])
			g[3] = x[0]*x[5] - 833.33252*x[3] - 100.0*x[0] + 83333.333
			g[4] = x[1]*x[6] - 1250.0*x[4] - x[1]*x[3] + 1250.0*x[3]
			g[5] = x[2]*x[7] - x[2]*x[4] + 2500.0*x[4] - 1250000
		}
		opt.RptFmtFdev = "%.2e"
		opt.RptDesc = `\cite{deb:00}-T4`

	// problem # 9 -- Deb's problem 7 -- Z. Michalewicz 1995
	//   has 5 variables, 3 equality constraints and is a very difficult problem.
	//   The objective function is the exponential of all variables multiplied together.
	//   The equality constraints are nonlinear (quadratic and cubic) expressions.
	//   The problem corresponds to
	//   \emph{problem 80 (page 100)} in \citep{hock:81},
	//   \emph{case 4} in \citep{mich:95} (page 146) and
	//   \emph{test 7} in \citep{deb:00}.
	case 9:
		opt.Tmax = 7000
		opt.EpsH = 1e-3
		opt.RptName = "9"
		opt.RptFref = []float64{0.0539498478}
		opt.RptXref = []float64{-1.717143, 1.595709, 1.827247, -0.7636413, -0.7636450}
		opt.FltMin = []float64{-2.3, -2.3, -3.2, -3.2, -3.2}
		opt.FltMax = []float64{+2.3, +2.3, +3.2, +3.2, +3.2}
		nh = 3
		fcn = func(f, g, h, x []float64, y []int, cpu int) {
			f[0] = math.Exp(x[0] * x[1] * x[2] * x[3] * x[4])
			h[0] = x[0]*x[0] + x[1]*x[1] + x[2]*x[2] + x[3]*x[3] + x[4]*x[4] - 10.0
			h[1] = x[1]*x[2] - 5.0*x[3]*x[4]
			h[2] = math.Pow(x[0], 3.0) + math.Pow(x[1], 3.0) + 1.0
		}
		opt.RptFmtFdev = "%.2e"
		opt.RptDesc = `\cite{deb:00}-T7`

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
	opt.RunMany("", "", constantSeed)
	opt.PrintStatF(0)
	io.PfMag("Tsys{tot} = %v\n", opt.SysTimeTot)
	io.PfYel("Tsys{ave} = %v\n", opt.SysTimeAve)

	// check
	goga.CheckFront0(opt, true)
	return
}

// check checks problem constraints
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

// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"math"
	"time"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

const TOLMINLOG = 1e-15

func main() {

	// Problems 1 to 5 have all variables as standard variables => μ=0 and σ=1 => y = x
	//
	// References
	//  [1] Santos SR, Matioli LC and Beck AT. New optimization algorithms for structural
	//      reliability analysis. Computer Modeling in Engineering & Sciences, 83(1):23-56; 2012
	//      doi:10.3970/cmes.2012.083.023
	//  [2] Borri A and Speranzini E. Structural reliability analysis using a standard deterministic
	//      finite element code. Structural Safety, 19(4):361-382; 1997
	//      doi:10.1016/S0167-4730(97)00017-9
	//  [3] Grooteman F.  Adaptive radial-based importance sampling method or structural
	//      reliability. Structural safety, 30:533-542; 2008
	//      doi:10.1016/j.strusafe.2007.10.002
	//  [4] Wang L and Grandhi RV. Higher-order failure probability calculation using nonlinear
	//      approximations. Computer Methods in Applied Mechanics and Engineering, 168(1-4):185-206;
	//      1999 doi:10.1016/S0045-7825(98)00140-6
	//  [5] Santosh TV, Saraf RK, Ghosh AK and KushwahaHS. Optimum step length selection rule in
	//      modified HL–RF method for structural reliability. International Journal of Pressure
	//      Vessels and Piping, 83(10):742-748; 2006 doi:10.1016/j.ijpvp.2006.07.004
	//  [6] Haldar and Mahadevan. Probability, reliability and statistical methods in engineering
	//      and design. John Wiley & Sons. 304p; 2000.
	//
	// Strategies
	//  1: operates on x:
	//       argmin_x{ β(y(x)) | g(x) ≤ 0 }
	//  2: operates on x with equality constraint:
	//       argmin_x{ β(y(x)) + c(x) | |g(x)| ≤ ϵ }
	//  3: operates on y:
	//       argmin_y{ β(y) | g(x(y)) ≤ 0 }
	//  4: operates on y with equality constraint:
	//       argmin_y{ β(y) + c(y(x)) | |g(x(y))| ≤ ϵ }
	//  1 and 3:
	//       ova ← y dot y
	//       oor ← 0 ≥ g(x)
	//  2 and 4, equality constraint:
	//       c ← ϵ ≥ |g(x)|
	//       ova ← y dot y + c
	//       oor ← c

	// read parameters
	fn := "reliability-simple"
	fn, _ = io.ArgToFilename(0, fn, ".json", true)
	C := goga.ReadConfParams(fn)
	io.Pf("\n%s\nproblem # %v\n", utl.PrintThickLine(80), C.Problem)

	// initialise random numbers generator
	rnd.Init(C.Seed)

	// problems's data
	ϵ := 0.1                        // tolerance for strategy 2
	npts := 101                     // number of points for contour plot
	var g func(x []float64) float64 // limit state function
	var βref float64                // reference β (if available)
	var vars rnd.Variables          // random variables data
	axequal := true                 // plot with axis.equal

	// set problem
	switch C.Problem {

	// problem # 1 of [1] and Eq. (A.5) of [2]
	case 1:
		g = func(x []float64) float64 {
			return 0.1*math.Pow(x[0]-x[1], 2) - (x[0]+x[1])/math.Sqrt2 + 2.5
		}
		βref = 2.5 // from [1] xref = []float64{1.7677, 1.7677}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
		}

	// problem # 2 of [1] and Eq. (A.6) of [2]
	case 2:
		g = func(x []float64) float64 {
			return -0.5*math.Pow(x[0]-x[1], 2) - (x[0]+x[1])/math.Sqrt2 + 3
		}
		βref = 1.658 // from [2] xref = []float64{-0.7583, 1.4752}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
		}

	// problem # 3 from [1] and # 6 from [3]
	case 3:
		g = func(x []float64) float64 {
			return 2 - x[1] - 0.1*math.Pow(x[0], 2) + 0.06*math.Pow(x[0], 3)
		}
		βref = 2 // from [1] xref = []float64{0, 2}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
		}

	// problem # 4 from [1] and # 8 from [3]
	case 4:
		g = func(x []float64) float64 {
			return 3 - x[1] + 256*math.Pow(x[0], 4)
		}
		βref = 3 // from [1] xref = []float64{0, 3}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
		}

	// problem # 5 from [1] and # 1 from [4] (modified)
	case 5:
		shift := 0.0 //0.1
		g = func(x []float64) float64 {
			return 1 + math.Pow(x[0]+x[1]+shift, 2)/4 - 4*math.Pow(x[0]-x[1]+shift, 2)
		}
		βref = 0.3536 // from [1] xref = []float64{-βref * math.Sqrt2 / 2, βref * math.Sqrt2 / 2}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1},
		}

	// problem # 7 from [1] and example # 1 (case1) from [5]
	// x1 and x2 are normally distributed and statistically independent with
	case 6:
		g = func(x []float64) float64 {
			return math.Pow(x[0], 3) + math.Pow(x[1], 3) - 18
		}
		βref = 2.2401 // from [1]
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5},
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5},
		}

	// problem # 8 from [1] and example # 1 (case2) from [5]
	// x1 and x2 are normally distributed and statistically independent with
	case 7:
		g = func(x []float64) float64 {
			return math.Pow(x[0], 3) + math.Pow(x[1], 3) - 18
		}
		βref = 2.2260 // from [1]
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5},
			&rnd.VarData{D: rnd.D_Normal, M: 9.9, S: 5},
		}

	// problem # 9 from [1] and case 7 of [3]
	// x1 and x2 are normally distributed and statistically independent with
	case 8:
		g = func(x []float64) float64 {
			return 2.5 - 0.2357*(x[0]-x[1]) + 0.0046*math.Pow(x[0]+x[1]-20, 4)
		}
		βref = 2.5 // from [1]
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 3},
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 3},
		}

	// problem # 10 from [1] and example # 2 from [5]
	// x1 and x2 are normally distributed and statistically independent with
	case 9:
		g = func(x []float64) float64 {
			return math.Pow(x[0], 3) + math.Pow(x[1], 3) - 67.5
		}
		βref = 1.9003 // from [1]
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5},
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5},
		}

	// problem # 11 from [1] and case 2 of [3]
	// x1 and x2 are normally distributed and statistically independent with
	case 10:
		g = func(x []float64) float64 {
			return x[0]*x[1] - 146.14
		}
		βref = 5.4280 // from [1]
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 78064.4, S: 11709.7},
			&rnd.VarData{D: rnd.D_Normal, M: 0.0104, S: 0.00156},
		}
		axequal = false

	// problem # 12 from [1] and example 2 of [4]
	// x1 and x2 are normally distributed and statistically independent with
	case 11:
		g = func(x []float64) float64 {
			return 2.2257 - 0.025*math.Sqrt2*math.Pow(x[0]+x[1]-20, 3)/27 + 0.2357*(x[0]-x[1])
		}
		βref = 2.2257 // from [1]
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 3},
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 3},
		}

	// problem # 14 from [1]
	// x1 and x2 are log-normally distributed and statistically independent with
	case 12:
		g = func(x []float64) float64 {
			return x[0]*x[1] - 1140
		}
		βref = 5.2127 // from [1] // from here: 5.210977819456551
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Log, M: 38, S: 3.8, Std: true},
			&rnd.VarData{D: rnd.D_Log, M: 54, S: 2.7, Std: true},
		}

	// problem # 6 of [1] and case # 3 of [3]
	case 13:
		g = func(x []float64) float64 {
			sum := 0.0
			for i := 0; i < 9; i++ {
				sum += x[i] * x[i]
			}
			return 2.0 - 0.015*sum - x[9]
		}
		βref = 2.0 // from [1]
		vars = make([]*rnd.VarData, 10)
		for i := 0; i < 10; i++ {
			vars[i] = &rnd.VarData{D: rnd.D_Normal, M: 0, S: 1}
		}

	default:
		chk.Panic("problem number %d is invalid", C.Problem)
	}

	// initialise distributions
	err := vars.Init()
	if err != nil {
		chk.Panic("cannot initialise distributions:\n%v", err)
	}

	// guess search space
	nvars := len(vars)
	vmin, vmax := make([]float64, nvars), make([]float64, nvars)
	for i, d := range vars {
		vmin[i] = d.M - 4*d.S
		vmax[i] = d.M + 4*d.S
		if d.D == rnd.D_Log && vmin[i] < rnd.TOLMINLOG {
			vmin[i] = rnd.TOLMINLOG
		}
	}

	// objective value function
	ovfunc := func(ind *goga.Individual, idIsland, t int, report *bytes.Buffer) (ova, oor float64) {

		// original and normalised variables
		x := ind.GetFloats()
		y, invalid := vars.Transform(x)
		if invalid {
			oor = 1e+8
			return
		}

		// squared distance from origin to limit state curve in normalised space
		b := la.VecDot(y, y)

		// compute objective value
		switch C.Strategy {
		case 1: // argmin_{x,y}{ β(y(x)) | g(x) ≤ 0 }
			ova = b                          // ova ← y dot y
			oor = utl.GtePenalty(0, g(x), 1) // oor ← 0 ≥ g(x)
		case 2: // argmin_{x,y}{ β(y(x)) + c(x) | |g(x)| ≤ ϵ }
			c := utl.GtePenalty(ϵ, math.Abs(g(x)), 1) // c ← ϵ ≥ |g(x)|
			ova = b + c                               // ova ← y dot y + c
			oor = c                                   // oor ← c
		}
		return
	}

	// evolver
	evo := goga.NewEvolverFloatChromo(C, vmin, vmax, ovfunc, goga.NewBingoFloats(vmin, vmax))

	// benchmarking
	cpu0 := time.Now()

	// for a number of trials
	betas := make([]float64, C.Ntrials)
	for i := 0; i < C.Ntrials; i++ {

		// reset population
		if i > 0 {
			for _, isl := range evo.Islands {
				isl.Pop.GenFloatRandom(C, vmin, vmax)
			}
		}
		pop0 := evo.Islands[0].Pop.GetCopy()

		// run
		check := i == C.Ntrials-1
		evo.Run(check, check)

		// compute β
		x := evo.Best.GetFloats()
		y, invalid := vars.Transform(x)
		if invalid {
			chk.Panic("best individual has invalid values")
		}
		β := math.Sqrt(la.VecDot(y, y))
		betas[i] = β
		if check {
			io.Pf("β = %g (%g)\n\n", β, βref)
		}

		// plot contour
		if check && C.DoPlot && len(vars) == 2 {
			pop1 := evo.Islands[0].Pop
			extra := func() { plt.SetXnticks(11); plt.SetYnticks(11) }
			goga.PlotTwoVarsContour("/tmp/goga", io.Sf("rel-prob%d-ST%d-orig", C.Problem, C.Strategy), pop0, pop1, evo.Best, npts, extra, axequal,
				vmin, vmax, false, false, nil, nil, g, g)
			goga.PlotTwoVarsContour("/tmp/goga", io.Sf("rel-prob%d-ST%d-tran", C.Problem, C.Strategy), pop0, pop1, evo.Best, npts, extra, axequal,
				vmin, vmax, false, true, vars.Transform, nil, g, g)
		}
	}

	// benchmarking
	io.Pfcyan("\nelapsed time = %v\n", time.Now().Sub(cpu0))

	// analysis
	if C.Ntrials > 1 {
		βmin, βave, βmax, βdev := rnd.StatBasic(betas, true)
		io.Pf("\nβmin = %v\n", βmin)
		io.PfYel("βave = %v\n", βave)
		io.Pf("βmax = %v\n", βmax)
		io.Pf("βdev = %v\n\n", βdev)
		io.Pf(rnd.BuildTextHist(nice_num(βmin-0.005), nice_num(βmax+0.005), 11, betas, "%.3f", 60))
	}
}

// nice_num returns a truncated float
func nice_num(x float64) float64 {
	s := io.Sf("%.2f", x)
	return io.Atof(s)
}

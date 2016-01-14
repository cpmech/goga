// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	PI = math.Pi
	NU = 21
	NV = 21
)

func plot_plane() {
	N := []float64{1, 1, 1}   // normal
	P := []float64{0.5, 0, 0} // point on plane
	d := -N[0]*P[0] - N[1]*P[1] - N[2]*P[2]
	X, Y := utl.MeshGrid2D(0, 0.5, 0, 0.5, NU, NV)
	Z := la.MatAlloc(NV, NU)
	for j := 0; j < NU; j++ {
		for i := 0; i < NV; i++ {
			Z[i][j] = (-d - N[0]*X[i][j] - N[1]*Y[i][j]) / N[2]
			if Z[i][j] < -0.01 {
				Z[i][j] = math.NaN()
			}
		}
	}
	plt.Wireframe(X, Y, Z, "color='k', zmin=0, zmax=0.5")
}

func plot_sphere() {
	R := 1.0
	U, V := utl.MeshGrid2D(0, PI/2.0, 0, PI/2.0, NU, NV)
	X, Y, Z := la.MatAlloc(NV, NU), la.MatAlloc(NV, NU), la.MatAlloc(NV, NU)
	for j := 0; j < NU; j++ {
		for i := 0; i < NV; i++ {
			X[i][j] = R * math.Cos(U[i][j]) * math.Sin(V[i][j])
			Y[i][j] = R * math.Sin(U[i][j]) * math.Sin(V[i][j])
			Z[i][j] = R * math.Cos(V[i][j])
		}
	}
	plt.Wireframe(X, Y, Z, "color='k'")
}

func plot_convex() {
	X, Y := utl.MeshGrid2D(0, 1, 0, 1, NU, NV)
	Z := la.MatAlloc(NV, NU)
	for j := 0; j < NU; j++ {
		for i := 0; i < NV; i++ {
			Z[i][j] = 1.0 - math.Sqrt(X[i][j]) - math.Sqrt(Y[i][j])
			if Z[i][j] < -0.01 {
				Z[i][j] = math.NaN()
			}
		}
	}
	plt.Wireframe(X, Y, Z, "color='k', zmin=0, zmax=1.0")
}

func solve_problem(problem int) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 200
	opt.Ncpu = 4
	opt.Tf = 500

	// problem variables
	var nf, ng, nh int       // number of functions
	var fcn goga.MinProb_t   // functions
	var plot_solution func() // plot solution in 3D

	// problems
	switch problem {

	// problem # 1: DTLZ1
	case 1:
		opt.RptName = "DTLZ1"
		opt.FltMin = make([]float64, 7)
		opt.FltMax = make([]float64, 7)
		for i := 0; i < 7; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c := 5.0
			for i := 2; i < 7; i++ {
				c += math.Pow((x[i]-0.5), 2.0) - math.Cos(20.0*PI*(x[i]-0.5))
			}
			c *= 100.0
			f[0] = 0.5 * x[0] * x[1] * (1.0 + c)
			f[1] = 0.5 * x[0] * (1.0 - x[1]) * (1.0 + c)
			f[2] = 0.5 * (1.0 - x[0]) * (1.0 + c)
		}
		plot_solution = func() { plot_plane() }

	// problem # 2: DTLZ2
	case 2:
		opt.RptName = "DTLZ2"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
		}
		plot_solution = func() { plot_sphere() }

	// problem # 3: DTLZ3
	case 3:
		opt.Tf = 2000
		opt.RptName = "DTLZ3"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c := 10.0
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i]-0.5), 2.0) - math.Cos(20.0*PI*(x[i]-0.5))
			}
			c *= 100.0
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
		}
		plot_solution = func() { plot_sphere() }

	// problem # 4: DTLZ4
	case 4:
		opt.RptName = "DTLZ4"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			α := 100.0
			f[0] = (1.0 + c) * math.Cos(math.Pow(x[0], α)*PI/2.0) * math.Cos(math.Pow(x[1], α)*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(math.Pow(x[0], α)*PI/2.0) * math.Sin(math.Pow(x[1], α)*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(math.Pow(x[0], α)*PI/2.0)
		}
		plot_solution = func() { plot_sphere() }

	// problem # 5: DTLZ2x (convex)
	case 5:
		opt.RptName = "DTLZ2x"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 0, 0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
			f[0] = math.Pow(f[0], 4.0)
			f[1] = math.Pow(f[1], 4.0)
			f[2] = math.Pow(f[2], 2.0)
		}
		plot_solution = func() { plot_convex() }

	default:
		chk.Panic("problem %d is not available", problem)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.Solve()

	// plot solution
	if true {
		plt.SetForEps(1.0, 400)
		plot_solution()
	}

	// plot results
	if true {
		X, Y, Z := make([]float64, opt.Nsol), make([]float64, opt.Nsol), make([]float64, opt.Nsol)
		for i, sol := range opt.Solutions {
			X[i], Y[i], Z[i] = sol.Ova[0], sol.Ova[1], sol.Ova[2]
		}
		plt.Plot3dPoints(X, Y, Z, "s=5, color='r', facecolor='r', edgecolor='r', preservePrev=1, xlbl='$f_0$', ylbl='$f_1$', zlbl='$f_2$'")
		plt.Camera(10, 45, "")
		plt.AxDist(11.0)
		plt.SaveD("/tmp/goga", io.Sf("multiobj3_%s_A.eps", opt.RptName))
		plt.Camera(10, -45, "")
		plt.AxDist(11.0)
		plt.SaveD("/tmp/goga", io.Sf("multiobj3_%s_B.eps", opt.RptName))
	}

	//plt.Show()
	return
}

func main() {
	P := utl.IntRange2(1, 6)
	//P := []int{3}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem)
	}
}

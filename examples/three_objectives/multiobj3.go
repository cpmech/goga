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
	"github.com/cpmech/gosl/vtk"
)

const (
	SQ2 = math.Sqrt2
	PI  = math.Pi
	NU  = 21
	NV  = 21
)

func cone_angle(s []float64) float64 {
	den := s[0] + s[1] + s[2]
	if den < 1e-14 {
		return 1e30
	}
	return math.Sqrt(math.Pow(s[0]-s[1], 2.0)+math.Pow(s[1]-s[2], 2.0)+math.Pow(s[2]-s[0], 2.0)) / den
}

func rot_matrix() [][]float64 {
	if false {
		return [][]float64{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		}
	}
	SQ3 := math.Sqrt(3.0)
	return [][]float64{
		{0, SQ2 / SQ3, 1 / SQ3},
		{-1 / SQ2, -1 / (SQ2 * SQ3), 1 / SQ3},
		{1 / SQ2, -1 / (SQ2 * SQ3), 1 / SQ3},
	}
}

func plot_cone(α float64, preservePrev bool) {
	nu, nv := 11, 21
	l := 1.2
	r := math.Tan(α) * l
	S, T := utl.MeshGrid2D(0, l, 0, 2.0*PI, nu, nv)
	X := la.MatAlloc(nv, nu)
	Y := la.MatAlloc(nv, nu)
	Z := la.MatAlloc(nv, nu)
	u := make([]float64, 3)
	v := make([]float64, 3)
	L := rot_matrix()
	for j := 0; j < nu; j++ {
		for i := 0; i < nv; i++ {
			u[0] = S[i][j] * r * math.Cos(T[i][j])
			u[1] = S[i][j] * r * math.Sin(T[i][j])
			u[2] = S[i][j]
			la.MatVecMul(v, 1, L, u)
			X[i][j], Y[i][j], Z[i][j] = v[0], v[1], v[2]
		}
	}
	pp := 0
	if preservePrev {
		pp = 1
	}
	plt.Wireframe(X, Y, Z, io.Sf("color='b', lw=0.5, preservePrev=%d", pp))
}

func plot_plane(preservePrev bool) {
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
	pp := 0
	if preservePrev {
		pp = 1
	}
	plt.Wireframe(X, Y, Z, io.Sf("color='k', lw=0.5, preservePrev=%d", pp))
}

func plot_sphere(preservePrev bool) {
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
	pp := 0
	if preservePrev {
		pp = 1
	}
	plt.Wireframe(X, Y, Z, io.Sf("color='k', lw=0.5, preservePrev=%d", pp))
}

func plot_convex(level float64, preservePrev bool) {
	X, Y := utl.MeshGrid2D(0, 1, 0, 1, NU, NV)
	Z := la.MatAlloc(NV, NU)
	for j := 0; j < NU; j++ {
		for i := 0; i < NV; i++ {
			Z[i][j] = level - math.Sqrt(X[i][j]) - math.Sqrt(Y[i][j])
			if Z[i][j] < -0.01 {
				Z[i][j] = math.NaN()
			}
		}
	}
	pp := 0
	if preservePrev {
		pp = 1
	}
	plt.Wireframe(X, Y, Z, io.Sf("color='k', lw=0.5, preservePrev=%d", pp))
}

func solve_problem(problem int) (opt *goga.Optimiser) {

	io.Pf("\n\n------------------------------------- problem = %d ---------------------------------------\n", problem)

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.Nsol = 200
	opt.Ncpu = 8
	opt.Tf = 500
	opt.Ntrials = 10

	// options for report
	opt.RptFmtE = "%.4e"
	opt.RptFmtEdev = "%.4e"

	// problem variables
	var nf, ng, nh int       // number of functions
	var fcn goga.MinProb_t   // functions
	var plot_solution func() // plot solution in 3D
	rng := []float64{0, 1, 0, 1, 0, 1}

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
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0] + f[1] + f[2] - 0.5
		}
		plot_solution = func() { plot_plane(false) }
		rng = []float64{0, 0.5, 0, 0.5, 0, 0.5}

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
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() { plot_sphere(false) }

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
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() { plot_sphere(false) }

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
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() { plot_sphere(false) }

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
		opt.Multi_fcnErr = func(f []float64) float64 {
			return math.Pow(math.Abs(f[0]), 0.5) + math.Pow(math.Abs(f[1]), 0.5) + f[2] - 1.0
		}
		plot_solution = func() { plot_convex(1.0, false) }

	// problem # 2: DTLZ2c
	case 6:
		opt.RptName = "DTLZ2c"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 1, 0
		//α := math.Atan(1.0 / SQ2) // <<< touches lower plane
		//α = PI/2.0 - α // <<< touches upper plane
		α := 15.0 * PI / 180.0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			var c float64
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i] - 0.5), 2.0)
			}
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
			g[0] = math.Tan(α) - cone_angle(f)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() {
			plot_sphere(false)
			plot_cone(α, true)
		}

	// problem # 3: DTLZ3c
	case 7:
		opt.Tf = 2000
		opt.RptName = "DTLZ3c"
		opt.FltMin = make([]float64, 12)
		opt.FltMax = make([]float64, 12)
		for i := 0; i < 12; i++ {
			opt.FltMin[i], opt.FltMax[i] = 0, 1
		}
		nf, ng, nh = 3, 1, 0
		α := 15.0 * PI / 180.0
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			c := 10.0
			for i := 2; i < 12; i++ {
				c += math.Pow((x[i]-0.5), 2.0) - math.Cos(20.0*PI*(x[i]-0.5))
			}
			c *= 100.0
			f[0] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Cos(x[1]*PI/2.0)
			f[1] = (1.0 + c) * math.Cos(x[0]*PI/2.0) * math.Sin(x[1]*PI/2.0)
			f[2] = (1.0 + c) * math.Sin(x[0]*PI/2.0)
			g[0] = math.Tan(α) - cone_angle(f)
		}
		opt.Multi_fcnErr = func(f []float64) float64 {
			return f[0]*f[0] + f[1]*f[1] + f[2]*f[2] - 1.0
		}
		plot_solution = func() {
			plot_sphere(false)
			plot_cone(α, true)
		}

	default:
		chk.Panic("problem %d is not available", problem)
	}

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf, ng, nh)

	// solve
	opt.RunMany("", "")
	goga.StatMulti(opt, true)

	// check
	var failed bool
	for _, sol := range opt.Solutions {
		for _, oor := range sol.Oor {
			if oor > 0 {
				failed = true
			}
		}
	}
	if failed {
		io.PfRed("failed\n")
	} else {
		io.PfGreen("OK\n")
	}

	// results
	var X, Y, Z []float64
	onlyFront0 := true
	if onlyFront0 {
		for _, sol := range opt.Solutions {
			if sol.Feasible() && sol.FrontId == 0 {
				X = append(X, sol.Ova[0])
				Y = append(Y, sol.Ova[1])
				Z = append(Z, sol.Ova[2])
			}
		}
	} else {
		X, Y, Z = make([]float64, opt.Nsol), make([]float64, opt.Nsol), make([]float64, opt.Nsol)
		for i, sol := range opt.Solutions {
			X[i], Y[i], Z[i] = sol.Ova[0], sol.Ova[1], sol.Ova[2]
		}
	}

	// plot results
	if false {
		plt.SetForEps(1.0, 400)
		plot_solution()
		plt.Plot3dPoints(X, Y, Z, "s=7, color='r', facecolor='r', edgecolor='r', preservePrev=1, xlbl='$f_0$', ylbl='$f_1$', zlbl='$f_2$'")
		e, a := 10.0, 45.0
		if problem == 6 {
			e, a = 15, 30
		}
		plt.Camera(e, a, "")
		plt.AxDist(11.0)
		plt.AxisRange3d(rng[0], rng[1], rng[2], rng[3], rng[4], rng[5])
		plt.SaveD("/tmp/goga", io.Sf("py_%s_A.eps", opt.RptName))
		e, a = 10, -45
		if problem == 6 {
			e, a = 10, -45
		}
		plt.Camera(e, a, "")
		plt.AxDist(11.0)
		plt.AxisRange3d(rng[0], rng[1], rng[2], rng[3], rng[4], rng[5])
		plt.SaveD("/tmp/goga", io.Sf("py_%s_B.eps", opt.RptName))
		//plt.Show()
	}

	// vtk
	if false {

		// create a new VTK Scene
		scn := vtk.NewScene()
		scn.HydroLine = false
		scn.FullAxes = false
		scn.AxesLen = 1.1
		scn.WithPlanes = false
		scn.LblX = "f0"
		scn.LblY = "f1"
		scn.LblZ = "f2"
		scn.LblSz = 20
		if problem == 1 {
			scn.AxesLen = 0.6
		}

		// optimal Pareto front
		front := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
			f = opt.Multi_fcnErr(x)
			return
		})
		front.Limits = []float64{0, 1, 0, 1, 0, 1}
		front.Color = []float64{0.45098039, 0.70588235, 1., 0.8}
		front.CmapNclrs = 0 // use this to use specified color
		front.Ndiv = []int{61, 61, 61}
		front.AddTo(scn)

		// cone
		α := 15.0
		kα := math.Tan(α * math.Pi / 180.0)
		cone := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
			f = cone_angle(x) - kα
			return
		})
		if ng > 0 {
			cone.Limits = []float64{0, -1, 0, 1, 0, 360}
			cone.Ndiv = []int{61, 61, 81}
			cone.OctRotate = true
			cone.GridShowPts = false
			//cone.Color = []float64{0.61960784, 0.74117647, 0.54117647, 0.8}
			//cone.Color = []float64{0.94901961, 0.83921569, 0.66666667, 0.8}
			cone.Color = []float64{0.96862745, 0.75294118, 0.40784314, 0.5}
			cone.CmapNclrs = 0 // use this to use specified color
			cone.AddTo(scn)    // remember to add to Scene
		}

		// particles
		var P vtk.Spheres
		P.X = X
		P.Y = Y
		P.Z = Z
		P.R = utl.DblVals(len(X), 0.015)
		P.Color = []float64{1, 0, 0, 1}
		P.AddTo(scn)

		// start interactive mode
		scn.SaveEps = false
		scn.SavePng = true
		scn.PngMag = 2
		scn.Fnk = io.Sf("/tmp/goga/%s_A", opt.RptName)
		scn.Run()
		scn.Fnk = io.Sf("/tmp/goga/%s_B", opt.RptName)
		scn.Run()
	}

	return
}

func main() {
	P := utl.IntRange2(2, 7)
	//P := []int{6}
	opts := make([]*goga.Optimiser, len(P))
	for i, problem := range P {
		opts[i] = solve_problem(problem)
	}
	io.Pf("\n-------------------------- generating report --------------------------\nn")
	nRowPerTab := 10
	goga.TexMultiReport("/tmp/goga", "tmp_threeobj", "threeobj", nRowPerTab, true, opts)
	goga.TexMultiReport("/tmp/goga", "threeobj", "threeobj", nRowPerTab, false, opts)
}

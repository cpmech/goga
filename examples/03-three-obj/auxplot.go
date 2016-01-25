// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/fun"
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

func cosX(w, m float64) float64 { return fun.Sign(math.Cos(w)) * math.Pow(math.Abs(math.Cos(w)), m) }
func sinX(w, m float64) float64 { return fun.Sign(math.Sin(w)) * math.Pow(math.Abs(math.Sin(w)), m) }

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

func plot_plane_axis(loc float64, dir int, preservePrev bool) {
	if !preservePrev {
		plt.PyCmds("gcf().add_subplot(111, projection='3d')\n")
	}
	sdir := []string{"x", "y", "z"}
	plt.PyCmds(io.Sf(`
from matplotlib.collections import PolyCollection
xs = [0.0, 1.0, 1.0, 0.0]
ys = [0.0, 0.0, 1.0, 1.0]
verts = [zip(xs, ys)]
poly = PolyCollection(verts)
gca().add_collection3d(poly, zs=%g, zdir='%s')
`, loc, sdir[dir]))
	return
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

func plot_superquadric(a, b, c float64, preservePrev bool) {
	A, B, C := 2.0/a, 2.0/b, 2.0/c
	R := 1.0
	U, V := utl.MeshGrid2D(0, PI/2.0, 0, PI/2.0, NU, NV)
	X, Y, Z := la.MatAlloc(NV, NU), la.MatAlloc(NV, NU), la.MatAlloc(NV, NU)
	for j := 0; j < NU; j++ {
		for i := 0; i < NV; i++ {
			X[i][j] = R * cosX(U[i][j], A) * sinX(V[i][j], A)
			Y[i][j] = R * sinX(U[i][j], B) * sinX(V[i][j], B)
			Z[i][j] = R * cosX(V[i][j], C)
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

func py_plot3(iOva, jOva, kOva int, opt *goga.Optimiser, plot_solution func(), onlyFront0, twice bool) {

	// results
	var X, Y, Z []float64
	if onlyFront0 {
		for _, sol := range opt.Solutions {
			if sol.Feasible() && sol.FrontId == 0 {
				X = append(X, sol.Ova[iOva])
				Y = append(Y, sol.Ova[jOva])
				Z = append(Z, sol.Ova[kOva])
			}
		}
	} else {
		X, Y, Z = make([]float64, opt.Nsol), make([]float64, opt.Nsol), make([]float64, opt.Nsol)
		for i, sol := range opt.Solutions {
			X[i], Y[i], Z[i] = sol.Ova[iOva], sol.Ova[jOva], sol.Ova[kOva]
		}
	}

	// plot
	plt.SetForEps(1.0, 400)
	plot_solution()
	plt.Plot3dPoints(X, Y, Z, "s=7, color='r', facecolor='r', edgecolor='r', preservePrev=1, xlbl='$f_0$', ylbl='$f_1$', zlbl='$f_2$'")
	e, a := 10.0, 45.0
	if opt.RptName == "DTLZ2c" {
		e, a = 15, 30
	}
	//plt.Camera(e, a, "")
	//plt.AxDist(11.0)
	//plt.AxisRange3d(opt.RptFmin[iOva], opt.RptFmax[iOva], opt.RptFmin[jOva], opt.RptFmax[jOva], opt.RptFmin[kOva], opt.RptFmax[kOva])
	//plt.SaveD("/tmp/goga", io.Sf("py_%s_A.eps", opt.RptName))
	if twice {
		e, a = 10, -45
		if opt.RptName == "DTLZ2c" {
			e, a = 10, -45
		}
		plt.Camera(e, a, "")
		plt.AxDist(11.0)
		plt.AxisRange3d(opt.RptFmin[iOva], opt.RptFmax[iOva], opt.RptFmin[jOva], opt.RptFmax[jOva], opt.RptFmin[kOva], opt.RptFmax[kOva])
		plt.SaveD("/tmp/goga", io.Sf("py_%s_B.eps", opt.RptName))
	}
}

func vtk_plot3(opt *goga.Optimiser, αcone, ptRad float64, onlyFront0, twice bool) {

	// results
	var X, Y, Z []float64
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

	// create a new VTK Scene
	scn := vtk.NewScene()
	scn.HydroLine = false
	scn.FullAxes = false
	scn.AxesLen = 1.1
	scn.WithPlanes = false
	scn.LblX = io.Sf("f%d", 0)
	scn.LblY = io.Sf("f%d", 1)
	scn.LblZ = io.Sf("f%d", 2)
	scn.LblSz = 20
	if opt.RptName == "DTLZ1" {
		scn.AxesLen = 0.6
	}

	// optimal Pareto front
	front := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
		f = opt.Multi_fcnErr(x)
		return
	})
	front.Limits = []float64{opt.RptFmin[0], opt.RptFmax[0], opt.RptFmin[1], opt.RptFmax[1], opt.RptFmin[2], opt.RptFmax[2]}
	front.Color = []float64{0.45098039, 0.70588235, 1., 0.8}
	front.CmapNclrs = 0 // use this to use specified color
	front.Ndiv = []int{61, 61, 61}
	front.AddTo(scn)

	// cone
	if opt.RptName == "DTLZ2c" {
		cone := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
			f = cone_angle(x) - math.Tan(αcone)
			return
		})
		cone.Limits = []float64{0, -1, 0, 1, 0, 360}
		cone.Ndiv = []int{61, 61, 81}
		cone.OctRotate = true
		cone.GridShowPts = false
		cone.Color = []float64{0.96862745, 0.75294118, 0.40784314, 0.5}
		cone.CmapNclrs = 0 // use this to use specified color
		cone.AddTo(scn)    // remember to add to Scene
	}

	// particles
	var P vtk.Spheres
	P.X, P.Y, P.Z = X, Y, Z
	P.R = utl.DblVals(len(X), ptRad)
	P.Color = []float64{1, 0, 0, 1}
	P.AddTo(scn)

	// start interactive mode
	scn.SaveEps = false
	scn.SavePng = true
	scn.PngMag = 2
	scn.Fnk = io.Sf("/tmp/goga/vtk_%s_A", opt.RptName)
	scn.Run()
	if twice {
		scn.Fnk = io.Sf("/tmp/goga/vtk_%s_B", opt.RptName)
		scn.Run()
	}
}

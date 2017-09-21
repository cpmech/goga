// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
	"github.com/cpmech/gosl/vtk"
)

// PyPlot3 plots 3d space (Python)
func PyPlot3(iOva, jOva, kOva int, opt *goga.Optimiser, surface func(), onlyFront0 bool) {

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

	// lengths
	triad := 1.2
	min, max := -0.1, 1.1
	if opt.RptName == "DTLZ1" {
		triad = 0.65
		max = 0.6
	}

	// start
	plt.Reset(true, &plt.A{Dpi: 200, Eps: Eps})

	// axes
	xLbl := io.Sf("$f_{%d}$", iOva)
	yLbl := io.Sf("$f_{%d}$", jOva)
	zLbl := io.Sf("$f_{%d}$", kOva)
	plt.Triad(triad, xLbl, yLbl, zLbl, &plt.A{C: "#4d9d85"}, &plt.A{C: "#1b6650"})

	// plot surface
	if surface != nil {
		surface()
	}

	// plot points
	plt.Plot3dPoints(X, Y, Z, &plt.A{C: "#e86f16", Ec: "r"})

	// set axes labels
	plt.SetLabels3d("", "", "", nil)

	// set camera and limits
	plt.Default3dView(min, max, min, max, min, max, true)

	// save figure
	plt.Save("/tmp/goga", "fig-"+opt.RptName)
}

// VtkPlot3 plots 3d space (VTK)
func VtkPlot3(opt *goga.Optimiser, αcone, ptRad float64, onlyFront0, twice bool) {

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

	// set camera
	scn.Width = 600
	scn.Height = 600
	scn.Zoom = 1.4
	if opt.RptName == "DTLZ2c" {
		scn.Zoom = 1.5
		scn.SetCamera(0, 0, 1, 0, 0, 0, 2, 1, 1.5)
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
			f = plt.CalcDiagAngle(x) - math.Tan(αcone)
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
	P.R = utl.Vals(len(X), ptRad)
	P.Color = []float64{1, 0, 0, 1}
	P.AddTo(scn)

	// start interactive mode
	scn.SaveEps = true
	scn.SavePng = false
	scn.PngMag = 2
	scn.Fnk = io.Sf("/tmp/goga/vtk_%s_A", opt.RptName)
	scn.Run()
	if twice {
		scn.Fnk = io.Sf("/tmp/goga/vtk_%s_B", opt.RptName)
		scn.Run()
	}
}

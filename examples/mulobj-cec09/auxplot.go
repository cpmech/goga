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

func plot2(opt *goga.Optimiser, onlyFront0 bool) {

	// plot reference values
	f0 := utl.GetColumn(0, opt.Multi_fStar)
	f1 := utl.GetColumn(1, opt.Multi_fStar)
	plt.Plot(f0, f1, "'b-', label='reference'")

	// plot goga values
	fmt := &plt.Fmt{C: "r", M: "o", Ms: 3, L: "goga", Ls: "None"}
	opt.PlotAddOvaOva(0, 1, opt.Solutions, true, fmt)
	plt.Gll("$f_0$", "$f_1$", "")
	plt.SaveD("/tmp/goga", io.Sf("m2_%s.eps", opt.RptName))
}

func plot3(opt *goga.Optimiser, onlyFront0, twice bool, ptRad float64) {

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

	// particles
	var P vtk.Spheres
	P.X, P.Y, P.Z = X, Y, Z
	P.R = utl.Vals(len(X), ptRad)
	P.Color = []float64{1, 0, 0, 1}
	P.AddTo(scn)

	// start interactive mode
	scn.SaveEps = false
	scn.SavePng = true
	scn.PngMag = 2
	scn.Fnk = io.Sf("/tmp/goga/m3_%s_A", opt.RptName)
	scn.Run()
	if twice {
		scn.Fnk = io.Sf("/tmp/goga/m3_%s_B", opt.RptName)
		scn.Run()
	}
}

func plot3x(opt *goga.Optimiser, onlyFront0 bool, i, j, k int, ptRad float64) {

	// points
	var X, Y, Z []float64
	if onlyFront0 {
		for _, sol := range opt.Solutions {
			if sol.Feasible() && sol.FrontId == 0 {
				X = append(X, sol.Flt[i])
				Y = append(Y, sol.Flt[j])
				Z = append(Z, sol.Flt[k])
			}
		}
	} else {
		X, Y, Z = make([]float64, opt.Nsol), make([]float64, opt.Nsol), make([]float64, opt.Nsol)
		for m, sol := range opt.Solutions {
			X[m], Y[m], Z[m] = sol.Flt[i], sol.Flt[j], sol.Flt[k]
		}
	}

	// create a new VTK Scene
	scn := vtk.NewScene()
	scn.HydroLine = false
	scn.FullAxes = true
	scn.AxesLen = 1.1
	scn.WithPlanes = false
	scn.LblX = io.Sf("x%d", i)
	scn.LblY = io.Sf("x%d", j)
	scn.LblZ = io.Sf("x%d", k)
	scn.LblSz = 20

	// reference particles
	var Ps vtk.Spheres
	switch opt.RptName {
	case "UF3":
		np := 101
		nx := opt.Nsol
		c1 := 0.5 * (1.0 + 3.0*(float64(1)-2.0)/(float64(nx)-2.0))
		c2 := 0.5 * (1.0 + 3.0*(float64(2)-2.0)/(float64(nx)-2.0))
		Ps.X = utl.LinSpace(0, 1, np)
		Ps.Y = make([]float64, np)
		Ps.Z = make([]float64, np)
		for i := 0; i < np; i++ {
			Ps.Y[i] = math.Pow(Ps.X[i], c1)
			Ps.Z[i] = math.Pow(Ps.X[i], c2)
		}
		Ps.R = utl.DblVals(np, 0.7*ptRad)
		Ps.Color = []float64{0, 0, 1, 1}
		Ps.AddTo(scn)
		scn.FullAxes = false
	}

	// particles
	var P vtk.Spheres
	P.X, P.Y, P.Z = X, Y, Z
	P.R = utl.DblVals(len(X), ptRad)
	P.Color = []float64{1, 0, 0, 1}
	P.AddTo(scn)

	// start interactive mode
	scn.SaveEps = false
	scn.SavePng = false
	scn.PngMag = 2
	scn.Fnk = io.Sf("/tmp/goga/m3_pts_%s", opt.RptName)
	scn.Run()
}

// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cec09

import (
	"testing"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
	"github.com/cpmech/gosl/vtk"
)

func Test_data2d(tst *testing.T) {
	prob := "CF4"
	dat := PFdata(prob)
	X := utl.DblsGetColumn(0, dat)
	Y := utl.DblsGetColumn(1, dat)
	plt.SetForEps(1.0, 250)
	plt.Plot(X, Y, "'r.'")
	plt.Gll("$f_1$", "$f_2$", "")
	plt.SaveD("/tmp/goga", io.Sf("cec09-%s.eps", prob))
}

func Test_data3d(tst *testing.T) {

	// data
	prob := "CF9"
	dat := PFdata(prob)
	X := utl.DblsGetColumn(0, dat)
	Y := utl.DblsGetColumn(1, dat)
	Z := utl.DblsGetColumn(2, dat)

	// figure
	plt.SetForEps(1.0, 400)
	plt.Plot3dPoints(X, Y, Z, "s=0.05, color='r', facecolor='r', edgecolor='r', xlbl='$f_1$', ylbl='$f_2$', zlbl='$f_3$'")
	plt.AxisRange3d(0, 1, 0, 1, 0, 1)
	plt.Camera(10, -135, "")
	//plt.Camera(10, 45, "")
	plt.SaveD("/tmp/goga", io.Sf("cec09-%s.eps", prob))

	// interactive
	if false {
		r := 0.005
		scn := vtk.NewScene()
		P := vtk.Spheres{X: X, Y: Y, Z: Z, R: utl.DblVals(len(X), r), Color: []float64{1, 0, 0, 1}}
		P.AddTo(scn)
		scn.Run()
	}
}

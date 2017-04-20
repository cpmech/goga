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
	X := utl.GetColumn(0, dat)
	Y := utl.GetColumn(1, dat)
	plt.Reset(false, nil)
	plt.Plot(X, Y, nil)
	plt.Gll("$f_1$", "$f_2$", nil)
	plt.Save("/tmp/goga", io.Sf("cec09-%s", prob))
}

func Test_data3d(tst *testing.T) {

	// data
	prob := "CF9"
	dat := PFdata(prob)
	X := utl.GetColumn(0, dat)
	Y := utl.GetColumn(1, dat)
	Z := utl.GetColumn(2, dat)

	// figure
	plt.Reset(false, nil)
	plt.Plot3dPoints(X, Y, Z, nil)
	plt.AxisRange3d(0, 1, 0, 1, 0, 1)
	plt.Camera(10, -135, nil)
	//plt.Camera(10, 45, nil)
	plt.Save("/tmp/goga", io.Sf("cec09-%s", prob))

	// interactive
	if false {
		r := 0.005
		scn := vtk.NewScene()
		P := vtk.Spheres{X: X, Y: Y, Z: Z, R: utl.Vals(len(X), r), Color: []float64{1, 0, 0, 1}}
		P.AddTo(scn)
		scn.Run()
	}
}

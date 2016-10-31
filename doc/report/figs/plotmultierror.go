// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {
	res := [][]float64{
		{5, 1.4183e-13},
		{7, 1.4620e-11},
		{10, 3.2089e-9},
		{13, 1.9350e-7},
		{15, 1.7667e-6},
		{20, 1.2063e-4},
	}
	X := utl.DblsGetColumn(0, res)
	Y := utl.DblsGetColumn(1, res)
	plt.SetForEps(0.75, 220)
	plt.SetYlog()
	plt.Plot(X, Y, "'b-',marker='.',ms=5,clip_on=0")
	plt.Gll("$N_f$", "$E_{min}$", "")
	plt.SaveD("/tmp/goga", "multierror.eps")
}

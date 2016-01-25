// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	nsol := 500
	tf := 1000
	time := []float64{
		223.42398255, // 1
		82.864179318, // 2
		50.945049948, // 3
		29.547849719, // 4
		20.741909766, // 5
		17.188611472, // 6
		15.937424833, // 7
		13.919150335, // 8
		12.589523593, // 9
		12.078978314, // 10
		11.034417259, // 11
		9.542071936,  // 12
		9.298965819,  // 13
		9.182769212,  // 14
		8.610938487,  // 15
		8.482685187,  // 16
	}

	ncpu := utl.LinSpace(1, 16, 16)

	speedup := make([]float64, 16)
	speedup[0] = 1
	for i := 1; i < 16; i++ {
		speedup[i] = time[0] / time[i] // Told / Tnew
	}

	io.Pforan("ncpu = %v\n", ncpu)

	plt.SetForEps(0.75, 250)
	plt.Plot(ncpu, speedup, io.Sf("'b-',marker='.', label='speedup: $N_{sol}=%d,\\,t_f=%d$', clip_on=0, zorder=100", nsol, tf))
	plt.Plot([]float64{1, 16}, []float64{1, 16}, "'k--',zorder=50")
	plt.Gll("$N_{cpu}:\\;$ number of groups", "speedup", "")
	plt.DoubleYscale("$T_{sys}:\\;$ system time [s]")
	plt.Plot(ncpu, time, "'k-',color='gray', clip_on=0")
	plt.SaveD("/tmp/goga", "topology-speedup.eps")
}

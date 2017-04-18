// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "github.com/cpmech/gosl/plt"

func main() {
	Nf := []float64{5, 7, 10, 13, 15, 20}
	Eave := []float64{3.5998e-12, 2.9629e-10, 6.0300e-8, 3.3686e-6, 2.5914e-5, 1.1966e-3}
	plt.SetForEps(0.75, 200)
	plt.Plot(Nf, Eave, "'b-', marker='.', clip_on=0")
	plt.SetYlog()
	plt.Gll("$N_f$", "$E_{ave}$", "")
	plt.SaveD("/tmp/goga", "multierror.eps")
}

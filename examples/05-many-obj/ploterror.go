// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "github.com/cpmech/gosl/plt"

func main() {
	Nf := []float64{5, 7, 10, 13, 15, 20}
	Eave := []float64{2.33e-12, 2.39e-10, 5.76e-8, 2.39e-6, 2.58e-5, 1.12e-3}
	plt.Reset(true, &plt.A{Eps: true, Prop: 0.75, WidthPt: 220})
	plt.HideBorders(&plt.A{HideR: true, HideT: true})
	plt.Plot(Nf, Eave, &plt.A{C: "r", M: ".", Lw: 1.2, NoClip: true})
	plt.SetYlog()
	plt.Gll("$N_f$", "$E_{ave}$", nil)
	plt.Save("/tmp/goga", "multierror")
}

// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gofem/fem"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// input filename
	fn, fnkey := io.ArgToFilename(0, "ground10", ".sim", true)

	// fem analysis
	analysis := fem.NewMain(fn, "", false, false, false, false, false, 0)

	// domain
	domain := analysis.Domains[0]

	// set figure
	plt.Reset(true, &plt.A{Eps: true, Prop: 0.65, WidthPt: 340})
	plt.Grid(&plt.A{C: "grey"})

	// draw structure
	argsLins := make(map[int]*plt.A)
	argsLins[-1] = &plt.A{C: "#004cc9", Lw: 3, NoClip: true}
	domain.Msh.Draw2d(true, false, false, nil, argsLins, nil)

	// draw arrows
	argsArrow := &plt.A{Fc: "orange", Ec: "orange", Scale: 12, Z: 100, NoClip: true}
	plt.SetTicksXlist([]float64{0, 60, 120, 180, 240, 300, 360, 420, 480, 540, 600, 660, 720})
	plt.SetTicksYlist([]float64{0, 60, 120, 180, 240, 300, 360})
	plt.Arrow(360, 0, 360, -100, argsArrow)
	plt.Arrow(720, 0, 720, -100, argsArrow)

	// draw polygon
	L := 20.0
	l := L / 2.0
	H := 360.0
	argsPoly := &plt.A{Fc: "r", Ec: "r", Closed: true, Z: 10, NoClip: true}
	plt.Polyline([][]float64{{-l, -l}, {l, -l}, {l, l}, {-l, l}}, argsPoly)
	plt.Polyline([][]float64{{-l, H - l}, {l, H - l}, {l, H + l}, {-l, H + l}}, argsPoly)

	// text
	fsz := 9.0
	plt.Text(0, -20, "fully fixed", &plt.A{Ha: "left", Va: "top", Fsz: fsz})
	plt.Text(0, 380, "fully fixed", &plt.A{Ha: "left", Va: "bottom", Fsz: fsz})
	plt.Text(350, -50, "100", &plt.A{Ha: "right", Va: "center", Fsz: fsz})
	plt.Text(710, -50, "100", &plt.A{Ha: "right", Va: "center", Fsz: fsz})

	// configure figure and save file
	plt.HideAllBorders()
	plt.Equal()
	err := plt.Save("/tmp/goga", "mesh-"+fnkey)
	if err != nil {
		io.PfRed("save filed:\n%v\n", err)
	}
}

// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func Test_igd01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("igd. igd metric with star equal to trial => igd=0")

	// load star values
	prob := "UF1"
	fStar := io.ReadMatrix(io.Sf("./examples/mulobj-cec09/cec09/pf_data/%s.dat", prob))
	npts := len(fStar)

	// optimiser
	var opt Optimiser
	opt.Default()
	opt.Nsol = npts
	opt.Ncpu = 1
	opt.FltMin = []float64{0, 0} // used to store fStar
	opt.FltMax = []float64{1, 1} // used to store fStar
	nf, ng, nh := 2, 0, 0

	// generator (store fStar into Flt)
	gen := func(sols []*Solution, prms *Parameters, reset bool) {
		for i, sol := range sols {
			if reset {
				sol.Reset(i)
			}
			sol.Flt[0], sol.Flt[1] = fStar[i][0], fStar[i][1]
		}
	}

	// objective function (copy fStar from Flt into Ova)
	obj := func(f, g, h, x []float64, y []int, cpu int) {
		f[0], f[1] = x[0], x[1]
	}

	// initialise optimiser
	opt.Init(gen, nil, obj, nf, ng, nh)

	// compute igd
	igd := opt.calcIgd(fStar)
	io.Pforan("igd = %v\n", igd)
	chk.Float64(tst, "igd", 1e-15, igd, 0)

	// plot
	if chk.Verbose {
		fmt := &plt.A{C: "red", M: ".", Ms: 1, Ls: "None", L: "solutions"}
		fS0 := utl.GetColumn(0, fStar)
		fS1 := utl.GetColumn(1, fStar)
		io.Pforan("len(fS0) = %v\n", len(fS0))
		plt.Reset(false, nil)
		opt.PlotAddOvaOva(0, 1, opt.Solutions, true, fmt)
		plt.Plot(fS0, fS1, &plt.A{C: "b", Ms: 2, L: io.Sf("star(%s)", prob)})
		plt.Gll("$f_0$", "$f_1$", nil)
		plt.Save("/tmp/goga", "igd01")
	}
}

func Test_igd02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("igd. igd metric with numerical values")

	// load star values
	prob := "UF1"
	fStar := io.ReadMatrix(io.Sf("./examples/mulobj-cec09/cec09/pf_data/%s.dat", prob))

	// load reference numerical solution (MOEAD/CEC09)
	fNum := io.ReadMatrix("./data/MOEAD-CEC09_PF_UF1_01.dat")
	nNum := len(fNum)

	// optimiser
	var opt Optimiser
	opt.Default()
	opt.Nsol = nNum
	opt.Ncpu = 1
	opt.FltMin = []float64{0, 0} // used to store fStar
	opt.FltMax = []float64{1, 1} // used to store fStar
	nf, ng, nh := 2, 0, 0

	// generator (store fNum into Flt)
	gen := func(sols []*Solution, prms *Parameters, reset bool) {
		for i, sol := range sols {
			if reset {
				sol.Reset(i)
			}
			sol.Flt[0], sol.Flt[1] = fNum[i][0], fNum[i][1]
		}
	}

	// objective function (copy fStar from Flt into Ova)
	obj := func(f, g, h, x []float64, y []int, cpu int) {
		f[0], f[1] = x[0], x[1]
	}

	// initialise optimiser
	opt.Init(gen, nil, obj, nf, ng, nh)

	// compute igd
	igd := opt.calcIgd(fStar)
	io.Pforan("igd = %v\n", igd)
	chk.Float64(tst, "igd", 1e-15, igd, 1.5019007244180080e-03)

	// plot
	if chk.Verbose {
		fmt := &plt.A{C: "red", M: ".", Ms: 1, Ls: "None", L: "solutions"}
		fS0 := utl.GetColumn(0, fStar)
		fS1 := utl.GetColumn(1, fStar)
		io.Pforan("len(fS0) = %v\n", len(fS0))
		plt.Reset(false, nil)
		opt.PlotAddOvaOva(0, 1, opt.Solutions, true, fmt)
		plt.Plot(fS0, fS1, &plt.A{C: "b", Ms: 2, L: io.Sf("star(%s)", prob)})
		plt.Gll("$f_0$", "$f_1$", nil)
		plt.Save("/tmp/goga", "igd02")
	}
}

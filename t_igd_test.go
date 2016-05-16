// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
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
	fStar, err := io.ReadMatrix(io.Sf("./examples/mulobj-cec09/cec09/pf_data/%s.dat", prob))
	if err != nil {
		tst.Errorf("cannot read fStar matrix:\n%v", err)
		return
	}
	npts := len(fStar)

	// optimiser
	var opt Optimiser
	opt.Default()
	opt.Nsol = npts
	opt.Ncpu = 1
	opt.Xmin = []float64{0, 0} // used to store fStar
	opt.Xmax = []float64{1, 1} // used to store fStar
	nf, ng, nh := 2, 0, 0

	// generator (store fStar into Flt)
	gen := func(sols []*Solution, prms *Parameters) {
		for i, sol := range sols {
			sol.Flt[0], sol.Flt[1] = fStar[i][0], fStar[i][1]
		}
	}

	// objective function (copy fStar from Flt into Ova)
	obj := func(f, g, h, x []float64, ξ []int, cpu int) {
		f[0], f[1] = x[0], x[1]
	}

	// initialise optimiser
	opt.Init(gen, nil, obj, nf, ng, nh)

	// compute igd
	igd := StatIgd(&opt, fStar)
	io.Pforan("igd = %v\n", igd)
	chk.Scalar(tst, "igd", 1e-15, igd, 0)

	// plot
	if chk.Verbose {
		fmt := &plt.Fmt{C: "red", M: ".", Ms: 1, Ls: "None", L: "solutions"}
		fS0 := utl.DblsGetColumn(0, fStar)
		fS1 := utl.DblsGetColumn(1, fStar)
		io.Pforan("len(fS0) = %v\n", len(fS0))
		plt.SetForEps(0.75, 300)
		opt.PlotAddOvaOva(0, 1, opt.Solutions, true, fmt)
		plt.Plot(fS0, fS1, io.Sf("'b.', ms=2, label='star(%s)', clip_on=0", prob))
		plt.Gll("$f_0$", "$f_1$", "")
		plt.SaveD("/tmp/goga", "igd01.eps")
	}
}

func Test_igd02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("igd. igd metric with numerical values")

	// load star values
	prob := "UF1"
	fStar, err := io.ReadMatrix(io.Sf("./examples/mulobj-cec09/cec09/pf_data/%s.dat", prob))
	if err != nil {
		tst.Errorf("cannot read fStar matrix:\n%v", err)
		return
	}

	// load reference numerical solution (MOEAD/CEC09)
	fNum, err := io.ReadMatrix("./data/MOEAD-CEC09_PF_UF1_01.dat")
	if err != nil {
		tst.Errorf("cannot read fNum matrix:\n%v", err)
		return
	}
	nNum := len(fNum)

	// optimiser
	var opt Optimiser
	opt.Default()
	opt.Nsol = nNum
	opt.Ncpu = 1
	opt.Xmin = []float64{0, 0} // used to store fStar
	opt.Xmax = []float64{1, 1} // used to store fStar
	nf, ng, nh := 2, 0, 0

	// generator (store fNum into Flt)
	gen := func(sols []*Solution, prms *Parameters) {
		for i, sol := range sols {
			sol.Flt[0], sol.Flt[1] = fNum[i][0], fNum[i][1]
		}
	}

	// objective function (copy fStar from Flt into Ova)
	obj := func(f, g, h, x []float64, ξ []int, cpu int) {
		f[0], f[1] = x[0], x[1]
	}

	// initialise optimiser
	opt.Init(gen, nil, obj, nf, ng, nh)

	// compute igd
	igd := StatIgd(&opt, fStar)
	io.Pforan("igd = %v\n", igd)
	chk.Scalar(tst, "igd", 1e-15, igd, 1.5019007244180080e-03)

	// plot
	if chk.Verbose {
		fmt := &plt.Fmt{C: "red", M: ".", Ms: 1, Ls: "None", L: "solutions"}
		fS0 := utl.DblsGetColumn(0, fStar)
		fS1 := utl.DblsGetColumn(1, fStar)
		io.Pforan("len(fS0) = %v\n", len(fS0))
		plt.SetForEps(0.75, 300)
		opt.PlotAddOvaOva(0, 1, opt.Solutions, true, fmt)
		plt.Plot(fS0, fS1, io.Sf("'b.', ms=2, label='star(%s)', clip_on=0", prob))
		plt.Gll("$f_0$", "$f_1$", "")
		plt.SaveD("/tmp/goga", "igd02.eps")
	}
}

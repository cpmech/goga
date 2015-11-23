// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

func get_individual(id int) (ind *Individual) {
	nova := 2
	noor := 3
	switch id {
	case 0:
		ind = NewIndividual(nova, noor, []float64{4.4, 5.5, 666}, []int{1, 20, 300})
		ind.Ovas[0] = 123
		ind.Ovas[1] = 345
		ind.Oors[0] = 10
		ind.Oors[1] = 20
		ind.Oors[2] = 30
	case 1:
		ind = NewIndividual(nova, noor, []float64{104.4, 105.5, 6.66}, []int{-1, -20, -300})
		ind.Ovas[0] = 200
		ind.Ovas[1] = 100
		ind.Oors[0] = 15
		ind.Oors[1] = 25
		ind.Oors[2] = 35
	}
	return
}

func Test_ind01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind01. representation and copying")

	rnd.Init(0)

	A := get_individual(0)
	B := A.GetCopy()
	chk.Scalar(tst, "ova0", 1e-17, B.Ovas[0], 123)
	chk.Scalar(tst, "ova1", 1e-17, B.Ovas[1], 345)
	chk.Scalar(tst, "oor0", 1e-17, B.Oors[0], 10)
	chk.Scalar(tst, "oor1", 1e-17, B.Oors[1], 20)
	chk.Scalar(tst, "oor2", 1e-17, B.Oors[2], 30)
}

func Test_ind02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind02. copy into")

	rnd.Init(0)

	A := get_individual(0)
	B := get_individual(1)

	var ops OpsData
	ops.SetDefault()
	ops.IntPc = 1.0
	ops.Cuts = []int{1, 2}
	ops.Xrange = [][]float64{{0, 1}, {-20, 20}, {-300, 300}}

	a := A.GetCopy()
	b := A.GetCopy()
	IndCrossover(a, b, A, B, nil, nil, 0, &ops)

	chk.Ints(tst, "a.Ints   ", a.Ints, []int{1, -20, 300})
	chk.Ints(tst, "b.Ints   ", b.Ints, []int{-1, 20, -300})

	x := get_individual(0)
	x.Ovas = []float64{0, 0}
	x.Oors = []float64{0, 0, 0}
	B.CopyInto(x)

	chk.Scalar(tst, "ova0", 1e-17, x.Ovas[0], 200)
	chk.Scalar(tst, "ova1", 1e-17, x.Ovas[1], 100)
	chk.Scalar(tst, "oor0", 1e-17, x.Oors[0], 15)
	chk.Scalar(tst, "oor1", 1e-17, x.Oors[1], 25)
	chk.Scalar(tst, "oor2", 1e-17, x.Oors[2], 35)
}

func Test_ind03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind03. comparing")

	A := get_individual(0)
	B := get_individual(1)
	A_dominates, B_dominates := IndCompare(A, B)
	io.Pfblue2("A: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	io.Pforan("B_dominates = %v\n", B_dominates)
	if !A_dominates {
		tst.Errorf("test failed\n")
		return
	}
	if B_dominates {
		tst.Errorf("test failed\n")
		return
	}

	A.Oors = []float64{0, 0, 0}
	B.Oors = []float64{0, 0, 0}
	A_dominates, B_dominates = IndCompare(A, B)
	io.Pfblue2("\nA: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	io.Pforan("B_dominates = %v\n", B_dominates)
	if A_dominates {
		tst.Errorf("test failed\n")
		return
	}
	if B_dominates {
		tst.Errorf("test failed\n")
		return
	}

	A.Ovas = []float64{200, 100}
	A_dominates, B_dominates = IndCompare(A, B)
	io.Pfblue2("\nA: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	io.Pforan("B_dominates = %v\n", B_dominates)
	if A_dominates {
		tst.Errorf("test failed\n")
		return
	}
	if B_dominates {
		tst.Errorf("test failed\n")
		return
	}

	A.Ovas = []float64{200, 99}
	A_dominates, B_dominates = IndCompare(A, B)
	io.Pfblue2("\nA: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	io.Pforan("B_dominates = %v\n", B_dominates)
	if !A_dominates {
		tst.Errorf("test failed\n")
		return
	}
	if B_dominates {
		tst.Errorf("test failed\n")
		return
	}

	A.Ovas = []float64{200, 100}
	B.Ovas = []float64{199, 100}
	A_dominates, B_dominates = IndCompare(A, B)
	io.Pfblue2("\nA: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	io.Pforan("B_dominates = %v\n", B_dominates)
	if A_dominates {
		tst.Errorf("test failed\n")
		return
	}
	if !B_dominates {
		tst.Errorf("test failed\n")
		return
	}
}

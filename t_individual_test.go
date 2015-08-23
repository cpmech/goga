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

func get_individual(id, nbases int) (ind *Individual) {
	nova := 2
	noor := 3
	switch id {
	case 0:
		ind = NewIndividual(nova, noor, nbases,
			[]int{1, 20, 300},
			[]float64{4.4, 5.5, 666},
			[]string{"abc", "b", "c"},
			[]byte("SGA"),
			[][]byte{[]byte("ABC"), []byte("DEF"), []byte("GHI")},
			[]Func_t{
				func(g *Individual) string { return "f0" },
				func(g *Individual) string { return "f1" },
				func(g *Individual) string { return "f2" },
			},
		)
		ind.Ovas[0] = 123
		ind.Ovas[1] = 345
		ind.Oors[0] = 10
		ind.Oors[1] = 20
		ind.Oors[2] = 30
	case 1:
		ind = NewIndividual(nova, noor, nbases,
			[]int{-1, -20, -300},
			[]float64{104.4, 105.5, 6.66},
			[]string{"X", "Y", "Z"},
			[]byte("#.#"),
			[][]byte{[]byte("^.^"), []byte("-o-"), []byte("*|*")},
			[]Func_t{
				func(g *Individual) string { return "g0" },
				func(g *Individual) string { return "g1" },
				func(g *Individual) string { return "g2" },
			},
		)
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

	nbases := 3
	A := get_individual(0, nbases)
	B := A.GetCopy()
	chk.Scalar(tst, "ova0", 1e-17, B.Ovas[0], 123)
	chk.Scalar(tst, "ova1", 1e-17, B.Ovas[1], 345)
	chk.Scalar(tst, "oor0", 1e-17, B.Oors[0], 10)
	chk.Scalar(tst, "oor1", 1e-17, B.Oors[1], 20)
	chk.Scalar(tst, "oor2", 1e-17, B.Oors[2], 30)

	fmts := [][]string{{" %d"}, {" %.1f"}, {" %q"}, {" %x"}, {" %q"}, {" %q"}}
	oA := A.Output(fmts, false)
	oB := B.Output(fmts, false)
	io.Pfyel("\n%v\n", oA)
	io.Pfyel("%v\n\n", oB)
	chk.String(tst, oA, " 1 20 300 4.4 5.5 666.0 \"abc\" \"b\" \"c\" 53 47 41 \"ABC\" \"DEF\" \"GHI\" \"f0\" \"f1\" \"f2\"")
	chk.String(tst, oB, " 1 20 300 4.4 5.5 666.0 \"abc\" \"b\" \"c\" 53 47 41 \"ABC\" \"DEF\" \"GHI\" \"f0\" \"f1\" \"f2\"")

	A.SetFloat(1, 33)
	A.SetFloat(2, 88)
	oA = A.Output(fmts, false)
	io.Pfyel("\n%v\n", oA)
	chk.String(tst, oA, " 1 20 300 4.4 33.0 88.0 \"abc\" \"b\" \"c\" 53 47 41 \"ABC\" \"DEF\" \"GHI\" \"f0\" \"f1\" \"f2\"")
}

func Test_ind02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind02. copy into")

	rnd.Init(0)

	nbases := 1
	A := get_individual(0, nbases)
	B := get_individual(1, nbases)

	fmts := [][]string{
		{"%2d", "%4d", "%5d"}, // ints
		{"%6g", "%6g", "%5g"}, // floats
		{"%4s", "%2s", "%2s"}, // strings
		{"%3x", "%3x", "%3x"}, // keys
		{"%4s", "%4s", "%4s"}, // bytes
		{"%3s", "%3s", "%3s"}, // funcs
	}
	io.Pfpink("A = %v\n", A.Output(fmts, false))
	io.Pfcyan("B = %v\n", B.Output(fmts, false))

	cuts := map[string][]int{
		"int": []int{1, 2},
		"str": []int{1},
	}
	pc := map[string]float64{
		"int": 1,
		"str": 1,
	}

	a := A.GetCopy()
	b := A.GetCopy()
	IndCrossover(a, b, A, B, nil, cuts, pc, nil, nil, nil, nil, nil, nil)

	io.Pforan("a = %v\n", a.Output(fmts, false))
	io.Pfblue2("b = %v\n", b.Output(fmts, false))

	chk.Ints(tst, "a.Ints   ", a.Ints, []int{1, -20, 300})
	chk.Ints(tst, "b.Ints   ", b.Ints, []int{-1, 20, -300})
	chk.Strings(tst, "a.Strings", a.Strings, []string{"abc", "Y", "Z"})
	chk.Strings(tst, "b.Strings", b.Strings, []string{"X", "b", "c"})
	// TODO: add other tests here
	io.Pf("\n")

	x := get_individual(0, nbases)
	x.Ovas = []float64{0, 0}
	x.Oors = []float64{0, 0, 0}
	io.Pfblue2("x = %v\n", x.Output(fmts, false))
	B.CopyInto(x)

	chk.Scalar(tst, "ova0", 1e-17, x.Ovas[0], 200)
	chk.Scalar(tst, "ova1", 1e-17, x.Ovas[1], 100)
	chk.Scalar(tst, "oor0", 1e-17, x.Oors[0], 15)
	chk.Scalar(tst, "oor1", 1e-17, x.Oors[1], 25)
	chk.Scalar(tst, "oor2", 1e-17, x.Oors[2], 35)

	io.Pforan("x = %v\n", x.Output(fmts, false))
	chk.String(tst, x.Output(fmts, false), B.Output(fmts, false))
}

func Test_ind03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind03. comparing")

	rnd.Init(1113)

	nbases := 1
	A := get_individual(0, nbases)
	B := get_individual(1, nbases)
	A_dominates := IndCompare(A, B, 0)
	io.Pfblue2("A: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	if !A_dominates {
		tst.Errorf("test failed\n")
		return
	}

	A.Oors = []float64{0, 0, 0}
	B.Oors = []float64{0, 0, 0}
	A_dominates = IndCompare(A, B, 0)
	io.Pfblue2("\nA: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	if A_dominates {
		tst.Errorf("test failed\n")
		return
	}

	A.Ovas = []float64{200, 100}
	A_dominates = IndCompare(A, B, 0)
	io.Pfblue2("\nA: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	if A_dominates {
		tst.Errorf("test failed\n")
		return
	}

	A.Ovas = []float64{200, 99}
	A_dominates = IndCompare(A, B, 0)
	io.Pfblue2("\nA: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	if !A_dominates {
		tst.Errorf("test failed\n")
		return
	}

	A.Ovas = []float64{200, 100}
	B.Ovas = []float64{199, 100}
	A_dominates = IndCompare(A, B, 0)
	io.Pfblue2("\nA: ovas = %v\n", A.Ovas)
	io.Pfblue2("A: oors = %v\n", A.Oors)
	io.Pfcyan("B: ovas = %v\n", B.Ovas)
	io.Pfcyan("B: oors = %v\n", B.Oors)
	io.Pforan("A_dominates = %v\n", A_dominates)
	if A_dominates {
		tst.Errorf("test failed\n")
		return
	}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func get_individual(id, nbases int) *Individual {
	switch id {
	case 0:
		return NewIndividual(nbases,
			[]int{1, 20, 300},
			[]float64{4.4, 5.5, 666},
			[]string{"abc", "b", "c"},
			[]byte("SGA"),
			[][]byte{[]byte("ABC"), []byte("DEF"), []byte("GHI")},
			[]Func_tt{
				func(g *Individual) string { return "f0" },
				func(g *Individual) string { return "f1" },
				func(g *Individual) string { return "f2" },
			},
		)
	case 1:
		return NewIndividual(nbases,
			[]int{-1, -20, -300},
			[]float64{104.4, 105.5, 6.66},
			[]string{"XX", "YY", "ZZ"},
			[]byte("#.#"),
			[][]byte{[]byte("^.^"), []byte("-o-"), []byte("*|*")},
			[]Func_tt{
				func(g *Individual) string { return "g0" },
				func(g *Individual) string { return "g1" },
				func(g *Individual) string { return "g2" },
			},
		)
	}
	return nil
}

func Test_ind01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind01")

	nbases := 3
	A := get_individual(0, nbases)
	B := A.GetCopy()

	oA := A.Output(nil)
	oB := B.Output(nil)
	io.Pfyel("\n%v\n", oA)
	io.Pfyel("%v\n\n", oB)
	chk.String(tst, oA, " 1 20 300 4.4 5.5 666 \"abc\" \"b\" \"c\" 53 47 41 \"ABC\" \"DEF\" \"GHI\" \"f0\" \"f1\" \"f2\"")
	chk.String(tst, oB, " 1 20 300 4.4 5.5 666 \"abc\" \"b\" \"c\" 53 47 41 \"ABC\" \"DEF\" \"GHI\" \"f0\" \"f1\" \"f2\"")
}

func Test_ind02(tst *testing.T) {

	verbose()
	chk.PrintTitle("ind02")

	nbases := 3
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
	io.Pfpink("A = %v\n", A.Output(fmts))
	io.Pfcyan("B = %v\n", B.Output(fmts))

	cuts := [][]int{
		{-1, -1}, // ints
		{-1, -1}, // floats
		{-1, -1}, // strings
		{-1, -1}, // keys
		{-1, -1}, // bytes
		{-1, -1}, // funcs
	}

	a := A.GetCopy()
	b := A.GetCopy()
	Crossover(a, b, A, B, nil, cuts, nil, nil, nil, nil, nil, nil, nil)

	io.Pforan("a = %v\n", a.Output(fmts))
	io.Pfblue2("b = %v\n", b.Output(fmts))
}

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

func get_individual(id, nbases int) *Individual {
	switch id {
	case 0:
		return NewIndividual(nbases,
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
	case 1:
		return NewIndividual(nbases,
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
	}
	return nil
}

func Test_ind01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind01")

	rnd.Init(0)

	nbases := 3
	A := get_individual(0, nbases)
	B := A.GetCopy()

	fmts := [][]string{{" %d"}, {" %.1f"}, {" %q"}, {" %x"}, {" %q"}, {" %q"}}
	oA := A.Output(fmts)
	oB := B.Output(fmts)
	io.Pfyel("\n%v\n", oA)
	io.Pfyel("%v\n\n", oB)
	chk.String(tst, oA, " 1 20 300 4.4 5.5 666.0 \"abc\" \"b\" \"c\" 53 47 41 \"ABC\" \"DEF\" \"GHI\" \"f0\" \"f1\" \"f2\"")
	chk.String(tst, oB, " 1 20 300 4.4 5.5 666.0 \"abc\" \"b\" \"c\" 53 47 41 \"ABC\" \"DEF\" \"GHI\" \"f0\" \"f1\" \"f2\"")

	A.SetFloat(1, 33)
	A.SetFloat(2, 88)
	oA = A.Output(fmts)
	io.Pfyel("\n%v\n", oA)
	chk.String(tst, oA, " 1 20 300 4.4 33.0 88.0 \"abc\" \"b\" \"c\" 53 47 41 \"ABC\" \"DEF\" \"GHI\" \"f0\" \"f1\" \"f2\"")
}

func Test_ind02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind02")

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
	io.Pfpink("A = %v\n", A.Output(fmts))
	io.Pfcyan("B = %v\n", B.Output(fmts))

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
	Crossover(a, b, A, B, nil, cuts, pc, nil, nil, nil, nil, nil, nil)

	io.Pforan("a = %v\n", a.Output(fmts))
	io.Pfblue2("b = %v\n", b.Output(fmts))

	chk.Ints(tst, "a.Ints   ", a.Ints, []int{1, -20, 300})
	chk.Ints(tst, "b.Ints   ", b.Ints, []int{-1, 20, -300})
	chk.Strings(tst, "a.Strings", a.Strings, []string{"abc", "Y", "Z"})
	chk.Strings(tst, "b.Strings", b.Strings, []string{"X", "b", "c"})
	// TODO: add other tests here
	io.Pf("\n")

	x := get_individual(0, nbases)
	io.Pfblue2("x = %v\n", x.Output(fmts))
	B.CopyInto(x)
	io.Pforan("x = %v\n", x.Output(fmts))
	chk.String(tst, x.Output(fmts), B.Output(fmts))
}

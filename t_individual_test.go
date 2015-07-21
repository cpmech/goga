// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func get_individual(id, nbases int) *Individual {
	var ind Individual
	switch id {
	case 0:
		ind.InitChromo(nbases,
			[]int{1, 20, 300},
			[]float64{4.4, 5.5, 666},
			[]string{"abc", "b", "c"},
			[]byte("SGA"),
			[][]byte{[]byte("ABC"), []byte("DEF"), []byte("GHI")},
			[]Func_t{
				func(g *Gene) string { return "f0" },
				func(g *Gene) string { return "f1" },
				func(g *Gene) string { return "f2" },
			},
		)
	case 1:
		ind.InitChromo(nbases,
			[]int{-1, -20, -300},
			[]float64{104.4, 105.5, 6.66},
			[]string{"XX", "YY", "ZZ"},
			[]byte("#.#"),
			[][]byte{[]byte("^.^"), []byte("-o-"), []byte("*|*")},
			[]Func_t{
				func(g *Gene) string { return "g0" },
				func(g *Gene) string { return "g1" },
				func(g *Gene) string { return "g2" },
			},
		)
	}
	return &ind
}

func Test_ind01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind01")

	nbases := 3
	ind := get_individual(0, nbases)

	fmts := [][]string{
		{" %d", " %g", " %q", " %x", " %q", " %q"}, // use for all genes
	}
	out := ind.Output(fmts)
	io.Pfyel("\n%v\n\n", out)
	chk.String(tst, out, "[ 1 4.4 \"abc\" 53 \"ABC\" \"f0\"] [ 20 5.5 \"b\" 47 \"DEF\" \"f1\"] [ 300 666 \"c\" 41 \"GHI\" \"f2\"]")
}

func Test_ind02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind02")

	nbases := 3
	ind := get_individual(0, nbases)

	nint, nflt, nstr, nbyt, nbytes, nfuncs := ind.CountBases()

	ints := make([]int, nint)
	flts := make([]float64, nflt)
	strs := make([]string, nstr)
	byts := make([]byte, nbyt)
	bytes := make([]byte, nbytes)
	funcs := make([]Func_t, nfuncs)

	ind.GetBases(ints, flts, strs, byts, bytes, funcs)

	io.Pforan("ints  = %v\n", ints)
	io.Pforan("flts  = %v\n", flts)
	io.Pforan("strs  = %v\n", strs)
	io.Pforan("byts  = %v\n", byts)
	io.Pforan("bytes = %v\n", bytes)
	io.Pforan("funcs = %v\n", funcs)

	chk.Ints(tst, "ints", ints, []int{1, 20, 300})
	chk.Strings(tst, "strs", strs, []string{"abc", "b", "c"})
	// TODO: add other checks

	oth := get_individual(1, nbases)
	oth.SetBases(ints, flts, strs, byts, bytes, funcs)

	fmts := [][]string{
		{"%4d", " %5g", " %3s", " %x", " %3s", " %3s"}, // use for all genes
	}
	io.Pfpink("ind = %v\n", ind.Output(fmts))
	io.Pfcyan("oth = %v\n", oth.Output(fmts))
	for i, g := range ind.Chromo {
		if *g.Int != *oth.Chromo[i].Int {
			tst.Errorf("int: individuals are different\n")
			return
		}
		if math.Abs(*g.Flt-*oth.Chromo[i].Flt) > 1e-12 {
			tst.Errorf("flt: individuals are different. diff = %v", math.Abs(*g.Flt-*oth.Chromo[i].Flt))
			return
		}
		if *g.String != *oth.Chromo[i].String {
			tst.Errorf("str: individuals are different\n")
			return
		}
		if *g.Byte != *oth.Chromo[i].Byte {
			tst.Errorf("byte: individuals are different\n")
			return
		}
		if string(g.Bytes) != string(oth.Chromo[i].Bytes) {
			tst.Errorf("bytes: individuals are different\n")
			return
		}
		if g.Func(g) != oth.Chromo[i].Func(oth.Chromo[i]) {
			tst.Errorf("func: individuals are different\n")
			return
		}
	}
}

func Test_ind03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind03")

	nbases := 3
	A := get_individual(0, nbases)
	B := get_individual(1, nbases)

	fmts := [][]string{
		{"%4d", " %5g", " %3s", " %x", " %3s", " %3s"}, // use for all genes
	}

	cuts := []int{0}
	scuts := []int{-1}

	a := A.GetCopy()
	b := A.GetCopy()
	Crossover(a, b, A, B, cuts, scuts)

	io.Pfpink("A = %v\n", A.Output(fmts))
	io.Pfcyan("B = %v\n", B.Output(fmts))
	io.Pforan("a = %v\n", a.Output(fmts))
	io.Pfblue2("b = %v\n", b.Output(fmts))
}

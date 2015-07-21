// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func get_individual(nbases int) *Individual {
	var ind Individual
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
	return &ind
}

func Test_ind01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind01")

	nbases := 3
	ind := get_individual(nbases)

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
	ind := get_individual(nbases)

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
}

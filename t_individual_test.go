// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_ind01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ind01")

	var ind Individual
	nbases := 3
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

	out := ind.Output([]string{"%d", "%g", "%s", "%q"})
	io.Pfyel("\n%v\n\n", out)
	chk.String(tst, out, "(1,4.4,abc,53,\"ABC\",f0) (20,5.5,b,47,\"DEF\",f1) (300,666,c,41,\"GHI\",f2)")
}

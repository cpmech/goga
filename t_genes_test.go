// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_gene01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("gene01")

	nbases := 3
	g0 := NewGene(nbases)

	g0.SetInt(123)
	g0.SetFloat(666)
	g0.SetString("abc")
	g0.SetByte('S')
	g0.SetBytes([]byte("ABC"))
	g0.SetFunc(func(g *Gene) string { return "hello" })

	g1 := g0.GetCopy()

	fmts := []string{"%d", "%g", "%s", "%q"}
	r0 := g0.Output(fmts)
	r1 := g1.Output(fmts)
	io.Pforan("g0 = %s\n", r0)
	io.Pfcyan("g1 = %s\n\n", r0)

	chk.String(tst, r0, "123,666,abc,53,\"ABC\",hello")
	chk.String(tst, r1, "123,666,abc,53,\"ABC\",hello")
	chk.Vector(tst, "subfloats", 1e-17, g0.Fbases, g1.Fbases)
}

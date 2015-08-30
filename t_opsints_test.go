// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"sort"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_cxint01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cxint01")

	var ops OpsData
	ops.SetDefault()
	ops.Pc = 1
	ops.Ncuts = 1

	A := []int{1, 2}
	B := []int{-1, -2}
	a := make([]int, len(A))
	b := make([]int, len(A))
	IntCrossover(a, b, A, B, 0, &ops)
	io.Pfred("A = %2d\n", A)
	io.PfRed("B = %2d\n", B)
	io.Pfcyan("a = %2d\n", a)
	io.Pfblue2("b = %2d\n", b)
	chk.Ints(tst, "a", a, []int{1, -2})
	chk.Ints(tst, "b", b, []int{-1, 2})
	io.Pf("\n")

	A = []int{1, 2, 3, 4, 5, 6, 7, 8}
	B = []int{-1, -2, -3, -4, -5, -6, -7, -8}
	a = make([]int, len(A))
	b = make([]int, len(A))
	ops.Cuts = []int{1, 3}
	IntCrossover(a, b, A, B, 0, &ops)
	io.Pfred("A = %2v\n", A)
	io.PfRed("B = %2v\n", B)
	io.Pfcyan("a = %2v\n", a)
	io.Pfblue2("b = %2v\n", b)
	chk.Ints(tst, "a", a, []int{1, -2, -3, 4, 5, 6, 7, 8})
	chk.Ints(tst, "b", b, []int{-1, 2, 3, -4, -5, -6, -7, -8})

	ops.Cuts = []int{5, 7}
	IntCrossover(a, b, A, B, 0, &ops)
	io.Pfred("A = %2v\n", A)
	io.PfRed("B = %2v\n", B)
	io.Pfcyan("a = %2v\n", a)
	io.Pfblue2("b = %2v\n", b)
	chk.Ints(tst, "a", a, []int{1, 2, 3, 4, 5, -6, -7, 8})
	chk.Ints(tst, "b", b, []int{-1, -2, -3, -4, -5, 6, 7, -8})
}

func Test_intordcx01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("intordcx01")

	var ops OpsData
	ops.SetDefault()
	ops.Pc = 1

	rnd.Init(0)

	A := []int{1, 2, 3, 4, 5, 6, 7, 8}
	B := []int{2, 4, 6, 8, 7, 5, 3, 1}
	a := make([]int, len(A))
	b := make([]int, len(A))
	ops.Cuts = []int{2, 5}
	IntOrdCrossover(a, b, A, B, 0, &ops)
	io.Pforan("A = %v\n", A)
	io.Pfblue2("B = %v\n", B)
	io.Pfgreen("a = %v\n", a)
	io.Pfyel("b = %v\n", b)
	chk.Ints(tst, "A", A, []int{1, 2, 3, 4, 5, 6, 7, 8})
	chk.Ints(tst, "B", B, []int{2, 4, 6, 8, 7, 5, 3, 1})
	chk.Ints(tst, "a", a, []int{4, 5, 6, 8, 7, 1, 2, 3})
	chk.Ints(tst, "b", b, []int{8, 7, 3, 4, 5, 1, 2, 6})
	sort.Ints(a)
	sort.Ints(b)
	nums := utl.IntRange2(1, 9)
	chk.Ints(tst, "asorted = 12345678", a, nums)
	chk.Ints(tst, "bsorted = 12345678", b, nums)

	A = []int{1, 3, 5, 7, 6, 2, 4, 8}
	B = []int{5, 6, 3, 8, 2, 1, 4, 7}
	ops.Cuts = []int{3, 6}
	IntOrdCrossover(a, b, A, B, 0, &ops)
	io.Pforan("\nA = %v\n", A)
	io.Pfblue2("B = %v\n", B)
	io.Pfgreen("a = %v\n", a)
	io.Pfyel("b = %v\n", b)
	chk.Ints(tst, "A", A, []int{1, 3, 5, 7, 6, 2, 4, 8})
	chk.Ints(tst, "B", B, []int{5, 6, 3, 8, 2, 1, 4, 7})
	chk.Ints(tst, "a", a, []int{5, 7, 6, 8, 2, 1, 4, 3})
	chk.Ints(tst, "b", b, []int{3, 8, 1, 7, 6, 2, 4, 5})
	sort.Ints(a)
	sort.Ints(b)
	chk.Ints(tst, "asorted = 12345678", a, nums)
	chk.Ints(tst, "bsorted = 12345678", b, nums)

	A = []int{1, 2, 3, 4, 5, 6, 7, 8}
	B = []int{2, 4, 6, 8, 7, 5, 3, 1}
	ops.Cuts = []int{}
	IntOrdCrossover(a, b, A, B, 0, &ops)
	io.Pforan("\nA = %v\n", A)
	io.Pfblue2("B = %v\n", B)
	io.Pfgreen("a = %v\n", a)
	io.Pfyel("b = %v\n", b)
	sort.Ints(a)
	sort.Ints(b)
	chk.Ints(tst, "asorted = 12345678", a, nums)
	chk.Ints(tst, "bsorted = 12345678", b, nums)

	C := []int{1, 2, 3}
	D := []int{3, 1, 2}
	c := make([]int, len(C))
	d := make([]int, len(D))
	IntOrdCrossover(c, d, C, D, 0, &ops)
	io.Pforan("\nC = %v\n", C)
	io.Pfblue2("D = %v\n", D)
	io.Pfgreen("c = %v\n", c)
	io.Pfyel("d = %v\n", d)
	chk.Ints(tst, "c", c, []int{2, 1, 3})
	chk.Ints(tst, "d", d, []int{1, 2, 3})
	sort.Ints(c)
	sort.Ints(d)
	chk.Ints(tst, "csorted = 123", c, []int{1, 2, 3})
	chk.Ints(tst, "dsorted = 123", d, []int{1, 2, 3})
}

func Test_intordmut01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("intordmut01")

	var ops OpsData
	ops.SetDefault()
	ops.Pm = 1

	rnd.Init(0)

	a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	io.Pforan("before: a = %v\n", a)
	ops.OrdSti = []int{2, 5, 4}
	IntOrdMutation(a, 0, &ops)
	io.Pfcyan("after:  a = %v\n", a)
	chk.Ints(tst, "a", a, []int{1, 2, 6, 7, 3, 4, 5, 8})
	nums := utl.IntRange2(1, 9)
	sort.Ints(a)
	chk.Ints(tst, "asorted = 12345678", a, nums)

	a = []int{1, 2, 3, 4, 5, 6, 7, 8}
	io.Pforan("\nbefore: a = %v\n", a)
	ops.OrdSti = nil
	IntOrdMutation(a, 0, &ops)
	io.Pfcyan("after:  a = %v\n", a)
	sort.Ints(a)
	chk.Ints(tst, "asorted = 12345678", a, nums)
}

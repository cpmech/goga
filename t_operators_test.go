// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"sort"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func Test_simplechromo01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("simplechromo01")

	rnd.Init(0)
	nbases := 2
	for i := 0; i < 10; i++ {
		chromo := SimpleChromo([]float64{1, 10, 100}, nbases)
		io.Pforan("chromo = %v\n", chromo)
		chk.IntAssert(len(chromo), 3*nbases)
		chk.Scalar(tst, "gene0", 1e-14, chromo[0]+chromo[1], 1)
		chk.Scalar(tst, "gene1", 1e-14, chromo[2]+chromo[3], 10)
		chk.Scalar(tst, "gene2", 1e-13, chromo[4]+chromo[5], 100)
	}
}

func Test_rws01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("rws01. roulette whell selection")

	f := []float64{2.0, 1.8, 1.6, 1.4, 1.2, 1.0, 0.8, 0.6, 0.4, 0.2, 0.0}
	n := len(f)
	p := make([]float64, n)
	sum := la.VecAccum(f)
	for i := 0; i < n; i++ {
		p[i] = f[i] / sum
	}
	cs := make([]float64, len(p))
	utl.CumSum(cs, p)
	selinds := make([]int, 6)
	RouletteSelect(selinds, cs, []float64{0.81, 0.32, 0.96, 0.01, 0.65, 0.42})
	io.Pforan("selinds = %v\n", selinds)
	chk.Ints(tst, "selinds", selinds, []int{5, 1, 8, 0, 4, 2})
}

func Test_sus01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sus01. stochastic-universal-sampling")

	f := []float64{2.0, 1.8, 1.6, 1.4, 1.2, 1.0, 0.8, 0.6, 0.4, 0.2, 0.0}
	n := len(f)
	p := make([]float64, n)
	sum := la.VecAccum(f)
	for i := 0; i < n; i++ {
		p[i] = f[i] / sum
	}
	cs := make([]float64, len(p))
	utl.CumSum(cs, p)
	selinds := make([]int, 6)
	SUSselect(selinds, cs, 0.1)
	io.Pforan("selinds = %v\n", selinds)
	chk.Ints(tst, "selinds", selinds, []int{0, 1, 2, 3, 5, 7})
}

func Test_pairs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("pairs01")

	selinds := []int{11, 9, 1, 10, 0, 8, 7, 7, 3, 5, 4, 2, 6, 6, 6, 12}
	ninds := len(selinds)
	A := make([]int, ninds/2)
	B := make([]int, ninds/2)
	FilterPairs(A, B, selinds)
	m, n := B[3], B[6]
	io.Pforan("A = %v\n", A)
	io.Pforan("B = %v\n", B)
	chk.Ints(tst, "A", A, []int{11, 1, 0, 7, 3, 4, 6, 6})
	chk.Ints(tst, "B", B, []int{9, 10, 8, m, 5, 2, n, 12})
	for i, a := range A {
		if B[i] == a {
			tst.Errorf("there are repeated values in A and B: a=%d, b=%d", a, B[i])
			return
		}
	}
}

func Test_ends01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ends01")

	size := 8
	cuts := []int{5, 7}
	ends := GenerateCxEnds(size, 0, cuts)
	io.Pfpink("size=%v cuts=%v\n", size, cuts)
	io.Pfyel("ends = %v\n", ends)
	chk.IntAssert(len(ends), 3)
	chk.IntAssert(ends[len(ends)-1], size)
	checkRepeated(ends)
	io.Pf("\n")

	size = 2
	cuts = []int{}
	ends = GenerateCxEnds(size, 0, cuts)
	io.Pfpink("size=%v cuts=%v\n", size, cuts)
	io.Pforan("ends = %v\n", ends)
	chk.Ints(tst, "ends", ends, []int{1, 2})
	chk.IntAssert(ends[len(ends)-1], size)
	checkRepeated(ends)
	io.Pf("\n")

	size = 3
	ncuts := 3
	ends = GenerateCxEnds(size, ncuts, nil)
	io.Pfpink("size=%v ncuts=%v\n", size, ncuts)
	io.Pforan("ends = %v\n", ends)
	chk.Ints(tst, "ends", ends, []int{1, 2, 3})
	chk.IntAssert(ends[len(ends)-1], size)
	checkRepeated(ends)
	io.Pf("\n")

	size = 3
	ncuts = 2
	ends = GenerateCxEnds(size, ncuts, nil)
	io.Pfpink("size=%v ncuts=%v\n", size, ncuts)
	io.Pforan("ends = %v\n", ends)
	chk.Ints(tst, "ends", ends, []int{1, 2, 3})
	chk.IntAssert(ends[len(ends)-1], size)
	checkRepeated(ends)
	io.Pf("\n")

	size = 8
	cuts = []int{7}
	ends = GenerateCxEnds(size, 0, cuts)
	io.Pfpink("size=%v cuts=%v\n", size, cuts)
	io.Pforan("ends = %v\n", ends)
	chk.Ints(tst, "ends", ends, []int{7, 8})
	chk.IntAssert(ends[len(ends)-1], size)
	checkRepeated(ends)
	io.Pf("\n")

	size = 8
	cuts = []int{2, 5}
	ends = GenerateCxEnds(size, 0, cuts)
	io.Pfpink("size=%v cuts=%v\n", size, cuts)
	io.Pforan("ends = %v\n", ends)
	chk.Ints(tst, "ends", ends, []int{2, 5, 8})
	chk.IntAssert(ends[len(ends)-1], size)
	checkRepeated(ends)
	io.Pf("\n")

	size = 20
	cuts = []int{1, 5, 15, 17}
	ends = GenerateCxEnds(size, 0, cuts)
	io.Pfpink("size=%v cuts=%v\n", size, cuts)
	io.Pfyel("ends = %v\n", ends)
	chk.Ints(tst, "ends", ends, []int{1, 5, 15, 17, 20})
	chk.IntAssert(ends[len(ends)-1], size)
	checkRepeated(ends)
	io.Pf("\n")

	size = 20
	ncuts = 5
	ends = GenerateCxEnds(size, ncuts, cuts)
	io.Pfpink("size=%v cuts=%v\n", size, cuts)
	io.Pfyel("ends = %v\n", ends)
	chk.IntAssert(ends[len(ends)-1], size)
	checkRepeated(ends)
	io.Pf("\n")
}

func Test_ends02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ends02")

	rnd.Init(0)

	size := 20
	ncuts := 10
	nsamples := 100
	hist := rnd.IntHistogram{Stations: utl.IntRange(size + 3)}
	for i := 0; i < nsamples; i++ {
		ends := GenerateCxEnds(size, ncuts, nil)
		hist.Count(ends, false)
	}
	io.Pf("%s\n", rnd.TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}

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

func Test_cx01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cx01")

	var ops OpsData
	ops.SetDefault()
	ops.Pc = 1
	ops.Cuts = []int{2, 4}

	A := []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6}
	B := []float64{9.1, 9.2, 9.3, 9.4, 9.5, 9.6}
	a := make([]float64, len(A))
	b := make([]float64, len(A))
	FltCrossover(a, b, A, B, 0, &ops)
	io.Pfred("A = %3v\n", A)
	io.PfRed("B = %3v\n", B)
	io.Pfcyan("a = %3v\n", a)
	io.Pfblue2("b = %3v\n", b)
	chk.Vector(tst, "a", 1e-17, a, []float64{1.1, 2.2, 9.3, 9.4, 5.5, 6.6})
	chk.Vector(tst, "b", 1e-17, b, []float64{9.1, 9.2, 3.3, 4.4, 9.5, 9.6})
	io.Pf("\n")

	C := []string{"A", "B", "C", "D", "E", "F"}
	D := []string{"-", "o", "+", "@", "*", "&"}
	c := make([]string, len(A))
	d := make([]string, len(A))
	ops.Cuts = []int{1, 3}
	StrCrossover(c, d, C, D, 0, &ops)
	io.Pfred("C = %3v\n", C)
	io.PfRed("D = %3v\n", D)
	io.Pfcyan("c = %3v\n", c)
	io.Pfblue2("d = %3v\n", d)
	chk.Strings(tst, "c", c, []string{"A", "o", "+", "D", "E", "F"})
	chk.Strings(tst, "d", d, []string{"-", "B", "C", "@", "*", "&"})
	io.Pf("\n")

	E := [][]byte{[]byte("A"), []byte("B"), []byte("C"), []byte("D"), []byte("E"), []byte("F")}
	F := [][]byte{[]byte("-"), []byte("o"), []byte("+"), []byte("@"), []byte("*"), []byte("&")}
	e := make([][]byte, len(A))
	f := make([][]byte, len(A))
	for i := 0; i < len(A); i++ {
		e[i] = make([]byte, 1)
		f[i] = make([]byte, 1)
	}
	ops.Cuts = []int{1, 3}
	BytCrossover(e, f, E, F, 0, &ops)
	io.Pfred("E = %3s\n", E)
	io.PfRed("F = %3s\n", F)
	io.Pfcyan("e = %3s\n", e)
	io.Pfblue2("f = %3s\n", f)
	e_s := make([]string, len(A))
	f_s := make([]string, len(A))
	for i := 0; i < len(A); i++ {
		e_s[i] = string(e[i])
		f_s[i] = string(f[i])
	}
	chk.Strings(tst, "e_s", e_s, []string{"A", "o", "+", "D", "E", "F"})
	chk.Strings(tst, "f_s", f_s, []string{"-", "B", "C", "@", "*", "&"})
	io.Pf("\n")

	G := []byte("ABCDEF")
	H := []byte("-o+@*&")
	g := make([]byte, len(A))
	h := make([]byte, len(A))
	ops.Cuts = []int{1, 3}
	KeyCrossover(g, h, G, H, 0, &ops)
	io.Pfred("G = %3v\n", G)
	io.PfRed("H = %3v\n", H)
	io.Pfcyan("g = %3v\n", g)
	io.Pfblue2("h = %3v\n", h)
	g_s := make([]string, len(A))
	h_s := make([]string, len(A))
	for i := 0; i < len(A); i++ {
		g_s[i] = string(g[i])
		h_s[i] = string(h[i])
	}
	chk.Strings(tst, "g_s", g_s, []string{"A", "o", "+", "D", "E", "F"})
	chk.Strings(tst, "h_s", h_s, []string{"-", "B", "C", "@", "*", "&"})
	io.Pf("\n")

	M := []Func_t{func(i *Individual) string { return "A" }, func(i *Individual) string { return "B" }, func(i *Individual) string { return "C" }, func(i *Individual) string { return "D" }, func(i *Individual) string { return "E" }, func(i *Individual) string { return "F" }}
	N := []Func_t{func(i *Individual) string { return "-" }, func(i *Individual) string { return "o" }, func(i *Individual) string { return "+" }, func(i *Individual) string { return "@" }, func(i *Individual) string { return "*" }, func(i *Individual) string { return "&" }}
	m := make([]Func_t, len(A))
	n := make([]Func_t, len(A))
	ops.Cuts = []int{1, 3}
	FunCrossover(m, n, M, N, 0, &ops)
	io.Pfred("M = %3v\n", M)
	io.PfRed("N = %3v\n", N)
	io.Pfcyan("m = %3v\n", m)
	io.Pfblue2("n = %3v\n", n)
	m_s := make([]string, len(A))
	n_s := make([]string, len(A))
	for i := 0; i < len(A); i++ {
		m_s[i] = m[i](nil)
		n_s[i] = n[i](nil)
	}
	chk.Strings(tst, "m_s", m_s, []string{"A", "o", "+", "D", "E", "F"})
	chk.Strings(tst, "n_s", n_s, []string{"-", "B", "C", "@", "*", "&"})
	io.Pf("\n")
}

func checkRepeated(ends []int) {
	for i := 1; i < len(ends); i++ {
		if ends[i] == ends[i-1] {
			chk.Panic("there are repeated entries in ends = %v", ends)
		}
	}
}

func Test_mut01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mut01")

	var ops OpsData
	ops.SetDefault()
	ops.Pm = 1

	rnd.Init(0)

	A := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ops.Nchanges = 3
	ops.Mmax = 10
	io.Pforan("before: A = %v\n", A)
	IntMutation(A, 0, &ops)
	io.Pforan("after:  A = %v\n", A)
	io.Pf("\n")

	B := []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9}
	io.Pforan("before: B = %v\n", B)
	FltMutation(B, 0, &ops)
	io.Pforan("after:  B = %v\n", B)
	io.Pf("\n")

	C := []string{"a", "b", "c", "d", "e", "f"}
	ops.Nchanges = 2
	io.Pforan("before: C = %v\n", C)
	StrMutation(C, 0, &ops)
	io.Pforan("after:  C = %v\n", C)
	io.Pf("\n")

	D := []byte("abcdefghijklm")
	ops.Nchanges = 3
	io.Pforan("before: D = %s\n", D)
	KeyMutation(D, 0, &ops)
	io.Pforan("after:  D = %s\n", D)
	io.Pf("\n")

	E := [][]byte{[]byte("abc"), []byte("def"), []byte("ghi"), []byte("jkl")}
	io.Pforan("before: E = %s\n", E)
	BytMutation(E, 0, &ops)
	io.Pforan("after:  E = %s\n", E)
	io.Pf("\n")

	F := []Func_t{
		func(o *Individual) string { return "f0" },
		func(o *Individual) string { return "f1" },
		func(o *Individual) string { return "f2" },
		func(o *Individual) string { return "g0" },
		func(o *Individual) string { return "g1" },
		func(o *Individual) string { return "g2" },
	}
	io.Pforan("before: F =")
	for _, f := range F {
		io.Pforan(" %q", f(nil))
	}
	FunMutation(F, 0, &ops)
	io.Pforan("\nafter:  F =")
	for _, f := range F {
		io.Pforan(" %q", f(nil))
	}
	io.Pf("\n")
	io.Pf("\n")
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

func Test_blx01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("blx01. blended crossover")

	var ops OpsData
	ops.SetDefault()
	ops.Pc = 1.0
	ops.Xrange = [][]float64{{-1, 2}, {0, 3}, {1, 4}, {3, 6}, {4, 7}}

	rnd.Init(0)

	A := []float64{0, 1, 2, 4, 5}
	B := []float64{1, 2, 3, 5, 6}
	a := make([]float64, len(A))
	b := make([]float64, len(A))
	FltCrossoverBlx(a, b, A, B, 0, &ops)
	io.Pforan("A = %v\n", A)
	io.Pforan("B = %v\n", B)
	io.Pfcyan("a = %v\n", a)
	io.Pfcyan("b = %v\n", b)
}

func Test_mwicz01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mwicz01. Michalewicz mutation")

	var ops OpsData
	ops.SetDefault()
	ops.Pm = 1.0
	ops.Tmax = 10

	rnd.Init(0)

	ops.Xrange = [][]float64{{0, 2}, {1, 3}, {2, 4}, {3, 5}, {4, 6}}
	T := utl.IntRange(int(ops.Tmax))
	for _, t := range T {
		io.Pf("t=%v Δ=%v\n", t, ops.MwiczDelta(float64(t), 1))
	}
	for _, t := range T {
		A := []float64{0, 1, 2, 3, 4}
		FltMutationMwicz(A, t, &ops)
		io.Pforan("A = %.8f\n", A)
	}

	if chk.Verbose {
		b := 2.0
		f := func(r, tb float64) float64 {
			return math.Pow(r, math.Pow(1.0-tb, b))
		}
		np := 21
		r, tb := utl.MeshGrid2D(0, 1, 0, 1, np, np) // tb = t/tmax
		z := la.MatAlloc(np, np)
		for i := 0; i < np; i++ {
			for j := 0; j < np; j++ {
				z[i][j] = f(r[i][j], tb[i][j])
			}
		}
		plt.Surface(tb, r, z, "linewidth=0.8")
		plt.Gll("tb", "r", "")
		plt.SaveD("/tmp/goga", "test_mwicz01.eps")
		//plt.Show()
	}
}

func Test_binmut01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("binmut01. mutation: binary")

	var ops OpsData
	ops.SetDefault()
	ops.Pm = 1.0
	ops.Tmax = 10
	ops.Nchanges = 3

	rnd.Init(0)

	A := []int{0, 1, 1, 1, 0, 0, 1, 0, 1, 1}
	a := make([]int, len(A))
	copy(a, A)
	io.Pforan("before: A = %v\n", A)
	IntBinMutation(A, 0, &ops)
	io.Pforan("after:  A = %v\n", A)
	ndiff := 0
	for i := 0; i < len(A); i++ {
		if A[i] != a[i] {
			ndiff++
		}
	}
	io.Pforan("number of changes = %v\n", ndiff)
	if ndiff != ops.Nchanges {
		tst.Errorf("binary mutation failed\n")
	}
}

func Test_cxdeb01(tst *testing.T) {

	verbose()
	chk.PrintTitle("cxdeb01. Deb's crossover")

	var ops OpsData
	ops.SetDefault()
	ops.Pc = 1.0
	ops.Xrange = [][]float64{{-3, 3}, {-4, 4}}
	ops.EnfRange = true

	rnd.Init(0)

	A := []float64{-1, 1}
	B := []float64{1, 2}
	a := make([]float64, len(A))
	b := make([]float64, len(A))
	FltCrossoverDeb(a, b, A, B, 0, &ops)
	io.Pforan("A = %v\n", A)
	io.Pforan("B = %v\n", B)
	io.Pfcyan("a = %.6f\n", a)
	io.Pfcyan("b = %.6f\n", b)

	nsamples := 100
	a0s, a1s := make([]float64, nsamples), make([]float64, nsamples)
	b0s, b1s := make([]float64, nsamples), make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		FltCrossoverDeb(a, b, B, A, 0, &ops)
		a0s[i], a1s[i] = a[0], a[1]
		b0s[i], b1s[i] = b[0], b[1]
	}
	ha0 := rnd.Histogram{Stations: []float64{-4, -3.5, -3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1}}
	hb0 := rnd.Histogram{Stations: []float64{0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 5, 5.5, 6}}
	ha1 := rnd.Histogram{Stations: utl.LinSpace(-4, 4, 11)}
	hb1 := rnd.Histogram{Stations: utl.LinSpace(-4, 4, 11)}
	ha0.Count(a0s, true)
	hb0.Count(b0s, true)
	ha1.Count(a1s, true)
	hb1.Count(b1s, true)

	io.Pforan("\na0s\n")
	io.Pf("%s", rnd.TextHist(ha0.GenLabels("%.1f"), ha0.Counts, 60))
	io.Pforan("b0s\n")
	io.Pf("%s", rnd.TextHist(hb0.GenLabels("%.1f"), hb0.Counts, 60))

	io.Pforan("\na1s\n")
	io.Pf("%s", rnd.TextHist(ha1.GenLabels("%.1f"), ha1.Counts, 60))
	io.Pforan("b1s\n")
	io.Pf("%s", rnd.TextHist(hb1.GenLabels("%.1f"), hb1.Counts, 60))
}

func Test_mtdeb01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mtdeb01. Deb's mutation")

	var ops OpsData
	ops.SetDefault()
	ops.Pm = 1.0
	ops.Xrange = [][]float64{{-3, 3}, {-4, 4}}
	ops.EnfRange = true

	rnd.Init(0)

	A := []float64{-1, 1}
	io.Pforan("before: A = %v\n", A)
	FltMutationDeb(A, 10, &ops)
	io.Pforan("after:  A = %v\n", A)

	ha0 := rnd.Histogram{Stations: utl.LinSpace(-3, 3, 11)}

	nsamples := 100
	aa := make([]float64, len(A))
	a0s := make([]float64, nsamples)
	for _, t := range []int{0, 50, 100} {
		for i := 0; i < nsamples; i++ {
			copy(aa, A)
			FltMutationDeb(aa, t, &ops)
			a0s[i] = aa[0]
		}
		ha0.Count(a0s, true)
		io.Pf("\ntime = %d\n", t)
		io.Pf("%s", rnd.TextHist(ha0.GenLabels("%.1f"), ha0.Counts, 60))
	}
}

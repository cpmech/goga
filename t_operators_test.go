// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
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

func Test_fitness01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("fitness01")

	ovs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	f := make([]float64, len(ovs))
	Fitness(f, ovs)
	io.Pforan("f = %v\n", f)
	chk.Vector(tst, "f", 1e-15, f, utl.LinSpace(1, 0, 11))
}

func Test_ranking01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ranking01")

	f := Ranking(11, 2.0)
	io.Pforan("f = %v\n", f)
	chk.Vector(tst, "f", 1e-15, f, []float64{2, 1.8, 1.6, 1.4, 1.2, 1, 0.8, 0.6, 0.4, 0.2, 0})

	f = Ranking(11, 1.1)
	io.Pfblue2("f = %v\n", f)
	chk.Vector(tst, "f", 1e-15, f, []float64{1.1, 1.08, 1.06, 1.04, 1.02, 1, 0.98, 0.96, 0.94, 0.92, 0.9})
}

func Test_cumsum01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cumsum01")

	p := []float64{1, 2, 3, 4, 5}
	cs := make([]float64, len(p))
	CumSum(cs, p)
	io.Pforan("cs = %v\n", cs)
	chk.Vector(tst, "cumsum", 1e-17, cs, []float64{1, 3, 6, 10, 15})
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
	CumSum(cs, p)
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
	CumSum(cs, p)
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
	nsamples := 1000
	hist := rnd.IntHistogram{Stations: utl.IntRange(size + 3)}
	for i := 0; i < nsamples; i++ {
		ends := GenerateCxEnds(size, ncuts, nil)
		hist.Count(ends, false)
	}
	io.Pf("%s\n", rnd.TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}

func Test_cxint01(tst *testing.T) {

	verbose()
	chk.PrintTitle("cxint01")

	A := []int{1, 2}
	B := []int{-1, -2}
	a := make([]int, len(A))
	b := make([]int, len(A))
	IntCrossover(a, b, A, B, 1, nil, 1)
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
	cuts := []int{5, 5, 8}
	ends := IntCrossover(a, b, A, B, 0, cuts, 1)
	io.Pforan("ends = %v\n", ends)
	io.Pfred("A = %2v\n", A)
	io.PfRed("B = %2v\n", B)
	io.Pfcyan("a = %2v\n", a)
	io.Pfblue2("b = %2v\n", b)
	chk.Ints(tst, "a", a, A)
	chk.Ints(tst, "b", b, B)
}

/*
func Test_cxint03(tst *testing.T) {

	verbose()
	chk.PrintTitle("cxint03")

	A := []int{1, 2, 3, 4, 5, 6, 7, 8}
	B := []int{-1, -2, -3, -4, -5, -6, -7, -8}
	a := make([]int, len(A))
	b := make([]int, len(A))

	ends := []int{8, 8}
	IntCrossover(a, b, A, B, ends, 1)

	io.Pfred("A = %2v\n", A)
	io.PfRed("B = %2v\n", B)
	io.Pfcyan("a = %2v\n", a)
	io.Pfblue2("b = %2v\n", b)

	chk.Ints(tst, "a", a, A)
	chk.Ints(tst, "b", b, B)
}

func Test_cxint04(tst *testing.T) {

	verbose()
	chk.PrintTitle("cxint04")

	A := []int{1, 2, 3, 4, 5, 6, 7, 8}
	B := []int{-1, -2, -3, -4, -5, -6, -7, -8}
	a := make([]int, len(A))
	b := make([]int, len(A))

	ends := []int{1, 5, 8}
	IntCrossover(a, b, A, B, ends, 1)

	io.Pfred("A = %2v\n", A)
	io.PfRed("B = %2v\n", B)
	io.Pfcyan("a = %2v\n", a)
	io.Pfblue2("b = %2v\n", b)

	chk.Ints(tst, "a", a, []int{1, -2, -3, -4, -5, 6, 7, 8})
	chk.Ints(tst, "b", b, []int{-1, 2, 3, 4, 5, -6, -7, -8})
}
*/

func checkRepeated(ends []int) {
	for i := 1; i < len(ends); i++ {
		if ends[i] == ends[i-1] {
			chk.Panic("there are repeated entries in ends = %v", ends)
		}
	}
}

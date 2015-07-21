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

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

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

func checkRepeated(ends []int) {
	for i := 1; i < len(ends); i++ {
		if ends[i] == ends[i-1] {
			chk.Panic("there are repeated entries in ends = %v", ends)
		}
	}
}

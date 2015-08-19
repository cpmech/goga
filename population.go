// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"sort"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Population holds all individuals
type Population []*Individual

// generation functions
type PopIntGen_t func(pop Population, ninds, nbases int, noise float64, args interface{}, irange [][]int) Population     // generate population of integers
type PopOrdGen_t func(pop Population, ninds, nbases int, noise float64, args interface{}, nints int) Population          // generate population of ordered integers
type PopFltGen_t func(pop Population, ninds, nbases int, noise float64, args interface{}, frange [][]float64) Population // generate population of float point numbers
type PopStrGen_t func(pop Population, ninds, nbases int, noise float64, args interface{}, pool [][]string) Population    // generate population of strings
type PopKeyGen_t func(pop Population, ninds, nbases int, noise float64, args interface{}, pool [][]byte) Population      // generate population of keys (bytes)
type PopBytGen_t func(pop Population, ninds, nbases int, noise float64, args interface{}, pool [][]string) Population    // generate population of bytes
type PopFunGen_t func(pop Population, ninds, nbases int, noise float64, args interface{}, pool [][]Func_t) Population    // generate population of functions

// PopFltGen generates a population of individuals with float point numbers
// Notes: (1) ngenes = len(frange)
//        (2) this function can be used with existent population
func PopFltGen(pop Population, ninds, nbases int, noise float64, args interface{}, frange [][]float64) Population {
	o := pop
	if len(o) != ninds {
		o = make([]*Individual, ninds)
	}
	ngenes := len(frange)
	for i := 0; i < ninds; i++ {
		if o[i] == nil {
			o[i] = new(Individual)
		}
		if o[i].Nfltgenes != ngenes {
			o[i].Nfltgenes = ngenes
			o[i].Floats = make([]float64, ngenes)
		}
	}
	npts := int(math.Pow(float64(ninds), 1.0/float64(ngenes))) // num points in 'square' grid
	ntot := int(math.Pow(float64(npts), float64(ngenes)))      // total num of individuals in grid
	den := 1.0                                                 // denominator to calculate dx
	if npts > 1 {
		den = float64(npts - 1)
	}
	var lfto int // leftover, e.g. n % (nx*ny)
	var rdim int // reduced dimension, e.g. (nx*ny)
	var idx int  // index of gene in grid
	var dx, x, mul, xmin, xmax float64
	for i := 0; i < ninds; i++ {
		if i < ntot { // on grid
			lfto = i
			for j := 0; j < ngenes; j++ {
				rdim = int(math.Pow(float64(npts), float64(ngenes-1-j)))
				idx = lfto / rdim
				lfto = lfto % rdim
				xmin = frange[j][0]
				xmax = frange[j][1]
				dx = xmax - xmin
				x = xmin + float64(idx)*dx/den
				if noise > 0 {
					mul = rnd.Float64(0, noise)
					if rnd.FlipCoin(0.5) {
						x += mul * x
					} else {
						x -= mul * x
					}
					if x < xmin {
						x = xmin + (xmin - x)
					}
					if x > xmax {
						x = xmax - (x - xmax)
					}
				}
				o[i].SetFloat(j, x)
			}
		} else { // additional individuals
			for j := 0; j < ngenes; j++ {
				xmin = frange[j][0]
				xmax = frange[j][1]
				x = rnd.Float64(xmin, xmax)
				o[i].SetFloat(j, x)
			}
		}
	}
	return o
}

// PopOrdGen generates a population of individuals with ordered integers
// Notes: (1) ngenes = len(frange)
//        (2) this function can be used with existent population
func PopOrdGen(pop Population, ninds, nbases int, noise float64, args interface{}, nints int) Population {
	o := pop
	if len(o) != ninds {
		o = make([]*Individual, ninds)
	}
	ngenes := nints
	for i := 0; i < ninds; i++ {
		if o[i] == nil {
			o[i] = new(Individual)
		}
		if len(o[i].Ints) != ngenes {
			o[i].Ints = make([]int, ngenes)
		}
		for j := 0; j < nints; j++ {
			o[i].Ints[j] = j
		}
		rnd.IntShuffle(o[i].Ints)
	}
	return o
}

// methods of Population ///////////////////////////////////////////////////////////////////////////

// GetCopy returns a copy of this population
func (o Population) GetCopy() (pop Population) {
	ninds := len(o)
	pop = make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = o[i].GetCopy()
	}
	return
}

// Len returns the length of the population == number of individuals
func (o Population) Len() int {
	return len(o)
}

// Swap swaps two individuals
func (o Population) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Less returns true if 'i' is "less bad" than 'j'; therefore it can be used
// to sort the population in increasing order of demerits: from best to worst
func (o Population) Less(i, j int) bool {
	return o[i].Demerit < o[j].Demerit
}

// Sort sorts the population from best to worst individuals; i.e. decreasing fitness values
func (o *Population) Sort() {
	sort.Sort(o)
}

// Output generates a nice table with population data
//  Input:
//  fmts      -- [ngenes] formats for int, flt, string, byte, bytes, and func
//               use fmts == nil to choose default ones
//  showBases -- show bases, if any
func (o Population) Output(fmts [][]string, showBases bool) (buf *bytes.Buffer) {

	// check
	if len(o) < 1 {
		return
	}

	// compute sizes and generate formats list
	if fmts == nil {
		sizes := make([][]int, 6)
		for _, ind := range o {
			sz := ind.GetStringSizes()
			for i := 0; i < 6; i++ {
				if len(sizes[i]) == 0 {
					sizes[i] = make([]int, len(sz[i]))
				}
				for j, s := range sz[i] {
					sizes[i][j] = utl.Imax(sizes[i][j], s)
				}
			}
		}
		fmts = make([][]string, 6)
		for i, str := range []string{"d", "g", "s", "x", "s", "s"} {
			fmts[i] = make([]string, len(sizes[i]))
			for j, sz := range sizes[i] {
				fmts[i][j] = io.Sf("%%%d%s", sz+1, str)
			}
		}
	}

	// compute sizes of header items
	szova, szoor, szdem := 0, 0, 0
	for _, ind := range o {
		szova = utl.Imax(szova, len(io.Sf("%g", ind.Ova)))
		szoor = utl.Imax(szoor, len(io.Sf("%g", ind.Oor)))
		szdem = utl.Imax(szdem, len(io.Sf("%g", ind.Demerit)))
	}
	szova = utl.Imax(szova, 3) // 3 ==> len("Ova")
	szoor = utl.Imax(szoor, 3) // 3 ==> len("Oor")
	szdem = utl.Imax(szdem, 7) // 7 ==> len("Demerit")

	// print individuals
	fmtova := io.Sf("%%%d", szova+1)
	fmtoor := io.Sf("%%%d", szoor+1)
	fmtdem := io.Sf("%%%d", szdem+1)
	line, sza, szb := "", 0, 0
	for i, ind := range o {
		stra := io.Sf(fmtova+"g", ind.Ova)
		if ind.Oor > 0 {
			stra = io.Sf(fmtova+"s", "n/a")
			stra += io.Sf(fmtoor+"g", ind.Oor)
		} else {
			stra += io.Sf(fmtoor+"s", "n/a")
		}
		stra += io.Sf(fmtdem+"g", ind.Demerit) + " "
		strb := ind.Output(fmts, showBases)
		line += stra + strb + "\n"
		if i == 0 {
			sza, szb = len(stra), len(strb)
		}
	}

	// write to buffer
	fmtgenes := io.Sf(" %%%d.%ds\n", szb, szb)
	n := sza + szb
	buf = new(bytes.Buffer)
	io.Ff(buf, io.StrThickLine(n))
	io.Ff(buf, fmtova+"s", "Ova")
	io.Ff(buf, fmtoor+"s", "Oor")
	io.Ff(buf, fmtdem+"s", "Demerit")
	io.Ff(buf, fmtgenes, "Genes")
	io.Ff(buf, io.StrThinLine(n))
	io.Ff(buf, line)
	io.Ff(buf, io.StrThickLine(n))
	return
}

// OutFloatBases print bases of float genes
func (o Population) OutFloatBases(numFmt string) (l string) {
	for _, ind := range o {
		for _, val := range ind.Floats {
			l += io.Sf(numFmt, val)
		}
		l += "\n"
	}
	return
}

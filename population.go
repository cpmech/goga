// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"
	"math/rand"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Population holds all individuals
type Population []*Individual

// generation functions ////////////////////////////////////////////////////////////////////////////

// PopBinGen generates a population of binary numbers [0,1]
func PopBinGen(id int, C *ConfParams) Population {
	o := make([]*Individual, C.Ninds)
	genes := make([]int, C.NumInts)
	for i := 0; i < C.Ninds; i++ {
		for j := 0; j < C.NumInts; j++ {
			genes[j] = rand.Intn(2)
		}
		o[i] = NewIndividual(C.Nova, C.Noor, C.Nbases, genes)
	}
	return o
}

// PopOrdGen generates a population of individuals with ordered integers
// Notes: (1) ngenes = C.NumInts
func PopOrdGen(id int, C *ConfParams) Population {
	o := make([]*Individual, C.Ninds)
	ngenes := C.NumInts
	for i := 0; i < C.Ninds; i++ {
		o[i] = NewIndividual(C.Nova, C.Noor, C.Nbases, make([]int, ngenes))
		for j := 0; j < C.NumInts; j++ {
			o[i].Ints[j] = j
		}
		rnd.IntShuffle(o[i].Ints)
	}
	return o
}

// PopFltAndBinGen generates population with floats and integers
func PopFltAndBinGen(id int, C *ConfParams) Population {
	o := PopFltGen(id, C)
	for i := 0; i < C.Ninds; i++ {
		o[i].Ints = make([]int, C.NumInts)
		for j := 0; j < C.NumInts; j++ {
			o[i].Ints[j] = rand.Intn(2)
		}
	}
	return o
}

// PopFltGen generates a population of individuals with float point numbers
// Notes: (1) ngenes = len(C.RangeFlt)
func PopFltGen(id int, C *ConfParams) Population {
	o := make([]*Individual, C.Ninds)
	ngenes := len(C.RangeFlt)
	for i := 0; i < C.Ninds; i++ {
		o[i] = NewIndividual(C.Nova, C.Noor, C.Nbases, make([]float64, ngenes))
	}
	if C.Latin {
		K := rnd.LatinIHS(ngenes, C.Ninds, C.LatinDf)
		dx := make([]float64, ngenes)
		for i := 0; i < ngenes; i++ {
			dx[i] = (C.RangeFlt[i][1] - C.RangeFlt[i][0]) / float64(C.Ninds-1)
		}
		for i := 0; i < ngenes; i++ {
			for j := 0; j < C.Ninds; j++ {
				o[j].SetFloat(i, C.RangeFlt[i][0]+float64(K[i][j]-1)*dx[i])
			}
		}
		return o
	}
	npts := int(math.Pow(float64(C.Ninds), 1.0/float64(ngenes))) // num points in 'square' grid
	ntot := int(math.Pow(float64(npts), float64(ngenes)))        // total num of individuals in grid
	den := 1.0                                                   // denominator to calculate dx
	if npts > 1 {
		den = float64(npts - 1)
	}
	var lfto int // leftover, e.g. n % (nx*ny)
	var rdim int // reduced dimension, e.g. (nx*ny)
	var idx int  // index of gene in grid
	var dx, x, mul, xmin, xmax float64
	for i := 0; i < C.Ninds; i++ {
		if i < ntot { // on grid
			lfto = i
			for j := 0; j < ngenes; j++ {
				rdim = int(math.Pow(float64(npts), float64(ngenes-1-j)))
				idx = lfto / rdim
				lfto = lfto % rdim
				xmin = C.RangeFlt[j][0]
				xmax = C.RangeFlt[j][1]
				dx = xmax - xmin
				x = xmin + float64(idx+id)*dx/den
				if C.Noise > 0 {
					mul = rnd.Float64(0, C.Noise)
					if rnd.FlipCoin(0.5) {
						x += mul * x
					} else {
						x -= mul * x
					}
				}
				if x < xmin {
					x = xmin + (xmin - x)
				}
				if x > xmax {
					x = xmax - (x - xmax)
				}
				o[i].SetFloat(j, x)
			}
		} else { // additional individuals
			for j := 0; j < ngenes; j++ {
				xmin = C.RangeFlt[j][0]
				xmax = C.RangeFlt[j][1]
				x = rnd.Float64(xmin, xmax)
				o[i].SetFloat(j, x)
			}
		}
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

// Output generates a nice table with population data
func (o Population) Output(C *ConfParams) (buf *bytes.Buffer) {

	// check
	if len(o) < 1 {
		return
	}

	// compute sizes and generate formats list
	if C.NumFmts == nil {
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
		name := []string{"int", "flt", "str", "key", "byt", "fun"}
		C.NumFmts = make(map[string][]string)
		for i, str := range []string{"d", "g", "s", "x", "s", "s"} {
			C.NumFmts[name[i]] = make([]string, len(sizes[i]))
			for j, sz := range sizes[i] {
				C.NumFmts[name[i]][j] = io.Sf("%%%d%s", sz+1, str)
			}
		}
	}

	// compute sizes of header items
	nova := len(o[0].Ovas)
	noor := len(o[0].Oors)
	szova, szoor, szdem := make([]int, nova), make([]int, noor), 0
	for k, ind := range o {
		if C.ShowNinds > 0 && k >= C.ShowNinds {
			break
		}
		for i := 0; i < nova; i++ {
			szova[i] = utl.Imax(szova[i], len(io.Sf("%g", ind.Ovas[i])))
		}
		if C.ShowOor {
			for i := 0; i < noor; i++ {
				szoor[i] = utl.Imax(szoor[i], len(io.Sf("%g", ind.Oors[i])))
			}
		}
		szdem = utl.Imax(szdem, len(io.Sf("%g", ind.Demerit)))
	}
	for i := 0; i < nova; i++ {
		szova[i] = utl.Imax(szova[i], 5) // 5 ==> len("Ova##")
	}
	if C.ShowOor {
		for i := 0; i < noor; i++ {
			szoor[i] = utl.Imax(szoor[i], 5) // 5 ==> len("Oor####")
		}
	}
	szdem = utl.Imax(szdem, 7) // 7 ==> len("Demerit")

	// print individuals
	fmtova := make([]string, nova)
	fmtoor := make([]string, noor)
	for i := 0; i < nova; i++ {
		fmtova[i] = io.Sf("%%%d", szova[i]+1)
	}
	if C.ShowOor {
		for i := 0; i < noor; i++ {
			fmtoor[i] = io.Sf("%%%d", szoor[i]+1)
		}
	}
	fmtdem := io.Sf("%%%d", szdem+1)
	line, sza, szb := "", 0, 0
	first := true
	for i, ind := range o {
		if C.ShowNinds > 0 && i >= C.ShowNinds {
			break
		}
		stra := ""
		for j := 0; j < nova; j++ {
			stra += io.Sf(fmtova[j]+"g", ind.Ovas[j])
		}
		if C.ShowOor {
			for j := 0; j < noor; j++ {
				if ind.Oors[j] > 0 {
					stra += io.Sf(fmtoor[j]+"g", ind.Oors[j])
				} else {
					stra += io.Sf(fmtoor[j]+"s", "n/a")
				}
			}
		} else {
			unfeasible := false
			for j := 0; j < noor; j++ {
				if ind.Oors[j] > 0 {
					unfeasible = true
				}
			}
			if unfeasible {
				stra += " unfe."
			} else {
				stra += "      "
			}
		}
		if C.ShowDem {
			stra += io.Sf(fmtdem+"g", ind.Demerit) + " "
		}
		strb := ind.Output(C.NumFmts, C.ShowBases)
		line += stra + strb + "\n"
		if first {
			sza, szb = len(stra), len(strb)
			first = false
		}
	}

	// write to buffer
	fmtgenes := io.Sf("%%%d.%ds\n", szb, szb)
	n := sza + szb
	buf = new(bytes.Buffer)
	io.Ff(buf, io.StrThickLine(n))
	for i := 0; i < nova; i++ {
		io.Ff(buf, fmtova[i]+"s", io.Sf("Ova%d", i))
	}
	if C.ShowOor {
		for i := 0; i < noor; i++ {
			io.Ff(buf, fmtoor[i]+"s", io.Sf("Oor%d", i))
		}
	} else {
		io.Ff(buf, " check")
	}
	if C.ShowDem {
		io.Ff(buf, fmtdem+"s ", "Demerit")
	}
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

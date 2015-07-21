// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"sort"

	"github.com/cpmech/gosl/io"
)

// Population holds all individuals
type Population []*Individual

// NewPopFloatChromo allocates a population made entirely of float point numbers
//  Input:
//   genes -- all genes of all individuals [ninds][ngenes]
func NewPopFloatChromo(nbases int, genes [][]float64) (pop Population) {
	ninds := len(genes)
	pop = make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = new(Individual)
		pop[i].InitChromo(nbases, genes[i])
	}
	return
}

// NewPopReference creates a population based on a reference individual
func NewPopReference(ninds int, ref *Individual) (pop Population) {
	pop = make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = ref.GetCopy()
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
// to sort the population in decreasing fitness order: best => worst
func (o Population) Less(i, j int) bool {
	return o[i].Fitness > o[j].Fitness
}

// Sort sorts the population from best to worst individuals; i.e. decreasing fitness values
func (o *Population) Sort() {
	sort.Sort(o)
}

// Output generates a nice table with population data
//  Input:
//  fmts -- formats for int, flt, string, byte, bytes, and func
//          use fmts == nil to choose default ones
func (o Population) Output(fmts []string) (buf *bytes.Buffer) {

	// check
	if len(o) < 1 {
		return
	}
	if len(o[0].Chromo) < 1 {
		return
	}

	// compute sizes and generate formats list
	nfields := o[0].Chromo[0].Nfields()
	ngenes := len(o[0].Chromo)
	sizes := make([]int, 6)
	if fmts == nil {
		fmts = make([]string, 6)
		for _, ind := range o {
			sz := ind.GetStringSizes()
			for i := 0; i < 6; i++ {
				sizes[i] = imax(sizes[i], sz[i])
				if nfields == 1 {
					if sz[i] > 0 {
						if sizes[i]*ngenes < 5 { // 5 ==> len("Genes")
							sizes[i] = 3
							if ngenes == 1 {
								sizes[i] = 5
							}
						}
					}
				}
			}
		}
		for i, str := range []string{"d", "g", "s", "x", "s", "s"} {
			fmts[i] = io.Sf("%%%d%s", sizes[i]+1, str)
		}
	}

	// compute sizes of header items
	nOvl, nFit := 0, 0
	for _, ind := range o {
		nOvl = imax(nOvl, len(io.Sf("%g", ind.ObjValue)))
		nFit = imax(nFit, len(io.Sf("%g", ind.Fitness)))
	}
	nOvl = imax(nOvl, 6) // 6 ==> len("ObjVal")
	nFit = imax(nFit, 7) // 7 ==> len("Fitness")

	// print individuals
	fmtOvl := io.Sf("%%%d", nOvl+1)
	fmtFit := io.Sf("%%%d", nFit+1)
	line, sza, szb := "", 0, 0
	for i, ind := range o {
		stra := io.Sf(fmtOvl+"g", ind.ObjValue) + io.Sf(fmtFit+"g", ind.Fitness) + " "
		strb := ind.Output(fmts)
		line += stra + strb + "\n"
		if i == 0 {
			sza, szb = len(stra), len(strb)
		}
	}

	// write to buffer
	fmtGen := io.Sf(" %%%ds\n", szb)
	n := sza + szb
	buf = new(bytes.Buffer)
	io.Ff(buf, printThickLine(n))
	io.Ff(buf, fmtOvl+"s", "ObjVal")
	io.Ff(buf, fmtFit+"s", "Fitness")
	io.Ff(buf, fmtGen, "Genes")
	io.Ff(buf, printThinLine(n))
	io.Ff(buf, line)
	io.Ff(buf, printThickLine(n))
	return
}

// OutFloatBases print bases of float genes
func (o Population) OutFloatBases(numFmt string) (l string) {
	for _, ind := range o {
		for _, g := range ind.Chromo {
			l += io.Sf(numFmt, g.Fbases)
		}
		l += "\n"
	}
	return
}

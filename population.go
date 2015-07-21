// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
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
//  fmts -- formats for     int,     flt, string, byte,  bytes, and func
//          e.g: []string{"%4d", "%8.3f", "%.6s", "%x", "%.6s", "%.6s"}
//          use fmts == nil to choose default ones
func (o Population) Output(fmts []string) (l string) {

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
	nChr := nfields * ngenes // spaces between fields
	for i := 0; i < 6; i++ {
		nChr += sizes[i] * ngenes
	}
	if nfields > 1 {
		nChr += ngenes*2 + (ngenes - 1) // 2 ==> "(" and ")" and (ngens-1) ==> space between
	}
	nChr = imax(5, nChr) // 5 ==> len("Genes")

	// print individuals
	nOvl = imax(nOvl, 6)        // 6 ==> len("ObjVal")
	nFit = imax(nFit, 7)        // 7 ==> len("Fitness")
	n := nOvl + nFit + nChr + 3 // 3 ==> spaces beeen "ObjVal", "Fitness" and "Genes"
	fmtOvl := io.Sf("%%%d", nOvl+1)
	fmtFit := io.Sf("%%%d", nFit+1)
	fmtChr := io.Sf("%%%ds", nChr)
	l += printThickLine(n)
	l += io.Sf(fmtOvl+"s", "ObjVal")
	l += io.Sf(fmtFit+"s", "Fitness") + " "
	l += io.Sf(fmtChr, "Genes")
	l += "\n" + printThinLine(n)
	fmtOvl += "g"
	fmtFit += "g"
	for _, ind := range o {
		l += io.Sf(fmtOvl, ind.ObjValue) + io.Sf(fmtFit, ind.Fitness) + " " + ind.Output(fmts) + "\n"
	}
	l += printThickLine(n)
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

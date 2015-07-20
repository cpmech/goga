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
func (o Population) Output() (l string) {

	// check
	if len(o) < 1 {
		return
	}
	if len(o[0].Chromo) < 1 {
		return
	}

	// mixed genes type
	if o[0].Chromo[0].Nfields() > 1 {
		return "TODO: mixed genes type"
	}

	// single type in genes
	sizes := make([]int, 4) // int, float, string, byte
	nOvl, nFit, nGen := 0, 0, 0
	for _, ind := range o {
		nOvl = imax(nOvl, len(io.Sf("%g", ind.ObjValue)))
		nFit = imax(nFit, len(io.Sf("%g", ind.Fitness)))
		sz := ind.GetStringSizes()
		for i := 0; i < 4; i++ {
			sizes[i] = imax(sizes[i], sz[i])
			if sz[i] > 0 {
				if sizes[i]*len(ind.Chromo) < 5 { // 5 ==> "Genes"
					sizes[i] = 3
					if len(ind.Chromo) == 1 {
						sizes[i] = 5
					}
				}
				nGen = (sizes[i] + 1) * len(ind.Chromo)
			}
		}
	}
	nOvl = imax(nOvl, 6) // 6 ==> "ObjVal"
	nFit = imax(nFit, 7) // 7 ==> "Fitness"
	fmts := make([]string, 4)
	n := nOvl + nFit + 2 + nGen
	for i, str := range []string{"d", "g", "s", "s"} {
		fmts[i] = io.Sf("%%%d%s", sizes[i]+1, str)
	}
	fmtOvl := io.Sf("%%%d", nOvl+1)
	fmtFit := io.Sf("%%%d", nFit+1)
	fmtGen := io.Sf("%%%ds", nGen)
	l += printThickLine(n)
	l += io.Sf(fmtOvl+"s", "ObjVal")
	l += io.Sf(fmtFit+"s", "Fitness")
	l += io.Sf(fmtGen, "Genes")
	l += "\n" + printThinLine(n)
	fmtOvl += "g"
	fmtFit += "g"
	for _, ind := range o {
		l += io.Sf(fmtOvl, ind.ObjValue) + io.Sf(fmtFit, ind.Fitness) + ind.Output(fmts) + "\n"
	}
	l += printThickLine(n)
	return
}

// OutFloatBases print bases of float genes
func (o Population) OutFloatBases(numFmt string) (l string) {
	for _, ind := range o {
		for _, g := range ind.Chromo {
			l += io.Sf(numFmt, g.SubFloats)
		}
		l += "\n"
	}
	return
}

// allocators //////////////////////////////////////////////////////////////////////////////////////

// NewFloatChromoPop allocates a population made entirely of float point numbers
//  Input:
//   genes -- all genes of all individuals [ninds][ngenes]
func NewFloatChromoPop(nbases int, genes [][]float64) (pop Population) {
	ninds := len(genes)
	pop = make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = new(Individual)
		pop[i].InitChromo(nbases, genes[i])
	}
	return
}

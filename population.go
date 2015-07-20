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
//  Input:
//   prob    -- probabilities
//   cumprob -- cumulated probabilities
func (o Population) Output(prob, cumprob []float64) (l string) {
	nI, nF, nS, nB := 0, 0, 0, 0
	for _, ind := range o {
		sI, sF, sS, sB := ind.GetStringSizes()
		nI = imax(nI, sI)
		nF = imax(nF, sF)
		nS = imax(nS, sS)
		nB = imax(nB, sB)
	}
	nI, nF, nS, nB = nI+1, nF+1, nS+1, nB+1
	fI, fF, fS, fB := io.Sf("%%%dd", nI), io.Sf("%%%dg", nF), io.Sf("%%%ds", nS), io.Sf("%%%ds", nB)
	for _, ind := range o {
		l += ind.Output(fI, fF, fS, fB) + "\n"
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

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// Individual implements one individual in a population
type Individual struct {

	// essential
	Ngenes int       // number of genes
	Nbases int       // number of bases
	Chromo []float64 // chromosome [ngenes*nbases]

	// derived
	ObjValue float64 // objective value
	Fitness  float64 // fitness

	// scratchpad
	Genes []float64 // genes [ngenes]
}

// CalcGenes calculate gene values
//  Output:
//   o.Genes -- will hold the gene values
func (o *Individual) CalcGenes() {
	chk.IntAssertLessThan(0, o.Ngenes)
	chk.IntAssertLessThan(0, o.Nbases)
	chk.IntAssert(o.Ngenes*o.Nbases, len(o.Chromo))
	if len(o.Genes) != o.Ngenes {
		o.Genes = make([]float64, o.Ngenes)
	}
	if o.Nbases == 1 {
		copy(o.Genes, o.Chromo)
		return
	}
	for i := 0; i < o.Ngenes; i++ {
		o.Genes[i] = 0
		for j := 0; j < o.Nbases; j++ {
			o.Genes[i] += o.Chromo[i*o.Nbases+j]
		}
	}
}

// String returns a table-row representation of an individual
//  Input:
//   ovfmt -- objective value formatting string; use "" to skip this item
//   ftfmt -- fitness formatting string; use "" to skip this item
//   gefmt -- genes formatting string; use "" to skip this item
//   bsfmt -- bases formatting string; use "" to skip this item
func (o Individual) String(ovfmt, ftfmt, gefmt, bsfmt string) (line string) {
	if ovfmt != "" {
		line += io.Sf(ovfmt, o.ObjValue)
	}
	if ftfmt != "" {
		line += io.Sf(ftfmt, o.Fitness)
	}
	if gefmt != "" {
		for _, g := range o.Genes {
			line += io.Sf(gefmt, g)
		}
	}
	if bsfmt != "" {
		for _, v := range o.Chromo {
			line += io.Sf(bsfmt, v)
		}
	}
	return
}

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

// GenTable generates a nice table with population data
//  Input:
//   prob      -- probabilities
//   cumprob   -- cumulated probabilities
//   showbases -- show basis values
func (o Population) GenTable(prob, cumprob []float64, showbases bool) (line string) {

	// number of bases
	if len(o) < 1 {
		return
	}
	ngenes := o[0].Ngenes
	nbases := o[0].Nbases

	// find max sizes of strings
	var ovsz, ftsz, gesz, bssz int
	for _, ind := range o {
		ind.CalcGenes()
		ovsz = imax(ovsz, len(io.Sf("%g", ind.ObjValue)))
		ftsz = imax(ftsz, len(io.Sf("%g", ind.Fitness)))
		for _, g := range ind.Genes {
			gesz = imax(gesz, len(io.Sf("%g", g)))
		}
		if showbases {
			for _, v := range ind.Chromo {
				bssz = imax(bssz, len(io.Sf("%g", v)))
			}
		}
	}

	// lengths of fields
	ovsz, ftsz, gesz, bssz = ovsz+1, ftsz+1, gesz+1, bssz+1
	ovsz = imax(7, ovsz)
	ftsz = imax(8, ftsz)
	allge := gesz * ngenes
	allbs := bssz * ngenes * nbases
	if allge < 6 {
		gesz = 6
	}
	if allbs < 12 {
		bssz = 12
	}
	if !showbases {
		bssz, allbs = 0, 0
	}
	total := ovsz + ftsz + allge + allbs

	// define formats
	ovnum := io.Sf("%%%dg", ovsz)
	ftnum := io.Sf("%%%dg", ftsz)
	genum := io.Sf("%%%dg", gesz)
	bsnum := io.Sf("%%%dg", bssz)
	ovstr := io.Sf("%%%ds", ovsz)
	ftstr := io.Sf("%%%ds", ftsz)
	gestr := io.Sf("%%%ds", imax(6, allge))
	bsstr := io.Sf("%%%ds", imax(12, allbs))
	line += printThickLine(total)
	if showbases {
		line += io.Sf(ovstr+ftstr+gestr+bsstr+"\n", "ObjVal", "Fitness", "Genes", "Chromosomes")
	} else {
		line += io.Sf(ovstr+ftstr+gestr+"\n", "ObjVal", "Fitness", "Genes")
		bsnum = ""
	}
	line += printThinLine(total)

	// print individuals
	for _, ind := range o {
		line += ind.String(ovnum, ftnum, genum, bsnum) + "\n"
	}
	line += printThickLine(total)
	return
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// Evolver realises the evolutionary process
type Evolver struct {
	C       *ConfParams // configuration parameters
	Islands []*Island   // islands
	Best    *Individual // best individual among all in all islands
}

// NewEvolver creates a new evolver
//  Input:
//   nislands -- number of islands
//   ninds    -- number of individuals to be generated
//   ref      -- reference individual with chromosome structure already set
//   bingo    -- Bingo structure set with pool of values to draw gene values
//   ovfunc   -- objective function
func NewEvolver(C *ConfParams, ref *Individual, ovfunc ObjFunc_t, bingo *Bingo) (o *Evolver) {
	o = new(Evolver)
	o.C = C
	o.Islands = make([]*Island, o.C.Nisl)
	for i := 0; i < o.C.Nisl; i++ {
		o.Islands[i] = NewIsland(i, o.C, NewPopRandom(o.C.Ninds, ref, bingo), ovfunc, bingo)
	}
	return
}

// NewEvolverPop creates a new evolver based on a given population
//  Input:
//   pops   -- populations. len(pop) == nislands
//   ovfunc -- objective function
func NewEvolverPop(C *ConfParams, pops []Population, ovfunc ObjFunc_t, bingo *Bingo) (o *Evolver) {
	o = new(Evolver)
	o.C = C
	chk.IntAssert(C.Nisl, len(pops))
	o.Islands = make([]*Island, o.C.Nisl)
	for i, pop := range pops {
		o.Islands[i] = NewIsland(i, o.C, pop, ovfunc, bingo)
	}
	return
}

// NewEvolverFloatChromo creates a new evolver with float point individuals
//  Input:
//   C.Ninds  -- number of individuals to be generated
//   C.Nbases -- number of bases
//   C.Grid   -- whether or not to calc values based on grid;
//               otherwise select randomly between xmin and xmax
//   C.Noise  -- if noise>0, apply noise to move points away from grid nodes
//               noise is a multiplier; e.g. 0.2
//   xmin     -- min values of genes
//   xmax     -- max values of genes. len(xmin) = len(xmax) = ngenes
func NewEvolverFloatChromo(C *ConfParams, xmin, xmax []float64, ovfunc ObjFunc_t, bingo *Bingo) (o *Evolver) {
	pops := make([]Population, C.Nisl)
	for i := 0; i < C.Nisl; i++ {
		pops[i] = NewPopFloatRandom(C, xmin, xmax)
	}
	return NewEvolverPop(C, pops, ovfunc, bingo)
}

// Run runs the evolution process
func (o *Evolver) Run(verbose, doreport bool) {

	// check
	nislands := len(o.Islands)
	if nislands < 1 {
		return
	}
	if o.C.Ninds < nislands {
		chk.Panic("number of individuals must be greater than the number of islands")
	}

	// first output
	t := 0
	for _, isl := range o.Islands {
		isl.WritePopToReport(t)
	}
	if verbose {
		io.Pf("\nrunning ...\n")
	}

	// for migration
	iworst := o.C.Ninds - 1
	receiveFrom := utl.IntsAlloc(nislands, nislands-1)

	// time loop
	t = 1
	tmig := o.C.Dtmig
	for t < o.C.Tf {

		// evolve up to migration time
		if o.C.Pll {
		} else {
			for i, isl := range o.Islands {
				for time := t; time < tmig; time++ {
					regen := o.calc_regen(time)
					report := o.calc_report(time)
					isl.SelectReprodAndRegen(time, regen, report)
					if verbose && i == 0 {
						o.print_time(time, regen, report)
					}
				}
			}
		}

		// update time
		t = tmig
		tmig += o.C.Dtmig
		if tmig > o.C.Tf {
			tmig = o.C.Tf
		}

		// reset receiveFrom matrix
		for i := 0; i < nislands; i++ {
			for j := 0; j < nislands-1; j++ {
				receiveFrom[i][j] = -1
			}
		}

		// compute destinations
		for i := 0; i < nislands; i++ {
			Aworst := o.Islands[i].Pop[iworst]
			k := 0
			for j := 0; j < nislands; j++ {
				if i != j {
					Bbest := o.Islands[j].Pop[0]
					send := Bbest.Compare(Aworst)
					if send {
						receiveFrom[i][k] = j // i gets individual from j
						k++
					}
				}
			}
		}

		// migration
		if verbose {
			io.Pfyel("\n%d : migration\n", t)
		}
		for i, from := range receiveFrom {
			k := 0
			for _, j := range from {
				if j >= 0 {
					o.Islands[j].Pop[0].CopyInto(o.Islands[i].Pop[iworst-k])
					k++
				}
			}
			o.Islands[i].CalcDemeritsAndSort(o.Islands[i].Pop)
		}
	}

	// best individual
	o.FindBestFromAll()

	// message
	if verbose {
		io.Pf("... end\n\n")
	}

	// write reports
	if doreport {
		for _, isl := range o.Islands {
			isl.SaveReport(verbose)
		}
	}

	// plot evolution
	if o.C.DoPlot {
		for i, isl := range o.Islands {
			isl.PlotOvs(".eps", io.Sf("label='island %d'", i), o.C.PltTi, o.C.PltTf, false, "%.6f", i == 0, i == o.C.Nisl-1)
		}
		for i, isl := range o.Islands {
			isl.PlotOor(".eps", io.Sf("label='island %d'", i), o.C.PltTi, o.C.PltTf, false, "%.6f", i == 0, i == o.C.Nisl-1)
		}
	}
	return
}

// FindBestFromAll finds best individual from all islands
//  Output: o.Best will point to the best individual
func (o *Evolver) FindBestFromAll() {
	if len(o.Islands) < 1 {
		return
	}
	o.Best = o.Islands[0].Pop[0] // TODO: check case of oor individuals
	for _, isl := range o.Islands {
		if isl.Pop[0].Ova < o.Best.Ova {
			o.Best = isl.Pop[0]
		}
	}
}

func (o Evolver) calc_regen(t int) bool {
	if t == o.C.RegIni {
		return true
	}
	return t%o.C.Dtreg == 0
}

func (o Evolver) calc_report(t int) bool {
	return t%o.C.Dtout == 0
}

func (o Evolver) print_time(time int, regen, report bool) {
	io.Pf(" ")
	if regen {
		io.Pf("%v", time)
		return
	}
	if report {
		io.Pfblue("%v", time)
		return
	}
	io.Pfgrey("%v", time)
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/io"

// Evolver realises the evolutionary process
type Evolver struct {
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
func NewEvolver(nislands, ninds int, ref *Individual, bingo *Bingo, ovfunc ObjFunc_t) (o *Evolver) {
	o = new(Evolver)
	o.Islands = make([]*Island, nislands)
	for i := 0; i < nislands; i++ {
		o.Islands[i] = NewIsland(NewPopRandom(ninds, ref, bingo), ovfunc)
	}
	return
}

// NewEvolverPop creates a new evolver based on a given population
//  Input:
//   pops   -- populations. len(pop) == nislands
//   ovfunc -- objective function
func NewEvolverPop(pops []Population, ovfunc ObjFunc_t) (o *Evolver) {
	o = new(Evolver)
	nislands := len(pops)
	o.Islands = make([]*Island, nislands)
	for i, pop := range pops {
		o.Islands[i] = NewIsland(pop, ovfunc)
	}
	return
}

// Run runs the evolution process
//  Input:
//   tf    -- final time
//   dtout -- increment of time for output
//   dtmig -- increment of time for migration
func (o *Evolver) Run(tf, dtout, dtmig int) {

	// check
	nislands := len(o.Islands)
	if nislands < 1 {
		return
	}

	/*
		o.Dest = make([][]bool, nislands)
		for i := 0; i < nislands; i++ {
			o.Dest[i] = make([]bool, nislands)
		}
		o.Ids = utl.IntRange(nislands)
	*/

	// header
	lent := len(io.Sf("%d", tf))
	strt := io.Sf("%%%d", lent+2)
	io.Pf("%s", printThickLine(lent+2+11+25))
	io.Pf(strt+"s%11s%25s\n", "time", "migration", "objval")
	io.Pf("%s", printThinLine(lent+2+11+25))
	strt = strt + "d%11s%25g\n"

	// time control
	t := 0
	tout := dtout
	tmig := dtmig

	// best individual
	o.Best = o.Islands[0].Pop[0]

	// first output
	io.Pf(strt, t, "", o.Best.ObjValue)

	// time loop
	done := make(chan int, nislands)
	for t < tf {

		// reproduction in all islands
		for i := 0; i < nislands; i++ {
			go func(isl *Island) {
				for j := t; j < tout; j++ {
					isl.SelectAndReprod(j)
				}
				done <- 1
			}(o.Islands[i])
		}

		// listen to channels
		for i := 0; i < nislands; i++ {
			<-done
		}

		// current time and next cycle
		t += dtout
		tout = t + dtout

		// migration
		mig := ""
		if t >= tmig {
			mig = "true"
			tmig = t + dtmig
			for i := 0; i < nislands; i++ {
				for j := i + 1; j < nislands; j++ {
					last := len(o.Islands[j].Pop) - 1
					o.Islands[i].Pop[0].CopyInto(o.Islands[j].Pop[last]) // iBest => jWorst
					o.Islands[j].Pop[0].CopyInto(o.Islands[i].Pop[last]) // jBest => iWorst
					o.Islands[i].Pop.Sort()
					o.Islands[j].Pop.Sort()
				}
			}
		}

		// best individual from all islands
		o.Best = o.Islands[0].Pop[0]
		for i := 0; i < nislands; i++ {
			if o.Islands[i].Pop[0].ObjValue < o.Best.ObjValue {
				o.Best = o.Islands[i].Pop[0]
			}
		}

		// output
		io.Pf(strt, t, mig, o.Best.ObjValue)
	}

	// footer
	io.Pf("%s", printThickLine(lent+2+11+25))
}

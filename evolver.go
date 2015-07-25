// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/io"

// Evolver realises the evolutionary process
type Evolver struct {
	Islands []*Island
	BestOV  float64
	// TODO: store pointer to best individual as well
}

// NewEvolver creates a new evolver
//  Input:
//   nislands --  number of islands
//   ninds    -- number of individuals to be generated
//   ref      -- reference individual with chromosome structure already set
//   bingo    -- Bingo structure set with pool of values to draw gene values
//   ovfunc   -- objective function
//  Output:
//   new population
func NewEvolver(nislands, ninds int, ref *Individual, bingo *Bingo, ovfunc ObjFunc_t) (o *Evolver) {
	o = new(Evolver)
	o.Islands = make([]*Island, nislands)
	for i := 0; i < nislands; i++ {
		o.Islands[i] = NewIsland(NewPopRandom(ninds, ref, bingo), ovfunc)
	}
	return
}

func (o *Evolver) Run(tf, dtout, dtmig int) {

	// check
	nislands := len(o.Islands)
	if nislands < 1 {
		return
	}

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
	o.BestOV = o.Islands[0].Pop[0].ObjValue

	// first output
	io.Pf(strt, t, "", o.BestOV)

	// time loop
	var ov float64
	ovs := make(chan float64, nislands)
	for t < tf {

		// reproduction in all islands
		for i := 0; i < nislands; i++ {
			go func(isl *Island) {
				for j := t; j < tout; j++ {
					isl.SelectAndReprod(j)
				}
				ovs <- isl.Pop[0].ObjValue
			}(o.Islands[i])
		}

		// listen to channels and get best OV
		o.BestOV = <-ovs
		for i := 1; i < nislands; i++ {
			ov = <-ovs
			if ov < o.BestOV {
				o.BestOV = ov
			}
		}

		// current time and next cycle
		t += dtout
		tout = t + dtout

		// migration
		mig := ""
		if t >= tmig {
			// TODO: implement this
			mig = "true"
			tmig = t + dtmig
		}

		// output
		io.Pf(strt, t, mig, o.BestOV)
	}

	// footer
	io.Pf("%s", printThickLine(lent+2+11+25))
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/io"

// Evolver realises the evolutionary process
type Evolver struct {
	Islands []*Island
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
	bestov := o.Islands[0].Pop[0].ObjValue

	// first output
	io.Pf(strt, t, "", bestov)

	// time loop
	done := make(chan int, nislands)
	for t < tf {

		// reproduction in all islands
		for i := 0; i < nislands; i++ {
			// TODO: activate this
			//go func() {
			for j := t; j < tout; j++ {
				o.Islands[i].SelectAndReprod(j)
			}
			done <- 1
			//}()
		}
		for i := 0; i < nislands; i++ {
			<-done
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

		// best individual
		bestov = o.Islands[0].Pop[0].ObjValue
		// TODO: get best among all islands

		// output
		io.Pf(strt, t, mig, bestov)
	}

	// footer
	io.Pf("%s", printThickLine(lent+2+11+25))
}

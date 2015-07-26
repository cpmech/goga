// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"os"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// Evolver realises the evolutionary process
type Evolver struct {
	Islands []*Island   // islands
	Best    *Individual // best individual among all in all islands
	DirOut  string      // directory to save output files. "" means "/tmp/goga/"
	FnKey   string      // filename key for output files. "" means no output files
	Json    bool        // output results as .json files; not tables
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
//   tf      -- final time
//   dtout   -- increment of time for output
//   dtmig   -- increment of time for migration
//   verbose -- print information suring progress
func (o *Evolver) Run(tf, dtout, dtmig int, verbose bool) {

	// check
	nislands := len(o.Islands)
	if nislands < 1 {
		return
	}

	// time control
	if dtout < 1 {
		dtout = 1
	}
	t := 0
	tout := dtout
	tmig := dtmig

	// best individual and index of worst individual
	o.FindBestFromAll()
	iworst := len(o.Islands[0].Pop) - 1

	// saving results
	dosave := o.prepare_for_saving_results()

	// header
	lent := len(io.Sf("%d", tf))
	strt := io.Sf("%%%d", lent+2)
	if verbose {
		io.Pf("%s", printThickLine(lent+2+11+25))
		io.Pf(strt+"s%11s%25s\n", "time", "migration", "objval")
		io.Pf("%s", printThinLine(lent+2+11+25))
		strt = strt + "d%11s%25g\n"
		io.Pf(strt, t, "", o.Best.ObjValue)
	}

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
			for i := 0; i < nislands; i++ {
				for j := i + 1; j < nislands; j++ {
					o.Islands[i].Pop[0].CopyInto(o.Islands[j].Pop[iworst]) // iBest => jWorst
					o.Islands[j].Pop[0].CopyInto(o.Islands[i].Pop[iworst]) // jBest => iWorst
				}
			}
			for _, isl := range o.Islands {
				isl.Pop.Sort()
			}
			mig = "true"
			tmig = t + dtmig
		}

		// output
		o.FindBestFromAll()
		if verbose {
			io.Pf(strt, t, mig, o.Best.ObjValue)
		}
	}

	// save results
	if dosave {
		o.save_results("final", t)
	}

	// footer
	if verbose {
		io.Pf("%s", printThickLine(lent+2+11+25))
	}
	return
}

// FindBestFromAll finds best individual from all islands
//  Output: o.Best will point to the best individual
func (o *Evolver) FindBestFromAll() {
	if len(o.Islands) < 1 {
		return
	}
	o.Best = o.Islands[0].Pop[0]
	for _, isl := range o.Islands {
		if isl.Pop[0].ObjValue < o.Best.ObjValue {
			o.Best = isl.Pop[0]
		}
	}
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func (o *Evolver) prepare_for_saving_results() (dosave bool) {
	dosave = o.FnKey != ""
	if dosave {
		if o.DirOut == "" {
			o.DirOut = "/tmp/goga/"
		}
		err := os.MkdirAll(o.DirOut, 0777)
		if err != nil {
			chk.Panic("cannot create directory:%v", err)
		}
		io.RemoveAll(io.Sf("%s/%s*", o.DirOut, o.FnKey))
		o.save_results("initial", 0)
	}
	return
}

func (o Evolver) save_results(key string, t int) {
	var b bytes.Buffer
	for i, isl := range o.Islands {
		if i > 0 {
			if o.Json {
				io.Ff(&b, ",\n")
			} else {
				io.Ff(&b, "\n")
			}
		}
		isl.Write(&b, t, o.Json)
	}
	ext := "res"
	if o.Json {
		ext = "json"
	}
	io.WriteFile(io.Sf("%s/%s_%s.%s", o.DirOut, o.FnKey, key, ext), &b)
}

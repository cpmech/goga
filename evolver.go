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

// Run runs the evolution process
func (o *Evolver) Run(verbose bool) {

	// check
	nislands := len(o.Islands)
	if nislands < 1 {
		return
	}

	// time control
	t := 0
	tout := o.C.Dtout
	tmig := o.C.Dtmig
	treg := o.C.Dtreg

	// regeneration control
	idxreg := 0
	if o.C.RegNmax < 0 {
		o.C.RegNmax = o.C.Tf + 1
	}
	if o.Islands[0].Pop[0].Nfltgenes == 0 {
		o.C.RegTol = 0
	}

	// best individual and index of worst individual
	o.FindBestFromAll()
	iworst := o.C.Ninds - 1
	_, averho, _, _ := o.Islands[0].Stat()

	// saving results
	dosave := o.prepare_for_saving_results(verbose)

	// header
	lent := len(io.Sf("%d", o.C.Tf))
	if lent < 5 {
		lent = 5
	}
	strt := io.Sf("%%%d", lent)
	szline := lent + 6 + 6 + 11 + 25 + 25
	if verbose {
		io.Pf("%s", printThickLine(szline))
		io.Pf(strt+"s%6s%6s%11s%25s%25s\n", "time", "mig", "reg", "ave(rho)", "ova", "oor")
		io.Pf("%s", printThinLine(szline))
		strt = strt + "d%6s%6s%11.3e%25g%25g\n"
		io.Pf(strt, t, "", "", averho, o.Best.Ova, o.Best.Oor)
	}
	strreg := []string{"", "best", "lims"}

	// time loop
	var res Comm_t
	var regidx int
	ch := make(chan Comm_t, nislands)
	for t := 1; t < o.C.Tf+1; t++ {

		// perform regeneration?
		doregen := false
		if (t == 1 && o.C.RegIni) || (t >= treg && idxreg < o.C.RegNmax) {
			doregen = true
			treg = t + o.C.Dtreg
			idxreg += 1
		}

		// selection, reproduction, regeneration and reporting
		if o.C.Pll {
			for i := 0; i < nislands; i++ {
				go func(isl *Island) {
					comm := isl.SelectReprodAndRegen(t, tout, doregen)
					ch <- comm
				}(o.Islands[i])
			}
			res = <-ch
			averho = res.AveRho
			regidx = res.RegIdx
			for i := 1; i < nislands; i++ {
				res = <-ch
				averho = min(averho, res.AveRho)
				regidx = imax(regidx, res.RegIdx)
			}
		} else {
			for i, isl := range o.Islands {
				comm := isl.SelectReprodAndRegen(t, tout, doregen)
				if i == 0 {
					averho = comm.AveRho
					regidx = comm.RegIdx
				} else {
					averho = min(averho, comm.AveRho)
					regidx = imin(regidx, comm.RegIdx)
				}
			}
		}

		// migration
		mig := ""
		if t >= tmig && nislands > 1 {
			for i := 0; i < nislands; i++ {
				for j := i + 1; j < nislands; j++ {
					o.Islands[i].Pop[0].CopyInto(o.Islands[j].Pop[iworst]) // iBest => jWorst
					o.Islands[j].Pop[0].CopyInto(o.Islands[i].Pop[iworst]) // jBest => iWorst
				}
			}
			for _, isl := range o.Islands {
				isl.CalcDemeritsAndSort(isl.Pop)
			}
			mig = "true"
			tmig = t + o.C.Dtmig
		}

		// best individual from all islands
		o.FindBestFromAll()

		// output
		if verbose && t >= tout {
			io.Pf(strt, t, mig, strreg[regidx], averho, o.Best.Ova, o.Best.Oor)
			tout += o.C.Dtout
		}
	}

	// footer
	if verbose {
		io.Pf("%s", printThickLine(szline))
	}

	// save results
	if dosave {
		o.save_results("final", t, verbose)
		for i, isl := range o.Islands {
			isl.SaveReport(o.C.DirOut, io.Sf("%s-isl%d", o.C.FnKey, i), verbose)
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
	o.Best = o.Islands[0].Pop[0]
	for _, isl := range o.Islands {
		if isl.Pop[0].Ova < o.Best.Ova {
			o.Best = isl.Pop[0]
		}
	}
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func (o *Evolver) prepare_for_saving_results(verbose bool) (dosave bool) {
	dosave = o.C.FnKey != ""
	if dosave {
		if o.C.DirOut == "" {
			o.C.DirOut = "/tmp/goga"
		}
		err := os.MkdirAll(o.C.DirOut, 0777)
		if err != nil {
			chk.Panic("cannot create directory:%v", err)
		}
		io.RemoveAll(io.Sf("%s/%s*", o.C.DirOut, o.C.FnKey))
		o.save_results("initial", 0, verbose)
	}
	return
}

func (o Evolver) save_results(key string, t int, verbose bool) {
	var b bytes.Buffer
	for i, isl := range o.Islands {
		if i > 0 {
			if o.C.Json {
				io.Ff(&b, ",\n")
			} else {
				io.Ff(&b, "\n")
			}
		}
		isl.Write(&b, t, o.C.Json)
	}
	ext := "res"
	if o.C.Json {
		ext = "json"
	}
	write := io.WriteFile
	if t > 0 && verbose {
		write = io.WriteFileV
		io.Pf("\n")
	}
	write(io.Sf("%s/%s-%s.%s", o.C.DirOut, o.C.FnKey, key, ext), &b)
	if t > 0 {
		for i, isl := range o.Islands {
			if isl.Report.Len() > 0 {
				write(io.Sf("%s/%s-isl%d.rpt", o.C.DirOut, o.C.FnKey, i), &isl.Report)
			}
		}
	}
}

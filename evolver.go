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
		o.Islands[i] = NewIsland(i, NewPopRandom(ninds, ref, bingo), ovfunc)
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
		o.Islands[i] = NewIsland(i, pop, ovfunc)
	}
	return
}

// Run runs the evolution process
//  Input:
//   tf      -- final time
//   dtout   -- increment of time for output
//   dtmig   -- increment of time for migration
//   dtreg   -- increment of time for regeneration
//   nreg    -- number of regenerations allowed. -1 means unlimited
//   verbose -- print information suring progress
func (o *Evolver) Run(tf, dtout, dtmig, dtreg, nreg int, verbose bool) {

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
	treg := dtreg

	// regeneration control
	idxreg := 0
	if nreg < 0 {
		nreg = tf + 1
	}

	// best individual and index of worst individual
	o.FindBestFromAll()
	iworst := len(o.Islands[0].Pop) - 1
	minrho, averho, maxrho, devrho := o.calc_stat()

	// saving results
	dosave := o.prepare_for_saving_results(verbose)

	// header
	lent := len(io.Sf("%d", tf))
	if lent < 5 {
		lent = 5
	}
	strt := io.Sf("%%%d", lent)
	szline := lent + 6 + 6 + 11 + 11 + 11 + 11 + 25
	if verbose {
		io.Pf("%s", printThickLine(szline))
		io.Pf(strt+"s%6s%6s%11s%11s%11s%11s%25s\n", "time", "mig", "reg", "min(rho)", "ave(rho)", "max(rho)", "dev(rho)", "objval")
		io.Pf("%s", printThinLine(szline))
		strt = strt + "d%6s%6s%11.3e%11.3e%11.3e%11.3e%25g\n"
		io.Pf(strt, t, "", "", minrho, averho, maxrho, devrho, o.Best.ObjValue)
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

		// statistics
		minrho, averho, maxrho, devrho = o.calc_stat()

		// regeneration
		reg := ""
		if t >= treg && idxreg < nreg {
			for i := 0; i < nislands; i++ {
				go func(isl *Island) {
					isl.Regenerate(t)
					done <- 1
				}(o.Islands[i])
			}
			for i := 0; i < nislands; i++ {
				<-done
			}
			reg = "true"
			treg = t + dtreg
			idxreg += 1
		}

		// best individual
		o.FindBestFromAll()

		// output
		if verbose {
			io.Pf(strt, t, mig, reg, minrho, averho, maxrho, devrho, o.Best.ObjValue)
		}
	}

	// footer
	if verbose {
		io.Pf("%s", printThickLine(szline))
	}

	// save results
	if dosave {
		o.save_results("final", t, verbose)
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

// SetParams sets all islands with given paramters
func (o *Evolver) SetParams(params *Params) {
	o.FnKey = params.Fnkey
	pc, pm := params.Pc, params.Pm
	for _, isl := range o.Islands {
		isl.CxProbs = map[string]float64{"int": pc, "flt": pc, "str": pc, "key": pc, "byt": pc, "fun": pc}
		isl.MtProbs = map[string]float64{"int": pm, "flt": pm, "str": pm, "key": pm, "byt": pm, "fun": pm}
		isl.Elitism = params.Elite
		isl.UseRanking = params.Rnk
		isl.RnkPressure = params.RnkSP
		isl.Roulette = params.Rws
	}
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// calc_stat computes some statistical data from float bases
//  Note: avedev is actually the maximum average among all islands
func (o Evolver) calc_stat() (minrho, averho, maxrho, devrho float64) {
	nislands := len(o.Islands)
	type res_t struct{ xmin, xave, xmax, xdev float64 }
	results := make(chan res_t, nislands)
	for i := 0; i < nislands; i++ {
		go func(isl *Island) {
			xmin, xave, xmax, xdev := isl.Stat()
			results <- res_t{xmin, xave, xmax, xdev}
		}(o.Islands[i])
	}
	res := <-results
	minrho, averho, maxrho, devrho = res.xmin, res.xave, res.xmax, res.xdev
	for i := 1; i < nislands; i++ {
		res = <-results
		minrho = min(minrho, res.xmin)
		averho = min(averho, res.xave)
		maxrho = min(maxrho, res.xmax)
		devrho = min(devrho, res.xdev)
	}
	return
}

func (o *Evolver) prepare_for_saving_results(verbose bool) (dosave bool) {
	dosave = o.FnKey != ""
	if dosave {
		if o.DirOut == "" {
			o.DirOut = "/tmp/goga"
		}
		err := os.MkdirAll(o.DirOut, 0777)
		if err != nil {
			chk.Panic("cannot create directory:%v", err)
		}
		io.RemoveAll(io.Sf("%s/%s*", o.DirOut, o.FnKey))
		o.save_results("initial", 0, verbose)
	}
	return
}

func (o Evolver) save_results(key string, t int, verbose bool) {
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
	write := io.WriteFile
	if t > 0 && verbose {
		write = io.WriteFileV
		io.Pf("\n")
	}
	write(io.Sf("%s/%s-%s.%s", o.DirOut, o.FnKey, key, ext), &b)
	if t > 0 {
		for i, isl := range o.Islands {
			if isl.Report.Len() > 0 {
				write(io.Sf("%s/%s-isl%d.rpt", o.DirOut, o.FnKey, i), &isl.Report)
			}
		}
	}
}

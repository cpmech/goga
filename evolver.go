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

// NewEvolverPop creates a new evolver based on given populations
func NewEvolver(C *ConfParams) (o *Evolver) {
	o = new(Evolver)
	o.C = C
	o.Islands = make([]*Island, o.C.Nisl)
	for i := 0; i < o.C.Nisl; i++ {
		o.Islands[i] = NewIsland(i, o.C)
	}
	o.Best = o.Islands[0].Pop[0]
	return
}

// Run runs the evolution process
func (o *Evolver) Run() {

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
		isl.WritePopToReport(t, 0)
	}
	if o.C.Verbose {
		o.print_legend()
		io.Pf("\nrunning ...\n")
	}
	if o.C.PostProc != nil {
		for _, isl := range o.Islands {
			o.C.PostProc(isl.Id, 0, isl.Pop)
		}
	}

	// for migration
	iworst := o.C.Ninds - 1
	receiveFrom := utl.IntsAlloc(nislands, nislands-1)

	// time loop
	t = 1
	tmig := o.C.Dtmig
	nomig := false
	if tmig > o.C.Tf {
		tmig = o.C.Tf
		nomig = true
	}
	done := make(chan int, nislands)
	for t < o.C.Tf {

		// evolve up to migration time
		if o.C.Pll {
			for i := 0; i < nislands; i++ {
				go func(isl *Island) {
					for time := t; time < tmig; time++ {
						report := o.calc_report(time)
						isl.Run(time, report, (o.C.Verbose && isl.Id == 0))
						if o.C.Verbose && isl.Id == 0 {
							o.print_time(time, report)
						}
					}
					done <- 1
				}(o.Islands[i])
			}
			for i := 0; i < nislands; i++ {
				<-done
			}
		} else {
			for _, isl := range o.Islands {
				for time := t; time < tmig; time++ {
					report := o.calc_report(time)
					isl.Run(time, report, (o.C.Verbose && isl.Id == 0))
					if o.C.Verbose && isl.Id == 0 {
						o.print_time(time, report)
					}
				}
			}
		}

		// update time
		t = tmig
		tmig += o.C.Dtmig
		if tmig > o.C.Tf {
			tmig = o.C.Tf
			nomig = true
		}

		// skip migration
		if nomig {
			continue
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
					send, _ := IndCompareDet(Bbest, Aworst)
					if send {
						receiveFrom[i][k] = j // i gets individual from j
						k++
					}
				}
			}
		}

		// migration
		if o.C.Verbose {
			io.Pfyel(" %d", t)
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
	if o.C.Verbose {
		io.Pf("\n... end\n\n")
	}

	// write reports
	if o.C.DoReport {
		for _, isl := range o.Islands {
			isl.SaveReport(o.C.Verbose)
		}
	}

	// plot evolution
	if o.C.DoPlot {
		for i, isl := range o.Islands {
			PlotOvs(isl, ".eps", "", o.C.PltTi, o.C.PltTf, i == 0, i == o.C.Nisl-1)
		}
		for i, isl := range o.Islands {
			PlotOor(isl, ".eps", "", o.C.PltTi, o.C.PltTf, i == 0, i == o.C.Nisl-1)
		}
	}
	return
}

// ResetAllPop resets/re-generates all populations in all islands
func (o *Evolver) ResetAllPop() {
	for id, isl := range o.Islands {
		isl.Pop = o.C.PopFltGen(id, o.C, o.C.RangeFlt)
		isl.CalcOvs(isl.Pop, 0)
		isl.CalcDemeritsAndSort(isl.Pop)
	}
}

// FindBestFromAll finds best individual from all islands
//  Output: o.Best will point to the best individual
func (o *Evolver) FindBestFromAll() {
	if len(o.Islands) < 1 {
		return
	}
	o.Best = o.Islands[0].Pop[0]
	for i := 1; i < o.C.Nisl; i++ {
		_, other_is_better := IndCompareDet(o.Best, o.Islands[i].Pop[0])
		if other_is_better {
			o.Best = o.Islands[i].Pop[0]
		}
	}
}

// GetFeasible returns all feasible individuals from all islands
func (o *Evolver) GetFeasible() (feasible []*Individual) {
	for _, isl := range o.Islands {
		for _, ind := range isl.Pop {
			unfeasible := false
			for _, oor := range ind.Oors {
				if oor > 0 {
					unfeasible = true
				}
			}
			if !unfeasible {
				feasible = append(feasible, ind)
			}
		}
	}
	return
}

// GetResults returns all ovas and oors from a subset of individuals
//  Output:
//   ovas -- [len(subset)][nova] objective values
//   oors -- [len(subset)][noor] out-of-range values
func (o *Evolver) GetResults(subset []*Individual) (ovas, oors [][]float64) {
	ninds := len(subset)
	ovas = utl.DblsAlloc(ninds, o.C.Nova)
	oors = utl.DblsAlloc(ninds, o.C.Noor)
	for i, ind := range subset {
		for j := 0; j < o.C.Nova; j++ {
			ovas[i][j] = ind.Ovas[j]
		}
		for j := 0; j < o.C.Noor; j++ {
			oors[i][j] = ind.Oors[j]
		}
	}
	return
}

// GetParetoFront returns all feasible individuals on the Pareto front
// Note: input data can be obtained from GetFeasible and GetResults
func (o *Evolver) GetParetoFront(feasible []*Individual, ovas, oors [][]float64) (ovafront, oorfront []*Individual) {
	chk.IntAssert(len(feasible), len(ovas))
	ovaf := utl.ParetoFront(ovas)
	ovafront = make([]*Individual, len(ovaf))
	for i, id := range ovaf {
		ovafront[i] = feasible[id]
	}
	if len(oors) > 0 {
		chk.IntAssert(len(feasible), len(oors))
		oorf := utl.ParetoFront(oors)
		oorfront = make([]*Individual, len(oorf))
		for i, id := range oorf {
			oorfront[i] = feasible[id]
		}
	}
	return
}

// GetFrontOvas collects 2 ova results from Pareto front
//  Input:
//   r and s -- 2 selected objective functions; e.g. r=0 and s=1 for 2D problems
func (o *Evolver) GetFrontOvas(r, s int, front []*Individual) (x, y []float64) {
	x = make([]float64, len(front))
	y = make([]float64, len(front))
	for i, ind := range front {
		x[i] = ind.Ovas[r]
		y[i] = ind.Ovas[s]
	}
	return
}

//func (o *Evolver) GetCompromise(feasible []*Individual) (xova, yova, xoor, yoor []float64) {
//}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

func (o Evolver) calc_report(t int) bool {
	return t%o.C.Dtout == 0
}

func (o Evolver) print_legend() {
	io.Pf("\nLEGEND\n")
	io.Pfgrey(" 00 -- generation number (time)\n")
	io.Pfblue(" 00 -- reporting time\n")
	io.Pf(" 00 -- prescribed regeneration time\n")
	io.Pfyel(" 00 -- migration time\n")
	io.Pfmag("  . -- automatic regeneration time to improve diversity\n")
}

func (o Evolver) print_time(time int, report bool) {
	io.Pf(" ")
	if report {
		io.Pfblue("%v", time)
		return
	}
	io.Pfgrey("%v", time)
}

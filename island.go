// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/graph"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Island holds one population and performs the reproduction operation
type Island struct {

	// input
	Id  int         // index of this island
	C   *ConfParams // configuration parameters
	Pop Population  // pointer to current population
	Bkp Population  // backup population

	// results
	Report   bytes.Buffer // buffer to report results
	OutOvas  [][]float64  // [nova][ntimes] best objective values collected from multiple calls to SelectReprodAndRegen
	OutOors  [][]float64  // [noor][ntimes] best out-of-range values collected from multiple calls to SelectReprodAndRegen
	OutTimes []float64    // [ntimes] times corresponding to OutOvas and OutOors

	// auxiliary internal data
	ovamin  []float64   // min ovas
	ovamax  []float64   // max ovas
	oormin  []float64   // min oors
	oormax  []float64   // max oors
	ovas    [][]float64 // all ova values
	oors    [][]float64 // all oor values
	sovas   [][]float64 // scaled ova values
	soors   [][]float64 // scaled oor values
	fitness []float64   // all fitness values
	prob    []float64   // probabilities
	cumprob []float64   // cumulated probabilities
	selinds []int       // indices of selected individuals
	A, B    []int       // indices of selected parents

	// for statistics
	allbases [][]float64 // [ngenes*nbases][ninds] all bases
	devbases []float64   // [ngenes*nbases] deviations of bases
	larbases []float64   // [ngenes*nbases] largest bases; max(abs(bases))

	// for crowding
	indices   []int         // [ninds]
	crowds    [][]int       // [ninds/crowd_size][crowd_size]
	distR1    [][]float64   // [crowd_size][cowd_size] dist for round 1
	distR2    [][]float64   // [crowd_size][(crowd_size-1)*2] dist for round 2
	matchR1   graph.Munkres // matches for round 1
	matchR2   graph.Munkres // matches for round 2
	winners   []*Individual // winners
	offspring []*Individual // offspring
	nextround []int         // next tournament

	// limits
	intXmin, intXmax []int     // int genes range
	fltXmin, fltXmax []float64 // flt genes range
}

// NewIsland creates a new island
func NewIsland(id int, C *ConfParams) (o *Island) {

	// check
	if C.Ninds < 2 || (C.Ninds%2 != 0) {
		chk.Panic("size of population must be even and greater than 2. C.Ninds = %d is invalid", C.Ninds)
	}
	if C.OvaOor == nil {
		chk.Panic("objective function (OvaOor) must be non nil")
	}

	// allocate island
	o = new(Island)
	o.Id = id
	o.C = C

	// create population
	if o.C.PopIntGen != nil {
		o.Pop = o.C.PopIntGen(id, o.C)
	}
	if o.C.PopFltGen != nil {
		o.Pop = o.C.PopFltGen(id, o.C)
	}
	if o.C.PopStrGen != nil {
		o.Pop = o.C.PopStrGen(id, o.C)
	}
	if o.C.PopKeyGen != nil {
		o.Pop = o.C.PopKeyGen(id, o.C)
	}
	if o.C.PopBytGen != nil {
		o.Pop = o.C.PopBytGen(id, o.C)
	}
	if o.C.PopFunGen != nil {
		o.Pop = o.C.PopFunGen(id, o.C)
	}
	if len(o.Pop) != o.C.Ninds {
		chk.Panic("generation of population failed:\nat least one generator function in Params must be non nil")
	}

	// copy population
	o.Bkp = o.Pop.GetCopy()

	// auxiliary data
	o.ovamin = make([]float64, o.C.Nova)
	o.ovamax = make([]float64, o.C.Nova)
	o.oormin = make([]float64, o.C.Noor)
	o.oormax = make([]float64, o.C.Noor)
	o.ovas = la.MatAlloc(o.C.Nova, o.C.Ninds)
	o.oors = la.MatAlloc(o.C.Noor, o.C.Ninds)
	o.sovas = la.MatAlloc(o.C.Nova, o.C.Ninds)
	o.soors = la.MatAlloc(o.C.Noor, o.C.Ninds)
	o.fitness = make([]float64, o.C.Ninds)
	o.prob = make([]float64, o.C.Ninds)
	o.cumprob = make([]float64, o.C.Ninds)
	o.selinds = make([]int, o.C.Ninds)
	o.A = make([]int, o.C.Ninds/2)
	o.B = make([]int, o.C.Ninds/2)

	// compute objective values, demerits, and sort population
	o.CalcOvs(o.Pop, 0)
	o.CalcDemeritsAndSort(o.Pop)

	// results
	o.OutOvas = la.MatAlloc(o.C.Nova, o.C.Tf)
	o.OutOors = la.MatAlloc(o.C.Noor, o.C.Tf)
	o.OutTimes = make([]float64, o.C.Tf)
	for i := 0; i < o.C.Nova; i++ {
		o.OutOvas[i][0] = o.Pop[0].Ovas[i]
	}
	for i := 0; i < o.C.Noor; i++ {
		o.OutOors[i][0] = o.Pop[0].Oors[i]
	}

	// stat
	nflts := o.Pop[0].Nfltgenes * o.Pop[0].Nbases
	if o.Pop[0].Nfltgenes > 0 {
		o.allbases = la.MatAlloc(nflts, o.C.Ninds)
		o.devbases = make([]float64, nflts)
		o.larbases = make([]float64, nflts)
	}

	// for crowding
	n := o.C.CrowdSize
	m := (o.C.CrowdSize - 1) * 2
	if o.C.Ninds%n > 0 {
		chk.Panic("number of individuals must be multiple of crowd size")
	}
	o.indices = utl.IntRange(o.C.Ninds)
	o.crowds = utl.IntsAlloc(o.C.Ninds/n, n)
	o.distR1 = la.MatAlloc(n, m)
	o.matchR1.Init(n, m)
	if m-n > 0 {
		o.distR2 = la.MatAlloc(n, m-n)
		o.matchR2.Init(n, m-n)
		o.nextround = make([]int, m-n)
	}
	o.winners = make([]*Individual, n)
	o.offspring = make([]*Individual, m)
	for i := 0; i < n; i++ {
		o.winners[i] = o.Pop[0].GetCopy()
	}
	for i := 0; i < m; i++ {
		o.offspring[i] = o.Pop[0].GetCopy()
		for j := 0; j < len(o.offspring[i].Floats); j++ {
			o.offspring[i].Floats[j] = 0
		}
	}

	// limits
	nints := len(o.Pop[0].Ints)
	if nints > 0 {
		o.intXmin = make([]int, nints)
		o.intXmax = make([]int, nints)
	}
	if nflts > 0 {
		o.fltXmin = make([]float64, nflts)
		o.fltXmax = make([]float64, nflts)
	}
	return
}

// CalcOvs computes objective and out-of-range values
func (o *Island) CalcOvs(pop Population, time int) {
	for _, ind := range pop {
		o.C.OvaOor(ind, o.Id, time, &o.Report)
		for _, oor := range ind.Oors {
			if oor < 0 {
				chk.Panic("out-of-range values must be positive (or zero) indicating the positive distance to constraints. oor=%g is invalid", oor)
			}
		}
	}
}

// CalcDemeritsAndSort computes demerits and sort population
func (o *Island) CalcDemeritsAndSort(pop Population) {
	for i, ind := range pop {
		for j := 0; j < o.C.Nova; j++ {
			o.ovas[j][i] = ind.Ovas[j]
		}
		for j := 0; j < o.C.Noor; j++ {
			o.oors[j][i] = ind.Oors[j]
		}
	}
	for i := 0; i < o.C.Nova; i++ {
		o.ovamin[i], o.ovamax[i] = utl.Scaling(o.sovas[i], o.ovas[i], 0, 1e-16, false, true)
	}
	for i := 0; i < o.C.Noor; i++ {
		o.oormin[i], o.oormax[i] = utl.Scaling(o.soors[i], o.oors[i], 0, 1e-16, false, true)
	}
	for i, ind := range pop {
		ind.Demerit = 0
		for j := 0; j < o.C.Nova; j++ {
			ind.Demerit += o.sovas[j][i]
		}
	}
	shift := 2.0
	for i, ind := range pop {
		firstOor := true
		for j := 0; j < o.C.Noor; j++ {
			if ind.Oors[j] > 0 {
				if firstOor {
					ind.Demerit = shift
					firstOor = false
				}
				ind.Demerit += o.soors[j][i]
			}
		}
	}
	pop.Sort()
}

// Run runs evolutionary process
func (o *Island) Run(time int, doreport, verbose bool) {

	// run
	switch o.C.GAtype {
	case "crowd":
		o.update_crowding(time)
	default:
		o.update_standard(time)
	}

	// swap populations (Pop will always point to current one)
	o.Pop, o.Bkp = o.Bkp, o.Pop
	o.CalcDemeritsAndSort(o.Pop)

	// elitism
	if o.C.Elite {
		prev_best, cur_worst := o.Bkp[0], o.Pop[o.C.Ninds-1]
		prev_dominates, _ := IndCompareDet(prev_best, cur_worst)
		if prev_dominates {
			prev_best.CopyInto(cur_worst)
		}
	}

	// statistics and regeneration of float-point individuals
	var averho float64
	if o.Pop[0].Nfltgenes > 0 {
		_, averho, _, _ = o.FltStat()
		homogeneous := averho < o.C.RegTol
		if homogeneous {
			o.Regenerate(time)
			if doreport {
				io.Ff(&o.Report, "time=%d: regeneration\n", time)
			}
			if verbose {
				io.Pfmag(" .")
			}
		}
	}

	// report
	if doreport {
		o.WritePopToReport(time, averho)
	}

	// post-process
	if o.C.PostProc != nil {
		o.C.PostProc(o.Id, time, o.Pop)
	}

	// results
	for i := 0; i < o.C.Nova; i++ {
		o.OutOvas[i][time] = o.Pop[0].Ovas[i]
	}
	for i := 0; i < o.C.Noor; i++ {
		o.OutOors[i][time] = o.Pop[0].Oors[i]
	}
	o.OutTimes[time] = float64(time)
}

// update_crowding runs the evolutionary process with niching via crowding and tournament selection
func (o *Island) update_crowding(time int) {

	// select groups (crowds)
	rnd.IntGetGroups(o.crowds, o.indices)

	// compute float gene limits
	o.calc_float_lims()

	// auxiliary variables
	n := o.C.CrowdSize
	m := (o.C.CrowdSize - 1) * 2
	ncrowd := len(o.crowds)

	// run tournaments
	for icrowd, crowd := range o.crowds {

		// crossover, mutation and new objective values
		for r := 0; r < n-1; r++ {
			i, j := r, r+1
			k, l := r*2, r*2+1
			I, J := crowd[i], crowd[j]
			A, B := o.Pop[I], o.Pop[J]
			a, b := o.offspring[k], o.offspring[l]
			if o.C.DiffEvol {
				jcrowd := (icrowd + 1) % ncrowd
				C, D := o.Pop[o.crowds[jcrowd][0]], o.Pop[o.crowds[jcrowd][1]]
				o.diff_evol_crossover(a, b, A, B, C, D)
			} else {
				IndCrossover(a, b, A, B, time, &o.C.Ops)
			}
			IndMutation(a, time, &o.C.Ops)
			IndMutation(b, time, &o.C.Ops)
			o.C.OvaOor(a, o.Id, time+1, &o.Report)
			o.C.OvaOor(b, o.Id, time+1, &o.Report)
		}

		// round 1: compute distances
		for i := 0; i < n; i++ {
			I := crowd[i]
			A := o.Pop[I]
			for j := 0; j < m; j++ {
				B := o.offspring[j]
				o.distR1[i][j] = IndDistance(A, B, o.intXmin, o.intXmax, o.fltXmin, o.fltXmax)
			}
		}

		// round 1: match competitors
		o.matchR1.SetCostMatrix(o.distR1)
		o.matchR1.Run()

		// compute next round
		k := 0
		for i := 0; i < m; i++ {
			if utl.IntIndexSmall(o.matchR1.Links, i) < 0 {
				o.nextround[k] = i
				k++
			}
		}

		// round 1: tournament
		for i := 0; i < n; i++ {
			I := crowd[i]
			j := o.matchR1.Links[i]
			A, B := o.Pop[I], o.offspring[j]
			o.tournament(A, B, I)
		}

		// next round
		if m-n > 0 {

			// round 2: compute distances
			for i := 0; i < n; i++ {
				I := crowd[i]
				A := o.Bkp[I]
				for j := 0; j < m-n; j++ {
					J := o.nextround[j]
					B := o.offspring[J]
					o.distR2[i][j] = IndDistance(A, B, o.intXmin, o.intXmax, o.fltXmin, o.fltXmax)
				}
			}

			// round 2: match competitors
			o.matchR2.SetCostMatrix(o.distR2)
			o.matchR2.Run()

			// round 2: tournament
			for i := 0; i < n; i++ {
				I := crowd[i]
				k := o.matchR2.Links[i]
				if k >= 0 {
					j := o.nextround[k]
					A, B := o.Bkp[I], o.offspring[j]
					o.tournament(A, B, I)
				}
			}
		}
	}
}

// diff_evol_crossover implements the differential-evolution crossover
// TODO: move this to operators file
func (o *Island) diff_evol_crossover(a, b, A, B, C, D *Individual) {
	nflts := len(A.Floats)
	sa := rnd.Int(0, nflts-1)
	sb := rnd.Int(0, nflts-1)
	var x float64
	for s := 0; s < nflts; s++ {

		// a
		if rnd.FlipCoin(o.C.Ops.DEpc) || s == sa {
			x = B.Floats[s] + o.C.Ops.DEmult*(C.Floats[s]-D.Floats[s])
		} else {
			x = A.Floats[s]
		}
		a.Floats[s] = o.C.Ops.EnforceRange(s, x)

		// b
		if rnd.FlipCoin(o.C.Ops.DEpc) || s == sb {
			x = A.Floats[s] + o.C.Ops.DEmult*(C.Floats[s]-D.Floats[s])
		} else {
			x = B.Floats[s]
		}
		b.Floats[s] = o.C.Ops.EnforceRange(s, x)
	}
}

// tournament runs game between A and B
func (o *Island) tournament(A, B *Individual, saveInto int) {

	// probabilistic
	if o.C.CompProb {
		if IndCompareProb(A, B, o.C.ParetoPhi) {
			A.CopyInto(o.Bkp[saveInto]) // A wins
			return
		}
		B.CopyInto(o.Bkp[saveInto]) // B wins
		return
	}

	// deterministic
	A_dom, B_dom := IndCompareDet(A, B)
	if A_dom {
		A.CopyInto(o.Bkp[saveInto]) // A wins
		return
	}
	if B_dom {
		B.CopyInto(o.Bkp[saveInto]) // B wins
		return
	}
	if rnd.FlipCoin(0.5) { // tie => roll dice
		A.CopyInto(o.Bkp[saveInto]) // A wins by chance
		return
	}
	B.CopyInto(o.Bkp[saveInto]) // B wins by chance
}

// update_standard performs the selection, reproduction and regeneration processes
//  Note: this function considers a SORTED population already
func (o *Island) update_standard(time int) {

	// fitness
	ninds := len(o.Pop)
	var sumfit float64
	if o.C.Rnk { // ranking
		sp := o.C.RnkSp
		for i := 0; i < ninds; i++ {
			o.fitness[i] = 2.0 - sp + 2.0*(sp-1.0)*float64(ninds-i-1)/float64(ninds-1)
			sumfit += o.fitness[i]
		}
	} else {
		mindem := o.Pop[0].Demerit
		maxdem := mindem
		for i := 0; i < ninds; i++ {
			mindem = utl.Min(mindem, o.Pop[i].Demerit)
			maxdem = utl.Max(maxdem, o.Pop[i].Demerit)
		}
		for i, ind := range o.Pop {
			o.fitness[i] = (maxdem - ind.Demerit) / (maxdem - mindem)
			sumfit += o.fitness[i]
		}
	}

	// probabilities
	for i := 0; i < ninds; i++ {
		o.prob[i] = o.fitness[i] / sumfit
		if i == 0 {
			o.cumprob[i] = o.prob[i]
		} else {
			o.cumprob[i] = o.cumprob[i-1] + o.prob[i]
		}
	}

	// selection
	if o.C.Rws {
		RouletteSelect(o.selinds, o.cumprob, nil)
	} else {
		SUSselect(o.selinds, o.cumprob, -1)
	}
	FilterPairs(o.A, o.B, o.selinds)

	// reproduction
	h := ninds / 2
	for i := 0; i < ninds/2; i++ {
		IndCrossover(o.Bkp[i], o.Bkp[h+i], o.Pop[o.A[i]], o.Pop[o.B[i]], time, &o.C.Ops)
		IndMutation(o.Bkp[i], time, &o.C.Ops)
		IndMutation(o.Bkp[h+i], time, &o.C.Ops)
	}

	// compute objective values
	o.CalcOvs(o.Bkp, time+1) // +1 => this is an updated generation
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// Regenerate regenerates population with basis on best individual(s)
func (o *Island) Regenerate(time int) {
	ninds := len(o.Pop)
	start := ninds - int(o.C.RegPct*float64(ninds))
	for i := start; i < ninds; i++ {
		for j := 0; j < o.Pop[i].Nfltgenes; j++ {
			xmin, xmax := o.C.RangeFlt[j][0], o.C.RangeFlt[j][1]
			o.Pop[i].SetFloat(j, rnd.Float64(xmin, xmax))
		}
	}
	o.CalcOvs(o.Pop, time)
	o.CalcDemeritsAndSort(o.Pop)
	return
}

// FltStat computes some statistic information with float-point individuals
//  rho (ρ) is a normalised quantity measuring the deviation of bases of each gene
func (o *Island) FltStat() (minrho, averho, maxrho, devrho float64) {
	ngenes, nbases := o.Pop[0].Nfltgenes, o.Pop[0].Nbases
	for k, ind := range o.Pop {
		for i := 0; i < ngenes; i++ {
			for j := 0; j < nbases; j++ {
				x := ind.Floats[i*nbases+j]
				o.allbases[i*nbases+j][k] = x
				if k == 0 {
					o.larbases[i*nbases+j] = math.Abs(x)
				} else {
					o.larbases[i*nbases+j] = utl.Max(o.larbases[i*nbases+j], math.Abs(x))
				}
			}
		}
	}
	for i := 0; i < ngenes; i++ {
		for j := 0; j < nbases; j++ {
			normfactor := 1.0 + o.larbases[i*nbases+j]
			o.devbases[i*nbases+j] = rnd.StatDev(o.allbases[i*nbases+j], o.C.UseStdDev) / normfactor
		}
	}
	minrho, averho, maxrho, devrho = rnd.StatBasic(o.devbases, o.C.UseStdDev)
	return
}

// WritePopToReport writes population to report
func (o *Island) WritePopToReport(time int, averho float64) {
	io.Ff(&o.Report, "time=%d averho=%g\n", time, averho)
	o.Report.Write(o.Pop.Output(o.C).Bytes())
}

// SaveReport saves report to file
func (o Island) SaveReport(verbose bool) {
	dosave := o.C.FnKey != ""
	if dosave {
		if o.C.DirOut == "" {
			o.C.DirOut = "/tmp/goga"
		}
		if verbose {
			io.WriteFileVD(o.C.DirOut, io.Sf("%s-%d.rpt", o.C.FnKey, o.Id), &o.Report)
			return
		}
		io.WriteFileD(o.C.DirOut, io.Sf("%s-%d.rpt", o.C.FnKey, o.Id), &o.Report)
	}
}

// calc_float_lims find float genes limits
func (o *Island) calc_float_lims() {
	for i, ind := range o.Pop {
		for j, x := range ind.Ints {
			if i == 0 {
				o.intXmin[j], o.intXmax[j] = x, x
			} else {
				o.intXmin[j] = utl.Imin(o.intXmin[j], x)
				o.intXmax[j] = utl.Imax(o.intXmax[j], x)
			}
		}
		for j, x := range ind.Floats {
			if i == 0 {
				o.fltXmin[j], o.fltXmax[j] = x, x
			} else {
				o.fltXmin[j] = utl.Min(o.fltXmin[j], x)
				o.fltXmax[j] = utl.Max(o.fltXmax[j], x)
			}
		}
	}
}

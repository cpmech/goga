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
	Nova     int          // number of ovas
	Noor     int          // number of oors
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
	indices []int         // [ninds]
	crowds  [][]int       // [ninds/crowd_size][crowd_size]
	dist    [][]float64   // [crowd_size][cowd_size]
	match   graph.Munkres // matches
}

// NewIsland creates a new island
func NewIsland(id, nova, noor int, C *ConfParams) (o *Island) {

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
		o.Pop = o.C.PopIntGen(o.C.Ninds, nova, noor, o.C.Nbases, o.C.Noise, o.C.PopGenArgs, o.C.RangeInt)
	}
	if o.C.PopOrdGen != nil {
		o.Pop = o.C.PopOrdGen(o.C.Ninds, nova, noor, o.C.Nbases, o.C.Noise, o.C.PopGenArgs, o.C.OrdNints)
	}
	if o.C.PopFltGen != nil {
		o.Pop = o.C.PopFltGen(o.C.Ninds, nova, noor, o.C.Nbases, o.C.Noise, o.C.PopGenArgs, o.C.RangeFlt)
	}
	if o.C.PopStrGen != nil {
		o.Pop = o.C.PopStrGen(o.C.Ninds, nova, noor, o.C.Nbases, o.C.Noise, o.C.PopGenArgs, o.C.PoolStr)
	}
	if o.C.PopKeyGen != nil {
		o.Pop = o.C.PopKeyGen(o.C.Ninds, nova, noor, o.C.Nbases, o.C.Noise, o.C.PopGenArgs, o.C.PoolKey)
	}
	if o.C.PopBytGen != nil {
		o.Pop = o.C.PopBytGen(o.C.Ninds, nova, noor, o.C.Nbases, o.C.Noise, o.C.PopGenArgs, o.C.PoolByt)
	}
	if o.C.PopFunGen != nil {
		o.Pop = o.C.PopFunGen(o.C.Ninds, nova, noor, o.C.Nbases, o.C.Noise, o.C.PopGenArgs, o.C.PoolFun)
	}
	if len(o.Pop) != o.C.Ninds {
		chk.Panic("generation of population failed:\nat least one generator function in Params must be non nil")
	}

	// copy population
	o.Bkp = o.Pop.GetCopy()

	// auxiliary data
	o.Nova = len(o.Pop[0].Ovas)
	o.Noor = len(o.Pop[0].Oors)
	o.ovamin = make([]float64, o.Nova)
	o.ovamax = make([]float64, o.Nova)
	o.oormin = make([]float64, o.Noor)
	o.oormax = make([]float64, o.Noor)
	o.ovas = la.MatAlloc(o.Nova, o.C.Ninds)
	o.oors = la.MatAlloc(o.Noor, o.C.Ninds)
	o.sovas = la.MatAlloc(o.Nova, o.C.Ninds)
	o.soors = la.MatAlloc(o.Noor, o.C.Ninds)
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
	o.OutOvas = la.MatAlloc(o.Nova, o.C.Tf)
	o.OutOors = la.MatAlloc(o.Noor, o.C.Tf)
	o.OutTimes = make([]float64, o.C.Tf)
	for i := 0; i < o.Nova; i++ {
		o.OutOvas[i][0] = o.Pop[0].Ovas[i]
	}
	for i := 0; i < o.Noor; i++ {
		o.OutOors[i][0] = o.Pop[0].Oors[i]
	}

	// stat
	if o.Pop[0].Nfltgenes > 0 {
		size := o.Pop[0].Nfltgenes * o.Pop[0].Nbases
		o.allbases = la.MatAlloc(size, o.C.Ninds)
		o.devbases = make([]float64, size)
		o.larbases = make([]float64, size)
	}

	// for crowding
	n := o.C.CrowdSize
	if o.C.Ninds%n > 0 {
		chk.Panic("number of individuals must be multiple of crowd size")
	}
	o.indices = utl.IntRange(o.C.Ninds)
	o.crowds = utl.IntsAlloc(o.C.Ninds/n, n)
	o.dist = la.MatAlloc(n, n)
	o.match.Init(n, n)
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
		for j := 0; j < o.Nova; j++ {
			o.ovas[j][i] = ind.Ovas[j]
		}
		for j := 0; j < o.Noor; j++ {
			o.oors[j][i] = ind.Oors[j]
		}
	}
	for i := 0; i < o.Nova; i++ {
		o.ovamin[i], o.ovamax[i] = utl.Scaling(o.sovas[i], o.ovas[i], 0, 1e-16, false, true)
	}
	for i := 0; i < o.Noor; i++ {
		o.oormin[i], o.oormax[i] = utl.Scaling(o.soors[i], o.oors[i], 0, 1e-16, false, true)
	}
	for i, ind := range pop {
		ind.Demerit = 0
		for j := 0; j < o.Nova; j++ {
			ind.Demerit += o.sovas[j][i]
		}
	}
	shift := 2.0
	for i, ind := range pop {
		firstOor := true
		for j := 0; j < o.Noor; j++ {
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
	case "sharing":
		o.update_sharing(time)
	default:
		o.update_standard(time)
	}

	// swap populations (Pop will always point to current one)
	o.Pop, o.Bkp = o.Bkp, o.Pop
	o.CalcDemeritsAndSort(o.Pop)

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
	for i := 0; i < o.Nova; i++ {
		o.OutOvas[i][time] = o.Pop[0].Ovas[i]
	}
	for i := 0; i < o.Noor; i++ {
		o.OutOors[i][time] = o.Pop[0].Oors[i]
	}
	o.OutTimes[time] = float64(time)
}

// update_crowding runs the evolutionary process with niching via crowding and tournament selection
func (o *Island) update_crowding(time int) {

	// select groups (crowds)
	rnd.IntGetGroups(o.crowds, o.indices)

	// run tournaments
	n := o.C.CrowdSize
	for _, crowd := range o.crowds {

		// crossover, mutation and new objective values
		for k := 0; k < n; k += 2 { // NOTE: for odd crowd_size, one child is overwritten
			i, j := k, (k+1)%n
			I, J := crowd[i], crowd[j]
			A, B := o.Pop[I], o.Pop[J]
			a, b := o.Bkp[I], o.Bkp[J]
			IndCrossover(a, b, A, B, time, &o.C.Ops)
			IndMutation(a, time, &o.C.Ops)
			IndMutation(b, time, &o.C.Ops)
			o.C.OvaOor(a, o.Id, time+1, &o.Report)
			o.C.OvaOor(b, o.Id, time+1, &o.Report)
		}

		// compute distances
		for i := 0; i < n; i++ {
			A := o.Pop[crowd[i]]
			for j := 0; j < n; j++ {
				a := o.Bkp[crowd[j]]
				o.dist[i][j] = IndDistance(A, a)
			}
		}

		// match competitors
		o.match.SetCostMatrix(o.dist)
		o.match.Run()

		// perform tournament
		for i := 0; i < n; i++ {
			j := o.match.Links[i]
			A, a := o.Pop[crowd[i]], o.Bkp[crowd[j]]
			if IndCompareProb(A, a, o.C.ParetoPhi) {
				A.CopyInto(a) // parent wins
			}
		}
	}
}

// update_sharing runs the evolutionary process with pareto-niching and sharing
func (o *Island) update_sharing(time int) {

	// selection
	nsample := int(o.C.ShSize * float64(o.C.Ninds))
	for k := 0; k < o.C.Ninds; k++ {
		sample := rnd.IntGetUniqueN(0, o.C.Ninds, nsample)
		pair := rnd.IntGetUniqueN(0, o.C.Ninds, 2)
		i, j := pair[0], pair[1]
		A, B := o.Pop[i], o.Pop[j]
		A_is_dominated, B_is_dominated := false, false
		for _, s := range sample {
			if s != i {
				_, other_dominates := IndCompareDet(A, o.Pop[s])
				if other_dominates {
					A_is_dominated = true
					break
				}
			}
		}
		for _, s := range sample {
			if s != j {
				_, other_dominates := IndCompareDet(B, o.Pop[s])
				if other_dominates {
					B_is_dominated = true
					break
				}
			}
		}
		if !A_is_dominated && B_is_dominated {
			o.selinds[k] = i
			continue
		}
		if A_is_dominated && !B_is_dominated {
			o.selinds[k] = j
			continue
		}
		shA := o.calc_sharing(i)
		shB := o.calc_sharing(j)
		if shA < shB {
			o.selinds[k] = i
		} else {
			o.selinds[k] = j
		}
	}

	// filter pairs
	FilterPairs(o.A, o.B, o.selinds)

	// reproduction
	h := o.C.Ninds / 2
	for i := 0; i < o.C.Ninds/2; i++ {
		IndCrossover(o.Bkp[i], o.Bkp[h+i], o.Pop[o.A[i]], o.Pop[o.B[i]], time, &o.C.Ops)
		IndMutation(o.Bkp[i], time, &o.C.Ops)
		IndMutation(o.Bkp[h+i], time, &o.C.Ops)
	}

	// compute objective values
	o.CalcOvs(o.Bkp, time+1) // +1 => this is an updated generation

	// elitism
	if o.C.Elite {
		iold, inew := o.Pop[0], o.Bkp[o.C.Ninds-1]
		old_dominates, _ := IndCompareDet(iold, inew)
		if old_dominates {
			iold.CopyInto(inew)
		}
	}
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

	// elitism
	if o.C.Elite {
		iold, inew := o.Pop[0], o.Bkp[o.C.Ninds-1]
		old_dominates, _ := IndCompareDet(iold, inew)
		if old_dominates {
			iold.CopyInto(inew)
		}
	}
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// Regenerate regenerates population with basis on best individual(s)
func (o *Island) Regenerate(time int) {
	chk.Panic("regenerate is deactivated")
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
	o.Report.Write(o.Pop.Output(nil, o.C.ShowOor, o.C.ShowBases, o.C.ShowNinds).Bytes())
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

// calc_sharing computes sharing coefficient
func (o *Island) calc_sharing(idxind int) (sh float64) {
	A := o.Pop[idxind]
	var dova, door, d float64
	for _, B := range o.Pop {
		if o.C.ShPhen {
			d = IndDistance(A, B)
		} else {
			dova, door = 0, 0
			for i := 0; i < o.Nova; i++ {
				dova += math.Pow((A.Ovas[i]-B.Ovas[i])/(1+o.ovamax[i]-o.ovamin[i]), 2.0)
			}
			for i := 0; i < o.Noor; i++ {
				door += math.Pow((A.Oors[i]-B.Oors[i])/(1+o.oormax[i]-o.oormin[i]), 2.0)
			}
			d = math.Sqrt(dova) + math.Sqrt(door)
		}
		if d < o.C.ShSig {
			sh += 1.0 - math.Pow(d/o.C.ShSig, o.C.ShAlp)
		}
	}
	return
}

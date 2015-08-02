// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// ObjFunc_t defines the template for the objective function
type ObjFunc_t func(ind *Individual, idIsland, time int, report *bytes.Buffer) (ova, oor float64)

// Island holds one population and performs the reproduction operation
type Island struct {

	// input
	Id        int         // index of this island
	C         *ConfParams // configuration parameters
	Pop       Population  // pointer to current population
	BkpPop    Population  // backup population
	OvOorFunc ObjFunc_t   // function to compute objective and out-of-range values
	BingoGrid *Bingo      // bingo for regeneration with initial values from grid
	BingoBest *Bingo      // bingo for regeneration with values recomputed based on best individual

	// results
	Report bytes.Buffer // buffer to report results
	OVA    []float64    // best objective values collected from multiple calls to SelectReprodAndRegen
	OOR    []float64    // best out-of-range values collected from multiple calls to SelectReprodAndRegen

	// auxiliary internal data
	ovas    []float64 // all ova values
	oors    []float64 // all oor values
	sovas   []float64 // scaled ova values
	soors   []float64 // scaled oor values
	fitness []float64 // all fitness values
	prob    []float64 // probabilities
	cumprob []float64 // cumulated probabilities
	selinds []int     // indices of selected individuals
	A, B    []int     // indices of selected parents

	// for statistics
	maxabsgene []float64   // [ngenes] max absolute values of genes
	fltbases   [][]float64 // [ngenes*nbases][ninds] all bases
	devbases   []float64   // [ngenes*nbases] deviations of bases
}

// NewIsland allocates a new island but with a give population already allocated
// Input:
//  id     -- index of this island
//  pop    -- the population
//  ovfunc -- objective function
//  bingo  -- structure needed for regeneration of individuals
func NewIsland(id int, C *ConfParams, pop Population, ovfunc ObjFunc_t, bingo *Bingo) (o *Island) {

	// check
	ninds := len(pop)
	chk.IntAssert(C.Ninds, ninds)
	if ninds%2 != 0 {
		chk.Panic("size of population must be even")
	}

	// allocate
	o = new(Island)
	o.Id = id
	o.C = C
	o.Pop = pop
	o.BkpPop = pop.GetCopy()
	o.OvOorFunc = ovfunc
	o.BingoGrid = bingo
	o.BingoBest = bingo.GetCopy()

	// auxiliary data
	o.ovas = make([]float64, ninds)
	o.oors = make([]float64, ninds)
	o.sovas = make([]float64, ninds)
	o.soors = make([]float64, ninds)
	o.fitness = make([]float64, ninds)
	o.prob = make([]float64, ninds)
	o.cumprob = make([]float64, ninds)
	o.selinds = make([]int, ninds)
	o.A = make([]int, ninds/2)
	o.B = make([]int, ninds/2)

	// compute objective values, demerits, and sort population
	o.CalcOvs(o.Pop, 0)
	o.CalcDemeritsAndSort(o.Pop)

	// results
	o.OVA = make([]float64, o.C.Tf)
	o.OOR = make([]float64, o.C.Tf)
	o.OVA[0] = o.Pop[0].Ova
	o.OOR[0] = o.Pop[0].Oor

	// for statistics
	nfltgenes := o.Pop[0].Nfltgenes
	if nfltgenes > 0 {
		nbases := o.Pop[0].Nbases
		o.maxabsgene = make([]float64, nfltgenes)
		o.fltbases = la.MatAlloc(nfltgenes*nbases, ninds)
		o.devbases = make([]float64, nfltgenes*nbases)
	}
	return
}

// CalcOvs computes objective and out-of-range values
func (o *Island) CalcOvs(pop Population, time int) {
	for _, ind := range pop {
		ova, oor := o.OvOorFunc(ind, o.Id, time, &o.Report)
		if oor < 0 {
			chk.Panic("out-of-range values must be positive (or zero) indicating the positive distance to constraints. oor=%g is invalid", oor)
		}
		if oor > 0 { // infeasible solutions (out-of-range)
			ind.Ova = 0 // not used
			ind.Oor = oor
		} else { // feasible solutions
			ind.Ova = ova
			ind.Oor = 0 // not used
		}
	}
}

// CalcDemeritsAndSort computes demerits and sort population
func (o *Island) CalcDemeritsAndSort(pop Population) {

	// ovs and oors
	var iova, ioor int // indices of individuals with ova and with oor, respectively
	for _, ind := range pop {
		if ind.Oor > 0 { // infeasible solutions (out-of-range)
			o.oors[ioor] = ind.Oor
			ioor++
		} else { // feasible solutions
			o.ovas[iova] = ind.Ova
			iova++
		}
	}

	// scaled ovs and oors
	utl.Scaling(o.sovas[:iova], o.ovas[:iova], 0, 1e-16, false, true)
	utl.Scaling(o.soors[:ioor], o.oors[:ioor], 2, 1e-16, false, true)

	// set demerits in individuals (loop with the same comparisons as before)
	ioor, iova = 0, 0
	for _, ind := range pop {
		if ind.Oor > 0 { // infeasible solutions (out-of-range)
			ind.Demerit = o.soors[ioor]
			ioor++
		} else { // feasible solutions
			ind.Demerit = o.sovas[iova]
			iova++
		}
	}

	// sort population
	pop.Sort()
}

// SelectReprodAndRegen performs the selection, reproduction and regeneration processes
// It also peforms the output to files.
//  Note: this function considers a SORTED population already
func (o *Island) SelectReprodAndRegen(time int, doregen, doreport bool) {

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
			mindem = min(mindem, o.Pop[i].Demerit)
			maxdem = max(maxdem, o.Pop[i].Demerit)
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
		Crossover(o.BkpPop[i], o.BkpPop[h+i], o.Pop[o.A[i]], o.Pop[o.B[i]], o.C.CxNcuts, o.C.CxCuts, o.C.CxProbs, o.C.CxIntFunc, o.C.CxFltFunc, o.C.CxStrFunc, o.C.CxKeyFunc, o.C.CxBytFunc, o.C.CxFunFunc)
		Mutation(o.BkpPop[i], o.C.MtNchanges, o.C.MtProbs, o.C.MtExtra, o.C.MtIntFunc, o.C.MtFltFunc, o.C.MtStrFunc, o.C.MtKeyFunc, o.C.MtBytFunc, o.C.MtFunFunc)
		Mutation(o.BkpPop[h+i], o.C.MtNchanges, o.C.MtProbs, o.C.MtExtra, o.C.MtIntFunc, o.C.MtFltFunc, o.C.MtStrFunc, o.C.MtKeyFunc, o.C.MtBytFunc, o.C.MtFunFunc)
	}

	// compute objective values, demerits, and sort population
	o.CalcOvs(o.BkpPop, time+1) // +1 => this is an updated generation
	o.CalcDemeritsAndSort(o.BkpPop)

	// elitism
	if o.C.Elite {
		iold, inew := o.Pop[0], o.BkpPop[ninds-1]
		docopy := iold.Compare(inew)
		if docopy {
			iold.CopyInto(inew)
			o.CalcDemeritsAndSort(o.BkpPop)
		}
	}

	// swap populations (Pop will always point to current one)
	o.Pop, o.BkpPop = o.BkpPop, o.Pop

	// statistics
	minrho, averho, maxrho, devrho := o.Stat()

	// regeneration
	homogeneous := averho < o.C.RegTol
	if homogeneous || doregen {
		basedOnBest := !homogeneous
		method := o.Regenerate(time, basedOnBest)
		if doreport {
			io.Ff(&o.Report, "time=%d: regeneration: method=%s\n", time, method)
		}
	}

	// report
	if doreport {
		io.Ff(&o.Report, "time=%d: homogeneous=%v minrho=%g averho=%g maxrho=%g devrho=%g\n", time, homogeneous, minrho, averho, maxrho, devrho)
		o.WritePopToReport(time)
	}

	// results
	// Note: sometimes the best ova may be zero when its oor is non-zero
	o.OVA[time] = o.Pop[0].Ova
	o.OOR[time] = o.Pop[0].Oor
	return
}

// Regenerate regenerates population with basis on best individual(s)
func (o *Island) Regenerate(time int, basedOnBest bool) (method string) {
	bingo := o.BingoGrid
	method = "lims"
	if basedOnBest || o.C.RegBest {
		method = "best"
		o.BingoBest.ResetBasedOnRef(time, o.Pop[0], o.C.RegMmin, o.C.RegMmax)
		bingo = o.BingoBest
	}
	ninds := len(o.Pop)
	start := ninds - int(o.C.RegPct*float64(ninds))
	for i := start; i < ninds; i++ {
		for j := 0; j < o.Pop[i].Nfltgenes; j++ {
			o.Pop[i].SetFloat(j, bingo.DrawFloat(i, j, ninds))
		}
	}
	o.CalcOvs(o.Pop, time)
	o.CalcDemeritsAndSort(o.Pop)
	return
}

// Stat computes some statistic information
//  rho (ρ) is a normalised quantity measuring the deviation of bases of each gene
//  Note: OoR individuals are excluded
func (o *Island) Stat() (minrho, averho, maxrho, devrho float64) {
	ngenes := o.Pop[0].Nfltgenes
	if ngenes < 1 {
		return
	}
	nbases := o.Pop[0].Nbases
	iova := 0
	for _, ind := range o.Pop {
		if ind.Oor > 0 && o.C.StatOorSkip { // skip oor individuals
			continue
		}
		for i := 0; i < ngenes; i++ {
			x := math.Abs(ind.GetFloat(i))
			if iova == 0 {
				o.maxabsgene[i] = x
			} else {
				o.maxabsgene[i] = max(o.maxabsgene[i], x)
			}
			for j := 0; j < nbases; j++ {
				o.fltbases[i*nbases+j][iova] = ind.Floats[i*nbases+j]
			}
		}
		iova++
	}
	if iova < 2 {
		return
	}
	for i := 0; i < ngenes; i++ {
		x := 1.0 + o.maxabsgene[i]
		for j := 0; j < nbases; j++ {
			o.devbases[i*nbases+j] = rnd.StatDev(o.fltbases[i*nbases+j][:iova], o.C.UseStdDev) / x
		}
	}
	minrho, averho, maxrho, devrho = rnd.StatBasic(o.devbases, o.C.UseStdDev)
	return
}

// PlotOvs plots objective values versus time
func (o Island) PlotOvs(ext, args string, t0, tf int, withtxt bool, numfmt string, first, last bool) {
	if o.C.DoPlot == false || o.C.FnKey == "" {
		return
	}
	if first {
		plt.SetForEps(0.75, 250)
	}
	var y []float64
	if tf == -1 {
		y = o.OVA[t0:]
		tf = len(o.OVA)
	} else {
		y = o.OVA[t0:tf]
	}
	n := len(y)
	T := utl.LinSpace(float64(t0), float64(tf), n)
	plt.Plot(T, y, args)
	if withtxt {
		plt.Text(T[0], y[0], io.Sf(numfmt, y[0]), "ha='left'")
		plt.Text(T[n-1], y[n-1], io.Sf(numfmt, y[n-1]), "ha='right'")
	}
	if last {
		plt.Gll("time", "objective value", "")
		plt.SaveD(o.C.DirOut, o.C.FnKey+"_ova"+ext)
	}
}

// PlotOor plots out-of-range values versus time
func (o Island) PlotOor(ext, args string, t0, tf int, withtxt bool, numfmt string, first, last bool) {
	if o.C.DoPlot == false || o.C.FnKey == "" {
		return
	}
	if first {
		plt.SetForEps(0.75, 250)
	}
	var y []float64
	if tf == -1 {
		y = o.OOR[t0:]
		tf = len(o.OOR)
	} else {
		y = o.OOR[t0:tf]
	}
	n := len(y)
	T := utl.LinSpace(float64(t0), float64(tf), n)
	plt.Plot(T, y, args)
	if withtxt {
		plt.Text(T[0], y[0], io.Sf(numfmt, y[0]), "ha='left'")
		plt.Text(T[n-1], y[n-1], io.Sf(numfmt, y[n-1]), "ha='right'")
	}
	if last {
		plt.Gll("time", "out-of-range value", "")
		plt.SaveD(o.C.DirOut, o.C.FnKey+"_oor"+ext)
	}
}

// WritePopToReport writes population to report
func (o *Island) WritePopToReport(time int) {
	io.Ff(&o.Report, "time=%d: population:\n", time)
	o.Report.Write(o.Pop.Output(nil, o.C.ShowBases).Bytes())
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

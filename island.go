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
type ObjFunc_t func(ind *Individual, idIsland, time int, report *bytes.Buffer)

// Island holds one population and performs the reproduction operation
type Island struct {

	// index
	Id int // index of this island

	// selection/reproduction
	UseRanking  bool    // use ranking for selection process
	RnkPressure float64 // ranking pressure
	Roulette    bool    // use roulette wheel selection; otherwise use stochastic-universal-sampling selection
	Elitism     bool    // perform elitism: keep at least one best individual from previous generation

	// regeneration
	RegenBest bool    // enforce that regeneration is always based on based individual, regardless the population is homogeneous or not
	RegenPct  float64 // percentage of individuals to be regenerated
	RegenMmin float64 // multiplier to decrease reference value; e.g. 0.1
	RegenMmax float64 // multiplier to increase reference value; e.g. 10.0

	// crossover
	CxNcuts   map[string]int         // crossover number of cuts for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxCuts    map[string][]int       // crossover specific cuts for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxProbs   map[string]float64     // crossover probabilities for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxFuncs   map[string]interface{} // crossover functions for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxIntFunc CxIntFunc_t
	CxFltFunc CxFltFunc_t
	CxStrFunc CxStrFunc_t
	CxKeyFunc CxKeyFunc_t
	CxBytFunc CxBytFunc_t
	CxFunFunc CxFunFunc_t

	// mutation
	MtNchanges map[string]int         // mutation number of changes for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtProbs    map[string]float64     // mutation probabilities for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtExtra    map[string]interface{} // mutation extra parameters for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	MtIntFunc  MtIntFunc_t            // mutation function
	MtFltFunc  MtFltFunc_t            // mutation function
	MtStrFunc  MtStrFunc_t            // mutation function
	MtKeyFunc  MtKeyFunc_t            // mutation function
	MtBytFunc  MtBytFunc_t            // mutation function
	MtFunFunc  MtFunFunc_t            // mutation function

	// input
	Pop       Population // pointer to current population
	BkpPop    Population // backup population
	ObjFunc   ObjFunc_t  // objective function
	BingoGrid *Bingo     // bingo for regeneration with initial values from grid
	BingoBest *Bingo     // bingo for regeneration with values recomputed based on best individual

	// results
	UseStdDev bool         // use standard deviation (σ) instead of average deviation in Stat
	ShowBases bool         // show also bases when printing results (if any)
	Report    bytes.Buffer // buffer to report results
	OVS       []float64    // best objective values collected from multiple calls to SelectAndReprod

	// auxiliary internal data
	fitsrnk []float64 // all fitness values computed by ranking
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
func NewIsland(id int, pop Population, ovfunc ObjFunc_t, bingo *Bingo) (o *Island) {

	// check
	ninds := len(pop)
	if ninds%2 != 0 {
		chk.Panic("size of population must be even")
	}

	// allocate
	o = new(Island)
	o.Id = id
	o.Pop = pop
	o.BkpPop = pop.GetCopy()
	o.ObjFunc = ovfunc
	o.BingoGrid = bingo
	o.BingoBest = bingo.GetCopy()

	// set default control values
	o.UseRanking = true
	o.RnkPressure = 1.2
	o.Elitism = true
	o.RegenBest = true
	o.RegenPct = 0.3
	o.RegenMmin = 0.1
	o.RegenMmax = 10.0

	// compute objective values
	for _, ind := range o.Pop {
		o.ObjFunc(ind, o.Id, 0, &o.Report)
	}

	// sort
	o.Pop.Sort()

	// results
	o.OVS = []float64{o.Pop[0].ObjValue}

	// auxiliary data
	o.fitsrnk = make([]float64, ninds)
	o.fitness = make([]float64, ninds)
	o.prob = make([]float64, ninds)
	o.cumprob = make([]float64, ninds)
	o.selinds = make([]int, ninds)
	o.A = make([]int, ninds/2)
	o.B = make([]int, ninds/2)

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

// SelectAndReprod performs the selection and reproduction processes
func (o *Island) SelectAndReprod(time int) {

	// fitness and probabilities
	ninds := len(o.Pop)
	sumfit := 0.0
	if o.UseRanking {
		sp := o.RnkPressure
		if sp < 1.0 || sp > 2.0 {
			sp = 1.2
		}
		for i := 0; i < ninds; i++ {
			o.fitness[i] = 2.0 - sp + 2.0*(sp-1.0)*float64(ninds-i-1)/float64(ninds-1)
			sumfit += o.fitness[i]
		}
	} else {
		ovmin, ovmax := o.Pop[0].ObjValue, o.Pop[0].ObjValue
		for _, ind := range o.Pop {
			ovmin = min(ovmin, ind.ObjValue)
			ovmax = max(ovmax, ind.ObjValue)
		}
		if math.Abs(ovmax-ovmin) < 1e-14 {
			for i := 0; i < ninds; i++ {
				o.fitness[i] = float64(i) / float64(ninds-1)
				sumfit += o.fitness[i]
			}
		} else {
			for i, ind := range o.Pop {
				o.fitness[i] = (ovmax - ind.ObjValue) / (ovmax - ovmin)
				sumfit += o.fitness[i]
			}
		}
	}
	for i := 0; i < ninds; i++ {
		o.prob[i] = o.fitness[i] / sumfit
	}
	CumSum(o.cumprob, o.prob)

	// selection
	if o.Roulette {
		RouletteSelect(o.selinds, o.cumprob, nil)
	} else {
		SUSselect(o.selinds, o.cumprob, -1)
	}
	FilterPairs(o.A, o.B, o.selinds)

	// reproduction
	h := ninds / 2
	for i := 0; i < ninds/2; i++ {
		Crossover(o.BkpPop[i], o.BkpPop[h+i], o.Pop[o.A[i]], o.Pop[o.B[i]], o.CxNcuts, o.CxCuts, o.CxProbs, o.CxIntFunc, o.CxFltFunc, o.CxStrFunc, o.CxKeyFunc, o.CxBytFunc, o.CxFunFunc)
		Mutation(o.BkpPop[i], o.MtNchanges, o.MtProbs, o.MtExtra, o.MtIntFunc, o.MtFltFunc, o.MtStrFunc, o.MtKeyFunc, o.MtBytFunc, o.MtFunFunc)
		Mutation(o.BkpPop[h+i], o.MtNchanges, o.MtProbs, o.MtExtra, o.MtIntFunc, o.MtFltFunc, o.MtStrFunc, o.MtKeyFunc, o.MtBytFunc, o.MtFunFunc)
	}

	// compute objective values
	for _, ind := range o.BkpPop {
		o.ObjFunc(ind, o.Id, time+1, &o.Report) // +1 => this is an updated generation
	}

	// sort
	o.BkpPop.Sort()

	// elitism
	if o.Elitism {
		if o.Pop[0].ObjValue < o.BkpPop[0].ObjValue {
			o.Pop[0].CopyInto(o.BkpPop[ninds-1])
			o.BkpPop.Sort()
		}
	}

	// swap populations (Pop will always point to current one)
	o.Pop, o.BkpPop = o.BkpPop, o.Pop

	// results
	o.OVS = append(o.OVS, o.Pop[0].ObjValue)
}

// Regenerate regenerates population with basis on best individual(s)
func (o *Island) Regenerate(time int, basedOnBest bool) {
	bingo := o.BingoGrid
	if basedOnBest || o.RegenBest {
		o.BingoBest.ResetBasedOnRef(time, o.Pop[0], o.RegenMmin, o.RegenMmax)
		bingo = o.BingoBest
	}
	ninds := len(o.Pop)
	start := ninds - int(o.RegenPct*float64(ninds))
	for i := start; i < ninds; i++ {
		for j := 0; j < o.Pop[i].Nfltgenes; j++ {
			o.Pop[i].SetFloat(j, bingo.DrawFloat(i, j, ninds))
		}
		o.ObjFunc(o.Pop[i], o.Id, time, nil)
	}
	o.Pop.Sort()
}

// Stat computes some statistic information
//  rho (ρ) is a normalised quantity measuring the deviation of bases of each gene
func (o *Island) Stat() (minrho, averho, maxrho, devrho float64) {
	ngenes := o.Pop[0].Nfltgenes
	if ngenes < 1 {
		return
	}
	nbases := o.Pop[0].Nbases
	for k, ind := range o.Pop {
		for i := 0; i < ngenes; i++ {
			x := math.Abs(ind.GetFloat(i))
			if k == 0 {
				o.maxabsgene[i] = x
			} else {
				o.maxabsgene[i] = max(o.maxabsgene[i], x)
			}
			for j := 0; j < nbases; j++ {
				o.fltbases[i*nbases+j][k] = ind.Floats[i*nbases+j]
			}
		}
	}
	for i := 0; i < ngenes; i++ {
		x := 1.0 + o.maxabsgene[i]
		for j := 0; j < nbases; j++ {
			o.devbases[i*nbases+j] = rnd.StatDev(o.fltbases[i*nbases+j], o.UseStdDev) / x
		}
	}
	minrho, averho, maxrho, devrho = rnd.StatBasic(o.devbases, o.UseStdDev)
	return
}

// Write writes results to buffer
func (o Island) Write(buf *bytes.Buffer, t int, json bool) {
	if json {
		return
	}
	buf.Write(o.Pop.Output(nil, o.ShowBases).Bytes())
}

// PlotOvs plots objective values versus time
func (o Island) PlotOvs(dirout, fnkey, args string, tf int, withtxt bool, numfmt string, first, last bool) {
	if first {
		plt.SetForEps(0.75, 250)
	}
	n := len(o.OVS)
	T := utl.LinSpace(0, float64(tf), n)
	plt.Plot(T, o.OVS, args)
	if withtxt {
		plt.Text(T[0], o.OVS[0], io.Sf(numfmt, o.OVS[0]), "ha='left'")
		plt.Text(T[n-1], o.OVS[n-1], io.Sf(numfmt, o.OVS[n-1]), "ha='right'")
	}
	if last {
		plt.Gll("time", "objective value", "")
		plt.SaveD(dirout, fnkey+".eps")
	}
}

// SaveReport saves report to file
func (o Island) SaveReport(dirout, fnkey string) {
	if dirout == "" {
		dirout = "/tmp/goga"
	}
	io.WriteFileD(dirout, fnkey+".rpt", &o.Report)
}

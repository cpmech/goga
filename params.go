// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"encoding/json"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// ConfParams is an auxiliary structure to hold configuration parameters for setting the GA up
type ConfParams struct {

	// initialisation
	Seed   int     // seed to initialise random numbers generator. Seed ≤ 0 means use current time
	Pll    bool    // allow running islands in parallel (go-routines)
	Nisl   int     // number of islands
	Ninds  int     // number of individuals: population size
	Nbases int     // number of bases in chromosome
	Grid   bool    // generate individuals based on grid
	Noise  float64 // apply noise when generate based on grid (if Noise > 0)
	IntOrd bool    // integer chromossome is ordered list

	// time control
	Tf    int // number of generations
	Dtout int // increment of time for output
	Dtmig int // increment of time for migration

	// regeneration
	RegTol    float64 // tolerance for ρ to activate regeneration
	RegPct    float64 // percentage of individuals to be regenerated; e.g. 0.3
	UseStdDev bool    // use standard deviation (σ) instead of average deviation in Stat

	// selection and reproduction
	Pc    float64 // probability of crossover
	Pm    float64 // probability of mutation
	Elite bool    // use elitism
	Rws   bool    // use Roulette-Wheel selection method
	Rnk   bool    // ranking
	RnkSp float64 // selective pressure for ranking

	// output
	Verbose   bool   // show messages during optimisation
	DoReport  bool   // generate report
	Json      bool   // output results as .json files; not tables
	DirOut    string // directory to save output files. "" means "/tmp/goga"
	FnKey     string // filename key for output files. "" means no output files
	DoPlot    bool   // plot results
	PltTi     int    // initial time for plot
	PltTf     int    // final time for plot
	ShowBases bool   // show also bases when printing results (if any)

	// auxiliary
	Problem  int     // problem ID
	Strategy int     // strategy for implementing constraints
	Ntrials  int     // number of trials
	Eps1     float64 // tolerance # 1; e.g. for strategy # 2 in reliability analyses

	// objective function
	OvaOor Objectives_t // compute objective value (ova) and out-of-range value (oor)

	// crossover
	CxNcuts   map[string]int         // crossover number of cuts for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxCuts    map[string][]int       // crossover specific cuts for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxProbs   map[string]float64     // crossover probabilities for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxFuncs   map[string]interface{} // crossover functions for each 'int', 'flt', 'str', 'key', 'byt', 'fun' tag
	CxIntFunc CxIntFunc_t            // crossover function
	CxFltFunc CxFltFunc_t            // crossover function
	CxStrFunc CxStrFunc_t            // crossover function
	CxKeyFunc CxKeyFunc_t            // crossover function
	CxBytFunc CxBytFunc_t            // crossover function
	CxFunFunc CxFunFunc_t            // crossover function

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

	// generation of individuals
	OrdNints   int         // ordered integer populations: number of integers
	RangeInt   [][]int     // [ngene][2] min and max integers
	RangeFlt   [][]float64 // [ngene][2] min and max float point numbers
	PoolStr    [][]string  // [ngene][nsamples] pool of words to be used in Gene.String
	PoolKey    [][]byte    // [ngene][nsamples] pool of bytes to be used in Gene.Byte
	PoolByt    [][]string  // [ngene][nsamples] pool of byte-words to be used in Gene.Bytes
	PoolFun    [][]Func_t  // [ngene][nsamples] pool of functions
	PopGenArgs interface{} // extra arguments for generation of populations

	// generation of populations
	PopIntGen PopIntGen_t // generate population of integers
	PopOrdGen PopOrdGen_t // generate population of ordered integers
	PopFltGen PopFltGen_t // generate population of float point numbers
	PopStrGen PopStrGen_t // generate population of strings
	PopKeyGen PopKeyGen_t // generate population of keys (bytes)
	PopBytGen PopBytGen_t // generate population of bytes
	PopFunGen PopFunGen_t // generate population of functions
}

// SetDefault sets default parameters
func (o *ConfParams) SetDefault() {

	// initialisation
	o.Seed = 0
	o.Pll = true
	o.Nisl = 4
	o.Ninds = 20
	o.Nbases = 10
	o.Grid = true
	o.Noise = 0.3
	o.IntOrd = false

	// time control
	o.Tf = 100
	o.Dtout = 10
	o.Dtmig = 40

	// regeneration
	o.RegTol = 1e-2
	o.RegPct = 0.3
	o.UseStdDev = false

	// selection and reproduction
	o.Pc = 0.8
	o.Pm = 0.01
	o.Elite = true
	o.Rws = false
	o.Rnk = true
	o.RnkSp = 1.2

	// output
	o.Verbose = true
	o.DoReport = false
	o.Json = false
	o.DirOut = "/tmp/goga"
	o.FnKey = ""
	o.DoPlot = false
	o.PltTi = 0
	o.PltTf = -1
	o.ShowBases = false

	// auxiliary
	o.Problem = 1
	o.Strategy = 1
	o.Ntrials = 100
	o.Eps1 = 0.1

	// number of cuts in chromossome
	o.CxNcuts = map[string]int{"int": 2, "flt": 2, "str": 2, "key": 2, "byt": 2, "fun": 2}
}

// CalcDerived calculates derived quantities
func (o *ConfParams) CalcDerived() {

	// set probabilities
	pc, pm := o.Pc, o.Pm
	o.CxProbs = map[string]float64{"int": pc, "flt": pc, "str": pc, "key": pc, "byt": pc, "fun": pc}
	o.MtProbs = map[string]float64{"int": pm, "flt": pm, "str": pm, "key": pm, "byt": pm, "fun": pm}

	// set specific crossover and mutation functions
	if o.IntOrd {
		o.CxIntFunc = IntOrdCrossover
		o.MtIntFunc = IntOrdMutation
	}
}

// NewConfParams returns a new ConfParams structure, with default values set
func NewConfParams() *ConfParams {
	var o ConfParams
	o.SetDefault()
	o.CalcDerived()
	return &o
}

// ReadConfParams reads configuration parameters from JSON file
func ReadConfParams(filenamepath string) *ConfParams {

	// new params
	var o ConfParams
	o.SetDefault()

	// read file
	b, err := io.ReadFile(filenamepath)
	if err != nil {
		chk.Panic("cannot read parameters file %q", filenamepath)
	}

	// decode
	err = json.Unmarshal(b, &o)
	if err != nil {
		chk.Panic("cannot unmarshal parameters file %q", filenamepath)
	}

	// results
	o.CalcDerived()
	return &o
}

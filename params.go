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

	// essential
	Nova   int // number of objective values
	Noor   int // number of out-of-range variables
	Nbases int // number of bases in chromosome

	// initialisation
	Seed  int  // seed to initialise random numbers generator. Seed ≤ 0 means use current time
	Pll   bool // allow running islands in parallel (go-routines)
	Nisl  int  // number of islands
	Ninds int  // number of individuals: population size

	// time control
	Tf    int // number of generations
	Dtout int // increment of time for output
	Dtmig int // increment of time for migration

	// migration and regeneration
	RegTol    float64 // tolerance for ρ to activate regeneration
	RegPct    float64 // percentage of individuals to be regenerated; e.g. 0.3
	UseStdDev bool    // use standard deviation (σ) instead of average deviation in Stat

	// operators' data
	Ops OpsData // operators' data

	// selection and reproduction
	Elite     bool    // use elitism
	Rws       bool    // use Roulette-Wheel selection method
	Rnk       bool    // ranking
	RnkSp     float64 // selective pressure for ranking
	GAtype    string  // type of GA; e.g. "std", "crowd"
	CrowdSize int     // crowd size
	ParetoPhi float64 // φ coefficient for probabilistic Pareto comparison
	CompProb  bool    // use probabilistic comparision in crowding
	DiffEvol  bool    // use differential evolution-like crossover

	// output
	Verbose   bool       // show messages during optimisation
	DoReport  bool       // generate report
	Json      bool       // output results as .json files; not tables
	DirOut    string     // directory to save output files. "" means "/tmp/goga"
	FnKey     string     // filename key for output files. "" means no output files
	DoPlot    bool       // plot results
	PltTi     int        // initial time for plot
	PltTf     int        // final time for plot
	ShowOor   bool       // show oor values when printing results (if any)
	ShowDem   bool       // show demerits when printing individuals
	ShowBases bool       // show also bases when printing results (if any)
	ShowNinds int        // number of individuals to show. use -1 to show all individuals
	PostProc  PostProc_t // function to post-process results

	// number formats. use nil for default values
	// fmts=["int","flt","str","key","byt","fun"][ngenes] print formats for each gene
	NumFmts   map[string][]string // number formats used during printing of individuals.
	NumFmtOva string              // number format for ova. use "" for default value

	// auxiliary
	Problem  int     // problem ID
	Strategy int     // strategy for implementing constraints
	Ntrials  int     // number of trials
	Eps1     float64 // tolerance # 1; e.g. for strategy # 2 in reliability analyses
	Check    bool    // run checking code before GA

	// objective function
	OvaOor Objectives_t // compute objective value (ova) and out-of-range value (oor)

	// generation of individuals
	Latin    bool        // use latin hypercube during generation
	LatinDf  int         // duplication factor when using latin hypercube
	Noise    float64     // apply noise when generating based on grid (if Noise > 0)
	NumInts  int         // number of integers for "ordered" and "binary" populations
	RangeInt [][]int     // [ngene][2] min and max integers
	RangeFlt [][]float64 // [ngene][2] min and max float point numbers
	PoolStr  [][]string  // [ngene][nsamples] pool of words to be used in Gene.String
	PoolKey  [][]byte    // [ngene][nsamples] pool of bytes to be used in Gene.Byte
	PoolByt  [][]string  // [ngene][nsamples] pool of byte-words to be used in Gene.Bytes
	PoolFun  [][]Func_t  // [ngene][nsamples] pool of functions

	// generation of populations
	PopIntGen PopIntGen_t // generate population of integers
	PopFltGen PopFltGen_t // generate population of float point numbers
	PopStrGen PopStrGen_t // generate population of strings
	PopKeyGen PopKeyGen_t // generate population of keys (bytes)
	PopBytGen PopBytGen_t // generate population of bytes
	PopFunGen PopFunGen_t // generate population of functions
}

// SetDefault sets default parameters
func (o *ConfParams) SetDefault() {

	// essential
	o.Nova = 1
	o.Noor = 0
	o.Nbases = 1

	// initialisation
	o.Seed = 0
	o.Pll = true
	o.Nisl = 4
	o.Ninds = 20

	// time control
	o.Tf = 100
	o.Dtout = 10
	o.Dtmig = 25

	// migration and regeneration
	o.RegTol = 0
	o.RegPct = 0.2
	o.UseStdDev = false

	// operators' data
	o.Ops.SetDefault()

	// selection and reproduction
	o.Elite = false
	o.Rws = false
	o.Rnk = true
	o.RnkSp = 1.2
	o.GAtype = "crowd"
	o.CrowdSize = 3
	o.ParetoPhi = 0
	o.CompProb = false
	o.DiffEvol = false

	// output
	o.Verbose = false
	o.DoReport = false
	o.Json = false
	o.DirOut = "/tmp/goga"
	o.FnKey = ""
	o.DoPlot = false
	o.PltTi = 0
	o.PltTf = -1
	o.ShowOor = true
	o.ShowDem = false
	o.ShowBases = false
	o.ShowNinds = -1

	// auxiliary
	o.Problem = 1
	o.Strategy = 1
	o.Ntrials = 1
	o.Eps1 = 0.1
	o.Check = false

	// generation of individuals
	o.Latin = true
	o.LatinDf = 5
	o.Noise = 0.1
}

// SetNbases sets number of bases and fixes corresponding operators
func (o *ConfParams) SetNbasesFixOp(nbases int) {
	o.Nbases = nbases
	o.Ops.CxFlt = FltCrossover
	o.Ops.MtFlt = FltMutation
}

// SetBlxMwicz sets BLX-α (crossover) and Michaelewicz (mutation) operators
func (o *ConfParams) SetBlxMwicz() {
	o.Ops.CxFlt = FltCrossoverBlx
	o.Ops.MtFlt = FltMutationMwicz
}

// SetIntBin sets functions to handle binary numbers [0,1]
func (o *ConfParams) SetIntBin(size int) {
	o.NumInts = size
	o.PopIntGen = PopBinGen
	o.Ops.CxInt = IntCrossover
	o.Ops.MtInt = IntBinMutation
}

// SetIntOrd sets functions to handle ordered integers
func (o *ConfParams) SetIntOrd(nstations int) {
	o.NumInts = nstations
	o.PopIntGen = PopOrdGen
	o.Ops.CxInt = IntOrdCrossover
	o.Ops.MtInt = IntOrdMutation
}

// CalcDerived calculates derived quantities
func (o *ConfParams) CalcDerived() {
	o.Ops.CalcDerived(o.Tf, o.RangeFlt)
	//if o.DiffEvol {
	//o.CrowdSize = 4
	//}
}

// NewConfParams returns a new ConfParams structure, with default values set
func NewConfParams() *ConfParams {
	var o ConfParams
	o.SetDefault()
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

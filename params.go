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
	Nova int // number of objective values
	Noor int // number of out-of-range variables

	// initialisation
	Seed  int  // seed to initialise random numbers generator. Seed ≤ 0 means use current time
	Pll   bool // allow running islands in parallel (go-routines)
	Nisl  int  // number of islands
	Ninds int  // number of individuals: population size
	Nimig int  // number of individuals that migrate

	// time control
	Tf    int // number of generations
	Dtout int // increment of time for output
	Dtmig int // increment of time for migration

	// output
	Verbose bool   // show messages during optimisation
	DoPlot  bool   // plot results
	FnKey   string // filename key for output files. "" means no output files
	DirOut  string // directory to save output files. "" means "/tmp/goga"

	// auxiliary
	Problem  int     // problem ID
	Strategy int     // strategy for implementing constraints
	Ntrials  int     // number of trials
	Eps1     float64 // tolerance # 1; e.g. for strategy # 2 in reliability analyses
	Check    bool    // run checking code before GA

	// generation of individuals
	Latin    bool        // use latin hypercube during generation
	LatinDf  int         // duplication factor when using latin hypercube
	Noise    float64     // apply noise when generating based on grid (if Noise > 0)
	NumInts  int         // number of integers for "ordered" and "binary" populations
	RangeFlt [][]float64 // [ngene][2] min and max float point numbers
	RangeInt [][]int     // [ngene][2] min and max integers

	// operators' data
	Ops OpsData // operators' data

	// callback functions
	OvaOor    Objectives_t // compute objective value (ova) and out-of-range value (oor)
	PopFltGen PopFltGen_t  // generate population of float point numbers
	PopIntGen PopIntGen_t  // generate population of integers

	// auxiliary
	derived_called bool // flags whether CalcDerived was called or not
}

// SetDefault sets default parameters
func (o *ConfParams) SetDefault() {

	// essential
	o.Nova = 1
	o.Noor = 0

	// initialisation
	o.Seed = 0
	o.Pll = true
	o.Nisl = 4
	o.Ninds = 24
	o.Nimig = 4

	// time control
	o.Tf = 200
	o.Dtout = 20
	o.Dtmig = 20

	// output
	o.Verbose = true
	o.DoPlot = false
	o.FnKey = ""
	o.DirOut = "/tmp/goga"

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

	// operators' data
	o.Ops.SetDefault()
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
	o.derived_called = true
}

// global functions ////////////////////////////////////////////////////////////////////////////////

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

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// check_input checks whether parameters are consistent or not
func (o *ConfParams) check_input() {
	if !o.derived_called {
		chk.Panic("CalcDerived must be called before simulation")
	}
	if o.Nova < 1 {
		chk.Panic("number of objective values (nova) must be greater than 0")
	}
	if len(o.RangeFlt) != len(o.Ops.Xrange) {
		chk.Panic("number of genes in RangeFlt must be equal to number of genes in Ops.Xrange => ConfParams.CalcDerived must be called. %d != %d", len(o.RangeFlt), len(o.Ops.Xrange))
	}
	if o.Nisl < 1 {
		chk.Panic("at least one island must be defined. Nisl=%d is incorrect", o.Nisl)
	}
	if o.Ninds < 2 || (o.Ninds%2 != 0) {
		chk.Panic("size of population must be even and greater than 2. Ninds = %d is invalid", o.Ninds)
	}
	if o.OvaOor == nil {
		chk.Panic("objective function (OvaOor) must be non nil")
	}
	if o.PopIntGen == nil && o.PopFltGen == nil {
		chk.Panic("at least one generator function in Params must be non nil")
	}
}

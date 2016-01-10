// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"encoding/json"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

// Parameters hold all configuration parameters
type Parameters struct {

	// sizes
	Nova int // number of objective values
	Noor int // number of out-of-range values
	Nsol int // total number of solutions
	Ncpu int // number of cpus

	// time
	Tf    int // final time
	DtExc int // delta time for exchange
	DtOut int // delta time for output

	// options
	Pll        bool    // parallel
	Seed       int     // seed for random numbers generator
	GenType    string  // generation type: "latin", "halton", "rnd"
	LatinDup   int     // Latin Hypercube duplicates number
	EpsMinProb float64 // minimum value for 'h' constraints
	Verbose    bool    // show messages
	Problem    int     // problem index
	GenAll     bool    // generate all solutions together; i.e. not within each group/CPU
	Ntrials    int     // run many trials
	BinInt     int     // flag that integers represent binary numbers if BinInt > 0; thus Nint=BinInt
	ClearFlt   bool    // clear flt if corresponding int is 0
	ExcTour    bool    // use exchange via tournament
	ExcOne     bool    // use exchange one randomly
	ConvDova0  float64 // Δova[0] to decide on convergence

	// differential evolution
	DiffEvolC        float64 // crossover probability
	DiffEvolF        float64 // vector length multiplier
	DiffEvolPm       float64 // mutation probability. use rotation otherwise
	DiffEvolUseCmult bool    // use C random multiplier
	DiffEvolUseFmult bool    // use F random multiplier

	// crossover and mutation
	DebEtam float64 // Deb's mutation parameters
	PmFlt   float64 // probability of mutation for floats
	PmInt   float64 // probability of mutation for ints
	PcInt   float64 // probability of crossover for ints

	// range
	FltMin []float64 // minimum float allowed
	FltMax []float64 // maximum float allowed
	IntMin []int     // minimum int allowed
	IntMax []int     // maximum int allowed

	// derived
	Nflt   int       // number of floats
	Nint   int       // number of integers
	DelFlt []float64 // max float range
	DelInt []int     // max int range
}

// Default sets default parameters
func (o *Parameters) Default() {

	// sizes
	o.Nova = 1
	o.Noor = 0
	o.Nsol = 40
	o.Ncpu = 4

	// time
	o.Tf = 100
	o.DtExc = -1
	o.DtOut = -1

	// options
	o.Pll = true
	o.Seed = 0
	o.GenType = "latin"
	o.LatinDup = 2
	o.EpsMinProb = 0.1
	o.Verbose = true
	o.Problem = 1
	o.GenAll = false
	o.Ntrials = 10
	o.BinInt = 0
	o.ClearFlt = false
	o.ExcTour = true
	o.ExcOne = false
	o.ConvDova0 = 0.1

	// differential evolution
	o.DiffEvolC = 1.0
	o.DiffEvolF = 1.0
	o.DiffEvolPm = 1.0
	o.DiffEvolUseCmult = true
	o.DiffEvolUseFmult = true

	// crossover and mutation
	o.DebEtam = 1
	o.PmFlt = 0.0
	o.PmInt = 0.1
	o.PcInt = 0.8
}

// Read reads configuration parameters from JSON file
func (o *Parameters) Read(filenamepath string) {
	o.Default()
	b, err := io.ReadFile(filenamepath)
	if err != nil {
		chk.Panic("cannot read parameters file %q", filenamepath)
	}
	err = json.Unmarshal(b, o)
	if err != nil {
		chk.Panic("cannot unmarshal parameters file %q", filenamepath)
	}
	return
}

// CalcDerived computes derived variables and checks consistency
func (o *Parameters) CalcDerived() {

	// check
	if o.Nova < 1 {
		chk.Panic("number of objective values (nova) must be greater than 0")
	}
	if o.Nsol < 6 {
		chk.Panic("number of solutions must greater than 6. Nsol = %d is invalid", o.Nsol)
	}
	if o.Ncpu < 2 {
		o.Ncpu = 1
		o.Pll = false
		o.DtExc = 1
	}
	if o.Ncpu > o.Nsol/2 {
		chk.Panic("number of CPU must be smaller than or equal to half the number of solutions. Ncpu=%d > Nsol/2=%d", o.Ncpu, o.Nsol/2)
	}
	if o.Tf < 1 {
		o.Tf = 1
	}
	if o.DtExc < 1 {
		o.DtExc = o.Tf / 10
	}
	if o.DtOut < 1 {
		o.DtOut = o.Tf / 5
	}

	// derived
	o.Nflt = len(o.FltMin)
	o.Nint = len(o.IntMin)
	if o.BinInt > 0 {
		o.Nint = o.BinInt
	}
	if o.Nflt == 0 && o.Nint == 0 {
		chk.Panic("either floats and ints must be set (via FltMin/Max or IntMin/Max)")
	}

	// floats
	if o.Nflt > 0 {
		chk.IntAssert(len(o.FltMax), o.Nflt)
		o.DelFlt = make([]float64, o.Nflt)
		for i := 0; i < o.Nflt; i++ {
			o.DelFlt[i] = o.FltMax[i] - o.FltMin[i]
		}
	}

	// generic ints
	if o.BinInt == 0 && o.Nint > 0 {
		chk.IntAssert(len(o.IntMax), o.Nint)
		o.DelInt = make([]int, o.Nint)
		for i := 0; i < o.Nint; i++ {
			o.DelInt[i] = o.IntMax[i] - o.IntMin[i]
		}
	}
	if o.Nint != o.Nflt {
		o.ClearFlt = false
	}

	// initialise random numbers generator
	rnd.Init(o.Seed)
}

// EnforceRange makes sure x is within given range
func (o *Parameters) EnforceRange(i int, x float64) float64 {
	if x < o.FltMin[i] {
		return o.FltMin[i]
	}
	if x > o.FltMax[i] {
		return o.FltMax[i]
	}
	return x
}

// LogParams returns a log with current parameters
func (o *Parameters) LogParams() (l string) {

	// sizes
	l += io.ArgsTable("SIZES",
		"number of objective values", "Nova", o.Nova,
		"number of out-of-range values", "Noor", o.Noor,
		"total number of solutions", "Nsol", o.Nsol,
		"number of cpus", "Ncpu", o.Ncpu,
	)

	// time
	l += "\n"
	l += io.ArgsTable("TIME",
		"final time", "Tf", o.Tf,
		"delta time for exchange", "DtExc", o.DtExc,
		"delta time for output", "DtOut", o.DtOut,
	)

	// options
	l += "\n"
	l += io.ArgsTable("OPTIONS",
		"parallel", "Pll", o.Pll,
		"seed for random numbers generator", "Seed", o.Seed,
		"generation type: 'latin', 'halton', 'rnd'", "GenType", o.GenType,
		"Latin Hypercube duplicates number", "LatinDup", o.LatinDup,
		"minimum value for 'h' constraints", "EpsMinProb", o.EpsMinProb,
		"show messages", "Verbose", o.Verbose,
		"problem index", "Problem", o.Problem,
		"generate all solutions together", "GenAll", o.GenAll,
		"run many trials", "Ntrials", o.Ntrials,
		"integers represent binary numbers", "BinInt", o.BinInt,
		"clear flt if corresponding int is 0", "ClearFlt", o.ClearFlt,
		"use exchange via tournament", "ExcTour", o.ExcTour,
		"use exchange one randomly", "ExcOne", o.ExcOne,
		"Δova[0] to decide on convergence", "ConvDova0", o.ConvDova0,
	)

	// differential evolution
	l += "\n"
	l += io.ArgsTable("DIFFERENTIAL EVOLUTION",
		"crossover probability", "DiffEvolC", o.DiffEvolC,
		"vector length multiplier", "DiffEvolF", o.DiffEvolF,
		"mutation probability", "DiffEvolPm", o.DiffEvolPm,
		"use C random multiplier", "DiffEvolUseCmult", o.DiffEvolUseCmult,
		"use F random multiplier", "DiffEvolUseFmult", o.DiffEvolUseFmult,
	)

	// crossover and mutation
	l += "\n"
	l += io.ArgsTable("CROSSOVER AND MUTATION",
		"Deb's mutation parameters", "DebEtam", o.DebEtam,
		"probability of mutation for floats", "PmFlt", o.PmFlt,
		"probability of mutation for ints", "PmInt", o.PmInt,
		"probability of crossover for ints", "PcInt", o.PcInt,
	)

	// derived
	l += "\n"
	l += io.ArgsTable("DERIVED",
		"number of floats", "Nflt", o.Nflt,
		"number of integers", "Nint", o.Nint,
	)
	return
}

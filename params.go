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
	LatinDup   int     // Latin Hypercube duplicates number
	EpsMinProb float64 // minimum value for 'h' constraints
	Verbose    bool    // show messages
	Problem    int     // problem index
	GenAll     bool    // generate all solutions together; i.e. not within each group/CPU

	// crossover and mutation
	DEpc    float64 // differential evolution pc
	DEmult  float64 // differential evolution multiplier
	DebEtac float64 // Deb's crossover parameter
	DebEtam float64 // Deb's mutation parameters
	PmFlt   float64 // probability of mutation for floats
	PmInt   float64 // probability of mutation for ints

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
	o.DtExc = o.Tf / 10
	o.DtOut = o.Tf / 5

	// options
	o.Pll = true
	o.Seed = 0
	o.LatinDup = 2
	o.EpsMinProb = 0.1
	o.Verbose = true
	o.Problem = 1
	o.GenAll = false

	// crossover and mutation
	o.DEpc = 0.5
	o.DEmult = 0.1
	o.DebEtac = 1
	o.DebEtam = 1
	o.PmFlt = 0.0
	o.PmInt = 0.1
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
	if o.Nsol < 2 {
		chk.Panic("number of solutions must greater than 2. Nsol = %d is invalid", o.Nsol)
	}
	if o.Ncpu < 2 {
		o.Ncpu = 1
		o.Pll = false
	}
	if o.Ncpu > o.Nsol/2 {
		chk.Panic("number of CPU must be smaller than or equal to half the number of solutions. Ncpu=%d > Nsol/2=%d", o.Ncpu, o.Nsol/2)
	}
	if o.Tf < 1 {
		o.Tf = 1
	}
	if o.DtExc < 1 {
		o.DtExc = 1
	}

	// derived
	o.Nflt = len(o.FltMin)
	o.Nint = len(o.IntMin)
	if o.Nflt == 0 && o.Nint == 0 {
		chk.Panic("either floats and ints must be set (via FltMin/Max or IntMin/Max)")
	}
	chk.IntAssert(len(o.FltMax), o.Nflt)
	chk.IntAssert(len(o.IntMax), o.Nint)
	o.DelFlt = make([]float64, o.Nflt)
	o.DelInt = make([]int, o.Nint)
	for i := 0; i < o.Nflt; i++ {
		o.DelFlt[i] = o.FltMax[i] - o.FltMin[i]
	}
	for i := 0; i < o.Nint; i++ {
		o.DelInt[i] = o.IntMax[i] - o.IntMin[i]
	}
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

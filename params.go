// Copyright 2015 The Goga Authors. All rights reserved.
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
	Tmax  int // final time
	DtExc int // delta time for exchange
	DtOut int // delta time for output

	// options
	DEC      float64 // C-coefficient for differential evolution
	Pll      bool    // parallel
	Seed     int     // seed for random numbers generator
	GenType  string  // generation type: "latin", "halton", "rnd"
	LatinDup int     // Latin Hypercube duplicates number
	EpsH     float64 // minimum value for 'h' constraints
	Verbose  bool    // show messages
	VerbStat bool    // show messages in Stat
	VerbTime bool    // show time messages
	GenAll   bool    // generate all solutions together; i.e. not within each group/CPU
	Nsamples int     // run many samples
	BinInt   int     // flag that integers represent binary numbers if BinInt > 0; thus Nint=BinInt
	ClearFlt bool    // clear flt if corresponding int is 0
	ExcTour  bool    // use exchange via tournament
	ExcOne   bool    // use exchange one randomly
	NormFlt  bool    // normalise float values
	UseMesh  bool    // use meshes to control points movement
	Nbry     int     // number of points along boundary / per iFlt (only if UseMesh==true)

	// crossover and mutation of integers
	IntPc       float64 // probability of crossover for ints
	IntNcuts    int     // number of cuts in crossover of ints
	IntPm       float64 // probability of mutation for ints
	IntNchanges int     // number of changes during mutation of ints

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

	// extra variables not directly related to GOGA (for convenience of having a reader already)
	Strategy int  // strategy
	PlotSet1 bool // plot set of graphs 1
	PlotSet2 bool // plot set of graphs 2
	ProbNum  int  // problem number

	// for mesh method
	NumXiXjPairs  int // number of (Xi,Xj) pairs
	NumXiXjBryPts int // number of points along the boundaries of one (Xi,Xj) plane
	NumExtraSols  int // total number of extra solutions due to all (Xi,Xj) boundaries
}

// Default sets default parameters
func (o *Parameters) Default() {

	// sizes
	o.Nova = 1
	o.Noor = 0
	o.Nsol = 40
	o.Ncpu = 4

	// time
	o.Tmax = 100
	o.DtExc = -1
	o.DtOut = -1

	// options
	o.DEC = 0.8
	o.Pll = true
	o.Seed = 0
	o.GenType = "latin"
	o.LatinDup = 2
	o.EpsH = 0.1
	o.Verbose = true
	o.VerbStat = false
	o.VerbTime = false
	o.GenAll = false
	o.Nsamples = 10
	o.BinInt = 0
	o.ClearFlt = false
	o.ExcTour = true
	o.ExcOne = true
	o.NormFlt = false
	o.UseMesh = false
	o.Nbry = 3

	// crossover and mutation of integers
	o.IntPc = 0.8
	o.IntNcuts = 1
	o.IntPm = 0.01
	o.IntNchanges = 1
}

// Read reads configuration parameters from JSON file
func (o *Parameters) Read(filenamepath string) {
	o.Default()
	b := io.ReadFile(filenamepath)
	err := json.Unmarshal(b, o)
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
	if o.Tmax < 1 {
		o.Tmax = 1
	}
	if o.DtExc < 1 {
		o.DtExc = o.Tmax / 10
	}
	if o.DtOut < 1 {
		o.DtOut = o.Tmax / 5
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

	// mesh
	if o.Nflt < 2 {
		o.UseMesh = false
	}
	if o.UseMesh {
		if o.Nbry < 2 {
			o.Nbry = 2
		}
		o.NumXiXjPairs = (o.Nflt*o.Nflt - o.Nflt) / 2
		o.NumXiXjBryPts = (o.Nbry-2)*4 + 4
		o.NumExtraSols = o.NumXiXjPairs * o.NumXiXjBryPts
		io.PfYel("NumXiXjPairs=%d NumXiXjBryPts=%d NumExtraSols=%d\n", o.NumXiXjPairs, o.NumXiXjBryPts, o.NumExtraSols)
		o.Nsol += o.NumExtraSols
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

	// number of cuts and changes in ints
	if o.Nint > 0 {
		if o.IntNcuts > o.Nint {
			o.IntNcuts = o.Nint
		}
		if o.IntNchanges > o.Nint {
			o.IntNchanges = o.Nint
		}
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

// Normalise4 normalises x ∈ [xmin,xmax] values into r ∈ [0,1]
func (o *Parameters) Normalise4(x, x0, x1, x2 []float64) (r, r0, r1, r2 []float64) {
	if o.Nflt < 1 {
		return
	}
	r, r0, r1, r2 = make([]float64, o.Nflt), make([]float64, o.Nflt), make([]float64, o.Nflt), make([]float64, o.Nflt)
	if o.NormFlt {
		for i := 0; i < o.Nflt; i++ {
			r[i] = (x[i] - o.FltMin[i]) / o.DelFlt[i]
			r0[i] = (x0[i] - o.FltMin[i]) / o.DelFlt[i]
			r1[i] = (x1[i] - o.FltMin[i]) / o.DelFlt[i]
			r2[i] = (x2[i] - o.FltMin[i]) / o.DelFlt[i]
		}
	} else {
		for i := 0; i < o.Nflt; i++ {
			r[i], r0[i], r1[i], r2[i] = x[i], x0[i], x1[i], x2[i]
		}
	}
	return
}

// DeNormalise1 de-normalises r ∈ [0,1] values into x ∈ [xmin,xmax]
//  Output: r becomes x
func (o *Parameters) DeNormalise1(r []float64) {
	if !o.NormFlt {
		return
	}
	for i := 0; i < o.Nflt; i++ {
		r[i] = o.FltMin[i] + r[i]*o.DelFlt[i]
	}
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
		"final time", "Tmax", o.Tmax,
		"delta time for exchange", "DtExc", o.DtExc,
		"delta time for output", "DtOut", o.DtOut,
	)

	// options
	l += "\n"
	l += io.ArgsTable("OPTIONS",
		"C-coefficient for differential evolution", "DEC", o.DEC,
		"parallel", "Pll", o.Pll,
		"seed for random numbers generator", "Seed", o.Seed,
		"generation type: 'latin', 'halton', 'rnd'", "GenType", o.GenType,
		"Latin Hypercube duplicates number", "LatinDup", o.LatinDup,
		"minimum value for 'h' constraints", "EpsH", o.EpsH,
		"show messages", "Verbose", o.Verbose,
		"show messages in Stat", "VerbStat", o.VerbStat,
		"show time messages", "VerbTime", o.VerbTime,
		"generate all solutions together", "GenAll", o.GenAll,
		"run many trials", "Nsamples", o.Nsamples,
		"integers represent binary numbers", "BinInt", o.BinInt,
		"clear flt if corresponding int is 0", "ClearFlt", o.ClearFlt,
		"use exchange via tournament", "ExcTour", o.ExcTour,
		"use exchange one randomly", "ExcOne", o.ExcOne,
		"normalise float values", "NormFlt", o.NormFlt,
		"use meshes to control points movement", "UseMesh", o.UseMesh,
		"number of points along boundary / per iFlt (only if UseMesh==true)", "Nbry", o.Nbry,
	)

	// crossover and mutation of integers
	l += "\n"
	l += io.ArgsTable("CROSSOVER AND MUTATION OF INTS",
		"probability of crossover for ints", "IntPc", o.IntPc,
		"number of cuts in crossover of ints", "IntNcuts", o.IntNcuts,
		"probability of mutation for ints", "IntPm", o.IntPm,
		"number of changes during mutation of ints", "IntNchanges", o.IntNchanges,
	)

	// derived
	l += "\n"
	l += io.ArgsTable("DERIVED",
		"number of floats", "Nflt", o.Nflt,
		"number of integers", "Nint", o.Nint,
		"number of (Xi,Xj) pairs", "NumXiXjPairs", o.NumXiXjPairs,
		"number of points along the boundaries of one (Xi,Xj) plane", "NumXiXjBryPts", o.NumXiXjBryPts,
		"total number of extra solutions due to all (Xi,Xj) boundaries", "NumExtraSols", o.NumExtraSols,
	)

	// extra
	l += "\n"
	l += io.ArgsTable("EXTRA",
		"strategy", "Strategy", o.Strategy,
		"plot set of graphs 1", "PlotSet1", o.PlotSet1,
		"plot set of graphs 2", "PlotSet2", o.PlotSet2,
		"problem number", "ProbNum", o.ProbNum,
	)
	return
}

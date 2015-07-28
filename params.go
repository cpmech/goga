// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"encoding/json"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// Params is an auxiliary structure to hold input parameters for setting the GA up
type Params struct {
	Nislands  int     // number of islands
	Ninds     int     // number of individuals: population size
	Nbases    int     // number of bases in chromosome
	Tf        int     // number of generations
	Dtout     int     // increment of time for output
	Dtmig     int     // increment of time for migration
	Pc        float64 // probability of crossover
	Pm        float64 // probability of mutation
	Elite     bool    // use elitism
	Rws       bool    // use Roulette-Wheel selection method
	Rnk       bool    // ranking
	RnkSP     float64 // selective pressure for ranking
	Fnkey     string  // key for reports' filenames
	UseIntRnd bool    // generate random integers instead of selecting from grid
	UseFltRnd bool    // generate random float point numbers instead of selecting from grid
}

// ReadParams reads parameters from JSON file
func ReadParams(filenamepath string) *Params {

	// new params
	var o Params

	// set default values
	o.Nislands = 1
	o.Ninds = 10
	o.Nbases = 10
	o.Tf = 100
	o.Dtout = 10
	o.Dtmig = 40
	o.Pc = 0.8
	o.Pm = 0.01
	o.Elite = true
	o.Rws = false
	o.Rnk = true
	o.RnkSP = 1.2
	o.Fnkey = ""
	o.UseIntRnd = false
	o.UseFltRnd = false

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
	return &o
}

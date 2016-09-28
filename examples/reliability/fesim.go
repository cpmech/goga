// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"math"

	"github.com/cpmech/gofem/fem"
	"github.com/cpmech/gofem/inp"
	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// global variables
var (
	FEMDATA []*FemData
)

// get_femsim_data allocates FE solvers and returns FORM data
func get_femsim_data(opt *goga.Optimiser, fnkey string) (lsft LSF_T, vars rnd.Variables) {
	FEMDATA = make([]*FemData, opt.Ncpu)
	for i := 0; i < opt.Ncpu; i++ {
		FEMDATA[i] = NewData(fnkey, i)
	}
	lsft = RunFEA
	vars = FEMDATA[0].Vars
	opt.RptName = fnkey
	opt.RptFref = []float64{FEMDATA[0].BetRef}
	return
}

// FemData structure
type FemData struct {

	// input
	VtagU  int     // tag of vertex to track deflection; 0 means all vertices
	Ukey   string  // "ux", "uy" or "uz"
	AwdU   float64 // allowed max(∥u∥);   ≤ 0 means unconstrained
	AwdM22 float64 // allowed max(∥M22∥); ≤ 0 means unconstrained
	AwdM11 float64 // allowed max(∥M11∥); ≤ 0 means unconstrained
	BetRef float64 // reference reliability index

	// derived
	Analysis *fem.Main       // fem structure
	Dom      *fem.Domain     // domain
	Sim      *inp.Simulation // simulation
	Vars     rnd.Variables   // random variables
	EqsU     []int           // equations to track deflection
}

// NewData allocates and initialises a new FEM data structure
func NewData(fnkey string, cpu int) *FemData {

	// load data
	var o FemData
	buf, err := io.ReadFile("od-" + fnkey + ".json")
	if err != nil {
		chk.Panic("cannot load data:\n%v", err)
	}
	err = json.Unmarshal(buf, &o)
	if err != nil {
		chk.Panic("cannot unmarshal data:\n%v", err)
	}

	// load FEM data
	o.Analysis = fem.NewMain(fnkey+".sim", io.Sf("cpu%d", cpu), false, false, false, false, false, cpu)
	o.Dom = o.Analysis.Domains[0]
	o.Sim = o.Dom.Sim
	o.Vars = o.Dom.Sim.AdjRandom

	// backup dependent variables
	for _, prm := range o.Sim.AdjDependent {
		prm.S = prm.V // copy V into S
	}

	// set stage
	err = o.Analysis.SetStage(0)
	if err != nil {
		chk.Panic("cannot set stage:\n%v", err)
	}

	// equations to track U
	verts := o.Dom.Msh.Verts
	if o.VtagU < 0 {
		verts = o.Dom.Msh.VertTag2verts[o.VtagU]
		if len(verts) < 1 {
			chk.Panic("cannot find vertices with tag = %d\n", o.VtagU)
		}
	}
	o.EqsU = make([]int, len(verts))
	for i, vert := range verts {
		eq := o.Dom.Vid2node[vert.Id].GetEq(o.Ukey)
		if eq < 0 {
			chk.Panic("cannot find equation corresponding to vertex id=%d and ukey=%q", vert.Id, o.Ukey)
		}
		o.EqsU[i] = eq
	}
	return &o
}

// RunFEA runs FE analysis.
func RunFEA(x []float64, cpu int) (lsf, failed float64) {

	// FemData
	o := FEMDATA[cpu]

	// check for NaNs
	defer func() {
		if math.IsNaN(failed) || math.IsNaN(lsf) {
			io.PfRed("x = %+#v\n", x)
			chk.Panic("NaN: failed=%v lsf=%v\n", failed, lsf)
		}
	}()

	// adjust parameters
	for i, v := range o.Vars {
		v.Prm.Set(x[i])
	}
	for _, prm := range o.Sim.AdjDependent {
		if prm.N == "I22" {
			prm.Set(prm.S * prm.Other.V * prm.Other.V)
		}
	}
	o.Dom.RecomputeKM()

	// run
	err := o.Analysis.SolveOneStage(0, true)
	if err != nil {
		failed = 1
		return
	}

	// displacement based limit-state-function
	if o.AwdU > 0 {
		δ := o.Dom.Sol.Y[o.EqsU[0]]
		for i := 1; i < len(o.EqsU); i++ {
			δ = utl.Max(δ, o.Dom.Sol.Y[o.EqsU[i]])
		}
		lsf = o.AwdU - δ
	}
	return
}

func (o FemData) Info() (l string) {
	l += io.Sf("VtagU  = %v\n", o.VtagU)
	l += io.Sf("AwdU   = %v\n", o.AwdU)
	l += io.Sf("AwdM22 = %v\n", o.AwdM22)
	l += io.Sf("AwdM11 = %v\n", o.AwdM11)
	return
}

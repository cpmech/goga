// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"math"

	"github.com/cpmech/gofem/ele/solid"
	"github.com/cpmech/gofem/fem"
	"github.com/cpmech/gofem/inp"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Tag2val implements a pair of tag and value
type Tag2val struct {
	Tag int
	Val float64
}

// OptData optimisation data
type OptData struct {
	Amin     float64     // range: minimum cross sectional area
	Amax     float64     // range: maximum cross sectional area
	Uawd     float64     // allowed vertical deflection
	VtagU    int         // tag of vertex to track deflection; 0 means all vertices
	Groups   bool        // use groups == cell tags; otherwise all cells are considered
	BinInt   bool        // use integers as binary values
	Mobility bool        // use mobility analysis
	ReqVtags []int       // required vertex tags
	ScompAwd []Tag2val   // allowed compressive stress (positive constant)
	StensAwd []Tag2val   // allowed tensile stress (positive constant)
	aEps     float64     // [derived] area-limit == Amin when not using BinInt
	t2i      map[int]int // [derived] maps tag to index in ScompAwd and StensAwd
}

// CalcErrSig computes error: |sig| - sig_allowed
func (o *OptData) CalcErrSig(ctag int, sig float64) float64 {
	i := o.t2i[ctag]
	if len(o.StensAwd) == len(o.ScompAwd) {
		if sig > 0 { // tensile
			return sig - o.StensAwd[i].Val
		}
	}
	return math.Abs(sig) - o.ScompAwd[i].Val
}

// CalcDerived computes derived variables
func (o *OptData) CalcDerived() {
	o.t2i = make(map[int]int)
	for i, pair := range o.ScompAwd {
		o.t2i[pair.Tag] = i
	}
	if !o.BinInt {
		o.aEps = o.Amin
		o.Amin = -o.Amax
	}
}

// FemData holds derived data for FE analyses, one per CPU
type FemData struct {
	Opt       OptData     // optimisation data
	Analysis  *fem.Main   // fem structure
	Reg       *inp.Region // region
	Dom       *fem.Domain // domain
	ReqVids   []int       // required vertex ids
	Cid2xid   []int       // element id => areas array id (group)
	Ncells    int         // len(Cells)
	Nareas    int         // len(Areas)
	VidU      int         // vertex id corresponding to VtagU; -1 means all vertices
	MaxWeight float64     // max weight will all rods active
}

// NewData allocates and initialises a new FEM data structure
func NewData(filename, fnkey string, cpu int) *FemData {

	// load opt data
	var o FemData
	buf, err := io.ReadFile("od-" + fnkey + ".json")
	if err != nil {
		chk.Panic("cannot load opt data:\n%v", err)
	}
	err = json.Unmarshal(buf, &o.Opt)
	if err != nil {
		chk.Panic("cannot unmarshal opt data:\n%v", err)
	}
	o.Opt.CalcDerived()

	// load FEM data
	o.Analysis = fem.NewMain(filename, io.Sf("cpu%d", cpu), false, false, false, false, false, cpu)
	o.Reg = o.Analysis.Sim.Regions[0]
	o.Dom = o.Analysis.Domains[0]

	// required vertices
	for _, vtag := range o.Opt.ReqVtags {
		verts := o.Dom.Msh.VertTag2verts[vtag]
		for _, vert := range verts {
			o.ReqVids = append(o.ReqVids, vert.Id)
		}
	}

	// groups
	o.Ncells = len(o.Dom.Msh.Cells)
	o.Cid2xid = make([]int, o.Ncells)
	if o.Opt.Groups {
		o.Nareas = len(o.Dom.Msh.CellTag2cells)
		for tag, cells := range o.Dom.Msh.CellTag2cells {
			xid := -tag - 1
			for _, cell := range cells {
				o.Cid2xid[cell.Id] = xid
			}
		}
	} else { // all cells
		o.Nareas = o.Ncells
		for i, _ := range o.Dom.Msh.Cells {
			o.Cid2xid[i] = i
		}
	}

	// vertex to track deflection
	o.VidU = -1 // allvertices
	if o.Opt.VtagU < 0 {
		verts := o.Dom.Msh.VertTag2verts[o.Opt.VtagU]
		if len(verts) != 1 {
			chk.Panic("cannot set vertex to track deflection: verts = %v", verts)
		}
		o.VidU = verts[0].Id
	}

	// compute max weight
	o.Analysis.SetStage(0)
	for _, elem := range o.Dom.Elems {
		ele := elem.(*solid.ElastRod)
		o.MaxWeight += ele.Mdl.Rho * ele.Mdl.A * ele.L
	}
	return &o
}

// RunFEM runs FE analysis.
func (o *FemData) RunFEM(Enabled []int, Areas []float64, draw int, debug bool) (mobility, failed, weight, umax, smax, errU, errS float64) {

	// check for NaNs
	defer func() {
		if math.IsNaN(mobility) || math.IsNaN(failed) || math.IsNaN(weight) || math.IsNaN(umax) || math.IsNaN(smax) || math.IsNaN(errU) || math.IsNaN(errS) {
			io.PfRed("enabled := %+#v\n", Enabled)
			io.PfRed("areas := %+#v\n", Areas)
			chk.Panic("NaN: mobility=%v failed=%v weight=%v umax=%v smax=%v errU=%v errS=%v\n", mobility, failed, weight, umax, smax, errU, errS)
		}
	}()

	// set connectivity
	if o.Opt.BinInt {
		for cid, ena := range Enabled {
			o.Dom.Msh.Cells[cid].Disabled = true
			if ena == 1 {
				o.Dom.Msh.Cells[cid].Disabled = false
			}
		}
	} else {
		for _, cell := range o.Dom.Msh.Cells {
			cid := cell.Id
			xid := o.Cid2xid[cid]
			o.Dom.Msh.Cells[cid].Disabled = true
			if Areas[xid] >= o.Opt.aEps {
				o.Dom.Msh.Cells[cid].Disabled = false
			}
		}
	}

	// set stage
	o.Analysis.SetStage(0)

	// check for required vertices
	nnod := len(o.Dom.Nodes)
	for _, vid := range o.ReqVids {
		if o.Dom.Vid2node[vid] == nil {
			//io.Pforan("required vertex (%d) missing\n", vid)
			mobility, failed, errU, errS = float64(1+2*nnod), 1, 1, 1
			return
		}
	}

	// compute mobility
	if o.Opt.Mobility {
		m := len(o.Dom.Elems)
		d := len(o.Dom.EssenBcs.Bcs)
		M := 2*nnod - m - d
		if M > 0 {
			//io.Pforan("full mobility: M=%v\n", M)
			mobility, failed, errU, errS = float64(M), 1, 1, 1
			return
		}
	}

	// set elements' cross-sectional areas and compute weight
	for _, elem := range o.Dom.Elems {
		ele := elem.(*solid.ElastRod)
		cid := ele.Cell.Id
		xid := o.Cid2xid[cid]
		ele.Mdl.A = Areas[xid]
		ele.Recompute(false)
		weight += ele.Mdl.Rho * ele.Mdl.A * ele.L
	}

	// run FE analysis
	err := o.Analysis.SolveOneStage(0, true)
	if err != nil {
		//io.Pforan("analysis failed\n")
		mobility, failed, errU, errS = 0, 1, 1, 1
		return
	}

	// find maximum deflection
	// Note that sometimes a mechanism can happen that makes the tip to move upwards
	if o.VidU >= 0 {
		eq := o.Dom.Vid2node[o.VidU].GetEq("uy")
		uy := o.Dom.Sol.Y[eq]
		umax = math.Abs(uy)
	} else {
		for _, nod := range o.Dom.Nodes {
			eq := nod.GetEq("uy")
			uy := o.Dom.Sol.Y[eq]
			umax = utl.Max(umax, math.Abs(uy))
		}
	}
	errU = fun.Ramp(umax - o.Opt.Uawd)

	// find max stress
	for _, elem := range o.Dom.Elems {
		ele := elem.(*solid.ElastRod)
		tag := ele.Cell.Tag
		sig := ele.CalcSig(o.Dom.Sol)
		smax = utl.Max(smax, math.Abs(sig))
		errS = utl.Max(errS, o.Opt.CalcErrSig(tag, sig))
	}
	errS = fun.Ramp(errS)

	// draw
	if draw > 0 {
		args := make(map[int]*plt.A)
		for _, elem := range o.Dom.Elems {
			ele := elem.(*solid.ElastRod)
			cid := ele.Cell.Id
			args[cid] = &plt.A{C: "#004cc9", Lw: 0.3 + ele.Mdl.A/15.0}
		}
		o.Dom.Msh.Draw2d(true, false, false, nil, args, nil)
	}

	// debug
	if false && (umax > 0 && errS < 1e-10) {
		io.PfYel("enabled := %+#v\n", Enabled)
		io.Pf("areas := %+#v\n", Areas)
		io.Pf("weight  = %v\n", weight)
		io.Pf("umax    = %v\n", umax)
		io.Pf("smax    = %v\n", smax)
		io.Pf("errU    = %v\n", errU)
		io.Pf("errS    = %v\n", errS)

		// post-processing
		msh := o.Dom.Msh
		vid := msh.VertTag2verts[-4][0].Id
		nod := o.Dom.Vid2node[vid]
		eqy := nod.GetEq("uy")
		uy := o.Dom.Sol.Y[eqy]
		io.Pfblue2("%2d : uy = %g\n", vid, uy)
	}
	return
}

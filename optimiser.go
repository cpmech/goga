// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// constants
const (
	INF = 1e+30 // infinite distance
)

// Generator_t defines callback function to generate trial solutions
type Generator_t func(sols []*Solution, prms *Parameters)

// ObjFunc_t defines the objective fluction
type ObjFunc_t func(sol *Solution, grp int)

// MinProb_t defines objective functon for specialised minimisation problem
type MinProb_t func(f, g, h, x []float64, ξ []int, grp int)

// CxFlt_t defines crossover function for floats
type CxFlt_t func(a, b, A, B, C, D []float64, prms *Parameters)

// CxInt_t defines crossover function for ints
type CxInt_t func(a, b, A, B, C, D []int, prms *Parameters)

// MtFlt_t defines mutation function for floats
type MtFlt_t func(a []float64, prms *Parameters)

// MtInt_t defines mutation function for ints
type MtInt_t func(a []int, prms *Parameters)

// Optimiser solves optimisation problems:
//  Solve:
//   min  {Ova[0](x), Ova[1](x), ...} objective values
//    x
//   s.t. Oor[0](x) = 0
//        Oor[1](x) = 0  out-of-range values
//
// A specialised version is also available
//  Solve:
//    min  {f0(x), f1(x), f2(x), ...}  nf functions
//     x   g0(x) ≥ 0
//         g1(x) ≥ 0  ng inequalities
//    s.t. h0(x) = 0
//         h1(x) = 0  nh equalities
//
//  x = xFlt  or  x = xInt   or  x = {xFlt, Xint}
//
type Optimiser struct {
	Group                  // Group holds solutions and parameters
	Groups     []*Group    // view to Solutions in Group
	ObjFunc    ObjFunc_t   // objective function
	MinProb    MinProb_t   // minimisation problem function
	Nf, Ng, Nh int         // number of f, g, h functions
	F, G, H    [][]float64 // temporary
	CxFlt      CxFlt_t     // crossover function for floats
	CxInt      CxInt_t     // crossover function for ints
	MtFlt      MtFlt_t     // mutation function for floats
	MtInt      MtInt_t     // mutation function for ints
}

// Initialises continues initialisation by generating individuals
//  Optional: fcn XOR obj, nf, ng, nh
func (o *Optimiser) Init(gen Generator_t, obj ObjFunc_t, fcn MinProb_t, nf, ng, nh int) {

	// generic or minimisation problem
	if obj != nil {
		o.ObjFunc = obj
	} else {
		if fcn == nil {
			chk.Panic("either ObjFunc or MinProb must be provided")
		}
		o.Nf, o.Ng, o.Nh, o.MinProb = nf, ng, nh, fcn
		o.ObjFunc = func(sol *Solution, grp int) {
			o.MinProb(o.F[grp], o.G[grp], o.H[grp], sol.Flt, sol.Int, grp)
			for i, f := range o.F[grp] {
				sol.Ova[i] = f
			}
			for i, g := range o.G[grp] {
				sol.Oor[i] = utl.GtePenalty(g, 0.0, 1) // g[i] ≥ 0
			}
			for i, h := range o.H[grp] {
				h = math.Abs(h)
				sol.Ova[0] += h
				sol.Oor[o.Ng+i] = utl.GtePenalty(o.EpsMinProb, h, 1) // ϵ ≥ |h[i]|
			}
		}
		o.F = la.MatAlloc(o.Ngrp, o.Nf)
		o.G = la.MatAlloc(o.Ngrp, o.Ng)
		o.H = la.MatAlloc(o.Ngrp, o.Nh)
		o.Nova = o.Nf
		o.Noor = o.Ng + o.Nh
	}

	// calc derived parameters
	o.CalcDerived()

	// trial solutions
	o.InitGroup(0, nil)
	o.Groups = make([]*Group, o.Ngrp)
	for i := 0; i < o.Ngrp; i++ {
		start := i * o.Nsol
		o.Groups[i] = new(Group)
		o.Groups[i].Parameters = o.Parameters
		o.Groups[i].InitGroup(i, o.Solutions[start:start+o.Nsol])
	}
	gen(o.Solutions, &o.Parameters)

	// crossover and mutation functions
	if o.CxFlt == nil {
		o.CxFlt = CxFltDE
	}
	if o.MtFlt == nil {
		o.MtFlt = MtFltDeb
	}
}

// Solve solves optimisation problem
func (o *Optimiser) Solve() {

	time := 1
	tmig := o.DtMig
	done := make(chan int, o.Ngrp)
	for time < o.Tf {

		// evolve up to migration time
		if o.Pll {
			for i := 0; i < o.Ngrp; i++ {
				go func(grp *Group) {
					for t := time; t < tmig; t++ {
						o.evolve(grp)
						o.print_time(t, grp.Id)
					}
					done <- 1
				}(o.Groups[i])
			}
			for i := 0; i < o.Ngrp; i++ {
				<-done
			}
		}

		// migration
		time = tmig
		tmig += o.DtMig
		if o.Verbose {
			io.Pfyel(" %d", time)
		}
		o.evolve(&o.Group)
	}

	if o.Verbose {
		io.Pf("\n")
	}
}

func (o *Optimiser) evolve(grp *Group) {

	// copy current population
	for i, sol := range grp.Solutions {
		sol.CopyInto(grp.Competitors[i])
	}

	// compute random pairs
	rnd.IntGetGroups(grp.Pairs, grp.Indices)

	// create new solutions
	idx := len(grp.Solutions)
	var a, b, A, B, C, D *Solution
	for k, pair := range grp.Pairs {
		l := (k + 1) % len(grp.Pairs)
		A = grp.Competitors[pair[0]]
		B = grp.Competitors[pair[1]]
		C = grp.Competitors[grp.Pairs[l][0]]
		D = grp.Competitors[grp.Pairs[l][1]]
		a = grp.Competitors[idx]
		b = grp.Competitors[idx+1]
		idx += 2
		o.crossover(a, b, A, B, C, D)
		o.mutation(a)
		o.mutation(b)
		o.ObjFunc(a, grp.Id)
		o.ObjFunc(b, grp.Id)
		grp.Nfeval += 2
	}

	// metrics
	grp.Metrics(true)

	// tournaments
	idx = len(grp.Solutions)
	for _, pair := range grp.Pairs {
		A = grp.Competitors[pair[0]]
		B = grp.Competitors[pair[1]]
		a = grp.Competitors[idx]
		b = grp.Competitors[idx+1]
		idx += 2
		dAa := A.Distance(a, grp.Fmin, grp.Fmax, grp.Imin, grp.Imax)
		dAb := A.Distance(b, grp.Fmin, grp.Fmax, grp.Imin, grp.Imax)
		dBa := B.Distance(a, grp.Fmin, grp.Fmax, grp.Imin, grp.Imax)
		dBb := B.Distance(b, grp.Fmin, grp.Fmax, grp.Imin, grp.Imax)
		if dAa+dBb < dAb+dBa {
			o.tournament(grp.Solutions[pair[0]], A, a)
			o.tournament(grp.Solutions[pair[1]], B, b)
		} else {
			o.tournament(grp.Solutions[pair[0]], A, b)
			o.tournament(grp.Solutions[pair[1]], B, a)
		}
	}
}

func (o *Optimiser) tournament(placeholder, p, q *Solution) {
	if p.Fight(q) {
		p.CopyInto(placeholder)
	} else {
		q.CopyInto(placeholder)
	}
}

func (o *Optimiser) crossover(a, b, A, B, C, D *Solution) {
	if o.Nflt > 0 {
		o.CxFlt(a.Flt, b.Flt, A.Flt, B.Flt, C.Flt, D.Flt, &o.Parameters)
	}
	if o.Nint > 0 {
		o.CxInt(a.Int, b.Int, A.Int, B.Int, C.Int, D.Int, &o.Parameters)
	}
}

func (o *Optimiser) mutation(a *Solution) {
	if o.Nflt > 0 && o.PmFlt > 0 {
		o.MtFlt(a.Flt, &o.Parameters)
	}
	if o.Nint > 0 && o.PmInt > 0 {
		o.MtInt(a.Int, &o.Parameters)
	}
}

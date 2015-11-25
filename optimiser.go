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

	// input
	Parameters           // input parameters
	ObjFunc    ObjFunc_t // objective function
	MinProb    MinProb_t // minimisation problem function
	CxFlt      CxFlt_t   // crossover function for floats
	CxInt      CxInt_t   // crossover function for ints
	MtFlt      MtFlt_t   // mutation function for floats
	MtInt      MtInt_t   // mutation function for ints

	// essential
	Competitors []*Solution // current and future solutions
	Solutions   []*Solution // current solutions => view to Competitors
	Indices     []int       // indices of individuals in Solutions
	Pairs       [][]int     // randomly selected pairs from Indices

	// metrics
	Omin   []float64     // current min ova
	Omax   []float64     // current max ova
	Fmin   []float64     // current min float
	Fmax   []float64     // current max float
	Imin   []int         // current min int
	Imax   []int         // current max int
	Fsizes []int         // front sizes
	Fronts [][]*Solution // non-dominated fronts

	// auxiliary
	Nf, Ng, Nh int         // number of f, g, h functions
	F, G, H    [][]float64 // temporary

	// stat
	Nfeval int // number of function evaluations
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

	// crossover and mutation functions
	if o.CxFlt == nil {
		o.CxFlt = CxFltDE
	}
	if o.MtFlt == nil {
		o.MtFlt = MtFltDeb
	}

	// essential
	npairs := o.NsolTot / 2
	ncomps := o.NsolTot * 2
	o.Competitors = NewSolutions(ncomps, &o.Parameters)
	o.Solutions = o.Competitors[:o.NsolTot]
	o.Indices = utl.IntRange(o.NsolTot)
	o.Pairs = utl.IntsAlloc(npairs, 2)

	// generate trial solutions
	grp := 0
	gen(o.Solutions, &o.Parameters)
	for i := 0; i < o.NsolTot; i++ {
		o.ObjFunc(o.Solutions[i], grp)
		o.Nfeval++
	}

	// metrics
	o.Omin = make([]float64, o.Nova)
	o.Omax = make([]float64, o.Nova)
	o.Fmin = make([]float64, o.Nflt)
	o.Fmax = make([]float64, o.Nflt)
	o.Imin = make([]int, o.Nint)
	o.Imax = make([]int, o.Nint)
	o.Fsizes = make([]int, ncomps)
	o.Fronts = make([][]*Solution, ncomps)
	for i := 0; i < ncomps; i++ {
		o.Fronts[i] = make([]*Solution, ncomps)
	}
}

// GetSolutionsCopy returns a copy of Solutions
func (o *Optimiser) GetSolutionsCopy() (res []*Solution) {
	res = NewSolutions(len(o.Solutions), &o.Parameters)
	for i, sol := range o.Solutions {
		sol.CopyInto(res[i])
	}
	return
}

// Solve solves optimisation problem
func (o *Optimiser) Solve() {
	for time := 1; time <= o.Tf; time++ {
		o.evolve()
		if o.Verbose {
			io.PfWhite("time = %10d\r", time)
		}
	}
	if o.Verbose {
		io.Pf("\n")
	}
}

func (o *Optimiser) evolve() {

	// compute random pairs
	rnd.IntGetGroups(o.Pairs, o.Indices)

	// create new solutions
	grp := 0
	idx := o.NsolTot
	var a, b, A, B, C, D *Solution
	for k, pair := range o.Pairs {
		l := (k + 1) % len(o.Pairs)
		A = o.Competitors[pair[0]]
		B = o.Competitors[pair[1]]
		C = o.Competitors[o.Pairs[l][0]]
		D = o.Competitors[o.Pairs[l][1]]
		a = o.Competitors[idx]
		b = o.Competitors[idx+1]
		idx += 2
		o.crossover(a, b, A, B, C, D)
		o.mutation(a)
		o.mutation(b)
		o.ObjFunc(a, grp)
		o.ObjFunc(b, grp)
		o.Nfeval += 2
	}

	// metrics
	o.metrics(o.Competitors)

	// tournaments
	idx = o.NsolTot
	for _, pair := range o.Pairs {
		A = o.Competitors[pair[0]]
		B = o.Competitors[pair[1]]
		a = o.Competitors[idx]
		b = o.Competitors[idx+1]
		idx += 2
		dAa := A.Distance(a, o.Fmin, o.Fmax, o.Imin, o.Imax)
		dAb := A.Distance(b, o.Fmin, o.Fmax, o.Imin, o.Imax)
		dBa := B.Distance(a, o.Fmin, o.Fmax, o.Imin, o.Imax)
		dBb := B.Distance(b, o.Fmin, o.Fmax, o.Imin, o.Imax)
		if dAa+dBb < dAb+dBa {
			if a.Fight(A) {
				a.CopyInto(A)
			}
			if b.Fight(B) {
				b.CopyInto(B)
			}
		} else {
			if b.Fight(A) {
				b.CopyInto(A)
			}
			if a.Fight(B) {
				a.CopyInto(B)
			}
		}
	}
}

func (o *Optimiser) metrics(sols []*Solution) (nfronts int) {

	// reset counters and find limits
	fz := o.Fsizes
	nsol := len(sols)
	for i, sol := range sols {

		// reset values
		sol.Nwins = 0
		sol.Nlosses = 0
		sol.FrontId = 0
		sol.DistCrowd = 0
		sol.DistNeigh = INF
		fz[i] = 0

		// ovas range
		for j := 0; j < o.Nova; j++ {
			x := sol.Ova[j]
			if math.IsNaN(x) {
				chk.Panic("NaN found in objective value array\n\txFlt = %v\n\txInt = %v\n\tova = %v\n\toor = %v", sol.Flt, sol.Int, sol.Ova, sol.Oor)
			}
			if i == 0 {
				o.Omin[j] = x
				o.Omax[j] = x
			} else {
				o.Omin[j] = utl.Min(o.Omin[j], x)
				o.Omax[j] = utl.Max(o.Omax[j], x)
			}
		}

		// floats range
		for j := 0; j < o.Nflt; j++ {
			x := sol.Flt[j]
			if i == 0 {
				o.Fmin[j] = x
				o.Fmax[j] = x
			} else {
				o.Fmin[j] = utl.Min(o.Fmin[j], x)
				o.Fmax[j] = utl.Max(o.Fmax[j], x)
			}
		}

		// ints range
		for j := 0; j < o.Nint; j++ {
			x := sol.Int[j]
			if i == 0 {
				o.Imin[j] = x
				o.Imax[j] = x
			} else {
				o.Imin[j] = utl.Imin(o.Imin[j], x)
				o.Imax[j] = utl.Imax(o.Imax[j], x)
			}
		}
	}

	// compute neighbour distances and dominance data
	for i := 0; i < nsol; i++ {
		A := sols[i]
		for j := i + 1; j < nsol; j++ {
			B := sols[j]
			dist := A.Distance(B, o.Fmin, o.Fmax, o.Imin, o.Imax)
			A.DistNeigh = utl.Min(A.DistNeigh, dist)
			B.DistNeigh = utl.Min(B.DistNeigh, dist)
			A_dom, B_dom := A.Compare(B)
			if A_dom {
				A.WinOver[A.Nwins] = B // i dominates j
				A.Nwins++              // i has another dominated item
				B.Nlosses++            // j is being dominated by i
			}
			if B_dom {
				B.WinOver[B.Nwins] = A // j dominates i
				B.Nwins++              // j has another dominated item
				A.Nlosses++            // i is being dominated by j
			}
		}
	}

	// first front
	for _, sol := range sols {
		if sol.Nlosses == 0 {
			o.Fronts[0][fz[0]] = sol
			fz[0]++
		}
	}

	// next fronts
	for r, front := range o.Fronts {
		if fz[r] == 0 {
			break
		}
		nfronts++
		for s := 0; s < fz[r]; s++ {
			A := front[s]
			for k := 0; k < A.Nwins; k++ {
				B := A.WinOver[k]
				B.Nlosses--
				if B.Nlosses == 0 { // B belongs to next front
					B.FrontId = r + 1
					o.Fronts[r+1][fz[r+1]] = B
					fz[r+1]++
				}
			}
		}
	}

	// crowd distances
	for r := 0; r < nfronts; r++ {
		l, m, n := fz[r], fz[r]-1, fz[r]-2
		if l == 1 {
			o.Fronts[r][0].DistCrowd = -1
			continue
		}
		F := o.Fronts[r][:l]
		for j := 0; j < o.Nova; j++ {
			SortByOva(F, j)
			δ := o.Omax[j] - o.Omin[j] + 1e-15
			if false {
				F[0].DistCrowd += math.Pow((F[1].Ova[j]-F[0].Ova[j])/δ, 2.0)
				F[m].DistCrowd += math.Pow((F[m].Ova[j]-F[n].Ova[j])/δ, 2.0)
			} else {
				F[0].DistCrowd = INF
				F[m].DistCrowd = INF
			}
			for i := 1; i < m; i++ {
				F[i].DistCrowd += ((F[i].Ova[j] - F[i-1].Ova[j]) / δ) * ((F[i+1].Ova[j] - F[i].Ova[j]) / δ)
			}
		}
	}
	return
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

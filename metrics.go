// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

// Metrics holds metric data such as non-dominated Pareto fronts
type Metrics struct {
	prms   *Parameters   // parameters
	Omin   []float64     // current min ova
	Omax   []float64     // current max ova
	Fmin   []float64     // current min float
	Fmax   []float64     // current max float
	Imin   []int         // current min int
	Imax   []int         // current max int
	Fsizes []int         // front sizes
	Fronts [][]*Solution // non-dominated fronts
}

// Init initialises Metrics
func (o *Metrics) Init(nsol int, prms *Parameters) {
	o.prms = prms
	o.Omin = make([]float64, prms.Nova)
	o.Omax = make([]float64, prms.Nova)
	o.Fmin = make([]float64, prms.Nflt)
	o.Fmax = make([]float64, prms.Nflt)
	o.Imin = make([]int, prms.Nint)
	o.Imax = make([]int, prms.Nint)
	o.Fsizes = make([]int, nsol)
	o.Fronts = make([][]*Solution, nsol)
	for i := 0; i < nsol; i++ {
		o.Fronts[i] = make([]*Solution, nsol)
	}
}

// Compute computes limits, find non-dominated Pareto fronts, and compute crowd distances
func (o *Metrics) Compute(sols []*Solution) (nfronts int) {

	// reset variables and find limits
	z := o.Fsizes
	nsol := len(sols)
	for i, sol := range sols {

		// reset values
		sol.Nwins = 0
		sol.Nlosses = 0
		sol.FrontId = 0
		sol.DistCrowd = 0
		sol.DistNeigh = INF
		z[i] = 0

		// check oors
		for j := 0; j < o.prms.Noor; j++ {
			if math.IsNaN(sol.Oor[j]) {
				chk.Panic("NaN found in out-of-range value array\n\txFlt = %v\n\txInt = %v\n\tova = %v\n\toor = %v", sol.Flt, sol.Int, sol.Ova, sol.Oor)
			}
		}

		// ovas range
		for j := 0; j < o.prms.Nova; j++ {
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
		for j := 0; j < o.prms.Nflt; j++ {
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
		for j := 0; j < o.prms.Nint; j++ {
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

	// compute neighbour distance
	for i := 0; i < nsol; i++ {
		A := sols[i]
		for j := i + 1; j < nsol; j++ {
			B := sols[j]
			o.closest(A, B)
		}
	}

	// skip if single-objective problem
	if o.prms.Nova < 2 {
		return
	}

	// compute wins/losses data
	for i := 0; i < nsol; i++ {
		A := sols[i]
		for j := i + 1; j < nsol; j++ {
			B := sols[j]
			A_win, B_win := A.Compare(B)
			if A_win {
				A.WinOver[A.Nwins] = B
				A.Nwins++
				B.Nlosses++
			}
			if B_win {
				B.WinOver[B.Nwins] = A
				B.Nwins++
				A.Nlosses++
			}
		}
	}

	// first front
	for _, sol := range sols {
		if sol.Nlosses == 0 {
			o.Fronts[0][z[0]] = sol
			z[0]++
		}
	}

	// next fronts
	for r, front := range o.Fronts {
		if z[r] == 0 {
			break
		}
		s := r + 1
		nfronts++
		for i := 0; i < z[r]; i++ {
			A := front[i]
			for j := 0; j < A.Nwins; j++ {
				B := A.WinOver[j]
				B.Nlosses--
				if B.Nlosses == 0 { // B belongs to next front
					B.FrontId = s
					o.Fronts[s][z[s]] = B
					z[s]++
				}
			}
		}
	}

	// crowd distances
	for r := 0; r < nfronts; r++ {
		l, m := z[r], z[r]-1
		if l == 1 {
			o.Fronts[r][0].DistCrowd = -1
			continue
		}
		F := o.Fronts[r][:l]
		for j := 0; j < o.prms.Nova; j++ {
			sortByOva(F, j)
			δ := o.Omax[j] - o.Omin[j] + 1e-15
			F[0].DistCrowd = INF
			F[m].DistCrowd = INF
			for i := 1; i < m; i++ {
				F[i].DistCrowd += ((F[i].Ova[j] - F[i-1].Ova[j]) / δ) * ((F[i+1].Ova[j] - F[i].Ova[j]) / δ)
			}
		}
	}
	return
}

// closest computes distance and set closest neighbours
func (o *Metrics) closest(A, B *Solution) {
	dist := A.Distance(B, o.Fmin, o.Fmax, o.Imin, o.Imax)
	if dist < A.DistNeigh {
		A.DistNeigh = dist
		A.Closest = B
	}
	if dist < B.DistNeigh {
		B.DistNeigh = dist
		B.Closest = A
	}
}

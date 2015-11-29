// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
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

// ComputeLimitsAndNeighDist computes ova,flt,int limits and neighbour distances
func (o *Metrics) ComputeLimitsAndNeighDist(sols []*Solution) {

	// find limits
	nsol := len(sols)
	for i, sol := range sols {

		// reset values
		sol.Repeated = false
		sol.DistNeigh = INF

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

	// compute neighbour distances
	dmax := o.calcdmax()
	for i := 0; i < nsol; i++ {
		A := sols[i]
		for j := i + 1; j < nsol; j++ {
			B := sols[j]
			o.closest(A, B, dmax, nsol)
		}
	}
}

// ComputeFrontsAndCrowdDist computes non-dominated Pareto fronts and crowd distances
func (o *Metrics) ComputeFrontsAndCrowdDist(sols []*Solution) (nfronts int) {

	chk.Panic("stop")

	// reset counters
	fz := o.Fsizes
	nsol := len(sols)
	for i, sol := range sols {
		sol.Nwins = 0
		sol.Nlosses = 0
		sol.FrontId = 0
		sol.DistCrowd = 0
		fz[i] = 0
	}

	// compute dominance data
	for i := 0; i < nsol; i++ {
		A := sols[i]
		if A.Repeated {
			continue
		}
		for j := i + 1; j < nsol; j++ {
			B := sols[j]
			if B.Repeated {
				continue
			}
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
		if sol.Repeated {
			continue
		}
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
			if A.Repeated {
				chk.Panic("_internal_error_ repeated is wrong 1")
			}
			for k := 0; k < A.Nwins; k++ {
				B := A.WinOver[k]
				if B.Repeated {
					chk.Panic("_internal_error_ repeatd is wrong 2")
				}
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
		for j := 0; j < o.prms.Nova; j++ {
			SortByOva(F, j)
			δ := o.Omax[j] - o.Omin[j] + 1e-15
			if o.prms.use_metrics_inf_crowd_dist {
				F[0].DistCrowd = INF
				F[m].DistCrowd = INF
			} else {
				F[0].DistCrowd += math.Pow((F[1].Ova[j]-F[0].Ova[j])/δ, 2.0)
				F[m].DistCrowd += math.Pow((F[m].Ova[j]-F[n].Ova[j])/δ, 2.0)
			}
			for i := 1; i < m; i++ {
				F[i].DistCrowd += ((F[i].Ova[j] - F[i-1].Ova[j]) / δ) * ((F[i+1].Ova[j] - F[i].Ova[j]) / δ)
			}
		}
	}
	return
}

// closest computes neighbour distance and sets repeated flag
func (o *Metrics) closest(A, B *Solution, dmax float64, nsol int) {
	var dist float64
	if o.prms.use_metrics_ovadistance {
		dist = A.OvaDistance(B, o.Omin, o.Omax)
	} else {
		dist = A.Distance(B, o.Fmin, o.Fmax, o.Imin, o.Imax)
	}
	if dist < A.DistNeigh {
		A.DistNeigh = dist
		A.Closest = B
	}
	if dist < B.DistNeigh {
		B.DistNeigh = dist
		B.Closest = A
	}
	if dist < o.prms.Mdmin*dmax {
		B.Repeated = o.prms.use_metrics_repeated_enabled
		if o.prms.use_metrics_repeated_enabled {
			B.FrontId = nsol
		}
	}
}

// calcdmax computes ova or flt max distance
func (o *Metrics) calcdmax() (dmax float64) {
	dmax = 1e-15
	if o.prms.use_metrics_ovadistance {
		for i := 0; i < o.prms.Nova; i++ {
			dmax += math.Pow(o.Omax[i]-o.Omin[i], 2.0)
		}
	} else {
		for i := 0; i < o.prms.Nflt; i++ {
			dmax += math.Pow(o.Fmax[i]-o.Fmin[i], 2.0)
		}
	}
	dmax = math.Sqrt(dmax)
	return
}

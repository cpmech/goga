// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Solution holds solution values
type Solution struct {

	// essential
	prms *Parameters // pointer to parameters
	Id   int         // identifier
	Ova  []float64   // objective values
	Oor  []float64   // out-of-range values
	Flt  []float64   // floats
	Int  []int       // ints

	// metrics
	WinOver   []*Solution // solutions dominated by this solution
	Nwins     int         // number of wins => current len(WinOver)
	Nlosses   int         // number of solutions dominating this solution
	FrontId   int         // Pareto front rank
	DistCrowd float64     // crowd distance
	DistNeigh float64     // closest neighbour distance
	Closest   *Solution   // closest neighbour
}

// NewSolution allocates new Solution
func NewSolution(id, nsol int, prms *Parameters) (o *Solution) {
	o = new(Solution)
	o.prms = prms
	o.Id = id
	o.Ova = make([]float64, prms.Nova)
	o.Oor = make([]float64, prms.Noor)
	o.Flt = make([]float64, prms.Nflt)
	o.Int = make([]int, prms.Nint)
	o.WinOver = make([]*Solution, nsol*2)
	return o
}

// NewSolutions allocates a number of Solutions
func NewSolutions(nsol int, prms *Parameters) (res []*Solution) {
	res = make([]*Solution, nsol)
	for i := 0; i < nsol; i++ {
		res[i] = NewSolution(i, nsol, prms)
	}
	return
}

// Feasible tells whether this solution is feasible or not
func (o *Solution) Feasible() bool {
	for _, oor := range o.Oor {
		if oor > 0 {
			return false
		}
	}
	return true
}

// CopyInto copies essential data into B
func (A *Solution) CopyInto(B *Solution) {
	B.Id = A.Id
	copy(B.Ova, A.Ova)
	copy(B.Oor, A.Oor)
	copy(B.Flt, A.Flt)
	copy(B.Int, A.Int)
}

// Distance computes (genotype) distance between A and B
func (A *Solution) Distance(B *Solution, fmin, fmax []float64, imin, imax []int) (dist float64) {
	nflt := len(A.Flt)
	if nflt > 0 {
		dflt := 0.0
		for i := 0; i < nflt; i++ {
			dflt += math.Abs(A.Flt[i]-B.Flt[i]) / (fmax[i] - fmin[i] + 1e-15)
		}
		dist += dflt / float64(nflt)
	}
	nint := len(A.Int)
	if nint > 0 {
		dint := 0.0
		for i := 0; i < nint; i++ {
			dint += math.Abs(float64(A.Int[i]-B.Int[i])) / (float64(imax[i]-imin[i]) + 1e-15)
		}
		dist += dint / float64(nint)
	}
	if nflt > 0 && nint > 0 {
		dist /= 2.0
	}
	return
}

// OvaDistance computes (phenotype) distance between A and B
func (A *Solution) OvaDistance(B *Solution, omin, omax []float64) (dist float64) {
	for i := 0; i < len(A.Ova); i++ {
		dist += math.Abs(A.Ova[i]-B.Ova[i]) / (omax[i] - omin[i] + 1e-15)
	}
	dist /= float64(len(A.Ova))
	return
}

// Compare compares two solutions
func (A *Solution) Compare(B *Solution) (A_dominates, B_dominates bool) {
	var A_nviolations, B_nviolations int
	for i := 0; i < len(A.Oor); i++ {
		if A.Oor[i] > 0 {
			A_nviolations++
		}
		if B.Oor[i] > 0 {
			B_nviolations++
		}
	}
	if A_nviolations > 0 {
		if B_nviolations > 0 {
			if A_nviolations < B_nviolations {
				A_dominates = true
				return
			}
			if B_nviolations < A_nviolations {
				B_dominates = true
				return
			}
			A_dominates, B_dominates = utl.DblsParetoMin(A.Oor, B.Oor)
			if !A_dominates && !B_dominates {
				A_dominates, B_dominates = utl.DblsParetoMin(A.Ova, B.Ova)
			}
			return
		}
		B_dominates = true
		return
	}
	if B_nviolations > 0 {
		A_dominates = true
		return
	}
	A_dominates, B_dominates = utl.DblsParetoMin(A.Ova, B.Ova)
	return
}

// Fight implements the competition between A and B
func (A *Solution) Fight(B *Solution) (A_wins bool) {
	A_dom, B_dom := A.Compare(B)
	if A_dom {
		return true
	}
	if B_dom {
		return false
	}
	if A.prms.Nova < 2 {
		if A.DistNeigh > B.DistNeigh {
			return true
		}
		if B.DistNeigh > A.DistNeigh {
			return false
		}
		if rnd.FlipCoin(0.5) {
			return true
		}
		return false
	}
	if A.FrontId == B.FrontId {
		if A.DistCrowd > B.DistCrowd {
			return true
		}
		if B.DistCrowd > A.DistCrowd {
			return false
		}
	}
	if A.FrontId == 0 {
		return true
	}
	if B.FrontId == 0 {
		return false
	}
	m := float64(B.FrontId) / float64(A.FrontId)
	prob := (1.0 - math.Exp(-10.0*m))
	if rnd.FlipCoin(prob) {
		return true
	}
	return false
}

// GetCopyResults returns a copy of results (x vectors)
func (o *Solution) GetCopyResults() (xFlt []float64, xInt []int) {
	if o.prms.Nflt > 0 {
		xFlt = make([]float64, o.prms.Nflt)
		copy(xFlt, o.Flt)
	}
	if o.prms.Nint > 0 {
		xInt = make([]int, o.prms.Nint)
		copy(xInt, o.Int)
	}
	return
}

// sorting /////////////////////////////////////////////////////////////////////////////////////////

type solByOva0 []*Solution
type solByOva1 []*Solution
type solByOva2 []*Solution
type solByBest []*Solution

func (o solByOva0) Len() int           { return len(o) }
func (o solByOva0) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva0) Less(i, j int) bool { return o[i].Ova[0] < o[j].Ova[0] }

func (o solByOva1) Len() int           { return len(o) }
func (o solByOva1) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva1) Less(i, j int) bool { return o[i].Ova[1] < o[j].Ova[1] }

func (o solByOva2) Len() int           { return len(o) }
func (o solByOva2) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva2) Less(i, j int) bool { return o[i].Ova[2] < o[j].Ova[2] }

func (o solByBest) Len() int      { return len(o) }
func (o solByBest) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o solByBest) Less(i, j int) bool {
	if o[i].FrontId == o[j].FrontId {
		return o[i].DistCrowd > o[j].DistCrowd
	}
	return o[i].FrontId < o[j].FrontId
}

// SortByOva sorts slice of solutions in ascending order of ova
func SortByOva(s []*Solution, idxOva int) {
	switch idxOva {
	case 0:
		sort.Sort(solByOva0(s))
	case 1:
		sort.Sort(solByOva1(s))
	case 2:
		sort.Sort(solByOva2(s))
	default:
		chk.Panic("this code can only handle Nova ≤ 3 for now")
	}
}

// SortByBest sorts slice of solutions with best solutions first
func SortByBest(s []*Solution) {
	sort.Sort(solByBest(s))
}

// TODO
func SortByTradeoff(s []*Solution) {
	chk.Panic("SortByTradeoff: TODO")
}

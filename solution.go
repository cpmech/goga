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
	Id  int       // identifier
	Ova []float64 // objective values
	Oor []float64 // out-of-range values
	Flt []float64 // floats
	Int []int     // ints

	// metrics
	WinOver   []*Solution // solutions dominated by this solution
	Nwins     int         // number of wins => current len(WinOver)
	Nlosses   int         // number of solutions dominating this solution
	FrontId   int         // Pareto front rank
	DistCrowd float64     // crowd distance
	DistNeigh float64     // minimum distance to any neighbouring solution
}

// NewSolution allocates new Solution
func NewSolution(id, nsol int, prms *Parameters) (o *Solution) {
	o = new(Solution)
	o.Id = id
	o.Ova = make([]float64, prms.Nova)
	o.Oor = make([]float64, prms.Noor)
	o.Flt = make([]float64, prms.Nflt)
	o.Int = make([]int, prms.Nint)
	o.WinOver = make([]*Solution, nsol)
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
	for i := 0; i < len(A.Flt); i++ {
		dist += math.Abs(A.Flt[i]-B.Flt[i]) / (fmax[i] - fmin[i] + 1e-15)
	}
	for i := 0; i < len(A.Int); i++ {
		dist += math.Abs(float64(A.Int[i]-B.Int[i])) / (float64(imax[i]-imin[i]) + 1e-15)
	}
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
	if A.FrontId == B.FrontId {
		if A.DistCrowd > B.DistCrowd {
			return true
		}
		if B.DistCrowd > A.DistCrowd {
			return false
		}
	}
	if A.DistNeigh > B.DistNeigh {
		return true
	}
	if B.DistNeigh > A.DistNeigh {
		return true
	}
	if rnd.FlipCoin(0.5) {
		return true
	}
	return false
}

// sorting /////////////////////////////////////////////////////////////////////////////////////////

type solByOva0 []*Solution
type solByOva1 []*Solution
type solByOva2 []*Solution

func (o solByOva0) Len() int           { return len(o) }
func (o solByOva0) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva0) Less(i, j int) bool { return o[i].Ova[0] < o[j].Ova[0] }

func (o solByOva1) Len() int           { return len(o) }
func (o solByOva1) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva1) Less(i, j int) bool { return o[i].Ova[1] < o[j].Ova[1] }

func (o solByOva2) Len() int           { return len(o) }
func (o solByOva2) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva2) Less(i, j int) bool { return o[i].Ova[2] < o[j].Ova[2] }

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

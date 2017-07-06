// Copyright 2015 The Goga Authors. All rights reserved.
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
	prms  *Parameters // pointer to parameters
	Id    int         // identifier (for debugging). future solutions have negative ids
	Fixed bool        // cannot be changed
	Ova   []float64   // objective values
	Oor   []float64   // out-of-range values
	Flt   []float64   // floats
	Int   []int       // ints

	// metrics
	WinOver   []*Solution // solutions dominated by this solution
	Nwins     int         // number of wins => current len(WinOver)
	Nlosses   int         // number of solutions dominating this solution
	FrontId   int         // Pareto front rank
	DistCrowd float64     // crowd distance
	DistNeigh float64     // closest neighbour distance
	Closest   *Solution   // closest neighbour

	// auxiliary
	Aux float64 // auxiliary data to be stored at each solution; e.g. limit state function value
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

// Reset rests state; i.e. zeroes all values
func (o *Solution) Reset(id int) {

	// essential
	o.Id = id
	o.Fixed = false
	utl.Fill(o.Ova, 0)
	utl.Fill(o.Oor, 0)
	utl.Fill(o.Flt, 0)
	utl.IntFill(o.Int, 0)

	// metrics
	for i := 0; i < len(o.WinOver); i++ {
		o.WinOver[i] = nil
	}
	o.Nwins = 0
	o.Nlosses = 0
	o.FrontId = 0
	o.DistCrowd = 0
	o.DistNeigh = 0
	o.Closest = nil

	// auxiliary
	o.Aux = 0
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
			A_dominates, B_dominates = utl.ParetoMin(A.Oor, B.Oor)
			if !A_dominates && !B_dominates {
				A_dominates, B_dominates = utl.ParetoMin(A.Ova, B.Ova)
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
	A_dominates, B_dominates = utl.ParetoMin(A.Ova, B.Ova)
	return
}

// Fight implements the competition between A and B
func (A *Solution) Fight(B *Solution) (A_wins bool) {

	// compare solutions
	A_dom, B_dom := A.Compare(B)
	if A_dom {
		return true
	}
	if B_dom {
		return false
	}

	// tie: single-objective problems
	if A.prms.Nova < 2 {
		if A.DistNeigh > B.DistNeigh {
			return true
		}
		if B.DistNeigh > A.DistNeigh {
			return false
		}
		return rnd.FlipCoin(0.5)
	}

	// tie: multi-objective problems: same Pareto front
	if A.FrontId == B.FrontId {
		if A.DistCrowd > B.DistCrowd {
			return true
		}
		if B.DistCrowd > A.DistCrowd {
			return false
		}
		return rnd.FlipCoin(0.5)
	}

	// tie: multi-objective problems: different Pareto fronts
	if A.FrontId < B.FrontId {
		return true
	}
	if B.FrontId < A.FrontId {
		return false
	}
	if A.DistNeigh > B.DistNeigh {
		return true
	}
	if B.DistNeigh > A.DistNeigh {
		return false
	}
	return rnd.FlipCoin(0.5)
}

// sorting /////////////////////////////////////////////////////////////////////////////////////////

// SortSolutions sort solutions either by OVA (single-obj) or Pareto front (multi-obj)
func SortSolutions(s []*Solution, idxOva int) {
	if len(s) > 0 {
		nova := len(s[0].Ova)
		if nova > 1 { // multi-objective
			sortByFrontThenOva(s, idxOva)
		} else { // single-objective
			sortByOva(s, idxOva)
		}
	}
}

////////////////////////////////////////////////////////////
// TODO: Improve this part to handle any number of Ovas ////
////////////////////////////////////////////////////////////

type solByOva0 []*Solution
type solByOva1 []*Solution
type solByOva2 []*Solution
type solByOva3 []*Solution
type solByOva4 []*Solution
type solByOva5 []*Solution
type solByOva6 []*Solution
type solByOva7 []*Solution
type solByOva8 []*Solution
type solByOva9 []*Solution
type solByOva10 []*Solution
type solByOva11 []*Solution
type solByOva12 []*Solution
type solByOva13 []*Solution
type solByOva14 []*Solution
type solByOva15 []*Solution
type solByOva16 []*Solution
type solByOva17 []*Solution
type solByOva18 []*Solution
type solByOva19 []*Solution

func (o solByOva0) Len() int           { return len(o) }
func (o solByOva0) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva0) Less(i, j int) bool { return o[i].Ova[0] < o[j].Ova[0] }

func (o solByOva1) Len() int           { return len(o) }
func (o solByOva1) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva1) Less(i, j int) bool { return o[i].Ova[1] < o[j].Ova[1] }

func (o solByOva2) Len() int           { return len(o) }
func (o solByOva2) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva2) Less(i, j int) bool { return o[i].Ova[2] < o[j].Ova[2] }

func (o solByOva3) Len() int           { return len(o) }
func (o solByOva3) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva3) Less(i, j int) bool { return o[i].Ova[3] < o[j].Ova[3] }

func (o solByOva4) Len() int           { return len(o) }
func (o solByOva4) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva4) Less(i, j int) bool { return o[i].Ova[4] < o[j].Ova[4] }

func (o solByOva5) Len() int           { return len(o) }
func (o solByOva5) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva5) Less(i, j int) bool { return o[i].Ova[5] < o[j].Ova[5] }

func (o solByOva6) Len() int           { return len(o) }
func (o solByOva6) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva6) Less(i, j int) bool { return o[i].Ova[6] < o[j].Ova[6] }

func (o solByOva7) Len() int           { return len(o) }
func (o solByOva7) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva7) Less(i, j int) bool { return o[i].Ova[7] < o[j].Ova[7] }

func (o solByOva8) Len() int           { return len(o) }
func (o solByOva8) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva8) Less(i, j int) bool { return o[i].Ova[8] < o[j].Ova[8] }

func (o solByOva9) Len() int           { return len(o) }
func (o solByOva9) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva9) Less(i, j int) bool { return o[i].Ova[9] < o[j].Ova[9] }

func (o solByOva10) Len() int           { return len(o) }
func (o solByOva10) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva10) Less(i, j int) bool { return o[i].Ova[10] < o[j].Ova[10] }

func (o solByOva11) Len() int           { return len(o) }
func (o solByOva11) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva11) Less(i, j int) bool { return o[i].Ova[11] < o[j].Ova[11] }

func (o solByOva12) Len() int           { return len(o) }
func (o solByOva12) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva12) Less(i, j int) bool { return o[i].Ova[12] < o[j].Ova[12] }

func (o solByOva13) Len() int           { return len(o) }
func (o solByOva13) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva13) Less(i, j int) bool { return o[i].Ova[13] < o[j].Ova[13] }

func (o solByOva14) Len() int           { return len(o) }
func (o solByOva14) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva14) Less(i, j int) bool { return o[i].Ova[14] < o[j].Ova[14] }

func (o solByOva15) Len() int           { return len(o) }
func (o solByOva15) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva15) Less(i, j int) bool { return o[i].Ova[15] < o[j].Ova[15] }

func (o solByOva16) Len() int           { return len(o) }
func (o solByOva16) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva16) Less(i, j int) bool { return o[i].Ova[16] < o[j].Ova[16] }

func (o solByOva17) Len() int           { return len(o) }
func (o solByOva17) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva17) Less(i, j int) bool { return o[i].Ova[17] < o[j].Ova[17] }

func (o solByOva18) Len() int           { return len(o) }
func (o solByOva18) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva18) Less(i, j int) bool { return o[i].Ova[18] < o[j].Ova[18] }

func (o solByOva19) Len() int           { return len(o) }
func (o solByOva19) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o solByOva19) Less(i, j int) bool { return o[i].Ova[19] < o[j].Ova[19] }

type solByFrontThenOva0 []*Solution

func (o solByFrontThenOva0) Len() int      { return len(o) }
func (o solByFrontThenOva0) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o solByFrontThenOva0) Less(i, j int) bool {
	if o[i].FrontId == o[j].FrontId {
		return o[i].Ova[0] < o[j].Ova[0]
	}
	return o[i].FrontId < o[j].FrontId
}

// sortByOva sorts slice of solutions in ascending order of ova
func sortByOva(s []*Solution, idxOva int) {
	switch idxOva {
	case 0:
		sort.Sort(solByOva0(s))
	case 1:
		sort.Sort(solByOva1(s))
	case 2:
		sort.Sort(solByOva2(s))
	case 3:
		sort.Sort(solByOva3(s))
	case 4:
		sort.Sort(solByOva4(s))
	case 5:
		sort.Sort(solByOva5(s))
	case 6:
		sort.Sort(solByOva6(s))
	case 7:
		sort.Sort(solByOva7(s))
	case 8:
		sort.Sort(solByOva8(s))
	case 9:
		sort.Sort(solByOva9(s))
	case 10:
		sort.Sort(solByOva10(s))
	case 11:
		sort.Sort(solByOva11(s))
	case 12:
		sort.Sort(solByOva12(s))
	case 13:
		sort.Sort(solByOva13(s))
	case 14:
		sort.Sort(solByOva14(s))
	case 15:
		sort.Sort(solByOva15(s))
	case 16:
		sort.Sort(solByOva16(s))
	case 17:
		sort.Sort(solByOva17(s))
	case 18:
		sort.Sort(solByOva18(s))
	case 19:
		sort.Sort(solByOva19(s))
	default:
		chk.Panic("this code can only handle Nova ≤ 20 for now")
	}
}

// sortByFrontThenOva sorts solutions first by front and then by ova
func sortByFrontThenOva(s []*Solution, idxOva int) {
	switch idxOva {
	case 0:
		sort.Sort(solByFrontThenOva0(s))
	default:
		chk.Panic("this code can only handle Nova ≤ 1 for now")
	}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Individual implements one individual in a population
type Individual struct {

	// essential data
	Ovas   []float64 // objective values
	Oors   []float64 // out-of-range values: sum of positive distances from constraints
	Floats []float64 // [optional] floats
	Ints   []int     // [optional] integers

	// auxiliary
	Id int // identifier; e.g. for debugging

	// comparison, diversity and non-dominance sorting data
	WinOver   []*Individual // [ninds] individuals dominated by this individual
	Nwins     int           // number of wins => current len(WinOver)
	Nlosses   int           // number of individuals dominating this individual
	FrontId   int           // Pareto front id
	DistCrowd float64       // crowd distance
	DistNeigh float64       // minimum distance to any neighbouring individual
}

// NewIndividual allocates a new individual
func NewIndividual(nova, noor int, floats []float64, ints []int) (o *Individual) {
	o = new(Individual)
	o.Ovas = make([]float64, nova)
	o.Oors = make([]float64, noor)
	if len(floats) > 0 {
		o.Floats = make([]float64, len(floats))
		copy(o.Floats, floats)
	}
	if len(ints) > 0 {
		o.Ints = make([]int, len(ints))
		copy(o.Ints, ints)
	}
	return
}

// GetCopy returns a copy of this individual
func (o Individual) GetCopy() (x *Individual) {
	x = new(Individual)
	x.Ovas = make([]float64, len(o.Ovas))
	x.Oors = make([]float64, len(o.Oors))
	copy(x.Ovas, o.Ovas)
	copy(x.Oors, o.Oors)
	if len(o.Floats) > 0 {
		x.Floats = make([]float64, len(o.Floats))
		copy(x.Floats, o.Floats)
	}
	if len(o.Ints) > 0 {
		x.Ints = make([]int, len(o.Ints))
		copy(x.Ints, o.Ints)
	}
	x.Id = o.Id
	return
}

// CopyInto copies this individual's data into another individual
func (o Individual) CopyInto(x *Individual) {
	copy(x.Ovas, o.Ovas)
	copy(x.Oors, o.Oors)
	if len(o.Floats) > 0 {
		copy(x.Floats, o.Floats)
	}
	if len(o.Ints) > 0 {
		copy(x.Ints, o.Ints)
	}
	x.Id = o.Id
	return
}

// Feasible returns whether this individual is feasible or not
func (o Individual) Feasible() bool {
	for _, oor := range o.Oors {
		if oor > 0 {
			return false
		}
	}
	return true
}

// Fight implements the competition between A and B
func (A *Individual) Fight(B *Individual) (A_wins bool) {
	A_dom, B_dom := IndCompare(A, B)
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
	//if false {
	if true {
		//io.Pforan("A.dist=%v  B.dist=%v\n", A.DistNeigh, B.DistNeigh)
		if A.DistNeigh > B.DistNeigh {
			return true
		}
		if B.DistNeigh > A.DistNeigh {
			return true
		}
	}
	if rnd.FlipCoin(0.5) {
		return true
	}
	return false
}

// global functions ////////////////////////////////////////////////////////////////////////////////

// IndCompare compares individual 'A' with another one 'B'. Deterministic method
func IndCompare(A, B *Individual) (A_dominates, B_dominates bool) {
	var A_nviolations, B_nviolations int
	for i := 0; i < len(A.Oors); i++ {
		if A.Oors[i] > 0 {
			A_nviolations++
		}
		if B.Oors[i] > 0 {
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
			A_dominates, B_dominates = utl.DblsParetoMin(A.Oors, B.Oors)
			if !A_dominates && !B_dominates {
				A_dominates, B_dominates = utl.DblsParetoMin(A.Ovas, B.Ovas)
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
	A_dominates, B_dominates = utl.DblsParetoMin(A.Ovas, B.Ovas)
	return
}

// IndDistance computes a distance measure from individual 'A' to another individual 'B'
func IndDistance(A, B *Individual, omin, omax []float64) (dist float64) {
	for i := 0; i < len(A.Ovas); i++ {
		δ := omax[i] - omin[i] + 1e-15
		dist += math.Pow((A.Ovas[i]-B.Ovas[i])/δ, 2.0)
	}
	return math.Sqrt(dist)
}

// IndCrossover performs the crossover between chromosomes of two individuals A and B
// resulting in the chromosomes of other two individuals a and b
func IndCrossover(a, b, A, B, C, D *Individual, time int, ops *OpsData) {
	if len(A.Floats) > 0 {
		if D == nil {
			ops.CxFlt(a.Floats, b.Floats, A.Floats, B.Floats, nil, nil, time, ops)
		} else {
			ops.CxFlt(a.Floats, b.Floats, A.Floats, B.Floats, C.Floats, D.Floats, time, ops)
		}
	}
	if len(A.Ints) > 0 {
		if D == nil {
			ops.CxInt(a.Ints, b.Ints, A.Ints, B.Ints, nil, nil, time, ops)
		} else {
			ops.CxInt(a.Ints, b.Ints, A.Ints, B.Ints, C.Ints, D.Ints, time, ops)
		}
	}
}

// IndMutation performs the mutation operation in the chromosomes of an individual
func IndMutation(A *Individual, time int, ops *OpsData) {
	if len(A.Floats) > 0 {
		ops.MtFlt(A.Floats, time, ops)
	}
	if len(A.Ints) > 0 {
		ops.MtInt(A.Ints, time, ops)
	}
}

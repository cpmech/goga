// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/graph"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// constants
const (
	INF = 1e+30 // infinite distance
)

// Island holds one population and performs the reproduction operation
type Island struct {

	// input
	Id  int         // index of this island
	C   *ConfParams // configuration parameters
	Pop Population  // population

	// crowding
	Indices     []int           // [ninds] all indices of individuals
	Groups      [][]int         // [ngroups][nparents] indices defining groups of individuals
	Competitors []*Individual   // [ngroups*nparents*nparents] all competitors
	Parents     [][]*Individual // [ngroups][nparents] all parents (view to Competitors)
	Offspring   [][]*Individual // [ngroups][noffspri] all offspring (view to Competitors)
	Mdist       [][]float64     // [nparents][noffspring] matching distances
	Match       graph.Munkres   // matches

	// metrics
	OvaMin []float64       // min ova
	OvaMax []float64       // max ova
	Fsizes []int           // front sizes
	Fronts [][]*Individual // non-dominated fronts

	// results
	Nfeval int // number of objective function evaluations
}

// NewIsland creates a new island
func NewIsland(id int, C *ConfParams) (o *Island) {

	// allocate island
	o = new(Island)
	o.Id = id
	o.C = C
	o.GenPop()
	o.CalcOvs()

	// constants
	ni := o.C.Ninds        // number of individuals
	np := 2                // number of parents in group
	no := np * (np - 1)    // number of offspring in group
	ng := ni / np          // number of groups
	nr := np * np          // number of individuals in round (parents + offspring)
	nc := ng * nr          // total number of competitors
	nn := utl.Imax(ni, nc) // max between ninds and ncompetitors
	nv := o.C.Nova         // number of objective values

	// crowding
	o.Indices = utl.IntRange(ni)
	o.Groups = utl.IntsAlloc(ng, np)
	o.Competitors = make([]*Individual, nc)
	for i := 0; i < nc; i++ {
		o.Competitors[i] = o.Pop[0].GetCopy()
	}
	o.Parents = make([][]*Individual, ng)
	o.Offspring = make([][]*Individual, ng)
	r := 0
	for k := 0; k < ng; k++ {
		o.Parents[k] = make([]*Individual, np)
		o.Offspring[k] = make([]*Individual, no)
		s := 0
		for i := 0; i < np; i++ {
			o.Parents[k][i], r = o.Competitors[r], r+1
			for j := i + 1; j < np; j++ {
				o.Offspring[k][s], r, s = o.Competitors[r], r+1, s+1
				o.Offspring[k][s], r, s = o.Competitors[r], r+1, s+1
			}
		}
	}
	o.Mdist = la.MatAlloc(np, no)
	o.Match.Init(np, no)

	// metrics
	o.OvaMin, o.OvaMax = make([]float64, nv), make([]float64, nv)
	o.Fsizes = make([]int, nn)
	o.Fronts = make([][]*Individual, nn)
	for i := 0; i < nn; i++ {
		o.Fronts[i] = make([]*Individual, nn)
	}
	return
}

// GenPop generates population
func (o *Island) GenPop() {
	if o.C.PopIntGen != nil {
		o.Pop = o.C.PopIntGen(o.Id, o.C)
	}
	if o.C.PopFltGen != nil {
		o.Pop = o.C.PopFltGen(o.Id, o.C)
	}
	if len(o.Pop) != o.C.Ninds {
		chk.Panic("generation of population failed")
	}
}

// CalcOvs computes objective values
func (o *Island) CalcOvs() {
	o.Nfeval = 0
	for _, ind := range o.Pop {
		o.C.OvaOor(o.Id, ind)
		o.Nfeval += 1
		for _, oor := range ind.Oors {
			if oor < 0 {
				chk.Panic("out-of-range values must be positive (or zero) indicating the positive distance to constraints. oor=%g is invalid", oor)
			}
		}
	}
}

// Reset resets population
func (o *Island) Reset() {
	o.GenPop()
	o.CalcOvs()
}

// Run runs evolutionary process with niching via crowding and tournament selection
func (o *Island) Run(time int) {

	// select groups
	rnd.IntGetGroups(o.Groups, o.Indices)

	// constants
	ni := o.C.Ninds     // number of individuals
	np := 2             // number of parents in group
	no := np * (np - 1) // number of offspring in group
	ng := ni / np       // number of groups

	// set parents
	for k := 0; k < ng; k++ {
		for i := 0; i < np; i++ {
			o.Pop[o.Groups[k][i]].CopyInto(o.Parents[k][i])
		}
	}

	// create offspring and set competitors
	var a, b, A, B, C, D *Individual
	for k := 0; k < ng; k++ {
		s := 0
		for i := 0; i < np; i++ {
			A = o.Parents[k][i]
			for j := i + 1; j < np; j++ {
				B = o.Parents[k][j]
				if o.C.Ops.Use4inds {
					knext := (k + 1) % ng
					C = o.Pop[o.Groups[knext][0]]
					D = o.Pop[o.Groups[knext][1]]
				}
				a, s = o.Offspring[k][s], s+1
				b, s = o.Offspring[k][s], s+1
				IndCrossover(a, b, A, B, C, D, time, &o.C.Ops)
				IndMutation(a, time, &o.C.Ops)
				IndMutation(b, time, &o.C.Ops)
				o.C.OvaOor(o.Id, a)
				o.C.OvaOor(o.Id, b)
				o.Nfeval += 2
			}
		}
	}

	// metrics: competitors
	Metrics(o.OvaMin, o.OvaMax, o.Fsizes, o.Fronts, o.Competitors)

	// tournaments
	idxnew := 0
	for k := 0; k < ng; k++ {

		// compute match distances
		for i := 0; i < np; i++ {
			A = o.Parents[k][i]
			for j := 0; j < no; j++ {
				a = o.Offspring[k][j]
				o.Mdist[i][j] = IndDistance(A, a, o.OvaMin, o.OvaMax)
			}
		}
		//la.PrintMat("mdist", o.Mdist, "%8.5f", false)

		// match competitors
		o.Match.SetCostMatrix(o.Mdist)
		o.Match.Run()
		//io.Pforan("links = %v\n", o.Match.Links)

		// matches
		for i := 0; i < np; i++ {
			A = o.Parents[k][i]
			B = o.Offspring[k][o.Match.Links[i]]
			if A.Fight(B) {
				A.CopyInto(o.Pop[idxnew]) // A wins
			} else {
				B.CopyInto(o.Pop[idxnew]) // B wins
			}
			idxnew++
		}
	}
}

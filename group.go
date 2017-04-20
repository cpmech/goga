// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/utl"

// Group holds a group of solutions
type Group struct {
	Ncur    int         // number of current solutions == len(All) / 2
	All     []*Solution // current and future solutions. half part is a view to Solutions
	Indices []int       // indices of current solutions
	Pairs   [][]int     // randomly selected pairs from Indices
	Metrics *Metrics    // metrics
}

// Init initialises group
func (o *Group) Init(cpu, ncpu int, solutions []*Solution, prms *Parameters) {
	nsol := len(solutions)
	start, endp1 := (cpu*nsol)/ncpu, ((cpu+1)*nsol)/ncpu
	o.Ncur = endp1 - start
	o.All = make([]*Solution, o.Ncur*2)
	o.Indices = make([]int, o.Ncur)
	o.Pairs = utl.IntAlloc(o.Ncur/2, 2)
	for i := 0; i < o.Ncur; i++ {
		o.All[i] = solutions[start+i]
		o.All[o.Ncur+i] = NewSolution(-(1 + i), nsol, prms) // the index is for debugging
		o.Indices[i] = i
	}
	o.Metrics = new(Metrics)
	o.Metrics.Init(len(o.All), prms)
}

// Reset resets group data
func (o *Group) Reset(cpu, ncpu int, solutions []*Solution) {
	nsol := len(solutions)
	start := (cpu * nsol) / ncpu
	for i, pp := range o.Pairs {
		for j, _ := range pp {
			o.Pairs[i][j] = 0
		}
	}
	for i := 0; i < o.Ncur; i++ {
		o.All[i] = solutions[start+i]
		o.All[o.Ncur+i].Reset(-(1 + i)) // there is no real need for this; but helps with debugging
	}
}

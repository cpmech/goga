// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"sort"

	"github.com/cpmech/gosl/utl"
)

// ValIndPair holds a value and a pointer to an individual
type ValIndPair struct {
	Val float64
	Ind *Individual
}

// collections
type ValIndAsc []ValIndPair // ascending
type ValIndDes []ValIndPair // descending

// for sorting
func (o ValIndAsc) Len() int           { return len(o) }
func (o ValIndAsc) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o ValIndAsc) Less(i, j int) bool { return o[i].Val < o[j].Val }
func (o ValIndDes) Len() int           { return len(o) }
func (o ValIndDes) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o ValIndDes) Less(i, j int) bool { return o[i].Val > o[j].Val }
func (o *ValIndAsc) Sort()             { sort.Sort(o) }
func (o *ValIndDes) Sort()             { sort.Sort(o) }

// NonDomSort finds Pareto fronts and sort individuals according to domination ranks
func (o *Island) NonDomSort(pop Population) {

	// reset counters
	ninds := len(pop)
	for i := 0; i < ninds; i++ {
		o.fsizes[i], o.sdom[i], o.ndby[i] = 0, 0, 0
	}

	// compute dominance data
	for i := 0; i < ninds; i++ {
		A := pop[i]
		for j := i + 1; j < ninds; j++ {
			B := pop[j]
			Adom, Bdom := IndCompareDet(A, B)
			if Adom {
				o.idom[i][o.sdom[i]] = j // i dominates j
				o.sdom[i]++              // i has another dominated item
				o.ndby[j]++              // j is being dominated by i
			}
			if Bdom {
				o.idom[j][o.sdom[j]] = i // j dominates i
				o.sdom[j]++              // j has another dominated item
				o.ndby[i]++              // i is being dominated by j
			}
		}
	}

	// first front
	for i := 0; i < ninds; i++ {
		if o.ndby[i] == 0 {
			o.fronts[0][o.fsizes[0]] = i
			o.fsizes[0]++
		}
	}

	// next fronts
	o.nfronts = 0
	for r, front := range o.fronts {
		if o.fsizes[r] == 0 {
			break
		}
		o.nfronts++
		for s := 0; s < o.fsizes[r]; s++ {
			i := front[s]
			for m := 0; m < o.sdom[i]; m++ {
				j := o.idom[i][m]
				o.ndby[j]--
				if o.ndby[j] == 0 { // j belongs to the next front
					o.fronts[r+1][o.fsizes[r+1]] = j
					o.fsizes[r+1]++
				}
			}
		}
	}

	// set front id
	for r := 0; r < o.nfronts; r++ {
		for s := 0; s < o.fsizes[r]; s++ {
			i := o.fronts[r][s]
			pop[i].FrontId = r
		}
	}
}

// CalcMinMaxOva computes min and max values of ovas
func (o *Island) CalcMinMaxOva(pop Population) {
	ninds := len(pop)
	for i := 0; i < ninds; i++ {
		for j := 0; j < o.C.Nova; j++ {
			ova := pop[i].Ovas[j]
			if i == 0 {
				o.ovamin[j] = ova
				o.ovamax[j] = ova
			} else {
				o.ovamin[j] = utl.Min(o.ovamin[j], ova)
				o.ovamax[j] = utl.Max(o.ovamax[j], ova)
			}
		}
	}
}

// CalcCrowdDist compute crowd distances considering each Pareto front
//  Note: (1) non-dominated Pareto fronts must be computed first
//        (2) min/max values of objective values must be computed first
func (o *Island) CalcCrowdDist(pop Population) {
	for _, ind := range pop {
		ind.Cdist = 0
	}
	for r := 0; r < o.nfronts; r++ {
		var F ValIndAsc
		F = make([]ValIndPair, o.fsizes[r])
		l, m, n := o.fsizes[r], o.fsizes[r]-1, o.fsizes[r]-2
		if l == 1 {
			//pop[o.fronts[r][0]].Cdist = INF
			continue
		}
		for j := 0; j < o.C.Nova; j++ {
			for s := 0; s < o.fsizes[r]; s++ {
				i := o.fronts[r][s]
				F[s].Ind = pop[i]
				F[s].Val = pop[i].Ovas[j]
			}
			δ := o.ovamax[j] - o.ovamin[j] + 1e-15
			F.Sort()
			//F[0].Ind.Cdist += math.Pow((F[1].Ind.Ovas[j]-F[0].Ind.Ovas[j])/δ, 2.0)
			//F[m].Ind.Cdist += math.Pow((F[m].Ind.Ovas[j]-F[n].Ind.Ovas[j])/δ, 2.0)
			_ = n
			F[0].Ind.Cdist = INF
			F[m].Ind.Cdist = INF
			for i := 1; i < m; i++ {
				F[i].Ind.Cdist += ((F[i].Ind.Ovas[j] - F[i-1].Ind.Ovas[j]) / δ) * ((F[i+1].Ind.Ovas[j] - F[i].Ind.Ovas[j]) / δ)
			}
		}
	}
	//for _, ind := range pop {
	//io.Pforan("d = %v\n", ind.Cdist)
	//}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"sort"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
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

func Metrics(ovamin, ovamax, fltmin, fltmax []float64, intmin, intmax, fsizes []int, fronts [][]*Individual, pop Population) (nfronts int) {

	// reset counters and find limits
	ninds := len(pop)
	for i, ind := range pop {

		// allocate slice
		if len(ind.WinOver) != ninds {
			ind.WinOver = make([]*Individual, ninds)
		}

		// reset values
		ind.Repeated = false
		ind.Nwins = 0
		ind.Nlosses = 0
		ind.FrontId = 0
		ind.DistCrowd = 0
		ind.DistNeigh = INF
		fsizes[i] = 0

		// calc ova limits
		for j, x := range ind.Ovas {
			if i == 0 {
				ovamin[j], ovamax[j] = x, x
			} else {
				ovamin[j] = utl.Min(ovamin[j], x)
				ovamax[j] = utl.Max(ovamax[j], x)
			}
		}

		// calc flt limits
		for j, x := range ind.Floats {
			if i == 0 {
				fltmin[j], fltmax[j] = x, x
			} else {
				fltmin[j] = utl.Min(fltmin[j], x)
				fltmax[j] = utl.Max(fltmax[j], x)
			}
		}

		// calc int limits
		for j, x := range ind.Ints {
			if i == 0 {
				intmin[j], intmax[j] = x, x
			} else {
				intmin[j] = utl.Imin(intmin[j], x)
				intmax[j] = utl.Imax(intmax[j], x)
			}
		}
	}

	// compute distance (genotype) and mark repeated
	DMIN := 1e-8
	for i := 0; i < ninds; i++ {
		A := pop[i]
		for j := i + 1; j < ninds; j++ {
			B := pop[j]
			d := IndDistance(A, B, intmin, intmax, fltmin, fltmax, nil, nil, false)
			if d < A.DistNeigh {
				A.DistNeigh = d
			}
			if d < B.DistNeigh {
				B.DistNeigh = d
			}
			if B.DistNeigh < DMIN {
				B.Repeated = true
			}
		}
	}

	// compute dominance data
	for i := 0; i < ninds; i++ {
		A := pop[i]
		if A.Repeated {
			continue
		}
		for j := i + 1; j < ninds; j++ {
			B := pop[j]
			if B.Repeated {
				continue
			}
			Adom, Bdom := IndCompareDet(A, B)
			if Adom {
				A.WinOver[A.Nwins] = B // i dominates j
				A.Nwins++              // i has another dominated item
				B.Nlosses++            // j is being dominated by i
			}
			if Bdom {
				B.WinOver[B.Nwins] = A // j dominates i
				B.Nwins++              // j has another dominated item
				A.Nlosses++            // i is being dominated by j
			}
		}
	}

	// first front
	for _, ind := range pop {
		if ind.Repeated {
			continue
		}
		if ind.Nlosses == 0 {
			fronts[0][fsizes[0]] = ind
			fsizes[0]++
		}
	}

	// next fronts
	for r, front := range fronts {
		if fsizes[r] == 0 {
			break
		}
		nfronts++
		for s := 0; s < fsizes[r]; s++ {
			A := front[s]
			for k := 0; k < A.Nwins; k++ {
				B := A.WinOver[k]
				B.Nlosses--
				if B.Nlosses == 0 { // B belongs to next front
					B.FrontId = r + 1
					fronts[r+1][fsizes[r+1]] = B
					fsizes[r+1]++
				}
			}
		}
	}

	// crowd distances
	nova := len(pop[0].Ovas)
	for r := 0; r < nfronts; r++ {
		l, m, n := fsizes[r], fsizes[r]-1, fsizes[r]-2
		if l == 1 {
			fronts[r][0].DistCrowd = -1
			continue
		}
		F := fronts[r][:l]
		for j := 0; j < nova; j++ {
			SortPopByOva(F, j)
			δ := ovamax[j] - ovamin[j] + 1e-15
			if true {
				F[0].DistCrowd += math.Pow((F[1].Ovas[j]-F[0].Ovas[j])/δ, 2.0)
				F[m].DistCrowd += math.Pow((F[m].Ovas[j]-F[n].Ovas[j])/δ, 2.0)
			} else {
				F[0].DistCrowd = INF
				F[m].DistCrowd = INF
			}
			for i := 1; i < m; i++ {
				F[i].DistCrowd += ((F[i].Ovas[j] - F[i-1].Ovas[j]) / δ) * ((F[i+1].Ovas[j] - F[i].Ovas[j]) / δ)
			}
		}
	}
	return
}

// NonDomSort finds Pareto fronts and sort individuals according to domination ranks.
// This function also sets FrontId in individuals
func (o *Island) NonDomSortAndSetFrontId(pop Population) {

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

// CalcAndSetDistances computes minimum (neighbour) and crowd distances.
// This function also sets Ndist and Cdist in individuals
//  Notes: 1. non-dominated Pareto fronts must be computed first
//         2. min/max values of objective values must be computed first
func (o *Island) CalcAndSetDistances(pop Population) {

	// reset variables
	for _, ind := range pop {
		ind.DistNeigh = INF
		ind.DistCrowd = 0
	}

	// minimum distances
	ninds := len(pop)
	//nova := len(pop[0].Ovas)
	for i := 0; i < ninds; i++ {
		o.ndist[i][i] = INF
		for j := i + 1; j < ninds; j++ {
			o.ndist[i][j] = IndDistance(pop[i], pop[j], nil, nil, nil, nil, o.ovamin, o.ovamax, true)
			//d := 0.0
			//for k := 0; k < nova; k++ {
			//δ := o.ovamax[k] - o.ovamin[k] + 1e-15
			//d += math.Pow((pop[i].Ovas[k]-pop[j].Ovas[k])/δ, 2.0)
			//}
			//o.ndist[i][j] = math.Sqrt(d)
			o.ndist[j][i] = o.ndist[i][j]
		}
	}
	for i, ind := range pop {
		ind.DistNeigh = la.VecMin(o.ndist[i])
	}

	// crowd distances
	for r := 0; r < o.nfronts; r++ {
		var F ValIndAsc
		F = make([]ValIndPair, o.fsizes[r])
		l, m, n := o.fsizes[r], o.fsizes[r]-1, o.fsizes[r]-2
		if l == 1 {
			pop[o.fronts[r][0]].DistCrowd = INF
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
			if false {
				F[0].Ind.DistCrowd += math.Pow((F[1].Ind.Ovas[j]-F[0].Ind.Ovas[j])/δ, 2.0)
				F[m].Ind.DistCrowd += math.Pow((F[m].Ind.Ovas[j]-F[n].Ind.Ovas[j])/δ, 2.0)
			} else {
				F[0].Ind.DistCrowd = INF
				F[m].Ind.DistCrowd = INF
			}
			for i := 1; i < m; i++ {
				F[i].Ind.DistCrowd += ((F[i].Ind.Ovas[j] - F[i-1].Ind.Ovas[j]) / δ) * ((F[i+1].Ind.Ovas[j] - F[i].Ind.Ovas[j]) / δ)
			}
		}
	}

	// debug
	//if true {
	if false {
		for _, ind := range pop {
			io.Pforan("ndist=%25g cdist=%25g\n", ind.DistNeigh, ind.DistCrowd)
		}
	}
}

// NomDomSortAndCalcDistances calculates non-dominated front, sort population and compute distances
// This function also does:
//  1) computes the min and max Ova values
//  2) set Ndist, Cdist and FrontId of individuals
func (o *Island) NomDomSortAndCalcDistances(pop Population) {
	o.NonDomSortAndSetFrontId(pop)
	o.CalcMinMaxOva(pop)
	o.CalcAndSetDistances(pop)
}

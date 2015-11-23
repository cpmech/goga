// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

func Metrics(ovamin, ovamax []float64, fsizes []int, fronts [][]*Individual, pop Population) (nfronts int) {

	// reset counters and find limits
	nv := len(ovamin)
	ninds := len(pop)
	for i, ind := range pop {

		// allocate slice
		if len(ind.WinOver) != ninds {
			ind.WinOver = make([]*Individual, ninds)
		}

		// reset values
		ind.Nwins = 0
		ind.Nlosses = 0
		ind.FrontId = 0
		ind.DistCrowd = 0
		ind.DistNeigh = INF
		fsizes[i] = 0

		// ovas range
		for j := 0; j < nv; j++ {
			x := ind.Ovas[j]
			if math.IsNaN(x) {
				chk.Panic("NaN found in objective value array\n\tx = %v\n\tovas = %v", ind.Floats, ind.Ovas)
			}
			if i == 0 {
				ovamin[j] = x
				ovamax[j] = x
			} else {
				ovamin[j] = utl.Min(ovamin[j], x)
				ovamax[j] = utl.Max(ovamax[j], x)
			}
		}
	}

	// compute neighbour distances and dominance data
	for i := 0; i < ninds; i++ {
		A := pop[i]
		for j := i + 1; j < ninds; j++ {
			B := pop[j]
			dist := IndDistance(A, B, ovamin, ovamax)
			A.DistNeigh = utl.Min(A.DistNeigh, dist)
			B.DistNeigh = utl.Min(B.DistNeigh, dist)
			Adom, Bdom := IndCompare(A, B)
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
			if false {
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

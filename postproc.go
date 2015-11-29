// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/utl"

// GetFeasible returns all feasible solutions
func GetFeasible(sols []*Solution) (feasible []*Solution) {
	for _, sol := range sols {
		if sol.Feasible() {
			feasible = append(feasible, sol)
		}
	}
	return
}

// GetResults returns all ovas and oors
//  Output:
//   ova -- [nsol][nova] objective values
//   oor -- [nsol][noor] out-of-range values
func GetResults(sols []*Solution, ovaOnly bool) (ova, oor [][]float64) {
	nsol := len(sols)
	nova := len(sols[0].Ova)
	noor := len(sols[0].Oor)
	ova = utl.DblsAlloc(nsol, nova)
	if !ovaOnly {
		oor = utl.DblsAlloc(nsol, noor)
	}
	for i, sol := range sols {
		for j := 0; j < nova; j++ {
			ova[i][j] = sol.Ova[j]
		}
		if !ovaOnly {
			for j := 0; j < noor; j++ {
				oor[i][j] = sol.Oor[j]
			}
		}
	}
	return
}

// GetParetoFrontRes returns results on Pareto front
//  Input:
//   p   -- first column in res
//   q   -- second column in res
//   res -- e.g. can be either ova or oor
func GetParetoFrontRes(p, q int, res [][]float64) (fp, fq []float64, front []int) {
	front = utl.ParetoFront(res)
	fp = make([]float64, len(front))
	fq = make([]float64, len(front))
	for i, id := range front {
		fp[i] = res[id][p]
		fq[i] = res[id][q]
	}
	return
}

// GetParetoFront returns Pareto front
func GetParetoFront(p, q int, all []*Solution, feasibleOnly bool) (fp, fq []float64, front []int) {
	var sols []*Solution
	if feasibleOnly {
		sols = GetFeasible(all)
	} else {
		sols = all
	}
	ova, _ := GetResults(sols, true)
	return GetParetoFrontRes(p, q, ova)
}

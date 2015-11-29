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
//   j1  -- first column in res
//   j2  -- second column in res
//   res -- e.g. can be either ova or oor
func GetParetoFrontRes(j1, j2 int, res [][]float64) (x, y []float64) {
	front := utl.ParetoFront(res)
	x = make([]float64, len(front))
	y = make([]float64, len(front))
	for i, id := range front {
		x[i] = res[id][j1]
		y[i] = res[id][j2]
	}
	return
}

// GetParetoFront returns Pareto front
func GetParetoFront(iOva, jOva int, all []*Solution, feasibleOnly bool) (x, y []float64) {
	var sols []*Solution
	if feasibleOnly {
		sols = GetFeasible(all)
	} else {
		sols = all
	}
	ova, _ := GetResults(sols, true)
	return GetParetoFrontRes(iOva, jOva, ova)
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// StatMinProb prints statistical analysis when using MinProb
func (o *Optimiser) StatMinProb(idxF, hlen int, Fref float64, verbose bool) (fmin, fave, fmax, fdev float64, F []float64) {
	if o.MinProb == nil {
		io.Pfred("_warning_ MinProb is <nil>\n")
		return
	}
	nfb := len(o.XfltBest)
	nib := len(o.XintBest)
	if nfb+nib == 0 {
		fmin, fave, fmax, fdev = 1e30, 1e30, 1e30, 1e30
		io.Pfred("_warning_ XfltBest and XintBest are not available. Call RunMany first.\n")
		return
	}
	nbest := utl.Imax(nfb, nib)
	var xf []float64
	var xi []int
	F = make([]float64, nbest)
	cpu := 0
	for i := 0; i < nbest; i++ {
		if nfb > 0 {
			xf = o.XfltBest[i]
		}
		if nib > 0 {
			xi = o.XintBest[i]
		}
		o.MinProb(o.F[cpu], o.G[cpu], o.H[cpu], xf, xi, cpu)
		F[i] = o.F[cpu][idxF]
	}
	if nbest < 2 {
		fmin, fave, fmax = F[0], F[0], F[0]
		return
	}
	fmin, fave, fmax, fdev = rnd.StatBasic(F, true)
	if verbose {
		io.Pf("fmin = %v\n", fmin)
		io.PfYel("fave = %v (%v)\n", fave, Fref)
		io.Pf("fmax = %v\n", fmax)
		io.Pf("fdev = %v\n\n", fdev)
		io.Pf(rnd.BuildTextHist(nice_num(fmin-0.05, 2), nice_num(fmax+0.05, 2), 11, F, "%.2f", hlen))
	}
	return
}

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

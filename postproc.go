// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// GetFeasible returns all feasible solutions
func GetFeasible(sols []*Solution) (feasible []*Solution) {
	for _, sol := range sols {
		if sol.Feasible() {
			feasible = append(feasible, sol)
		}
	}
	return
}

// GetBestFeasible returns the best and list of feasible candidates
// Note: feasible array is sorted by iOva
func GetBestFeasible(opt *Optimiser, iOvaSort int) (best *Solution, feasible []*Solution) {
	feasible = GetFeasible(opt.Solutions)
	if len(feasible) == 0 {
		return
	}
	SortSolutions(feasible, iOvaSort)
	best = feasible[0]
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
	ova = utl.Alloc(nsol, nova)
	if !ovaOnly {
		oor = utl.Alloc(nsol, noor)
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

// CheckFront0 returns front0 and number of failed/success
func CheckFront0(opt *Optimiser, verbose bool) (nfailed int, front0 []*Solution) {
	front0 = make([]*Solution, 0)
	var nsuccess int
	for _, sol := range opt.Solutions {
		var failed bool
		for _, oor := range sol.Oor {
			if oor > 0 {
				failed = true
				break
			}
		}
		if failed {
			nfailed++
		} else {
			nsuccess++
			if sol.FrontId == 0 {
				front0 = append(front0, sol)
			}
		}
	}
	if verbose {
		if nfailed > 0 {
			io.PfRed("N failed = %d out of %d\n", nfailed, opt.Nsol)
		} else {
			io.PfGreen("N success = %d out of %d\n", nsuccess, opt.Nsol)
		}
		io.PfYel("N front 0 = %d\n", len(front0))
	}
	return
}

// FormatXFGH formats x, f, g and h to a string
// Note: (1) only CPU=0 must call this function
//       (2) use fmtX=="" to prevent output of X values
func FormatXFGH(opt *Optimiser, fmtX, fmtFGH string) (l string) {
	if opt.MinProb == nil {
		chk.Panic("FormatXFGH needs the definition of a MinProb")
	}
	fmtFGHok := fmtFGH + " "
	fmtFGHwrong := fmtFGH + "!"
	cpu := 0
	var infeasible bool
	for _, sol := range opt.Solutions {
		if fmtX != "" {
			for j := 0; j < opt.Nflt; j++ {
				l += io.Sf(fmtX, sol.Flt[j])
			}
		}
		infeasible = false
		opt.MinProb(opt.F[cpu], opt.G[cpu], opt.H[cpu], sol.Flt, sol.Int, cpu)
		for j := 0; j < opt.Nf; j++ {
			l += io.Sf(fmtFGH, opt.F[cpu][j])
		}
		for j := 0; j < opt.Ng; j++ {
			if opt.G[cpu][j] < 0 {
				l += io.Sf(fmtFGHwrong, opt.G[cpu][j])
				infeasible = true
			} else {
				l += io.Sf(fmtFGHok, opt.G[cpu][j])
			}
		}
		for j := 0; j < opt.Nh; j++ {
			if math.Abs(opt.H[cpu][j]) > opt.EpsH {
				l += io.Sf(fmtFGHwrong, opt.H[cpu][j])
				infeasible = true
			} else {
				l += io.Sf(fmtFGHok, opt.H[cpu][j])
			}
		}
		if infeasible {
			l += io.Sf(" !!\n")
		} else {
			l += io.Sf(" ok\n")
		}
	}
	return
}

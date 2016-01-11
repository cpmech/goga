// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	gotime "time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Optimiser solves optimisation problems:
//  Solve:
//   min  {Ova[0](x), Ova[1](x), ...} objective values
//    x
//   s.t. Oor[0](x) = 0
//        Oor[1](x) = 0  out-of-range values
//
// A specialised version is also available
//  Solve:
//    min  {f0(x), f1(x), f2(x), ...}  nf functions
//     x   g0(x) ≥ 0
//         g1(x) ≥ 0  ng inequalities
//    s.t. h0(x) = 0
//         h1(x) = 0  nh equalities
//
//  x = xFlt  or  x = xInt   or  x = {xFlt, Xint}
//
type Optimiser struct {

	// input
	Parameters           // input parameters
	ObjFunc    ObjFunc_t // objective function
	MinProb    MinProb_t // minimisation problem function
	CxInt      CxInt_t   // crossover function for ints
	MtInt      MtInt_t   // mutation function for ints

	// essential
	Generator Generator_t // generate solutions
	Solutions []*Solution // current solutions
	Groups    []*Group    // [cpu] competitors per CPU. pointers to current and future solutions
	Metrics   *Metrics    // metrics

	// auxiliary
	Stat                   // structure holding stat data
	Nf, Ng, Nh int         // number of f, g, h functions
	F, G, H    [][]float64 // [cpu] temporary
	tmp        *Solution   // temporary solution
	cpupairs   [][]int     // pairs of CPU ids. for exchanging solutions
	iova0      int         // index of current item in ova[0]
	ova0       []float64   // last ova[0] values to assess convergence
}

// Initialises continues initialisation by generating individuals
//  Optional: fcn XOR obj, nf, ng, nh
func (o *Optimiser) Init(gen Generator_t, obj ObjFunc_t, fcn MinProb_t, nf, ng, nh int) {

	// generic or minimisation problem
	if obj != nil {
		o.ObjFunc = obj
	} else {
		if fcn == nil {
			chk.Panic("either ObjFunc or MinProb must be provided")
		}
		o.Nf, o.Ng, o.Nh, o.MinProb = nf, ng, nh, fcn
		o.ObjFunc = func(sol *Solution, cpu int) {
			o.MinProb(o.F[cpu], o.G[cpu], o.H[cpu], sol.Flt, sol.Int, cpu)
			for i, f := range o.F[cpu] {
				sol.Ova[i] = f
			}
			for i, g := range o.G[cpu] {
				sol.Oor[i] = utl.GtePenalty(g, 0.0, 1) // g[i] ≥ 0
			}
			for i, h := range o.H[cpu] {
				h = math.Abs(h)
				sol.Ova[0] += h
				sol.Oor[o.Ng+i] = utl.GtePenalty(o.EpsH, h, 1) // ϵ ≥ |h[i]|
			}
		}
		o.F = la.MatAlloc(o.Ncpu, o.Nf)
		o.G = la.MatAlloc(o.Ncpu, o.Ng)
		o.H = la.MatAlloc(o.Ncpu, o.Nh)
		o.Nova = o.Nf
		o.Noor = o.Ng + o.Nh
	}

	// calc derived parameters
	o.Generator = gen
	o.CalcDerived()

	// allocate solutions
	o.Solutions = NewSolutions(o.Nsol, &o.Parameters)
	o.Groups = make([]*Group, o.Ncpu)
	for cpu := 0; cpu < o.Ncpu; cpu++ {
		o.Groups[cpu] = new(Group)
		o.Groups[cpu].Init(cpu, o.Ncpu, o.Solutions, &o.Parameters)
	}

	// metrics
	o.Metrics = new(Metrics)
	o.Metrics.Init(o.Nsol, &o.Parameters)

	// auxiliary
	o.tmp = NewSolution(0, 0, &o.Parameters)
	o.cpupairs = utl.IntsAlloc(o.Ncpu/2, 2)
	o.iova0 = -1
	o.ova0 = make([]float64, o.Tf)

	// generate trial solutions
	o.generate_solutions(0)
}

// GetSolutionsCopy returns a copy of Solutions
func (o *Optimiser) GetSolutionsCopy() (res []*Solution) {
	res = NewSolutions(len(o.Solutions), &o.Parameters)
	for i, sol := range o.Solutions {
		sol.CopyInto(res[i])
	}
	return
}

// Solve solves optimisation problem
func (o *Optimiser) Solve() {

	// benchmark
	if o.Verbose {
		t0 := gotime.Now()
		defer func() {
			io.Pf("\nnfeval = %d\n", o.Nfeval)
			io.Pfblue2("cpu time = %v\n", gotime.Now().Sub(t0))
		}()
	}

	// perform evolution
	done := make(chan int, o.Ncpu)
	time := 0
	texc := time + o.DtExc
	for time < o.Tf {

		// run groups in parallel. up to exchange time
		for icpu := 0; icpu < o.Ncpu; icpu++ {
			go func(cpu int) {
				nfeval := 0
				for t := time; t < texc; t++ {
					if cpu == 0 && o.Verbose {
						io.Pf("time = %10d\r", t+1)
					}
					nfeval += o.EvolveOneGroup(cpu)
				}
				done <- nfeval
			}(icpu)
		}
		for cpu := 0; cpu < o.Ncpu; cpu++ {
			o.Nfeval += <-done
		}

		// compute metrics with all solutions included
		o.Metrics.Compute(o.Solutions)

		// exchange via tournament
		if o.Ncpu > 1 {
			if o.ExcTour {
				for i := 0; i < o.Ncpu; i++ {
					j := (i + 1) % o.Ncpu
					I := rnd.IntGetUnique(o.Groups[i].Indices, 2)
					J := rnd.IntGetUnique(o.Groups[j].Indices, 2)
					A, B := o.Groups[i].All[I[0]], o.Groups[i].All[I[1]]
					a, b := o.Groups[j].All[J[0]], o.Groups[j].All[J[1]]
					o.Tournament(A, B, a, b, o.Metrics)
				}
			}

			// exchange one randomly
			if o.ExcOne {
				rnd.IntGetGroups(o.cpupairs, utl.IntRange(o.Ncpu))
				for _, pair := range o.cpupairs {
					i, j := pair[0], pair[1]
					n := utl.Imin(o.Groups[i].Ncur, o.Groups[j].Ncur)
					k := rnd.Int(0, n)
					A := o.Groups[i].All[k]
					B := o.Groups[j].All[k]
					B.CopyInto(o.tmp)
					A.CopyInto(B)
					o.tmp.CopyInto(A)
				}
			}
		}

		// update time variables
		time += o.DtExc
		texc += o.DtExc
		time = utl.Imin(time, o.Tf)
		texc = utl.Imin(texc, o.Tf)
	}
}

// EvolveOneGroup evolves one group (CPU)
func (o *Optimiser) EvolveOneGroup(cpu int) (nfeval int) {

	// auxiliary
	G := o.Groups[cpu].All // competitors (old and new)
	I := o.Groups[cpu].Indices
	P := o.Groups[cpu].Pairs

	// compute random pairs
	rnd.IntGetGroups(P, I)
	np := len(P)

	// create new solutions
	z := o.Groups[cpu].Ncur // index of first new solution
	for k := 0; k < np; k++ {
		l := (k + 1) % np
		m := (k + 2) % np
		n := (k + 3) % np

		A := G[P[k][0]]
		A0 := G[P[l][0]]
		A1 := G[P[m][0]]
		A2 := G[P[n][0]]

		B := G[P[k][1]]
		B0 := G[P[l][1]]
		B1 := G[P[m][1]]
		B2 := G[P[n][1]]

		a := G[z+P[k][0]]
		b := G[z+P[k][1]]

		if o.Nflt > 0 {
			DiffEvol(a.Flt, A.Flt, A0.Flt, A1.Flt, A2.Flt, &o.Parameters)
			DiffEvol(b.Flt, B.Flt, B0.Flt, B1.Flt, B2.Flt, &o.Parameters)
		}

		if o.Nint > 0 {
			o.CxInt(a.Int, b.Int, A.Int, B.Int, &o.Parameters)
			o.MtInt(a.Int, &o.Parameters)
			o.MtInt(b.Int, &o.Parameters)
		}

		if o.BinInt > 0 && o.ClearFlt {
			for i := 0; i < o.Nint; i++ {
				if a.Int[i] == 0 {
					a.Flt[i] = 0
				}
				if b.Int[i] == 0 {
					b.Flt[i] = 0
				}
			}
		}

		o.ObjFunc(a, cpu)
		o.ObjFunc(b, cpu)
		nfeval += 2
	}

	// metrics
	o.Groups[cpu].Metrics.Compute(G)

	// tournaments
	for k := 0; k < np; k++ {
		A := G[P[k][0]]
		B := G[P[k][1]]
		a := G[z+P[k][0]]
		b := G[z+P[k][1]]
		o.Tournament(A, B, a, b, o.Groups[cpu].Metrics)
	}
	return
}

// Tournament performs the tournament among 4 individuals
func (o *Optimiser) Tournament(A, B, a, b *Solution, m *Metrics) {
	dAa := A.Distance(a, m.Fmin, m.Fmax, m.Imin, m.Imax)
	dAb := A.Distance(b, m.Fmin, m.Fmax, m.Imin, m.Imax)
	dBa := B.Distance(a, m.Fmin, m.Fmax, m.Imin, m.Imax)
	dBb := B.Distance(b, m.Fmin, m.Fmax, m.Imin, m.Imax)
	if dAa+dBb < dAb+dBa {
		if !A.Fight(a) {
			a.CopyInto(A)
		}
		if !B.Fight(b) {
			b.CopyInto(B)
		}
		return
	}
	if !A.Fight(b) {
		b.CopyInto(A)
	}
	if !B.Fight(a) {
		a.CopyInto(B)
	}
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// generate_solutions generate solutions
func (o *Optimiser) generate_solutions(itrial int) {

	// benchmark
	if o.Verbose && itrial == 0 {
		t0 := gotime.Now()
		defer func() {
			io.Pfblue2("trial solutions generated in %v\n", gotime.Now().Sub(t0))
		}()
	}

	// generate
	if o.GenAll {
		o.Generator(o.Solutions, &o.Parameters)
		for _, sol := range o.Solutions {
			o.ObjFunc(sol, 0)
		}
	} else {
		done := make(chan int, o.Ncpu)
		for icpu := 0; icpu < o.Ncpu; icpu++ {
			go func(cpu int) {
				start, endp1 := (cpu*o.Nsol)/o.Ncpu, ((cpu+1)*o.Nsol)/o.Ncpu
				sols := o.Solutions[start:endp1]
				o.Generator(sols, &o.Parameters)
				for _, sol := range sols {
					o.ObjFunc(sol, cpu)
				}
				done <- 1
			}(icpu)
		}
		for cpu := 0; cpu < o.Ncpu; cpu++ {
			<-done
		}
	}

	// metrics
	o.iova0 = -1
	o.Nfeval = o.Nsol
	o.Metrics.Compute(o.Solutions)
}

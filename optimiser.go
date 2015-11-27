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
	CxFlt      CxFlt_t   // crossover function for floats
	CxInt      CxInt_t   // crossover function for ints
	MtFlt      MtFlt_t   // mutation function for floats
	MtInt      MtInt_t   // mutation function for ints

	// essential
	Solutions   []*Solution    // current solutions
	FutureSols  []*Solution    // future solutions
	Competitors []*Competitors // [cpu] competitors per CPU. pointers to current and future solutions

	// auxiliary
	Nf, Ng, Nh int         // number of f, g, h functions
	F, G, H    [][]float64 // [cpu] temporary

	// stat
	Nfeval int // number of function evaluations
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
				sol.Oor[o.Ng+i] = utl.GtePenalty(o.EpsMinProb, h, 1) // ϵ ≥ |h[i]|
			}
		}
		o.F = la.MatAlloc(o.Ncpu, o.Nf)
		o.G = la.MatAlloc(o.Ncpu, o.Ng)
		o.H = la.MatAlloc(o.Ncpu, o.Nh)
		o.Nova = o.Nf
		o.Noor = o.Ng + o.Nh
	}

	// calc derived parameters
	o.CalcDerived()

	// crossover and mutation functions
	if o.CxFlt == nil {
		o.CxFlt = CxFltDE
	}
	if o.MtFlt == nil {
		o.MtFlt = MtFltDeb
	}

	// essential
	o.Solutions = NewSolutions(o.Nsol, &o.Parameters)
	o.FutureSols = NewSolutions(o.Nsol, &o.Parameters)
	o.Competitors = make([][]*Solution, o.Ncpu)
	o.Indices = make([][]int, o.Ncpu)
	o.Pairs = make([][][]int, o.Ncpu)
	o.Metrics = make([]*Metrics, o.Ncpu)
	for cpu := 0; cpu < o.Ncpu; cpu++ {
		start, endp1 := (cpu*o.Nsol)/o.Ncpu, ((cpu+1)*o.Nsol)/o.Ncpu
		nsol := endp1 - start
		o.Competitors[cpu] = make([]*Solution, nsol*2)
		o.Indices[cpu] = make([]int, nsol)
		o.Pairs[cpu] = utl.IntsAlloc(nsol/2, 2)
		for i := 0; i < nsol; i++ {
			o.Competitors[cpu][i] = o.Solutions[start+i]
			o.Competitors[cpu][nsol+i] = o.FutureSols[start+i]
			o.Indices[cpu][i] = i
		}
		o.Metrics[cpu] = new(Metrics)
		o.Metrics[cpu].Init(o.Competitors[cpu])
	}

	// generate trial solutions
	t0 := gotime.Now()
	done := make(chan int, o.Ncpu)
	for icpu := 0; icpu < o.Ncpu; icpu++ {
		go func(cpu int) {
			nsol := len(o.Indices[cpu])
			sols := o.Competitors[cpu][:nsol]
			gen(sols, &o.Parameters)
			for _, sol := range sols {
				o.ObjFunc(sol, cpu)
			}
			done <- 1
		}(icpu)
	}
	for cpu := 0; cpu < o.Ncpu; cpu++ {
		<-done
	}
	o.Nfeval = o.Nsol
	if o.Verbose {
		io.Pf(". . . trial solutions generated in %v . . .\n", gotime.Now().Sub(t0))
	}
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
	done := make(chan int, o.Ncpu)
	time := 0
	texc := time + o.DtExc
	for time < o.Tf {
		if o.Verbose {
			io.Pf("time = %10d\r", time)
		}
		for icpu := 0; icpu < o.Ncpu; icpu++ {
			go func(cpu int) {
				for t := time; t < texc; t++ {
					if cpu == 0 && o.Verbose {
						io.Pf("time = %10d\r", t)
					}
					o.evolve(cpu)
				}
				done <- 1
			}(icpu)
		}
		for cpu := 0; cpu < o.Ncpu; cpu++ {
			<-done
		}
		time += o.DtExc
		texc += o.DtExc
		time = utl.Imin(time, o.Tf)
		texc = utl.Imin(texc, o.Tf)
		if o.Verbose {
			io.Pf("\n")
		}
	}
	if o.Verbose {
		io.PfWhite("time = %10d\n", time)
	}
}

func (o *Optimiser) evolve(cpu int) {

	// auxiliary
	competitors := o.Competitors[cpu]
	indices := o.Indices[cpu]
	pairs := o.Pairs[cpu]

	// compute random pairs
	rnd.IntGetGroups(pairs, indices)

	// create new solutions
	nsol := len(indices)
	idx := nsol
	var a, b, A, B, C, D *Solution
	for k, pair := range pairs {
		l := (k + 1) % len(pairs)
		A = competitors[pair[0]]
		B = competitors[pair[1]]
		C = competitors[pairs[l][0]]
		D = competitors[pairs[l][1]]
		a = competitors[idx]
		b = competitors[idx+1]
		idx += 2
		o.crossover(a, b, A, B, C, D)
		o.mutation(a)
		o.mutation(b)
		o.ObjFunc(a, 0)
		o.ObjFunc(b, 0)
		o.Nfeval += 2
	}

	// metrics
	o.metrics(competitors)

	// tournaments
	idx = nsol
	for _, pair := range pairs {
		A = competitors[pair[0]]
		B = competitors[pair[1]]
		a = competitors[idx]
		b = competitors[idx+1]
		idx += 2
		dAa := A.Distance(a, o.Fmin, o.Fmax, o.Imin, o.Imax)
		dAb := A.Distance(b, o.Fmin, o.Fmax, o.Imin, o.Imax)
		dBa := B.Distance(a, o.Fmin, o.Fmax, o.Imin, o.Imax)
		dBb := B.Distance(b, o.Fmin, o.Fmax, o.Imin, o.Imax)
		if dAa+dBb < dAb+dBa {
			if a.Fight(A) {
				a.CopyInto(A)
			}
			if b.Fight(B) {
				b.CopyInto(B)
			}
		} else {
			if b.Fight(A) {
				b.CopyInto(A)
			}
			if a.Fight(B) {
				a.CopyInto(B)
			}
		}
	}
}

func (o *Optimiser) crossover(a, b, A, B, C, D *Solution) {
	if o.Nflt > 0 {
		o.CxFlt(a.Flt, b.Flt, A.Flt, B.Flt, C.Flt, D.Flt, &o.Parameters)
	}
	if o.Nint > 0 {
		o.CxInt(a.Int, b.Int, A.Int, B.Int, C.Int, D.Int, &o.Parameters)
	}
}

func (o *Optimiser) mutation(a *Solution) {
	if o.Nflt > 0 && o.PmFlt > 0 {
		o.MtFlt(a.Flt, &o.Parameters)
	}
	if o.Nint > 0 && o.PmInt > 0 {
		o.MtInt(a.Int, &o.Parameters)
	}
}

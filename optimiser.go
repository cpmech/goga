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
	Generator  Generator_t // generate solutions
	Solutions  []*Solution // current solutions
	FutureSols []*Solution // future solutions
	Groups     []*Group    // [cpu] competitors per CPU. pointers to current and future solutions
	Metrics    *Metrics    // metrics

	// auxiliary
	Nf, Ng, Nh int         // number of f, g, h functions
	F, G, H    [][]float64 // [cpu] temporary
	tmp        *Solution   // temporary solution
	cpupairs   [][]int     // pairs of CPU ids. for exchanging solutions

	// stat
	Nfeval   int         // number of function evaluations
	XfltBest [][]float64 // best results after RunMany
	XintBest [][]int     // best results after RunMany
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
	o.Generator = gen
	o.CalcDerived()

	// crossover and mutation functions
	if o.CxFlt == nil {
		o.CxFlt = CxFltDE
	}
	if o.MtFlt == nil {
		o.MtFlt = MtFltDeb
	}

	// allocate solutions
	o.Solutions = NewSolutions(o.Nsol, &o.Parameters)
	o.FutureSols = NewSolutions(o.Nsol, &o.Parameters)
	o.Groups = make([]*Group, o.Ncpu)
	for cpu := 0; cpu < o.Ncpu; cpu++ {
		o.Groups[cpu] = new(Group)
		o.Groups[cpu].Init(cpu, o.Ncpu, o.Solutions, o.FutureSols, &o.Parameters)
	}

	// metrics
	o.Metrics = new(Metrics)
	o.Metrics.Init(o.Nsol, &o.Parameters)

	// auxiliary
	o.tmp = NewSolution(0, 0, &o.Parameters)
	o.cpupairs = utl.IntsAlloc(o.Ncpu/2, 2)

	// generate trial solutions
	o.gensolutions(0)
}

func (o *Optimiser) gensolutions(itrial int) {
	t0 := gotime.Now()
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
	o.Nfeval = o.Nsol
	if o.Verbose && itrial > 0 {
		io.Pf(". . . trial solutions generated in %v . . .\n", gotime.Now().Sub(t0))
	}
	o.Metrics.Compute(o.Solutions)
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
					if o.UseTriples {
						nfeval += o.evolve_with_triples(cpu)
					} else {
						nfeval += o.evolve(cpu)
					}
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
			if o.use_exchange_via_tournament {
				for i := 0; i < o.Ncpu; i++ {
					j := (i + 1) % o.Ncpu
					o.exchange_via_tournament(i, j)
				}
			}

			// exchange one randomly
			if o.use_exchange_one_randomly {
				rnd.IntGetGroups(o.cpupairs, utl.IntRange(o.Ncpu))
				for _, pair := range o.cpupairs {
					o.exchange_one_randomly(pair[0], pair[1])
				}
			}
		}

		// update time variables
		time += o.DtExc
		texc += o.DtExc
		time = utl.Imin(time, o.Tf)
		texc = utl.Imin(texc, o.Tf)
		if o.Verbose {
			io.Pf("\n")
		}
	}

	// message
	if o.Verbose {
		io.Pf("nfeval = %d\n", o.Nfeval)
	}
}

func (o *Optimiser) RunMany() {
	if o.Verbose {
		t0 := gotime.Now()
		defer func() {
			io.Pfblue2("\ncpu time = %v\n", gotime.Now().Sub(t0))
		}()
	}
	for itrial := 0; itrial < o.Ntrials; itrial++ {
		if itrial > 0 {
			o.gensolutions(itrial)
		}
		o.Solve()
		if o.Nova < 2 {
			SortByOva(o.Solutions, 0)
		} else {
			SortByTradeoff(o.Solutions)
		}
		if o.Nflt > 0 {
			o.XfltBest = append(o.XfltBest, o.Solutions[0].Flt)
		}
		if o.Nflt > 0 {
			o.XintBest = append(o.XintBest, o.Solutions[0].Int)
		}
	}
	return
}

// StatMinProb prints statistical analysis when using MinProb
func (o *Optimiser) StatMinProb(idxF, hlen int, Fref float64, verbose bool) (fmin, fave, fmax, fdev float64) {
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
	F := make([]float64, nbest)
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
		io.Pf(rnd.BuildTextHist(nice_num(fmin-0.05), nice_num(fmax+0.05), 11, F, "%.2f", hlen))
	}
	return
}

func (o *Optimiser) exchange_via_tournament(i, j int) {
	selI := rnd.IntGetUnique(o.Groups[i].Indices, 2)
	selJ := rnd.IntGetUnique(o.Groups[j].Indices, 2)
	A, B := o.Groups[i].All[selI[0]], o.Groups[i].All[selI[1]]
	a, b := o.Groups[j].All[selJ[0]], o.Groups[j].All[selJ[1]]
	o.tournament(A, B, a, b, o.Metrics)
}

func (o *Optimiser) exchange_one_randomly(i, j int) {
	n := utl.Imin(o.Groups[i].Ncur, o.Groups[j].Ncur)
	k := rnd.Int(0, n)
	A := o.Groups[i].All[k]
	B := o.Groups[j].All[k]
	B.CopyInto(o.tmp)
	A.CopyInto(B)
	o.tmp.CopyInto(A)
}

// evolve evolves one group
func (o *Optimiser) evolve(cpu int) (nfeval int) {

	// auxiliary
	competitors := o.Groups[cpu].All
	indices := o.Groups[cpu].Indices
	pairs := o.Groups[cpu].Pairs

	// compute random pairs
	rnd.IntGetGroups(pairs, indices)

	// create new solutions
	z := o.Groups[cpu].Ncur
	for k := 0; k < len(pairs); k++ {
		l := (k + 1) % len(pairs)
		m := (k + 2) % len(pairs)
		A := competitors[pairs[k][0]]
		B := competitors[pairs[k][1]]
		C := competitors[pairs[l][0]]
		D := competitors[pairs[l][1]]
		E := competitors[pairs[m][0]]
		F := competitors[pairs[m][1]]
		a := competitors[z+pairs[k][0]]
		b := competitors[z+pairs[k][1]]
		o.crossover(a, b, A, B, C, D, E, F)
		o.mutation(a)
		o.mutation(b)
		o.ObjFunc(a, cpu)
		o.ObjFunc(b, cpu)
		nfeval += 2
	}

	// metrics
	o.Groups[cpu].Metrics.Compute(competitors)

	// tournaments
	for k := 0; k < len(pairs); k++ {
		A := competitors[pairs[k][0]]
		B := competitors[pairs[k][1]]
		a := competitors[z+pairs[k][0]]
		b := competitors[z+pairs[k][1]]
		o.tournament(A, B, a, b, o.Groups[cpu].Metrics)
	}
	return
}

// evolve_with_triples evolves one group with triples of solutions
func (o *Optimiser) evolve_with_triples(cpu int) (nfeval int) {

	// auxiliary
	competitors := o.Groups[cpu].All
	indices := o.Groups[cpu].Indices
	triples := o.Groups[cpu].Triples
	mdist := o.Groups[cpu].Mdist
	match := o.Groups[cpu].Match

	// compute random triples
	rnd.IntGetGroups(triples, indices)

	// create new solutions
	z := o.Groups[cpu].Ncur
	news := make([]*Solution, 3)
	main := make([]*Solution, 3)
	auxi := make([]*Solution, 3)
	for k := 0; k < len(triples); k++ {
		l := (k + 1) % len(triples)
		for i := 0; i < 3; i++ {
			news[i] = competitors[triples[k][i]+z]
			main[i] = competitors[triples[k][i]]
			auxi[i] = competitors[triples[l][i]]
		}
		CxFltDE_triple(news[0].Flt, news[1].Flt, news[2].Flt, main[0].Flt, main[1].Flt, main[2].Flt, auxi[0].Flt, auxi[1].Flt, auxi[2].Flt, &o.Parameters)
		for i := 0; i < 3; i++ {
			o.ObjFunc(news[i], cpu)
		}
		nfeval += 3
	}

	// metrics
	o.Groups[cpu].Metrics.Compute(competitors)

	// tournaments
	m := o.Groups[cpu].Metrics
	for k := 0; k < len(triples); k++ {
		for i := 0; i < 3; i++ {
			news[i] = competitors[triples[k][i]+z]
			main[i] = competitors[triples[k][i]]
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				mdist[i][j] = main[i].Distance(news[j], m.Fmin, m.Fmax, m.Imin, m.Imax)
			}
		}
		match.SetCostMatrix(mdist)
		match.Run()
		for i := 0; i < 3; i++ {
			j := match.Links[i]
			if !main[i].Fight(news[j]) {
				news[j].CopyInto(main[i])
			}
		}
	}
	return
}

// crossover performs crossover in A,B,C,D,E,F to obtain a and b
func (o *Optimiser) crossover(a, b, A, B, C, D, E, F *Solution) {
	if o.Nflt > 0 {
		o.CxFlt(a.Flt, b.Flt, A.Flt, B.Flt, C.Flt, D.Flt, E.Flt, F.Flt, &o.Parameters)
	}
	if o.Nint > 0 {
		o.CxInt(a.Int, b.Int, A.Int, B.Int, C.Int, D.Int, E.Int, F.Int, &o.Parameters)
	}
}

// mutation performs mutation in a
func (o *Optimiser) mutation(a *Solution) {
	if o.Nflt > 0 && o.PmFlt > 0 {
		o.MtFlt(a.Flt, &o.Parameters)
	}
	if o.Nint > 0 && o.PmInt > 0 {
		o.MtInt(a.Int, &o.Parameters)
	}
}

// tournament performs the tournament among 4 individuals
func (o *Optimiser) tournament(A, B, a, b *Solution, m *Metrics) {
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
	} else {
		if !A.Fight(b) {
			b.CopyInto(A)
		}
		if !B.Fight(a) {
			a.CopyInto(B)
		}
	}
}

// nice_num returns a truncated float
func nice_num(x float64) float64 {
	s := io.Sf("%.2f", x)
	return io.Atof(s)
}

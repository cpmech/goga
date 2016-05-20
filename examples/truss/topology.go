// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"strings"
	"time"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

// main function
func main() {

	// flags
	benchmark := false
	ncpuMax := 16

	// benchmarking
	if benchmark {
		var nsol, tf int
		var et time.Duration
		X := make([]float64, ncpuMax)
		T := make([]float64, ncpuMax)
		S := make([]float64, ncpuMax) // speedup
		S[0] = 1
		for i := 0; i < ncpuMax; i++ {
			io.Pf("\n\n")
			nsol, tf, et = runone(i + 1)
			io.PfYel("elaspsedTime = %v\n", et)
			X[i] = float64(i + 1)
			T[i] = et.Seconds()
			if i > 0 {
				S[i] = T[0] / T[i] // Told / Tnew
			}
		}

		plt.SetForEps(0.75, 250)
		plt.Plot(X, S, io.Sf("'b-',marker='.', label='speedup: $N_{sol}=%d,\\,t_f=%d$', clip_on=0, zorder=100", nsol, tf))
		plt.Plot([]float64{1, 16}, []float64{1, 16}, "'k--',zorder=50")
		plt.Gll("$N_{cpu}:\\;$ number of groups", "speedup", "leg_out=1")
		plt.DoubleYscale("$T_{sys}:\\;$ system time [s]")
		plt.Plot(X, T, "'k-',color='gray', clip_on=0")
		plt.SaveD("/tmp/goga", "topology-speedup.eps")
		return
	}

	// normal run
	runone(-1)
}

func runone(ncpu int) (nsol, tf int, elaspsedTime time.Duration) {

	// input filename
	fn, fnkey := io.ArgToFilename(0, "ground10", ".sim", true)

	// GA parameters
	var opt goga.Optimiser
	opt.Read("ga-" + fnkey + ".json")
	opt.GenType = "rnd"
	nsol, tf = opt.Nsol, opt.Tf
	postproc := true
	if ncpu > 0 {
		opt.Ncpu = ncpu
		postproc = false
	}

	// FEM
	data := make([]*FemData, opt.Ncpu)
	for i := 0; i < opt.Ncpu; i++ {
		data[i] = NewData(fn, fnkey, i)
	}
	io.Pforan("MaxWeight = %v\n", data[0].MaxWeight)

	// set integers
	if data[0].Opt.BinInt {
		opt.CxInt = goga.CxInt
		opt.MtInt = goga.MtIntBin
		opt.BinInt = data[0].Ncells
	}

	// set floats
	opt.FltMin = make([]float64, data[0].Nareas)
	opt.FltMax = make([]float64, data[0].Nareas)
	for i := 0; i < data[0].Nareas; i++ {
		opt.FltMin[i] = data[0].Opt.Amin
		opt.FltMax[i] = data[0].Opt.Amax
	}

	// initialise optimiser
	opt.Nova = 2 // weight and deflection
	opt.Noor = 4 // mobility, feasibility, maxdeflection, stress
	opt.Init(goga.GenTrialSolutions, func(sol *goga.Solution, cpu int) {
		mob, fail, weight, umax, _, errU, errS := data[cpu].RunFEM(sol.Int, sol.Flt, 0, false)
		sol.Ova[0] = weight
		sol.Ova[1] = umax
		sol.Oor[0] = mob
		sol.Oor[1] = fail
		sol.Oor[2] = errU
		sol.Oor[3] = errS
	}, nil, 0, 0, 0)

	// initial solutions
	var sols0 []*goga.Solution
	if false {
		sols0 = opt.GetSolutionsCopy()
	}

	// benchmark
	initialTime := time.Now()
	defer func() {
		elaspsedTime = time.Now().Sub(initialTime)
	}()

	// solve
	opt.Verbose = true
	opt.Solve()
	goga.SortByOva(opt.Solutions, 0)

	// post processing
	if !postproc {
		return
	}

	// check
	nfailed, front0 := goga.CheckFront0(&opt, true)

	// save results
	var log, res bytes.Buffer
	io.Ff(&log, opt.LogParams())
	io.Ff(&res, PrintSolutions(data[0], opt.Solutions))
	io.Ff(&res, io.Sf("\n\nnfailed = %d\n", nfailed))
	io.WriteFileVD("/tmp/goga", fnkey+".log", &log)
	io.WriteFileVD("/tmp/goga", fnkey+".res", &res)

	// plot Pareto-optimal front
	feasibleOnly := true
	plt.SetForEps(0.8, 350)
	if strings.HasPrefix(fnkey, "ground10") {
		_, ref, _ := io.ReadTable("p460_fig300.dat")
		plt.Plot(ref["w"], ref["u"], "'b-', label='reference'")
	}
	fmtAll := &plt.Fmt{L: "final solutions", M: ".", C: "orange", Ls: "none", Ms: 3}
	fmtFront := &plt.Fmt{L: "final Pareto front", C: "r", M: "o", Ms: 3, Ls: "none"}
	goga.PlotOvaOvaPareto(&opt, sols0, 0, 1, feasibleOnly, fmtAll, fmtFront)
	plt.Gll("weight ($f_0$)", "deflection ($f_1)$", "") //, "leg_out=1, leg_ncol=4, leg_hlen=1.5")
	if strings.HasPrefix(fnkey, "ground10") {
		plt.AxisRange(1800, 14000, 1, 6)
	}

	// plot selected results
	ia, ib, ic, id, ie := 0, 0, 0, 0, 0
	nfront0 := len(front0)
	io.Pforan("nfront0 = %v\n", nfront0)
	if nfront0 > 4 {
		ib = nfront0 / 10
		ic = nfront0 / 5
		id = nfront0 / 2
		ie = nfront0 - 1
	}
	A := front0[ia]
	B := front0[ib]
	C := front0[ic]
	D := front0[id]
	E := front0[ie]
	wid, hei := 0.20, 0.10
	draw_truss(data[0], "A", A, 0.17, 0.75, wid, hei)
	draw_truss(data[0], "B", B, 0.20, 0.55, wid, hei)
	draw_truss(data[0], "C", C, 0.28, 0.33, wid, hei)
	draw_truss(data[0], "D", D, 0.47, 0.22, wid, hei)
	draw_truss(data[0], "E", E, 0.70, 0.18, wid, hei)

	// save figure
	plt.SaveD("/tmp/goga", fnkey+".eps")

	// tex file
	title := "Shape and topology optimisation. Results."
	label := "topoFront"
	document := true
	compact := true
	tex_results("/tmp/goga", "tmp_"+fnkey, title, label, data[0], A, B, C, D, E, document, compact)
	document = false
	tex_results("/tmp/goga", fnkey, title, label, data[0], A, B, C, D, E, document, compact)
	return
}

type FltFormatter []float64

func (o FltFormatter) String() (l string) {
	for _, val := range o {
		if val < 1e-9 {
			l += "       "
		} else {
			l += io.Sf("%7.2f", val)
		}
	}
	return l
}

func PrintSolutions(fed *FemData, sols []*goga.Solution) (l string) {
	goga.SortByOva(sols, 0)
	l = io.Sf("%8s%6s%6s |%s\n", "weight", "umax", "smax", "areas")
	for _, sol := range sols {
		mob, fail, weight, umax, smax, errU, errS := fed.RunFEM(sol.Int, sol.Flt, 0, false)
		if mob > 0 || fail > 0 || errU > 0 || errS > 0 {
			l += io.Sf("%20s |%s\n", "unfeasible    ", FltFormatter(sol.Flt))
			continue
		}
		l += io.Sf("%8.1f%6.2f%6.2f |%s\n", weight, umax, smax, FltFormatter(sol.Flt))
	}
	return
}

func draw_truss(dat *FemData, key string, A *goga.Solution, lef, bot, wid, hei float64) (weight, deflection float64) {
	gap := 0.1
	plt.PyCmds(io.Sf(`
from pylab import axes, setp, sca
ax_current = gca()
ax_new = axes([%g, %g, %g, %g], axisbg='#dcdcdc')
setp(ax_new, xticks=[0,720], yticks=[0,360])
axis('equal')
axis('off')
`, lef, bot, wid, hei))
	_, _, weight, deflection, _, _, _ = dat.RunFEM(A.Int, A.Flt, 1, false)
	plt.PyCmds("sca(ax_current)\n")
	plt.PlotOne(weight, deflection, "'g*', zorder=1000, clip_on=0")
	plt.Text(weight, deflection+gap, key, "")
	return
}

func tex_results(dirout, fnkey, title, label string, dat *FemData, A, B, C, D, E *goga.Solution, document, compact bool) {
	if len(A.Flt) != 10 {
		chk.Panic("tex_results works with len(Areas)==10 only\n")
		return
	}
	buf := new(bytes.Buffer)
	if document {
		io.Ff(buf, `\documentclass[a4paper]{article}

\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{booktabs}
\usepackage[margin=1.5cm,footskip=0.5cm]{geometry}

\title{GOGA Report}
\author{Dorival Pedroso}

\begin{document}

`)
	}
	io.Ff(buf, `\begin{table} \centering
\caption{%s}
`, title)
	if compact {
		io.Ff(buf, `\begin{tabular}[c]{cccccccc} \toprule
point & weight & deflection &  $A_0$ & $A_1$ & $A_2$ & $A_3$ & $A_4$   \\
      &        &            &  $A_5$ & $A_6$ & $A_7$ & $A_8$ & $A_9$   \\ \hline
`)
	} else {
		io.Ff(buf, `\begin{tabular}[c]{ccccccccccccc} \toprule
point & weight & deflection & $A_0$ & $A_1$ & $A_2$ & $A_3$ & $A_4$ & $A_5$ & $A_6$ & $A_7$ & $A_8$ & $A_9$ \\ \hline
`)
	}

	writeline := func(pt string, E []int, A []float64) {
		_, _, weight, deflection, _, _, _ := dat.RunFEM(E, A, 0, false)
		if compact {
			io.Ff(buf, "%s & $%.2f$ & $%.6f$ &  $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ \\\\\n", pt, weight, deflection, A[0], A[1], A[2], A[3], A[4])
			io.Ff(buf, "   &        &        &  $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ \\\\\n", A[5], A[6], A[7], A[8], A[9])
		} else {
			io.Ff(buf, "%s & $%.2f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ \\\\\n", pt, weight, deflection, A[0], A[1], A[2], A[3], A[4], A[5], A[6], A[7], A[8], A[9])
		}
	}

	writeline("A", A.Int, A.Flt)
	writeline("B", B.Int, B.Flt)
	writeline("C", C.Int, C.Flt)
	writeline("D", D.Int, D.Flt)
	writeline("E", E.Int, E.Flt)

	io.Ff(buf, `
\bottomrule
\end{tabular}
\label{tab:%s}
\end{table}
`, label)
	if document {
		io.Ff(buf, ` \end{document}`)
	}

	tex := fnkey + ".tex"
	if document {
		io.WriteFileD(dirout, tex, buf)
		_, err := io.RunCmd(true, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory=/tmp/goga/", tex)
		if err != nil {
			chk.Panic("%v", err)
		}
		io.PfBlue("file <%s/%s.pdf> generated\n", dirout, fnkey)
	} else {
		io.WriteFileVD(dirout, tex, buf)
	}
}

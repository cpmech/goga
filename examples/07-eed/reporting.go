// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

var selectedKeys = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

func PrintResults(sys *System, selected []*goga.Solution, Pref []float64, costRef, emisRef float64) (l string) {
	n := 9*2 + 8*6 + 12 + 3
	l += io.Sf("%s", io.StrThickLine(n))
	l += io.Sf("%3s%9s%9s%8s%8s%8s%8s%8s%8s%12s\n", "pt", "cost", "emis", "P1", "P2", "P3", "P4", "P5", "P6", "bal.err")
	l += io.Sf("%s", io.StrThinLine(n))
	if Pref != nil {
		l += io.Sf("%3s%9.4f%9.5f%8.5f%8.5f%8.5f%8.5f%8.5f%8.5f%12.4e\n", "ref", costRef, emisRef, Pref[0], Pref[1], Pref[2], Pref[3], Pref[4], Pref[5], sys.Balance(Pref))
	}
	writeline := func(pt string, P []float64) {
		l += io.Sf("%3s%9.4f%9.5f%8.5f%8.5f%8.5f%8.5f%8.5f%8.5f%12.4e\n", pt, sys.FuelCost(P), sys.Emission(P), P[0], P[1], P[2], P[3], P[4], P[5], sys.Balance(P))
	}
	for i, sel := range selected {
		writeline(selectedKeys[i], sel.Flt)
	}
	l += io.Sf("%s", io.StrThickLine(n))
	io.Pf(l)
	return
}

func texResults(dirout, fnkey, title, label string, sys *System, selected []*goga.Solution, document, compact bool) {
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
		io.Ff(buf, `\begin{tabular}[c]{ccccccc} \toprule
point & cost & emission & $h_0$  &  $P_0$ & $P_1$ & $P_2$ \\
      &      &          &        &  $P_3$ & $P_4$ & $P_5$ \\ \hline
`)
	} else {
		io.Ff(buf, `\begin{tabular}[c]{cccccccccc} \toprule
point & cost & emission & $P_0$ & $P_1$ & $P_2$ & $P_3$ & $P_4$ & $P_5$ & $h_0$ \\ \hline
`)
	}

	writeline := func(pt string, P []float64) {
		strh0 := io.TexNum("%.2e", sys.Balance(P), true)
		if compact {
			io.Ff(buf, "%s & $%.4f$ & $%.6f$ & $%s$  &  $%.6f$ & $%.6f$ & $%.6f$ \\\\\n", pt, sys.FuelCost(P), sys.Emission(P), strh0, P[0], P[1], P[2])
			io.Ff(buf, "   &        &        &       &  $%.6f$ & $%.6f$ & $%.6f$ \\\\\n", P[3], P[4], P[5])
		} else {
			io.Ff(buf, "%s & $%.4f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%.6f$ & $%s$ \\\\\n", pt, sys.FuelCost(P), sys.Emission(P), P[0], P[1], P[2], P[3], P[4], P[5], strh0)
		}
	}

	for i, sel := range selected {
		writeline(selectedKeys[i], sel.Flt)
	}

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
		_, err := io.RunCmd(false, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory=/tmp/goga/", tex)
		if err != nil {
			chk.Panic("%v", err)
		}
		io.PfBlue("file <%s/%s.pdf> generated\n", dirout, fnkey)
	} else {
		io.WriteFileVD(dirout, tex, buf)
	}
}

func check(fcn goga.MinProb_t, ng, nh int, xs []float64, fs, ϵ float64) {
	f := make([]float64, 1)
	g := make([]float64, ng)
	h := make([]float64, nh)
	cpu := 0
	fcn(f, g, h, xs, nil, cpu)
	io.Pfblue2("xs = %v\n", xs)
	io.Pfblue2("f(x)=%g  (%g)  diff=%g\n", f[0], fs, math.Abs(fs-f[0]))
	for i, v := range g {
		unfeasible := false
		if v < 0 {
			unfeasible = true
		}
		if unfeasible {
			io.Pfred("g%d(x) = %g\n", i, v)
		} else {
			io.Pfgreen("g%d(x) = %g\n", i, v)
		}
	}
	for i, v := range h {
		unfeasible := false
		if math.Abs(v) > ϵ {
			unfeasible = true
		}
		if unfeasible {
			io.Pfred("h%d(x) = %g\n", i, v)
		} else {
			io.Pfgreen("h%d(x) = %g\n", i, v)
		}
	}
}

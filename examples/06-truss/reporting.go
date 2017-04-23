// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

// FltFormatter helps with formatting float64 numbers
type FltFormatter []float64

// String returns the string representing numbers
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

// PrintSolutions prints solutions
func PrintSolutions(fed *FemData, sols []*goga.Solution) (l string) {
	goga.SortSolutions(sols, 0)
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

// drawTruss draws truss
func drawTruss(dat *FemData, key string, A *goga.Solution, lef, bot, wid, hei float64) (weight, deflection float64) {
	gap := 0.1
	axMain, _ := plt.ZoomWindow(lef, bot, wid, hei, nil)
	_, _, weight, deflection, _, _, _ = dat.RunFEM(A.Int, A.Flt, 1, false)
	plt.Equal()
	plt.AxisRange(0, 720, 0, 360)
	plt.AxisOff()
	plt.PlotOne(weight, deflection, &plt.A{C: "g", M: "*"})
	plt.Sca(axMain)
	plt.Text(weight, deflection+gap, key, &plt.A{C: "blue"})
	return
}

// texResults generates files with results
func texResults(dirout, fnkey, title, label string, dat *FemData, A, B, C, D, E *goga.Solution, document, compact bool) {
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
		_, err := io.RunCmd(false, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory=/tmp/goga/", tex)
		if err != nil {
			chk.Panic("%v", err)
		}
		io.PfBlue("file <%s/%s.pdf> generated\n", dirout, fnkey)
	} else {
		io.WriteFileVD(dirout, tex, buf)
	}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

// TeX document ////////////////////////////////////////////////////////////////////////////////////

// TexDocumentStart starts TeX document
func TexDocumentStart() (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	io.Ff(buf, `\documentclass[a4paper]{article}

\usepackage{mydefaults}
\usepackage[margin=1.5cm,footskip=0.5cm]{geometry}

\title{GOGA Report}
\author{Dorival Pedroso}

\begin{document}
`)
	return
}

// TexDocumentEnd ends TeX document
func TexDocumentEnd(buf *bytes.Buffer) {
	io.Ff(buf, `
\end{document}`)
}

// TexWrite writes and compiles TeX document
func TexWrite(dirout, fnkey string, buf *bytes.Buffer, dorun bool) {
	tex := fnkey + ".tex"
	io.WriteFileVD(dirout, tex, buf)
	if dorun {
		_, err := io.RunCmd(true, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory=/tmp/goga/", tex)
		if err != nil {
			chk.Panic("%v", err)
		}
		io.PfBlue("file <%s/%s.pdf> generated\n", dirout, fnkey)
	}
}

// parameters only table ///////////////////////////////////////////////////////////////////////////

// TexPrmsTableStart starts table with parameters
func TexPrmsTableStart(buf *bytes.Buffer) {
	io.Ff(buf, `
\begin{table} \centering
\caption{goga: Parameters}
\begin{tabular}[c]{cccccc} \toprule
P & $N_{sol}$ & $N_{cpu}$ & $t_{max}$ & $\Delta t_{exc}$ & $N_{eval}$ \\ \hline
`)
}

// TexPrmsTableEnd ends table with parameters
func TexPrmsTableEnd(buf *bytes.Buffer) {
	io.Ff(buf, `\end{tabular}
\label{tab:prms}
\end{table}`)
}

// TexPrmsTableItem adds item to table with parameters
func TexPrmsTableItem(o *Optimiser, buf *bytes.Buffer, problem int) {
	io.Ff(buf, "%d & %d & %d & %d & %d & %d \\\\\n", problem, o.Nsol, o.Ncpu, o.Tf, o.DtExc, o.Nfeval)
}

// TexPrmsReport generates TeX report with parameters
//  nRowPerTab -- number of rows per table
func TexPrmsReport(dirout, fnkey string, opts []*Optimiser, nRowPerTab int) {
	buf := TexDocumentStart()
	for i, opt := range opts {
		if i%nRowPerTab == 0 {
			if i > 0 {
				io.Ff(buf, `\bottomrule`)
				TexPrmsTableEnd(buf) // end previous table
				io.Ff(buf, "\n")
			}
			TexPrmsTableStart(buf) // begin new table
		} else {
			if i > 0 {
				io.Ff(buf, `\hline`)
			}
		}
		TexPrmsTableItem(opt, buf, i+1)
	}
	io.Ff(buf, `\bottomrule`)
	TexPrmsTableEnd(buf) // end previous table
	io.Ff(buf, "\n")
	TexDocumentEnd(buf)
	TexWrite(dirout, fnkey, buf, true)
}

// single objective tables /////////////////////////////////////////////////////////////////////////

// TexSingleObjTableStart starts table for single-objective optimisation results with ntrials
func TexSingleObjTableStart(buf *bytes.Buffer, ntrials int) {
	io.Ff(buf, `
\begin{table*} \centering
\caption{Constrained single objective problems: Results}
\begin{tabular}[c]{cccc} \toprule
P & settings & results & histogram ($N_{trials}=%d$) \\ \hline
`, ntrials)
}

// TexSingleObjTableEnd ends table for single-objective optimisation results with ntrials
func TexSingleObjTableEnd(buf *bytes.Buffer) {
	io.Ff(buf, `\end{tabular}
\label{tab:singleobj}
\end{table*}`)
}

// TexSingleObjTableItem adds item to table for single-objective optimisation results with ntrials
func TexSingleObjTableItem(o *Optimiser, buf *bytes.Buffer, problem int, fref float64, nDigitsF, nDigitsX, nDigitsHist int) {
	hlen := 20
	SortByOva(o.Solutions, 0)
	best := o.Solutions[0]
	fmin, fave, fmax, fdev, F := o.StatMinProb(0, 20, fref, false)
	fmtF := "%g"
	fmtX := io.Sf("%%.%df", nDigitsX)
	fmtHist := io.Sf("%%.%df", nDigitsHist)
	io.Ff(buf, `%d
&
{$\!\begin{aligned}
    N_{sol}        & = %d \ACR
	N_{cpu}        & = %d \ACR
	t_{max}        & = %d \ACR
	\Delta t_{exc} & = %d \ACR
	N_{eval}       & = %d
\end{aligned}$}
&
{$\!\begin{aligned}
    f_{min}  &= `+fmtF+`  \ACR
             &\phantom{=}( `+fmtF+`) \ACR
    f_{ave}  &= `+fmtF+`  \ACR
    f_{max}  &= `+fmtF+` \ACR
    f_{dev}  &= {\bf `+fmtF+`} \ACR
    T_{sys}  &= %v
\end{aligned}$}
&
\begin{minipage}{7cm} \scriptsize
\begin{verbatim}
%s
\end{verbatim}
\end{minipage} \\
\multicolumn{4}{c}{$X_{best}$=`+fmtX+`} \\
`, problem, o.Nsol, o.Ncpu, o.Tf, o.DtExc, o.Nfeval,
		nice_num(fmin, nDigitsF), fref, nice_num(fave, nDigitsF), nice_num(fmax, nDigitsF), fdev, o.SysTime,
		rnd.BuildTextHist(nice_num(fmin-0.05, nDigitsHist), nice_num(fmax+0.05, nDigitsHist), 11, F, fmtHist, hlen),
		best.Flt)
}

// TexSingleObjReport produces Single-Objective table TeX report
//  nRowPerTab -- number of rows per table
func TexSingleObjReport(dirout, fnkey string, ntrials, nRowPerTab int, opts []*Optimiser, frefs []float64, nDigitsF, nDigitsX, nDigitsHist []int) {
	nprob := len(opts)
	if nRowPerTab < 1 {
		chk.Panic("number of rows per table must be greater than 0")
	}
	if len(nDigitsHist) < nprob {
		chk.Panic("size of slice with number of digits for histogram must be equal to or greater than the number of problems")
	}
	chk.IntAssert(len(frefs), nprob)
	buf := TexDocumentStart()
	for i, opt := range opts {
		if i%nRowPerTab == 0 {
			if i > 0 {
				io.Ff(buf, `\bottomrule`)
				TexSingleObjTableEnd(buf) // end previous table
				io.Ff(buf, "\n")
			}
			TexSingleObjTableStart(buf, ntrials) // begin new table
		} else {
			if i > 0 {
				io.Ff(buf, `\hline`)
			}
		}
		TexSingleObjTableItem(opt, buf, i+1, frefs[i], nDigitsF[i], nDigitsX[i], nDigitsHist[i])
	}
	io.Ff(buf, `\bottomrule`)
	TexSingleObjTableEnd(buf) // end previous table
	io.Ff(buf, "\n")
	TexDocumentEnd(buf)
	TexWrite(dirout, fnkey, buf, true)
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// nice_num returns a truncated float
func nice_num(x float64, ndigits int) float64 {
	s := io.Sf("%."+io.Sf("%d", ndigits)+"f", x)
	return io.Atof(s)
}

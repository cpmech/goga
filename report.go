// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"strings"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

const (
	SCIENTIFIC_NOTATION_TEX = true // convert scientific notation to tex
)

// TeX document ////////////////////////////////////////////////////////////////////////////////////

// TexDocumentStart starts TeX document
func TexDocumentStart() (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	io.Ff(buf, `\documentclass[a4paper]{article}

\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{booktabs}
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
	io.Ff(buf, `
\end{tabular}
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
func TexSingleObjTableEnd(buf *bytes.Buffer, label string) {
	io.Ff(buf, `\end{tabular}
\label{tab:%s}
\end{table*}
`, label)
}

// TexSingleObjTableItem adds item to table for single-objective optimisation results with ntrials
func TexSingleObjTableItem(o *Optimiser, buf *bytes.Buffer) {
	o.fix_formatting_data()
	FrefTxt := "N/A"
	if len(o.RptFref) > 0 {
		FrefTxt = tex(o.RptFmtF, o.RptFref[0])
	}
	Fmin, Fave, Fmax, Fdev, F := StatF(o, 0, false)
	FminTxt, FaveTxt, FmaxTxt, FdevTxt := tex(o.RptFmtF, Fmin), tex(o.RptFmtF, Fave), tex(o.RptFmtF, Fmax), tex(o.RptFmtFdev, Fdev)
	hist := rnd.BuildTextHist(nice(Fmin-0.05, o.HistNdig), nice(Fmax+0.05, o.HistNdig), o.HistNsta, F, o.HistFmt, o.HistLen)
	io.Ff(buf,
		`%s
&
{$\!\begin{aligned}
    N_{sol}        &= %d \\
	N_{cpu}        &= %d \\
	t_{max}        &= %d \\
	\Delta t_{exc} &= %d \\
	N_{eval}       &= %d
\end{aligned}$}
&
{$\!\begin{aligned}
    f_{min}  &= %s \\
             &\phantom{=} (%s) \\
    f_{ave}  &= %s \\
    f_{max}  &= %s \\
    f_{dev}  &= {\bf %s} \\
    T_{sys}  &= %v
\end{aligned}$}
&
\begin{minipage}{7cm} \scriptsize
\begin{verbatim}
%s
\end{verbatim}
\end{minipage} \\
\multicolumn{4}{c}{$X_{best}$=`+o.RptFmtX+`} \\
`,
		o.RptName,
		o.Nsol, o.Ncpu, o.Tf, o.DtExc, o.Nfeval,
		FminTxt, FrefTxt, FaveTxt, FmaxTxt, FdevTxt, o.SysTime,
		hist, o.BestOfBestFlt)
	if len(o.RptXref) == o.Nflt {
		io.Ff(buf, ` \multicolumn{4}{c}{$X_{ref}$=`+o.RptFmtX+`} \\`, o.RptXref)
	}
}

// TexSingleObjReport produces Single-Objective table TeX report
//  nRowPerTab -- number of rows per table
func TexSingleObjReport(dirout, fnkey, label string, nRowPerTab int, docHeader bool, opts []*Optimiser) {
	if nRowPerTab < 1 {
		nRowPerTab = len(opts)
	}
	var buf *bytes.Buffer
	if docHeader {
		buf = TexDocumentStart()
	} else {
		buf = new(bytes.Buffer)
	}
	idxtab := 0
	for i, opt := range opts {
		if i%nRowPerTab == 0 {
			if i > 0 {
				io.Ff(buf, `\bottomrule`)
				TexSingleObjTableEnd(buf, lbl(idxtab, label)) // end previous table
				io.Ff(buf, "\n\n\n")
				idxtab++
			}
			TexSingleObjTableStart(buf, opt.Ntrials) // begin new table
		} else {
			if i > 0 {
				io.Ff(buf, `\hline`)
				io.Ff(buf, "\n\n")
			}
		}
		TexSingleObjTableItem(opt, buf)
	}
	io.Ff(buf, `\bottomrule`)
	TexSingleObjTableEnd(buf, lbl(idxtab, label)) // end previous table
	io.Ff(buf, "\n")
	if docHeader {
		TexDocumentEnd(buf)
	}
	TexWrite(dirout, fnkey, buf, docHeader)
}

// F1F0 tables /////////////////////////////////////////////////////////////////////////////////////

// TexF1F0TableStart starts table for single-objective optimisation results with ntrials
func TexF1F0TableStart(buf *bytes.Buffer, ntrials int) {
	io.Ff(buf, `
\begin{table*} \centering
\caption{Constrained multiple objective problems. $N_{trials}=%d$}
\begin{tabular}[c]{cccc} \toprule
P & settings & error & spread \\ \hline
`, ntrials)
}

// TexF1F0TableEnd ends table for single-objective optimisation results with ntrials
func TexF1F0TableEnd(buf *bytes.Buffer, label string) {
	io.Ff(buf, `\end{tabular}
\label{tab:%s}
\end{table*}
`, label)
}

// TexF1F0TableItem adds item to table for single-objective optimisation results with ntrials
func TexF1F0TableItem(o *Optimiser, buf *bytes.Buffer) {
	o.fix_formatting_data()
	Emin, Eave, Emax, Edev, E, Lmin, Lave, Lmax, Ldev, _ := StatF1F0(o, false)
	EminTxt, EaveTxt, EmaxTxt, EdevTxt := tex(o.RptFmtE, Emin), tex(o.RptFmtE, Eave), tex(o.RptFmtE, Emax), tex(o.RptFmtEdev, Edev)
	LminTxt, LaveTxt, LmaxTxt, LdevTxt := tex(o.RptFmtL, Lmin), tex(o.RptFmtL, Lave), tex(o.RptFmtL, Lmax), tex(o.RptFmtLdev, Ldev)
	io.Ff(buf,
		`%s
&
{$\!\begin{aligned}
    N_{sol}        &= %d \\
	N_{cpu}        &= %d \\
	t_{max}        &= %d \\
	\Delta t_{exc} &= %d \\
	N_{eval}       &= %d
\end{aligned}$}
&
{$\!\begin{aligned}
    E_{min} &= %s \\
    E_{ave} &= %s \\
    E_{max} &= %s \\
    E_{dev} &= {\bf %s} \\
	T_{sys} &= %v
\end{aligned}$}
&
{$\!\begin{aligned}
    L_{min} &= %s \\
    L_{ave} &= %s \\
    L_{max} &= %s \\
    L_{dev} &= {\bf %s} \\
	count   &= %d
\end{aligned}$} \\
`,
		o.RptName,
		o.Nsol, o.Ncpu, o.Tf, o.DtExc, o.Nfeval,
		EminTxt, EaveTxt, EmaxTxt, EdevTxt, o.SysTime,
		LminTxt, LaveTxt, LmaxTxt, LdevTxt, len(E))
}

// TexF1F0Report produces multi-objective (f1f0) table TeX report
//  nRowPerTab -- number of rows per table
func TexF1F0Report(dirout, fnkey, label string, nRowPerTab int, docHeader bool, opts []*Optimiser) {
	if nRowPerTab < 1 {
		nRowPerTab = len(opts)
	}
	var buf *bytes.Buffer
	if docHeader {
		buf = TexDocumentStart()
	} else {
		buf = new(bytes.Buffer)
	}
	idxtab := 0
	for i, opt := range opts {
		if i%nRowPerTab == 0 {
			if i > 0 {
				io.Ff(buf, `\bottomrule`)
				TexF1F0TableEnd(buf, lbl(idxtab, label)) // end previous table
				io.Ff(buf, "\n\n\n")
				idxtab++
			}
			TexF1F0TableStart(buf, opt.Ntrials) // begin new table
		} else {
			if i > 0 {
				io.Ff(buf, `\hline`)
				io.Ff(buf, "\n\n")
			}
		}
		TexF1F0TableItem(opt, buf)
	}
	io.Ff(buf, `\bottomrule`)
	TexF1F0TableEnd(buf, lbl(idxtab, label)) // end previous table
	io.Ff(buf, "\n")
	if docHeader {
		TexDocumentEnd(buf)
	}
	TexWrite(dirout, fnkey, buf, docHeader)
}

// Multi tables /////////////////////////////////////////////////////////////////////////////////////

// TexMultiTableStart starts table for single-objective optimisation results with ntrials
func TexMultiTableStart(buf *bytes.Buffer, ntrials int) {
	io.Ff(buf, `
\begin{table*} \centering
\caption{Constrained multiple objective problems.}
\begin{tabular}[c]{cccc} \toprule
P & settings & error & histogram ($N_{trials}=%d$) \\ \hline
`, ntrials)
}

// TexMultiTableEnd ends table for single-objective optimisation results with ntrials
func TexMultiTableEnd(buf *bytes.Buffer, label string) {
	io.Ff(buf, `\end{tabular}
\label{tab:%s}
\end{table*}
`, label)
}

// TexMultiTableItem adds item to table for single-objective optimisation results with ntrials
func TexMultiTableItem(o *Optimiser, buf *bytes.Buffer) {
	o.fix_formatting_data()
	Emin, Eave, Emax, Edev, E := StatMulti(o, false)
	EminTxt, EaveTxt, EmaxTxt, EdevTxt := tex(o.RptFmtE, Emin), tex(o.RptFmtE, Eave), tex(o.RptFmtE, Emax), tex(o.RptFmtEdev, Edev)
	hist := rnd.BuildTextHist(nice(Emin-0.05, o.HistNdig), nice(Emax+0.05, o.HistNdig), o.HistNsta, E, o.HistFmt, o.HistLen)
	io.Ff(buf,
		`%s
&
{$\!\begin{aligned}
    N_{sol}        &= %d \\
	N_{cpu}        &= %d \\
	t_{max}        &= %d \\
	\Delta t_{exc} &= %d \\
	N_{eval}       &= %d
\end{aligned}$}
&
{$\!\begin{aligned}
    E_{min} &= %s \\
    E_{ave} &= %s \\
    E_{max} &= %s \\
    E_{dev} &= {\bf %s} \\
	T_{sys} &= %v
\end{aligned}$}
&
\begin{minipage}{7cm} \scriptsize
\begin{verbatim}
%s
\end{verbatim}
\end{minipage} \\
`,
		o.RptName,
		o.Nsol, o.Ncpu, o.Tf, o.DtExc, o.Nfeval,
		EminTxt, EaveTxt, EmaxTxt, EdevTxt, o.SysTime,
		hist)
}

// TexMultiReport produces multi-objective (f1f0) table TeX report
//  nRowPerTab -- number of rows per table
func TexMultiReport(dirout, fnkey, label string, nRowPerTab int, docHeader bool, opts []*Optimiser) {
	if nRowPerTab < 1 {
		nRowPerTab = len(opts)
	}
	var buf *bytes.Buffer
	if docHeader {
		buf = TexDocumentStart()
	} else {
		buf = new(bytes.Buffer)
	}
	idxtab := 0
	for i, opt := range opts {
		if i%nRowPerTab == 0 {
			if i > 0 {
				io.Ff(buf, `\bottomrule`)
				TexMultiTableEnd(buf, lbl(idxtab, label)) // end previous table
				io.Ff(buf, "\n\n\n")
				idxtab++
			}
			TexMultiTableStart(buf, opt.Ntrials) // begin new table
		} else {
			if i > 0 {
				io.Ff(buf, `\hline`)
				io.Ff(buf, "\n\n")
			}
		}
		TexMultiTableItem(opt, buf)
	}
	io.Ff(buf, `\bottomrule`)
	TexMultiTableEnd(buf, lbl(idxtab, label)) // end previous table
	io.Ff(buf, "\n")
	if docHeader {
		TexDocumentEnd(buf)
	}
	TexWrite(dirout, fnkey, buf, docHeader)
}

// write all values ////////////////////////////////////////////////////////////////////////////////

func WriteAllValues(dirout, fnkey string, opt *Optimiser) {
	var buf bytes.Buffer
	io.Ff(&buf, "%5s", "front")
	for i := 0; i < opt.Nova; i++ {
		io.Ff(&buf, "%24s", io.Sf("f%d", i))
	}
	for i := 0; i < opt.Noor; i++ {
		io.Ff(&buf, "%24s", io.Sf("u%d", i))
	}
	for i := 0; i < opt.Nflt; i++ {
		io.Ff(&buf, "%24s", io.Sf("x%d", i))
	}
	for i := 0; i < opt.Nint; i++ {
		io.Ff(&buf, "%24s", io.Sf("y%d", i))
	}
	io.Ff(&buf, "\n")
	for _, sol := range opt.Solutions {
		io.Ff(&buf, "%5d", sol.FrontId)
		for i := 0; i < opt.Nova; i++ {
			io.Ff(&buf, "%24g", sol.Ova[i])
		}
		for i := 0; i < opt.Noor; i++ {
			io.Ff(&buf, "%24g", sol.Oor[i])
		}
		for i := 0; i < opt.Nflt; i++ {
			io.Ff(&buf, "%24g", sol.Flt[i])
		}
		for i := 0; i < opt.Nint; i++ {
			io.Ff(&buf, "%24g", sol.Int[i])
		}
		io.Ff(&buf, "\n")
	}
	io.WriteFileVD(dirout, fnkey+".res", &buf)
}

func tex(fmt string, num float64) (l string) {
	if fmt == "" {
		fmt = "%g"
	}
	l = io.Sf(fmt, num)
	if SCIENTIFIC_NOTATION_TEX {
		s := strings.Split(l, "e-")
		if len(s) == 2 {
			e := s[1]
			if e[0] == '0' {
				e = string(e[1])
			}
			l = s[0] + "\\cdot 10^{-" + e + "}"
		}
	}
	return
}

func lbl(i int, label string) string {
	C := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return io.Sf("%s:%s", label, string(C[i%len(C)]))
}

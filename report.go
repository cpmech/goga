// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"strings"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

const (
	SCIENTIFIC_NOTATION_TEX = true // convert scientific notation to tex
)

// TeX document ////////////////////////////////////////////////////////////////////////////////////

// TexDocumentStart starts TeX document
func TexDocumentStart(useGeom bool) (buf *bytes.Buffer) {
	str := ""
	if useGeom {
		str = `\usepackage[margin=1.5cm,footskip=0.5cm]{geometry}`
	}
	buf = new(bytes.Buffer)
	io.Ff(buf, `\documentclass[a4paper]{article}

\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{booktabs}
%s

\title{GOGA Report}
\author{Dorival Pedroso}

\begin{document}
`, str)
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
		_, err := io.RunCmd(false, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory=/tmp/goga/", tex)
		if err != nil {
			chk.Panic("%v", err)
		}
		io.PfBlue("file <%s/%s.pdf> generated\n", dirout, fnkey)
	}
}

// TexTableStart starts table
func TexTableStart(buf *bytes.Buffer, title, col4, col5 string, textSize string) {
	io.Ff(buf, `
\begin{table*} [!t] \centering
\caption{%s}
%s

\begin{tabular}[c]{ccccc} \toprule
P & settings & settings/info & %s & %s \\ \hline
`, title, textSize, col4, col5)
}

// TexTableEnd ends table
func TexTableEnd(buf *bytes.Buffer, label string) {
	io.Ff(buf, `
\end{tabular}
\label{tab:%s}
\end{table*}
`, label)
}

// TexReport produces table TeX report
//  Type:
//     1 -- one objective; with histogram of OVA
//     2 -- two objective; no histogram; with E(error) and L(spread)
//     3 -- multi objective; with histogram of error
func TexReport(dirout, fnkey, title, label string, Type, nRowPerTab int, docHeader, useGeom bool, textSize, miniPageSz, histTextSize string, opts []*Optimiser) {
	col4, col5 := "error", io.Sf("histogram ($N_{samples}=%d$)", opts[0].Nsamples)
	var addrow func(opt *Optimiser, buf *bytes.Buffer, miniPageSz, histTextSize string)
	switch Type {
	case 1:
		col4 = "objective"
		addrow = TexOneObjTableItem
	case 2:
		col5 = "spread"
		addrow = TexTwoObjTableItem
	default:
		addrow = TexMultiTableItem
	}
	if nRowPerTab < 1 {
		nRowPerTab = len(opts)
	}
	var buf *bytes.Buffer
	if docHeader {
		buf = TexDocumentStart(useGeom)
	} else {
		buf = new(bytes.Buffer)
	}
	idxtab := 0
	contd := ""
	for i, opt := range opts {
		if i%nRowPerTab == 0 {
			if i > 0 {
				io.Ff(buf, "\n\\bottomrule\n\n")
				TexTableEnd(buf, lbl(idxtab, label)) // end previous table
				io.Ff(buf, "\n\n\n\n\n\n")
				contd = " (contd.)"
				idxtab++
			}
			TexTableStart(buf, title+contd, col4, col5, textSize) // begin new table
		} else {
			if i > 0 {
				io.Ff(buf, "\n\\hline\n")
				io.Ff(buf, "\n\n")
			}
		}
		addrow(opt, buf, miniPageSz, histTextSize)
	}
	io.Ff(buf, "\n\\bottomrule\n\n")
	TexTableEnd(buf, lbl(idxtab, label)) // end previous table
	io.Ff(buf, "\n\n\n\n\n\n")
	if docHeader {
		TexDocumentEnd(buf)
	}
	TexWrite(dirout, fnkey, buf, docHeader)
}

// one-obj table ///////////////////////////////////////////////////////////////////////////////////

// TexOneObjTableItem adds item to one-obj table
func TexOneObjTableItem(o *Optimiser, buf *bytes.Buffer, miniPageSz, histTextSize string) {
	o.fix_formatting_data()
	FrefTxt := "N/A"
	if len(o.RptFref) > 0 {
		FrefTxt = tex(o.RptFmtF, o.RptFref[0])
	}
	Fmin, Fave, Fmax, Fdev, F := StatF(o, 0, false)
	FminTxt, FaveTxt, FmaxTxt, FdevTxt := tex(o.RptFmtF, Fmin), tex(o.RptFmtF, Fave), tex(o.RptFmtF, Fmax), tex(o.RptFmtFdev, Fdev)
	hist := rnd.BuildTextHist(nice(Fmin, o.HistNdig)-o.HistDelFmin, nice(Fmax, o.HistNdig)+o.HistDelFmax, o.HistNsta, F, o.HistFmt, o.HistLen)
	io.Ff(buf, `
%s
&
{$\!\begin{aligned}
    N_{sol}        &= %d \\
	N_{cpu}        &= %d \\
	t_{max}        &= %d \\
	\Delta t_{exc} &= %d
\end{aligned}$}
&
{$\!\begin{aligned}
	%s_{ref} &= {\bf %s} \\
	C_{DE}   &= %g \\
	N_{eval} &= %d \\
	T_{sys}  &= %v
\end{aligned}$}
&
{$\!\begin{aligned}
    %s_{min} &= {\bf %s} \\
    %s_{ave} &= %s \\
    %s_{max} &= %s \\
    %s_{dev} &= {\bf %s}
\end{aligned}$}
&
\begin{minipage}{%s} %s
\begin{verbatim}
%s \end{verbatim}
\end{minipage} \\
\multicolumn{5}{c}{{\scriptsize $x_{best}$=`+o.RptFmtX+`}} \\`,
		o.RptName,
		o.Nsol, o.Ncpu, o.Tf, o.DtExc,
		o.RptWordF, FrefTxt, o.DEC, o.Nfeval, dround(o.SysTime, 0.001e9),
		o.RptWordF, FminTxt, o.RptWordF, FaveTxt, o.RptWordF, FmaxTxt, o.RptWordF, FdevTxt,
		miniPageSz, histTextSize, hist,
		o.BestOfBestFlt)
	if len(o.RptXref) == o.Nflt {
		io.Ff(buf, `
\multicolumn{5}{c}{{\scriptsize $x_{ref.}$=`+o.RptFmtX+`}} \\`, o.RptXref)
	}
	io.Ff(buf, "\n")
}

// two-obj table ///////////////////////////////////////////////////////////////////////////////////

// TexTwoObjTableItem adds item to two-obj table
func TexTwoObjTableItem(o *Optimiser, buf *bytes.Buffer, miniPageSz, histTextSize string) {
	o.fix_formatting_data()
	Emin, Eave, Emax, Edev, E, Lmin, Lave, Lmax, Ldev, _ := StatF1F0(o, false)
	EminTxt, EaveTxt, EmaxTxt, EdevTxt := tex(o.RptFmtE, Emin), tex(o.RptFmtE, Eave), tex(o.RptFmtE, Emax), tex(o.RptFmtEdev, Edev)
	LminTxt, LaveTxt, LmaxTxt, LdevTxt := tex(o.RptFmtL, Lmin), tex(o.RptFmtL, Lave), tex(o.RptFmtL, Lmax), tex(o.RptFmtLdev, Ldev)
	io.Ff(buf, `
%s
&
{$\!\begin{aligned}
    N_{sol}        &= %d \\
	N_{cpu}        &= %d \\
	t_{max}        &= %d \\
	\Delta t_{exc} &= %d
\end{aligned}$}
&
{$\!\begin{aligned}
	count    &= %d \\
	C_{DE}   &= %g \\
	N_{eval} &= %d \\
	T_{sys}  &= %v
\end{aligned}$}
&
{$\!\begin{aligned}
    E_{min} &= %s \\
    E_{ave} &= %s \\
    E_{max} &= %s \\
    E_{dev} &= {\bf %s}
\end{aligned}$}
&
{$\!\begin{aligned}
    L_{min} &= %s \\
    L_{ave} &= %s \\
    L_{max} &= %s \\
    L_{dev} &= {\bf %s}
\end{aligned}$} \\
`,
		o.RptName,
		o.Nsol, o.Ncpu, o.Tf, o.DtExc,
		len(E), o.DEC, o.Nfeval, dround(o.SysTime, 0.001e9),
		EminTxt, EaveTxt, EmaxTxt, EdevTxt,
		LminTxt, LaveTxt, LmaxTxt, LdevTxt)
}

// multi-obj table //////////////////////////////////////////////////////////////////////////////////

// TexMultiTableItem adds item to multi-obj table
func TexMultiTableItem(o *Optimiser, buf *bytes.Buffer, miniPageSz, histTextSize string) {
	o.fix_formatting_data()
	Ekey, Emin, Eave, Emax, Edev, E := StatMulti(o, false)
	EminTxt, EaveTxt, EmaxTxt, EdevTxt := tex(o.RptFmtE, Emin), tex(o.RptFmtE, Eave), tex(o.RptFmtE, Emax), tex(o.RptFmtEdev, Edev)
	hist := rnd.BuildTextHist(nice(Emin, o.HistNdig)-o.HistDelEmin, nice(Emax, o.HistNdig)+o.HistDelEmax, o.HistNsta, E, o.HistFmt, o.HistLen)
	io.Ff(buf, `
%s
&
{$\!\begin{aligned}
    N_{sol}        &= %d \\
	N_{cpu}        &= %d \\
	t_{max}        &= %d \\
	\Delta t_{exc} &= %d
\end{aligned}$}
&
{$\!\begin{aligned}
	N_{f}    &= %d \\
	C_{DE}   &= %g \\
	N_{eval} &= %d \\
	T_{sys}  &= %v
\end{aligned}$}
&
{$\!\begin{aligned}
    %s_{min} &= %s \\
    %s_{ave} &= %s \\
    %s_{max} &= %s \\
    %s_{dev} &= {\bf %s}
\end{aligned}$}
&
\begin{minipage}{%s} %s
\begin{verbatim}
%s \end{verbatim}
\end{minipage} \\
`,
		o.RptName,
		o.Nsol, o.Ncpu, o.Tf, o.DtExc,
		o.Nova, o.DEC, o.Nfeval, dround(o.SysTime, 0.001e9),
		Ekey, EminTxt, Ekey, EaveTxt, Ekey, EmaxTxt, Ekey, EdevTxt,
		miniPageSz, histTextSize, hist)
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
			if e == "00" {
				l = s[0]
				return
			}
			if e[0] == '0' {
				e = string(e[1])
			}
			l = s[0] + "\\cdot 10^{-" + e + "}"
		}
		s = strings.Split(l, "e+")
		if len(s) == 2 {
			e := s[1]
			if e == "00" {
				l = s[0]
				return
			}
			if e[0] == '0' {
				e = string(e[1])
			}
			l = s[0] + "\\cdot 10^{+" + e + "}"
		}
	}
	return
}

func lbl(i int, label string) string {
	C := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return io.Sf("%s:%s", label, string(C[i%len(C)]))
}

func dround(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}

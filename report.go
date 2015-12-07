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

func TexDocumentEnd(buf *bytes.Buffer) {
	io.Ff(buf, `\end{document}`)
}

func TexSingleObjTableStart(buf *bytes.Buffer, ntrials int) {
	io.Ff(buf, `\begin{table} \centering
\caption{Constrained single objective problems: Results}
\begin{tabular}[c]{cccc} \toprule
problem & settings & results & histogram ($N_{trials}=%d$) \\ \hline
`, ntrials)
}

func TexSingleObjTableEnd(buf *bytes.Buffer) {
	io.Ff(buf, `\end{tabular}
\label{tab:singleobj}
\end{table}`)
}

func (o *Optimiser) TexSingleObjTableItem(buf *bytes.Buffer, problem int, fref float64, fmtFave, fmtFdev, fmtHist, fmtX string) {
	hlen := 28
	SortByOva(o.Solutions, 0)
	best := o.Solutions[0]
	fmin, fave, fmax, fdev, F := o.StatMinProb(0, 20, fref, false)
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
    f_{min}  &= `+fmtFave+`  \ACR
    f_{ave}  &= {\bf `+fmtFave+`}  \ACR
             &\phantom{=}( `+fmtFave+`) \ACR
    f_{max}  &= `+fmtFave+`  \ACR
    f_{dev}  &= `+fmtFdev+`
\end{aligned}$}
&
\begin{minipage}{7cm} \scriptsize
\begin{verbatim}
%s\end{verbatim}
\end{minipage} \\
\multicolumn{4}{c}{best=`+fmtX+`} \\
`, problem, o.Nsol, o.Ncpu, o.Tf, o.DtExc, o.Nfeval, fmin, fave, fref, fmax, fdev,
		rnd.BuildTextHist(nice_num(fmin-0.05), nice_num(fmax+0.05), 11, F, fmtHist, hlen),
		best.Flt)
}

func TexWrite(fnkey string, buf *bytes.Buffer, dorun bool) {
	tex := fnkey + ".tex"
	io.WriteFileD("/tmp/goga", tex, buf)
	if dorun {
		_, err := io.RunCmd(true, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory=/tmp/goga/", tex)
		if err != nil {
			chk.Panic("%v", err)
		}
	}
}

// nice_num returns a truncated float
func nice_num(x float64) float64 {
	s := io.Sf("%.2f", x)
	return io.Atof(s)
}

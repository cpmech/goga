// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// TexReport produces TeX table report
//  Type:
//     1 -- one objective. with histogram of OVA
//     2 -- two objective. no histogram. with E(error) and L(spread)
//     3 -- multi objective. with histogram of error
//     4 -- one objective. no histogram. compact table
type TexReport struct {

	// options
	DirOut       string // output directory
	Fnkey        string // filename key
	Title        string // title of table
	RefLabel     string // lable of table for referencing
	TextSize     string // formatting string for size of text
	MiniPageSz   string // string for minipage size
	HistTextSize string // formatting string for histogram size
	Type         int    // type of table (see above)
	NRowPerTab   int    // number of rows per table. -1 means all rows
	UseGeom      bool   // use TeX geometry package
	RunPDF       bool   // generate PDF
	ShowNsol     bool   // show Nsol
	ShowNcpu     bool   // show Ncpu
	ShowTmax     bool   // show Tf
	ShowDtExc    bool   // show Dtexc
	ShowDEC      bool   // show DE coefficient
	ShowX01      bool   // show x[0] and x[1] in table

	// constants
	DroundCte time.Duration // constant for dround(duration) function

	// input
	Opts []*Optimiser // all optimisers

	// derived
	symbF     string
	col4      string
	col5      string
	buf       *bytes.Buffer
	binp      *bytes.Buffer // buffer for input data
	bxres     *bytes.Buffer // buffer for x results
	singleObj bool
}

// SetDefault sets default options for report
func NewTexReport(opts []*Optimiser) (o *TexReport) {

	// new struct
	o = new(TexReport)

	// options
	o.DirOut = "/tmp/goga"
	o.Fnkey = "texreport"
	o.Title = "Goga Report"
	o.RefLabel = ""
	o.TextSize = `\scriptsize  \setlength{\tabcolsep}{0.5em}`
	o.MiniPageSz = "4.1cm"
	o.HistTextSize = `\fontsize{5pt}{6pt}`
	o.Type = 4
	o.NRowPerTab = -1
	o.UseGeom = true
	o.RunPDF = true
	o.ShowNsol = false
	o.ShowNcpu = false
	o.ShowTmax = false
	o.ShowDtExc = false
	o.ShowDEC = false
	o.ShowX01 = true

	// constants
	o.DroundCte = 0.001e9 // 0.0001e9

	// input
	o.Opts = opts

	// buffers
	o.buf = new(bytes.Buffer)
	o.binp = new(bytes.Buffer)
	o.bxres = new(bytes.Buffer)

	// check
	if len(o.Opts) < 1 {
		chk.Panic("slice Opts must be set with at least one item")
	}
	if o.Type < 1 || o.Type > 4 {
		chk.Panic("type of table must be in [1,4]")
	}

	// symbol for f
	o.symbF = o.Opts[0].RptWordF
	if o.symbF == "" {
		o.symbF = "f"
	}

	// flag
	o.singleObj = o.Opts[0].Nova == 1
	return
}

// input data table //////////////////////////////////////////////////////////////////////////

// inputHeader adds table header for input data table
func (o *TexReport) inputHeader() {
	io.Ff(o.binp, `
\begin{table*} [!t] \centering
\caption{%s.}
%s

\begin{tabular}[c]{cccccc} \toprule
P & $N_{sol}$ & $N_{cpu}$ & $t_{max}$ & $\Delta t_{exc}$ & $C_{DE}$ \\ \hline
`, o.Title+": input parameters", o.TextSize)
}

// inputRow adds row to input data table
func (o *TexReport) inputRow(opt *Optimiser) {
	io.Ff(o.binp, "%s & %d & %d & %d & %d & %g \\\\ \n", opt.RptName, opt.Nsol, opt.Ncpu, opt.Tf, opt.DtExc, opt.DEC)
}

// inputFooter adds foote to input data table
func (o *TexReport) inputFooter() {
	io.Ff(o.binp, `
\bottomrule
\end{tabular}
\label{tab:%s}
\end{table*}
`, o.RefLabel+"Inp")
}

// x-results table ///////////////////////////////////////////////////////////////////////////

// xResHeader adds table header for X results
func (o *TexReport) xResHeader() {
	if !o.singleObj {
		return
	}
	io.Ff(o.bxres, `
\begin{table*} [!t] \centering
\caption{%s.}
%s

\begin{tabular}[c]{cl} \toprule
P &
x \\ \hline
`, o.Title+": x values", o.TextSize)
}

// xResRow adds row to table with X results
func (o *TexReport) xResRow(opt *Optimiser) {
	if !o.singleObj {
		return
	}
	io.Ff(o.bxres, "%s &\n", opt.RptName)
	io.Ff(o.bxres, "  $x_{best} = $"+opt.RptFmtX+" \\\\ \n", opt.BestOfBestFlt)
	if len(opt.RptXref) == opt.Nflt {
		io.Ff(o.bxres, "& $x_{ref.} = $"+opt.RptFmtX+" \\\\ \n", opt.RptXref)
	}
}

// xResFooter adds footer to rable with X results
func (o *TexReport) xResFooter() {
	if !o.singleObj {
		return
	}
	io.Ff(o.bxres, `
\bottomrule
\end{tabular}
\label{tab:%s}
\end{table*}
`, o.RefLabel+"Xres")
}

// compact and normal tables /////////////////////////////////////////////////////////////////

// compactTableHeader adds table header for compact table
func (o *TexReport) compactTableHeader(contd string) {
	txtCols := "cccccccc"
	txtNsol := ""
	if o.ShowNsol {
		txtCols += "c"
		txtNsol = "& $N_{sol}$"
	}
	txtNcpu := ""
	if o.ShowNcpu {
		txtCols += "c"
		txtNcpu = "& $N_{cpu}$"
	}
	txtTmax := ""
	if o.ShowTmax {
		txtCols += "c"
		txtTmax = "& $t_{max}$"
	}
	txtDtExc := ""
	if o.ShowDtExc {
		txtCols += "c"
		txtDtExc = "& $\\Delta t_{exc}$"
	}
	txtDEC := ""
	if o.ShowDEC {
		txtCols += "c"
		txtDEC = "& $C_{DE}$"
	}
	txtX01 := ""
	if o.ShowX01 {
		txtCols += "cccc"
		txtX01 = "& $x_0$ & $x_0^{ref.}$ & $x_1$ & $x_1^{ref.}$"
	}
	io.Ff(o.buf, `
\begin{table*} [!t] \centering
\caption{%s: results.}
%s

\begin{tabular}[c]{%s} \toprule
P
%s  %s  %s  %s
%s & $N_{eval}$ & $T_{sys}$ & $%s_{ref}$ &
$%s_{min}$ & $%s_{ave}$ & $%s_{max}$ & $%s_{dev}$ 
%s
\\ \hline
`, o.Title+contd, o.TextSize, txtCols,
		txtNsol, txtNcpu, txtTmax, txtDtExc,
		txtDEC,
		o.symbF, o.symbF, o.symbF, o.symbF, o.symbF,
		txtX01)
}

// normalTableHeader adds table header for normal table
func (o *TexReport) normalTableHeader(contd string) {
	io.Ff(o.buf, `
\begin{table*} [!t] \centering
\caption{%s.}
%s

\begin{tabular}[c]{ccccc} \toprule
P & settings & settings/info & %s & %s \\ \hline
`, o.Title+contd, o.TextSize, o.col4, o.col5)
}

// tableFooter add table footer
func (o *TexReport) tableFooter(idxtab int) {
	io.Ff(o.buf, `
\bottomrule
\end{tabular}
\label{tab:%s}
\end{table*}
`, io.Sf("%s%d", o.RefLabel, idxtab))
}

// generate report ///////////////////////////////////////////////////////////////////////////////////

func (o *TexReport) Clear() {
	if o.binp != nil {
		o.binp.Reset()
	}
	if o.buf != nil {
		o.buf.Reset()
	}
	if o.bxres != nil {
		o.bxres.Reset()
	}
}

// Generate generates report
func (o *TexReport) Generate() {

	// functions
	o.col4 = "error"
	o.col5 = io.Sf("histogram ($N_{samples}=%d$)", o.Opts[0].Nsamples)
	addHeader := o.normalTableHeader
	addRow := o.oneNormalAddRow
	switch o.Type {
	case 1:
		o.col4 = "objective"
	case 2:
		o.col5 = "spread"
		addRow = o.twoAddRow
	case 3:
		addRow = o.multiAddRow
	case 4:
		addHeader = o.compactTableHeader
		addRow = o.oneCompactAddRow
	}

	// number of rows per table
	nRowPerTab := o.NRowPerTab
	if nRowPerTab < 1 {
		nRowPerTab = len(o.Opts)
	}
	if o.RefLabel == "" {
		o.RefLabel = o.Fnkey
	}

	// input and xres tables
	o.inputHeader()
	o.xResHeader()

	// add rows
	idxtab := 0
	contd := ""
	for i, opt := range o.Opts {
		if i%nRowPerTab == 0 {
			if i > 0 {
				o.tableFooter(idxtab) // end previous table
				io.Ff(o.buf, "\n\n\n")
				contd = " (contd.)"
				idxtab++
			}
			addHeader(contd) // begin new table
		} else {
			if i > 0 {
				if o.Type != 4 {
					io.Ff(o.buf, "\\hline\n")
				}
				//io.Ff(o.bxres, "\n\\hline\n")
				io.Ff(o.buf, "\n")
			}
		}
		addRow(opt)
		o.inputRow(opt)
		o.xResRow(opt)
	}

	// close tables
	o.inputFooter()
	io.Ff(o.binp, "\n\n\n")
	o.xResFooter()
	o.tableFooter(idxtab) // end previous table
	io.Ff(o.buf, "\n\n\n")

	// write table
	tex := o.Fnkey + ".tex"
	io.WriteFileVD(o.DirOut, tex, o.buf, o.binp, o.bxres)

	// generate PDF
	if o.RunPDF {
		header := new(bytes.Buffer)
		footer := new(bytes.Buffer)
		str := ""
		if o.UseGeom {
			str = `\usepackage[margin=1.5cm,footskip=0.5cm]{geometry}`
		}
		io.Ff(header, `\documentclass[a4paper]{article}

\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{booktabs}
%s

\title{GOGA Report}
\author{Dorival Pedroso}

\begin{document}
`, str)
		io.Ff(footer, `
\end{document}`)

		// write temporary TeX file
		tex = "tmp_" + tex
		io.WriteFileD(o.DirOut, tex, header, o.buf, o.binp, o.bxres, footer)

		// run pdflatex
		_, err := io.RunCmd(false, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory="+o.DirOut, tex)
		if err != nil {
			io.PfRed("pdflatex failed: %v\n", err)
			return
		}
		io.PfBlue("file <%s/tmp_%s.pdf> generated\n", o.DirOut, o.Fnkey)
	}
}

// one-obj table ///////////////////////////////////////////////////////////////////////////////////

// oneCompactAddRow adds row to compact one-obj table
func (o *TexReport) oneCompactAddRow(opt *Optimiser) {
	opt.fix_formatting_data()
	FrefTxt := "N/A"
	if len(opt.RptFref) > 0 {
		FrefTxt = tx(opt.RptFmtF, opt.RptFref[0])
	}
	Fmin, Fave, Fmax, Fdev, _ := StatF(opt, 0, false)
	FminTxt, FaveTxt, FmaxTxt, FdevTxt := tx(opt.RptFmtF, Fmin), tx(opt.RptFmtF, Fave), tx(opt.RptFmtF, Fmax), tx(opt.RptFmtFdev, Fdev)
	txtNsol := ""
	if o.ShowNsol {
		txtNsol = io.Sf("& %d", opt.Nsol)
	}
	txtNcpu := ""
	if o.ShowNcpu {
		txtNcpu = io.Sf("& %d", opt.Ncpu)
	}
	txtTmax := ""
	if o.ShowTmax {
		txtTmax = io.Sf("& %d", opt.Tf)
	}
	txtDtExc := ""
	if o.ShowDtExc {
		txtDtExc = io.Sf("& %d", opt.DtExc)
	}
	txtDEC := ""
	if o.ShowDEC {
		txtDEC = io.Sf("& %g", opt.DEC)
	}
	txtX01 := ""
	if o.ShowX01 {
		x0, x1, x0ref, x1ref := "N/A", "N/A", "N/A", "N/A"
		if len(opt.BestOfBestFlt) > 1 {
			x0 = io.Sf(opt.RptFmtX, opt.BestOfBestFlt[0])
			x1 = io.Sf(opt.RptFmtX, opt.BestOfBestFlt[1])
		}
		if len(opt.RptXref) > 1 {
			x0ref = io.Sf(opt.RptFmtX, opt.RptXref[0])
			x1ref = io.Sf(opt.RptFmtX, opt.RptXref[1])
		}
		txtX01 = io.Sf("& %s & (%s) & %s & (%s)", x0, x0ref, x1, x1ref)
	}
	io.Ff(o.buf, `%s  %s %s %s %s   %s   & %d & %v & (%s) &   %s & %s & %s & $%s$   %s \\`,
		opt.RptName,
		txtNsol, txtNcpu, txtTmax, txtDtExc,
		txtDEC,
		opt.Nfeval, dround(opt.SysTimeAve, o.DroundCte), FrefTxt,
		FminTxt, FaveTxt, FmaxTxt, FdevTxt,
		txtX01)
}

// oneNormalAddRow adds row to normal one-obj table
func (o *TexReport) oneNormalAddRow(opt *Optimiser) {
	opt.fix_formatting_data()
	FrefTxt := "N/A"
	if len(opt.RptFref) > 0 {
		FrefTxt = tx(opt.RptFmtF, opt.RptFref[0])
	}
	Fmin, Fave, Fmax, Fdev, F := StatF(opt, 0, false)
	FminTxt, FaveTxt, FmaxTxt, FdevTxt := tx(opt.RptFmtF, Fmin), tx(opt.RptFmtF, Fave), tx(opt.RptFmtF, Fmax), tx(opt.RptFmtFdev, Fdev)
	hist := rnd.BuildTextHist(nice(Fmin, opt.HistNdig)-opt.HistDelFmin, nice(Fmax, opt.HistNdig)+opt.HistDelFmax, opt.HistNsta, F, opt.HistFmt, opt.HistLen)
	io.Ff(o.buf, `
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
`,
		opt.RptName,
		opt.Nsol, opt.Ncpu, opt.Tf, opt.DtExc,
		opt.RptWordF, FrefTxt, opt.DEC, opt.Nfeval, dround(opt.SysTimeAve, o.DroundCte),
		opt.RptWordF, FminTxt, opt.RptWordF, FaveTxt, opt.RptWordF, FmaxTxt, opt.RptWordF, FdevTxt,
		o.MiniPageSz, o.HistTextSize, hist)
	io.Ff(o.buf, "\n")
}

// two-obj table ///////////////////////////////////////////////////////////////////////////////////

// twoAddRow adds row to two-obj table
func (o *TexReport) twoAddRow(opt *Optimiser) {
	opt.fix_formatting_data()
	Emin, Eave, Emax, Edev, E, Lmin, Lave, Lmax, Ldev, _ := StatF1F0(opt, false)
	EminTxt, EaveTxt, EmaxTxt, EdevTxt := tx(opt.RptFmtE, Emin), tx(opt.RptFmtE, Eave), tx(opt.RptFmtE, Emax), tx(opt.RptFmtEdev, Edev)
	LminTxt, LaveTxt, LmaxTxt, LdevTxt := tx(opt.RptFmtL, Lmin), tx(opt.RptFmtL, Lave), tx(opt.RptFmtL, Lmax), tx(opt.RptFmtLdev, Ldev)
	io.Ff(o.buf, `
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
		opt.RptName,
		opt.Nsol, opt.Ncpu, opt.Tf, opt.DtExc,
		len(E), opt.DEC, opt.Nfeval, dround(opt.SysTimeAve, o.DroundCte),
		EminTxt, EaveTxt, EmaxTxt, EdevTxt,
		LminTxt, LaveTxt, LmaxTxt, LdevTxt)
}

// multi-obj table //////////////////////////////////////////////////////////////////////////////////

// multiAddRow adds row to multi-obj table
func (o *TexReport) multiAddRow(opt *Optimiser) {
	opt.fix_formatting_data()
	Ekey, Emin, Eave, Emax, Edev, E := StatMulti(opt, false)
	EminTxt, EaveTxt, EmaxTxt, EdevTxt := tx(opt.RptFmtE, Emin), tx(opt.RptFmtE, Eave), tx(opt.RptFmtE, Emax), tx(opt.RptFmtEdev, Edev)
	hist := rnd.BuildTextHist(nice(Emin, opt.HistNdig)-opt.HistDelEmin, nice(Emax, opt.HistNdig)+opt.HistDelEmax, opt.HistNsta, E, opt.HistFmt, opt.HistLen)
	io.Ff(o.buf, `
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
		opt.RptName,
		opt.Nsol, opt.Ncpu, opt.Tf, opt.DtExc,
		opt.Nova, opt.DEC, opt.Nfeval, dround(opt.SysTimeAve, o.DroundCte),
		Ekey, EminTxt, Ekey, EaveTxt, Ekey, EmaxTxt, Ekey, EdevTxt,
		o.MiniPageSz, o.HistTextSize, hist)
}

// other reporting functions ///////////////////////////////////////////////////////////////////////

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

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func tx(fmt string, num float64) string {
	return utl.TexNum(fmt, num, true)
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

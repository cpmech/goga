// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"bytes"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// TexReport produces TeX table report
type TexReport struct {

	// constants
	DroundCte time.Duration // constant for dround(duration) function

	// input
	Opts []*Optimiser // all optimisers

	// options
	Title      string // title of table
	TextSize   string // formatting string for size of text
	NrowPerTab int    // number of rows per table. -1 means all rows
	UseGeom    bool   // use TeX geometry package
	Landscape  bool   // landscape paper
	RunPDF     bool   // generate PDF

	// options for histogram
	MiniPageSz   string // string for minipage size
	HistTextSize string // formatting string for histogram size

	// data for preparing the title of table
	ShowNsamples bool // show Nsamples

	// input data columns
	ShowNsol  bool // show Nsol
	ShowNcpu  bool // show Ncpu
	ShowTmax  bool // show Tmax
	ShowDtExc bool // show Dtexc
	ShowDEC   bool // show DE coefficient

	// stat columns
	ShowNfeval     bool // show Nfeval
	ShowSysTimeAve bool // show SysTimeAve
	ShowSysTimeTot bool // show SysTimeTot

	// results columns
	ShowFref bool // show Fref
	ShowFmin bool // show Fmin
	ShowFave bool // show Fave
	ShowFmax bool // show Fmax
	ShowFdev bool // show Fdev
	ShowX01  bool // show x[0] and x[1] in table
	ShowAllX bool // show all x values in table
	ShowXref bool // show X references as well as X values

	// multi-objective columns
	ShowEmin   bool // show Emin
	ShowEave   bool // show Eave
	ShowEmax   bool // show Emax
	ShowEdev   bool // show Edev
	ShowLmin   bool // show Lmin
	ShowLave   bool // show Lave
	ShowLmax   bool // show Lmax
	ShowLdev   bool // show Ldev
	ShowIGDmin bool // show IGDmin
	ShowIGDave bool // show IGDave
	ShowIGDmax bool // show IGDmax
	ShowIGDdev bool // show IGDdev

	// derived
	nsamples  int           // number of samples
	nfltMax   int           // max number of floats from all opt problems
	nintMax   int           // max number of integers from fall opt problems
	singleObj bool          // is single obj problem
	symbF     string        // symbol for F function
	buf       *bytes.Buffer // buffer with text to be written
}

// NewTexReport allocates new TexReport object
func NewTexReport(opts []*Optimiser) (o *TexReport) {

	// new struct
	o = new(TexReport)

	// constants
	o.DroundCte = 0.001e9 // 0.0001e9

	// input
	o.Opts = opts

	// options
	o.Title = "Goga Report"
	o.TextSize = `\scriptsize  \setlength{\tabcolsep}{0.5em}`
	o.NrowPerTab = 10
	o.UseGeom = true
	o.Landscape = false
	o.RunPDF = true

	// options for histogram
	o.MiniPageSz = "4.1cm"
	o.HistTextSize = `\fontsize{5pt}{6pt}`

	// check
	if len(o.Opts) < 1 {
		chk.Panic("slice Opts must be set with at least one item")
	}

	// derived
	o.nsamples = o.Opts[0].Nsamples
	o.nfltMax = o.Opts[0].Nflt
	o.nintMax = o.Opts[0].Nint
	for _, opt := range o.Opts {
		o.nfltMax = utl.Imax(o.nfltMax, opt.Nflt)
		o.nintMax = utl.Imax(o.nintMax, opt.Nint)
	}
	o.singleObj = o.Opts[0].Nova == 1
	o.buf = new(bytes.Buffer)

	// symbol for F
	o.symbF = o.Opts[0].RptWordF
	if o.symbF == "" {
		o.symbF = "f"
	}

	// set default data for tables
	o.SetColumnsDefault()
	return
}

// SetColumnsDefault sets default flags for table
func (o *TexReport) SetColumnsDefault() {
	if o.singleObj {
		o.SetColumnsSingleObj(true, false)
	} else {
		o.SetColumnsMultiObj(false)
	}
}

// SetColumnsInput sets flags to generate a table with input parameters
func (o *TexReport) SetColumnsInputData() {
	o.SetColumnsAll(false)
	o.ShowNsol = true
	o.ShowNcpu = true
	o.ShowTmax = true
	o.ShowDtExc = true
	o.ShowDEC = true
}

// SetColumnsXvalues sets flags to generate a table with xvalues
func (o *TexReport) SetColumnsXvalues() {
	o.SetColumnsAll(false)
	o.ShowAllX = true
	o.ShowXref = true
}

// SetColumnsSingleObj sets flags to generate a table for single-objective problems
func (o *TexReport) SetColumnsSingleObj(showX01, showAllX bool) {
	o.SetColumnsAll(false)
	o.ShowNsamples = true
	o.ShowSysTimeAve = true
	o.ShowFref = true
	o.ShowFmin = true
	o.ShowFave = true
	o.ShowFmax = true
	o.ShowFdev = true
	o.ShowX01 = showX01
	o.ShowAllX = showAllX
	o.ShowXref = true
}

// SetColumnsMultiObj sets flags to generate a table for multi-objective problems
func (o *TexReport) SetColumnsMultiObj(usingIGD bool) {
	o.SetColumnsAll(false)
	o.ShowNsamples = true
	o.ShowSysTimeAve = true
	if usingIGD {
		o.ShowIGDmin = true
		o.ShowIGDave = true
		o.ShowIGDmax = true
		o.ShowIGDdev = true
	} else {
		o.ShowEmin = true
		o.ShowEave = true
		o.ShowEmax = true
		o.ShowEdev = true
		o.ShowLmin = true
		o.ShowLave = true
		o.ShowLmax = true
		o.ShowLdev = true
	}
}

// SetColumnsAll sets flags to generate all columns
func (o *TexReport) SetColumnsAll(flag bool) {

	// data for preparing the title of table
	o.ShowNsamples = flag

	// input data columns
	o.ShowNsol = flag
	o.ShowNcpu = flag
	o.ShowTmax = flag
	o.ShowDtExc = flag
	o.ShowDEC = flag

	// stat columns
	o.ShowNfeval = flag
	o.ShowSysTimeAve = flag
	o.ShowSysTimeTot = flag

	// results columns
	o.ShowFref = flag
	o.ShowFmin = flag
	o.ShowFave = flag
	o.ShowFmax = flag
	o.ShowFdev = flag
	o.ShowX01 = flag
	o.ShowAllX = flag
	o.ShowXref = flag

	// multi-objective columns
	o.ShowEmin = flag
	o.ShowEave = flag
	o.ShowEmax = flag
	o.ShowEdev = flag
	o.ShowLmin = flag
	o.ShowLave = flag
	o.ShowLmax = flag
	o.ShowLdev = flag
	o.ShowIGDmin = flag
	o.ShowIGDave = flag
	o.ShowIGDmax = flag
	o.ShowIGDdev = flag
}

// Generate generates report
func (o *TexReport) Generate(dirout, fnkey string) {

	// clear buffer
	o.buf.Reset()

	// number of rows per table
	nRowPerTab := o.NrowPerTab
	if nRowPerTab < 1 {
		nRowPerTab = len(o.Opts)
	}

	// generate results table
	o.addHeader(". Results")
	idxtab := 0
	for i, opt := range o.Opts {
		if i%nRowPerTab == 0 {
			if i > 0 {
				o.addFooter(fnkey, idxtab) // end previous table
				io.Ff(o.buf, "\n")
				o.addHeader(" (contd.)") // begin new table
				idxtab++
			}
		} else {
			if i > 0 {
				io.Ff(o.buf, "\n")
			}
		}
		o.addRow(opt, false, false)
	}
	o.addFooter(fnkey, idxtab)

	// generate input data table
	io.Ff(o.buf, "\n")
	o.SetColumnsInputData()
	o.addHeader(". Input data")
	for _, opt := range o.Opts {
		o.addRow(opt, true, false)
	}
	o.addFooter(fnkey+"-inp", 0)

	// generate xvalues table
	io.Ff(o.buf, "\n")
	o.SetColumnsXvalues()
	o.addHeader(". X values")
	for _, opt := range o.Opts {
		o.addRow(opt, false, true)
	}
	o.addFooter(fnkey+"-xvals", 0)

	// write file
	fn := fnkey + ".tex"
	io.WriteFileVD(dirout, fn, o.buf)

	// generate PDF
	if o.RunPDF {
		pdf := new(bytes.Buffer)
		if o.Landscape {
			io.Ff(pdf, "\\documentclass[a4paper,landscape]{article}\n")
		} else {
			io.Ff(pdf, "\\documentclass[a4paper]{article}\n")
		}
		io.Ff(pdf, "\\usepackage{amsmath}\n")
		io.Ff(pdf, "\\usepackage{amssymb}\n")
		io.Ff(pdf, "\\usepackage{booktabs}\n")
		if o.UseGeom {
			io.Ff(pdf, "\\usepackage[margin=1.5cm,footskip=0.5cm]{geometry}\n")
		}
		io.Ff(pdf, "\n\\begin{document}\n\n")
		io.Ff(pdf, "%s\n", o.buf)
		io.Ff(pdf, "\\end{document}\n")

		// write temporary TeX file
		fn = "tmp_" + fnkey + ".tex"
		io.WriteFileD(dirout, fn, pdf)

		// run pdflatex
		_, err := io.RunCmd(false, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory="+dirout, fn)
		if err != nil {
			io.PfBlue("file <%s/%s> generated\n", dirout, fn)
			io.PfRed("pdflatex failed: %v\n", err)
			return
		}
		io.PfBlue("file <%s/tmp_%s.pdf> generated\n", dirout, fnkey)
	}
}

// add header, rows and footer /////////////////////////////////////////////////////////////////////

// addHeader adds table header
func (o *TexReport) addHeader(titleExtra string) {

	// begin TeX table
	tex1 := `\begin{table*} [!t] \centering` + "\n"

	// title
	if o.ShowNsamples {
		tex1 += io.Sf(`\caption{%s%s ($N_{samples}=%d$).}`, o.Title, titleExtra, o.nsamples) + "\n"
	} else {
		tex1 += io.Sf(`\caption{%s%s.}`, o.Title, titleExtra) + "\n"
	}

	// text size formatting
	tex1 += o.TextSize + "\n"

	// column descriptors
	txtc := "c" // first column: "P"

	// input data columns
	tex2 := "P" // problem id
	if o.ShowNsol {
		txtc += "c"
		tex2 += io.Sf(` & $N_{sol}$`)
	}
	if o.ShowNcpu {
		txtc += "c"
		tex2 += io.Sf(` & $N_{cpu}$`)
	}
	if o.ShowTmax {
		txtc += "c"
		tex2 += io.Sf(` & $t_{max}$`)
	}
	if o.ShowDtExc {
		txtc += "c"
		tex2 += io.Sf(` & ${\Delta t_{exc}}$`)
	}
	if o.ShowDEC {
		txtc += "c"
		tex2 += io.Sf(` & $C_{DE}$`)
	}

	// stat columns
	if o.ShowNfeval {
		txtc += "c"
		tex2 += io.Sf(` & $N_{eval}$`)
	}
	if o.ShowSysTimeAve {
		txtc += "c"
		tex2 += io.Sf(` & $T_{sys}^{ave}$`)
	}
	if o.ShowSysTimeTot {
		txtc += "c"
		tex2 += io.Sf(` & $T_{sys}^{tot}$`)
	}

	// results columns
	if o.ShowFref {
		txtc += "c"
		tex2 += io.Sf(` & $%s_{ref}$`, o.symbF)
	}
	if o.ShowFmin {
		txtc += "c"
		tex2 += io.Sf(` & $%s_{min}$`, o.symbF)
	}
	if o.ShowFave {
		txtc += "c"
		tex2 += io.Sf(` & $%s_{ave}$`, o.symbF)
	}
	if o.ShowFmax {
		txtc += "c"
		tex2 += io.Sf(` & $%s_{max}$`, o.symbF)
	}
	if o.ShowFdev {
		txtc += "c"
		tex2 += io.Sf(` & $%s_{dev}$`, o.symbF)
	}
	if o.ShowX01 && !o.ShowAllX {
		if o.ShowXref {
			txtc += "cccc"
			tex2 += io.Sf(` & $x_0$ & $x_0^{ref}$ & $x_1$ & $x_1^{ref}$`)
		} else {
			txtc += "cc"
			tex2 += io.Sf(` & $x_0$ & $x_1$`)
		}
	}
	if o.ShowAllX {
		for i := 0; i < o.nfltMax; i++ {
			txtc += "c"
			tex2 += io.Sf(` & $x_{%d}$`, i)
		}
	}

	// multi-objective columns
	if o.ShowEmin {
		txtc += "c"
		tex2 += io.Sf(` & $E_{min}$`)
	}
	if o.ShowEave {
		txtc += "c"
		tex2 += io.Sf(` & $E_{ave}$`)
	}
	if o.ShowEmax {
		txtc += "c"
		tex2 += io.Sf(` & $E_{max}$`)
	}
	if o.ShowEdev {
		txtc += "c"
		tex2 += io.Sf(` & $E_{dev}$`)
	}
	if o.ShowLmin {
		txtc += "c"
		tex2 += io.Sf(` & $L_{dev}$`)
	}
	if o.ShowLave {
		txtc += "c"
		tex2 += io.Sf(` & $L_{ave}$`)
	}
	if o.ShowLmax {
		txtc += "c"
		tex2 += io.Sf(` & $L_{max}$`)
	}
	if o.ShowLdev {
		txtc += "c"
		tex2 += io.Sf(` & $L_{dev}$`)
	}
	if o.ShowIGDmin {
		txtc += "c"
		tex2 += io.Sf(` & ${IGD}_{dev}$`)
	}
	if o.ShowIGDave {
		txtc += "c"
		tex2 += io.Sf(` & ${IGD}_{ave}$`)
	}
	if o.ShowIGDmax {
		txtc += "c"
		tex2 += io.Sf(` & ${IGD}_{max}$`)
	}
	if o.ShowIGDdev {
		txtc += "c"
		tex2 += io.Sf(` & ${IGD}_{dev}$`)
	}

	// new line
	tex2 += ` \\ \hline` + "\n"

	// begin TeX tabular
	tex1 += io.Sf(`\begin{tabular}[c]{%s} \toprule`, txtc) + "\n"

	// write to buffer
	io.Ff(o.buf, "%s%s", tex1, tex2)
}

// addRow adds row to table
func (o *TexReport) addRow(opt *Optimiser, inpDatTable, xvalsTable bool) {

	// results
	var Fmin, Fave, Fmax, Fdev float64
	var Emin, Eave, Emax, Edev float64
	var Lmin, Lave, Lmax, Ldev float64
	var IGDmin, IGDave, IGDmax, IGDdev float64
	if !inpDatTable && !xvalsTable {

		// compute F stat values
		if o.ShowFmin || o.ShowFave || o.ShowFmax || o.ShowFdev {
			Fmin, Fave, Fmax, Fdev, _ = StatF(opt, 0, false)
		}

		// compute E and L stat values
		if o.ShowLmin || o.ShowLave || o.ShowLmax || o.ShowLdev {
			Emin, Eave, Emax, Edev, _, Lmin, Lave, Lmax, Ldev, _ = StatF1F0(opt, false)
		} else {
			if o.ShowEmin || o.ShowEave || o.ShowEmax || o.ShowEdev {
				Emin, Eave, Emax, Edev, _ = StatMultiE(opt, false)
			}
		}

		// compute IGD stat values
		if o.ShowIGDmin || o.ShowIGDave || o.ShowIGDmax || o.ShowIGDdev {
			IGDmin, IGDave, IGDmax, IGDdev, _ = StatMultiIGD(opt, false)
		}
	}

	// fix formatting strings
	opt.fix_formatting_data()

	// input data columns
	tex := opt.RptName // problem id
	if o.ShowNsol {
		tex += io.Sf(` & %d`, opt.Nsol)
	}
	if o.ShowNcpu {
		tex += io.Sf(` & %d`, opt.Ncpu)
	}
	if o.ShowTmax {
		tex += io.Sf(` & %d`, opt.Tmax)
	}
	if o.ShowDtExc {
		tex += io.Sf(` & %d`, opt.DtExc)
	}
	if o.ShowDEC {
		tex += io.Sf(` & %g`, opt.DEC)
	}

	// stat columns
	if o.ShowNfeval {
		tex += io.Sf(` & %d`, opt.Nfeval)
	}
	if o.ShowSysTimeAve {
		tex += io.Sf(` & %v`, dround(opt.SysTimeAve, o.DroundCte))
	}
	if o.ShowSysTimeTot {
		tex += io.Sf(` & %v`, dround(opt.SysTimeTot, o.DroundCte))
	}

	// results columns
	if o.ShowFref {
		str := "N/A"
		if len(opt.RptFref) > 0 {
			str = tx(opt.RptFmtF, opt.RptFref[0])
		}
		tex += io.Sf(` & %s`, str)
	}
	if o.ShowFmin {
		tex += io.Sf(` & %s`, tx(opt.RptFmtF, Fmin))
	}
	if o.ShowFave {
		tex += io.Sf(` & %s`, tx(opt.RptFmtF, Fave))
	}
	if o.ShowFmax {
		tex += io.Sf(` & %s`, tx(opt.RptFmtF, Fmax))
	}
	if o.ShowFdev {
		tex += io.Sf(` & %s`, tx(opt.RptFmtFdev, Fdev))
	}
	if o.ShowX01 && !o.ShowAllX {
		if o.ShowXref {
			x0, x1, x0ref, x1ref := "N/A", "N/A", "N/A", "N/A"
			if len(opt.BestOfBestFlt) == opt.Nflt {
				x0 = io.Sf(opt.RptFmtX, opt.BestOfBestFlt[0])
				x1 = io.Sf(opt.RptFmtX, opt.BestOfBestFlt[1])
			}
			if len(opt.RptXref) == opt.Nflt {
				x0ref = io.Sf(opt.RptFmtX, opt.RptXref[0])
				x1ref = io.Sf(opt.RptFmtX, opt.RptXref[1])
			}
			tex += io.Sf(` & %s & (%s) & %s & (%s)`, x0, x0ref, x1, x1ref)
		} else {
			x0, x1 := "N/A", "N/A"
			if len(opt.BestOfBestFlt) == opt.Nflt {
				x0 = io.Sf(opt.RptFmtX, opt.BestOfBestFlt[0])
				x1 = io.Sf(opt.RptFmtX, opt.BestOfBestFlt[1])
			}
			tex += io.Sf(` & %s & %s`, x0, x1)
		}
	}
	if o.ShowAllX {
		for i := 0; i < o.nfltMax; i++ {
			if i >= opt.Nflt {
				tex += " & "
			} else {
				str := io.Sf(opt.RptFmtX, opt.BestOfBestFlt[i])
				tex += io.Sf(` & %s`, str)
			}
		}
		if xvalsTable {
			tex += " \\\\\n ref"
			for i := 0; i < o.nfltMax; i++ {
				if i >= opt.Nflt {
					tex += " & "
				} else {
					if len(opt.RptXref) == opt.Nflt {
						str := io.Sf(opt.RptFmtX, opt.RptXref[i])
						tex += io.Sf(` & %s`, str)
					} else {
						tex += " & N/A"
					}
				}
			}
		}
	}

	// multi-objective columns
	if o.ShowEmin {
		tex += io.Sf(` & %s`, tx(opt.RptFmtE, Emin))
	}
	if o.ShowEave {
		tex += io.Sf(` & %s`, tx(opt.RptFmtE, Eave))
	}
	if o.ShowEmax {
		tex += io.Sf(` & %s`, tx(opt.RptFmtE, Emax))
	}
	if o.ShowEdev {
		tex += io.Sf(` & %s`, tx(opt.RptFmtE, Edev))
	}
	if o.ShowLmin {
		tex += io.Sf(` & %s`, tx(opt.RptFmtL, Lmin))
	}
	if o.ShowLave {
		tex += io.Sf(` & %s`, tx(opt.RptFmtL, Lave))
	}
	if o.ShowLmax {
		tex += io.Sf(` & %s`, tx(opt.RptFmtL, Lmax))
	}
	if o.ShowLdev {
		tex += io.Sf(` & %s`, tx(opt.RptFmtL, Ldev))
	}
	if o.ShowIGDmin {
		tex += io.Sf(` & %s`, tx(opt.RptFmtE, IGDmin))
	}
	if o.ShowIGDave {
		tex += io.Sf(` & %s`, tx(opt.RptFmtE, IGDave))
	}
	if o.ShowIGDmax {
		tex += io.Sf(` & %s`, tx(opt.RptFmtE, IGDmax))
	}
	if o.ShowIGDdev {
		tex += io.Sf(` & %s`, tx(opt.RptFmtE, IGDdev))
	}

	// new line
	tex += " \\\\\n"

	// write to buffer
	io.Ff(o.buf, "%s", tex)
}

// addFooter add table footer
func (o *TexReport) addFooter(tableLabel string, idxtab int) {
	io.Ff(o.buf, "\n\\bottomrule\n")
	io.Ff(o.buf, "\\end{tabular}\n")
	if idxtab > 0 {
		io.Ff(o.buf, "\\label{tab:%s}\n", io.Sf("%s%d", tableLabel, idxtab))
	} else {
		io.Ff(o.buf, "\\label{tab:%s}\n", tableLabel)
	}
	io.Ff(o.buf, "\\end{table*}\n")
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
	return "$" + io.TexNum(fmt, num, true) + "$"
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

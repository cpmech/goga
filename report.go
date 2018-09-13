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
	Title        string  // title of table
	RtableNotes  string  // footnotes for table with results
	ItableNotes  string  // footnotes for table with input data
	XtableNotes  string  // footnotes for table with x values
	XtableFontSz string  // formatting string for size of text in x-values table
	TableColSep  float64 // column separation in table
	UseGeom      bool    // use TeX geometry package
	Landscape    bool    // landscape paper
	RunPDF       bool    // generate PDF
	DescHeader   string  // description header

	// title and description
	ShowNsamples    bool // show Nsamples
	ShowDescription bool // show description

	// columns: input data columns
	ShowNsol  bool // show Nsol
	ShowNcpu  bool // show Ncpu
	ShowTmax  bool // show Tmax
	ShowDtExc bool // show Dtexc
	ShowDEC   bool // show DE coefficient

	// columns: stat columns
	ShowNfeval     bool // show Nfeval
	ShowSysTimeAve bool // show SysTimeAve
	ShowSysTimeTot bool // show SysTimeTot

	// columns: results columns
	ShowFref bool // show Fref
	ShowFmin bool // show Fmin
	ShowFave bool // show Fave
	ShowFmax bool // show Fmax
	ShowFdev bool // show Fdev
	ShowX01  bool // show x[0] and x[1] in table
	ShowAllX bool // show all x values in table
	ShowXref bool // show X references as well as X values

	// columns: multi-objective columns
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
	nsamples  int    // number of samples
	nfltMax   int    // max number of floats from all opt problems
	nintMax   int    // max number of integers from fall opt problems
	singleObj bool   // is single obj problem
	symbF     string // symbol for F function
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
	o.XtableFontSz = "\\scriptsize"
	o.TableColSep = 0.5
	o.UseGeom = true
	o.Landscape = false
	o.RunPDF = true
	o.DescHeader = "desc"

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
	o.ShowNfeval = true
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
	o.ShowDescription = true
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
	o.ShowDescription = true
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
	o.ShowDescription = flag

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

	// TeX report
	rpt := io.Report{
		Title:       o.Title,
		Author:      "Goga Authors",
		Landscape:   o.Landscape,
		TableColSep: o.TableColSep,
	}

	// generate results table
	K, nrows, F, M := o.GenTable()
	rpt.AddTableF(o.Title+". Results.", fnkey+"-res", o.RtableNotes, K, nrows, F, M)

	// generate input data table
	o.SetColumnsInputData()
	K, nrows, F, M = o.GenTable()
	rpt.AddTableF(o.Title+". Input data.", fnkey+"-inp", o.ItableNotes, K, nrows, F, M)

	// generate xvalues table
	if o.singleObj {
		o.SetColumnsXvalues()
		K, nrows, F, M = o.GenTable()
		rpt.TableFontSz = o.XtableFontSz
		if o.ShowXref {
			rpt.RowGapPt = 14
			rpt.RowGapStep = 1
		}
		rpt.AddTableF(o.Title+". Solutions.", fnkey+"-sol", o.XtableNotes, K, nrows, F, M)
	}

	// save file
	rpt.WriteTexPdf(dirout, fnkey, nil)

	// save tables
	if o.singleObj {
		rpt.WriteTexTables(dirout, map[string]string{
			fnkey + "-res": "table-" + fnkey + "-results",
			fnkey + "-inp": "table-" + fnkey + "-inputdata",
			fnkey + "-sol": "table-" + fnkey + "-solutions",
		})
	} else {
		rpt.WriteTexTables(dirout, map[string]string{
			fnkey + "-res": "table-" + fnkey + "-results",
			fnkey + "-inp": "table-" + fnkey + "-inputdata",
		})
	}
}

// GenTable generates table for TeX report
//   K -- column keys
//   nrows -- number of rows
//   F -- map of functions: key => function returning formatted row value
//   M -- maps key to tex formatted text of this key (i.e. equation)
func (o *TexReport) GenTable() (K []string, nrows int, F map[string]io.FcnRow, M map[string]string) {

	// allocate maps
	F = make(map[string]io.FcnRow)
	M = make(map[string]string)
	nrows = len(o.Opts)

	// column: problem name
	K = append(K, "P")
	M["P"] = "P"
	F["P"] = func(i int) string { return o.Opts[i].RptName }

	// columns: input data columns
	if o.ShowNsol {
		K = append(K, "Nsol")
		M["Nsol"] = `$N_{sol}$`
		F["Nsol"] = func(i int) string { return io.Sf("%d", o.Opts[i].Nsol) }
	}
	if o.ShowNcpu {
		K = append(K, "Ncpu")
		M["Ncpu"] = `$N_{cpu}$`
		F["Ncpu"] = func(i int) string { return io.Sf("%d", o.Opts[i].Ncpu) }
	}
	if o.ShowTmax {
		K = append(K, "Tmax")
		M["Tmax"] = `$t_{max}$`
		F["Tmax"] = func(i int) string { return io.Sf("%d", o.Opts[i].Tmax) }
	}
	if o.ShowDtExc {
		K = append(K, "DtExc")
		M["DtExc"] = `${\Delta t_{exc}}$`
		F["DtExc"] = func(i int) string { return io.Sf("%d", o.Opts[i].DtExc) }
	}
	if o.ShowDEC {
		K = append(K, "DEC")
		M["DEC"] = `$C_{DE}$`
		F["DEC"] = func(i int) string { return io.Sf("%g", o.Opts[i].DEC) }
	}

	// columns: stat columns
	if o.ShowNfeval {
		K = append(K, "Nfeval")
		M["Nfeval"] = `$N_{eval}$`
		F["Nfeval"] = func(i int) string { return io.Sf("%d", o.Opts[i].Nfeval) }
	}
	if o.ShowSysTimeAve {
		K = append(K, "SysTimeAve")
		M["SysTimeAve"] = `$T_{sys}^{ave}$`
		F["SysTimeAve"] = func(i int) string { return io.Sf("%v", dround(o.Opts[i].SysTimeAve, o.DroundCte)) }
	}
	if o.ShowSysTimeTot {
		K = append(K, "SysTimeTot")
		M["SysTimeTot"] = `$T_{sys}^{tot}$`
		F["SysTimeTot"] = func(i int) string { return io.Sf("%v", dround(o.Opts[i].SysTimeTot, o.DroundCte)) }
	}

	// columns: results columns
	if o.ShowFref {
		K = append(K, "Fref")
		M["Fref"] = io.Sf(`${%s}_{ref}$`, o.symbF)
		F["Fref"] = func(i int) string {
			if len(o.Opts[i].RptFref) > 0 {
				return tx(o.Opts[i].RptFmtF, o.Opts[i].RptFref[0])
			} else {
				return "N/A"
			}
		}
	}
	if o.ShowFmin {
		K = append(K, "Fmin")
		M["Fmin"] = io.Sf(`${%s}_{min}$`, o.symbF)
		F["Fmin"] = func(i int) string { return tx(o.Opts[i].RptFmtF, o.Opts[i].Fmin[0]) }
	}
	if o.ShowFave {
		K = append(K, "Fave")
		M["Fave"] = io.Sf(`${%s}_{ave}$`, o.symbF)
		F["Fave"] = func(i int) string { return tx(o.Opts[i].RptFmtF, o.Opts[i].Fave[0]) }
	}
	if o.ShowFmax {
		K = append(K, "Fmax")
		M["Fmax"] = io.Sf(`${%s}_{max}$`, o.symbF)
		F["Fmax"] = func(i int) string { return tx(o.Opts[i].RptFmtF, o.Opts[i].Fmax[0]) }
	}
	if o.ShowFdev {
		K = append(K, "Fdev")
		M["Fdev"] = io.Sf(`${%s}_{dev}$`, o.symbF)
		F["Fdev"] = func(i int) string { return tx(o.Opts[i].RptFmtFdev, o.Opts[i].Fdev[0]) }
	}

	// x-values
	if o.ShowAllX {
		for j := 0; j < o.nfltMax; j++ {
			key := io.Sf("x%d", j)
			K = append(K, key)
			M[key] = io.Sf(`$x_{%d}$`, j)
			jcopy := j // must create copy because when the function is called, j will be changed
			F[key] = func(i int) string {
				nflt := o.Opts[i].Nflt
				if jcopy >= nflt {
					return ""
				}
				xval := io.Sf(o.Opts[i].RptFmtX, o.Opts[i].BestOfBestFlt[jcopy])
				if o.ShowXref {
					xref := "N/A"
					if len(o.Opts[i].RptXref) == nflt {
						xref = io.Sf(o.Opts[i].RptFmtX, o.Opts[i].RptXref[jcopy])
					}
					return io.Sf(`$\begin{matrix} %s \\ %s \end{matrix}$`, xval, xref)
				} else {
					return xval
				}
			}
		}
	}

	// columns: multi-objective columns
	if o.ShowEmin {
		K = append(K, "Emin")
		M["Emin"] = `$E_{min}$`
		F["Emin"] = func(i int) string { return tx(o.Opts[i].RptFmtE, o.Opts[i].Emin) }
	}
	if o.ShowEave {
		K = append(K, "Eave")
		M["Eave"] = `$E_{ave}$`
		F["Eave"] = func(i int) string { return tx(o.Opts[i].RptFmtE, o.Opts[i].Eave) }
	}
	if o.ShowEmax {
		K = append(K, "Emax")
		M["Emax"] = `$E_{max}$`
		F["Emax"] = func(i int) string { return tx(o.Opts[i].RptFmtE, o.Opts[i].Emax) }
	}
	if o.ShowEdev {
		K = append(K, "Edev")
		M["Edev"] = `$E_{dev}$`
		F["Edev"] = func(i int) string { return tx(o.Opts[i].RptFmtEdev, o.Opts[i].Edev) }
	}
	if o.ShowLmin {
		K = append(K, "Lmin")
		M["Lmin"] = `$L_{dev}$`
		F["Lmin"] = func(i int) string { return tx(o.Opts[i].RptFmtL, o.Opts[i].Lmin) }
	}
	if o.ShowLave {
		K = append(K, "Lave")
		M["Lave"] = `$L_{ave}$`
		F["Lave"] = func(i int) string { return tx(o.Opts[i].RptFmtL, o.Opts[i].Lave) }
	}
	if o.ShowLmax {
		K = append(K, "Lmax")
		M["Lmax"] = `$L_{max}$`
		F["Lmax"] = func(i int) string { return tx(o.Opts[i].RptFmtL, o.Opts[i].Lmax) }
	}
	if o.ShowLdev {
		K = append(K, "Ldev")
		M["Ldev"] = `$L_{dev}$`
		F["Ldev"] = func(i int) string { return tx(o.Opts[i].RptFmtLdev, o.Opts[i].Ldev) }
	}
	if o.ShowIGDmin {
		K = append(K, "IGDmin")
		M["IGDmin"] = `${IGD}_{dev}$`
		F["IGDmin"] = func(i int) string { return tx(o.Opts[i].RptFmtE, o.Opts[i].IGDmin) }
	}
	if o.ShowIGDave {
		K = append(K, "IGDave")
		M["IGDave"] = `${IGD}_{ave}$`
		F["IGDave"] = func(i int) string { return tx(o.Opts[i].RptFmtE, o.Opts[i].IGDave) }
	}
	if o.ShowIGDmax {
		K = append(K, "IGDmax")
		M["IGDmax"] = `${IGD}_{max}$`
		F["IGDmax"] = func(i int) string { return tx(o.Opts[i].RptFmtE, o.Opts[i].IGDmax) }
	}
	if o.ShowIGDdev {
		K = append(K, "IGDdev")
		M["IGDdev"] = `${IGD}_{dev}$`
		F["IGDdev"] = func(i int) string { return tx(o.Opts[i].RptFmtEdev, o.Opts[i].IGDdev) }
	}

	// column: problem description
	if o.ShowDescription {
		K = append(K, "desc")
		M["desc"] = o.DescHeader
		F["desc"] = func(i int) string { return o.Opts[i].RptDesc }
	}
	return
}

// other reporting functions ///////////////////////////////////////////////////////////////////////

// WriteAllValues export all values to file
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
			io.Ff(&buf, "%24d", sol.Int[i])
		}
		io.Ff(&buf, "\n")
	}
	io.WriteFileVD(dirout, fnkey+".res", &buf)
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func tx(fmt string, num float64) string {
	return "$" + io.TexNum(fmt, num, true) + "$"
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

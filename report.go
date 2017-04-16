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
	K, M, T, C := o.GenTable()
	rpt.AddTable(K, T, o.Title, fnkey, M, C)

	// generate input data table
	o.SetColumnsInputData()
	K, M, T, C = o.GenTable()
	rpt.AddTable(K, T, o.Title+". Input data", fnkey, M, C)

	// generate xvalues table
	o.SetColumnsXvalues()
	K, M, T, C = o.GenTable()
	rpt.TableFontSz = o.XtableFontSz
	rpt.AddTable(K, T, o.Title+". X values", fnkey, M, C)

	// save file
	err := rpt.WriteTexPdf("/tmp/goga", fnkey, nil)
	if err != nil {
		io.PfRed("pdflatex failed: %v\n", err)
	}
}

// GenTable generates table for TeX report
//   K -- keys
//   M -- key2tex map: converts key into formatted text (i.e. equation)
//   T -- table with results; a map
//   C -- key2convert map: converts numbers at column 'key' to string
func (o *TexReport) GenTable() (K []string, M map[string]string, T map[string][]float64, C map[string]io.FcnConvertNum) {

	// allocate maps
	M = make(map[string]string)
	T = make(map[string][]float64)
	C = make(map[string]io.FcnConvertNum)
	nrows := len(o.Opts)

	// column: problem name and desc
	//K=append(K,"P")
	//M["P"] = "P"
	//T["P"] = make([]float64,nrows)
	//C["P"] = func(i int, v float64) string { return io.Sf("%g", v) }

	// columns: input data columns
	if o.ShowNsol {
		K = append(K, "Nsol")
		M["Nsol"] = `$N_{sol}$`
		T["Nsol"] = make([]float64, nrows)
		C["Nsol"] = func(i int, v float64) string { return io.Sf("%g", v) }
	}
	if o.ShowNcpu {
		K = append(K, "Ncpu")
		M["Ncpu"] = `$N_{cpu}$`
		T["Ncpu"] = make([]float64, nrows)
		C["Ncpu"] = func(i int, v float64) string { return io.Sf("%g", v) }
	}
	if o.ShowTmax {
		K = append(K, "Tmax")
		M["Tmax"] = `$t_{max}$`
		T["Tmax"] = make([]float64, nrows)
		C["Tmax"] = func(i int, v float64) string { return io.Sf("%g", v) }
	}
	if o.ShowDtExc {
		K = append(K, "DtExc")
		M["DtExc"] = `${\Delta t_{exc}}$`
		T["DtExc"] = make([]float64, nrows)
		C["DtExc"] = func(i int, v float64) string { return io.Sf("%g", v) }
	}
	if o.ShowDEC {
		K = append(K, "DEC")
		M["DEC"] = `$C_{DE}$`
		T["DEC"] = make([]float64, nrows)
		C["DEC"] = func(i int, v float64) string { return io.Sf("%g", v) }
	}

	// columns: stat columns
	if o.ShowNfeval {
		K = append(K, "Nfeval")
		M["Nfeval"] = `$N_{eval}$`
		T["Nfeval"] = make([]float64, nrows)
		C["Nfeval"] = func(i int, v float64) string { return io.Sf("%g", v) }
	}
	if o.ShowSysTimeAve {
		K = append(K, "SysTimeAve")
		M["SysTimeAve"] = `$T_{sys}^{ave}$`
		T["SysTimeAve"] = make([]float64, nrows)
		C["SysTimeAve"] = func(i int, v float64) string {
			d := time.Duration(int(v))
			return io.Sf(`$\begin{matrix} %v \\ (%v) \end{matrix}$`, dround(d, o.DroundCte), dround(o.Opts[i].SysTimeAve, o.DroundCte))
		}
	}
	if o.ShowSysTimeTot {
		K = append(K, "SysTimeTot")
		M["SysTimeTot"] = `$T_{sys}^{tot}$`
		T["SysTimeTot"] = make([]float64, nrows)
		C["SysTimeTot"] = func(i int, v float64) string {
			d := time.Duration(int(v))
			return io.Sf("%v", dround(d, o.DroundCte))
		}
	}

	// columns: results columns
	if o.ShowFref {
		K = append(K, "Fref")
		M["Fref"] = io.Sf(`${%s}_{ref}$`, o.symbF)
		T["Fref"] = make([]float64, nrows)
		C["Fref"] = func(i int, v float64) string {
			if len(o.Opts[i].RptFref) > 0 {
				return tx(o.Opts[i].RptFmtF, v)
			} else {
				return "N/A"
			}
		}
	}
	if o.ShowFmin {
		K = append(K, "Fmin")
		M["Fmin"] = io.Sf(`${%s}_{min}$`, o.symbF)
		T["Fmin"] = make([]float64, nrows)
		C["Fmin"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtF, v) }
	}
	if o.ShowFave {
		K = append(K, "Fave")
		M["Fave"] = io.Sf(`${%s}_{ave}$`, o.symbF)
		T["Fave"] = make([]float64, nrows)
		C["Fave"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtF, v) }
	}
	if o.ShowFmax {
		K = append(K, "Fmax")
		M["Fmax"] = io.Sf(`${%s}_{max}$`, o.symbF)
		T["Fmax"] = make([]float64, nrows)
		C["Fmax"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtF, v) }
	}
	if o.ShowFdev {
		K = append(K, "Fdev")
		M["Fdev"] = io.Sf(`${%s}_{dev}$`, o.symbF)
		T["Fdev"] = make([]float64, nrows)
		C["Fdev"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtF, v) }
	}

	// x-values
	if o.ShowAllX {
		for j := 0; j < o.nfltMax; j++ {
			key := io.Sf("x%d", j)
			K = append(K, key)
			M[key] = io.Sf(`$x_{%d}$`, j)
			T[key] = make([]float64, nrows)
			jcopy := j // must create copy because when the function is called, j will be changed
			C[key] = func(i int, v float64) string {
				nflt := o.Opts[i].Nflt
				if jcopy >= nflt {
					return ""
				}
				xval := io.Sf(o.Opts[i].RptFmtX, v)
				if o.ShowXref {
					xref := "N/A"
					if len(o.Opts[i].RptXref) == nflt {
						xref = io.Sf(o.Opts[i].RptFmtX, o.Opts[i].RptXref[jcopy])
					}
					return io.Sf(`$\begin{matrix} %s \\ (%s) \end{matrix}$`, xval, xref)
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
		T["Emin"] = make([]float64, nrows)
		C["Emin"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtE, v) }
	}
	if o.ShowEave {
		K = append(K, "Eave")
		M["Eave"] = `$E_{ave}$`
		T["Eave"] = make([]float64, nrows)
		C["Eave"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtE, v) }
	}
	if o.ShowEmax {
		K = append(K, "Emax")
		M["Emax"] = `$E_{max}$`
		T["Emax"] = make([]float64, nrows)
		C["Emax"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtE, v) }
	}
	if o.ShowEdev {
		K = append(K, "Edev")
		M["Edev"] = `$E_{dev}$`
		T["Edev"] = make([]float64, nrows)
		C["Edev"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtE, v) }
	}
	if o.ShowLmin {
		K = append(K, "Lmin")
		M["Lmin"] = `$L_{dev}$`
		T["Lmin"] = make([]float64, nrows)
		C["Lmin"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtL, v) }
	}
	if o.ShowLave {
		K = append(K, "Lave")
		M["Lave"] = `$L_{ave}$`
		T["Lave"] = make([]float64, nrows)
		C["Lave"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtL, v) }
	}
	if o.ShowLmax {
		K = append(K, "Lmax")
		M["Lmax"] = `$L_{max}$`
		T["Lmax"] = make([]float64, nrows)
		C["Lmax"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtL, v) }
	}
	if o.ShowLdev {
		K = append(K, "Ldev")
		M["Ldev"] = `$L_{dev}$`
		T["Ldev"] = make([]float64, nrows)
		C["Ldev"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtL, v) }
	}
	if o.ShowIGDmin {
		K = append(K, "IGDmin")
		M["IGDmin"] = `${IGD}_{dev}$`
		T["IGDmin"] = make([]float64, nrows)
		C["IGDmin"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtE, v) }
	}
	if o.ShowIGDave {
		K = append(K, "IGDave")
		M["IGDave"] = `${IGD}_{ave}$`
		T["IGDave"] = make([]float64, nrows)
		C["IGDave"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtE, v) }
	}
	if o.ShowIGDmax {
		K = append(K, "IGDmax")
		M["IGDmax"] = `${IGD}_{max}$`
		T["IGDmax"] = make([]float64, nrows)
		C["IGDmax"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtE, v) }
	}
	if o.ShowIGDdev {
		K = append(K, "IGDdev")
		M["IGDdev"] = `${IGD}_{dev}$`
		T["IGDdev"] = make([]float64, nrows)
		C["IGDdev"] = func(i int, v float64) string { return tx(o.Opts[i].RptFmtE, v) }
	}

	// add rows
	for i, opt := range o.Opts {
		o.tableRow(T, i, opt)
	}
	return
}

// tableRow generates one row in table
func (o *TexReport) tableRow(T map[string][]float64, i int, opt *Optimiser) {

	// compute F stat values
	var errF error
	var Fmin, Fave, Fmax, Fdev float64
	if o.ShowFmin || o.ShowFave || o.ShowFmax || o.ShowFdev {
		Fmin, Fave, Fmax, Fdev, _, errF = StatF(opt, 0, false)
	}

	// compute E and L stat values
	var errE error
	var Emin, Eave, Emax, Edev float64
	var Lmin, Lave, Lmax, Ldev float64
	if o.ShowLmin || o.ShowLave || o.ShowLmax || o.ShowLdev {
		Emin, Eave, Emax, Edev, _, Lmin, Lave, Lmax, Ldev, _, errE = StatF1F0(opt, false)
	} else {
		if o.ShowEmin || o.ShowEave || o.ShowEmax || o.ShowEdev {
			Emin, Eave, Emax, Edev, _, errE = StatMultiE(opt, false)
		}
	}

	// compute IGD stat values
	var errIGD error
	var IGDmin, IGDave, IGDmax, IGDdev float64
	if o.ShowIGDmin || o.ShowIGDave || o.ShowIGDmax || o.ShowIGDdev {
		IGDmin, IGDave, IGDmax, IGDdev, _, errIGD = StatMultiIGD(opt, false)
	}

	// columns: multi-objective columns
	if o.ShowNsol {
		T["Nsol"][i] = float64(opt.Nsol)
	}
	if o.ShowNcpu {
		T["Ncpu"][i] = float64(opt.Ncpu)
	}
	if o.ShowTmax {
		T["Tmax"][i] = float64(opt.Tmax)
	}
	if o.ShowDtExc {
		T["DtExc"][i] = float64(opt.DtExc)
	}
	if o.ShowDEC {
		T["DEC"][i] = opt.DEC
	}

	// columns: multi-objective columns
	if o.ShowNfeval {
		T["Nfeval"][i] = float64(opt.Nfeval)
	}
	if o.ShowSysTimeAve {
		T["SysTimeAve"][i] = float64(opt.SysTimeAve.Nanoseconds())
	}
	if o.ShowSysTimeTot {
		T["SysTimeTot"][i] = float64(opt.SysTimeTot.Nanoseconds())
	}

	// columns: multi-objective columns
	if o.ShowFref {
		if len(opt.RptFref) > 0 {
			T["Fref"][i] = opt.RptFref[0]
		}
	}
	if o.ShowFmin && errF == nil {
		T["Fmin"][i] = Fmin
	}
	if o.ShowFave && errF == nil {
		T["Fave"][i] = Fave
	}
	if o.ShowFmax && errF == nil {
		T["Fmax"][i] = Fmax
	}
	if o.ShowFdev && errF == nil {
		T["Fdev"][i] = Fdev
	}

	// x-values
	if o.ShowAllX {
		for j := 0; j < opt.Nflt; j++ {
			key := io.Sf("x%d", j)
			T[key][i] = opt.BestOfBestFlt[j]
		}
	}

	// columns: multi-objective columns
	if o.ShowEmin && errE == nil {
		T["Emin"][i] = Emin
	}
	if o.ShowEave && errE == nil {
		T["Eave"][i] = Eave
	}
	if o.ShowEmax && errE == nil {
		T["Emax"][i] = Emax
	}
	if o.ShowEdev && errE == nil {
		T["Edev"][i] = Edev
	}
	if o.ShowLmin && errE == nil {
		T["Lmin"][i] = Lmin
	}
	if o.ShowLave && errE == nil {
		T["Lave"][i] = Lave
	}
	if o.ShowLmax && errE == nil {
		T["Lmax"][i] = Lmax
	}
	if o.ShowLdev && errE == nil {
		T["Ldev"][i] = Ldev
	}
	if o.ShowIGDmin && errIGD == nil {
		T["IGDmin"][i] = IGDmin
	}
	if o.ShowIGDave && errIGD == nil {
		T["IGDave"][i] = IGDave
	}
	if o.ShowIGDmax && errIGD == nil {
		T["IGDmax"][i] = IGDmax
	}
	if o.ShowIGDdev && errIGD == nil {
		T["IGDdev"][i] = IGDdev
	}
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

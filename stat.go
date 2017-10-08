// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

type Stat struct {

	// stat
	Nfeval     int             // number of function evaluations
	SysTimes   []time.Duration // all system times for each run
	SysTimeAve time.Duration   // average of all system times
	SysTimeTot time.Duration   // total system (real/CPU) time

	// formatting data for reports
	RptName         string    // problem name
	RptFref         []float64 // reference OVAs
	RptXref         []float64 // reference flts
	RptFmin         []float64 // min OVAs for reports/graphs
	RptFmax         []float64 // max OVAs for reports/graphs
	RptFmtF         string    // format for fmin, fave and fmax
	RptFmtFdev      string    // format for fdev
	RptFmtE         string    // format for emin, eave and emax
	RptFmtEdev      string    // format for edev
	RptFmtL         string    // format for lmin, lave and lmax
	RptFmtLdev      string    // format for ldev
	RptFmtX         string    // format for x values
	RptWordF        string    // word to use for 'f'; e.g. '\beta'
	RptDesc         string    // description text
	HistFmt         string    // format in histogram
	HistDelFmin     float64   // Δf for minimum f value in histogram
	HistDelFmax     float64   // Δf for minimum f value in histogram
	HistDelEmin     float64   // Δe for minimum e value in histogram
	HistDelEmax     float64   // Δe for minimum e value in histogram
	HistDelFminZero bool      // use zero for Δf (min)
	HistDelFmaxZero bool      // use zero for Δf (max)
	HistDelEminZero bool      // use zero for Δe (min)
	HistDelEmaxZero bool      // use zero for Δe (max)
	HistNdig        int       // number of digits in histogram
	HistNsta        int       // number of stations in histogram
	HistLen         int       // number of characters (bar length) in histogram

	// RunMany: best solutions
	BestOvas      [][]float64 // best OVAs [nova][nsamples]
	BestFlts      [][]float64 // best flts [nflt][nsamples]
	BestInts      [][]int     // best ints [nint][nsamples]
	BestOfBestOva []float64   // [nova]
	BestOfBestFlt []float64   // [nflt]
	BestOfBestInt []int       // [nint]

	// RunMany: checking multi-obj problems
	F1F0_func      func(f0 float64) float64  // f1(f0) function
	F1F0_err       []float64                 // max(error(f1))
	F1F0_arcLen    []float64                 // arc-length: spreading on (f0,f1) space
	F1F0_arcLenRef float64                   // reference arc-length along f1(f0) curve
	F1F0_f0ranges  [][]float64               // ranges of f0 values to compute arc-length
	Multi_fcnErr   func(f []float64) float64 // computes Pareto-optimal front error with many OVAs
	Multi_err      []float64                 // max(error(f[i]))
	Multi_fStar    [][]float64               // reference points on Pareto front [npoints][nova]
	Multi_IGD      []float64                 // IGD metric

	// RunMany: statistics: F
	Fmin []float64 // minimum of each F [iOva]
	Fave []float64 // average of each F [iOva]
	Fmax []float64 // maximum of each F [iOva]
	Fdev []float64 // deviation of each F [iOva]

	// RunMany: statistics: E and L
	Emin float64 // minimum E
	Eave float64 // avarage E
	Emax float64 // maximum E
	Edev float64 // deviation in E
	Lmin float64 // minimum L
	Lave float64 // avarage L
	Lmax float64 // maximum L
	Ldev float64 // deviation in L

	// RunMany: statistics: IGD
	IGDmin float64 // minimum IGD
	IGDave float64 // avarage IGD
	IGDmax float64 // maximum IGD
	IGDdev float64 // deviation in IGD
}

// RunMany runs many trials in order to produce statistical data
//   Input:
//     dirout -- directory to write files with results [may be ""]
//     fnkey  -- filename key with results (will add .res) [may be ""]
func (o *Optimiser) RunMany(dirout, fnkey string, constantSeed bool) {

	// benchmark
	t0 := time.Now()
	defer func() {
		o.SysTimeTot = time.Now().Sub(t0)
		var tmp int64
		for _, dur := range o.SysTimes {
			tmp += dur.Nanoseconds()
		}
		tmp /= int64(o.Nsamples)
		o.SysTimeAve = time.Duration(tmp)
	}()

	// disable verbose flag temporarily
	if o.Verbose {
		defer func() {
			o.Verbose = true
		}()
		o.Verbose = false
	}

	// remove previous results
	if fnkey != "" {
		io.RemoveAll(dirout + "/" + fnkey + "-*.res")
	}

	// denominator for calculation of L metric
	denominatorL := 1.0
	if o.F1F0_arcLenRef > 0 {
		denominatorL = o.F1F0_arcLenRef
	}

	// allocate variables
	o.SysTimes = make([]time.Duration, o.Nsamples)
	o.BestOvas = make([][]float64, o.Nova)
	o.BestFlts = make([][]float64, o.Nflt)
	o.BestInts = make([][]int, o.Nint)
	o.BestOfBestOva = make([]float64, o.Nova)
	o.BestOfBestFlt = make([]float64, o.Nflt)
	o.BestOfBestInt = make([]int, o.Nint)

	// perform trials
	for itrial := 0; itrial < o.Nsamples; itrial++ {

		// re-generate solutions
		o.Nfeval = 0
		if itrial > 0 {
			o.Reset(constantSeed)
		}

		// save initial solutions
		if fnkey != "" {
			WriteAllValues(dirout, io.Sf("%s-%04d_ini", fnkey, itrial), o)
		}

		// message
		if o.VerbStat {
			io.Pf(". . . running trial # %d\n", itrial)
		}

		// solve
		timeIni := time.Now()
		o.Solve()
		o.SysTimes[itrial] = time.Now().Sub(timeIni)

		// sort
		SortSolutions(o.Solutions, 0)

		// feasible solution
		if o.Solutions[0].Feasible() {

			// best solution
			best := o.Solutions[0]
			for i := 0; i < o.Nova; i++ {
				o.BestOvas[i] = append(o.BestOvas[i], best.Ova[i])
			}
			for i := 0; i < o.Nflt; i++ {
				o.BestFlts[i] = append(o.BestFlts[i], best.Flt[i])
			}
			for i := 0; i < o.Nint; i++ {
				o.BestInts[i] = append(o.BestInts[i], best.Int[i])
			}

			// best of all trials
			first_best := len(o.BestOvas[0]) == 1
			if first_best {
				copy(o.BestOfBestOva, best.Ova)
				copy(o.BestOfBestFlt, best.Flt)
				copy(o.BestOfBestInt, best.Int)
			} else {
				if best.Ova[0] < o.BestOfBestOva[0] {
					copy(o.BestOfBestOva, best.Ova)
					copy(o.BestOfBestFlt, best.Flt)
					copy(o.BestOfBestInt, best.Int)
				}
			}

			// check multi-objective results
			if o.F1F0_func != nil {
				var rms_err float64
				var nfeasible int
				for _, sol := range o.Solutions {
					if sol.Feasible() && sol.FrontId == 0 {
						f0, f1 := sol.Ova[0], sol.Ova[1]
						f1_cor := o.F1F0_func(f0)
						rms_err += math.Pow(f1-f1_cor, 2.0)
						nfeasible++
					}
				}
				if nfeasible > 0 {
					rms_err = math.Sqrt(rms_err / float64(nfeasible))
					o.F1F0_err = append(o.F1F0_err, rms_err)
				}
			}

			// arc-length along Pareto front
			if o.Nova == 2 {
				if best.Feasible() && best.FrontId == 0 && o.Solutions[1].FrontId == 0 {
					dist := 0.0
					for i := 1; i < o.Nsol; i++ {
						if o.Solutions[i].FrontId == 0 {
							F0, F1 := o.Solutions[i-1].Ova[0], o.Solutions[i-1].Ova[1]
							f0, f1 := o.Solutions[i].Ova[0], o.Solutions[i].Ova[1]
							if o.F1F0_f0ranges != nil {
								a := o.find_f0_spot(F0)
								b := o.find_f0_spot(f0)
								if a == -1 || b == -1 {
									continue
								}
								if a != b {
									//io.Pforan("\nF0=%g is in [%g,%g]\n", F0, o.F1F0_f0ranges[a][0], o.F1F0_f0ranges[a][1])
									//io.Pfpink("f0=%g is in [%g,%g]\n", f0, o.F1F0_f0ranges[b][0], o.F1F0_f0ranges[b][1])
									continue
								}
							}
							dist += math.Sqrt(math.Pow(f0-F0, 2.0) + math.Pow(f1-F1, 2.0))
						}
					}
					o.F1F0_arcLen = append(o.F1F0_arcLen, dist/denominatorL)
				}
			}

			// multiple OVAs
			if o.Nova > 1 && o.Multi_fcnErr != nil {
				var rms_err float64
				var nfeasible int
				for _, sol := range o.Solutions {
					if sol.Feasible() && sol.FrontId == 0 {
						f_err := o.Multi_fcnErr(sol.Ova)
						rms_err += f_err * f_err
						nfeasible++
					}
				}
				if nfeasible > 0 {
					rms_err = math.Sqrt(rms_err / float64(nfeasible))
					o.Multi_err = append(o.Multi_err, rms_err)
				}
			}

			// IGD metric
			if o.Nova > 1 && len(o.Multi_fStar) > 0 {
				o.Multi_IGD = append(o.Multi_IGD, o.calcIgd(o.Multi_fStar))
			}

			// save final solutions
			if fnkey != "" {
				f0min := best.Ova[0]
				for _, sol := range o.Solutions {
					f0min = utl.Min(f0min, sol.Ova[0])
				}
				WriteAllValues(dirout, io.Sf("%s-%04d_f0min=%g", fnkey, itrial, f0min), o)
			}
		}
	}

	// statistics: F
	o.Fmin = make([]float64, o.Nova)
	o.Fave = make([]float64, o.Nova)
	o.Fmax = make([]float64, o.Nova)
	o.Fdev = make([]float64, o.Nova)
	for i := 0; i < o.Nova; i++ {
		o.Fmin[i], o.Fave[i], o.Fmax[i], o.Fdev[i] = INF, INF, INF, INF
		if len(o.BestOvas[i]) > 1 && o.Nova == 1 {
			o.Fmin[i], o.Fave[i], o.Fmax[i], o.Fdev[i] = rnd.StatBasic(o.BestOvas[i], true)
		}
	}

	// statistics: E and L
	o.Emin, o.Eave, o.Emax, o.Edev = INF, INF, INF, INF
	o.Lmin, o.Lave, o.Lmax, o.Ldev = INF, INF, INF, INF
	if o.F1F0_func != nil {
		o.Emin, o.Eave, o.Emax, o.Edev = rnd.StatBasic(o.F1F0_err, true)
		o.Lmin, o.Lave, o.Lmax, o.Ldev = rnd.StatBasic(o.F1F0_arcLen, true)
	}
	if o.Multi_fcnErr != nil {
		o.Emin, o.Eave, o.Emax, o.Edev = rnd.StatBasic(o.Multi_err, true)
	}

	// statistics: IGD
	o.IGDmin, o.IGDave, o.IGDmax, o.IGDdev = INF, INF, INF, INF
	if len(o.Multi_IGD) > 0 {
		o.IGDmin, o.IGDave, o.IGDmax, o.IGDdev = rnd.StatBasic(o.Multi_IGD, true)
	}
}

// PrintStatF print statistical information corresponding to objective function idxF
func (o *Optimiser) PrintStatF(idxF int) {
	if len(o.BestOvas[idxF]) == 0 {
		io.Pf("there are no samples for statistical analysis\n")
		return
	}
	str := "\n"
	if len(o.RptFref) == o.Nova {
		str = io.Sf(" (%g)\n", o.RptFref[idxF])
	}
	io.Pf("fmin = %g\n", o.Fmin[idxF])
	io.Pf("fave = %g"+str, o.Fave[idxF])
	io.Pf("fmax = %g\n", o.Fmax[idxF])
	io.Pf("fdev = %g\n", o.Fdev[idxF])
	o.fix_formatting_data()
	io.Pf(rnd.BuildTextHist(
		nice(o.Fmin[idxF], o.HistNdig)-o.HistDelFmin,
		nice(o.Fmax[idxF], o.HistNdig)+o.HistDelFmax,
		o.HistNsta, o.BestOvas[idxF], o.HistFmt, o.HistLen))
}

// PrintStatF1F0 prints statistical analysis for two-objective problems
//  emin, eave, emax, edev -- errors on f1(f0)
//  lmin, lave, lmax, ldev -- arc-lengths along f1(f0) curve
func (o *Optimiser) PrintStatF1F0() {
	if len(o.F1F0_err) == 0 && len(o.F1F0_arcLen) == 0 {
		io.Pf("there are no samples for statistical analysis\n")
		return
	}
	o.fix_formatting_data()
	io.Pf("\nerror on Pareto front\n")
	io.Pf("emin = %g\n", o.Emin)
	io.Pf("eave = %g\n", o.Eave)
	io.Pf("emax = %g\n", o.Emax)
	io.Pf("edev = %g\n", o.Edev)
	io.Pf(rnd.BuildTextHist(
		nice(o.Emin, o.HistNdig)-o.HistDelEmin,
		nice(o.Emax, o.HistNdig)+o.HistDelEmax,
		o.HistNsta, o.F1F0_err, o.HistFmt, o.HistLen))
	io.Pf("\nnormalised arc length along Pareto front (ref = %g)\n", o.F1F0_arcLenRef)
	io.Pf("lmin = %g\n", o.Lmin)
	io.Pf("lave = %g\n", o.Lave)
	io.Pf("lmax = %g\n", o.Lmax)
	io.Pf("ldev = %g\n", o.Ldev)
	io.Pf(rnd.BuildTextHist(
		nice(o.Lmin, o.HistNdig)-o.HistDelEmin,
		nice(o.Lmax, o.HistNdig)+o.HistDelEmax,
		o.HistNsta, o.F1F0_arcLen, o.HistFmt, o.HistLen))
}

// PrintStatMultiE prints statistical error analysis for multi-objective problems
func (o *Optimiser) PrintStatMultiE() {
	if len(o.Multi_err) < 2 {
		io.Pf("there are no samples for statistical analysis\n")
		return
	}
	o.fix_formatting_data()
	io.Pf("\nerror on Pareto front (multi)\n")
	io.Pf("Emin = %g\n", o.Emin)
	io.Pf("Eave = %g\n", o.Eave)
	io.Pf("Emax = %g\n", o.Emax)
	io.Pf("Edev = %g\n", o.Edev)
	io.Pf(rnd.BuildTextHist(
		nice(o.Emin, o.HistNdig)-o.HistDelEmin,
		nice(o.Emax, o.HistNdig)+o.HistDelEmax,
		o.HistNsta, o.Multi_err, o.HistFmt, o.HistLen))
}

// PrintStatIGD prints statistical IGD analysis for multi-objective problems
func (o *Optimiser) PrintStatIGD() {
	if len(o.Multi_IGD) < 2 {
		io.Pf("there are no samples for statistical analysis\n")
		return
	}
	o.fix_formatting_data()
	io.Pf("\nerror on Pareto front (multi)\n")
	io.Pf("IGDmin = %g\n", o.IGDmin)
	io.Pf("IGDave = %g\n", o.IGDave)
	io.Pf("IGDmax = %g\n", o.IGDmax)
	io.Pf("IGDdev = %g\n", o.IGDdev)
	io.Pf(rnd.BuildTextHist(
		nice(o.IGDmin, o.HistNdig)-o.HistDelEmin,
		nice(o.IGDmax, o.HistNdig)+o.HistDelEmax,
		o.HistNsta, o.Multi_IGD, o.HistFmt, o.HistLen))
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// calcIgd computes the IGD metric (smaller value means the Pareto front is wide and accurate).
//  fStar is a matrix with reference points [npoints][nova]
func (o *Optimiser) calcIgd(fStar [][]float64) (igd float64) {
	for _, point := range fStar {
		dmin := INF
		for _, sol := range o.Solutions {
			if sol.Feasible() {
				d := 0.0
				for j := 0; j < o.Nova; j++ {
					d += (point[j] - sol.Flt[j]) * (point[j] - sol.Flt[j])
				}
				if d < dmin {
					dmin = d
				}
			}
		}
		igd += math.Sqrt(dmin)
	}
	igd /= float64(len(fStar))
	return
}

// fix_formatting_data fixes formatting data and data for histograms
func (o *Stat) fix_formatting_data() {
	if o.RptFmtF == "" {
		o.RptFmtF = "%g"
	}
	if o.RptFmtFdev == "" {
		o.RptFmtFdev = "%g"
	}
	if o.RptFmtE == "" {
		o.RptFmtE = "%.8e"
	}
	if o.RptFmtEdev == "" {
		o.RptFmtEdev = "%.8e"
	}
	if o.RptFmtL == "" {
		o.RptFmtL = "%g"
	}
	if o.RptFmtLdev == "" {
		o.RptFmtLdev = "%.8e"
	}
	if o.RptFmtX == "" {
		o.RptFmtX = "%g"
	}
	if o.RptWordF == "" {
		o.RptWordF = "f"
	}
	if o.HistFmt == "" {
		o.HistFmt = "%.2f"
	}
	if math.Abs(o.HistDelFmin) < 1e-15 {
		o.HistDelFmin = 0.05
	}
	if math.Abs(o.HistDelFmax) < 1e-15 {
		o.HistDelFmax = 0.05
	}
	if math.Abs(o.HistDelEmin) < 1e-15 {
		o.HistDelEmin = 0.05
	}
	if math.Abs(o.HistDelEmax) < 1e-15 {
		o.HistDelEmax = 0.05
	}
	if o.HistDelFminZero {
		o.HistDelFmin = 0
	}
	if o.HistDelFmaxZero {
		o.HistDelFmax = 0
	}
	if o.HistDelEminZero {
		o.HistDelEmin = 0
	}
	if o.HistDelEmaxZero {
		o.HistDelEmax = 0
	}
	if o.HistNdig == 0 {
		o.HistNdig = 3
	}
	if o.HistNsta == 0 {
		o.HistNsta = 8
	}
	if o.HistLen == 0 {
		o.HistLen = 20
	}
}

// find_f0_spot finds where f0 falls in
func (o *Stat) find_f0_spot(f0 float64) (idx int) {
	for i, f0vals := range o.F1F0_f0ranges {
		if f0 >= f0vals[0] && f0 <= f0vals[1] {
			return i
		}
	}
	l := len(o.F1F0_f0ranges) - 1
	if f0 > o.F1F0_f0ranges[l][0] {
		return l
	}
	if f0 < o.F1F0_f0ranges[0][1] {
		return 0
	}
	return -1
}

// nice returns a truncated float
func nice(x float64, ndigits int) float64 {
	s := io.Sf("%."+io.Sf("%d", ndigits)+"f", x)
	return io.Atof(s)
}

// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
)

func main() {

	rnd.Init(0)

	// Simply supported beam
	// Analyse the max deflection at mid-span of simply supported beam
	// with uniform distributed load q and concentrated load at midspan
	//  Data:
	//   L    -- span
	//   EI   -- Young's modulus times cross-sectional moment of inertia
	//   p    -- x[0] concentrated load at mid-span
	//   q    -- x[1] distributed load
	//   δlim -- max deflection (vertical displacement) at mid-span
	// Reference
	//  Haldar A, Reliability-Based Structura Design, 2005

	// constants
	δlim := 0.0381 // [m] max allowed deflection
	L := 9.144     // [m] span
	EI := 182262.0 // [kN m²] flexural rigidity
	L3 := math.Pow(L, 3.0)

	// statistics of p=x[0] and q=x[1]
	μ := []float64{111.2, 35.03} // mean values
	σ := []float64{11.12, 5.25}  // deviation values
	lrv := []bool{true, false}   // is lognormal random variable?

	// limit state function
	gfcn := func(x []float64, args ...interface{}) (g float64, err error) {
		p, q := x[0], x[1]
		g = δlim - (p*L3/EI/48.0 + 5.0*q*L3*L/EI/384.0)
		return
	}

	// derivative of limit state function
	hfcn := func(dgdx, x []float64, args ...interface{}) (err error) {
		dgdx[0] = -L3 / EI / 48.0            // dg/dp
		dgdx[1] = -5.0 * L3 * L / EI / 384.0 // dg/dq
		return
	}

	// first order reliability method
	var form rnd.ReliabFORM
	form.Init(μ, σ, lrv, gfcn, hfcn)
	form.TolA = 0.005
	form.TolB = 0.005
	verbose := false // show messages
	βtrial := 3.0
	βform, _, _, xform := form.Run(βtrial, verbose)
	io.Pforan("βform = %v\n", βform)
	io.Pforan("xform = %v\n", xform)

	// objective function
	x := make([]float64, 2) // original random variables
	y := make([]float64, 2) // normalised random variables
	ovfunc := func(ind *goga.Individual, idIsland, time int, report *bytes.Buffer) (ova, oor float64) {
		x[0], x[1] = ind.GetFloat(0), ind.GetFloat(1)
		gx, err := gfcn(x)
		if err != nil {
			oor = 1e3
			return
		}
		for i := 0; i < len(x); i++ {
			y[i] = (x[i] - μ[i]) / σ[i]
		}
		ova = la.VecDot(y, y)
		oor = math.Abs(gx) // gx must be equal to zero
		return
	}

	// parameters
	C := goga.NewConfParams()
	C.Pll = false
	C.Nisl = 4
	C.Ninds = 20
	C.FnKey = "rel-simple-beam-form"
	C.DoPlot = true
	C.PltTi = 20
	C.CalcDerived()

	// bingo
	cf := 2.0
	xmin := []float64{
		μ[0] - cf*μ[0],
		μ[1] - cf*μ[1],
	}
	xmax := []float64{
		μ[0] + cf*μ[0],
		μ[1] + cf*μ[1],
	}
	bingo := goga.NewBingoFloats(xmin, xmax)

	// evolver
	evo := goga.NewEvolverFloatChromo(C, xmin, xmax, ovfunc, bingo)
	verbose = true
	doreport := true
	evo.Run(verbose, doreport)

	// results
	x[0], x[1] = evo.Best.GetFloat(0), evo.Best.GetFloat(1)
	for i := 0; i < len(x); i++ {
		y[i] = (x[i] - μ[i]) / σ[i]
	}
	β := math.Sqrt(la.VecDot(y, y))
	io.Pfgreen("\nx=%g\n", x)
	io.Pfgreen("β = %g\n", β)
	io.PfGreen("BestOV=%g\n", evo.Best.Ova)
}

// calc_normal_vars finds equivalent normal mean and std deviation for lognormal variables
func calc_normal_vars(μN, σN, x []float64, lrv []bool, μ, σ []float64) {
	var lnd rnd.DistLogNormal
	for i := 0; i < len(x); i++ {
		if lrv[i] {
			if true { // TODO: replace the following 2 hard-wired calculations
				lnd.Sig = σ[i] / μ[i]
				lnd.Mu = math.Log(μ[i]) - lnd.Sig*lnd.Sig/2.0
			} else {
				lnd.Init(μ[i], σ[i])
			}
			lnd.CalcDerived()
			Fi := lnd.Cdf(x[i])
			fi := lnd.Pdf(x[i])
			z := rnd.StdInvPhi(Fi)
			σN[i] = rnd.Stdphi(z) / fi
			μN[i] = x[i] - z*σN[i]
		}
	}
}

// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

// constants
const (
	INF = 1e+30 // infinite distance
)

// Generator_t defines callback function to generate trial solutions
type Generator_t func(sols []*Solution, prms *Parameters, reset bool)

// ObjFunc_t defines the objective fuction
type ObjFunc_t func(sol *Solution, cpu int)

// MinProb_t defines objective functon for specialised minimisation problem
type MinProb_t func(f, g, h, x []float64, y []int, cpu int)

// CxInt_t defines crossover function for ints
type CxInt_t func(a, b, A, B []int, prms *Parameters)

// MtInt_t defines mutation function for ints
type MtInt_t func(a []int, prms *Parameters)

// Output_t defines a function to perform output of data during the evolution
type Output_t func(time int, sols []*Solution)

// YfuncX_t defines the folowing simple function: y(x), usually used in plots
type YfuncX_t func(x float64) float64

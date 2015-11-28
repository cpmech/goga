// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

// constants
const (
	INF   = 1e+30 // infinite distance
	MDMIN = 0.001 // multiplier to compute minimum distance when flagging repeated solutions
)

// Generator_t defines callback function to generate trial solutions
type Generator_t func(sols []*Solution, prms *Parameters)

// ObjFunc_t defines the objective fluction
type ObjFunc_t func(sol *Solution, cpu int)

// MinProb_t defines objective functon for specialised minimisation problem
type MinProb_t func(f, g, h, x []float64, ξ []int, cpu int)

// CxFlt_t defines crossover function for floats
type CxFlt_t func(a, b, A, B, C, D, E, F []float64, prms *Parameters)

// CxInt_t defines crossover function for ints
type CxInt_t func(a, b, A, B, C, D, E, F []int, prms *Parameters)

// MtFlt_t defines mutation function for floats
type MtFlt_t func(a []float64, prms *Parameters)

// MtInt_t defines mutation function for ints
type MtInt_t func(a []int, prms *Parameters)

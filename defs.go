// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

// PostProc_t defines a function to post-process results
type PostProc_t func(idIsland, time int, pop Population)

// Func_t defines a type for a generic function to be used as a gene value
type Func_t func(ind *Individual) string

// Objectives_t defines the template for the objective functions and constraints.
// ind is set with ovas and oors
type Objectives_t func(isl int, ind *Individual)

// PopFltGen_t defines function to generate population of float point numbers
type PopFltGen_t func(isl int, C *ConfParams) Population

// PopIntGen_t defines function to generate population of integers
type PopIntGen_t func(isl int, C *ConfParams) Population

// crossover functions
type CxFltFunc_t func(a, b, A, B, C, D []float64, time int, dat *OpsData) (ends []int)
type CxIntFunc_t func(a, b, A, B, C, D []int, time int, dat *OpsData) (ends []int)

// mutation functions
type MtFltFunc_t func(a []float64, time int, dat *OpsData)
type MtIntFunc_t func(a []int, time int, dat *OpsData)

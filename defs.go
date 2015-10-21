// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "bytes"

// PostProc_t defines a function to post-process results
type PostProc_t func(idIsland, time int, pop Population)

// Func_t defines a type for a generic function to be used as a gene value
type Func_t func(ind *Individual) string

// Objectives_t defines the template for the objective functions and constraints.
// ind is set with ovas and oors
type Objectives_t func(ind *Individual, idIsland, time int, report *bytes.Buffer)

// PopIntGen_t defines function to generate population of integers
type PopIntGen_t func(isl int, C *ConfParams) Population

// PopFltGen_t defines function to generate population of float point numbers
type PopFltGen_t func(isl int, C *ConfParams) Population

// PopStrGen_t defines function to generate population of strings
type PopStrGen_t func(isl int, C *ConfParams) Population

// PopKeyGen_t defines function to generate population of keys (bytes)
type PopKeyGen_t func(isl int, C *ConfParams) Population

// PopBytGen_t defines function to generate population of bytes
type PopBytGen_t func(isl int, C *ConfParams) Population

// PopFunGen_t defines function to generate population of functions
type PopFunGen_t func(isl int, C *ConfParams) Population

// crossover functions
type CxIntFunc_t func(a, b, A, B, C, D []int, time int, dat *OpsData) (ends []int)
type CxFltFunc_t func(a, b, A, B, C, D []float64, time int, dat *OpsData) (ends []int)
type CxStrFunc_t func(a, b, A, B, C, D []string, time int, dat *OpsData) (ends []int)
type CxKeyFunc_t func(a, b, A, B, C, D []byte, time int, dat *OpsData) (ends []int)
type CxBytFunc_t func(a, b, A, B, C, D [][]byte, time int, dat *OpsData) (ends []int)
type CxFunFunc_t func(a, b, A, B, C, D []Func_t, time int, dat *OpsData) (ends []int)

// mutation functions
type MtIntFunc_t func(a []int, time int, dat *OpsData)
type MtFltFunc_t func(a []float64, time int, dat *OpsData)
type MtStrFunc_t func(a []string, time int, dat *OpsData)
type MtKeyFunc_t func(a []byte, time int, dat *OpsData)
type MtBytFunc_t func(a [][]byte, time int, dat *OpsData)
type MtFunFunc_t func(a []Func_t, time int, dat *OpsData)

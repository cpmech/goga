// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "bytes"

// Func_t defines a type for a generic function to be used as a gene value
type Func_t func(ind *Individual) string

// Objectives_t defines the template for the objective functions and constraints.
// ind is set with ovas and oors
type Objectives_t func(ind *Individual, idIsland, time int, report *bytes.Buffer)

// PopIntGen_t defines function to generate population of integers
type PopIntGen_t func(pop Population, ninds, nova, noor, nbases int, noise float64, args interface{}, irange [][]int) Population

// PopOrdGen_t defines function to generate population of ordered integers
type PopOrdGen_t func(pop Population, ninds, nova, noor, nbases int, noise float64, args interface{}, nints int) Population

// PopFltGen_t defines function to generate population of float point numbers
type PopFltGen_t func(pop Population, ninds, nova, noor, nbases int, noise float64, args interface{}, frange [][]float64) Population

// PopStrGen_t defines function to generate population of strings
type PopStrGen_t func(pop Population, ninds, nova, noor, nbases int, noise float64, args interface{}, pool [][]string) Population

// PopKeyGen_t defines function to generate population of keys (bytes)
type PopKeyGen_t func(pop Population, ninds, nova, noor, nbases int, noise float64, args interface{}, pool [][]byte) Population

// PopBytGen_t defines function to generate population of bytes
type PopBytGen_t func(pop Population, ninds, nova, noor, nbases int, noise float64, args interface{}, pool [][]string) Population

// PopFunGen_t defines function to generate population of functions
type PopFunGen_t func(pop Population, ninds, nova, noor, nbases int, noise float64, args interface{}, pool [][]Func_t) Population

// crossover functions
type CxIntFunc_t func(a, b, A, B []int, ncuts int, cuts []int, pc float64) (ends []int)
type CxFltFunc_t func(a, b, A, B []float64, ncuts int, cuts []int, pc float64) (ends []int)
type CxStrFunc_t func(a, b, A, B []string, ncuts int, cuts []int, pc float64) (ends []int)
type CxKeyFunc_t func(a, b, A, B []byte, ncuts int, cuts []int, pc float64) (ends []int)
type CxBytFunc_t func(a, b, A, B [][]byte, ncuts int, cuts []int, pc float64) (ends []int)
type CxFunFunc_t func(a, b, A, B []Func_t, ncuts int, cuts []int, pc float64) (ends []int)

// mutation functions
type MtIntFunc_t func(a []int, nchanges int, pm float64, extra interface{})
type MtFltFunc_t func(a []float64, nchanges int, pm float64, extra interface{})
type MtStrFunc_t func(a []string, nchanges int, pm float64, extra interface{})
type MtKeyFunc_t func(a []byte, nchanges int, pm float64, extra interface{})
type MtBytFunc_t func(a [][]byte, nchanges int, pm float64, extra interface{})
type MtFunFunc_t func(a []Func_t, nchanges int, pm float64, extra interface{})

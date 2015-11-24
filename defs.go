// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

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

// ParetoF1F0_t Pareto front solution
type ParetoF1F0_t func(f0 float64) float64

// MinProblem_t defines the minimisation problem. See Optimiser
type MinProblem_t func(f, g, h, x []float64, isl int)

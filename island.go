// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "math"

// ObjFunc_t defines the template for the objective function
type ObjFunc_t func(ind *Individual, time int, best *Individual)

// Control holds control parameters
type Control struct {
	UseRanking  bool
	RnkPressure float64
}

// Island holds one population and performs the reproduction operation
type Island struct {

	// control parameters
	C Control

	// population
	Pop     Population // pointer to current population
	PopA    Population // population holder
	PopB    Population // population holder
	ObjFunc ObjFunc_t  // objective function

	// auxiliary internal data
	fitsrnk []float64   // all fitness values computed by ranking
	fitness []float64   // all fitness values
	prob    []float64   // probabilities
	cumprob []float64   // cumulated probabilities
	oldbest *Individual // best individual (copy)
}

// NewIsland allocates a new island but with a give population already allocated
// Input:
//  pop   -- the population
//  ofunc -- objective function
func NewIsland(pop Population, ofunc ObjFunc_t) (o *Island) {

	// allocate
	o = new(Island)
	o.Pop = pop
	o.PopA = pop
	o.ObjFunc = ofunc

	// set default control values
	o.C.UseRanking = true
	o.C.RnkPressure = 1.2

	// compute objective values
	for _, ind := range o.Pop {
		o.ObjFunc(ind, 0, nil)
	}

	// sort
	o.Pop.Sort()

	// auxiliary data
	ninds := len(o.Pop)
	o.fitsrnk = make([]float64, ninds)
	o.fitness = make([]float64, ninds)
	o.prob = make([]float64, ninds)
	o.cumprob = make([]float64, ninds)
	o.oldbest = o.Pop[0].GetCopy()
	return
}

func (o *Island) Reproduction(time int) {

	// best individual: Note: population must be sorted already
	o.Pop[0].CopyInto(o.oldbest)

	// fitness and probabilities
	ninds := len(o.Pop)
	sumfit := 0.0
	if o.C.UseRanking {
		sp := o.C.RnkPressure
		if sp < 1.0 || sp > 2.0 {
			sp = 1.2
		}
		for i := 0; i < ninds; i++ {
			o.fitness[i] = 2.0 - sp + 2.0*(sp-1.0)*float64(ninds-i-1)/float64(ninds-1)
			sumfit += o.fitness[i]
		}
	} else {
		ovmin, ovmax := o.Pop[0].ObjValue, o.Pop[0].ObjValue
		for _, ind := range o.Pop {
			ovmin = min(ovmin, ind.ObjValue)
			ovmax = max(ovmax, ind.ObjValue)
		}
		if math.Abs(ovmax-ovmin) < 1e-14 {
			for i := 0; i < ninds; i++ {
				o.fitness[i] = float64(i) / float64(ninds-1)
				sumfit += o.fitness[i]
			}
		} else {
			for i, ind := range o.Pop {
				o.fitness[i] = (ovmax - ind.ObjValue) / (ovmax - ovmin)
				sumfit += o.fitness[i]
			}
		}
	}
	for i := 0; i < ninds; i++ {
		o.prob[i] = o.fitness[i] / sumfit
	}
	CumSum(o.cumprob, o.prob)
}

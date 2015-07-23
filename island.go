// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

// ObjFunc_t defines the template for the objective function
type ObjFunc_t func(ind *Individual, time int, best *Individual)

// Island holds one population and performs the reproduction operation
type Island struct {
	Pop     Population  // pointer to current population
	PopA    Population  // population holder
	PopB    Population  // population holder
	ObjFunc ObjFunc_t   // objective function
	BestInd *Individual // best individual (copy)
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

	// compute objective values
	for _, ind := range o.Pop {
		o.ObjFunc(ind, 0, nil)
	}

	// sort
	o.Pop.Sort()

	// best individual (copy)
	o.BestInd = o.Pop[0].GetCopy()
	return
}

func (o *Island) Reproduction(time int) {

}

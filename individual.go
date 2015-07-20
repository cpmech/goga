// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

// Individual implements one individual in a population
type Individual struct {
	Chromo   []*Gene // chromosome [ngenes*nbases]
	ObjValue float64 // objective value
	Fitness  float64 // fitness
}

// InitChromo initialises chromosome with all genes
//  Input:
//   nbases -- used to split genes of floats into smaller parts
//   slices -- slices of ints, floats, strings, bytes, and Func_t
//  Notes:
//   1) the slices in 'genes' can all be combined to define genes with mixed data;
//   2) the slices can also be nil, except for one of them.
//  Example
func (o *Individual) InitChromo(nbases int, slices ...interface{}) {

	// auxiliary function
	newgenes := func(ngenes int) {
		o.Chromo = make([]*Gene, ngenes)
		for i := 0; i < ngenes; i++ {
			o.Chromo[i] = NewGene(nbases)
		}
	}

	// set genes
	ngenes := 0
	for _, slice := range slices {
		switch s := slice.(type) {
		case []int:
			if ngenes < 1 {
				ngenes = len(s)
				newgenes(ngenes)
			}
			for i, value := range s {
				o.Chromo[i].SetInt(value)
			}
		case []float64:
			if ngenes < 1 {
				ngenes = len(s)
				newgenes(ngenes)
			}
			for i, value := range s {
				o.Chromo[i].SetFloat(value)
			}
		case []string:
			if ngenes < 1 {
				ngenes = len(s)
				newgenes(ngenes)
			}
			for i, value := range s {
				o.Chromo[i].SetString(value)
			}
		case []byte:
			if ngenes < 1 {
				ngenes = len(s)
				newgenes(ngenes)
			}
			for i, value := range s {
				o.Chromo[i].SetByte(value)
			}
		case [][]byte:
			if ngenes < 1 {
				ngenes = len(s)
				newgenes(ngenes)
			}
			for i, value := range s {
				o.Chromo[i].SetBytes(value)
			}
		case []Func_t:
			if ngenes < 1 {
				ngenes = len(s)
				newgenes(ngenes)
			}
			for i, value := range s {
				o.Chromo[i].SetFunc(value)
			}
		}
	}
}

func (o Individual) Output(fmtInt, fmtFloat, fmtString, fmtBytes string) (l string) {
	l = "("
	for i, g := range o.Chromo {
		if i > 0 {
			l += ") ("
		}
		l += g.Output(fmtInt, fmtFloat, fmtString, fmtBytes)
	}
	l += ")"
	return
}

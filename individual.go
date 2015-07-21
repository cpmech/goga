// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/io"

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

// GetCopy returns a copy of this individual
func (o Individual) GetCopy() (x *Individual) {
	x = new(Individual)
	ngenes := len(o.Chromo)
	x.Chromo = make([]*Gene, ngenes)
	for i := 0; i < ngenes; i++ {
		x.Chromo[i] = o.Chromo[i].GetCopy()
	}
	x.ObjValue = o.ObjValue
	x.Fitness = o.Fitness
	return
}

// output //////////////////////////////////////////////////////////////////////////////////////////

// GetStringSizes returns the sizes of strings represent each gene type
//  sizes -- []int{szInt, szFloat, szString, szBytes}
func (o Individual) GetStringSizes() (sizes []int) {
	sizes = make([]int, 4)
	for _, g := range o.Chromo {
		if g.Int != nil {
			sizes[0] = imax(sizes[0], len(io.Sf("%d", g.GetInt())))
		}
		if g.Flt != nil {
			sizes[1] = imax(sizes[1], len(io.Sf("%g", g.GetFloat())))
		}
		if g.String != nil {
			sizes[2] = imax(sizes[2], len(io.Sf("%s", g.GetString())))
		}
		if g.Bytes != nil {
			sizes[3] = imax(sizes[3], len(io.Sf("%s", string(g.GetBytes()))))
		}
	}
	return
}

// Output returns a string representation of this individual
//  fmts -- []string{formatInt, formatFloat, formatString, formatBytes}
func (o Individual) Output(fmts []string) (l string) {
	if len(o.Chromo) < 1 {
		return
	}
	nfields := o.Chromo[0].Nfields()
	if nfields > 1 {
		l = "("
	}
	for i, g := range o.Chromo {
		if i > 0 && nfields > 1 {
			l += ") ("
		}
		l += g.Output(fmts)
	}
	if nfields > 1 {
		l += ")"
	}
	return
}

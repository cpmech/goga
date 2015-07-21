// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

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

// get methods /////////////////////////////////////////////////////////////////////////////////////

func (o Individual) CountBases() (nint, nflt, nstr, nbyt, nbytes, nfuncs int) {
	for _, g := range o.Chromo {

		// ints
		if g.Int != nil {
			nint++
		}

		// floats
		nbases := len(g.Fbases)
		if nbases > 1 {
			nflt += nbases
		} else {
			if g.Flt != nil {
				nflt++
			}
		}

		// strings
		if g.String != nil {
			nstr++
		}

		// byte
		if g.Byte != nil {
			nbyt++
		}

		// bytes
		nbytes += len(g.Bytes)

		// functions
		if g.Func != nil {
			nfuncs++
		}
	}
	return
}

// GetBases returns all bases from all genes
func (o Individual) GetBases(ints []int, flts []float64, strs []string, byts []byte, bytes []byte, funcs []Func_t) {
	var kint, kflt, kstr, kbyt, kbytes, kfuncs int
	for _, g := range o.Chromo {

		// ints
		if g.Int != nil {
			ints[kint] = *g.Int
			kint++
		}

		// floats
		nbases := len(g.Fbases)
		if nbases > 1 {
			for j := 0; j < nbases; j++ {
				flts[kflt] = g.Fbases[j]
				kflt++
			}
		} else {
			if g.Flt != nil {
				flts[kflt] = *g.Flt
				kflt++
			}
		}

		// strings
		if g.String != nil {
			strs[kstr] = *g.String
			kstr++
		}

		// byte
		if g.Byte != nil {
			byts[kbyt] = *g.Byte
			kbyt++
		}

		// bytes
		for j := 0; j < len(g.Bytes); j++ {
			bytes[kbytes] = g.Bytes[j]
			kbytes++
		}

		// functions
		if g.Func != nil {
			funcs[kfuncs] = g.Func
			kfuncs++
		}
	}
}

// genetic algorithm routines //////////////////////////////////////////////////////////////////////

// output //////////////////////////////////////////////////////////////////////////////////////////

// GetStringSizes returns the sizes of strings represent each gene type
//  sizes -- [ngenes] sizes of strings for {int, flt, string, byte, bytes, func}
func (o Individual) GetStringSizes() (sizes [][]int) {
	ngenes := len(o.Chromo)
	sizes = utl.IntsAlloc(ngenes, 6)
	for i, g := range o.Chromo {
		if g.Int != nil {
			sizes[i][0] = imax(sizes[i][0], len(io.Sf("%v", g.GetInt())))
		}
		if g.Flt != nil {
			sizes[i][1] = imax(sizes[i][1], len(io.Sf("%v", g.GetFloat())))
		}
		if g.String != nil {
			sizes[i][2] = imax(sizes[i][2], len(io.Sf("%v", g.GetString())))
		}
		if g.Byte != nil {
			sizes[i][3] = imax(sizes[i][3], len(io.Sf("%v", g.GetByte())))
		}
		if g.Bytes != nil {
			sizes[i][4] = imax(sizes[i][4], len(io.Sf("%v", string(g.GetBytes()))))
		}
		if g.Func != nil {
			sizes[i][5] = imax(sizes[i][5], len(io.Sf("%v", g.GetFunc()(g))))
		}
	}
	return
}

// Output returns a string representation of this individual
//  fmts -- [ngenes] formats of strings for {int, flt, string, byte, bytes, func}
//          use fmts == nil to choose default ones
func (o Individual) Output(fmts [][]string) (l string) {
	ngenes := len(o.Chromo)
	if ngenes < 1 {
		return
	}
	for i, g := range o.Chromo {
		if i > 0 {
			l += " "
		}
		nfields := g.Nfields()
		if nfields > 1 {
			l += "["
		}
		if len(fmts) == ngenes {
			l += g.Output(fmts[i])
		} else {
			if len(fmts) == 1 {
				l += g.Output(fmts[0])
			} else {
				l += g.Output(nil)
			}
		}
		if nfields > 1 {
			l += "]"
		}
	}
	return
}

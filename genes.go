// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
)

// Func_t defines a type for a generic function to be used as a gene value
type Func_t func(g *Gene) string

// Gene defines the gene type
type Gene struct {
	Int    *int      // int gene
	Flt    *float64  // float64 gene
	Fbases []float64 // subdivisions of float64 gene == bases
	String *string   // string gene
	Byte   *byte     // byte gene
	Bytes  []byte    // bytes gene
	Func   Func_t    // function gene
}

// NewGene allocates a new gene
func NewGene(nbases int) *Gene {
	gene := new(Gene)
	if nbases > 1 {
		gene.Fbases = make([]float64, nbases)
	}
	return gene
}

// GetCopy returns a copy of this gene
func (o Gene) GetCopy() (x *Gene) {
	nbases := len(o.Fbases)
	x = NewGene(nbases)
	if o.Int != nil {
		x.SetInt(*o.Int)
	}
	if o.Flt != nil {
		x.SetFloat(*o.Flt)
		if nbases > 1 {
			copy(x.Fbases, o.Fbases)
		}
	}
	if o.String != nil {
		x.SetString(*o.String)
	}
	if o.Byte != nil {
		x.SetByte(*o.Byte)
	}
	if o.Bytes != nil {
		x.SetBytes(o.Bytes)
	}
	if o.Func != nil {
		x.SetFunc(o.Func)
	}
	return
}

// genetic algorithm operators /////////////////////////////////////////////////////////////////////

// TODO
//func (o *Gene)

// set methods /////////////////////////////////////////////////////////////////////////////////////

// SetInt sets an integer as gene value
func (o *Gene) SetInt(value int) {
	if o.Int == nil {
		o.Int = new(int)
	}
	*o.Int = value
}

// SetFloat sets a float point number as gene value
//  Note: if nbases > 1, basis values will be randomly computed
func (o *Gene) SetFloat(value float64) {
	if o.Flt == nil {
		o.Flt = new(float64)
	}
	*o.Flt = value
	nbases := len(o.Fbases)
	if nbases > 1 {
		rnd.Float64s(o.Fbases, 0, 1)
		sum := la.VecAccum(o.Fbases)
		for j := 0; j < nbases; j++ {
			o.Fbases[j] = value * o.Fbases[j] / sum
		}
	}
}

// SetFbases sets sub-floats (divisions of Float)
//  Input:
//   start  -- start position in SubFloats
//   values -- values to be copied into SubFloats
//  Example:
//   SubFloats (before) = [0, 1, 2, 3, 4, 5]
//   values             =       [6, 7, 8]
//   start              =        2
//   SubFloats (after   = [0, 1, 6, 7, 8, 5]
//  Note: Float will be computed accordingly; i.e. Float = sum(SubFloats)
func (o *Gene) SetFbases(start int, values []float64) {
	nbases := len(o.Fbases)
	if nbases < 2 {
		if len(values) > 0 {
			*o.Flt = values[0]
			return
		}
		return
	}
	chk.IntAssertLessThan(start, nbases)
	chk.IntAssertLessThan(len(values), nbases+1)
	for i, v := range values {
		o.Fbases[start+i] = v
	}
	*o.Flt = la.VecAccum(o.Fbases)
}

// SetString sets a string as gene value
func (o *Gene) SetString(value string) {
	if o.String == nil {
		o.String = new(string)
	}
	*o.String = value
}

// SetByte sets a byte as gene value
func (o *Gene) SetByte(value byte) {
	if o.Byte == nil {
		o.Byte = new(byte)
	}
	*o.Byte = value
}

// SetBytes sets a slice of bytes as gene value
func (o *Gene) SetBytes(value []byte) {
	if len(o.Bytes) != len(value) {
		o.Bytes = make([]byte, len(value))
	}
	copy(o.Bytes, value)
}

// SetFunc sets a function as gene value
func (o *Gene) SetFunc(value Func_t) {
	o.Func = value
}

// get methods /////////////////////////////////////////////////////////////////////////////////////

// GetInt returns the int value, if any
func (o Gene) GetInt() int {
	if o.Int != nil {
		return *o.Int
	}
	return 0
}

// GetFloat returns the float point number value, if any
func (o Gene) GetFloat() float64 {
	if o.Flt != nil {
		return *o.Flt
	}
	return 0
}

// GetString returns the string value, if any
func (o Gene) GetString() string {
	if o.String != nil {
		return *o.String
	}
	return ""
}

// GetByte returns the byte value, if any
func (o Gene) GetByte() byte {
	if o.Byte != nil {
		return *o.Byte
	}
	return 0
}

// GetBytes returns the slice of bytes, if any
func (o Gene) GetBytes() []byte {
	return o.Bytes
}

// GetFunc returns the function, if any
func (o Gene) GetFunc() Func_t {
	return o.Func
}

// output //////////////////////////////////////////////////////////////////////////////////////////

// Output returns a string representation of this gene
//  fmts -- formats for     int,     flt, string, byte,  bytes, and func
//          use fmts == nil to choose default ones
func (o *Gene) Output(fmts []string) (l string) {
	if len(fmts) != 6 {
		fmts = []string{" %4d", " %8.3f", " %6.6s", " %x", " %6.6s", " %6.6s"}
	}
	if o.Int != nil {
		l += io.Sf(fmts[0], *o.Int)
	}
	if o.Flt != nil {
		l += io.Sf(fmts[1], *o.Flt)
	}
	if o.String != nil {
		l += io.Sf(fmts[2], *o.String)
	}
	if o.Byte != nil {
		l += io.Sf(fmts[3], *o.Byte)
	}
	if o.Bytes != nil {
		l += io.Sf(fmts[4], string(o.Bytes))
	}
	if o.Func != nil {
		l += io.Sf(fmts[5], o.Func(o))
	}
	return
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func (o Gene) Nfields() (n int) {
	if o.Int != nil {
		n++
	}
	if o.Flt != nil {
		n++
	}
	if o.String != nil {
		n++
	}
	if o.Byte != nil {
		n++
	}
	if o.Bytes != nil {
		n++
	}
	if o.Func != nil {
		n++
	}
	return
}

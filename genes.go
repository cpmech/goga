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
type Func_t func() string

// Gene defines the gene type
type Gene struct {
	Int       *int      // int gene
	Float     *float64  // float64 gene
	SubFloats []float64 // subdivisions of float64 gene == bases
	String    *string   // string gene
	Byte      *byte     // byte gene
	Bytes     []byte    // bytes gene
	Func      Func_t    // function gene
}

// NewGene allocates a new gene
func NewGene(nbases int) *Gene {
	gene := new(Gene)
	if nbases > 1 {
		gene.SubFloats = make([]float64, nbases)
	}
	return gene
}

// GetCopy returns a copy of this gene
func (o Gene) GetCopy() (x *Gene) {
	nbases := len(o.SubFloats)
	x = NewGene(nbases)
	if o.Int != nil {
		x.SetInt(*o.Int)
	}
	if o.Float != nil {
		x.SetFloat(*o.Float)
		if nbases > 1 {
			copy(x.SubFloats, o.SubFloats)
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
	if o.Float == nil {
		o.Float = new(float64)
	}
	*o.Float = value
	nbases := len(o.SubFloats)
	if nbases > 1 {
		rnd.Float64s(o.SubFloats, 0, 1)
		sum := la.VecAccum(o.SubFloats)
		for j := 0; j < nbases; j++ {
			o.SubFloats[j] = value * o.SubFloats[j] / sum
		}
	}
}

// SetSubFloats sets sub-floats (divisions of Float)
//  Input:
//   start  -- start position in SubFloats
//   values -- values to be copied into SubFloats
//  Example:
//   SubFloats (before) = [0, 1, 2, 3, 4, 5]
//   values             =       [6, 7, 8]
//   start              =        2
//   SubFloats (after   = [0, 1, 6, 7, 8, 5]
//  Note: Float will be computed accordingly; i.e. Float = sum(SubFloats)
func (o *Gene) SetSubFloats(start int, values []float64) {
	nbases := len(o.SubFloats)
	if nbases < 2 {
		if len(values) > 0 {
			*o.Float = values[0]
			return
		}
		return
	}
	chk.IntAssertLessThan(start, nbases)
	chk.IntAssertLessThan(len(values), nbases+1)
	for i, v := range values {
		o.SubFloats[start+i] = v
	}
	*o.Float = la.VecAccum(o.SubFloats)
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
	if o.Float != nil {
		return *o.Float
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
//  fmts -- []string{formatInt, formatFloat, formatString, formatBytes}
func (o Gene) Output(fmts []string) (l string) {
	comma := ","
	if o.Nfields() == 1 {
		comma = ""
	}
	if o.Int != nil {
		l += io.Sf(fmts[0], *o.Int)
	}
	if o.Float != nil {
		l += io.Sf(comma+fmts[1], *o.Float)
	}
	if o.String != nil {
		l += io.Sf(comma+fmts[2], *o.String)
	}
	if o.Byte != nil {
		l += io.Sf(comma+"%x", *o.Byte)
	}
	if o.Bytes != nil {
		l += io.Sf(comma+fmts[3], string(o.Bytes))
	}
	if o.Func != nil {
		l += comma + o.Func()
	}
	return
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func (o Gene) Nfields() (n int) {
	if o.Int != nil {
		n++
	}
	if o.Float != nil {
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

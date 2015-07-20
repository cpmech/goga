// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
)

// Func_t defines a type for a generic function to be used as a gene value
type Func_t func() string

// Gene defines the gene type
type Gene struct {
	Int      *int      // int gene
	Float    *float64  // float64 gene
	SubFloat []float64 // subdivisions of float64 gene == bases
	String   *string   // string gene
	Byte     *byte     // byte gene
	Bytes    []byte    // bytes gene
	Func     Func_t    // function gene
}

// NewGene allocates a new gene
func NewGene(nbases int) *Gene {
	gene := new(Gene)
	if nbases > 1 {
		gene.SubFloat = make([]float64, nbases)
	}
	return gene
}

// set methods /////////////////////////////////////////////////////////////////////////////////////

// SetInt sets an integer as gene value
func (o *Gene) SetInt(value int) {
	if o.Int == nil {
		o.Int = new(int)
	}
	*o.Int = value
}

// SetFloat sets a float point number as gene value
func (o *Gene) SetFloat(value float64) {
	if o.Float == nil {
		o.Float = new(float64)
	}
	*o.Float = value
	nbases := len(o.SubFloat)
	if nbases > 1 {
		rnd.Float64s(o.SubFloat, 0, 1)
		sum := la.VecAccum(o.SubFloat)
		for j := 0; j < nbases; j++ {
			o.SubFloat[j] = value * o.SubFloat[j] / sum
		}
	}
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
func (o Gene) Output(fmtInt, fmtFloat, fmtString, fmtBytes string) (l string) {
	comma := ","
	if o.Nfields() == 1 {
		comma = ""
	}
	if o.Int != nil {
		l += io.Sf(fmtInt, *o.Int)
	}
	if o.Float != nil {
		l += io.Sf(comma+fmtFloat, *o.Float)
	}
	if o.String != nil {
		l += io.Sf(comma+fmtString, *o.String)
	}
	if o.Byte != nil {
		l += io.Sf(comma+"%x", *o.Byte)
	}
	if o.Bytes != nil {
		l += io.Sf(comma+fmtBytes, string(o.Bytes))
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

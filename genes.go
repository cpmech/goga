// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
)

type Func_t func() string

type Gene struct {
	Int      *int      // int gene
	Float    *float64  // float64 gene
	SubFloat []float64 // subdivisions of float64 gene == bases
	String   *string   // string gene
	Byte     *byte     // byte gene
	Bytes    []byte    // bytes gene
	Func     Func_t    // function gene
}

func NewGene(nbases int) *Gene {
	gene := new(Gene)
	if nbases > 1 {
		gene.SubFloat = make([]float64, nbases)
	}
	return gene
}

// set methods /////////////////////////////////////////////////////////////////////////////////////

func (o *Gene) SetInt(value int) {
	if o.Int == nil {
		o.Int = new(int)
	}
	*o.Int = value
}

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

func (o *Gene) SetString(value string) {
	if o.String == nil {
		o.String = new(string)
	}
	*o.String = value
}

func (o *Gene) SetByte(value byte) {
	if o.Byte == nil {
		o.Byte = new(byte)
	}
	*o.Byte = value
}

func (o *Gene) SetBytes(value []byte) {
	if len(o.Bytes) != len(value) {
		o.Bytes = make([]byte, len(value))
	}
	copy(o.Bytes, value)
}

func (o *Gene) SetFunc(value Func_t) {
	o.Func = value
}

// get methods /////////////////////////////////////////////////////////////////////////////////////

func (o Gene) GetInt() int {
	if o.Int != nil {
		return *o.Int
	}
	return 0
}

func (o Gene) GetFloat() float64 {
	if o.Float != nil {
		return *o.Float
	}
	return 0
}

func (o Gene) GetString() string {
	if o.String != nil {
		return *o.String
	}
	return ""
}

func (o Gene) GetByte() byte {
	if o.Byte != nil {
		return *o.Byte
	}
	return 0
}

func (o Gene) GetBytes() []byte {
	return o.Bytes
}

func (o Gene) GetFunc() Func_t {
	return o.Func
}

// output //////////////////////////////////////////////////////////////////////////////////////////

func (o Gene) Output(fmtInt, fmtFloat, fmtString, fmtBytes string) (l string) {
	if o.Int != nil {
		l += io.Sf(fmtInt, *o.Int)
	}
	if o.Float != nil {
		l += io.Sf(","+fmtFloat, *o.Float)
	}
	if o.String != nil {
		l += io.Sf(","+fmtString, *o.String)
	}
	if o.Byte != nil {
		l += io.Sf(",%x", *o.Byte)
	}
	if o.Bytes != nil {
		l += io.Sf(","+fmtBytes, string(o.Bytes))
	}
	if o.Func != nil {
		l += "," + o.Func()
	}
	return
}

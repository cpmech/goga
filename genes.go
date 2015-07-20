// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
)

type Gene struct {
	Int      *int      // int gene
	Float    *float64  // float64 gene
	SubFloat []float64 // subdivisions of float64 gene == bases
	String   *string   // string gene
}

func NewGene(nbases int, values ...interface{}) *Gene {
	gene := new(Gene)
	for _, value := range values {
		switch v := value.(type) {
		case int:
			gene.SetInt(v)
		case float64:
			if nbases > 1 {
				gene.SubFloat = make([]float64, nbases)
			}
			gene.SetFloat(v)
		case string:
			gene.SetString(v)
		default:
			return nil
		}
	}
	return gene
}

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
	if len(o.SubFloat) > 1 {
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

func (o Gene) GetRepresentation(szInt, szFloat, szString int) string {
	return ""
}

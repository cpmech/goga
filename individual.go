// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/chk"

// Individual implements one individual in a population
type Individual struct {
	Chromo   []*Gene // chromosome [ngenes*nbases]
	ObjValue float64 // objective value
	Fitness  float64 // fitness
}

func (o *Individual) InitIntChromo(nbases int, genes []int) (err error) {
	o.Chromo = make([]*Gene, len(genes))
	for i, g := range o.Chromo {
		g = NewGene(nbases, genes[i])
		if g == nil {
			return chk.Err("cannot create chromosome of integers")
		}
	}
	return
}

func (o *Individual) InitFloatChromo(nbases int, genes []float64) (err error) {
	o.Chromo = make([]*Gene, len(genes))
	for i, g := range o.Chromo {
		g = NewGene(nbases, genes[i])
		if g == nil {
			return chk.Err("cannot create chromosome of float point numbers")
		}
	}
	return
}

func (o *Individual) InitStringChromo(nbases int, genes []string) (err error) {
	o.Chromo = make([]*Gene, len(genes))
	for i, g := range o.Chromo {
		g = NewGene(nbases, genes[i])
		if g == nil {
			return chk.Err("cannot create chromosome of float point numbers")
		}
	}
	return
}

// ints, floats and strings can be nil
func (o *Individual) InitMixedChromo(nbases int, ints []int, floats []float64, strings []string) (err error) {

	// check
	if ints == nil && floats == nil && strings == nil {
		return chk.Err("at least one of the 'ints', 'floats' or 'strings' slices must be non-nil")
	}

	// strings only
	if ints == nil && floats == nil {
		o.Chromo = make([]*Gene, len(strings))
		for i, g := range o.Chromo {
			g = NewGene(nbases, strings[i])
			if g == nil {
				return chk.Err("cannot create chromosome of mixed gene types: strings only")
			}
		}
		return
	}

	// floats only
	if ints == nil && strings == nil {
		o.Chromo = make([]*Gene, len(floats))
		for i, g := range o.Chromo {
			g = NewGene(nbases, floats[i])
			if g == nil {
				return chk.Err("cannot create chromosome of mixed gene types: floats only")
			}
		}
		return
	}

	// ints only
	if floats == nil && strings == nil {
		o.Chromo = make([]*Gene, len(ints))
		for i, g := range o.Chromo {
			g = NewGene(nbases, ints[i])
			if g == nil {
				return chk.Err("cannot create chromosome of mixed gene types: ints only")
			}
		}
		return
	}

	// floats and strings
	if ints == nil {
		chk.IntAssert(len(floats), len(strings))
		o.Chromo = make([]*Gene, len(floats))
		for i, g := range o.Chromo {
			g = NewGene(nbases, floats[i], strings[i])
			if g == nil {
				return chk.Err("cannot create chromosome of mixed gene types. floats and strings")
			}
		}
		return
	}

	// ints and strings
	if floats == nil {
		chk.IntAssert(len(ints), len(strings))
		o.Chromo = make([]*Gene, len(ints))
		for i, g := range o.Chromo {
			g = NewGene(nbases, ints[i], strings[i])
			if g == nil {
				return chk.Err("cannot create chromosome of mixed gene types. ints and strings")
			}
		}
		return
	}

	// ints and floats
	if strings == nil {
		chk.IntAssert(len(ints), len(floats))
		o.Chromo = make([]*Gene, len(ints))
		for i, g := range o.Chromo {
			g = NewGene(nbases, ints[i], floats[i])
			if g == nil {
				return chk.Err("cannot create chromosome of mixed gene types. ints and floats")
			}
		}
	}
	return
}

// String returns a table-row representation of an individual
//  Input:
//   ovfmt -- objective value formatting string; use "" to skip this item
//   ftfmt -- fitness formatting string; use "" to skip this item
//   gefmt -- genes formatting string; use "" to skip this item
//   bsfmt -- bases formatting string; use "" to skip this item
func (o Individual) String(ovfmt, ftfmt, gefmt, bsfmt string) (line string) {
	/*
		if ovfmt != "" {
			line += io.Sf(ovfmt, o.ObjValue)
		}
		if ftfmt != "" {
			line += io.Sf(ftfmt, o.Fitness)
		}
		if gefmt != "" {
			for _, g := range o.Genes {
				line += io.Sf(gefmt, g)
			}
		}
		if bsfmt != "" {
			for _, v := range o.Chromo {
				line += io.Sf(bsfmt, v)
			}
		}
	*/
	return
}

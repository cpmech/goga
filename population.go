// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"math"
	"math/rand"

	"github.com/cpmech/gosl/rnd"
)

// Population holds all individuals
type Population []*Individual

// GetCopy returns a copy of this population
func (o Population) GetCopy() (pop Population) {
	ninds := len(o)
	pop = make([]*Individual, ninds)
	for i := 0; i < ninds; i++ {
		pop[i] = o[i].GetCopy()
	}
	return
}

// generation functions ////////////////////////////////////////////////////////////////////////////

// PopBinGen generates a population of binary numbers [0,1]
func PopBinGen(id int, C *ConfParams) Population {
	o := make([]*Individual, C.Ninds)
	genes := make([]int, C.NumInts)
	for i := 0; i < C.Ninds; i++ {
		for j := 0; j < C.NumInts; j++ {
			genes[j] = rand.Intn(2)
		}
		o[i] = NewIndividual(C.Nova, C.Noor, nil, genes)
	}
	return o
}

// PopOrdGen generates a population of individuals with ordered integers
// Notes: (1) ngenes = C.NumInts
func PopOrdGen(id int, C *ConfParams) Population {
	o := make([]*Individual, C.Ninds)
	ngenes := C.NumInts
	for i := 0; i < C.Ninds; i++ {
		o[i] = NewIndividual(C.Nova, C.Noor, nil, make([]int, ngenes))
		for j := 0; j < C.NumInts; j++ {
			o[i].Ints[j] = j
		}
		rnd.IntShuffle(o[i].Ints)
	}
	return o
}

// PopFltAndBinGen generates population with floats and integers
func PopFltAndBinGen(id int, C *ConfParams) Population {
	o := PopFltGen(id, C)
	for i := 0; i < C.Ninds; i++ {
		o[i].Ints = make([]int, C.NumInts)
		for j := 0; j < C.NumInts; j++ {
			o[i].Ints[j] = rand.Intn(2)
		}
	}
	return o
}

// PopFltGen generates a population of individuals with float point numbers
// Notes: (1) ngenes = len(C.RangeFlt)
func PopFltGen(id int, C *ConfParams) Population {
	o := make([]*Individual, C.Ninds)
	ngenes := len(C.RangeFlt)
	for i := 0; i < C.Ninds; i++ {
		o[i] = NewIndividual(C.Nova, C.Noor, make([]float64, ngenes), nil)
	}
	if C.Latin {
		K := rnd.LatinIHS(ngenes, C.Ninds, C.LatinDf)
		dx := make([]float64, ngenes)
		for i := 0; i < ngenes; i++ {
			dx[i] = (C.RangeFlt[i][1] - C.RangeFlt[i][0]) / float64(C.Ninds-1)
		}
		for i := 0; i < ngenes; i++ {
			for j := 0; j < C.Ninds; j++ {
				o[j].Floats[i] = C.RangeFlt[i][0] + float64(K[i][j]-1)*dx[i]
			}
		}
		return o
	}
	npts := int(math.Pow(float64(C.Ninds), 1.0/float64(ngenes))) // num points in 'square' grid
	ntot := int(math.Pow(float64(npts), float64(ngenes)))        // total num of individuals in grid
	den := 1.0                                                   // denominator to calculate dx
	if npts > 1 {
		den = float64(npts - 1)
	}
	var lfto int // leftover, e.g. n % (nx*ny)
	var rdim int // reduced dimension, e.g. (nx*ny)
	var idx int  // index of gene in grid
	var dx, x, mul, xmin, xmax float64
	for i := 0; i < C.Ninds; i++ {
		if i < ntot { // on grid
			lfto = i
			for j := 0; j < ngenes; j++ {
				rdim = int(math.Pow(float64(npts), float64(ngenes-1-j)))
				idx = lfto / rdim
				lfto = lfto % rdim
				xmin = C.RangeFlt[j][0]
				xmax = C.RangeFlt[j][1]
				dx = xmax - xmin
				x = xmin + float64(idx+id)*dx/den
				if C.Noise > 0 {
					mul = rnd.Float64(0, C.Noise)
					if rnd.FlipCoin(0.5) {
						x += mul * x
					} else {
						x -= mul * x
					}
				}
				if x < xmin {
					x = xmin + (xmin - x)
				}
				if x > xmax {
					x = xmax - (x - xmax)
				}
				if x < xmin {
					x = xmin
				}
				if x > xmax {
					x = xmax
				}
				o[i].Floats[j] = x
			}
		} else { // additional individuals
			for j := 0; j < ngenes; j++ {
				xmin = C.RangeFlt[j][0]
				xmax = C.RangeFlt[j][1]
				x = rnd.Float64(xmin, xmax)
				o[i].Floats[j] = x
			}
		}
	}
	return o
}

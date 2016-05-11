// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cec09

var (
	Nx   map[string]int
	Nf   map[string]int
	Ng   map[string]int
	Nh   map[string]int
	Xmin map[string][]float64
	Xmax map[string][]float64
)

func init() {
	Nx = map[string]int{
		"UF1":  30,
		"UF2":  30,
		"UF3":  30,
		"UF4":  30,
		"UF5":  30,
		"UF6":  30,
		"UF7":  30,
		"UF8":  30,
		"UF9":  30,
		"UF10": 30,
		"CF1":  10,
		"CF2":  10,
		"CF3":  10,
		"CF4":  10,
		"CF5":  10,
		"CF6":  10,
		"CF7":  10,
		"CF8":  10,
		"CF9":  10,
		"CF10": 10,
	}
	Nf = map[string]int{
		"UF1":  2,
		"UF2":  2,
		"UF3":  2,
		"UF4":  2,
		"UF5":  2,
		"UF6":  2,
		"UF7":  2,
		"UF8":  3,
		"UF9":  3,
		"UF10": 3,
		"CF1":  2,
		"CF2":  2,
		"CF3":  2,
		"CF4":  2,
		"CF5":  2,
		"CF6":  2,
		"CF7":  2,
		"CF8":  3,
		"CF9":  3,
		"CF10": 3,
	}
	Ng = map[string]int{
		"UF1":  0,
		"UF2":  0,
		"UF3":  0,
		"UF4":  0,
		"UF5":  0,
		"UF6":  0,
		"UF7":  0,
		"UF8":  0,
		"UF9":  0,
		"UF10": 0,
		"CF1":  1,
		"CF2":  1,
		"CF3":  1,
		"CF4":  1,
		"CF5":  1,
		"CF6":  2,
		"CF7":  2,
		"CF8":  1,
		"CF9":  1,
		"CF10": 1,
	}
	Nh = map[string]int{
		"UF1":  0,
		"UF2":  0,
		"UF3":  0,
		"UF4":  0,
		"UF5":  0,
		"UF6":  0,
		"UF7":  0,
		"UF8":  0,
		"UF9":  0,
		"UF10": 0,
		"CF1":  0,
		"CF2":  0,
		"CF3":  0,
		"CF4":  0,
		"CF5":  0,
		"CF6":  0,
		"CF7":  0,
		"CF8":  0,
		"CF9":  0,
		"CF10": 0,
	}
	Xmin = map[string][]float64{
		"UF1": []float64{0, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		"UF2": []float64{0, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		"UF3": []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	Xmax = map[string][]float64{
		"UF1": []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		"UF2": []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		"UF3": []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
}

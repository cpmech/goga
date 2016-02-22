// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "math"

const (
	SQ2 = math.Sqrt2
)

// limit state function
type LSF_T func(x []float64, cpu int) (lsf float64, failed float64)

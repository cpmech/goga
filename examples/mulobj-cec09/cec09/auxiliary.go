// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cec09

import (
	"github.com/cpmech/gosl/io"
)

// PFdata loads data
func PFdata(problem string) (dat [][]float64) {
	return io.ReadMatrix(io.Sf("pf_data/%s.dat", problem))
}

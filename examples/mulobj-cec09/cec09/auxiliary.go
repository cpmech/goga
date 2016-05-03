// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cec09

/*
#include "IGD.h"
*/
import "C"

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func PFdata(problem string) (dat [][]float64) {
	dat, err := io.ReadMatrix(io.Sf("pf_data/%s.dat", problem))
	if err != nil {
		chk.Panic("cannot load data for %q\n%v", problem, err)
	}
	return
}

func IGD() {
}

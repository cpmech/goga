// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cec09

import "testing"

func Test_dims(tst *testing.T) {
	if len(Xmin["UF1"]) != Nx["UF1"] {
		tst.Errorf("UF1: size of Xmin is incorrect\n")
		return
	}
}

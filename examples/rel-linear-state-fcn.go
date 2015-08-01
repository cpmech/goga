// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func g(x []float64) float64 {
	return x[0] + 2.0*x[1] + 2.0*x[2] + x[3] - 5.0*x[4] - 5.0*x[5]
}

func main() {

	// Example 1 from:
	// Cheng J and Xiao RC (2005) Serviceability reliability analysis of cable-stayed bridges.
	// Structural Engineering and Mechanics, Vol 20, No 6, 609-630

	// statistics
	mean := []float64{120, 120, 120, 120, 50, 40}
	devi := []float64{12, 12, 12, 12, 12, 12}
	dist := []string{"log", "log", "log", "log", "log", "log"}

	// FOSM

}

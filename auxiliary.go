// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

// GtPenalty implements a 'greater than' penalty function where
// x must be greater than b; otherwise the error is magnified
func GtPenalty(x, b, penaltyM float64) float64 {
	if x > b {
		return 0.0
	}
	return penaltyM*(b-x) + 1e-16 // must add small number because x must be greater than b
}

// GtePenalty implements a 'greater than or equal' penalty function where
// x must be greater than b or equal to be; otherwise the error is magnified
func GtePenalty(x, b, penaltyM float64) float64 {
	if x >= b {
		return 0.0
	}
	return penaltyM * (b - x)
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func printThickLine(n int) (l string) {
	for i := 0; i < n; i++ {
		l += "="
	}
	return l + "\n"
}

func printThinLine(n int) (l string) {
	for i := 0; i < n; i++ {
		l += "-"
	}
	return l + "\n"
}

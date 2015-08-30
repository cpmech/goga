// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/rnd"

// StrCrossover performs the crossover of genetic data from A and B
//  Output:
//   a and b -- offspring
//  Example:
//         0 1 2 3 4 5 6 7
//     A = a b c d e f g h    size = 8
//     B = * . . . . * * *    cuts = [1, 5]
//          ↑       ↑     ↑   ends = [1, 5, 8]
//          1       5     8
//     a = a . . . . f g h
//     b = * b c d e * * *
func StrCrossover(a, b, A, B []string, time int, ops *OpsData) (ends []int) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends = GenerateCxEnds(size, ops.Ncuts, ops.Cuts)
	swap := false
	start := 0
	for _, end := range ends {
		if swap {
			for j := start; j < end; j++ {
				b[j], a[j] = A[j], B[j]
			}
		} else {
			for j := start; j < end; j++ {
				a[j], b[j] = A[j], B[j]
			}
		}
		start = end
		swap = !swap
	}
	return
}

// StrMutation performs the mutation of genetic data from A
//  Output: modified individual 'A'
func StrMutation(A []string, time int, ops *OpsData) {
	size := len(A)
	if !rnd.FlipCoin(ops.Pm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, ops.Nchanges)
	for _, i := range pos {
		A[i] = "TODO" // TODO: improve this
	}
}

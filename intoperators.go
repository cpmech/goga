// Copyright 2015 The Goga Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/rnd"
)

// CxInt performs the crossover of genetic data from A and B
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
func CxInt(a, b, A, B []int, prms *Parameters) {
	size := len(A)
	if !rnd.FlipCoin(prms.IntPc) || size < 2 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	ends := GenerateCxEnds(size, prms.IntNcuts, nil)
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

// CxIntOrd performs the crossover in a pair of individuals with integer numbers
// that correspond to a ordered sequence, e.g. for traveling salesman problem
//  Output:
//    a and b -- offspring chromosomes
//  Note: using OX1 method explained in [1] (proposed in [2])
//  References:
//   [1] Larrañaga P, Kuijpers CMH, Murga RH, Inza I and Dizdarevic S. Genetic Algorithms for the
//       Travelling Salesman Problem: A Review of Representations and Operators. Artificial
//       Intelligence Review, 13:129-170; 1999. doi:10.1023/A:1006529012972
//   [2] Davis L. Applying Adaptive Algorithms to Epistatic Domains. Proceedings of International
//       Joint Conference on Artificial Intelligence, 162-164; 1985.
//  Example:
//   data:
//         0 1   2 3 4   5 6 7
//     A = a b | c d e | f g h        size = 8
//     B = b d | f h g | e c a        cuts = [2, 5]
//             ↑       ↑       ↑      ends = [2, 5, 8]
//             2       5       8
//   first step: copy subtours
//     a = . . | f h g | . . .
//     b = . . | c d e | . . .
//   second step: copy unique from subtour's end, position 5
//               start adding here
//                       ↓                           5 6 7   0 1   2 3 4
//     a = d e | f h g | a b c         get from A: | f̶ g̶ h̶ | a b | c d e
//     b = h g | c d e | a b f         get from B: | e̶ c̶ a | b d̶ | f h g
func CxIntOrd(a, b, A, B []int, prms *Parameters) {
	size := len(A)
	if !rnd.FlipCoin(prms.IntPc) || size < 3 {
		for i := 0; i < len(A); i++ {
			a[i], b[i] = A[i], B[i]
		}
		return
	}
	var s, t int
	var cuts []int
	if len(cuts) == 2 {
		s, t = cuts[0], cuts[1]
	} else {
		s = rnd.Int(1, size-2)
		t = rnd.Int(s+1, size-1)
	}
	chk.IntAssertLessThan(s, t)
	acore := B[s:t]
	bcore := A[s:t]
	ncore := t - s
	acorehas := make(map[int]bool) // TODO: check if map can be replaced => improve efficiency
	bcorehas := make(map[int]bool)
	for i := 0; i < ncore; i++ {
		a[s+i] = acore[i]
		b[s+i] = bcore[i]
		acorehas[acore[i]] = true
		bcorehas[bcore[i]] = true
	}
	ja, jb := t, t
	for i := 0; i < size; i++ {
		k := (i + t) % size
		if !acorehas[A[k]] {
			a[ja] = A[k]
			ja++
			if ja == size {
				ja = 0
			}
		}
		if !bcorehas[B[k]] {
			b[jb] = B[k]
			jb++
			if jb == size {
				jb = 0
			}
		}
	}
	return
}

// mutation ////////////////////////////////////////////////////////////////////////////////////////

// MtInt performs the mutation of genetic data from A
//  Output: modified individual 'A'
func MtInt(A []int, prms *Parameters) {
	size := len(A)
	if !rnd.FlipCoin(prms.IntPm) || size < 1 {
		return
	}
	mmax := 2
	pos := rnd.IntGetUniqueN(0, size, prms.IntNchanges)
	for _, i := range pos {
		m := rnd.Int(1, mmax)
		if rnd.FlipCoin(0.5) {
			A[i] += m * A[i]
		} else {
			A[i] -= m * A[i]
		}
	}
}

// MtIntBin performs the mutation of a binary chromosome
//  Output: modified individual 'A'
func MtIntBin(A []int, prms *Parameters) {
	size := len(A)
	if !rnd.FlipCoin(prms.IntPm) || size < 1 {
		return
	}
	pos := rnd.IntGetUniqueN(0, size, prms.IntNchanges)
	for _, i := range pos {
		if A[i] == 0 {
			A[i] = 1
		} else {
			A[i] = 0
		}
	}
}

// MtIntOrd performs the mutation of genetic data from a ordered list of integers A
//  Output: modified individual 'A'
//  Note: using DM method as explained in [1] (citing [2])
//  References:
//   [1] Larrañaga P, Kuijpers CMH, Murga RH, Inza I and Dizdarevic S. Genetic Algorithms for the
//       Travelling Salesman Problem: A Review of Representations and Operators. Artificial
//       Intelligence Review, 13:129-170; 1999. doi:10.1023/A:1006529012972
//   [2] Michalewicz Z. Genetic Algorithms + Data Structures = Evolution Programs. Berlin
//       Heidelberg: Springer Verlag; 1992
//       Joint Conference on Artificial Intelligence, 162-164; 1985.
//
//  DM displacement mutation method:
//   Ex:
//           0 1 2 3 4 5 6 7
//       A = a b c d e f g h   s = 2
//              ↑     ↑        t = 5
//              2     5
//
//       core = c d e  (subtour)  ncore = t - s = 5 - 2 = 3
//
//                0 1 2 3 4
//       remain = a b f g h  (remaining)  nrem = size - ncore = 8 - 3 = 5
//                       ↑
//                       4 = ins
func MtIntOrd(A []int, prms *Parameters) {
	size := len(A)
	if !rnd.FlipCoin(prms.IntPm) || size < 3 {
		if size == 2 {
			A[0], A[1] = A[1], A[0]
		}
		return
	}
	var indices []int
	var s, t, ncore, nrem, ins int
	if indices != nil {
		s, t, ins = indices[0], indices[1], indices[2]
		ncore = t - s
		nrem = size - ncore
	} else {
		s = rnd.Int(1, size-2)
		t = rnd.Int(s+1, size-1)
		ncore = t - s
		nrem = size - ncore
		ins = rnd.Int(1, nrem)
	}
	core := make([]int, ncore)
	remain := make([]int, nrem)
	var jc, jr int
	for i := 0; i < size; i++ {
		if i >= s && i < t {
			core[jc] = A[i]
			jc++
		} else {
			remain[jr] = A[i]
			jr++
		}
	}
	jc, jr = 0, 0
	for i := 0; i < size; i++ {
		if i < ins {
			A[i] = remain[jr]
			jr++
		} else {
			if jc < ncore {
				A[i] = core[jc]
				jc++
			} else {
				A[i] = remain[jr]
				jr++
			}
		}
	}
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// GenerateCxEnds randomly computes the end positions of cuts in chromosomes
//  Input:
//   size  -- size of chromosome
//   ncuts -- number of cuts to be used, unless cuts != nil
//   cuts  -- cut positions. can be nil => use ncuts instead
//  Output:
//   ends -- end positions where the last one equals size
//  Example:
//        0 1 2 3 4 5 6 7
//    A = a b c d e f g h    size = 8
//         ↑       ↑     ↑   cuts = [1, 5]
//         1       5     8   ends = [1, 5, 8]
func GenerateCxEnds(size, ncuts int, cuts []int) (ends []int) {

	// handle small slices
	if size < 2 {
		return
	}
	if size == 2 {
		return []int{1, size}
	}

	// cuts slice is given
	if len(cuts) > 0 {
		ncuts = len(cuts)
		ends = make([]int, ncuts+1)
		ends[ncuts] = size
		for i, cut := range cuts {
			if cut < 1 || cut >= size {
				chk.Panic("cut=%d is outside the allowed range: 1 ≤ cut ≤ size-1", cut)
			}
			if i > 0 {
				if cut == cuts[i-1] {
					chk.Panic("repeated cut values are not allowed: cuts=%v", cuts)
				}
			}
			ends[i] = cut
		}
		sort.Ints(ends)
		return
	}

	// randomly generate cuts
	if ncuts < 1 {
		ncuts = 1
	}
	if ncuts >= size {
		ncuts = size - 1
	}
	ends = make([]int, ncuts+1)
	ends[ncuts] = size

	// pool of values for selections
	pool := rnd.IntGetUniqueN(1, size, ncuts)
	sort.Ints(pool)
	for i := 0; i < ncuts; i++ {
		ends[i] = pool[i]
	}
	return
}

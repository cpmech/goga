// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cec09

/*
#cgo CXXFLAGS: -O3 -I.
#cgo CFLAGS: -O3 -I.
#cgo LDFLAGS: -lm -ldl
#include "cec09.h"
#define UINT unsigned int
*/
import "C"

import "unsafe"

func UF1(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF1((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF2(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF2((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF3(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF3((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF4(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF4((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF5(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF5((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF6(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF6((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF7(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF7((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF8(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF8((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF9(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF9((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func UF10(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.UF10((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (C.UINT)(nx))
}
func CF1(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF1((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF2(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF2((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF3(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF3((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF4(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF4((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF5(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF5((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF6(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF6((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF7(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF7((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF8(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF8((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF9(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF9((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}
func CF10(f, g, h, x []float64, ξ []int, cpu int) {
	nx := len(x)
	C.CF10((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (C.UINT)(nx))
}

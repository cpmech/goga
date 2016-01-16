// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/cpmech/gemlab"
	"github.com/cpmech/gosl/io"
)

func main() {

	L, H := 1000.0, 300.0
	ncol, nrow := 6, 2
	xmin, ymin := 0.0, 0.0

	nx, ny := ncol+1, nrow+1
	nhoriz := ny * ncol
	nverti := nx * nrow
	ndiago := ncol * nrow * 2
	ncells := nhoriz + nverti + ndiago
	npoints := nx * ny
	dx, dy := L/float64(ncol), H/float64(nrow)

	io.Pforan("nx=%v ny=%v nhoriz=%v nverti=%v ncells=%v npoints=%v\n", nx, ny, nhoriz, nverti, ncells, npoints)

	var dat gemlab.InData
	dat.AddCells = new(gemlab.AddCells)
	C := dat.AddCells
	C.Tags = make([]int, ncells)
	C.Points = make([][]float64, npoints)
	C.Types = make([]string, ncells)
	C.Conn = make([][]int, ncells)

	var ip, ic int
	for j := 0; j < ny; j++ {
		y := ymin + float64(j)*dy
		for i := 0; i < nx; i++ {
			x := xmin + float64(i)*dx
			C.Points[ip] = []float64{x, y}
			if i > 0 && j == 0 {
				C.Conn[ic] = []int{ip - 1, ip} // horizontal
				C.Types[ic] = "lin2"
				C.Tags[ic] = -1
				ic++
			}
			if i == 0 && j > 0 {
				C.Conn[ic] = []int{ip - nx, ip} // vertical
				C.Types[ic] = "lin2"
				C.Tags[ic] = -1
				ic++
			}
			if i > 0 && j > 0 {
				C.Conn[ic+0] = []int{ip - 1, ip}      // horizontal
				C.Conn[ic+1] = []int{ip - nx, ip}     // vertical
				C.Conn[ic+2] = []int{ip - nx - 1, ip} // diagonal
				C.Conn[ic+3] = []int{ip - 1, ip - nx} // diagonal
				C.Types[ic+0] = "lin2"
				C.Types[ic+1] = "lin2"
				C.Types[ic+2] = "lin2"
				C.Types[ic+3] = "lin2"
				C.Tags[ic+0] = -1
				C.Tags[ic+1] = -1
				C.Tags[ic+2] = -1
				C.Tags[ic+3] = -1
				ic += 4
			}
			ip++
		}
	}

	// -2 are places of applied load
	// -5 are the second row
	dat.VtagsL = &gemlab.VtagsL{
		Tags: []int{-2, -5},
		Xxa: [][]float64{
			{xmin, ymin},
			{xmin + dx, ymin + dy},
		},
		Xxb: [][]float64{
			{xmin + L, ymin},
			{xmin + L - dx, ymin + dy},
		},
	}

	// -1 is the vertex to track deflection
	// -3 and -4 are supports
	dat.Vtags = &gemlab.Vtags{
		Tags: []int{-1, -3, -4},
		Coords: [][]float64{
			{xmin + L/2, ymin},
			{xmin, ymin},
			{xmin + L, ymin},
		},
	}

	if err := gemlab.Generate(io.Sf("ground%d", ncells), &dat); err != nil {
		io.Pfred("gemlab failed:%v\n", err)
	}
}

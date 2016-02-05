// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gemlab"
	"github.com/cpmech/gosl/io"
)

func main() {

	L, H := 720.0, 360.0
	ncol, nrow := 4, 2
	xmin, ymin := 0.0, 0.0
	skipleft := true

	nx, ny := ncol+1, nrow+1
	nhoriz := ny * ncol
	nverti := nx * nrow
	ndiago := ncol * nrow * 2
	ncells := nhoriz + nverti + ndiago
	npoints := nx * ny
	dx, dy := L/float64(ncol), H/float64(nrow)
	if skipleft {
		ncells -= nrow
	}

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
			if i == 0 && j > 0 && !skipleft {
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

	// -10 are left vertical points
	// -20 are the lower horizontal bars
	// -30 are the second horizontal bars
	if true {
		dat.VtagsL = &gemlab.VtagsL{
			Tags: []int{-10, -20, -30},
			Xxa: [][]float64{
				{xmin, ymin},
				{xmin, ymin},
				{xmin + dx, ymin + dy},
			},
			Xxb: [][]float64{
				{xmin, ymin + H},
				{xmin + L, ymin},
				{xmin + L - dx, ymin + dy},
			},
		}
	}

	// -1 and -2 are supports
	// -3 are places to apply load
	// -4 is the vertex to track deflection
	dat.Vtags = &gemlab.Vtags{
		Tags: []int{-1, -2, -3, -4},
		Coords: [][]float64{
			{xmin, ymin},
			{xmin, ymin + H},
			{xmin + L/2, ymin},
			{xmin + L, ymin},
		},
	}

	if err := gemlab.Generate(io.Sf("ground%d", ncells), &dat); err != nil {
		io.Pfred("gemlab failed:%v\n", err)
	}
}

// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

// PlotOvs plots objective values versus time
func PlotOvs(isl *Island, ext, args string, t0, tf int, first, last bool) {
	if isl.C.DoPlot == false || isl.C.FnKey == "" {
		return
	}
	if first {
		plt.SetForEps(0.75, 250)
	}
	me := (tf-t0)/20 + isl.Id
	if me < 1 {
		me = 1
	}
	if len(args) > 0 {
		args += ","
	}
	nova := len(isl.Pop[0].Ovas)
	for i := 0; i < nova; i++ {
		plt.Plot(isl.OutTimes[t0:], isl.OutOvas[i][t0:], io.Sf("%s marker='%s', markersize=%d, markevery=%d, zorder=10, clip_on=0", args, get_marker(isl.Id), get_mrksz(isl.Id), me))
	}
	if last {
		plt.Gll("time", "objective value", "")
		plt.SaveD(isl.C.DirOut, isl.C.FnKey+"_ova"+ext)
	}
}

// PlotOor plots out-of-range values versus time
func PlotOor(isl *Island, ext, args string, t0, tf int, first, last bool) {
	if isl.C.DoPlot == false || isl.C.FnKey == "" {
		return
	}
	if first {
		plt.SetForEps(0.75, 250)
	}
	me := (tf-t0)/20 + isl.Id
	if me < 1 {
		me = 1
	}
	if len(args) > 0 {
		args += ","
	}
	noor := len(isl.Pop[0].Oors)
	for i := 0; i < noor; i++ {
		plt.Plot(isl.OutTimes[t0:], isl.OutOors[i][t0:], io.Sf("%s marker='%s', markersize=%d, markevery=%d, zorder=10, clip_on=0", args, get_marker(isl.Id), get_mrksz(isl.Id), me))
	}
	if last {
		plt.Gll("time", "out-of-range value", "")
		plt.SaveD(isl.C.DirOut, isl.C.FnKey+"_oor"+ext)
	}
}

// get_marker returns a marker for graphs
func get_marker(i int) string {
	pool := []string{"", "+", ".", "x", "s", "o", "*"}
	return pool[i%len(pool)]
}

// get_mrksz returns a marker size for graphs
func get_mrksz(i int) int {
	pool := []int{6, 6, 6, 3, 6, 6, 6}
	return pool[i%len(pool)]
}

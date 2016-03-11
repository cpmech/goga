// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

func main() {

	buf := goga.TexDocumentStart()

	io.Ff(buf, `
\begin{table} [!t] \centering
\caption{Random variables used in benchmark tests.}

\begin{tabular}[c]{ccccccc} \toprule
P & var & $\mu$ & $\sigma$ & distr & min & max \\ \hline
`)

	for i := 1; i < 19; i++ {
		opt := new(goga.Optimiser)
		opt.ProbNum = i
		_, vars := get_simple_data(opt)
		v := vars[0]
		io.Pforan("vars = %v\n", vars)
		io.Ff(buf, `%d & $x_0$ & %g & %g & %s & %g & %g \\`, i, v.M, v.S, rnd.GetDistrName(v.D), v.Min, v.Max)
		for j := 1; j < len(vars); j++ {
			v = vars[j]
			io.Ff(buf, ` & $x_%d$ & %g & %g & %s & %g & %g \\`, j, v.M, v.S, rnd.GetDistrName(v.D), v.Min, v.Max)
			io.Ff(buf, "\n")
		}
		io.Ff(buf, " \\hline\n\n")
	}

	io.Ff(buf, `
\end{tabular}
\label{tab:prms-simple}
\end{table}
`)

	goga.TexDocumentEnd(buf)
	io.WriteFileVD("/tmp/goga", "prms-table-simple.tex", buf)
}

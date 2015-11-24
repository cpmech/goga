// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import "github.com/cpmech/gosl/io"

func (o Solution) String() (l string) {
	l = io.Sf("%2d : %g : %g : %g : %d", o.Id, o.Ova, o.Oor, o.Flt, o.Int)
	return
}

func (o Group) String() (l string) {
	l = io.Sf("Group %d:\n", o.Id)
	for _, sol := range o.Solutions {
		l += io.Sf("%v\n", sol)
	}
	return
}

func (o Optimiser) String() (l string) {
	l = io.Sf("Parameters:\n%+v\n\nSolutions:\n", o.Parameters)
	for _, sol := range o.Solutions {
		l += io.Sf("%v\n", sol)
	}
	l += "\nGroups:\n"
	for _, grp := range o.Groups {
		l += io.Sf("%v\n", grp)
	}
	return l
}

func (o Optimiser) print_time(time, grp int) {
	if o.Verbose && grp == 0 {
		io.Pf(" ")
		if time%o.DtOut == 0 {
			io.Pfblue("%v", time)
			return
		}
		io.Pfgrey("%v", time)
	}
}

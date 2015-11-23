// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goga

import (
	"sort"

	"github.com/cpmech/gosl/chk"
)

// popByDemerit implements sort.Interface for Population based on Demerit
// Note: sorting population in increasing order of demerits: from best to worst
type popByDemerit []*Individual

// popById implements sort.Interface for Population based on Id. Note: ascending order
type popById []*Individual

// popByRank implements sort.Interface for Population based on Rank: FrontId followed by DistCrowd.
// Note: ascending order => best to worst
type popByRank []*Individual

// popByNwins implements sort.Interface for Population based on Nwins. Note: ascending order
type popByNwins []*Individual

// popByDistNeigh implements sort.Interface for Population based on DistNeigh
// Note: sorting population in decreasing order of DistNeigh: from diverse to less-diverse
type popByDistNeigh []*Individual

// popByOva0 implements sort.Interface for Population based on Ovas[0]. Note: ascending order
type popByOva0 []*Individual

// popByOva1 implements sort.Interface for Population based on Ovas[1]. Note: ascending order
type popByOva1 []*Individual

// popByOva2 implements sort.Interface for Population based on Ovas[2]. Note: ascending order
type popByOva2 []*Individual

// sorting functions: ByDemerit
func (o popByDemerit) Len() int           { return len(o) }
func (o popByDemerit) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o popByDemerit) Less(i, j int) bool { return o[i].Demerit < o[j].Demerit }

// sorting functions: ById
func (o popById) Len() int           { return len(o) }
func (o popById) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o popById) Less(i, j int) bool { return o[i].Id < o[j].Id }

// sorting functions: ByRank
func (o popByRank) Len() int      { return len(o) }
func (o popByRank) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o popByRank) Less(i, j int) bool {
	if o[i].FrontId == o[j].FrontId {
		return o[i].DistCrowd > o[j].DistCrowd
	}
	return o[i].FrontId < o[j].FrontId
}

// sorting functions: ByNwins
func (o popByNwins) Len() int           { return len(o) }
func (o popByNwins) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o popByNwins) Less(i, j int) bool { return o[i].Nwins < o[j].Nwins }

// sorting functions: ByDistNeigh
func (o popByDistNeigh) Len() int           { return len(o) }
func (o popByDistNeigh) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o popByDistNeigh) Less(i, j int) bool { return o[i].DistNeigh > o[j].DistNeigh }

// sorting functions: ByOva0
func (o popByOva0) Len() int           { return len(o) }
func (o popByOva0) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o popByOva0) Less(i, j int) bool { return o[i].Ovas[0] < o[j].Ovas[0] }

// sorting functions: ByOva1
func (o popByOva1) Len() int           { return len(o) }
func (o popByOva1) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o popByOva1) Less(i, j int) bool { return o[i].Ovas[1] < o[j].Ovas[1] }

// sorting functions: ByOva2
func (o popByOva2) Len() int           { return len(o) }
func (o popByOva2) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o popByOva2) Less(i, j int) bool { return o[i].Ovas[2] < o[j].Ovas[2] }

// SortByDemerit sorts population in incresing order of demerits: from best to worst
func (o Population) SortByDemerit() {
	sort.Sort(popByDemerit(o))
}

// SortById sorts population in incresing order of Id
func (o Population) SortById() {
	sort.Sort(popById(o))
}

// SortByRank sorts population in incresing order of Rank
func (o Population) SortByRank() {
	sort.Sort(popByRank(o))
}

// SortByNwins sorts population in incresing order of Nwins
func (o Population) SortByNwins() {
	sort.Sort(popByNwins(o))
}

// SortByDistNeigh sorts population in decresing order of DistNeigh: from diverse to less-diverse
func (o Population) SortByDistNeigh() {
	sort.Sort(popByDistNeigh(o))
}

// SortPopByOva sorts population in ascending order of ova
func SortPopByOva(pop Population, idxOva int) {
	switch idxOva {
	case 0:
		sort.Sort(popByOva0(pop))
	case 1:
		sort.Sort(popByOva1(pop))
	case 2:
		sort.Sort(popByOva2(pop))
	default:
		chk.Panic("this code can only handle Nova ≤ 3")
	}
}

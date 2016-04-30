// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore
package main

/*
#cgo LDFLAGS: -lm
#include "fcnsuite.h"
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
)

func getfcn(problem int) (opt *goga.Optimiser) {

	// GA parameters
	opt = new(goga.Optimiser)
	opt.Default()
	opt.EpsH = 0.0001

	// dims
	nx := []int{13, 20, 10, 5, 4, 2, 10, 2, 7, 8, 2, 3, 5, 10, 3, 5, 6, 9, 15, 24, 7, 22, 9, 2}
	nf := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	//          1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16  17  18 19 20 21 22 23 24
	ng := []int{9, 2, 0, 6, 2, 2, 8, 2, 4, 6, 0, 1, 0, 0, 0, 38, 0, 13, 5, 6, 1, 1, 2, 2}
	nh := []int{0, 0, 1, 0, 3, 0, 0, 0, 0, 0, 1, 0, 3, 3, 2, 0, 4, 0, 0, 14, 5, 19, 4, 0}

	// get fcn
	idx := problem - 1
	var fcn goga.MinProb_t // functions
	switch problem {
	case 1:
		opt.RptFref = []float64{-15}
		opt.RptXref = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 3, 3, 1}
		opt.FltMin = []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		opt.FltMax = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 100, 100, 100, 1}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g01((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 2:
		opt.RptFref = []float64{-0.80361910412559}
		opt.RptXref = []float64{3.16246061572185, 3.12833142812967, 3.09479212988791, 3.06145059523469, 3.02792915885555, 2.99382606701730, 2.95866871765285, 2.92184227312450, 0.49482511456933, 0.48835711005490, 0.48231642711865, 0.47664475092742, 0.47129550835493, 0.46623099264167, 0.46142004984199, 0.45683664767217, 0.45245876903267, 0.44826762241853, 0.44424700958760, 0.44038285956317}
		opt.FltMin = []float64{1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100}
		opt.FltMax = []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g02((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 3:
		opt.RptFref = []float64{-1.00050010001000}
		opt.RptXref = []float64{0.31624357647283069, 0.316243577414338339, 0.316243578012345927, 0.316243575664017895, 0.316243578205526066, 0.31624357738855069, 0.316243575472949512, 0.316243577164883938, 0.316243578155920302, 0.316243576147374916}
		opt.FltMin = []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		opt.FltMax = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g03((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), nil, (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 4:
		opt.RptFref = []float64{-3.066553867178332e+04}
		opt.RptXref = []float64{78, 33, 29.9952560256815985, 45, 36.7758129057882073}
		opt.FltMin = []float64{78, 33, 27, 27, 27}
		opt.FltMax = []float64{102, 45, 45, 45, 45}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g04((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 5:
		opt.RptFref = []float64{5126.4967140071}
		opt.RptXref = []float64{679.945148297028709, 1026.06697600004691, 0.118876369094410433, -0.39623348521517826}
		opt.FltMin = []float64{0, 0, -0.55, -0.55}
		opt.FltMax = []float64{1200, 1200, 0.55, 0.55}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g05((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 6:
		opt.RptFref = []float64{-6961.81387558015}
		opt.RptXref = []float64{14.0950000000000006400000000, 0.8429607892154795668000000}
		opt.FltMin = []float64{13, 0}
		opt.FltMax = []float64{100, 100}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g06((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 7:
		opt.RptFref = []float64{24.30620906818}
		opt.RptXref = []float64{2.17199634142692, 2.3636830416034, 8.77392573913157, 5.09598443745173, 0.990654756560493, 1.43057392853463, 1.32164415364306, 9.82872576524495, 8.2800915887356, 8.3759266477347}
		opt.FltMin = []float64{-10, -10, -10, -10, -10, -10, -10, -10, -10, -10}
		opt.FltMax = []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g07((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 8:
		opt.RptFref = []float64{-0.0958250414180359}
		opt.RptXref = []float64{1.22797135260752599, 4.24537336612274885}
		opt.FltMin = []float64{0, 0}
		opt.FltMax = []float64{10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g08((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 9:
		opt.RptFref = []float64{680.630057374402}
		opt.RptXref = []float64{2.33049935147405174, 1.95137236847114592, -0.477541399510615805, 4.36572624923625874, -0.624486959100388983, 1.03813099410962173, 1.5942266780671519}
		opt.FltMin = []float64{-10, -10, -10, -10, -10, -10, -10}
		opt.FltMax = []float64{10, 10, 10, 10, 10, 10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g09((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 10:
		opt.RptFref = []float64{}
		opt.RptXref = []float64{579.306685017979589, 1359.97067807935605, 5109.97065743133317, 182.01769963061534, 295.601173702746792, 217.982300369384632, 286.41652592786852, 395.601173702746735}
		opt.FltMin = []float64{100, 1000, 1000, 10, 10, 10, 10, 10}
		opt.FltMax = []float64{10000, 10000, 10000, 1000, 1000, 1000, 1000, 1000}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g10((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 11:
		opt.RptFref = []float64{0.7499}
		opt.RptXref = []float64{-0.707036070037170616, 0.500000004333606807}
		opt.FltMin = []float64{-1, -1}
		opt.FltMax = []float64{1, 1}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g11((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), nil, (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 12:
		opt.RptFref = []float64{-1}
		opt.RptXref = []float64{5, 5, 5}
		opt.FltMin = []float64{0, 0, 0}
		opt.FltMax = []float64{10, 10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g12((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 13:
		opt.RptFref = []float64{0.053941514041898}
		opt.RptXref = []float64{-1.71714224003, 1.59572124049468, 1.8272502406271, -0.763659881912867}
		opt.FltMin = []float64{-2.3, -2.3, -3.2, -3.2, -3.2}
		opt.FltMax = []float64{2.3, 2.3, 3.2, 3.2, 3.2}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g13((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), nil, (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 14:
		opt.RptFref = []float64{-47.7648884594915}
		opt.RptXref = []float64{0.0406684113216282, 0.147721240492452, 0.783205732104114, 0.00141433931889084, 0.485293636780388, 0.000693183051556082, 0.0274052040687766, 0.0179509660214818, 0.0373268186859717, 0.0968844604336845}
		opt.FltMin = []float64{1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100, 1e-100}
		opt.FltMax = []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g14((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), nil, (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 15:
		opt.RptFref = []float64{961.715022289961}
		opt.RptXref = []float64{3.51212812611795133, 0.216987510429556135, 3.55217854929179921}
		opt.FltMin = []float64{0, 0, 0}
		opt.FltMax = []float64{10, 10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g15((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), nil, (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 16:
		opt.RptFref = []float64{-1.90515525853479}
		opt.RptXref = []float64{705.174537070090537, 68.5999999999999943, 102.899999999999991, 282.324931593660324}
		opt.FltMin = []float64{704.4148, 68.6, 0, 193, 25}
		opt.FltMax = []float64{906.3855, 288.88, 134.75, 287.0966, 84.1988}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g16((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 17:
		opt.RptFref = []float64{8853.53967480648}
		opt.RptXref = []float64{201.784467214523659, 99.9999999999999005, 383.071034852773266, 420, -10.9076584514292652}
		opt.FltMin = []float64{0, 0, 340, 340, -1000, 0}
		opt.FltMax = []float64{400, 1000, 420, 420, 1000, 0.5236}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g17((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), nil, (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 18:
		opt.RptFref = []float64{-0.866025403784439}
		opt.RptXref = []float64{-0.657776192427943163, -0.153418773482438542, 0.323413871675240938, -0.946257611651304398}
		opt.FltMin = []float64{-10, -10, -10, -10, -10, -10, -10, -10, 0}
		opt.FltMax = []float64{10, 10, 10, 10, 10, 10, 10, 10, 20}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g18((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 19:
		opt.RptFref = []float64{32.6555929502463}
		opt.RptXref = []float64{1.66991341326291344e-17, 3.95378229282456509e-16, 3.94599045143233784, 1.06036597479721211e-16}
		opt.FltMin = []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		opt.FltMax = []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g19((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 20:
		opt.RptFref = []float64{0}
		opt.RptXref = []float64{1.28582343498528086e-18, 4.83460302526130664e-34, 0, 0, 6.30459929660781851e-18, 7.57192526201145068e-34, 5.03350698372840437e-34, 9.28268079616618064e-34, 0, 1.76723384525547359e-17, 3.55686101822965701e-34, 2.99413850083471346e-34, 0.158143376337580827, 2.29601774161699833e-19, 1.06106938611042947e-18, 1.31968344319506391e-18, 0.530902525044209539, 0, 2.89148310257773535e-18, 3.34892126180666159e-18, 0, 0.310999974151577319, 5.41244666317833561e-05, 4.84993165246959553e-16}
		opt.FltMin = []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		opt.FltMin[0] = 1e-100
		opt.FltMin[23] = 1e-100
		opt.FltMax = []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g20((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 21:
		opt.RptFref = []float64{193.724510070035}
		opt.RptXref = []float64{193.724510070034967, 5.56944131553368433e-27, 17.3191887294084914, 100.047897801386839}
		opt.FltMin = []float64{0, 0, 0, 100, 6.3, 5.9, 4.5}
		opt.FltMax = []float64{1000, 40, 40, 300, 6.7, 6.4, 6.25}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g21((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 22:
		opt.RptFref = []float64{236.430975504001}
		opt.RptXref = []float64{236.430975504001054, 135.82847151732463, 204.818152544824585, 6446.54654059436416, 3007540.83940215595, 4074188.65771341929, 32918270.5028952882, 130.075408394314167, 170.817294970528621, 299.924591605478554, 399.258113423595205, 330.817294971142758, 184.51831230897065, 248.64670239647424, 127.658546694545862, 269.182627528746707, 160.000016724090955, 5.29788288102680571, 5.13529735903945728, 5.59531526444068827, 5.43444479314453499, 5.07517453535834395}
		opt.FltMin = []float64{0, 0, 0, 0, 0, 0, 0, 100, 100, 100.01, 100, 100, 0, 0, 0, 0.01, 0.01, -4.7, -4.7, -4.7, -4.7, -4.7}
		opt.FltMax = []float64{20000, 1e+6, 1e+6, 1e+6, 4e+7, 4e+7, 4e+7, 299.99, 399.99, 300, 400, 600, 500, 500, 500, 300, 400, 6.25, 6.25, 6.25, 6.25, 6.25}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g22((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 23:
		opt.RptFref = []float64{-400.055099999999584}
		opt.RptXref = []float64{0.00510000000000259465, 99.9947000000000514, 9.01920162996045897e-18, 99.9999000000000535, 0.000100000000027086086, 2.75700683389584542e-14, 99.9999999999999574, 200, 0.0100000100000100008}
		opt.FltMin = []float64{0, 0, 0, 0, 0, 0, 0, 0, 0.01}
		opt.FltMax = []float64{300, 300, 100, 200, 100, 300, 100, 200, 0.03}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g23((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), (*C.double)(unsafe.Pointer(&h[0])), (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	case 24:
		opt.RptFref = []float64{-5.50801327159536}
		opt.RptXref = []float64{2.32952019747762, 3.17849307411774}
		opt.FltMin = []float64{0, 0}
		opt.FltMax = []float64{3, 4}
		fcn = func(f, g, h, x []float64, ξ []int, cpu int) {
			C.g24((*C.double)(unsafe.Pointer(&x[0])), (*C.double)(unsafe.Pointer(&f[0])), (*C.double)(unsafe.Pointer(&g[0])), nil, (C.int)(nx[idx]), (C.int)(nf[idx]), (C.int)(ng[idx]), (C.int)(nh[idx]))
			for i := 0; i < ng[idx]; i++ {
				g[i] = -g[i]
			}
		}
	default:
		chk.Panic("problem %d is not available", problem)
	}

	// number of trial solutions
	opt.Nsol = nx[idx] * 10

	// initialise optimiser
	opt.Init(goga.GenTrialSolutions, nil, fcn, nf[idx], ng[idx], nh[idx])
	return
}

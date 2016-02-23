// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/goga"
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/rnd"
)

// Problems 1 to 5 have all variables as standard variables => μ=0 and σ=1 => y = x
//
// References
//  SMB  Santos SR, Matioli LC and Beck AT. New optimization algorithms for structural
//       reliability analysis. Computer Modeling in Engineering & Sciences, 83(1):23-56; 2012
//       doi:10.3970/cmes.2012.083.023
//  BS   Borri A and Speranzini E. Structural reliability analysis using a standard
//       deterministic finite element code. Structural Safety, 19(4):361-382; 1997
//       doi:10.1016/S0167-4730(97)00017-9
//  GR   Grooteman F. Adaptive radial-based importance sampling method or structural
//       reliability. Structural safety, 30:533-542; 2008
//       doi:10.1016/j.strusafe.2007.10.002
//  GW   Grandhi RV and Wang L. Higher-order failure probability calculation using nonlinear
//       approximations. Computer Methods in Applied Mechanics and Engineering, 168:185-206; 1999
//  SSGK Santosh TV, Saraf RK, Ghosh AK and Kushwaha HS. Optimum step length selection rule in
//       modified HL-RF method for structural reliability. International Journal of Pressure
//       Vessels and Piping, 83(10):742-748; 2006
//       doi:10.1016/j.ijpvp.2006.07.004
//  HM   Haldar and Mahadevan. Probability, reliability and statistical methods in engineering
//       and design. John Wiley & Sons. 304p; 2000.
//  KLH  Der Kiureghian A, Lin H and Hwang S. Second‐Order Reliability Approximations.
//       Journal of Engineering Mechanics 113(8):1208-1225; 1987
//       doi:10.1061/(ASCE)0733-9399(1987)113:8(1208)
//  CX   Cheng J and Xiao RC. Serviceability reliability analysis of cable-stayed bridges.
//       Structural Engineering and Mechanics, 20(6):609-630; 2005
//       doi:10.12989/sem.2005.20.6.609
//  MS   Mahadevan S and Shi P. Multiple Linearization Method for Nonlinear Reliability
//       Analysiss, Journal of Engineering Mechanics, 127(11):1165-1173; 2001
//       doi:10.1061/(ASCE)0733-9399(2001)127:11(1165)

// Output:
//  lsf  -- limit state function
//  βref -- reference β (if available)
//  vars -- random variables data
func get_simple_data(opt *goga.Optimiser) (lsf LSF_T, vars rnd.Variables) {

	var desc string
	var βref float64
	var xref []float64
	switch opt.ProbNum {

	case 1:
		desc = "SMB1/BS5"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return 0.1*math.Pow(x[0]-x[1], 2) - (x[0]+x[1])/SQ2 + 2.5, 0.0
		}
		βref = 2.5 // from SMB
		xref = []float64{1.7677, 1.7677}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
		}

	case 2:
		desc = "SMB2/BS6"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return -0.5*math.Pow(x[0]-x[1], 2) - (x[0]+x[1])/SQ2 + 3, 0.0
		}
		βref = 1.658 // from BS6
		xref = []float64{-0.7583, 1.4752}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
		}

	case 3:
		desc = "SMB3/GR6"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return 2 - x[1] - 0.1*math.Pow(x[0], 2) + 0.06*math.Pow(x[0], 3), 0.0
		}
		βref = 2 // from SMB
		xref = []float64{0, 2}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
		}

	case 4:
		desc = "SMB4/GR8"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return 3 - x[1] + 256*math.Pow(x[0], 4), 0.0
		}
		βref = 3 // from SMB
		xref = []float64{0, 3}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
		}

	case 5:
		desc = "SMB5/GW1" // modified GW1
		shift := 0.0      // 0.1
		lsf = func(x []float64, cpu int) (float64, float64) {
			return 1 + math.Pow(x[0]+x[1]+shift, 2)/4 - 4*math.Pow(x[0]-x[1]+shift, 2), 0.0
		}
		βref = 0.3536 // from SMB
		xref = []float64{-βref * SQ2 / 2.0, βref * SQ2 / 2.0}
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
			&rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20},
		}

	case 6:
		desc = "SMB7/SSGK1a" // SSGK case 1
		lsf = func(x []float64, cpu int) (float64, float64) {
			return math.Pow(x[0], 3) + math.Pow(x[1], 3) - 18, 0.0
		}
		βref = 2.2401 // from SMB
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5, Min: -50, Max: 50},
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5, Min: -50, Max: 50},
		}

	case 7:
		desc = "SMB8/SSGK1b" // SSGK case 2
		lsf = func(x []float64, cpu int) (float64, float64) {
			return math.Pow(x[0], 3) + math.Pow(x[1], 3) - 18, 0.0
		}
		βref = 2.2260 // from SMB
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5, Min: -50, Max: 50},
			&rnd.VarData{D: rnd.D_Normal, M: 9.9, S: 5, Min: -50, Max: 50},
		}

	case 8:
		desc = "SMB9/GR7"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return 2.5 - 0.2357*(x[0]-x[1]) + 0.0046*math.Pow(x[0]+x[1]-20, 4), 0.0
		}
		βref = 2.5 // from SMB
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 3, Min: -50, Max: 50},
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 3, Min: -50, Max: 50},
		}

	case 9:
		desc = "SMB10/SSGK2"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return math.Pow(x[0], 3) + math.Pow(x[1], 3) - 67.5, 0.0
		}
		βref = 1.9003 // from SMB
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5, Min: -50, Max: 50},
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 5, Min: -50, Max: 50},
		}

	case 10:
		desc = "SMB11/GR2"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return x[0]*x[1] - 146.14, 0.0
		}
		βref = 5.4280 // from SMB
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 78064.4, S: 11709.7, Min: 1000, Max: 150000},
			&rnd.VarData{D: rnd.D_Normal, M: 0.0104, S: 0.00156, Min: -0.05, Max: 0.05},
		}

	case 11:
		desc = "SMB12/GW2"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return 2.2257 - 0.025*SQ2*math.Pow(x[0]+x[1]-20, 3)/27 + 0.2357*(x[0]-x[1]), 0.0
		}
		βref = 2.2257 // from SMB
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 3, Min: -50, Max: 50},
			&rnd.VarData{D: rnd.D_Normal, M: 10, S: 3, Min: -50, Max: 50},
		}

	case 12:
		desc = "SMB14/HM7.6"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return x[0]*x[1] - 1140, 0.0
		}
		βref = 5.2127 // from SMB // from here: 5.210977819456551
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Lognormal, M: 38, S: 3.8, Min: 20, Max: 60},
			&rnd.VarData{D: rnd.D_Lognormal, M: 54, S: 2.7, Min: 40, Max: 70},
		}

	// more than 2 variables -----------------------------------------------------------------

	case 13:
		desc = "SMB6/GR3"
		lsf = func(x []float64, cpu int) (float64, float64) {
			sum := 0.0
			for i := 0; i < 9; i++ {
				sum += x[i] * x[i]
			}
			return 2.0 - 0.015*sum - x[9], 0.0
		}
		βref = 2.0 // from SMB
		vars = make([]*rnd.VarData, 10)
		for i := 0; i < 10; i++ {
			vars[i] = &rnd.VarData{D: rnd.D_Normal, M: 0, S: 1, Min: -20, Max: 20}
		}
		opt.Nsol = 120
		opt.Ncpu = 4

	case 14:
		desc = "KLH1/CX1"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return x[0] + 2.0*x[1] + 2.0*x[2] + x[3] - 5.0*x[4] - 5.0*x[5], 0.0
		}
		βref = 2.348 // from CX1
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Lognormal, M: 120, S: 12, Min: 50, Max: 200},
			&rnd.VarData{D: rnd.D_Lognormal, M: 120, S: 12, Min: 50, Max: 200},
			&rnd.VarData{D: rnd.D_Lognormal, M: 120, S: 12, Min: 50, Max: 200},
			&rnd.VarData{D: rnd.D_Lognormal, M: 120, S: 12, Min: 50, Max: 200},
			&rnd.VarData{D: rnd.D_Lognormal, M: 50, S: 15, Min: 5, Max: 150},
			&rnd.VarData{D: rnd.D_Lognormal, M: 40, S: 12, Min: 5, Max: 120},
		}
		opt.Nsol = 60
		opt.Ncpu = 2

	case 15:
		desc = "SMB16/KLH2"
		lsf = func(x []float64, cpu int) (float64, float64) {
			s := x[0] + 2.0*x[1] + 2.0*x[2] + x[3] - 5.0*x[4] - 5.0*x[5]
			for i := 0; i < 6; i++ {
				s += 0.001 * math.Sin(1000*x[i])
			}
			return s, 0.0
		}
		βref = 2.3482 // from SMB
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Lognormal, M: 120, S: 12, Min: 50, Max: 200},
			&rnd.VarData{D: rnd.D_Lognormal, M: 120, S: 12, Min: 50, Max: 200},
			&rnd.VarData{D: rnd.D_Lognormal, M: 120, S: 12, Min: 50, Max: 200},
			&rnd.VarData{D: rnd.D_Lognormal, M: 120, S: 12, Min: 50, Max: 200},
			&rnd.VarData{D: rnd.D_Lognormal, M: 50, S: 15, Min: 5, Max: 150},
			&rnd.VarData{D: rnd.D_Lognormal, M: 40, S: 12, Min: 5, Max: 120},
		}
		opt.Nsol = 60
		opt.Ncpu = 2

	case 16:
		desc = "SMB17/MS5"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return -240758.1777 + 10467.364*x[0] + 11410.63*x[1] +
				3505.3015*x[2] - 246.81*x[0]*x[0] - 285.3275*x[1]*x[1] - 195.46*x[2]*x[2], 0.0
		}
		βref = 0.8292 // from SMB
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Lognormal, M: 21.2, S: 0.1, Min: 20, Max: 22},
			&rnd.VarData{D: rnd.D_Lognormal, M: 20.0, S: 0.2, Min: 19, Max: 21},
			&rnd.VarData{D: rnd.D_Lognormal, M: 9.2, S: 0.1, Min: 8, Max: 10},
		}
		opt.Nsol = 30

	case 17:
		desc = "SMB18/SSGK4a" // SSGK case 1
		lsf = func(x []float64, cpu int) (float64, float64) {
			return x[0]*x[1] - 78.12*x[2], 0.0
		}
		βref = 3.3221 // from SMB or 3.31819 from SSGK
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Normal, M: 2e7, S: 5e6, Min: 1e6, Max: 4e7},
			&rnd.VarData{D: rnd.D_Normal, M: 1e-4, S: 2e-5, Min: 1e-5, Max: 2e-4},
			&rnd.VarData{D: rnd.D_Gumbel, M: 4, S: 1.0, Min: 1, Max: 15},
		}
		opt.Nsol = 30

	case 18:
		desc = "SMB19/SSGK4b" // SSGK case 2
		lsf = func(x []float64, cpu int) (float64, float64) {
			return x[0]*x[1] - 78.12*x[2], 0.0
		}
		βref = 4.45272 // from SSGK
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Lognormal, M: 2e7, S: 5e6, Min: 1e6, Max: 5e7},
			&rnd.VarData{D: rnd.D_Lognormal, M: 1e-4, S: 2e-5, Min: 1e-5, Max: 3e-4},
			&rnd.VarData{D: rnd.D_Gumbel, M: 4, S: 1.0, Min: 1, Max: 15},
		}
		opt.Nsol = 30

	case 19:
		desc = "SMB20/SSGK5"
		lsf = func(x []float64, cpu int) (float64, float64) {
			return 1.1 - 0.00115*x[0]*x[1] + 0.00157*x[1]*x[1] + 0.00117*x[0]*x[0] +
				+0.0135*x[1]*x[2] - 0.0705*x[1] - 0.00534*x[0] - 0.0149*x[0]*x[2] +
				-0.0611*x[1]*x[3] + 0.0717*x[0]*x[3] - 0.226*x[2] + 0.0333*x[2]*x[2] +
				-0.558*x[2]*x[3] + 0.998*x[3] - 1.339*x[3]*x[3], 0.0
		}
		βref = 2.42031 // from SSGK
		vars = rnd.Variables{
			&rnd.VarData{D: rnd.D_Frechet, L: 8.782275, A: 4.095645, Min: 8, Max: 12},
			&rnd.VarData{D: rnd.D_Normal, M: 25, S: 5, Min: 5, Max: 50},
			&rnd.VarData{D: rnd.D_Normal, M: 0.8, S: 0.2, Min: 0.1, Max: 2.0},
			&rnd.VarData{D: rnd.D_Lognormal, M: 0.0625, S: 0.0625, Min: 0.001, Max: 0.4},
		}
		opt.Nsol = 40

	default:
		chk.Panic("simple problem number %d is invalid", opt.ProbNum)
	}
	opt.RptName = desc
	opt.RptName = io.Sf("%d", opt.ProbNum)
	opt.RptFref = []float64{βref}
	opt.RptXref = xref
	return
}

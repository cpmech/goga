 /*
 * =============================================================
 * IGD.cpp
 *
 * Matlab-C source codes
 *
 * IGD performance metric for CEC 2009 MOO Competition
 *
 * Usage: igd = IGD(PF*, PF, C) or igd = IGD(PF*, PF)
 *
 * Calculate the distance from the ideal Pareto Front (PF*) to an
 * obtained nondominated front (PF), C is the constraints
 *
 * PF*, PF, C MUST be columnwise, 
 *
 * Please refer to the report for more information.
 * =============================================================
 */

#ifndef IGD_H
#define IGD_H

#ifdef __cplusplus
extern "C" {
#endif

double IGD( double* S, double* Q, double *C, int row, int Scol, int Qcol, int Crow );

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif //IGD_H

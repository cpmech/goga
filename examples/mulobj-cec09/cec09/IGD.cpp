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

#include "math.h"

double IGD( double* S, double* Q, double *C, int row, int Scol, int Qcol, int Crow )
{
	int i, j, k, *feasible, nof;
	double d,min,dis;
	
	/* Step 1: remove the infeasible points, i.e, constraint<-1.0E-6 */
	feasible = new int[Qcol];
	for(i=0; i<Qcol; i++) feasible[i] = 1;
	nof = Qcol;
	if(C>0)
	{
		for(i=0; i<Qcol; i++)
		{
			for(k=0; k<Crow; k++)
			{ 
				if(C[i*Crow+k] < -1E-6)
				{
					feasible[i] = 0;
					nof --;
					break;
				}
			}
		}
	}
	if(nof==0) {
        delete feasible;
        return 1.0E6;
    }
	
	/* Step 2: calculate the IGD value for feasible points */
	dis=0.0;
	for( i=0; i<Scol; i++ )
	{
		min = 1.0E200;
		for( j=0; j<Qcol; j++ ) if( feasible[j]>0 )
		{
			d=0.0;
			for( k=0; k<row; k++ )
				d += ( S[i*row+k] - Q[j*row+k] ) * ( S[i*row+k] - Q[j*row+k] );
			if( d < min ) min = d;
		}
		dis += sqrt(min);
	}
    delete feasible;
	return dis/(double)(Scol);
}


/* The gateway routine */
// void mexFunction(int nlhs, mxArray *plhs[],
//                  int nrhs, const mxArray *prhs[])
// {
//   double *S, *Q, *C, dis;
//   int status, Srows, Scols, Qrows, Qcols, Crows, Ccols;
//   
//   S = Q = C = 0;
//   Srows = Scols = Qrows = Qcols = Crows = Ccols = 0;
// 
//   /*  Check for proper number of arguments. */
//   if (nrhs < 2 || nrhs >3 || nlhs != 1) mexErrMsgTxt("Usage: igd = IGD(PF*,PF,C)");
// 
//   /* Create pointers to the input matrix S, Q, and C */
//   S 	= mxGetPr(prhs[0]);  
//   Srows = mxGetM(prhs[0]);
//   Scols = mxGetN(prhs[0]);
//   Q 	= mxGetPr(prhs[1]);
//   Qrows = mxGetM(prhs[1]);
//   Qcols = mxGetN(prhs[1]);
//   if(nrhs == 3)
//   {
// 	  C 	= mxGetPr(prhs[2]);
// 	  Crows = mxGetM(prhs[2]);
// 	  Ccols = mxGetN(prhs[2]);  	  
//   }
//   
//   /* Check for dimension of input matrix. */
//   if (Srows != Qrows) mexErrMsgTxt("Input should be columnwise and dimension of PF* and PF should be the same.");
//   if (nrhs == 3 && Qcols != Ccols) mexErrMsgTxt("Input should be columnwise and dimension of PF and C should be the number.");
//   
//   /* Calculate the upsilon measure */
//   plhs[0] = mxCreateDoubleMatrix(1,1, mxREAL);
//   *(mxGetPr(plhs[0])) = IGD( S, Q, C, Srows, Scols, Qcols, Crows);
// }

/*
  Test function suite (last update:
  Thomas Philip Runarsson (email: tpr@hi.is) 
  Time-stamp: "2005-11-19 18:01:18 tpr"

  for linux:
  
  gcc -c fcnsuite.c
  ld -o fcnsuite.so -shared fcnsuite.o

  for windows:

  gcc -c fcnsuite.c -DWINDOWS
  dllwrap -o fcnsuite.dll fcnsuite.o
*/

#ifdef WINDOWS
#define DLLIMPORT __declspec (dllexport)
#else
#define DLLIMPORT
#endif

#include <math.h>

DLLIMPORT void
g01 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  int j;

  /* objective function */
  f[0] = 5.0 * (x[0] + x[1] + x[2] + x[3]) - 5.0 * (x[0] * x[0] + x[1] * x[1] + x[2] * x[2] + x[3] * x[3]);
  for (j = 4; j < 13; j++)
    f[0] = f[0] - x[j];
  /* constraints g<=0 */
  g[0] = 2.0 * x[0] + 2.0 * x[1] + x[9] + x[10] - 10.;
  g[1] = 2.0 * x[0] + 2.0 * x[2] + x[9] + x[11] - 10.;
  g[2] = 2.0 * x[1] + 2.0 * x[2] + x[10] + x[11] - 10.;
  g[3] = -8.0 * x[0] + x[9];
  g[4] = -8.0 * x[1] + x[10];
  g[5] = -8.0 * x[2] + x[11];
  g[6] = -2.0 * x[3] - x[4] + x[9];
  g[7] = -2.0 * x[5] - x[6] + x[10];
  g[8] = -2.0 * x[7] - x[8] + x[11];
}

DLLIMPORT void
g02 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  int j;
  double f1, f2, f3, g1, g2;

  /* objective function */
  f1 = 0.;
  f2 = 1.;
  f3 = 0.;
  g1 = 1.;
  g2 = 0.;
  for (j = 0; j < nx; j++)
    {
      f1 = f1 + pow (cos (x[j]), 4);
      f2 = f2 * cos (x[j]) * cos (x[j]);
      f3 = f3 + ((double) (j + 1)) * x[j] * x[j];
      g1 = g1 * x[j];
      g2 = g2 + x[j];
    }
  f[0] = fabs ((f1 - 2 * f2) / sqrt (f3));
  f[0] = -f[0]; /* Max-->Min, Modified by Jane,Nov 22 2005 */
  /* constraints g<=0 */
  g[0] = 0.75 - g1;
  g[1] = g2 - 7.5 * ((double) nx);
}

DLLIMPORT void
g03 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  int j;
  double f1, f2, f3 = sqrt ((double) nx);

  /* objective function */
  f1 = 1.;
  f2 = 0.;
  for (j = 0; j < nx; j++)
    {
      f1 = f3 * f1 * x[j];
      f2 = f2 + x[j] * x[j];
    }
  f[0] = f1;
  f[0] = -f[0]; /* Max-->Min, Modified by Jane,Nov 22 2005 */
  /* constraints g<=0 */
  h[0] = f2 - 1.0;
}

DLLIMPORT void
g04 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = 5.3578547 * x[2] * x[2] + 0.8356891 * x[0] * x[4] + 37.293239 * x[0] - 40792.141;
  /* constraints g<=0 */
  g[0] = 85.334407 + 0.0056858 * x[1] * x[4] + 0.0006262 * x[0] * x[3] - 0.0022053 * x[2] * x[4] - 92.;
  g[1] = -85.334407 - 0.0056858 * x[1] * x[4] - 0.0006262 * x[0] * x[3] + 0.0022053 * x[2] * x[4];
  g[2] = 80.51249 + 0.0071317 * x[1] * x[4] + 0.0029955 * x[0] * x[1] + 0.0021813 * x[2] * x[2] - 110.;
  g[3] = -80.51249 - 0.0071317 * x[1] * x[4] - 0.0029955 * x[0] * x[1] - 0.0021813 * x[2] * x[2] + 90.;
  g[4] = 9.300961 + 0.0047026 * x[2] * x[4] + 0.0012547 * x[0] * x[2] + 0.0019085 * x[2] * x[3] - 25.;
  g[5] = -9.300961 - 0.0047026 * x[2] * x[4] - 0.0012547 * x[0] * x[2] - 0.0019085 * x[2] * x[3] + 20.;
}

DLLIMPORT void
g05 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = 3.0 * x[0] + 0.000001 * pow (x[0], 3) + 2.0 * x[1] + (0.000002 / 3.0) * pow (x[1], 3);
  /* constraints h=0 => g=|h|-delta<=0 */
  g[0] = -x[3] + x[2] - 0.55;
  g[1] = -x[2] + x[3] - 0.55;
  h[0] = 1000.0 * sin (-x[2] - 0.25) + 1000.0 * sin (-x[3] - 0.25) + 894.8 - x[0];
  h[1] = 1000.0 * sin (x[2] - 0.25) + 1000.0 * sin (x[2] - x[3] - 0.25) + 894.8 - x[1];
  h[2] = 1000.0 * sin (x[3] - 0.25) + 1000.0 * sin (x[3] - x[2] - 0.25) + 1294.8;
}

DLLIMPORT void
g06 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = pow ((x[0] - 10.), 3) + pow ((x[1] - 20.), 3);
  /* constraints */
  g[0] = 100. - (x[0] - 5.) * (x[0] - 5.) - (x[1] - 5.) * (x[1] - 5.);
  g[1] = (x[0] - 6.) * (x[0] - 6.) + (x[1] - 5.) * (x[1] - 5.) - 82.81;
}

DLLIMPORT void
g07 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] =
    x[0] * x[0] + x[1] * x[1] + x[0] * x[1] - 14.0 * x[0] - 16.0 * x[1] + (x[2] - 10.0) * (x[2] - 10.0) + 4.0 * (x[3] -
														 5.0) *
    (x[3] - 5.0) + (x[4] - 3.0) * (x[4] - 3.0) + 2.0 * (x[5] - 1.0) * (x[5] - 1.0) + 5.0 * x[6] * x[6] + 7.0 * (x[7] -
														11) *
    (x[7] - 11) + 2.0 * (x[8] - 10.0) * (x[8] - 10.0) + (x[9] - 7.0) * (x[9] - 7.0) + 45.;
  /* constraints g<=0 */
  g[0] = -105.0 + 4.0 * x[0] + 5.0 * x[1] - 3.0 * x[6] + 9.0 * x[7];
  g[1] = 10.0 * x[0] - 8.0 * x[1] - 17.0 * x[6] + 2.0 * x[7];
  g[2] = -8.0 * x[0] + 2.0 * x[1] + 5.0 * x[8] - 2.0 * x[9] - 12.0;
  g[3] = 3.0 * (x[0] - 2.0) * (x[0] - 2.0) + 4.0 * (x[1] - 3.0) * (x[1] - 3.0) + 2.0 * x[2] * x[2] - 7.0 * x[3] - 120.0;
  g[4] = 5.0 * x[0] * x[0] + 8.0 * x[1] + (x[2] - 6.0) * (x[2] - 6.0) - 2.0 * x[3] - 40.0;
  g[5] = x[0] * x[0] + 2.0 * (x[1] - 2.0) * (x[1] - 2.0) - 2.0 * x[0] * x[1] + 14.0 * x[4] - 6.0 * x[5];
  g[6] = 0.5 * (x[0] - 8.0) * (x[0] - 8.0) + 2.0 * (x[1] - 4.0) * (x[1] - 4.0) + 3.0 * x[4] * x[4] - x[5] - 30.0;
  g[7] = -3.0 * x[0] + 6.0 * x[1] + 12.0 * (x[8] - 8.0) * (x[8] - 8.0) - 7.0 * x[9];
}

DLLIMPORT void
g08 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{

  double pi = 4.0 * atan (1.0);

/* objective function */
  f[0] = pow (sin (2 * pi * x[0]), 3) * sin (2 * pi * x[1]) / (pow (x[0], 3) * (x[0] + x[1]));
  f[0] = -f[0]; /* Max-->Min, Modified by Jane,Nov 22 2005 */
/* constraints */
  g[0] = x[0] * x[0] - x[1] + 1.0;
  g[1] = 1.0 - x[0] + (x[1] - 4.0) * (x[1] - 4.0);
}

DLLIMPORT void
g09 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] =
    (x[0] - 10.0) * (x[0] - 10.0) + 5.0 * (x[1] - 12.0) * (x[1] - 12.0) + pow (x[2],
									       4) + 3.0 * (x[3] - 11.0) * (x[3] -
													   11.0) +
    10.0 * pow (x[4], 6) + 7.0 * x[5] * x[5] + pow (x[6], 4) - 4.0 * x[5] * x[6] - 10.0 * x[5] - 8.0 * x[6];
  /* constraints g<=0 */
  g[0] = -127.0 + 2 * x[0] * x[0] + 3.0 * pow (x[1], 4) + x[2] + 4.0 * x[3] * x[3] + 5.0 * x[4];
  g[1] = -282.0 + 7.0 * x[0] + 3.0 * x[1] + 10.0 * x[2] * x[2] + x[3] - x[4];
  g[2] = -196.0 + 23.0 * x[0] + x[1] * x[1] + 6.0 * x[5] * x[5] - 8.0 * x[6];
  g[3] = 4.0 * x[0] * x[0] + x[1] * x[1] - 3.0 * x[0] * x[1] + 2.0 * x[2] * x[2] + 5.0 * x[5] - 11.0 * x[6];
}

DLLIMPORT void
g10 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = x[0] + x[1] + x[2];
  /* constraints g<=0 */
  g[0] = -1.0 + 0.0025 * (x[3] + x[5]);
  g[1] = -1.0 + 0.0025 * (x[4] + x[6] - x[3]);
  g[2] = -1.0 + 0.01 * (x[7] - x[4]);
  g[3] = -x[0] * x[5] + 833.33252 * x[3] + 100.0 * x[0] - 83333.333;
  g[4] = -x[1] * x[6] + 1250.0 * x[4] + x[1] * x[3] - 1250.0 * x[3];
  g[5] = -x[2] * x[7] + 1250000.0 + x[2] * x[4] - 2500.0 * x[4];
}

DLLIMPORT void
g11 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = x[0] * x[0] + (x[1] - 1.0) * (x[1] - 1.0);
  /* constraints g<=0 */
  h[0] = x[1] - x[0] * x[0];
}

DLLIMPORT void
g12 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  double gt;
  int i, j, k;

/* objective function */
  f[0] = (100. - (x[0] - 5.) * (x[0] - 5.) - (x[1] - 5.) * (x[1] - 5.) - (x[2] - 5.) * (x[2] - 5.)) / 100.;
  f[0] = -f[0]; /* Max-->Min, Modified by Jane,Nov 22 2005 */
/* constraints */
  g[0] = (x[0] - 1.) * (x[0] - 1.) + (x[1] - 1.) * (x[1] - 1.) + (x[2] - 1.) * (x[2] - 1.) - 0.0625;
  for (i = 1; i <= 9; i++)
    {
      for (j = 1; j <= 9; j++)
	{
	  for (k = 1; k <= 9; k++)
	    {
	      gt = (x[0] - i) * (x[0] - i) + (x[1] - j) * (x[1] - j) + (x[2] - k) * (x[2] - k) - 0.0625;
	      if (gt < g[0])
		g[0] = gt;
	    }
	}
    }
}

DLLIMPORT void
g13 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = exp (x[0] * x[1] * x[2] * x[3] * x[4]);
  /* constraints h(x) = 0 */
  h[0] = x[0] * x[0] + x[1] * x[1] + x[2] * x[2] + x[3] * x[3] + x[4] * x[4] - 10.0;
  h[1] = x[1] * x[2] - 5.0 * x[3] * x[4];
  h[2] = pow (x[0], 3) + pow (x[1], 3) + 1.0;
}

DLLIMPORT void
g14 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  int i;
  double sumlog = 0.0, sum = 0.0;
  double C[10] = { -6.089, -17.164, -34.054, -5.914, -24.721, -14.986, -24.100, -10.708, -26.662, -22.179 };

  /* objective function */
  for (i = 0; i < 10; i++)
    sumlog += x[i];
  for (i = 0; i < 10; i++)
    sum += x[i] * (C[i] + log (x[i] / sumlog));
  f[0] = sum;
  /* constraints h=0 */
  h[0] = x[0] + 2.0 * x[1] + 2.0 * x[2] + x[5] + x[9] - 2.0;
  h[1] = x[3] + 2.0 * x[4] + x[5] + x[6] - 1.0;
  h[2] = x[2] + x[6] + x[7] + 2.0 * x[8] + x[9] - 1.0;
}

DLLIMPORT void
g15 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = 1000.0 - pow (x[0], 2.0) - 2.0 * x[1] * x[1] - x[2] * x[2] - x[0] * x[1] - x[0] * x[2];
  /* constraints h=0 */
  h[0] = pow (x[0], 2.0) + pow (x[1], 2.0) + pow (x[2], 2.0) - 25.0;
  h[1] = 8.0 * x[0] + 14.0 * x[1] + 7.0 * x[2] - 56.0; /* BUG FIXED was 8.0 * x[1] */
}

DLLIMPORT void
g16 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  double x1, x2, x3, x4, x5;
  double C[17], Y[17];

  x1 = x[0];
  x2 = x[1];
  x3 = x[2];
  x4 = x[3];
  x5 = x[4];

  Y[0] = x2 + x3 + 41.6;
  C[0] = 0.024 * x4 - 4.62;
  Y[1] = (12.5 / C[0]) + 12.0;
  C[1] = 0.0003535 * pow (x1, 2.0) + 0.5311 * x1 + 0.08705 * Y[1] * x1;
  C[2] = 0.052 * x1 + 78.0 + 0.002377 * Y[1] * x1;
  Y[2] = C[1] / C[2];
  Y[3] = 19.0 * Y[2];
  C[3] = 0.04782 * (x1 - Y[2]) + ((0.1956 * pow (x1 - Y[2], 2.0)) / x2) + 0.6376 * Y[3] + 1.594 * Y[2];
  C[4] = 100 * x2;
  C[5] = x1 - Y[2] - Y[3];
  C[6] = 0.950 - (C[3] / C[4]);
  Y[4] = C[5] * C[6];
  Y[5] = x1 - Y[4] - Y[3] - Y[2];
  C[7] = (Y[4] + Y[3]) * 0.995;
  Y[6] = C[7] / Y[0];
  Y[7] = C[7] / 3798.0;
  C[8] = Y[6] - (0.0663 * Y[6] / Y[7]) - 0.3153;
  Y[8] = (96.82 / C[8]) + 0.321 * Y[0];
  Y[9] = 1.29 * Y[4] + 1.258 * Y[3] + 2.29 * Y[2] + 1.71 * Y[5];
  Y[10] = 1.71 * x1 - 0.452 * Y[3] + 0.580 * Y[2];
  C[9] = 12.3 / 752.3;
  C[10] = 1.75 * Y[1] * 0.995 * x1;
  C[11] = 0.995 * Y[9] + 1998.0;
  Y[11] = C[9] * x1 + (C[10] / C[11]);
  Y[12] = C[11] - 1.75 * Y[1];
  Y[13] = 3623.0 + 64.4 * x2 + 58.4 * x3 + (146312.0 / (Y[8] + x5));
  C[12] = 0.995 * Y[9] + 60.8 * x2 + 48 * x4 - 0.1121 * Y[13] - 5095.0;
  Y[14] = Y[12] / C[12];
  Y[15] = 148000.0 - 331000.0 * Y[14] + 40.0 * Y[12] - 61.0 * Y[14] * Y[12];
  C[13] = 2324 * Y[9] - 28740000 * Y[1];
  Y[16] = 14130000 - 1328.0 * Y[9] - 531.0 * Y[10] + (C[13] / C[11]);
  C[14] = (Y[12] / Y[14]) - (Y[12] / 0.52);
  C[15] = 1.104 - 0.72 * Y[14];
  C[16] = Y[8] + x5;

  /* objective function */
  f[0] = 0.0000005843 * Y[16] - 0.000117 * Y[13] - 0.1365 - 0.00002358 * Y[12] - 0.000001502 * Y[15] - 0.0321 * Y[11] - 0.004324 * Y[4] - 0.0001 * (C[14] / C[15]) - 37.48 * (Y[1] / C[11]);
  f[0] = -f[0]; /* Max-->Min, Modified by Jane,Nov 22 2005 */
  /* constraints g(x) <= 0 */
  g[0] = -Y[3] + (0.28 / 0.72) * Y[4];
  g[1] = -1.5 * x2 + x3;
  g[2] = -21.0 + 3496.0 * (Y[1] / C[11]);
  g[3] = -(62212.0 / C[16]) + 110.6 + Y[0];
  g[4] = 213.1 - Y[0];
  g[5] = Y[0] - 405.23;
  g[6] = 17.505 - Y[1];
  g[7] = Y[1] - 1053.6667;
  g[8] = 11.275 - Y[2];
  g[9] = Y[2] - 35.03;
  g[10] = 214.228 - Y[3];
  g[11] = Y[3] - 665.585;
  g[12] = 7.458 - Y[4];
  g[13] = Y[4] - 584.463;
  g[14] = 0.961 - Y[5];
  g[15] = Y[5] - 265.916;
  g[16] = 1.612 - Y[6];
  g[17] = Y[6] - 7.046;
  g[18] = 0.146 - Y[7];
  g[19] = Y[7] - 0.222;
  g[20] = 107.99 - Y[8];
  g[21] = Y[8] - 273.366;
  g[22] = 922.693 - Y[9];
  g[23] = Y[9] - 1286.105;
  g[24] = 926.832 - Y[10];
  g[25] = Y[10] - 1444.046;
  g[26] = 18.766 - Y[11];
  g[27] = Y[11] - 537.141;
  g[28] = 1072.163 - Y[12];
  g[29] = Y[12] - 3247.039;
  g[30] = 8961.448 - Y[13];
  g[31] = Y[13] - 26844.086;
  g[32] = 0.063 - Y[14];
  g[33] = Y[14] - 0.386;
  g[34] = 71084.33 - Y[15];
  g[35] = Y[15] - 140000.0;
  g[36] = 2802713.0 - Y[16];
  g[37] = Y[16] - 12146108.0;
}

DLLIMPORT void
g17 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  double f1, f2, x1, x2, x3, x4, x5, x6;
  double aux1, aux2, aux5, aux4;

  x1 = x[0];
  x2 = x[1];
  x3 = x[2];
  x4 = x[3];
  x5 = x[4];
  x6 = x[5];

  aux1 = 300.0 - (((x3 * x4) * cos (1.48477 - x6)) - ((0.90798 * pow (x3, 2.0)) * cos (1.47588))) / 131.078;
  aux2 = -(((x3 * x4)  * cos (1.48477 + x6)) - ((0.90798 * pow (x4, 2.0))  * cos (1.47588)))/ 131.078;
  aux5 = -(((x3 * x4)  * sin (1.48477 + x6)) - ((0.90798 * pow (x4, 2.0))  * sin (1.47588)))/ 131.078;
  aux4 = 200.0 - (((x3 * x4)  * sin (1.48477 - x6)) - ((0.90798 * pow (x3, 2.0))  * sin (1.47588)))/ 131.078;


  /* objective fucntion */
  if (x1 >= 0.0 && x1 < 300.0)
    {
      f1 = 30.0 * aux1;
    }
  else
    {
      if (x1 >= 300.0 && x1 <= 400.0)
	{
	  f1 = 31.0 * aux1;
	}
    }
  if (x2 >= 0.0 && x2 < 100.0)
    {
      f2 = 28.0 * aux2;
    }
  else
    {
      if (x2 >= 100.0 && x2 < 200.0)
	{
	  f2 = 29.0 * aux2;
	}
      else
	{
	  if (x2 >= 200.0 && x2 <= 1000.0)
	    {
	      f2 = 30.0 * aux2;
	    }
	}
    }
  f[0] = f1 + f2;

  /* constraint function h = 0 */
  h[0] = aux1 - x1;
  h[1] = aux2 - x2;
  h[2] = aux5 - x5;
  h[3] = aux4;
}

DLLIMPORT void
g18 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = 0.5 * (x[0] * x[3] - x[1] * x[2] + x[2] * x[8] - x[4] * x[8] + x[4] * x[7] - x[5] * x[6]);
  f[0] = -f[0]; /* Max-->Min, Modified by Jane,Nov 22 2005 */
  /* constraint function g <= 0 */
  g[0] = -1.0 + pow (x[2], 2.0) + pow (x[3], 2.0);
  g[1] = -1.0 + pow (x[8], 2.0);
  g[2] = -1.0 + pow (x[4], 2.0) + pow (x[5], 2.0);
  g[3] = -1.0 + pow (x[0], 2.0) + pow (x[1] - x[8], 2.0);
  g[4] = -1.0 + pow (x[0] - x[4], 2.0) + pow (x[1] - x[5], 2.0);
  g[5] = -1.0 + pow (x[0] - x[6], 2.0) + pow (x[1] - x[7], 2.0);
  g[6] = -1.0 + pow (x[2] - x[4], 2.0) + pow (x[3] - x[5], 2.0);
  g[7] = -1.0 + pow (x[2] - x[6], 2.0) + pow (x[3] - x[7], 2.0);
  g[8] = -1.0 + pow (x[6], 2.0) + pow (x[7] - x[8], 2.0);/* bug fixed by Jane,Nov 25 2005 */
  g[9] = -x[0] * x[3] + x[1] * x[2];
  g[10] = -x[2] * x[8];
  g[11] = x[4] * x[8];
  g[12] = -x[4] * x[7] + x[5] * x[6];
}

DLLIMPORT void
g19 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  int i, j;
  double sum1 = 0.0, sum2 = 0.0, sum3 = 0.0;
  double A[10][5] = { -16.0, 2.0, 0.0, 1.0, 0.0,
		      0.0, -2.0, 0.0, 0.4, 2.0,  
		      -3.5, 0.0, 2.0, 0.0, 0.0,
		      0.0, -2.0, 0.0, -4.0, -1.0,
		      0.0, -9.0, -2.0, 1.0, -2.8,
		      2.0, 0.0, -4.0, 0.0, 0.0,
		      -1.0, -1.0, -1.0, -1.0, -1.0,
		      -1.0, -2.0, -3.0, -2.0, -1.0,
		      1.0, 2.0, 3.0, 4.0, 5.0,
		      1.0, 1.0, 1.0, 1.0, 1.0
  };
  double B[10] = { -40.0, -2.0, -0.25, -4.0, -4.0, -1.0, -40.0, -60.0, 5.0, 1.0 };
  double C[5][5] = { 30.0, -20.0, -10.0, 32.0, -10.0, 
		     -20.0, 39.0, -6.0, -31.0, 32.0,
		     -10.0, -6.0, 10.0, -6.0, -10.0,
		     32.0, -31.0, -6.0, 39.0, -20.0,
		     -10.0, 32.0, -10.0, -20.0, 30.0
  };
  double D[5] = { 4.0, 8.0, 10.0, 6.0, 2.0 };
  double E[5] = { -15.0, -27.0, -36.0, -18.0, -12.0 };

  /* objective function */
  for (i = 0; i < 10; i++)
    {
      sum1 += B[i] * x[i];
    }
  for (i = 0; i < 5; i++)
    {
      for (j = 0; j < 5; j++)
	{
	  sum2 += C[i][j] * x[10 + i] * x[10 + j];
	}
    }
  for (i = 0; i < 5; i++)
    {
      sum3 += D[i] * pow (x[10 + i], 3.0);
    }
  f[0] = sum1 - sum2 - 2.0 * sum3;
  f[0] = -f[0]; /* Max-->Min, Modified by Jane,Nov 22 2005 */
  /* constraints g <= 0 */
  for (j = 0; j < 5; j++)
    {
      sum1 = 0.0;
      for (i = 0; i < 5; i++)
	sum1 += C[i][j] * x[10 + i];
      sum2 = 0.0;
      for (i = 0; i < 10; i++)
	sum2 += A[i][j] * x[i];
      g[j] = -((2.0 * sum1) + (3.0 * D[j] * pow (x[10 + j], 2.0)) + E[j] - sum2);
    }
}

DLLIMPORT void
g20 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  double sum1, sum2, sumtotal;
  int i, j;
  double A[24] = { 0.0693, 0.0577, 0.05, 0.2, 0.26, 0.55, 0.06, 0.1, 0.12, 0.18, 0.1, 0.09, 0.0693, 0.0577, 0.05, 0.2, 0.26, 0.55, 0.06, 0.1, 0.12, 0.18, 0.1, 0.09 }; 
  double B[24] = { 44.094, 58.12, 58.12, 137.4, 120.9, 170.9, 62.501, 84.94, 133.425, 82.507, 46.07, 60.097, 44.094, 58.12, 58.12, 137.4, 120.9, 170.9, 62.501, 84.94, 133.425, 82.507, 46.07, 60.097 };
  double C[12] = { 123.7, 31.7, 45.7, 14.7, 84.7, 27.7, 49.7, 7.1, 2.1, 17.7, 0.85, 0.64 };
  double D[12] = { 31.244, 36.12, 34.784, 92.7, 82.7, 91.6, 56.708, 82.7, 80.8, 64.517, 49.4, 49.1 };
  double E[6] = { 0.1, 0.3, 0.4, 0.3, 0.6, 0.3 };

  /* objective function */
  f[0] = 0.0;
  for (j = 0; j < 24; j++)
    {
      f[0] += A[j] * x[j];
    }
  /* constraints h(x) = 0 and g(x) <= 0 */
  sum1 = 0.0;
  for (j = 0; j < 12; j++)
    sum1 += x[j] / B[j];
  sum2 = 0.0;
  for (j = 12; j < 24; j++)
    sum2 += x[j] / B[j];
  for (i = 0; i < 12; i++)
    h[i] = (x[i + 12] / (B[i + 12] * sum2)) - ((C[i] * x[i]) / (40.0 * B[i] * sum1));
  sumtotal = 0.0;
  for (j = 0; j < 24; j++)
    sumtotal += x[j];
  h[12] = sumtotal - 1.0;
  sum1 = 0.0;
  for (j = 0; j < 12; j++)
    sum1 += x[j] / D[j];
  sum2 = 0.0;
  for (j = 12; j < 24; j++)
    sum2 += x[j] / B[j];
  h[13] = sum1 + (0.7302 * 530.0 * (14.7 / 40)) * sum2 - 1.671;
  for (j = 0; j < 3; j++)
    g[j] = (x[j] + x[j + 12]) / (sumtotal + E[j]);
  for (j = 3; j < 6; j++)
    g[j] = (x[j + 3] + x[j + 15]) / (sumtotal + E[j]);
}

DLLIMPORT void
g21 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  int i;

  /* objective function */
  f[0] = x[0];
  /* constraint functions h(x) and g(x) <= 0 */
  g[0] = -x[0] + 35.0 * pow (x[1], 0.6) + 35.0 * pow (x[2], 0.6);
  h[0] = -300.0 * x[2] + 7500 * x[4] - 7500 * x[5] - 25.0 * x[3] * x[4] + 25.0 * x[3] * x[5] + x[2] * x[3];
  h[1] = 100.0 * x[1] + 155.365 * x[3] + 2500 * x[6] - x[1] * x[3] - 25.0 * x[3] * x[6] - 15536.5;
  h[2] = -x[4] + log (-x[3] + 900.0);
  h[3] = -x[5] + log (x[3] + 300.0);
  h[4] = -x[6] + log (-2.0 * x[3] + 700.0);
}

DLLIMPORT void
g22 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  int i;

  /* objective function */
  f[0] = x[0];
  /* constraint functions */
  g[0] = -x[0] + pow (x[1], 0.6) + pow (x[2], 0.6) + pow (x[3], 0.6);
  h[0] = x[4] - 100000.0 * x[7] + 10000000.0;
  h[1] = x[5] + 100000.0 * x[7] - 100000.0 * x[8];
  h[2] = x[6] + 100000.0 * x[8] - 50000000.0;
  h[3] = x[4] + 100000.0 * x[9] - 33000000.0;
  h[4] = x[5] + 100000 * x[10] - 44000000.0;
  h[5] = x[6] + 100000 * x[11] - 66000000.0;
  h[6] = x[4] - 120.0 * x[1] * x[12];
  h[7] = x[5] - 80.0 * x[2] * x[13];
  h[8] = x[6] - 40.0 * x[3] * x[14];
  h[9] = x[7] - x[10] + x[15];
  h[10] = x[8] - x[11] + x[16];
  h[11] = -x[17] + log (x[9] - 100.0);/*bug fixed by Jane,Nov 25 2005*/
  h[12] = -x[18] + log (-x[7] + 300.0);
  h[13] = -x[19] + log (x[15]);
  h[14] = -x[20] + log (-x[8] + 400.0);
  h[15] = -x[21] + log (x[16]);
  h[16] = -x[7] - x[9] + x[12] * x[17] - x[12] * x[18] + 400.0;
  h[17] = x[7] - x[8] - x[10] + x[13] * x[19] - x[13] * x[20] + 400.0;
  h[18] = x[8] - x[11] - 4.60517 * x[14] + x[14] * x[21] + 100.0;
}

DLLIMPORT void
g23 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  int i;

  /* objective function */
  f[0] = -9.0 * x[4] - 15.0 * x[7] + 6.0 * x[0] + 16.0 * x[1] + 10.0 * (x[5] + x[6]);
  /* constraint function */
  g[0] = x[8] * x[2] + 0.02 * x[5] - 0.025 * x[4];
  g[1] = x[8] * x[3] + 0.02 * x[6] - 0.015 * x[7];
  h[0] = x[0] + x[1] - x[2] - x[3];
  h[1] = 0.03 * x[0] + 0.01 * x[1] - x[8] * (x[2] + x[3]);
  h[2] = x[2] + x[5] - x[4];
  h[3] = x[3] + x[6] - x[7];
}

DLLIMPORT void
g24 (double *x, double *f, double *g, double *h, int nx, int nf, int ng, int nh)
{
  /* objective function */
  f[0] = -x[0] - x[1];
  /* constraint function */
  g[0] = -2.0 * pow (x[0], 4.0) + 8.0 * pow (x[0], 3.0) - 8.0 * pow (x[0], 2.0) + x[1] - 2.0;
  g[1] = -4.0 * pow (x[0], 4.0) + 32.0 * pow (x[0], 3.0) - 88.0 * pow (x[0], 2.0) + 96.0 * x[0] + x[1] - 36.0;
}



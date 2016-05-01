/*
 * Copyright © 2005 The Walking Fish Group (WFG).
 *
 * This material is provided "as is", with no warranty expressed or implied.
 * Any use is at your own risk. Permission to use or copy this software for
 * any purpose is hereby granted without fee, provided this notice is
 * retained on all copies. Permission to modify the code and to distribute
 * modified code is granted, provided a notice that the code was modified is
 * included with the above copyright notice.
 *
 * http://www.wfg.csse.uwa.edu.au/
 */


/*
 * TransFunctions.h
 *
 * Defines the transformation functions used by the WFG toolkit. For further
 * documentation, including the effects and arguments of each transformation
 * function, refer to Table 2 the EMO 2005 paper (available from the WFG web
 * site).
 */


#ifndef TRANS_FUNCTIONS_H
#define TRANS_FUNCTIONS_H


//// Standard includes. /////////////////////////////////////////////////////

#include <vector>


//// Definitions/namespaces. ////////////////////////////////////////////////

namespace WFG
{

namespace Toolkit
{

namespace TransFunctions
{

//** The polynomial bias transformation function. ***************************
double b_poly( const double& y, const double& alpha );

//** The flat region bias transformation function. **************************
double b_flat
(
  const double& y,
  const double& A,
  const double& B,
  const double& C
);

//** The parameter dependent bias transformation function. ******************
double b_param
(
  const double& y,
  const double& u,
  const double& A,
  const double& B,
  const double& C
);

//** The linear shift transformation function. ******************************
double s_linear( const double& y, const double& A );

//** The deceptive shift transformation function. ***************************
double s_decept
(
  const double& y,
  const double& A,
  const double& B,
  const double& C
);

//** The multi-modal shift transformation function. *************************
double s_multi
(
  const double& y,
  const int A,
  const double& B,
  const double& C
);

//** The weighted sum reduction transformation function. ********************
double r_sum
(
  const std::vector< double >& y,
  const std::vector< double >& w
);

//** The non-separable reduction transformation function. *******************
double r_nonsep( const std::vector< double >& y, const int A );

}  // TransFunctions namespace

}  // Toolkit namespace

}  // WFG namespace

#endif
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
 * ShapeFunctions.h
 *
 * Defines the shape functions used by the WFG toolkit. For further
 * documentation, including the nature and arguments of each shape function,
 * refer to Table 1 of the EMO 2005 paper (available from the WFG web site).
 */


#ifndef SHAPE_FUNCTIONS_H
#define SHAPE_FUNCTIONS_H


//// Standard includes. /////////////////////////////////////////////////////

#include <vector>


//// Definitions/namespaces. ////////////////////////////////////////////////

namespace WFG
{

namespace Toolkit
{

namespace ShapeFunctions
{

//** The linear shape function. (m is indexed from 1.) **********************
double linear ( const std::vector< double >& x, const int m );

//** The convex shape function. (m is indexed from 1.) **********************
double convex ( const std::vector< double >& x, const int m );

//** The concave shape function. (m is indexed from 1.) *********************
double concave( const std::vector< double >& x, const int m );

//** The mixed convex/concave shape function. *******************************
double mixed
(
  const std::vector< double >& x,
  const int A,
  const double& alpha
);

//** The disconnected shape function. ***************************************
double disc
(
  const std::vector< double >& x,
  const int A,
  const double& alpha,
  const double& beta
);

}  // ShapeFunctions namespace

}  // Toolkit namespace

}  // WFG namespace


#endif
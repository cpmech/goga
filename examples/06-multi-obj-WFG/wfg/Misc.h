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
 * Misc.h
 *
 * Defines some general convenience functions and values.
 */


#ifndef MISC_H
#define MISC_H


//// Standard includes. /////////////////////////////////////////////////////

#include <vector>


//// Definitions/namespaces. ////////////////////////////////////////////////

namespace WFG
{

namespace Toolkit
{

namespace Misc
{


//// Constants. /////////////////////////////////////////////////////////////

extern const double PI;


//// Functions. /////////////////////////////////////////////////////////////

//** Used to correct values in [-epislon,0] to 0, and [1,epsilon] to 1. *****
double correct_to_01( const double& a, const double& epsilon = 1.0e-10 );

//** Returns true if all elements of "x" are in [0,1], false otherwise. *****
bool vector_in_01( const std::vector< double >& x );

}  // ShapeFunctions namespace

}  // Toolkit namespace

}  // WFG namespace


#endif
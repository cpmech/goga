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
 * ExampleShapes.h
 *
 * Defines functions that take final transition vectors as arguments, from
 * which overall fitness values are calculated. In effect, each function
 * degenerates the transition vector as appropriate, calculates shape
 * function values, and then scales and shifts the resultant vector to
 * produce the desired fitness vector. Functions are provided for each of the
 * problems described in the EMO 2005 paper, namely WFG1--WFG9 and I1--I5.
 * For specific details, refer to the EMO 2005 paper (available from the WFG
 * web site).
 */


#ifndef EXAMPLE_SHAPES_H
#define EXAMPLE_SHAPES_H


//// Standard includes. /////////////////////////////////////////////////////

#include <vector>


//// Definitions/namespaces. ////////////////////////////////////////////////

namespace WFG
{

namespace Toolkit
{

namespace Examples
{

namespace Shapes
{

//** Given the last transition vector, get the fitness values for WFG1. *****
std::vector< double > WFG1_shape( const std::vector< double >& t_p );

//** Given the last transition vector, get the fitness values for WFG2. *****
std::vector< double > WFG2_shape( const std::vector< double >& t_p );

//** Given the last transition vector, get the fitness values for WFG3. *****
std::vector< double > WFG3_shape( const std::vector< double >& t_p );

//** Given the last transition vector, get the fitness values for WFG4. *****
std::vector< double > WFG4_shape( const std::vector< double >& t_p );

//** Given the last transition vector, get the fitness values for I1. *******
std::vector< double > I1_shape( const std::vector< double >& t_p );

}  // Shapes namespace

}  // Examples namespace

}  // Toolkit namespace

}  // WFG namespace

#endif
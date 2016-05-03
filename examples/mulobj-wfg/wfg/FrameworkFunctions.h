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
 * FrameworkFunctions.h
 *
 * Defines the "framework" functions of the WFG toolkit, effectively
 * mirroring the format described at the start of Section 4 of the EMO 2005
 * paper, updated as per Section VIII of the IEEE TEC review paper.
 */


/*
 * 2006-03-28: Added distance scaling constant D to calculate_f().
 */


#ifndef FRAMEWORK_FUNCTIONS_H
#define FRAMEWORK_FUNCTIONS_H


//// Standard includes. /////////////////////////////////////////////////////

#include <vector>


//// Definitions/namespaces. ////////////////////////////////////////////////

namespace WFG
{

namespace Toolkit
{

namespace FrameworkFunctions
{

//** Normalise the elements of "z" to the domain [0,1]. *********************
std::vector< double > normalise_z
(
  const std::vector< double >& z,
  const std::vector< double >& z_max
);

//** Degenerate the values of "t_p" based on the degeneracy vector "A". *****
std::vector< double > calculate_x
(
  const std::vector< double >& t_p,
  const std::vector< short >& A
);

//** Calculate the fitness vector using the distance scaling constant D, ****
//** the distance parameter in "x", the shape function values in "h",    ****
//** and the scaling constants in "S".                                   ****
std::vector< double > calculate_f
(
  const double&                D,
  const std::vector< double >& x,
  const std::vector< double >& h,
  const std::vector< double >& S
);

}  // FrameworkFunctions namespace

}  // Toolkit namespace

}  // WFG namespace

#endif
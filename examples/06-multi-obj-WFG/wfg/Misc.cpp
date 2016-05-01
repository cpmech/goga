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
 * Misc.cpp
 *
 * Implementation of Misc.h.
 */


#include "Misc.h"


//// Standard includes. /////////////////////////////////////////////////////

#include <vector>
#include <cassert>


//// Used namespaces. ///////////////////////////////////////////////////////

using namespace WFG::Toolkit;
using std::vector;


//// Implemented constants. /////////////////////////////////////////////////

const double WFG::Toolkit::Misc::PI = 3.1415926535897932384626433832795;


//// Implemented functions. /////////////////////////////////////////////////

double Misc::correct_to_01( const double& a, const double& epsilon )
{
  assert( epsilon >= 0.0 );

  const double min = 0.0;
  const double max = 1.0;

  const double min_epsilon = min - epsilon;
  const double max_epsilon = max + epsilon;

  if ( a <= min && a >= min_epsilon )
  {
    return min;
  }
  else if ( a >= max && a <= max_epsilon )
  {
    return max;
  }
  else
  {
    return a;
  }
}

bool Misc::vector_in_01( const vector< double >& x )
{
  for( int i = 0; i < static_cast< int >( x.size() ); i++ )
  {
    if( x[i] < 0.0 || x[i] > 1.0 )
    {
      return false;
    }
  }

  return true;
}
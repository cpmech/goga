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
 * Implementation of FrameworkFunctions.h.
 */


/*
 * 2006-03-28: Added distance scaling constant D to calculate_f().
 */


#include "FrameworkFunctions.h"


//// Standard includes. /////////////////////////////////////////////////////

#include <cassert>


//// Toolkit includes. //////////////////////////////////////////////////////

#include "Misc.h"


//// Used namespaces. ///////////////////////////////////////////////////////

using namespace WFG::Toolkit;
using namespace WFG::Toolkit::Misc;
using std::vector;


//// Implemented functions. /////////////////////////////////////////////////

vector< double > FrameworkFunctions::normalise_z
(
  const vector< double >& z,
  const vector< double >& z_max
)
{
  vector< double > result;

  for( int i = 0; i < static_cast< int >( z.size() ); i++ )
  {
    assert( z[i] >= 0.0 );
    assert( z[i] <= z_max[i] );
    assert( z_max[i] > 0.0 );

    result.push_back( z[i] / z_max[i] );
  }

  return result;
}

vector< double >  FrameworkFunctions::calculate_x
(
  const vector< double >& t_p,
  const vector< short >& A
)
{
  assert( vector_in_01( t_p ) );
  assert( t_p.size() != 0 );
  assert( A.size() == t_p.size()-1 );

  vector< double > result;

  for( int i = 0; i < static_cast< int >( t_p.size() ) - 1; i++ )
  {
    assert( A[i] == 0 || A[i] == 1 );

    const double tmp1 = std::max< double >( t_p.back(), A[i] );
    result.push_back( tmp1*( t_p[i] - 0.5 ) + 0.5 );
  }

  result.push_back( t_p.back() );

  return result;
}

vector< double >  FrameworkFunctions::calculate_f
(
  const double&           D,
  const vector< double >& x,
  const vector< double >& h,
  const vector< double >& S
)
{
  assert( D > 0.0 );
  assert( vector_in_01( x ) );
  assert( vector_in_01( h ) );
  assert( x.size() == h.size() );
  assert( h.size() == S.size() );

  vector< double > result;

  for( int i = 0; i < static_cast< int >( h.size() ); i++ )
  {
    assert( S[i] > 0.0 );

    result.push_back( D*x.back() + S[i]*h[i] );
  }

  return result;
}
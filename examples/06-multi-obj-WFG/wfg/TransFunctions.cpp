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
 * TransFunctions.cpp
 *
 * Implementation of TransFunctions.h.
 *
 * Changelog:
 *   2005.06.07 (Simon Huband)
 *     - Replaced a call to abs by fabs (this was causing problems).
 */


#include "TransFunctions.h"


//// Standard includes. /////////////////////////////////////////////////////

#include <cassert>
#include <cmath>


//// Toolkit includes. //////////////////////////////////////////////////////

#include "Misc.h"


//// Used namespaces. ///////////////////////////////////////////////////////

using namespace WFG::Toolkit;
using namespace WFG::Toolkit::Misc;
using WFG::Toolkit::Misc::correct_to_01;
using std::vector;


//// Local functions. ///////////////////////////////////////////////////////

namespace
{

//** Calculate the minimum of two doubles. **********************************
inline double min( const double& a, const double& b )
{
  return std::min< const double >( a, b );
}

}  // unnamed namespace


//// Implemented functions. /////////////////////////////////////////////////

double TransFunctions::b_poly( const double& y, const double& alpha )
{
  assert( y >= 0.0 );
  assert( y <= 1.0 );
  assert( alpha >  0.0 );
  assert( alpha != 1.0 );

  return correct_to_01( pow( y, alpha ) );
}

double TransFunctions::b_flat
(
  const double& y,
  const double& A,
  const double& B,
  const double& C
)
{
  assert( y >= 0.0 );
  assert( y <= 1.0 );
  assert( A >= 0.0 );
  assert( A <= 1.0 );
  assert( B >= 0.0 );
  assert( B <= 1.0 );
  assert( C >= 0.0 );
  assert( C <= 1.0 );
  assert( B < C );
  assert( B != 0.0 || A == 0.0 );
  assert( B != 0.0 || C != 1.0 );
  assert( C != 1.0 || A == 1.0 );
  assert( C != 1.0 || B != 0.0 );

  const double tmp1 = min( 0.0, floor( y-B ) ) * A*( B-y )/B;
  const double tmp2 = min( 0.0, floor( C-y ) ) * ( 1.0-A )*( y-C )/( 1.0-C );

  return correct_to_01( A+tmp1-tmp2 );
}

double TransFunctions::b_param
(
  const double& y,
  const double& u,
  const double& A,
  const double& B,
  const double& C
)
{
  assert( y >= 0.0 );
  assert( y <= 1.0 );
  assert( u >= 0.0 );
  assert( u <= 1.0 );
  assert( A > 0.0 );
  assert( A < 1.0 );
  assert( B > 0.0 );
  assert( B < C );

  const double v = A - ( 1.0-2.0*u )*fabs( floor( 0.5-u )+A );

  return correct_to_01( pow( y, B + ( C-B )*v ) );
}

double TransFunctions::s_linear( const double& y, const double& A )
{
  assert( y >= 0.0 );
  assert( y <= 1.0 );
  assert( A > 0.0 );
  assert( A < 1.0 );

  return correct_to_01( fabs( y-A )/fabs( floor( A-y )+A ) );
}

double TransFunctions::s_decept
(
  const double& y,
  const double& A,
  const double& B,
  const double& C
)
{
  assert( y >= 0.0 );
  assert( y <= 1.0 );
  assert( A > 0.0 );
  assert( A < 1.0 );
  assert( B > 0.0 );
  assert( B < 1.0 );
  assert( C > 0.0 );
  assert( C < 1.0 );
  assert( A - B > 0.0 );
  assert( A + B < 1.0 );

  const double tmp1 = floor( y-A+B )*( 1.0-C+( A-B )/B )/( A-B );
  const double tmp2 = floor( A+B-y )*( 1.0-C+( 1.0-A-B )/B )/( 1.0-A-B );

  return correct_to_01( 1.0 + ( fabs( y-A )-B )*( tmp1 + tmp2 + 1.0/B ) );
}

double TransFunctions::s_multi
(
  const double& y,
  const int A,
  const double& B,
  const double& C
)
{
  assert( y >= 0.0 );
  assert( y <= 1.0 );
  assert( A >= 1 );
  assert( B >= 0.0 );
  assert( ( 4.0*A+2.0 )*Misc::PI >= 4.0*B );
  assert( C > 0.0 );
  assert( C < 1.0 );

  const double tmp1 = fabs( y-C )/( 2.0*( floor( C-y )+C ) );
  const double tmp2 = ( 4.0*A+2.0 )*Misc::PI*( 0.5-tmp1 );

  return correct_to_01( ( 1.0 + cos( tmp2 ) + 4.0*B*pow( tmp1, 2.0 ) )/( B+2.0 ) );
}

double TransFunctions::r_sum
(
  const vector< double >& y,
  const vector< double >& w
)
{
  assert( y.size() != 0        );
  assert( w.size() == y.size() );
  assert( Misc::vector_in_01( y ) );

  double numerator   = 0.0;
  double denominator = 0.0;

  for( int i = 0; i < static_cast< int >( y.size() ); i++ )
  {
    assert( w[i] > 0.0 );

    numerator   += w[i]*y[i];
    denominator += w[i];
  }

  return correct_to_01( numerator / denominator );
}

double TransFunctions::r_nonsep( const std::vector< double >& y, const int A )
{
  const int y_len = static_cast< int >( y.size() );

  assert( y_len != 0 );
  assert( Misc::vector_in_01( y ) );
  assert( A >= 1 );
  assert( A <= y_len );
  assert( y.size() % A == 0 );

  double numerator = 0.0;

  for( int j = 0; j < y_len; j++ )
  {
    numerator += y[j];

    for( int k = 0; k <= A-2; k++ )
    {
      numerator += fabs( y[j] - y[( j+k+1 ) % y_len] );
    }
  }

  const double tmp = ceil( A/2.0 );
  const double denominator = y_len*tmp*( 1.0 + 2.0*A - 2.0*tmp )/A;

  return correct_to_01( numerator / denominator );
}
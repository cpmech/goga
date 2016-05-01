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
 * Implementation of ShapeFunctions.h.
 */


#include "ShapeFunctions.h"


//// Standard includes. /////////////////////////////////////////////////////

#include <cassert>
#include <cmath>


//// Toolkit includes. //////////////////////////////////////////////////////

#include "Misc.h"
#include "Misc.h"


//// Used namespaces. ///////////////////////////////////////////////////////

using namespace WFG::Toolkit;
using namespace WFG::Toolkit::Misc;
using WFG::Toolkit::Misc::correct_to_01;
using std::vector;


//// Local functions. ///////////////////////////////////////////////////////

namespace
{

//** True if all elements of "x" are in [0,1], and m is in [1, x.size()]. ***
bool shape_args_ok( const vector< double >& x, const int m )
{
  const int M = static_cast< int >( x.size() );

  return Misc::vector_in_01( x ) && m >= 1 && m <= M;
}

}  // unnamed namespace


//// Implemented functions. /////////////////////////////////////////////////

double ShapeFunctions::linear( const vector< double >& x, const int m )
{
  assert( shape_args_ok( x, m ) );

  const int M = static_cast< int >( x.size() );
  double result = 1.0;

  for( int i=1; i <= M-m; i++ )
  {
    result *= x[i-1];
  }

  if( m != 1 )
  {
    result *= 1 - x[M-m];
  }

  return correct_to_01( result );
}

double ShapeFunctions::convex( const vector< double >& x, const int m )
{
  assert( shape_args_ok( x, m ) );

  const int M = static_cast< int >( x.size() );
  double result = 1.0;

  for( int i=1; i <= M-m; i++ )
  {
    result *= 1.0 - cos( x[i-1]*Misc::PI/2.0 );
  }

  if( m != 1 )
  {
    result *= 1.0 - sin( x[M-m]*Misc::PI/2.0 );
  }

  return correct_to_01( result );
}

double ShapeFunctions::concave( const vector< double >& x, const int m )
{
  assert( shape_args_ok( x, m ) );

  const int M = static_cast< int >( x.size() );
  double result = 1.0;

  for( int i=1; i <= M-m; i++ )
  {
    result *= sin( x[i-1]*Misc::PI/2.0 );
  }

  if( m != 1 )
  {
    result *= cos( x[M-m]*Misc::PI/2.0 );
  }

  return correct_to_01( result );
}

double ShapeFunctions::mixed
(
  const vector< double >& x,
  const int A,
  const double& alpha
)
{
  assert( Misc::vector_in_01( x ) );
  assert( x.size() != 0   );
  assert( A        >= 1   );
  assert( alpha    >  0.0 );

  const double tmp = 2.0*A*Misc::PI;

  return correct_to_01( pow( 1.0-x[0]-cos( tmp*x[0] + Misc::PI/2.0 )/tmp, alpha ) );
}

double ShapeFunctions::disc
(
  const vector< double >& x,
  const int A,
  const double& alpha,
  const double& beta
)
{
  assert( Misc::vector_in_01( x ) );
  assert( x.size() != 0   );
  assert( A        >= 1   );
  assert( alpha    >  0.0 );
  assert( beta     >  0.0 );

  const double tmp1 = A*pow( x[0], beta )*Misc::PI;
  return correct_to_01( 1.0 - pow( x[0], alpha )*pow( cos( tmp1 ), 2.0 ) );
}
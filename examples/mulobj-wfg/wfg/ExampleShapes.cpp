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
 * Implementation of ExampleShapes.h.
 */


/*
 * 2006-03-28: Updated WFG_calculate_f() and I1_shape to employ a distance
 *             scaling constant D value of 1.0, as per changes introduced in
 *             the IEEE TEC review paper.
 */


#include "ExampleShapes.h"


//// Standard includes. /////////////////////////////////////////////////////

#include <cassert>


//// Toolkit includes. //////////////////////////////////////////////////////

#include "ShapeFunctions.h"
#include "FrameworkFunctions.h"
#include "Misc.h"


//// Used namespaces. ///////////////////////////////////////////////////////

using namespace WFG::Toolkit;
using namespace WFG::Toolkit::Examples;
using namespace WFG::Toolkit::Misc;
using std::vector;


//// Local functions. ///////////////////////////////////////////////////////

namespace
{

//** Construct a vector of length M-1, with values "1,0,0,..." if ***********
//** "degenerate" is true, otherwise with values "1,1,1,..." if   ***********
//** "degenerate" is false.                                       ***********
vector< short > WFG_create_A( const int M, const bool degenerate )
{
  assert( M >= 2 );

  if ( degenerate )
  {
    vector< short > A( M-1, 0 );
    A[0] = 1;

    return A;
  }
  else
  {
    return vector< short >( M-1, 1 );
  }
}

//** Given the vector "x" (the last value of which is the sole distance ****
//** parameter), and the shape function results in "h", calculate the   ****
//** scaled fitness values for a WFG problem.                           ****
vector< double > WFG_calculate_f
(
  const vector< double >& x,
  const vector< double >& h
)
{
  assert( vector_in_01( x ) );
  assert( vector_in_01( h ) );
  assert( x.size() == h.size() );

  const int M = static_cast< int >( h.size() );

  vector< double > S;

  for( int m = 1; m <= M; m++ )
  {
    S.push_back( m*2.0 );
  }

  return FrameworkFunctions::calculate_f( 1.0, x, h, S );
}

}  // unnamed namespace


//// Implemented functions. /////////////////////////////////////////////////

vector< double > Shapes::WFG1_shape( const vector< double >& t_p )
{
  assert( vector_in_01( t_p ) );
  assert( t_p.size() >= 2 );

  const int M = static_cast< int >( t_p.size() );

  const vector< short >&  A = WFG_create_A( M, false );
  const vector< double >& x = FrameworkFunctions::calculate_x( t_p, A );

  vector< double > h;

  for( int m = 1; m <= M-1; m++ )
  {
    h.push_back( ShapeFunctions::convex( x, m ) );
  }
  h.push_back( ShapeFunctions::mixed( x, 5, 1.0 ) );

  return WFG_calculate_f( x, h );
}

vector< double > Shapes::WFG2_shape( const vector< double >& t_p )
{
  assert( vector_in_01( t_p ) );
  assert( t_p.size() >= 2 );

  const int M = static_cast< int >( t_p.size() );

  const vector< short >&  A = WFG_create_A( M, false );
  const vector< double >& x = FrameworkFunctions::calculate_x( t_p, A );

  vector< double > h;

  for( int m = 1; m <= M-1; m++ )
  {
    h.push_back( ShapeFunctions::convex( x, m ) );
  }
  h.push_back( ShapeFunctions::disc( x, 5, 1.0, 1.0 ) );

  return WFG_calculate_f( x, h );
}

vector< double > Shapes::WFG3_shape( const vector< double >& t_p )
{
  assert( vector_in_01( t_p ) );
  assert( t_p.size() >= 2 );

  const int M = static_cast< int >( t_p.size() );

  const vector< short >&  A = WFG_create_A( M, true );
  const vector< double >& x = FrameworkFunctions::calculate_x( t_p, A );

  vector< double > h;

  for( int m = 1; m <= M; m++ )
  {
    h.push_back( ShapeFunctions::linear( x, m ) );
  }

  return WFG_calculate_f( x, h );
}

vector< double > Shapes::WFG4_shape( const vector< double >& t_p )
{
  assert( vector_in_01( t_p ) );
  assert( t_p.size() >= 2 );

  const int M = static_cast< int >( t_p.size() );

  const vector< short >&  A = WFG_create_A( M, false );
  const vector< double >& x = FrameworkFunctions::calculate_x( t_p, A );

  vector< double > h;

  for( int m = 1; m <= M; m++ )
  {
    h.push_back( ShapeFunctions::concave( x, m ) );
  }

  return WFG_calculate_f( x, h );
}

vector< double > Shapes::I1_shape( const vector< double >& t_p )
{
  assert( vector_in_01( t_p ) );
  assert( t_p.size() >= 2 );

  const int M = static_cast< int >( t_p.size() );

  const vector< short >&  A = WFG_create_A( M, false );
  const vector< double >& x = FrameworkFunctions::calculate_x( t_p, A );

  vector< double > h;

  for( int m = 1; m <= M; m++ )
  {
    h.push_back( ShapeFunctions::concave( x, m ) );
  }

  return FrameworkFunctions::calculate_f( 1.0, x, h, vector< double >( M, 1.0 ) );
}
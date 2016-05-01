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
 * ExampleProblems.h
 *
 * Implementation of ExampleProblems.h.
 */


#include "ExampleProblems.h"


//// Standard includes. /////////////////////////////////////////////////////

#include <cassert>


//// Toolkit includes. //////////////////////////////////////////////////////

#include "ExampleTransitions.h"
#include "ExampleShapes.h"


//// Used namespaces. ///////////////////////////////////////////////////////

using namespace WFG::Toolkit;
using namespace WFG::Toolkit::Examples;
using std::vector;


//// Local functions. ///////////////////////////////////////////////////////

namespace
{

//** True if "k" in [1,z.size()), "M" >= 2, and "k" mod ("M"-1) == 0. *******
bool ArgsOK( const vector< double >& z, const int k, const int M )
{
  const int n = static_cast< int >( z.size() );

  return k >= 1 && k < n && M >= 2 && k % ( M-1 ) == 0;
}

//** Reduces each paramer in "z" to the domain [0,1]. ***********************
vector< double > WFG_normalise_z( const vector< double >& z )
{
  vector< double > result;

  for( int i = 0; i < static_cast< int >( z.size() ); i++ )
  {
    const double bound = 2.0*( i+1 );

    assert( z[i] >= 0.0   );
    assert( z[i] <= bound );

    result.push_back( z[i] / bound );
  }

  return result;
}

}  // unnamed namespace


//// Implemented functions. /////////////////////////////////////////////////

vector< double > Problems::WFG1
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG1_t1( y, k );
  y = Transitions::WFG1_t2( y, k );
  y = Transitions::WFG1_t3( y );
  y = Transitions::WFG1_t4( y, k, M );

  return Shapes::WFG1_shape( y );
}

vector< double > Problems::WFG2
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );
  assert( ( static_cast< int >( z.size() )-k ) % 2 == 0 );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG1_t1( y, k );
  y = Transitions::WFG2_t2( y, k );
  y = Transitions::WFG2_t3( y, k, M );

  return Shapes::WFG2_shape( y );
}

vector< double > Problems::WFG3
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );
  assert( ( static_cast< int >( z.size() )-k ) % 2 == 0 );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG1_t1( y, k );
  y = Transitions::WFG2_t2( y, k );
  y = Transitions::WFG2_t3( y, k, M );

  return Shapes::WFG3_shape( y );
}

vector< double > Problems::WFG4
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG4_t1( y );
  y = Transitions::WFG2_t3( y, k, M );

  return Shapes::WFG4_shape( y );
}

vector< double > Problems::WFG5
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG5_t1( y );
  y = Transitions::WFG2_t3( y, k, M );

  return Shapes::WFG4_shape( y );
}

vector< double > Problems::WFG6
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG1_t1( y, k );
  y = Transitions::WFG6_t2( y, k, M );

  return Shapes::WFG4_shape( y );
}

vector< double > Problems::WFG7
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG7_t1( y, k );
  y = Transitions::WFG1_t1( y, k );
  y = Transitions::WFG2_t3( y, k, M );

  return Shapes::WFG4_shape( y );
}

vector< double > Problems::WFG8
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG8_t1( y, k );
  y = Transitions::WFG1_t1( y, k );
  y = Transitions::WFG2_t3( y, k, M );

  return Shapes::WFG4_shape( y );
}

vector< double > Problems::WFG9
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = WFG_normalise_z( z );

  y = Transitions::WFG9_t1( y );
  y = Transitions::WFG9_t2( y, k );
  y = Transitions::WFG6_t2( y, k, M );

  return Shapes::WFG4_shape( y );
}

vector< double > Problems::I1
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = z;

  y = Transitions::I1_t2( y, k );
  y = Transitions::I1_t3( y, k, M );

  return Shapes::I1_shape( y );
}

vector< double > Problems::I2
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = z;

  y = Transitions::I2_t1( y );
  y = Transitions::I1_t2( y, k );
  y = Transitions::I1_t3( y, k, M );

  return Shapes::I1_shape( y );
}

vector< double > Problems::I3
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = z;

  y = Transitions::I3_t1( y );
  y = Transitions::I1_t2( y, k );
  y = Transitions::I1_t3( y, k, M );

  return Shapes::I1_shape( y );
}

vector< double > Problems::I4
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = z;

  y = Transitions::I1_t2( y, k );
  y = Transitions::I4_t3( y, k, M );

  return Shapes::I1_shape( y );
}

vector< double > Problems::I5
(
  const vector< double >& z,
  const int k,
  const int M
)
{
  assert( ArgsOK( z, k, M ) );

  vector< double > y = z;

  y = Transitions::I3_t1( y );
  y = Transitions::I1_t2( y, k );
  y = Transitions::I4_t3( y, k, M );

  return Shapes::I1_shape( y );
}
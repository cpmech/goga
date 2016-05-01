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
 * ExampleTransitions.h
 *
 * Implementation of ExampleTransitions.h.
 */


#include "ExampleTransitions.h"


//// Standard includes. /////////////////////////////////////////////////////

#include <cassert>


//// Toolkit includes. //////////////////////////////////////////////////////

#include "Misc.h"
#include "TransFunctions.h"


//// Used namespaces. ///////////////////////////////////////////////////////

using namespace WFG::Toolkit;
using namespace WFG::Toolkit::Examples;
using namespace WFG::Toolkit::Misc;
using std::vector;


//// Local functions. ///////////////////////////////////////////////////////

namespace
{

//** Construct a vector with the elements v[head], ..., v[tail-1]. **********
vector< double > subvector
(
  const vector< double >& v,
  const int head,
  const int tail
)
{
  assert( head >= 0 );
  assert( head < tail );
  assert( tail <= static_cast< int >( v.size() ) );

  vector< double > result;

  for( int i = head; i < tail; i++ )
  {
    result.push_back( v[i] );
  }

  return result;
}

}  // unnamed namespace


//// Implemented functions. /////////////////////////////////////////////////

vector< double > Transitions::WFG1_t1
(
  const vector< double >& y,
  const int k
)
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );

  vector< double > t;

  for( int i = 0; i < k; i++ )
  {
    t.push_back( y[i] );
  }

  for( int i = k; i < n; i++ )
  {
    t.push_back( TransFunctions::s_linear( y[i], 0.35 ) );
  }

  return t;
}

vector< double > Transitions::WFG1_t2
(
  const vector< double >& y,
  const int k
)
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );

  vector< double > t;

  for( int i = 0; i < k; i++ )
  {
    t.push_back( y[i] );
  }

  for( int i = k; i < n; i++ )
  {
    t.push_back( TransFunctions::b_flat( y[i], 0.8, 0.75, 0.85 ) );
  }

  return t;
}

vector< double > Transitions::WFG1_t3( const vector< double >& y )
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );

  vector< double > t;

  for( int i = 0; i < n; i++ )
  {
    t.push_back( TransFunctions::b_poly( y[i], 0.02 ) );
  }

  return t;
}

vector< double > Transitions::WFG1_t4
(
  const vector< double >& y,
  const int k,
  const int M
)
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );
  assert( M >= 2 );
  assert( k % ( M-1 ) == 0 );

  vector< double > w;

  for( int i = 1; i <= n; i++ )
  {
    w.push_back( 2.0*i );
  }

  vector< double > t;

  for( int i = 1; i <= M-1; i++ )
  {
    const int head = ( i-1 )*k/( M-1 );
    const int tail = i*k/( M-1 );

    const vector< double >& y_sub = subvector( y, head, tail );
    const vector< double >& w_sub = subvector( w, head, tail );

    t.push_back( TransFunctions::r_sum( y_sub, w_sub ) );
  }

  const vector< double >& y_sub = subvector( y, k, n );
  const vector< double >& w_sub = subvector( w, k, n );

  t.push_back( TransFunctions::r_sum( y_sub, w_sub ) );

  return t;
}

vector< double > Transitions::WFG2_t2
(
  const vector< double >& y,
  const int k
)
{
  const int n = static_cast< int >( y.size() );
  const int l = n-k;

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );
  assert( l % 2 == 0 );

  vector< double > t;

  for( int i = 0; i < k; i++ )
  {
    t.push_back( y[i] );
  }

  for( int i = k+1; i <= k+l/2; i++ )
  {
    const int head = k+2*( i-k )-2;
    const int tail = k+2*( i-k );

    t.push_back( TransFunctions::r_nonsep( subvector( y, head, tail ), 2 ) );
  }

  return t;
}

vector< double > Transitions::WFG2_t3
(
  const vector< double >& y,
  const int k,
  const int M
)
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );
  assert( M >= 2 );
  assert( k % ( M-1 ) == 0 );

  const vector< double > w( n, 1.0 );

  vector< double > t;

  for( int i = 1; i <= M-1; i++ )
  {
    const int head = ( i-1 )*k/( M-1 );
    const int tail = i*k/( M-1 );

    const vector< double >& y_sub = subvector( y, head, tail );
    const vector< double >& w_sub = subvector( w, head, tail );

    t.push_back( TransFunctions::r_sum( y_sub, w_sub ) );
  }

  const vector< double >& y_sub = subvector( y, k, n );
  const vector< double >& w_sub = subvector( w, k, n );

  t.push_back( TransFunctions::r_sum( y_sub, w_sub ) );

  return t;
}

vector< double > Transitions::WFG4_t1( const vector< double >& y )
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );

  vector< double > t;

  for( int i = 0; i < n; i++ )
  {
    t.push_back( TransFunctions::s_multi( y[i], 30, 10, 0.35 ) );
  }

  return t;
}

vector< double > Transitions::WFG5_t1( const vector< double >& y )
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );

  vector< double > t;

  for( int i = 0; i < n; i++ )
  {
    t.push_back( TransFunctions::s_decept( y[i], 0.35, 0.001, 0.05 ) );
  }

  return t;
}

vector< double > Transitions::WFG6_t2
(
  const vector< double >& y,
  const int k,
  const int M
)
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );
  assert( M >= 2 );
  assert( k % ( M-1 ) == 0 );

  vector< double > t;

  for( int i = 1; i <= M-1; i++ )
  {
    const int head = ( i-1 )*k/( M-1 );
    const int tail = i*k/( M-1 );

    const vector< double >& y_sub = subvector( y, head, tail );

    t.push_back( TransFunctions::r_nonsep( y_sub, k/( M-1 ) ) );
  }

  const vector< double >& y_sub = subvector( y, k, n );

  t.push_back( TransFunctions::r_nonsep( y_sub, n-k ) );

  return t;
}

vector< double > Transitions::WFG7_t1
(
  const vector< double >& y,
  const int k
)
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );

  const vector< double > w( n, 1.0 );

  vector< double > t;

  for( int i = 0; i < k; i++ )
  {
    const vector< double >& y_sub = subvector( y, i+1, n );
    const vector< double >& w_sub = subvector( w, i+1, n );

    const double u = TransFunctions::r_sum( y_sub, w_sub );

    t.push_back( TransFunctions::b_param( y[i], u, 0.98/49.98, 0.02, 50 ) );
  }

  for( int i = k; i < n; i++ )
  {
    t.push_back( y[i] );
  }

  return t;
}

vector< double > Transitions::WFG8_t1
(
  const vector< double >& y,
  const int k
)
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );

  const vector< double > w( n, 1.0 );

  vector< double > t;

  for( int i = 0; i < k; i++ )
  {
    t.push_back( y[i] );
  }

  for( int i = k; i < n; i++ )
  {
    const vector< double >& y_sub = subvector( y, 0, i );
    const vector< double >& w_sub = subvector( w, 0, i );

    const double u = TransFunctions::r_sum( y_sub, w_sub );

    t.push_back( TransFunctions::b_param( y[i], u, 0.98/49.98, 0.02, 50 ) );
  }

  return t;
}

vector< double > Transitions::WFG9_t1( const vector< double >& y )
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );

  const vector< double > w( n, 1.0 );

  vector< double > t;

  for( int i = 0; i < n-1; i++ )
  {
    const vector< double >& y_sub = subvector( y, i+1, n );
    const vector< double >& w_sub = subvector( w, i+1, n );

    const double u = TransFunctions::r_sum( y_sub, w_sub );

    t.push_back( TransFunctions::b_param( y[i], u, 0.98/49.98, 0.02, 50 ) );
  }

  t.push_back( y.back() );

  return t;
}

vector< double > Transitions::WFG9_t2
(
  const vector< double >& y,
  const int k
)
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );
  assert( k >= 1 );
  assert( k <  n );

  vector< double > t;

  for( int i = 0; i < k; i++ )
  {
    t.push_back( TransFunctions::s_decept( y[i], 0.35, 0.001, 0.05 ) );
  }

  for( int i = k; i < n; i++ )
  {
    t.push_back( TransFunctions::s_multi( y[i], 30, 95, 0.35 ) );
  }

  return t;
}

vector< double > Transitions::I1_t2
(
  const vector< double >& y,
  const int k
)
{
  return WFG1_t1( y, k );
}

vector< double > Transitions::I1_t3
(
  const vector< double >& y,
  const int k,
  const int M
)
{
  return WFG2_t3( y, k, M );
}

vector< double > Transitions::I2_t1( const vector< double >& y )
{
  return WFG9_t1( y );
}

vector< double > Transitions::I3_t1( const vector< double >& y )
{
  const int n = static_cast< int >( y.size() );

  assert( vector_in_01( y ) );

  const vector< double > w( n, 1.0 );

  vector< double > t;

  t.push_back( y.front() );

  for( int i = 1; i < n; i++ )
  {
    const vector< double >& y_sub = subvector( y, 0, i );
    const vector< double >& w_sub = subvector( w, 0, i );

    const double u = TransFunctions::r_sum( y_sub, w_sub );

    t.push_back( TransFunctions::b_param( y[i], u, 0.98/49.98, 0.02, 50 ) );
  }

  return t;
}

vector< double > Transitions::I4_t3
(
  const vector< double >& y,
  const int k,
  const int M
)
{
  return WFG6_t2( y, k, M );
}
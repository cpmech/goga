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
 * main.cpp
 *
 * This file contains a simple driver for testing the WFG problems and
 * transformation functions from the WFG test problem toolkit.
 *
 * Changelog:
 *   2005.06.01 (Simon Huband)
 *     - Corrected commments to indicate k and l are the number of position
 *       and distance parameters, respectively (not the other way around).
 */


//// Standard includes. /////////////////////////////////////////////////////

#include <iostream>
#include <cstdlib>
#include <string>
#include <sstream>
#include <vector>
#include <cassert>
#include <cmath>


//// Toolkit includes. //////////////////////////////////////////////////////

#include "ExampleProblems.h"
#include "TransFunctions.h"


//// Used namespaces. ///////////////////////////////////////////////////////

using namespace WFG::Toolkit;
using namespace WFG::Toolkit::Examples;
using std::vector;
using std::string;


//// Local functions. ///////////////////////////////////////////////////////

namespace
{

//** Using a uniform random distribution, generate a number in [0,bound]. ***

double next_double( const double bound = 1.0 )
{
  assert( bound > 0.0 );

  return bound * rand() / static_cast< double >( RAND_MAX );
}


//** Create a random Pareto optimal solution for WFG1. **********************

vector< double > WFG_1_random_soln( const int k, const int l )
{
  vector< double > result;  // the result vector


  //---- Generate a random set of position parameters.

  for( int i = 0; i < k; i++ )
  {
    // Account for polynomial bias.
    result.push_back( pow( next_double(), 50.0 ) );
  }


  //---- Set the distance parameters.

  for( int i = k; i < k+l; i++ )
  {
    result.push_back( 0.35 );
  }


  //---- Scale to the correct domains.

  for( int i = 0; i < k+l; i++ )
  {
    result[i] *= 2.0*(i+1);
  }


  //---- Done.

  return result;
}


//** Create a random Pareto optimal solution for WFG2-WFG7. *****************

vector< double > WFG_2_thru_7_random_soln( const int k, const int l )
{
  vector< double > result;  // the result vector


  //---- Generate a random set of position parameters.

  for( int i = 0; i < k; i++ )
  {
    result.push_back( next_double() );
  }


  //---- Set the distance parameters.

  for( int i = k; i < k+l; i++ )
  {
    result.push_back( 0.35 );
  }


  //---- Scale to the correct domains.

  for( int i = 0; i < k+l; i++ )
  {
    result[i] *= 2.0*(i+1);
  }


  //---- Done.

  return result;
}


//** Create a random Pareto optimal solution for WFG8. **********************

vector< double > WFG_8_random_soln( const int k, const int l )
{
  vector< double > result;  // the result vector


  //---- Generate a random set of position parameters.

  for( int i = 0; i < k; i++ )
  {
    result.push_back( next_double() );
  }


  //---- Calculate the distance parameters.

  for( int i = k; i < k+l; i++ )
  {
    const vector< double >  w( result.size(), 1.0 );
    const double u = TransFunctions::r_sum( result, w  );

    const double tmp1 = fabs( floor( 0.5 - u ) + 0.98/49.98 );
    const double tmp2 = 0.02 + 49.98*( 0.98/49.98 - ( 1.0 - 2.0*u )*tmp1 );

    result.push_back( pow( 0.35, pow( tmp2, -1.0 ) ));
  }


  //---- Scale to the correct domains.

  for( int i = 0; i < k+l; i++ )
  {
    result[i] *= 2.0*(i+1);
  }


  //---- Done.

  return result;
}


//** Create a random Pareto optimal solution for WFG9. **********************

vector< double > WFG_9_random_soln( const int k, const int l )
{
  vector< double > result( k+l );  // the result vector


  //---- Generate a random set of position parameters.

  for( int i = 0; i < k; i++ )
  {
    result[i] = next_double();
  }


  //---- Calculate the distance parameters.

  result[k+l-1] = 0.35;  // the last distance parameter is easy

  for( int i = k+l-2; i >= k; i-- )
  {
    vector< double > result_sub;
    for( int j = i+1; j < k+l; j++ )
    {
      result_sub.push_back( result[j] );
    }

    const vector< double > w( result_sub.size(), 1.0 );
    const double tmp1 = TransFunctions::r_sum( result_sub, w  );

    result[i] = pow( 0.35, pow( 0.02 + 1.96*tmp1, -1.0 ) );
  }


  //---- Scale to the correct domains.

  for( int i = 0; i < k+l; i++ )
  {
    result[i] *= 2.0*(i+1);
  }


  //---- Done.

  return result;
}


//** Create a random Pareto optimal solution for I1. *****************

vector< double > I1_random_soln( const int k, const int l )
{
  vector< double > result;  // the result vector


  //---- Generate a random set of position parameters.

  for( int i = 0; i < k; i++ )
  {
    result.push_back( next_double() );
  }


  //---- Set the distance parameters.

  for( int i = k; i < k+l; i++ )
  {
    result.push_back( 0.35 );
  }


  //---- Done.

  return result;
}


//** Create a random Pareto optimal solution for I2. **********************

vector< double > I2_random_soln( const int k, const int l )
{
  vector< double > result( k+l );  // the result vector


  //---- Generate a random set of position parameters.

  for( int i = 0; i < k; i++ )
  {
    result[i] = next_double();
  }


  //---- Calculate the distance parameters.

  result[k+l-1] = 0.35;  // the last distance parameter is easy

  for( int i = k+l-2; i >= k; i-- )
  {
    vector< double > result_sub;
    for( int j = i+1; j < k+l; j++ )
    {
      result_sub.push_back( result[j] );
    }

    const vector< double > w( result_sub.size(), 1.0 );
    const double tmp1 = TransFunctions::r_sum( result_sub, w  );

    result[i] = pow( 0.35, pow( 0.02 + 1.96*tmp1, -1.0 ) );
  }


  //---- Done.

  return result;
}


//** Create a random Pareto optimal solution for I3. **********************

vector< double > I3_random_soln( const int k, const int l )
{
  vector< double > result;  // the result vector


  //---- Generate a random set of position parameters.

  for( int i = 0; i < k; i++ )
  {
    result.push_back( next_double() );
  }


  //---- Calculate the distance parameters.

  for( int i = k; i < k+l; i++ )
  {
    const vector< double >  w( result.size(), 1.0 );
    const double u = TransFunctions::r_sum( result, w  );

    const double tmp1 = fabs( floor( 0.5 - u ) + 0.98/49.98 );
    const double tmp2 = 0.02 + 49.98*( 0.98/49.98 - ( 1.0 - 2.0*u )*tmp1 );

    result.push_back( pow( 0.35, pow( tmp2, -1.0 ) ));
  }


  //---- Done.

  return result;
}


//** Create a random Pareto optimal solution for I4. **********************

vector< double > I4_random_soln( const int k, const int l )
{
  return I1_random_soln( k, l );
}


//** Create a random Pareto optimal solution for I5. **********************

vector< double > I5_random_soln( const int k, const int l )
{
  return I3_random_soln( k, l );
}


//** Generate a random solution for a given problem. ************************

vector< double > problem_random_soln
(
  const int k,
  const int l,
  const std::string fn
)
{
  if ( fn == "WFG1" )
  {
    return WFG_1_random_soln( k, l );
  }
  else if
  (
    fn == "WFG2" ||
    fn == "WFG3" ||
    fn == "WFG4" ||
    fn == "WFG5" ||
    fn == "WFG6" ||
    fn == "WFG7"
  )
  {
    return WFG_2_thru_7_random_soln( k, l );
  }
  else if ( fn == "WFG8" )
  {
    return WFG_8_random_soln( k, l );
  }
  else if ( fn == "WFG9" )
  {
    return WFG_9_random_soln( k, l );
  }
  else if ( fn == "I1" )
  {
    return I1_random_soln( k, l );
  }
  else if ( fn == "I2" )
  {
    return I2_random_soln( k, l );
  }
  else if ( fn == "I3" )
  {
    return I3_random_soln( k, l );
  }
  else if ( fn == "I4" )
  {
    return I4_random_soln( k, l );
  }
  else if ( fn == "I5" )
  {
    return I5_random_soln( k, l );
  }
  else
  {
    assert( false );
    return vector< double >();
  }
}


//** Calculate the fitness for a problem given some parameter set. **********

vector< double > problem_calc_fitness
(
  const vector< double >& z,
  const int k,
  const int M,
  const std::string fn
)
{
  if ( fn == "WFG1" )
  {
    return Problems::WFG1( z, k, M );
  }
  else if ( fn == "WFG2" )
  {
    return Problems::WFG2( z, k, M );
  }
  else if ( fn == "WFG3" )
  {
    return Problems::WFG3( z, k, M );
  }
  else if ( fn == "WFG4" )
  {
    return Problems::WFG4( z, k, M );
  }
  else if ( fn == "WFG5" )
  {
    return Problems::WFG5( z, k, M );
  }
  else if ( fn == "WFG6" )
  {
    return Problems::WFG6( z, k, M );
  }
  else if ( fn == "WFG7" )
  {
    return Problems::WFG7( z, k, M );
  }
  else if ( fn == "WFG8" )
  {
    return Problems::WFG8( z, k, M );
  }
  else if ( fn == "WFG9" )
  {
    return Problems::WFG9( z, k, M );
  }
  else if ( fn == "I1" )
  {
    return Problems::I1( z, k, M );
  }
  else if ( fn == "I2" )
  {
    return Problems::I2( z, k, M );
  }
  else if ( fn == "I3" )
  {
    return Problems::I3( z, k, M );
  }
  else if ( fn == "I4" )
  {
    return Problems::I4( z, k, M );
  }
  else if ( fn == "I5" )
  {
    return Problems::I5( z, k, M );
  }
  else
  {
    assert( false );
    return vector< double >();
  }
}


//** Convert a double vector into a string. *********************************

string make_string( const vector< double >& v )
{
  std::ostringstream result;

  if( !v.empty() )
  {
    result << v.front();
  }

  for( int i = 1; i < static_cast< int >( v.size() ); i++ )
  {
    result << " " << v[i];
  }

  return result.str();
}

}  // unnamed namespace


//// Standard functions. ////////////////////////////////////////////////////

//** Main. ******************************************************************

int WfgFunctions(std::string fn)
{

  //---- Generate values for desired function.

  if
  (
    fn == "WFG1" ||
    fn == "WFG2" ||
    fn == "WFG3" ||
    fn == "WFG4" ||
    fn == "WFG5" ||
    fn == "WFG6" ||
    fn == "WFG7" ||
    fn == "WFG8" ||
    fn == "WFG9" ||
    fn == "I1"   ||
    fn == "I2"   ||
    fn == "I3"   ||
    fn == "I4"   ||
    fn == "I5"
  )
  {
    const int count = 10000;  // the number of random fitness values to print
    const int M = 3;          // the number of objectives
    const int k_factor = 2;   // k (# position parameters) = k_factor*( M-1 ) 
    const int l_factor = 2;   // l (# distance parameters) = l_factor*2

    srand( 0 );  // seed the random number generator

    // Generate count random fitness values.
    for( int i = 0; i <= count; i++ )
    {
      const int k = k_factor*( M-1 );
      const int l = l_factor*2;

      const vector< double >& z = problem_random_soln( k, l, fn );
      const vector< double >& f = problem_calc_fitness( z, k, M, fn );

      std::cout << make_string( f ) << std::endl;
    }
  }
  else if
  (
    fn == "b_poly" ||
    fn == "b_flat" ||
    fn == "s_linear" ||
    fn == "s_decept" ||
    fn == "s_multi"
  )
  {
    const int count = 10000;  // the number of times (-1) to sample the function

    // Sample the transformation function count+1 times.
    for( int i = 0; i <= count; i++ )
    {
      const double y = static_cast< double >( i ) / count;
      double new_y;

      if ( fn == "b_poly" )
      {
        new_y = TransFunctions::b_poly( y, 20.0 );
      }
      else if ( fn == "b_flat" )
      {
        new_y = TransFunctions::b_flat( y, 0.7, 0.4, 0.5 );
      }
      else if ( fn == "s_linear" )
      {
        new_y = TransFunctions::s_linear( y, 0.35 );
      }
      else if ( fn == "s_decept" )
      {
        new_y = TransFunctions::s_decept( y, 0.35, 0.005, 0.05 );
      }
      else if ( fn == "s_multi" )
      {
        new_y = TransFunctions::s_multi( y, 5, 10, 0.35 );
      }
      else
      {
        assert( false );
        return 1;
      }

      std::cout << y << " " << new_y << std::endl;
    }
  }
  else if( fn == "b_param" || fn == "r_sum" || fn == "r_nonsep" )
  {
    srand( 0 );

    const int count = 10000;

    // Randomly sample the transformation count times.
    for( int i = 0; i < count; i++ )
    {
      vector< double > y;

      y.push_back( next_double() );
      y.push_back( next_double() );
      double new_y;

      if ( fn == "b_param" )
      {
        new_y = TransFunctions::b_param( y[0], y[1], 0.5, 2, 10 );
      }
      else if ( fn == "r_sum" )
      {
        vector< double > w;

        w.push_back( 1.0 );
        w.push_back( 5.0 );

        new_y = TransFunctions::r_sum( y, w );
      }
      else if ( fn == "r_nonsep" )
      {
        new_y = TransFunctions::r_nonsep( y, 2 );
      }
      else
      {
        assert( false );
        return 1;
      }

      std::cout << y[0] << " " << y[1] << " " << new_y << std::endl;
    }
  }
  else
  {
    std::cout << "Invalid fn.\n";
    return 1;
  }

  return 0;
}

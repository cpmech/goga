#include <iostream>
#include <cstdlib>
#include <string>
#include <sstream>
#include <vector>
#include <cassert>
#include <cmath>

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

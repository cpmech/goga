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
 * Defines as functions all of the transition vectors described in the EMO
 * 2005 paper, namely the transition vectors employed by WFG1--WFG9 and the
 * transition vectors employed by I1--I5. For details of the effects and
 * requirements of each transition vector, refer to the EMO 2005 paper
 * (available from the WFG web site).
 */


#ifndef EXAMPLE_TRANSITIONS_H
#define EXAMPLE_TRANSITIONS_H


//// Standard includes. /////////////////////////////////////////////////////

#include <vector>


//// Definitions/namespaces. ////////////////////////////////////////////////

namespace WFG
{

namespace Toolkit
{

namespace Examples
{

namespace Transitions
{

/*
 * For any transition function, the first "k" elements of "y" are
 * position-related parameters, and "M" is the number of objectives.
 */

//** t1 from WFG1. **********************************************************
std::vector< double > WFG1_t1( const std::vector< double >& y, const int k );

//** t2 from WFG1. **********************************************************
std::vector< double > WFG1_t2( const std::vector< double >& y, const int k );

//** t3 from WFG1. **********************************************************
std::vector< double > WFG1_t3( const std::vector< double >& y );

//** t4 from WFG1. **********************************************************
std::vector< double > WFG1_t4
(
  const std::vector< double >& y,
  const int k,
  const int M
);

//** t2 from WFG2. **********************************************************
std::vector< double > WFG2_t2( const std::vector< double >& y, const int k );

//** t3 from WFG2. Effectively as per WFG4, t2. *****************************
std::vector< double > WFG2_t3
(
  const std::vector< double >& y,
  const int k,
  const int M
);

//** t1 from WFG4. **********************************************************
std::vector< double > WFG4_t1( const std::vector< double >& y );

//** t1 from WFG5. **********************************************************
std::vector< double > WFG5_t1( const std::vector< double >& y );

//** t2 from WFG6. **********************************************************
std::vector< double > WFG6_t2
(
  const std::vector< double >& y,
  const int k,
  const int M
);

//** t1 from WFG7. **********************************************************
std::vector< double > WFG7_t1( const std::vector< double >& y, const int k );

//** t1 from WFG8. **********************************************************
std::vector< double > WFG8_t1( const std::vector< double >& y, const int k );

//** t1 from WFG9. **********************************************************
std::vector< double > WFG9_t1( const std::vector< double >& y );

//** t2 from WFG9. **********************************************************
std::vector< double > WFG9_t2( const std::vector< double >& y, const int k );

//** t2 from I1. ************************************************************
std::vector< double > I1_t2( const std::vector< double >& y, const int k );

//** t3 from I1. ************************************************************
std::vector< double > I1_t3
(
  const std::vector< double >& y,
  const int k,
  const int M
);

//** t1 from I2. ************************************************************
std::vector< double > I2_t1( const std::vector< double >& y );

//** t1 from I3. ************************************************************
std::vector< double > I3_t1( const std::vector< double >& y );

//** t3 from I4. ************************************************************
std::vector< double > I4_t3
(
  const std::vector< double >& y,
  const int k,
  const int M
);

}  // Transitions namespace

}  // Examples namespace

}  // Toolkit namespace

}  // WFG namespace

#endif
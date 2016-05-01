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
 * Defines the problems described in the EMO 2005 paper, namely WFG1--WFG9
 * and I1--I5. For the specifics of each problem, refer to the EMO 2005 paper
 * (available from the WFG web site).
 */


#ifndef EXAMPLE_PROBLEMS_H
#define EXAMPLE_PROBLEMS_H


//// Standard includes. /////////////////////////////////////////////////////

#include <vector>


//// Definitions/namespaces. ////////////////////////////////////////////////

namespace WFG
{

namespace Toolkit
{

namespace Examples
{

namespace Problems
{

/*
 * For all problems, the first "k" elements of "z" are the position-related
 * parameters, and "M" is the number of objectives.
 */

//** The WFG1 problem. ******************************************************
std::vector< double > WFG1
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The WFG2 problem. ******************************************************
std::vector< double > WFG2
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The WFG3 problem. ******************************************************
std::vector< double > WFG3
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The WFG4 problem. ******************************************************
std::vector< double > WFG4
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The WFG5 problem. ******************************************************
std::vector< double > WFG5
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The WFG6 problem. ******************************************************
std::vector< double > WFG6
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The WFG7 problem. ******************************************************
std::vector< double > WFG7
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The WFG8 problem. ******************************************************
std::vector< double > WFG8
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The WFG9 problem. ******************************************************
std::vector< double > WFG9
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The I1 problem. ********************************************************
std::vector< double > I1
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The I2 problem. ********************************************************
std::vector< double > I2
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The I3 problem. ********************************************************
std::vector< double > I3
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The I4 problem. ********************************************************
std::vector< double > I4
(
  const std::vector< double >& z,
  const int k,
  const int M
);

//** The I5 problem. ********************************************************
std::vector< double > I5
(
  const std::vector< double >& z,
  const int k,
  const int M
);

}  // Problems namespace

}  // Examples namespace

}  // Toolkit namespace

}  // WFG namespace

#endif
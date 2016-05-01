/* File : wfg.i */
%module wfg

%include <typemaps.i>
%include "std_string.i"
%include "std_vector.i"

namespace std {
   %template(DoubleVector) vector<double>;
}

%{
#include "ExampleProblems.h"
%}

/* Let's just grab the original header file here */
%include "ExampleProblems.h"

%inline %{
extern int WfgFunctions(std::string fn);
%}

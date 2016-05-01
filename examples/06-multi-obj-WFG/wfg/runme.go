// This example illustrates how C++ classes can be used from Go using SWIG.

package main

import (
	"fmt"

	. "./wfg"
)

func main() {
	ret := WfgFunctions("WFG1")
	fmt.Printf("ret = %v\n", ret)
}

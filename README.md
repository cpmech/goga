# Goga &ndash; Go Evolutionary/Genetic Algorithm

Goga is a computer library for developing evolutionary algorithms and, in particular, genetic
algorithms. The goal of these algorithms is to solve optimisation problems with (or not) many
constraints and many objectives. Also, problems with mixed-type representations such as those with
real numbers and integers are considered by Goga.



## Examples

[Check out examples here](https://github.com/cpmech/goga/blob/master/examples/README.md)



## Installation

1 Install dependencies:

Goga depends on the [Gosl Go Scientific Library](https://github.com/cpmech/gosl), therefore, please
install Gosl first.

2 Install Goga:

```
go get github.com/cpmech/goga
```


## Documentation

Here, we call _user-defined_ types as _structures_. These are simply Go `types` defined as `struct`.
Some may think of these _structures_ as _classes_. Goga has several _global_ functions as well and
tries to avoid complicated (unnecessary?) programming techniques.

An allocated structure is called here an **object** and functions attached to this object are called
**methods**. The variable holding the pointer to an object is always named **o** in Goga (e.g.
like `self` or `this`).

Some objects need to be initialised before usage. In this case, functions named `Init` have to be
called (e.g. like `constructors`). Some structures require an explicit call to a function to release
allocated memory. This function is named `Free`. Functions that allocate a pointer to a structure
are prefixed with `New`; for instance: `NewIsoSurf`.

Goga has several functions and _structures_. Check the **[developer's
documentation](http://rawgit.com/cpmech/goga/master/doc/index.html)** to see what's available and
how to call functions and methods.




## Authors and license

See the AUTHORS file.

Unless otherwise noted, the Goga source files are distributed under the BSD-style license found in the LICENSE file.

# Goga &ndash; Go Evolutionary/Genetic Algorithm

Goga is a computer library for developing evolutionary algorithms based on the _differential
evolution_ and/or _genetic algorithm_ concepts. The goal of these algorithms is to solve
optimisation problems with (or not) many constraints and many objectives. Also, problems with
mixed-type representations such as those with real numbers and integers are considered by Goga.

[See the documentation](https://godoc.org/github.com/cpmech/goga) for more details (e.g. how to call
functions and use structures).

[![GoDoc](https://godoc.org/github.com/cpmech/goga?status.svg)](https://godoc.org/github.com/cpmech/goga)

The core algorithms in Goga are well explained in [this journal paper](doc/goga.pdf); see also [1, 2]


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
Some may think of these _structures_ as _classes_. Goga has several global functions as well and
tries to avoid complicated constructions.

An allocated structure is called here an **object** and functions attached to this object are called
**methods**. The variable holding the pointer to an object is always named **o** in Goga (e.g.
like `self` or `this`).

Some objects need to be initialised before usage. In this case, functions named `Init` have to be
called (e.g. like `constructors`).



## Bibliography

Goga is included in the following works:

1. Pedroso DM, Bonyadi MR, Gallagher M (2017) Parallel evolutionary algorithm for single and multi-objective optimisation: differential evolution and constraints handling, Applied Soft Computing http://dx.doi.org/10.1016/j.asoc.2017.09.006
2. Pedroso DM (2017) FORM reliability analysis using a parallel evolutionary algorithm, Structural Safety 65:84-99 http://dx.doi.org/10.1016/j.strusafe.2017.01.001


## Authors and license

See the AUTHORS file.

Unless otherwise noted, the Goga source files are distributed under the BSD-style license found in the LICENSE file.

OR-Tools for Go
---

[![Go Reference](https://pkg.go.dev/badge/github.com/irfansharif/solver.svg)](https://godocs.io/github.com/irfansharif/solver)

This project contains a cgo-based API for using Google's [Operations Research
Tools](https://developers.google.com/optimization/). It exposes a high-level
package for the [CP-SAT
Solver](https://developers.google.com/optimization/cp/cp_solver), targeting the
[v8.0](https://github.com/google/or-tools/releases/tag/v8.0) release.

### Examples

Here's a simple example solving for free integer variables, ensuring that
they're all different.

```go
model := NewModel()

var numVals int64 = 3
x := model.NewIntVar(0, numVals-1, "x")
y := model.NewIntVar(0, numVals-1, "y")
z := model.NewIntVar(0, numVals-1, "z")

ct := NewAllDifferentConstraint(x, y, z)
model.AddConstraints(ct)

result := model.Solve()
require.True(t, result.Optimal(), "expected solver to find solution")

{
  x := result.Value(x)
  y := result.Value(y)
  z := result.Value(z)

  for _, value := range []int64{x, y, z} {
    require.Truef(t, value >= 0 && value <= numVals-1,
      "expected %d to be in domain [%d, %d]", value, 0, numVals-1)
  }

  require.Falsef(t, x == y || x == z || y == z,
    "all different constraint violated, both x=%d y=%d z=%d", x, y, z)
}
```

Here's another solving with a few linear constraints and a maximization
objective.

```go
model := NewModel()
x := model.NewIntVar(0, 100, "x")
y := model.NewIntVar(0, 100, "y")

// Constraint 1: x + 2y <= 14.
ct1 := NewLinearConstraint(
  NewLinearExpr([]IntVar{x, y}, []int64{1, 2}, 0),
  NewDomain(math.MinInt64, 14),
)

// Constraint 2: 3x - y >= 0.
ct2 := NewLinearConstraint(
  NewLinearExpr([]IntVar{x, y}, []int64{3, -1}, 0),
  NewDomain(0, math.MaxInt64),
)

// Constraint 3: x - y <= 2.
ct3 := NewLinearConstraint(
  NewLinearExpr([]IntVar{x, y}, []int64{1, -1}, 0),
  NewDomain(0, 2),
)

model.AddConstraints(ct1, ct2, ct3)

// Objective function: 3x + 4y.
model.Maximize(NewLinearExpr([]IntVar{x, y}, []int64{3, 4}, 0))

result := model.Solve()
require.True(t, result.Optimal(), "expected solver to find solution")

{
  x := result.Value(x)
  y := result.Value(y)

  require.Equal(t, int64(6), x)
  require.Equal(t, int64(4), y)
  require.Equal(t, float64(34), result.ObjectiveValue())
}
```

Finally, an example solving for arbitrary boolean constraints.

```go
model := NewModel()

a := model.NewLiteral("a")
b := model.NewLiteral("b")
c := model.NewLiteral("c")
d := model.NewLiteral("d")
e := model.NewLiteral("e")
f := model.NewLiteral("f")

model.AddConstraints(
  NewBooleanAndConstraint(a, b), // a && b
  NewBooleanOrConstraint(c, d),  // c || d
  NewBooleanXorConstraint(e, f), // e != f
)

result := model.Solve()
require.True(t, result.Optimal(), "expected solver to find solution")

{
  a := result.BooleanValue(a)
  b := result.BooleanValue(b)
  c := result.BooleanValue(c)
  d := result.BooleanValue(d)
  e := result.BooleanValue(e)
  f := result.BooleanValue(f)

  require.True(t, a && b)
  require.True(t, c || d)
  require.True(t, e != f)
}
```

For more, look through the package tests and the
[docs](https://godocs.io/github.com/irfansharif/solver).

### Contributing

The Go/C++ binding code is generated using [SWIG](http://www.swig.org) and can
be found under `internal/`. SWIG generated code is ugly and difficult to work
with; a sanitized API is exposed via the top-level package.

Because of the C++ dependencies, the library is compiled/tested using
[Bazel](https://bazel.build). The top-level Makefiles packages most things
you'd need.

```sh
# ensure that the submodules are initialized:
#   git submodule update --init --recursive
#
# supported bazel version >= 4.0.0
# supported swig version == 4.0.2
# supported protoc version == 3.14.0
# supported protoc-gen-go version == 1.27.1

$ make help
Supported commands: build, test, generate, rewrite

$ make generate
--- generating go:generate files
--- generating swig files
--- generating proto files
--- generating bazel files
ok

$ make build
ok

$ make test
...
INFO: Build completed successfully, 4 total actions
```

#### Testing

This library is tested using the (awesome)
[datadriven](https://github.com/cockroachdb/datadriven) library + a tiny
testing grammar. See `testdata/` for what that looks like.

```
sat
model.name(ex)
model.literals(x, y, z)
constrain.at-most-k(x to z | 2)
model.print()
----
model=ex
  literals (num = 3)
    x, y, z
  constraints (num = 1)
    at-most-k: x, y, z | 2

sat
model.solve()
----
optimal

sat
result.bools(x to z)
----
x = false
y = false
z = false
```

```sh
# to update the testdata files
$ make rewrite

# to run specific tests
$ bazel test :all internal/... --test_output=all \
  --cache_test_results=no \
  --test_arg='-test.v' \
  --test_filter='Test.*'
```

### Acknowledgements

The SWIG interface files to work with protobufs was cribbed from
[AirspaceTechnologies/or-tools](https://github.com/AirspaceTechnologies/or-tools).
To figure out how to structure this package as a stand-alone bazel target, I
looked towards from
[gonzojive/or-tools-go](https://github.com/gonzojive/or-tools-go). The CP-SAT
stuff was then mostly pattern matching.

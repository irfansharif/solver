OR-Tools for Go
---

[![Go Reference](https://pkg.go.dev/badge/github.com/irfansharif/or-tools.svg)](https://godocs.io/github.com/irfansharif/or-tools/cpsatsolver)


This project contains a cgo-based API for using Google's [Operations Research
Tools](https://developers.google.com/optimization/). It exposes high-level
packages for the [CP-SAT
Solver](https://developers.google.com/optimization/cp/cp_solver) and the [Glop
Linear Solver](https://developers.google.com/optimization/lp/glop).

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

and := NewBooleanAndConstraint(a, b) // a && b
or := NewBooleanOrConstraint(c, d)   // c || d
xor := NewBooleanXorConstraint(e, f) // e != f
model.AddConstraints(and, or, xor)

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

For more, look through the package tests.

### Contributing

The Go/C++ binding code is generated using [SWIG](http://www.swig.org), and can
be found under `internal/swig`. SWIG generated code is ugly and difficult to
work with, which is why the package is internal; a sanitized API is exposed via
the top-level package.

Because of the C++ dependencies, the library is compiled/tested using
[Bazel](https://bazel.build).

```sh
# supported bazel version >= 4.0.0
bazel build cpsatsolver/... linearsolver...
bazel test cpsatsolver:all --test_output=all \
  --cache_test_results=no \
  --test_arg='-test.v' \
  --test_filter='Test.*'
bazel test linearsolver/...

# to update the BUILD.bazel files
bazel run //:gazelle
bazel run //:gazelle -- update-repos -from_file=go.mod -prune=true
```

The generated files are checked in. They can be regenerated using the
following:

```sh
# ensure that the submodules are initialized:
#   git submodule update --init --recursive
#
# supported swig version == 4.0.2
# supported protoc version == 3.13.0
# supported protoc-gen-go version == 1.27.1

# to generate the C++/Go wrapper files
swig -v -go -cgo -c++ -intgosize 64 \
  -Ic-deps/or-tools \
  -Ic-deps/abseil-cpp \
  -o internal/linearsolver/linearsolver_wrapper.cc \
  -module linearsolver \
  internal/linearsolver/linearsolver.i

swig -v -go -cgo -c++ -intgosize 64 \
  -Ic-deps/or-tools \
  -Ic-deps/abseil-cpp \
  -o internal/cpsatsolver/cpsatsolver_wrapper.cc \
  -module cpsatsolver \
  internal/cpsatsolver/cpsatsolver.i

# to generate the protobuf files
protoc --proto_path=internal/cpsatsolver/pb \
  --go_out=internal/cpsatsolver/pb \
  --go_opt=Mcp_model.proto=github.com/irfansharif/or-tools/internal/cpsatsolver/pb \
  --go_opt=Msat_parameters.proto=github.com/irfansharif/or-tools/internal/cpsatsolver/pb \
  --go_opt=paths=source_relative \
  cp_model.proto sat_parameters.proto
```

---

NB: This repo was originally forked from
[gonzojive/or-tools-go](https://github.com/gonzojive/or-tools-go), which
exposed the linear solver. Bits of it were cribbed from
[AirspaceTechnologies/or-tools](https://github.com/AirspaceTechnologies/or-tools);
they authored the SWIG interface files to generate code dealing with protobufs.
The CP-SAT stuff was then mostly pattern matching.

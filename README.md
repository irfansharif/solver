OR-Tools for Go
---

This project contains a cgo-based API for using Google's [Operations Research
Tools](https://developers.google.com/optimization/). The Go/C++ binding code is
generated using [SWIG](http://www.swig.org), and can be found under
`internal/swig`. SWIG generated code is ugly and difficult to work with, which
is why the package is internal; a sanitized API is exposed via the
top-level package.

Due to the C++ dependencies, the library is compiled/tested using
[Bazel](https://bazel.build). 

```sh
# supported bazel version >= 4.0.0
bazel build //:or-tools
bazel test //:all --test_output=all \
  --cache_test_results=no \
  --test_arg='-test.v' \
  --test_filter='TestNew.*'
bazel run //:gazelle # to update the BUILD.bazel files
```

### Regenerating the SWIG bindings

The generated files are checked in. They can be regenerated using the
following:

```sh
# ensure that the submodules are initialized:
#   git submodule update --init --recursive
#
# supported swig version == 4.0.2
swig -v -go -cgo -c++ -intgosize 64 \
  -Ic-deps/or-tools \
  -Ic-deps/abseil-cpp \
  -o internal/swig/linear_solver_go_wrap.cc \
  -module swig \
  internal/swig/linear_solver.i
```

---

NB: This repo was originally forked from
[gonzojive/or-tools-go](https://github.com/gonzojive/or-tools-go).

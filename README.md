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
bazel test :all internal/...

# to update the BUILD.bazel files
bazel run //:gazelle
```

### Regenerating the SWIG bindings

The generated files are checked in. They can be regenerated using the
following:

```sh
# ensure that the submodules are initialized:
#   git submodule update --init --recursive
#
# supported swig version == 4.0.2
# supported protoc version == 3.13.0
# supported protoc-gen-go version == 1.27.1

# to generate to C++/Go wrapper files
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

```sh
# to run `gofmt` against everything
gofmt -s -w .
```

---

NB: This repo was originally forked from
[gonzojive/or-tools-go](https://github.com/gonzojive/or-tools-go). Bits of it
were cribbed from
[AirspaceTechnologies/or-tools](https://github.com/AirspaceTechnologies/or-tools).

## OR tools for golang

This project contains a cgo-based API for using Google's Operations Research
tools. Code is generated in the `ortoolsswig` folder, but the generated code is
ugly, so most people will want to use the `ortoolsgo` package, which is written
on top of the swig bindings.

The library compiles with bazel. For example, the tests can be run with this
command:

```shell
ibazel test //ortoolsswig:go_default_test
```

### Regenerating the SWIG bindings

For now, some generated files are checked in. They can be regenerated using this
command:

```shell
swig -v -go -cgo -c++ -intgosize 64 \
  -I/home/red/git/or-tools \
  -I/home/red/git/abseil-cpp \
  -o linear_solver_go_wrap.cc \
  -module ortoolsswig \
  linear_solver.i
```

It will be necessary to clone absl and or-tools:

https://github.com/abseil/abseil-cpp.git

https://github.com/google/or-tools.git

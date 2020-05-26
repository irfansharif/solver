```shell
swig -go -cgo -c++ -intgosize 64 -o linear_solver_go_wrap.cc -module ortools linear_solver.i
```


```shell
bazel build @ortools//examples/cpp:linear_programming
```

```
 //%pythoncode {
 //def setup_variable_operator(opname):
 //  setattr(Variable, opname,
 //          lambda self, *args: getattr(VariableExpr(self), opname)(*args))
 //for opname in LinearExpr.OVERRIDDEN_OPERATOR_METHODS:
 //  setup_variable_operator(opname)
 //}  // %pythoncode
```

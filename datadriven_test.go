// Copyright 2021 Irfan Sharif.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package solver_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/cockroachdb/datadriven"
	"github.com/irfansharif/solver"
	"github.com/irfansharif/solver/internal/testutils"
	"github.com/irfansharif/solver/internal/testutils/bazel"
	"github.com/irfansharif/solver/internal/testutils/parser/ast"
	"github.com/stretchr/testify/require"
)

func TestDatadriven(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		path, implant := bazel.WritableSandboxPathFor(t, "", path)
		defer implant()

		model := solver.NewModel("") // instantiate a model

		itvM := make(map[string]solver.Interval)
		varM := make(map[string]solver.IntVar)
		litM := make(map[string]solver.Literal)

		var result solver.Result
		var solved bool

		getIntervals := func(s *testutils.Scanner, is ...string) []solver.Interval {
			var intervals []solver.Interval
			for _, v := range is {
				iv, ok := itvM[v]
				if !ok {
					s.Fatalf("unrecognized variable: %s", v)
				}

				intervals = append(intervals, iv)
			}
			return intervals
		}

		getIntVars := func(s *testutils.Scanner, vs ...string) []solver.IntVar {
			var intVars []solver.IntVar
			for _, v := range vs {
				iv, ok := varM[v]
				if !ok {
					s.Fatalf("unrecognized variable: %s", v)
				}

				intVars = append(intVars, iv)
			}
			return intVars
		}

		getLiterals := func(s *testutils.Scanner, ls ...string) []solver.Literal {
			var literals []solver.Literal
			for _, l := range ls {
				lit, ok := litM[l]
				if !ok {
					s.Fatalf("unrecognized literal: %s", l)
				}

				literals = append(literals, lit)
			}
			return literals
		}

		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			parts := strings.Split(d.Pos, ":")
			line, _ := strconv.Atoi(parts[1])
			s := testutils.NewScanner(t, strings.NewReader(d.Input), path, line)
			var out strings.Builder
			for s.Scan() {
				stmt := testutils.Compile(s, s.Text())
				if d.Cmd == "recognize" {
					continue
				}

				switch stmt.Method {
				case ast.NameMethod: // model.name(arg)
					argument := stmt.Argument.(*ast.VariablesArgument)
					model.TestingSetName(argument.Variables[0])
				case ast.VarsMethod: // m.vars(x,y,z in [0, 2])
					argument := stmt.Argument.(*ast.DomainArgument)
					dom := argument.AsSolverDomain()
					for _, v := range argument.Variables {
						varM[v] = model.NewIntVarFromDomain(dom, v)
					}
				case ast.LiteralsMethod: // model.literals(c,d)
					argument := stmt.Argument.(*ast.VariablesArgument)
					for _, l := range argument.Variables {
						litM[l] = model.NewLiteral(l)
					}

				case ast.ConstantsMethod: // model.constants(a,b == 42)
					argument := stmt.Argument.(*ast.ConstantsArgument)
					for _, c := range argument.Variables {
						varM[c] = model.NewConstant(int64(argument.Constant), c)
					}
				case ast.IntervalsMethod: // model.intervals(i as [s,e|sz], j as [e,s|sz]) if a
					var enforcement []solver.Literal
					if stmt.Enforcement != nil {
						enforcement = getLiterals(s, stmt.Enforcement.Literals...)
					}

					argument := stmt.Argument.(*ast.IntervalsArgument)
					for _, iv := range argument.Intervals {
						variables := getIntVars(s, iv.Start, iv.End, iv.Size)
						start, end, size := variables[0], variables[1], variables[2]
						itvM[iv.Name] = model.NewInterval(start, end, size, iv.Name)
						itvM[iv.Name].OnlyEnforceIf(enforcement...)
					}
				case ast.PrintMethod: // model.print()
					out.WriteString(model.String())
				case ast.ValidateMethod: // m.validate()
					ok, err := model.Validate()
					if ok {
						out.WriteString("ok")
					} else {
						out.WriteString(fmt.Sprintf("invalid: %v", err.Error()))
					}
				case ast.SolveMethod: // model.solve()
					result = model.Solve()
					switch {
					case result.Feasible():
						out.WriteString("feasible")
					case result.Infeasible():
						out.WriteString("infeasible")
					case result.Invalid():
						out.WriteString("invalid")
					case result.Optimal():
						out.WriteString("optimal")
						solved = true
					}

				case ast.AllDifferentMethod: // constrain.all-different(x,y,z)
					argument := stmt.Argument.(*ast.VariablesArgument)
					intVars := getIntVars(s, argument.Variables...)
					model.AddConstraints(solver.NewAllDifferentConstraint(intVars...))
				case ast.AllSameMethod: // constrain.all-same(x,y,z)
					argument := stmt.Argument.(*ast.VariablesArgument)
					intVars := getIntVars(s, argument.Variables...)
					model.AddConstraints(solver.NewAllSameConstraint(intVars...))
				case ast.ImplicationMethod: // constrain.boolean-and(x,y,z) [if a,b]
					argument := stmt.Argument.(*ast.ImplicationArgument)
					literals := getLiterals(s, argument.Left, argument.Right)
					model.AddConstraints(solver.NewImplicationConstraint(literals[0], literals[1]))
				case ast.BooleanAndMethod: // constrain.boolean-and(x,y,z) [if a,b]
					argument := stmt.Argument.(*ast.VariablesArgument)
					literals := getLiterals(s, argument.Variables...)
					var enforcement []solver.Literal
					if stmt.Enforcement != nil {
						enforcement = getLiterals(s, stmt.Enforcement.Literals...)
					}
					model.AddConstraints(solver.NewBooleanAndConstraint(literals...).OnlyEnforceIf(enforcement...))
				case ast.BooleanOrMethod: // constrain.boolean-or(x,y,z) [if a,b]
					argument := stmt.Argument.(*ast.VariablesArgument)
					literals := getLiterals(s, argument.Variables...)
					var enforcement []solver.Literal
					if stmt.Enforcement != nil {
						enforcement = getLiterals(s, stmt.Enforcement.Literals...)
					}
					model.AddConstraints(solver.NewBooleanOrConstraint(literals...).OnlyEnforceIf(enforcement...))
				case ast.BooleanXorMethod: // constrain.boolean-xor(x,y,z)
					argument := stmt.Argument.(*ast.VariablesArgument)
					literals := getLiterals(s, argument.Variables...)
					model.AddConstraints(solver.NewBooleanXorConstraint(literals...))
				case ast.AtMostKMethod: // constrain.at-most-k(x to z | K)
					argument := stmt.Argument.(*ast.KArgument)
					literals := getLiterals(s, argument.Literals...)
					model.AddConstraints(solver.NewAtMostKConstraint(argument.K, literals...))
				case ast.AtLeastKMethod: // constrain.at-least-k(x to z | K)
					argument := stmt.Argument.(*ast.KArgument)
					literals := getLiterals(s, argument.Literals...)
					model.AddConstraints(solver.NewAtLeastKConstraint(argument.K, literals...))
				case ast.ExactlyKMethod: // constrain.exactly-k(x to z | K)
					argument := stmt.Argument.(*ast.KArgument)
					literals := getLiterals(s, argument.Literals...)
					model.AddConstraints(solver.NewExactlyKConstraint(argument.K, literals...))
				case ast.AssignmentsMethod:
					argument := stmt.Argument.(*ast.AssignmentsArgument)
					if argument.ForLiterals() {
						literals := getLiterals(s, argument.Variables...)
						if argument.In {
							model.AddConstraints(solver.NewAllowedLiteralAssignmentsConstraint(literals, argument.AllowedLiteralAssignments))
						} else {
							model.AddConstraints(solver.NewForbiddenLiteralAssignmentsConstraint(literals, argument.AllowedLiteralAssignments))
						}
					} else {
						variables := getIntVars(s, argument.Variables...)
						if argument.In {
							model.AddConstraints(solver.NewAllowedAssignmentsConstraint(variables, argument.AsInt64s()))
						} else {
							model.AddConstraints(solver.NewForbiddenAssignmentsConstraint(variables, argument.AsInt64s()))
						}
					}
				case ast.CumulativeMethod: // constrain.cumulative(i: 12, j: 13 | 32)
					argument := stmt.Argument.(*ast.CumulativeArgument)
					capacity := getIntVars(s, argument.Capacity)[0]
					model.AddConstraints(
						solver.NewCumulativeConstraint(
							capacity,
							getIntervals(s, argument.Intervals()...),
							argument.Demands(),
						),
					)
				case ast.BinaryOpMethod: // constrain.binary-op(a % b == c)
					argument := stmt.Argument.(*ast.BinaryOpArgument)
					switch argument.Op {
					case "%":
						variables := getIntVars(s, argument.Left, argument.Right, argument.Target)
						dividend, divisor, target := variables[0], variables[1], variables[2]
						model.AddConstraints(
							solver.NewModuloConstraint(target, dividend, divisor),
						)
					case "/":
						variables := getIntVars(s, argument.Left, argument.Right, argument.Target)
						dividend, divisor, target := variables[0], variables[1], variables[2]
						model.AddConstraints(
							solver.NewDivisionConstraint(target, dividend, divisor),
						)
					case "*":
						variables := getIntVars(s, argument.Target, argument.Left, argument.Right)
						target, multiplands := variables[0], variables[1:]
						model.AddConstraints(
							solver.NewProductConstraint(target, multiplands...),
						)
					}
				case ast.BoolsMethod: // result.bool(x,y to z)
					require.True(t, solved)
					argument := stmt.Argument.(*ast.VariablesArgument)
					literals := getLiterals(s, argument.Variables...)
					for i, lit := range literals {
						val := result.BooleanValue(lit)
						out.WriteString(fmt.Sprintf("%s = %t", argument.Variables[i], val))
						if i != len(literals)-1 {
							out.WriteString("\n")
						}
					}
				case ast.ValuesMethod: // result.values(x,y to z)
					require.True(t, solved)
					argument := stmt.Argument.(*ast.VariablesArgument)
					variables := getIntVars(s, argument.Variables...)
					for i, iv := range variables {
						val := result.Value(iv)
						out.WriteString(fmt.Sprintf("%s = %d", argument.Variables[i], val))
						if i != len(variables)-1 {
							out.WriteString("\n")
						}
					}
				default:
					t.Fatalf("unrecognized method: %s", stmt.Method)
				}
			}

			return out.String()
		})
	})
}

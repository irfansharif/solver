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
	"strings"
	"testing"

	"github.com/cockroachdb/datadriven"
	"github.com/irfansharif/solver"
	"github.com/irfansharif/solver/internal/testutils"
	"github.com/irfansharif/solver/internal/testutils/bazel"
	"github.com/irfansharif/solver/internal/testutils/parser/ast"
)

func TestDatadriven(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		path, implant := bazel.WritableSandboxPathFor(t, "", path)
		defer implant()

		model := solver.NewModel("") // instantiate a model
		varMap := make(map[string]solver.IntVar)
		litMap := make(map[string]solver.Literal)
		var result solver.Result

		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			s := testutils.NewScanner(t, strings.NewReader(d.Input), path)
			var out strings.Builder
			for s.Scan() {
				stmt, err := testutils.Compile(s.Text())
				if err != nil {
					s.Fatal(err)
				}
				if d.Cmd == "recognize" {
					continue
				}

				switch stmt.Method {
				case ast.NameMethod: // model.name(arg)
					argument := stmt.Argument.(*ast.VariablesArgument)
					model.SetName(argument.Variables[0])
				case ast.VarsMethod: // m.vars(x,y,z in [0, 2])
					argument := stmt.Argument.(*ast.DomainArgument)
					dom := argument.AsSolverDomain()
					for _, v := range argument.Variables {
						varMap[v] = model.NewIntVarFromDomain(dom, v)
					}
				case ast.LiteralsMethod: // model.literals(c,d)
					argument := stmt.Argument.(*ast.VariablesArgument)
					for _, l := range argument.Variables {
						litMap[l] = model.NewLiteral(l)
					}
				case ast.ConstantsMethod: // model.constants(a,b == 42)
					argument := stmt.Argument.(*ast.ConstantsArgument)
					for _, v := range argument.Variables {
						iv := model.NewConstant(int64(argument.Constant), v)
						varMap[v] = iv
					}
				case ast.PrintMethod: // model.print()
					out.WriteString(model.String())
				case ast.ValidateMethod: // m.validate()
					var ok bool
					ok, err = model.Validate()
					if ok {
						out.WriteString("ok")
					} else {
						out.WriteString(fmt.Sprintf("validation error: %v", err.Error()))
					}
				case ast.SolveMethod: // model.solve()
					result = model.Solve()
				case ast.OptimalMethod: // result.optimal()
					out.WriteString(fmt.Sprintf("%t", result.Optimal()))
				case ast.InfeasibleMethod: // result.infeasible()
					out.WriteString(fmt.Sprintf("%t", result.Infeasible()))
				case ast.FeasibleMethod: // result.feasible()
					out.WriteString(fmt.Sprintf("%t", result.Feasible()))
				case ast.InvalidMethod: // result.invalid()
					out.WriteString(fmt.Sprintf("%t", result.Invalid()))
				case ast.BooleansMethod: // result.booleans(x,y to z)
					argument := stmt.Argument.(*ast.VariablesArgument)
					for i, l := range argument.Variables {
						lit, ok := litMap[l]
						if !ok {
							t.Fatalf("unrecognized literal: %s", l)
						}

						val := result.BooleanValue(lit)
						out.WriteString(fmt.Sprintf("%s = %t", l, val))
						if i != len(argument.Variables)-1 {
							out.WriteString("\n")
						}
					}
				case ast.ValuesMethod: // result.values(x,y to z)
					argument := stmt.Argument.(*ast.VariablesArgument)
					for i, v := range argument.Variables {
						iv, ok := varMap[v]
						if !ok {
							t.Fatalf("unrecognized variable: %s", v)
						}

						val := result.Value(iv)
						out.WriteString(fmt.Sprintf("%s = %d", v, val))
						if i != len(argument.Variables)-1 {
							out.WriteString("\n")
						}
					}

				case ast.AllDifferentMethod: // constrain.all-different(x,y,z)
					argument := stmt.Argument.(*ast.VariablesArgument)
					var intVars []solver.IntVar
					for _, v := range argument.Variables {
						iv, ok := varMap[v]
						if !ok {
							t.Fatalf("unrecognized variable: %s", v)
						}

						intVars = append(intVars, iv)
					}

					model.AddConstraints(
						solver.NewAllDifferentConstraint(intVars...),
					)
				case ast.AllSameMethod: // constrain.all-same(x,y,z)
					argument := stmt.Argument.(*ast.VariablesArgument)
					var intVars []solver.IntVar
					for _, v := range argument.Variables {
						iv, ok := varMap[v]
						if !ok {
							t.Fatalf("unrecognized variable: %s", v)
						}

						intVars = append(intVars, iv)
					}

					model.AddConstraints(
						solver.NewAllSameConstraint(intVars...),
					)
				case ast.BooleanAndMethod: // constrain.boolean -and(x,y,z) [if a,b]
					argument := stmt.Argument.(*ast.VariablesArgument)
					var literals []solver.Literal
					for _, l := range argument.Variables {
						lit, ok := litMap[l]
						if !ok {
							t.Fatalf("unrecognized literal: %s", l)
						}

						literals = append(literals, lit)
					}

					var enforcement []solver.Literal
					if stmt.Enforcement != nil {
						for _, l := range stmt.Enforcement.Variables {
							lit, ok := litMap[l]
							if !ok {
								t.Fatalf("unrecognized enforcement literal: %s", l)
							}

							enforcement = append(enforcement, lit)
						}
					}

					model.AddConstraints(
						solver.NewBooleanAndConstraint(literals...).OnlyEnforceIf(enforcement...),
					)
				case ast.AtMostKMethod: // constrain.at-most-k(x to z | K)
					argument := stmt.Argument.(*ast.KArgument)
					var literals []solver.Literal
					for _, l := range argument.Variables {
						lit, ok := litMap[l]
						if !ok {
							t.Fatalf("unrecognized literal: %s", l)
						}

						literals = append(literals, lit)
					}

					model.AddConstraints(
						solver.NewAtMostKConstraint(argument.K, literals...),
					)
				default:
					t.Fatalf("unrecognized method: %s", stmt.Method)
				}

			}

			return out.String()
		})
	})
}

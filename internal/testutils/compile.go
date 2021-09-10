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

package testutils

import (
	"testing"

	"github.com/irfansharif/solver/internal/testutils/parser"
	"github.com/irfansharif/solver/internal/testutils/parser/ast"
)

// Compile compiles the given statement and returns the corresponding AST node.
func Compile(tb testing.TB, input string) *ast.Statement {
	p := parser.New(tb, input)
	stmt := p.Statement()

	// TODO(irfansharif): Should we make a single receiver+method type? There
	// are only three receivers, and a static list of methods.
	switch stmt.Receiver {
	case "model":
		switch stmt.Method {
		case ast.ConstantsMethod, ast.IntervalsMethod, ast.LiteralsMethod,
			ast.MaximizeMethod, ast.MinimizeMethod, ast.NameMethod,
			ast.PrintMethod, ast.SolveMethod, ast.SolveAllMethod,
			ast.ValidateMethod, ast.VarsMethod:
		default:
			tb.Fatalf("unrecognized method: %s.%s", stmt.Receiver, stmt.Method)
		}
	case "constrain":
		switch stmt.Method {
		case ast.AllDifferentMethod, ast.AllSameMethod, ast.AssignmentsMethod,
			ast.AtLeastKMethod, ast.AtMostKMethod, ast.BinaryOpMethod,
			ast.BooleanAndMethod, ast.BooleanOrMethod, ast.BooleanXorMethod,
			ast.CumulativeMethod, ast.ElementMethod, ast.EqualityMethod,
			ast.ExactlyKMethod, ast.ImplicationMethod, ast.LinearExprsMethod,
			ast.NonOverlappingMethod, ast.NonOverlapping2DMethod:
		default:
			tb.Fatalf("unrecognized method: %s.%s", stmt.Receiver, stmt.Method)
		}
	case "result":
		switch stmt.Method {
		case ast.BoolsMethod, ast.ObjectiveValueMethod, ast.ValuesMethod:
		default:
			tb.Fatalf("unrecognized method: %s.%s", stmt.Receiver, stmt.Method)
		}
	default:
		tb.Fatalf("unrecognized receiver: %s", stmt.Receiver)
	}

	if stmt.Enforcement != nil {
		switch stmt.Method {
		case ast.BooleanOrMethod, ast.BooleanAndMethod, ast.LinearExprsMethod:
		case ast.IntervalsMethod:
			if len(stmt.Enforcement.Variables) > 1 {
				tb.Fatalf("only single enforcement literal supported for %s.%s", stmt.Receiver, stmt.Method)
			}
		default:
			tb.Fatalf("enforcement clause unsupported for %s.%s", stmt.Receiver, stmt.Method)
		}
	}

	if stmt.Argument != nil {
		switch t := stmt.Argument.(type) {
		case *ast.AssignmentsArgument:

			switch stmt.Method {
			case ast.AssignmentsMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.BinaryOpArgument:
			switch stmt.Method {
			case ast.BinaryOpMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.ConstantsArgument:
			switch stmt.Method {
			case ast.ConstantsMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.CumulativeArgument:
			switch stmt.Method {
			case ast.CumulativeMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.DomainArgument:
			switch stmt.Method {
			case ast.VarsMethod, ast.LinearExprsMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.ElementArgument:
			switch stmt.Method {
			case ast.ElementMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.ImplicationArgument:
			switch stmt.Method {
			case ast.ImplicationMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.IntervalsArgument:
			switch stmt.Method {
			case ast.IntervalsMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.KArgument:
			switch stmt.Method {
			case ast.AtMostKMethod, ast.AtLeastKMethod, ast.ExactlyKMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.LinearEqualityArgument:
			switch stmt.Method {
			case ast.EqualityMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.LinearExprsArgument:
			switch stmt.Method {
			case ast.MaximizeMethod, ast.MinimizeMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.NonOverlapping2DArgument:
			switch stmt.Method {
			case ast.NonOverlapping2DMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.VariableEqualityArgument:
			switch stmt.Method {
			case ast.EqualityMethod:
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		case *ast.VariablesArgument:
			switch stmt.Method {
			case ast.AllDifferentMethod, ast.AllSameMethod,
				ast.BooleanAndMethod, ast.BooleanOrMethod, ast.BooleanXorMethod,
				ast.BoolsMethod, ast.LiteralsMethod, ast.NameMethod,
				ast.NonOverlappingMethod, ast.ValuesMethod:
			case ast.MaximizeMethod, ast.MinimizeMethod:
				// There's ambiguity in the grammar, and we give precedence to
				// VariablesArgument during parsing. Let's fix up here.
				stmt.Argument = t.AsLinearExprsArgument()
			default:
				tb.Fatalf("unexpected type for %s.%s: %T", stmt.Receiver, stmt.Method, t)
			}
		default:
			tb.Fatalf("unrecognized type: %T", t)
		}
	}

	return stmt
}

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
	"fmt"

	"github.com/irfansharif/solver/internal/testutils/parser"
	"github.com/irfansharif/solver/internal/testutils/parser/ast"
)

// Compile compiles the given statement and returns the corresponding ast node.
func Compile(stmt string) (*ast.Statement, error) {
	p := parser.New(stmt)
	s, err := p.Statement()
	if err != nil {
		return nil, err
	}

	// TODO(irfansharif): Should we make a single receiver+method type? There
	// are only three receivers, and a static list of methods.
	switch s.Receiver {
	case "model":
		switch s.Method {
		case ast.ConstantsMethod, ast.IntervalsMethod, ast.LiteralsMethod,
			ast.MaximizeMethod, ast.MinimizeMethod, ast.NameMethod,
			ast.PrintMethod, ast.SolveMethod, ast.SolveAllMethod,
			ast.ValidateMethod, ast.VarsMethod:
		default:
			return nil, fmt.Errorf("unrecognized method: %s.%s", s.Receiver, s.Method)
		}
	case "constrain":
		switch s.Method {
		case ast.AllDifferentMethod, ast.AllSameMethod, ast.AssignmentsMethod,
			ast.AtLeastKMethod, ast.AtMostKMethod, ast.BinaryOpMethod,
			ast.BooleanAndMethod, ast.BooleanOrMethod, ast.BooleanXorMethod,
			ast.CumulativeMethod, ast.ElementMethod, ast.EqualityMethod,
			ast.ExactlyKMethod, ast.ImplicationMethod, ast.LinearExprsMethod,
			ast.NonOverlappingMethod, ast.NonOverlapping2DMethod:
		default:
			return nil, fmt.Errorf("unrecognized method: %s.%s", s.Receiver, s.Method)
		}
	case "result":
		switch s.Method {
		case ast.BooleansMethod, ast.FeasibleMethod, ast.InfeasibleMethod,
			ast.InvalidMethod, ast.ObjectiveValueMethod, ast.OptimalMethod,
			ast.ValuesMethod:
		default:
			return nil, fmt.Errorf("unrecognized method: %s.%s", s.Receiver, s.Method)
		}
	default:
		return nil, fmt.Errorf("unrecognized reciever: %s", s.Receiver)
	}

	if s.Enforcement != nil {
		switch s.Method {
		case ast.BooleanOrMethod, ast.BooleanAndMethod, ast.LinearExprsMethod:
		case ast.IntervalsMethod:
			if len(s.Enforcement.Variables) > 1 {
				return nil, fmt.Errorf("only single enforcement literal supported for %s.%s", s.Receiver, s.Method)
			}
		default:
			return nil, fmt.Errorf("enforcement clause unsupported for %s.%s", s.Receiver, s.Method)
		}
	}

	if s.Argument != nil {
		switch t := s.Argument.(type) {
		case *ast.AssignmentsArgument:
			switch s.Method {
			case ast.AssignmentsMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.BinaryOpArgument:
			switch s.Method {
			case ast.BinaryOpMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.ConstantsArgument:
			switch s.Method {
			case ast.ConstantsMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.CumulativeArgument:
			switch s.Method {
			case ast.CumulativeMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.DomainArgument:
			switch s.Method {
			case ast.VarsMethod, ast.LinearExprsMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.ElementArgument:
			switch s.Method {
			case ast.ElementMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.ImplicationArgument:
			switch s.Method {
			case ast.ImplicationMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.IntervalsArgument:
			switch s.Method {
			case ast.IntervalsMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.KArgument:
			switch s.Method {
			case ast.AtMostKMethod, ast.AtLeastKMethod, ast.ExactlyKMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.LinearEqualityArgument:
			switch s.Method {
			case ast.EqualityMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.LinearExprsArgument:
			switch s.Method {
			case ast.MaximizeMethod, ast.MinimizeMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.NonOverlapping2DArgument:
			switch s.Method {
			case ast.NonOverlapping2DMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.VariableEqualityArgument:
			switch s.Method {
			case ast.EqualityMethod:
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		case *ast.VariablesArgument:
			switch s.Method {
			case ast.AllDifferentMethod, ast.AllSameMethod,
				ast.BooleanAndMethod, ast.BooleanOrMethod, ast.BooleanXorMethod,
				ast.BooleansMethod, ast.LiteralsMethod, ast.NameMethod,
				ast.NonOverlappingMethod, ast.ValuesMethod:
			case ast.MaximizeMethod, ast.MinimizeMethod:
				// There's ambiguity in the grammar, and we give precedence to
				// VariablesArgument during parsing. Let's fix up here.
				s.Argument = t.AsLinearExprsArgument()
			default:
				return nil, fmt.Errorf("unexpected type for %s.%s: %T", s.Receiver, s.Method, t)
			}
		default:
			return nil, fmt.Errorf("unrecognized type: %T", t)
		}
	}

	return s, nil
}

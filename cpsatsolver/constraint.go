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

package cpsatsolver

import (
	"fmt"

	swigpb "github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

// Constraint is what a model attempts to satisfy when deciding on a solution.
type Constraint = *constraint

type constraint struct {
	proto *swigpb.ConstraintProto
}

// NewAllDifferentConstraint forces all variables to take different values.
func NewAllDifferentConstraint(vars ...IntVar) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_AllDiff{
				AllDiff: &swigpb.AllDifferentConstraintProto{
					Vars: intVars(vars).indexes(),
				},
			},
		},
	}
}

// NewAtMostOneConstraint enforces that no more than one literal is
// true at the same time.
func NewAtMostOneConstraint(literals ...Literal) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_AtMostOne{
				AtMostOne: &swigpb.BoolArgumentProto{
					Literals: intVars(literals).indexes(),
				},
			},
		},
	}
}

// NewBooleanAndConstraint forces all the literals to be true.
func NewBooleanAndConstraint(literals ...Literal) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_BoolAnd{
				BoolAnd: &swigpb.BoolArgumentProto{
					Literals: intVars(literals).indexes(),
				},
			},
		},
	}
}

// NewBooleanOrConstraint forces at least one literal to be true.
func NewBooleanOrConstraint(literals ...Literal) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_BoolOr{
				BoolOr: &swigpb.BoolArgumentProto{
					Literals: intVars(literals).indexes(),
				},
			},
		},
	}
}

// NewBooleanXorConstraint forces an odd number of the literals to be true.
func NewBooleanXorConstraint(literals ...Literal) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_BoolXor{
				BoolXor: &swigpb.BoolArgumentProto{
					Literals: intVars(literals).indexes(),
				},
			},
		},
	}
}

// NewImplicationConstraint ensures that the first literal implies the second.
func NewImplicationConstraint(a, b Literal) Constraint {
	return NewBooleanOrConstraint(a.negation(fmt.Sprintf("~%s", a.name())), b)
}

// NewAllowedLiteralAssignmentsConstraint forces the values of the n-tuple
// formed by the given literals to be among one of the listed n-tuple
// assignments.
func NewAllowedLiteralAssignmentsConstraint(literals []Literal, assignments [][]bool) Constraint {
	return newLiteralAssignmentsConstraintInternal(literals, assignments)
}

// NewForbiddenLiteralAssignmentsConstraint forbids the values of the n-tuple
// formed by the given literals to be among one of the listed n-tuple
// assignments.
func NewForbiddenLiteralAssignmentsConstraint(literals []Literal, assignments [][]bool) Constraint {
	constraint := newLiteralAssignmentsConstraintInternal(literals, assignments)
	constraint.proto.GetTable().Negated = true
	return constraint
}

// NewDivisionConstraint forces the target to equal numerator/denominator.
func NewDivisionConstraint(target, numerator, denominator IntVar) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntDiv{
				IntDiv: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars([]IntVar{numerator, denominator}).indexes(),
				},
			},
		},
	}
}

// NewProductConstraint forces the target to equal the product of all
// multiplicands.
func NewProductConstraint(target IntVar, multiplicands ...IntVar) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntProd{
				IntProd: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars(multiplicands).indexes(),
				},
			},
		},
	}
}

// NewMaximumConstraint forces the target to equal the maximum of all
// variables.
func NewMaximumConstraint(target IntVar, vars ...IntVar) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntMax{
				IntMax: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars(vars).indexes(),
				},
			},
		},
	}
}

// NewMinimumConstraint forces the target to equal the minimum of all
// variables.
func NewMinimumConstraint(target IntVar, vars ...IntVar) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntMin{
				IntMin: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars(vars).indexes(),
				},
			},
		},
	}
}

// NewModuloConstraint forces the target to equal dividend%divisor.
func NewModuloConstraint(target, dividend, divisor IntVar) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntMod{
				IntMod: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars([]IntVar{dividend, divisor}).indexes(),
				},
			},
		},
	}
}

// NewAllowedAssignmentsConstraint forces the values of the n-tuple
// formed by the given variables to be among one of the listed n-tuple
// assignments.
func NewAllowedAssignmentsConstraint(vars []IntVar, assignments [][]int64) Constraint {
	return newAssignmentsConstraintInternal(vars, assignments)
}

// NewForbiddenAssignmentsConstraint forbids the values of the n-tuple
// formed by the given variables to be among one of the listed n-tuple
// assignments.
func NewForbiddenAssignmentsConstraint(vars []IntVar, assignments [][]int64) Constraint {
	constraint := newAssignmentsConstraintInternal(vars, assignments)
	constraint.proto.GetTable().Negated = true
	return constraint
}

// NewLinearConstraint enforces a linear inequality among the variables,
// such as 0 <= x + 2y <= 10.
func NewLinearConstraint(e LinearExpr, d Domain) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_Linear{
				Linear: &swigpb.LinearConstraintProto{
					Vars:   e.vars(),
					Coeffs: e.coeffs(),
					Domain: d.list(e.offset()),
				},
			},
		},
	}
}

// NewLinearMaximumConstraint forces the target to equal the maximum of all
// linear expressions.
func NewLinearMaximumConstraint(target LinearExpr, exprs ...LinearExpr) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_LinMax{
				LinMax: &swigpb.LinearArgumentProto{
					Target: target.proto,
					Exprs:  linearExprs(exprs).protos(),
				},
			},
		},
	}
}

// NewLinearMinimumConstraint forces the target to equal the minimum of all
// linear expressions.
func NewLinearMinimumConstraint(target LinearExpr, exprs ...LinearExpr) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_LinMin{
				LinMin: &swigpb.LinearArgumentProto{
					Target: target.proto,
					Exprs:  linearExprs(exprs).protos(),
				},
			},
		},
	}
}

// NewElementConstraint forces the target to equal vars[index]. Implicitly,
// index takes on one of the values in [0, len(vars)).
func NewElementConstraint(target, index IntVar, vars ...IntVar) Constraint {
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_Element{
				Element: &swigpb.ElementConstraintProto{
					Target: target.index(),
					Index:  index.index(),
					Vars:   intVars(vars).indexes(),
				},
			},
		},
	}
}

// NewElementLiteralConstraint forces the target to equal literals[index].
// Implicitly, index takes on one of the values in [0, len(literals)).
func NewElementLiteralConstraint(target Literal, index IntVar, literals ...Literal) Constraint {
	return NewElementConstraint(target, index, literals...)
}

// OnlyEnforceIf enforces the constraint iff all literals listed are true. If
// not explicitly called, of if the list is empty, then the constraint will
// always be enforced.
//
// NB: Only a few constraint support enforcement:
// - NewBooleanOrConstraint
// - NewBooleanAndConstraint
// - NewLinearConstraint
func (c *constraint) OnlyEnforceIf(literals ...Literal) {
	c.proto.EnforcementLiteral = intVars(literals).indexes()
}

func newLiteralAssignmentsConstraintInternal(literals []Literal, assignments [][]bool) Constraint {
	var integerAssignments [][]int64
	for _, assignment := range assignments { // convert [][]bool to [][]int64
		var integerAssignment []int64
		for _, a := range assignment {
			i := 0
			if a {
				i = 1
			}
			integerAssignment = append(integerAssignment, int64(i))
		}
		integerAssignments = append(integerAssignments, integerAssignment)
	}

	return newAssignmentsConstraintInternal(literals, integerAssignments)
}

func newAssignmentsConstraintInternal(vars []IntVar, assignments [][]int64) Constraint {
	var values []int64
	for _, assignment := range assignments {
		if len(assignment) != len(vars) {
			panic("mismatched assignment and vars length")
		}
		values = append(values, assignment...)
	}
	return &constraint{
		proto: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_Table{
				Table: &swigpb.TableConstraintProto{
					Vars:   intVars(vars).indexes(),
					Values: values,
				},
			},
		},
	}
}

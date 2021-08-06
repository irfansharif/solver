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
	swigpb "github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

// Constraint is what a model attempts to satisfy when deciding on a solution.
type Constraint interface {
	// OnlyEnforceIf enforces the constraint iff all literals listed are true. If
	// not explicitly called, of if the list is empty, then the constraint will
	// always be enforced.
	//
	// NB: Only a few constraints support enforcement:
	// - NewBooleanOrConstraint
	// - NewBooleanAndConstraint
	// - NewLinearConstraint
	//
	// Intervals support enforcement too, but only with a single literal.
	OnlyEnforceIf(literals ...Literal) Constraint

	// protos returns the underlying CP-SAT constraint protobuf representations.
	protos() []*swigpb.ConstraintProto
}

type constraint struct {
	pb *swigpb.ConstraintProto
}

var _ Constraint = &constraint{}

// NewAllDifferentConstraint forces all variables to take different values.
func NewAllDifferentConstraint(vars ...IntVar) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_AllDiff{
				AllDiff: &swigpb.AllDifferentConstraintProto{
					Vars: intVars(vars).indexes(),
				},
			},
		},
	}
}

// NewAllSameConstraint forces all variables to take the same values.
func NewAllSameConstraint(vars ...IntVar) Constraint {
	var cs []Constraint
	for i := range vars {
		if i == 0 {
			continue
		}
		cs = append(cs, NewMaximumConstraint(vars[i-1], vars[i]))
	}
	return constraints(cs)
}

// NewAtMostKConstraint ensures that no more than k literals are true.
func NewAtMostKConstraint(k int, literals ...Literal) Constraint {
	if k == 1 {
		return newAtMostOneConstraint(literals...)
	}

	lb, ub := int64(0), int64(k)
	return NewLinearConstraint(Sum(asIntVars(literals)...), NewDomain(lb, ub))
}

// NewAtLeastKConstraint ensures that at least k literals are true.
func NewAtLeastKConstraint(k int, literals ...Literal) Constraint {
	if k == 1 {
		return NewBooleanOrConstraint(literals...)
	}

	lb, ub := int64(k), int64(len(literals))
	return NewLinearConstraint(Sum(asIntVars(literals)...), NewDomain(lb, ub))
}

// NewExactlyKConstraint ensures that exactly k literals are true.
func NewExactlyKConstraint(k int, literals ...Literal) Constraint {
	lb, ub := int64(k), int64(k)
	return NewLinearConstraint(Sum(asIntVars(literals)...), NewDomain(lb, ub))
}

// NewBooleanAndConstraint ensures that all literals are true.
func NewBooleanAndConstraint(literals ...Literal) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_BoolAnd{
				BoolAnd: &swigpb.BoolArgumentProto{
					Literals: asIntVars(literals).indexes(),
				},
			},
		},
	}
}

// NewBooleanOrConstraint ensures that at least one literal is true. It can be
// thought of as a special case of NewAtLeastKConstraint, but one that uses a
// more efficient internal encoding.
func NewBooleanOrConstraint(literals ...Literal) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_BoolOr{
				BoolOr: &swigpb.BoolArgumentProto{
					Literals: asIntVars(literals).indexes(),
				},
			},
		},
	}
}

// NewBooleanXorConstraint ensures that an odd number of the literals are true.
func NewBooleanXorConstraint(literals ...Literal) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_BoolXor{
				BoolXor: &swigpb.BoolArgumentProto{
					Literals: asIntVars(literals).indexes(),
				},
			},
		},
	}
}

// NewImplicationConstraint ensures that the first literal implies the second.
func NewImplicationConstraint(a, b Literal) Constraint {
	return NewBooleanOrConstraint(a.Not(), b)
}

// NewAllowedLiteralAssignmentsConstraint ensures that the values of the n-tuple
// formed by the given literals is one of the listed n-tuple assignments.
func NewAllowedLiteralAssignmentsConstraint(literals []Literal, assignments [][]bool) Constraint {
	return newLiteralAssignmentsConstraintInternal(literals, assignments)
}

// NewForbiddenLiteralAssignmentsConstraint ensures that the values of the
// n-tuple formed by the given literals is not one of the listed n-tuple
// assignments.
func NewForbiddenLiteralAssignmentsConstraint(literals []Literal, assignments [][]bool) Constraint {
	constraint := newLiteralAssignmentsConstraintInternal(literals, assignments)
	constraint.pb.GetTable().Negated = true
	return constraint
}

// NewDivisionConstraint ensures that the target is to equal to
// numerator/denominator.
func NewDivisionConstraint(target, numerator, denominator IntVar) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntDiv{
				IntDiv: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars([]IntVar{numerator, denominator}).indexes(),
				},
			},
		},
	}
}

// NewProductConstraint ensures that the target to equal to the product of all
// multiplicands.
func NewProductConstraint(target IntVar, multiplicands ...IntVar) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntProd{
				IntProd: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars(multiplicands).indexes(),
				},
			},
		},
	}
}

// NewMaximumConstraint ensures that the target is equal to the maximum of all
// variables.
func NewMaximumConstraint(target IntVar, vars ...IntVar) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntMax{
				IntMax: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars(vars).indexes(),
				},
			},
		},
	}
}

// NewMinimumConstraint ensures that the target is equal to the minimum of all
// variables.
func NewMinimumConstraint(target IntVar, vars ...IntVar) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntMin{
				IntMin: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars(vars).indexes(),
				},
			},
		},
	}
}

// NewModuloConstraint ensures that the target to equal to dividend%divisor.
func NewModuloConstraint(target, dividend, divisor IntVar) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_IntMod{
				IntMod: &swigpb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVars([]IntVar{dividend, divisor}).indexes(),
				},
			},
		},
	}
}

// NewAllowedAssignmentsConstraint ensures that the values of the n-tuple
// formed by the given variables is one of the listed n-tuple assignments.
func NewAllowedAssignmentsConstraint(vars []IntVar, assignments [][]int64) Constraint {
	return newAssignmentsConstraintInternal(vars, assignments)
}

// NewForbiddenAssignmentsConstraint ensures that the values of the n-tuple
// formed by the given variables is not one of the listed n-tuple assignments.
func NewForbiddenAssignmentsConstraint(vars []IntVar, assignments [][]int64) Constraint {
	constraint := newAssignmentsConstraintInternal(vars, assignments)
	constraint.pb.GetTable().Negated = true
	return constraint
}

// NewLinearConstraint ensures that the linear expression lies in the given
// domain. It can be used to express linear equalities of the form:
//
// 		0 <= x + 2y <= 10
//
func NewLinearConstraint(e LinearExpr, d Domain) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
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

// NewLinearMaximumConstraint ensures that the target is equal to the maximum of
// all linear expressions.
func NewLinearMaximumConstraint(target LinearExpr, exprs ...LinearExpr) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_LinMax{
				LinMax: &swigpb.LinearArgumentProto{
					Target: target.proto(),
					Exprs:  linearExprs(exprs).protos(),
				},
			},
		},
	}
}

// NewLinearMinimumConstraint ensures that the target is equal to the minimum of
// all linear expressions.
func NewLinearMinimumConstraint(target LinearExpr, exprs ...LinearExpr) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_LinMin{
				LinMin: &swigpb.LinearArgumentProto{
					Target: target.proto(),
					Exprs:  linearExprs(exprs).protos(),
				},
			},
		},
	}
}

// NewElementConstraint ensures that the target is equal to vars[index].
// Implicitly index takes on one of the values in [0, len(vars)).
func NewElementConstraint(target, index IntVar, vars ...IntVar) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
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

// NewNonOverlappingConstraint ensures that all the intervals are disjoint.
// More formally, there must exist a sequence such that for every pair of
// consecutive intervals, we have intervals[i].end <= intervals[i+1].start.
// Intervals of size zero matter for this constraint. This is also known as a
// disjunctive constraint in scheduling.
func NewNonOverlappingConstraint(intervals ...Interval) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_NoOverlap{
				NoOverlap: &swigpb.NoOverlapConstraintProto{
					Intervals: itrvals(intervals).indexes(),
				},
			},
		},
	}
}

// NewNonOverlapping2DConstraint ensures that the boxes defined by the following
// don't overlap:
//
// 		[xintervals[i].start, xintervals[i].end)
// 		[yintervals[i].start, yintervals[i].end)
//
// Intervals/boxes of size zero are considered for overlap if the last argument
// is true.
func NewNonOverlapping2DConstraint(
	xintervals []Interval,
	yintervals []Interval,
	boxesWithNoAreaCanOverlap bool,
) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_NoOverlap_2D{
				NoOverlap_2D: &swigpb.NoOverlap2DConstraintProto{
					XIntervals: itrvals(xintervals).indexes(),
					YIntervals: itrvals(yintervals).indexes(),

					BoxesWithNullAreaCanOverlap: boxesWithNoAreaCanOverlap,
				},
			},
		},
	}
}

// NewCumulativeConstraint ensures that the sum of the demands of the intervals
// (intervals[i]'s demand is specified in demands[i]) at each interval point
// cannot exceed a max capacity. The intervals are interpreted as [start, end).
// Intervals of size zero are ignored.
func NewCumulativeConstraint(capacity int32, intervals []Interval, demands []int32) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_Cumulative{
				Cumulative: &swigpb.CumulativeConstraintProto{
					Capacity:  capacity,
					Intervals: itrvals(intervals).indexes(),
					Demands:   demands,
				},
			},
		},
	}
}

// newAtMostOneConstraint is a special case of NewAtMostKConstraint that uses a
// more efficient internal encoding.
func newAtMostOneConstraint(literals ...Literal) Constraint {
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_AtMostOne{
				AtMostOne: &swigpb.BoolArgumentProto{
					Literals: asIntVars(literals).indexes(),
				},
			},
		},
	}
}

// OnlyEnforceIf is part of the Constraint interface.
func (c *constraint) OnlyEnforceIf(literals ...Literal) Constraint {
	c.pb.EnforcementLiteral = asIntVars(literals).indexes()
	return c
}

// protos is part of the Constraint interface.
func (c *constraint) protos() []*swigpb.ConstraintProto {
	return []*swigpb.ConstraintProto{c.pb}
}

type constraints []Constraint

var _ Constraint = &constraints{}

// OnlyEnforceIf is part of the Constraint interface.
func (cs constraints) OnlyEnforceIf(literals ...Literal) Constraint {
	for _, c := range cs {
		c.OnlyEnforceIf(literals...)
	}
	return cs
}

// protos is part of the Constraint interface.
func (cs constraints) protos() []*swigpb.ConstraintProto {
	var res []*swigpb.ConstraintProto
	for _, c := range cs {
		res = append(res, c.protos()...)
	}
	return res
}

func newLiteralAssignmentsConstraintInternal(literals []Literal, assignments [][]bool) *constraint {
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

	return newAssignmentsConstraintInternal(asIntVars(literals), integerAssignments)
}

func newAssignmentsConstraintInternal(vars []IntVar, assignments [][]int64) *constraint {
	var values []int64
	for _, assignment := range assignments {
		if len(assignment) != len(vars) {
			panic("mismatched assignment and vars length")
		}
		values = append(values, assignment...)
	}
	return &constraint{
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_Table{
				Table: &swigpb.TableConstraintProto{
					Vars:   intVars(vars).indexes(),
					Values: values,
				},
			},
		},
	}
}

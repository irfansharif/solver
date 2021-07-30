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

// Package cpsatsolver is a Go wrapper library for the CP-SAT Solver included
// as part Google's Operations Research tools.
package cpsatsolver

import (
	"fmt"

	swigpb "github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

type Model struct {
	proto           *swigpb.CpModelProto
	intVarToIdx     map[*intVar]int
	constraintToIdx map[*constraint]int
}

func NewModel() *Model {
	return &Model{
		proto:           &swigpb.CpModelProto{},
		intVarToIdx:     make(map[*intVar]int),
		constraintToIdx: make(map[*constraint]int),
	}
}

type Literal = *intVar
type IntVar = *intVar
type Domain = *domain

func (m *Model) NewLiteral(name string) Literal {
	return m.NewIntVarFromDomain(NewDomain(0, 1), name)
}

func (m *Model) NewIntVar(lb int64, ub int64, name string) IntVar {
	return m.NewIntVarFromDomain(NewDomain(lb, ub), name)
}

func (m *Model) NewIntVarFromDomain(domain Domain, name string) IntVar {
	intVar := newIntVar(domain, name)
	m.addIntVar(intVar)
	return intVar
}

func (m *Model) NewConstant(c int64) IntVar {
	return m.NewIntVarFromDomain(NewDomain(c, c), fmt.Sprintf("%d", c))
}

func (m *Model) AddBooleanOr(ls ...Literal) {
	c := newConstraint()
	literals := m.getIntVarIndexes(ls...)
	c.proto.Constraint = &swigpb.ConstraintProto_BoolOr{
		BoolOr: &swigpb.BoolArgumentProto{
			Literals: literals,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddBooleanAnd(ls ...Literal) {
	c := newConstraint()
	literals := m.getIntVarIndexes(ls...)
	c.proto.Constraint = &swigpb.ConstraintProto_BoolAnd{
		BoolAnd: &swigpb.BoolArgumentProto{
			Literals: literals,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddBooleanXor(ls ...Literal) {
	c := newConstraint()
	literals := m.getIntVarIndexes(ls...)
	c.proto.Constraint = &swigpb.ConstraintProto_BoolXor{
		BoolXor: &swigpb.BoolArgumentProto{
			Literals: literals,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddAtMostOne(ls ...Literal) {
	c := newConstraint()
	literals := m.getIntVarIndexes(ls...)
	c.proto.Constraint = &swigpb.ConstraintProto_AtMostOne{
		AtMostOne: &swigpb.BoolArgumentProto{
			Literals: literals,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddAllowedLiteralAssignments(ls []Literal, assignments [][]bool) {
	_ = m.addLiteralAssignmentsInternal(ls, assignments)
	return
}

func (m *Model) AddForbiddenLiteralAssignments(ls []Literal, assignments [][]bool) {
	c := m.addLiteralAssignmentsInternal(ls, assignments)
	c.proto.Constraint.(*swigpb.ConstraintProto_Table).Table.Negated = true
	return
}

func (m *Model) AddElementLiteral(target Literal, index IntVar, ls ...Literal) {
	m.AddElement(target, index, ls...)
}

func (m *Model) AddAllDifferent(is ...IntVar) {
	c := newConstraint()
	vars := m.getIntVarIndexes(is...)
	c.proto.Constraint = &swigpb.ConstraintProto_AllDiff{
		AllDiff: &swigpb.AllDifferentConstraintProto{
			Vars: vars,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddAllowedAssignments(is []IntVar, assignments [][]int64) {
	_ = m.addAssignmentsInternal(is, assignments)
	return
}

func (m *Model) AddForbiddenAssignments(is []IntVar, assignments [][]int64) {
	c := m.addAssignmentsInternal(is, assignments)
	c.proto.Constraint.(*swigpb.ConstraintProto_Table).Table.Negated = true
	return
}

func (m *Model) AddElement(target, index IntVar, is ...IntVar) {
	c := newConstraint()
	vars := m.getIntVarIndexes(is...)
	c.proto.Constraint = &swigpb.ConstraintProto_Element{
		Element: &swigpb.ElementConstraintProto{
			Target: m.getIntVarIndex(target),
			Index:  m.getIntVarIndex(index),
			Vars:   vars,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddDivision(target, numerator, denominator IntVar) {
	c := newConstraint()
	c.proto.Constraint = &swigpb.ConstraintProto_IntDiv{
		IntDiv: &swigpb.IntegerArgumentProto{
			Target: m.getIntVarIndex(target),
			Vars:   m.getIntVarIndexes(numerator, denominator),
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddModulo(target, dividend, divisor IntVar) {
	c := newConstraint()
	c.proto.Constraint = &swigpb.ConstraintProto_IntMod{
		IntMod: &swigpb.IntegerArgumentProto{
			Target: m.getIntVarIndex(target),
			Vars:   m.getIntVarIndexes(dividend, divisor),
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddMaximum(target IntVar, is ...IntVar) {
	c := newConstraint()
	c.proto.Constraint = &swigpb.ConstraintProto_IntMax{
		IntMax: &swigpb.IntegerArgumentProto{
			Target: m.getIntVarIndex(target),
			Vars:   m.getIntVarIndexes(is...),
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddMinimum(target IntVar, is ...IntVar) {
	c := newConstraint()
	c.proto.Constraint = &swigpb.ConstraintProto_IntMin{
		IntMin: &swigpb.IntegerArgumentProto{
			Target: m.getIntVarIndex(target),
			Vars:   m.getIntVarIndexes(is...),
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddProduct(target IntVar, is ...IntVar) {
	c := newConstraint()
	c.proto.Constraint = &swigpb.ConstraintProto_IntProd{
		IntProd: &swigpb.IntegerArgumentProto{
			Target: m.getIntVarIndex(target),
			Vars:   m.getIntVarIndexes(is...),
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddLinearConstraint(expr LinearExpr, domain Domain) {
	c := newConstraint()
	c.proto.Constraint = &swigpb.ConstraintProto_Linear{
		Linear: &swigpb.LinearConstraintProto{
			Vars:   m.getIntVarIndexes(expr.vars...),
			Coeffs: expr.coeffs,
			Domain: domain.list(expr.offset),
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddLinearMaximum(target LinearExpr, ls ...LinearExpr) {
	c := newConstraint()
	c.proto.Constraint = &swigpb.ConstraintProto_LinMax{
		LinMax: &swigpb.LinearArgumentProto{
			Target: m.asLinearExprProto(target),
			Exprs:  m.asLinearExprProtos(ls...),
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddLinearMinimum(target LinearExpr, ls ...LinearExpr) {
	c := newConstraint()
	c.proto.Constraint = &swigpb.ConstraintProto_LinMax{
		LinMax: &swigpb.LinearArgumentProto{
			Target: m.asLinearExprProto(target),
			Exprs:  m.asLinearExprProtos(ls...),
		},
	}
	m.addConstraint(c)
}

func (m *Model) Minimize(expr LinearExpr) {
	m.proto.Objective = &swigpb.CpObjectiveProto{
		Vars:   m.getIntVarIndexes(expr.vars...),
		Coeffs: expr.coeffs,
		Offset: float64(expr.offset),
	}
}

func (m *Model) Maximize(expr LinearExpr) {
	m.Minimize(expr)
	for i, coeff := range m.proto.GetObjective().GetCoeffs() {
		m.proto.GetObjective().GetCoeffs()[i] = -coeff
	}
	m.proto.GetObjective().ScalingFactor = -1
	m.proto.GetObjective().Offset = -m.proto.GetObjective().GetOffset()
}

func (m *Model) addLiteralAssignmentsInternal(ls []Literal, assignments [][]bool) *constraint {
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

	return m.addAssignmentsInternal(ls, integerAssignments)
}

func (m *Model) addAssignmentsInternal(is []IntVar, assignments [][]int64) *constraint {
	c := newConstraint()
	vars := m.getIntVarIndexes(is...)
	var values []int64
	for _, assignment := range assignments {
		if len(assignment) != len(is) {
			panic("mismatched assignment and int vars length")
		}
		values = append(values, assignment...)
	}
	c.proto.Constraint = &swigpb.ConstraintProto_Table{
		Table: &swigpb.TableConstraintProto{
			Vars:   vars,
			Values: values,
		},
	}
	m.addConstraint(c)
	return c
}

func (m *Model) asLinearExprProtos(exprs ...LinearExpr) []*swigpb.LinearExpressionProto {
	var ls []*swigpb.LinearExpressionProto
	for _, expr := range exprs {
		ls = append(ls, m.asLinearExprProto(expr))
	}
	return ls
}

func (m *Model) asLinearExprProto(expr LinearExpr) *swigpb.LinearExpressionProto {
	return &swigpb.LinearExpressionProto{
		Vars:   m.getIntVarIndexes(expr.vars...),
		Coeffs: expr.coeffs,
		Offset: expr.offset,
	}
}

func (m *Model) addIntVar(iv IntVar) {
	idx := len(m.proto.GetVariables())
	m.proto.Variables = append(m.proto.Variables, iv.proto)
	m.intVarToIdx[iv] = idx
}

func (m *Model) getIntVarIndexes(is ...IntVar) []int32 {
	var vars []int32
	for _, iv := range is {
		vars = append(vars, m.getIntVarIndex(iv))
	}
	return vars
}

func (m *Model) getIntVarIndex(iv IntVar) int32 {
	return int32(m.intVarToIdx[iv])
}

func (m *Model) addConstraint(c *constraint) {
	idx := len(m.proto.GetConstraints())
	m.proto.Constraints = append(m.proto.Constraints, c.proto)
	m.constraintToIdx[c] = idx
}

func (m *Model) constraintIndex(c *constraint) int {
	return m.constraintToIdx[c]
}

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
	proto            *swigpb.CpModelProto
	intVarIdxMap     map[*intVar]int
	constraintIdxMap map[*constraint]int
}

func NewModel() *Model {
	return &Model{
		proto:            &swigpb.CpModelProto{},
		intVarIdxMap:     make(map[*intVar]int),
		constraintIdxMap: make(map[*constraint]int),
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

func (m *Model) AddBoolOr(ls ...Literal) {
	c := newConstraint()
	literals := m.intVarIndexes(ls...)
	c.proto.Constraint = &swigpb.ConstraintProto_BoolOr{
		BoolOr: &swigpb.BoolArgumentProto{
			Literals: literals,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddBoolAnd(ls ...Literal) {
	c := newConstraint()
	literals := m.intVarIndexes(ls...)
	c.proto.Constraint = &swigpb.ConstraintProto_BoolAnd{
		BoolAnd: &swigpb.BoolArgumentProto{
			Literals: literals,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddBoolXor(ls ...Literal) {
	c := newConstraint()
	literals := m.intVarIndexes(ls...)
	c.proto.Constraint = &swigpb.ConstraintProto_BoolXor{
		BoolXor: &swigpb.BoolArgumentProto{
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

func (m *Model) AddLiteralElement(target Literal, index IntVar, ls ...Literal) {
	m.AddElement(target, index, ls...)
}

func (m *Model) AddAllDifferent(is ...IntVar) {
	c := newConstraint()
	vars := m.intVarIndexes(is...)
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
	vars := m.intVarIndexes(is...)
	c.proto.Constraint = &swigpb.ConstraintProto_Element{
		Element: &swigpb.ElementConstraintProto{
			Target: m.intVarIndex(target),
			Index:  m.intVarIndex(index),
			Vars:   vars,
		},
	}
	m.addConstraint(c)
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
	vars := m.intVarIndexes(is...)
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

func (m *Model) addIntVar(iv IntVar) {
	idx := len(m.proto.GetVariables())
	m.proto.Variables = append(m.proto.Variables, iv.proto)
	m.intVarIdxMap[iv] = idx
}

func (m *Model) intVarIndexes(is ...IntVar) []int32 {
	var vars []int32
	for _, iv := range is {
		vars = append(vars, m.intVarIndex(iv))
	}
	return vars
}

func (m *Model) intVarIndex(iv IntVar) int32 {
	return int32(m.intVarIdxMap[iv])
}

func (m *Model) addConstraint(c *constraint) {
	idx := len(m.proto.GetConstraints())
	m.proto.Constraints = append(m.proto.Constraints, c.proto)
	m.constraintIdxMap[c] = idx
}

func (m *Model) constraintIndex(c *constraint) int {
	return m.constraintIdxMap[c]
}

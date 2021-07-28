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

// Package cpsatsolver is a Go wrapper library for the CP-SAT solver included
// as part Google's Operations Research tools.
package cpsatsolver

import (
	"fmt"

	swigpb "github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

type Model struct {
	proto             *swigpb.CpModelProto
	intVarIndexes     map[*IntVar]int
	constraintIndexes map[*constraint]int
}

func NewModel() *Model {
	return &Model{
		proto:             &swigpb.CpModelProto{},
		intVarIndexes:     make(map[*IntVar]int),
		constraintIndexes: make(map[*constraint]int),
	}
}

func (m *Model) NewIntVar(lb int64, ub int64, name string) *IntVar {
	return m.NewIntVarFromDomain(NewDomain(lb, ub), name)
}

func (m *Model) NewBoolVar(name string) *IntVar {
	return m.NewIntVarFromDomain(NewDomain(0, 1), name)
}

func (m *Model) NewIntVarFromDomain(domain *Domain, name string) *IntVar {
	intVar := newIntVar(domain, name)
	m.addIntVar(intVar)
	return intVar
}

func (m *Model) NewConstant(c int64) *IntVar {
	return m.NewIntVarFromDomain(NewDomain(c, c), fmt.Sprintf("%d", c))
}

func (m *Model) AddAllDifferent(vs ...*IntVar) {
	c := newConstraint()
	var vars []int32
	for _, v := range vs {
		vars = append(vars, int32(m.intVarIndex(v)))
	}
	c.proto.Constraint = &swigpb.ConstraintProto_AllDiff{
		AllDiff: &swigpb.AllDifferentConstraintProto{
			Vars: vars,
		},
	}
	m.addConstraint(c)
}

func (m *Model) AddAllowedAssignment(vs []*IntVar, assignments ...[]int64) {
	_ = m.addAssignmentInternal(vs, assignments...)
	return
}

func (m *Model) AddForbiddenAssignment(vs []*IntVar, assignments ...[]int64) {
	c := m.addAssignmentInternal(vs, assignments...)
	c.proto.Constraint.(*swigpb.ConstraintProto_Table).Table.Negated = true
	return
}

func (m *Model) addAssignmentInternal(vs []*IntVar, assignments ...[]int64) *constraint {
	c := newConstraint()
	var vars []int32
	for _, v := range vs {
		vars = append(vars, int32(m.intVarIndex(v)))
	}
	var values []int64
	for _, assignment := range assignments {
		if len(assignment) != len(vs) {
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

func (m *Model) addIntVar(intVar *IntVar) {
	m.proto.Variables = append(m.proto.Variables, intVar.proto)
	m.intVarIndexes[intVar] = len(m.proto.Variables) - 1
}

func (m *Model) intVarIndex(intVar *IntVar) int {
	return m.intVarIndexes[intVar]
}

func (m *Model) addConstraint(constraint *constraint) {
	m.proto.Constraints = append(m.proto.Constraints, constraint.proto)
	m.constraintIndexes[constraint] = len(m.proto.Constraints) - 1
}

func (m *Model) constraintIndex(constraint *constraint) int {
	return m.constraintIndexes[constraint]
}

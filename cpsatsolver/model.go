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
	"errors"
	"fmt"

	swig "github.com/irfansharif/or-tools/internal/cpsatsolver"
	swigpb "github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

// Model is a constraint programming problem.
type Model struct {
	proto *swigpb.CpModelProto
}

// NewModel instantiates a new model.
func NewModel() *Model {
	return &Model{
		proto: &swigpb.CpModelProto{},
	}
}

// NewLiteral adds a new literal to the model.
func (m *Model) NewLiteral(name string) Literal {
	return m.NewIntVarFromDomain(NewDomain(0, 1), name)
}

// NewNegation adds a new literal to the model, one that's a negation of
// the given one. It uses a more efficient encoding than two literals with a
// model constraint xor-ing them together.
func (m *Model) NewNegation(l Literal, name string) Literal {
	return l.negation(name)
}

// NewIntVar adds a new integer variable to the model, one that's constrained to
// the given inclusive upper/lower bound.
func (m *Model) NewIntVar(lb int64, ub int64, name string) IntVar {
	return m.NewIntVarFromDomain(NewDomain(lb, ub), name)
}

// NewIntVarFromDomain adds a new integer variable to the model, one that's
// constrained to the given domain.
func (m *Model) NewIntVarFromDomain(d Domain, name string) IntVar {
	idx := len(m.proto.GetVariables())
	iv := newIntVar(d, int32(idx), name)
	m.proto.Variables = append(m.proto.Variables, iv.proto)
	return iv
}

// NewConstant adds a new constant to the model.
func (m *Model) NewConstant(c int64) IntVar {
	return m.NewIntVarFromDomain(NewDomain(c, c), fmt.Sprintf("%d", c))
}

// AddConstraints adds constraints to the model. When deciding on a solution,
// these constraints will need to be satisfied.
func (m *Model) AddConstraints(cs ...Constraint) {
	for _, c := range cs {
		m.proto.Constraints = append(m.proto.Constraints, c.proto)
	}
}

// Minimize sets a minimization objective for the model.
func (m *Model) Minimize(e LinearExpr) {
	m.proto.Objective = &swigpb.CpObjectiveProto{
		Vars:   e.vars(),
		Coeffs: e.coeffs(),
		Offset: float64(e.offset()),
	}
}

// Maximize sets a maximization objective for the model.
func (m *Model) Maximize(e LinearExpr) {
	m.Minimize(e)
	for i, coeff := range m.proto.GetObjective().GetCoeffs() {
		m.proto.GetObjective().GetCoeffs()[i] = -coeff
	}
	m.proto.GetObjective().ScalingFactor = -1
	m.proto.GetObjective().Offset = -m.proto.GetObjective().GetOffset()
}

// Validate checks whether the model is valid. If not, a descriptive error
// message is returned.
func (m *Model) Validate() (ok bool, _ error) {
	msg := swig.SatHelperValidateModel(*m.proto)
	if msg == "" {
		return true, nil
	}

	return false, errors.New(msg)
}

// Solve attempts to satisfy the model's constraints, if any, by deciding values
// for all the variables/literals that were instantiated into it. It returns the
// optimal result if an objective function is declared. If not, it returns
// the first found result that satisfies the model.
func (m *Model) Solve() Result {
	proto := swig.SatHelperSolve(*m.proto)
	return Result{proto: &proto}
}

// SolveAll returns all valid results that satisfy the model.
func (m *Model) SolveAll() []Result {
	var results []Result
	cb := &solutionCallback{
		cb: func(r Result) { results = append(results, r) },
	}
	cb.director = swig.NewDirectorSolutionCallback(cb)

	enumerate := true
	params := swigpb.SatParameters{EnumerateAllSolutions: &enumerate}
	swig.SatHelperSolveWithParametersAndSolutionCallback(*m.proto, params, cb.director)
	swig.DeleteDirectorSolutionCallback(cb.director)
	return results
}

type solutionCallback struct {
	cb       func(Result)
	director swig.SolutionCallback
}

func (p *solutionCallback) OnSolutionCallback() {
	proto := p.director.Response()
	p.cb(Result{proto: &proto})
}

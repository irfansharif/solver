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

// Package linearsolver is a Go wrapper library for the linear solver included
// as part of Google's Operations Research tools.
package linearsolver

import (
	"fmt"

	swig "github.com/irfansharif/or-tools/internal/linearsolver"
)

// Solver is the main type though which users build and solve problems.
//
// This is based on
// https://developers.google.com/optimization/reference/linear_solver/linear_solver/MPSolver.
type Solver struct {
	s swig.Solver
}

// New returns a new solver.
func New(name string, pt ProblemType) *Solver {
	return &Solver{
		s: swig.NewSolver(name, pt.asOptimizationProblemType()),
	}
}

// Close closes the solver. This must be called after New().
func (s *Solver) Close() error {
	swig.DeleteSolver(s.s)
	return nil
}

// Objective returns the objective object. Note that the objective is owned by
// the solver, and is initialized to its default value (see the MPObjective
// class) at construction.
func (s *Solver) Objective() *Objective {
	return &Objective{s.s.Objective()}
}

// NewVar creates a variable with the given bounds, integrality requirement and
// name. Bounds can be finite or +/- MPSolver::infinity(). The MPSolver owns the
// variable (i.e. the returned pointer is borrowed). Variable names are
// optional. If you give an empty name, name() will auto-generate one for you
// upon request.
func (s *Solver) NewVar(lowBound, upBound float64, integer bool, name string) *Variable {
	return &Variable{s.s.Var(lowBound, upBound, integer, name)}
}

// NewConstraintBounded returns a linear constraint with given bounds.
//
// Bounds can be finite or +/- MPSolver::infinity(). The MPSolver class assumes
// ownership of the constraint.
func (s *Solver) NewConstraintBounded(lowBound, upBound float64, name string) *Constraint {
	return &Constraint{s.s.Constraint(lowBound, upBound, name)}
}

// NumVariables returns the number of variable being optimized by the solver.
func (s *Solver) NumVariables() int {
	return s.s.NumVariables()
}

// NumConstraints returns the number of constraints.
func (s *Solver) NumConstraints() int {
	return s.s.NumConstraints()
}

// Solve solves the problem using the default parameter values.
func (s *Solver) Solve() error {
	code := s.s.Solve()
	switch code {
	case swig.SolverStatusOptimal:
		return nil
	case swig.SolverStatusAbnormal:
		return fmt.Errorf("abnormal status: this could be a numerical problem in the formulation or some other problem")
	case swig.SolverStatusFeasible:
	case swig.SolverStatusInfeasible:
	case swig.SolverStatusNotSolved:
	case swig.SolverStatusUnbounded:
	default:
	}
	return fmt.Errorf("unhandled status code %v", code)
}

// Variable is a variable to be optimized by the solver.
type Variable struct {
	v swig.Variable
}

// SolutionValue returns the value of the variable in the current solution.
//
// If the variable is integer, then the value will always be an integer (the
// underlying solver handles floating-point values only, but this function
// automatically rounds it to the nearest integer; see: man 3 round).
func (v *Variable) SolutionValue() float64 {
	return v.v.SolutionValue()
}

// Objective is the objective function to be optimized.
type Objective struct {
	o swig.Objective
}

// SetMaximization sets the optimization direction to maximize.
func (o *Objective) SetMaximization() {
	o.o.SetMaximization()
}

// SetMinimization sets the optimization direction to minimize.
func (o *Objective) SetMinimization() {
	o.o.SetMinimization()
}

// SetCoefficient sets the coefficient of the variable in the objective. If the
// variable does not belong to the solver, the function just returns, or crashes
// in non-opt mode.
func (o *Objective) SetCoefficient(v *Variable, coeff float64) {
	o.o.SetCoefficient(v.v, coeff)
}

// Value returns the objective value of the best solution found so far. It is
// the optimal objective value if the problem has been solved to optimality.
// Note: the objective value may be slightly different than what you could
// compute yourself using \c MPVariable::solution_value(); please use the
// --verify_solution flag to gain confidence about the numerical stability of
// your solution.
func (o *Objective) Value() float64 {
	return o.o.Value()
}

// Constraint is used for setting linear programming bounds.
type Constraint struct {
	c swig.Constraint
}

// SetCoefficient sets the coefficient of the variable on the constraint.
//
// If the variable does not belong to the solver, the function just returns, or
// crashes in non-opt mode.
func (c *Constraint) SetCoefficient(v *Variable, coeff float64) {
	c.c.SetCoefficient(v.v, coeff)
}

// ProblemType is a type of problem supported by the OR Tools linear solver.
type ProblemType int

const (
	GLOPLinearProgramming ProblemType = iota
)

// asOptimizationProblemType returns the SWIG version of the enum.
func (pt ProblemType) asOptimizationProblemType() swig.Operations_researchMPSolverOptimizationProblemType {
	switch pt {
	case GLOPLinearProgramming:
		return swig.SolverGLOP_LINEAR_PROGRAMMING
	default:
		panic("unknown problem type")
	}
}
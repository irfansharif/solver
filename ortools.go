// Package ortools is a go library for Google's Operations Research tools.
//
// Currently bindings are only provided for the linear solver.
package ortools

import (
	"fmt"

	ortoolsswig "github.com/irfansharif/or-tools/internal/swig"
)

// ProblemType is a type of OptimizationProblemType supported by the OR Tools library.
type ProblemType string

// swigEnum returns the SWIG version of the enum.
func (pt ProblemType) swigEnum() ortoolsswig.Operations_researchMPSolverOptimizationProblemType {
	switch pt {
	case LinearProgramming:
		return ortoolsswig.SolverGLOP_LINEAR_PROGRAMMING
	default:
		return 0
	}
}

// ProblemType definitions
const (
	LinearProgramming ProblemType = "LinearProgrammingProblemType"
)

// Solver is the main type though which users build and solve problems.
//
// This is based on
// https://developers.google.com/optimization/reference/linear_solver/linear_solver/MPSolver.
type Solver struct {
	s ortoolsswig.Solver
}

// NewSolver returns a new solver.
func NewSolver(name string, problemType ProblemType) *Solver {
	return &Solver{
		ortoolsswig.NewSolver(name, problemType.swigEnum()),
	}
}

// Close closes the solver. This must be called after NewSolver().
func (s *Solver) Close() error {
	ortoolsswig.DeleteSolver(s.s)
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
	case ortoolsswig.SolverStatusOptimal:
		return nil
	case ortoolsswig.SolverStatusAbnormal:
		return fmt.Errorf("solver returned abnormal status code; this could be a numerical problem in the formulation or some other problem")
	case ortoolsswig.SolverStatusFeasible:
	case ortoolsswig.SolverStatusInfeasible:
	case ortoolsswig.SolverStatusNotSolved:
	case ortoolsswig.SolverStatusUnbounded:
	default:
	}
	return fmt.Errorf("unhandled status code %v", code)
}

// Variable is a variable to be optimized by the solver.
type Variable struct {
	v ortoolsswig.Variable
}

// SolutionValue returns the value of the variable in the current solution.
//
// If the variable is integer, then the value will always be an integer (the
// underlying solver handles floating-point values only, but this function
// automatically rounds it to the nearest integer; see: man 3 round).
func (v *Variable) SolutionValue() float64 {
	return v.v.SolutionValue()
}

// Objective is the objective function to be optmized.
type Objective struct {
	o ortoolsswig.Objective
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
	c ortoolsswig.Constraint
}

// SetCoefficient sets the coefficient of the variable on the constraint.
//
// If the variable does not belong to the solver, the function just returns, or
// crashes in non-opt mode.
func (c *Constraint) SetCoefficient(v *Variable, coeff float64) {
	c.c.SetCoefficient(v.v, coeff)
}

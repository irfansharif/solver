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

package solver

import (
	"github.com/irfansharif/solver/internal/pb"
)

// Result is what's returned after attempting to solve a model.
type Result struct {
	pb *pb.CpSolverResponse
}

// Optimal is true iff a feasible solution has been found.
//
// More generally, this status represents a successful attempt at solving a
// model (if we've found a solution for a pure feasability problem, or if a gap
// limit has been specified and we've found solutions within the limit). In the
// case where we're iterating through all feasible solutions, this status will
// only be Feasible().
func (r Result) Optimal() bool {
	return r.pb.Status == pb.CpSolverStatus_OPTIMAL
}

// Infeasible is true iff the problem has been proven infeasible.
func (r Result) Infeasible() bool {
	return r.pb.Status == pb.CpSolverStatus_INFEASIBLE
}

// Feasible is true if a feasible solution has been found, and if we're
// enumerating through all solutions (if asked). See comment for Optimal for
// more details.
func (r Result) Feasible() bool {
	return r.pb.Status == pb.CpSolverStatus_FEASIBLE
}

func (r Result) Invalid() bool {
	return r.pb.Status == pb.CpSolverStatus_MODEL_INVALID
}

// Value returns the decided value of the given IntVar. This is only valid to
// use if the result is optimal or feasible.
func (r Result) Value(iv IntVar) int64 {
	return r.pb.GetSolution()[iv.index()]
}

// BooleanValue returns the decided value of the given Literal. This is only
// valid to use if the result is optimal or feasible.
func (r Result) BooleanValue(l Literal) bool {
	if l.isNegated() {
		return r.Value(l.Not()) == 0
	}

	return r.Value(l) == 1
}

// ObjectiveValue is the result of evaluating a model's objective function if
// the solution found is optimal or feasible. If no solution is found,
// then for a minimization problem, this will be an upper-bound of the objective
// of any feasible solution. For a maximization problem, it will be the
// lower-bound.
func (r Result) ObjectiveValue() float64 {
	return r.pb.GetObjectiveValue()
}

func (r Result) String() string {
	return "unimplemented" // XXX:
}

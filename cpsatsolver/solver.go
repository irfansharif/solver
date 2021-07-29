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
	swig "github.com/irfansharif/or-tools/internal/cpsatsolver"
	"github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

type Solver struct {
	model    *Model
	response Response
}

func NewSolver(m *Model) *Solver {
	return &Solver{
		model: m,
	}
}

func (s *Solver) Solve() Response {
	proto := swig.SatHelperSolve(*s.model.proto)
	s.response = Response{&proto}
	return s.response
}

func (s *Solver) Value(iv IntVar) int64 {
	return s.response.proto.GetSolution()[s.model.intVarIdxMap[iv]]
}

func (s *Solver) LiteralValue(l Literal) bool {
	return s.Value(l) == 1
}

type Response struct {
	proto *pb.CpSolverResponse
}

func (r Response) Optimal() bool {
	return r.proto.Status == pb.CpSolverStatus_OPTIMAL
}

func (r Response) Infeasible() bool {
	return r.proto.Status == pb.CpSolverStatus_INFEASIBLE
}

func (r Response) Feasible() bool {
	return r.proto.Status == pb.CpSolverStatus_FEASIBLE
}

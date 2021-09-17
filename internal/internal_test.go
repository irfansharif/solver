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

package internal

import (
	"testing"

	"github.com/irfansharif/solver/internal/pb"
)

func TestSimpleCPSAT(t *testing.T) {
	modelProto := pb.CpModelProto{}

	var numVals int64 = 3
	x := &pb.IntegerVariableProto{Name: "x", Domain: []int64{0, numVals - 1}}
	y := &pb.IntegerVariableProto{Name: "y", Domain: []int64{0, numVals - 1}}
	z := &pb.IntegerVariableProto{Name: "z", Domain: []int64{0, numVals - 1}}
	modelProto.Variables = append(modelProto.Variables, x) // idx: 0
	modelProto.Variables = append(modelProto.Variables, y) // idx: 1
	modelProto.Variables = append(modelProto.Variables, z) // idx: 2

	constraint := &pb.ConstraintProto{
		Name: "x != y",
		Constraint: &pb.ConstraintProto_AllDiff{
			AllDiff: &pb.AllDifferentConstraintProto{
				Vars: []int32{0, 1},
			},
		},
	}
	modelProto.Constraints = append(modelProto.Constraints, constraint)

	wrapper := NewSolveWrapper()
	response := wrapper.Solve(modelProto)
	if response.Status != pb.CpSolverStatus_FEASIBLE &&
		response.Status != pb.CpSolverStatus_OPTIMAL {
		t.Fatalf("expected solver to find solution")
	}

	{
		x := response.GetSolution()[0]
		y := response.GetSolution()[1]
		z := response.GetSolution()[2]

		for _, val := range []int64{x, y, z} {
			if val < 0 || val >= numVals {
				t.Fatalf("expected %d to be in domain [%d, %d)", val, 0, numVals)
			}
		}

		if x == y {
			t.Fatalf("x != y constrain violated, both found to be %d", x)
		}
	}
}

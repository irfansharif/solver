package cpsatsolver

import (
	"testing"

	"github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

func TestSimpleCPSAT(t *testing.T) {
	modelProto := pb.CpModelProto{
		Name: "model",
	}

	var numVals = 3
	x := &pb.IntegerVariableProto{Name: "x", Domain: []int64{0, int64(numVals - 1)}}
	y := &pb.IntegerVariableProto{Name: "y", Domain: []int64{0, int64(numVals - 1)}}
	z := &pb.IntegerVariableProto{Name: "z", Domain: []int64{0, int64(numVals - 1)}}

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
	response := SatHelperSolve(modelProto)
	if response.Status == pb.CpSolverStatus_FEASIBLE ||
		response.Status == pb.CpSolverStatus_OPTIMAL {
		t.Logf("x = %d", response.GetSolution()[0])
		t.Logf("y = %d", response.GetSolution()[1])
		t.Logf("z = %d", response.GetSolution()[2])
	}
}

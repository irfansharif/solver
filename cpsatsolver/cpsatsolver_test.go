package cpsatsolver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntVarAllDifferent(t *testing.T) {
	model := NewModel()

	var numVals int64 = 3
	x := model.NewIntVar(0, numVals-1, "x")
	y := model.NewIntVar(0, numVals-1, "y")
	z := model.NewIntVar(0, numVals-1, "z")
	model.AllDifferent(x, y, z)

	solver := NewSolver(model)
	response := solver.Solve()
	require.True(t, response.Optimal(), "expected solver to find solution")

	{
		x := solver.Value(x)
		y := solver.Value(y)
		z := solver.Value(z)

		for _, value := range []int64{x, y, z} {
			require.Truef(t, value >= 0 && value <= numVals-1,
				"expected %d to be in domain [%d, %d]", value, 0, numVals-1)
		}

		require.Falsef(t, x == y || x == z || y == z,
			"all different constraint violated, both x=%d y=%d z=%d", x, y, z)

		t.Logf("x = %d", x)
		t.Logf("y = %d", y)
		t.Logf("z = %d", z)
	}
}

func TestConstantRemainsSo(t *testing.T) {
	model := NewModel()
	value := int64(42)
	c := model.NewConstant(value)

	solver := NewSolver(model)
	response := solver.Solve()
	require.True(t, response.Optimal(), "expected solver to find solution")
	require.Equal(t, value, solver.Value(c))
}

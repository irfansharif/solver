package cpsatsolver

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntVarAllDifferent(t *testing.T) {
	model := NewModel()

	var numVals int64 = 3
	x := model.NewIntVar(0, numVals-1, "x")
	y := model.NewIntVar(0, numVals-1, "y")
	z := model.NewIntVar(0, numVals-1, "z")
	model.AddAllDifferent(x, y, z)

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
	}
}

func TestConstantRemainsSo(t *testing.T) {
	model := NewModel()
	value := int64(42)
	c := model.NewConstant(value)

	solver := NewSolver(model)
	require.True(t, solver.Solve().Optimal(), "expected solver to find solution")
	require.Equal(t, value, solver.Value(c))
}

func TestAllowedAssignments(t *testing.T) {
	model := NewModel()

	x := model.NewIntVar(0, 2, "x")
	y := model.NewIntVar(0, 2, "y")
	z := model.NewIntVar(0, 2, "z")

	assignments := [][]int64{
		{1, 2, 1},
		{2, 1, 0},
	}
	model.AddAllowedAssignments([]IntVar{x, y, z}, assignments)
	solver := NewSolver(model)

	require.True(t, solver.Solve().Optimal(), "expected solver to find solution")
	assignment := []int64{solver.Value(x), solver.Value(y), solver.Value(z)}
	require.True(t, reflect.DeepEqual(assignment, assignments[0]) ||
		reflect.DeepEqual(assignment, assignments[1]))
}

func TestForbiddenAssignments(t *testing.T) {
	model := NewModel()

	x := model.NewIntVar(1, 2, "x")
	y := model.NewIntVar(1, 2, "y")

	forbiddenAssignments := [][]int64{
		{1, 2},
		{2, 1},
	}
	model.AddForbiddenAssignments([]IntVar{x, y}, forbiddenAssignments)
	solver := NewSolver(model)

	require.True(t, solver.Solve().Optimal(), "expected solver to find solution")
	require.True(t, solver.Value(x) == solver.Value(y))
}

func TestConflictingAssignments(t *testing.T) {
	model := NewModel()

	x := model.NewIntVar(1, 2, "x")
	y := model.NewIntVar(1, 2, "y")

	assignments := [][]int64{
		{1, 2},
		{2, 1},
	}
	model.AddForbiddenAssignments([]IntVar{x, y}, assignments)
	model.AddAllowedAssignments([]IntVar{x, y}, assignments)

	solver := NewSolver(model)
	require.True(t, solver.Solve().Infeasible(), "didn't expect solver to find solution")
}

func TestBooleanConstraints(t *testing.T) {
	model := NewModel()

	a := model.NewLiteral("a")
	b := model.NewLiteral("b")
	c := model.NewLiteral("c")
	d := model.NewLiteral("d")
	e := model.NewLiteral("e")
	f := model.NewLiteral("f")

	model.AddBoolAnd(a, b) // a && b
	model.AddBoolOr(c, d)  // c || d
	model.AddBoolXor(e, f) // e != f

	solver := NewSolver(model)
	require.True(t, solver.Solve().Optimal(), "expected solver to find solution")

	{
		a := solver.LiteralValue(a)
		b := solver.LiteralValue(b)
		c := solver.LiteralValue(c)
		d := solver.LiteralValue(d)
		e := solver.LiteralValue(e)
		f := solver.LiteralValue(f)

		require.True(t, a && b)
		require.True(t, c || d)
		require.True(t, e != f)
	}
}

func TestAllowedBooleanAssignments(t *testing.T) {
	model := NewModel()

	a := model.NewLiteral("a")
	b := model.NewLiteral("b")

	assignments := [][]bool{
		{true, false},
		{false, true},
	}
	model.AddAllowedLiteralAssignments([]Literal{a, b}, assignments)

	solver := NewSolver(model)
	require.True(t, solver.Solve().Optimal(), "expected solver to find solution")

	{
		a := solver.LiteralValue(a)
		b := solver.LiteralValue(b)

		require.True(t, a != b)
	}
}

func TestForbiddenBooleanAssignments(t *testing.T) {
	model := NewModel()

	a := model.NewLiteral("a")
	b := model.NewLiteral("b")

	forbiddenAssignments := [][]bool{
		{true, false},
		{false, true},
	}
	model.AddForbiddenLiteralAssignments([]Literal{a, b}, forbiddenAssignments)

	solver := NewSolver(model)
	require.True(t, solver.Solve().Optimal(), "expected solver to find solution")

	{
		a := solver.LiteralValue(a)
		b := solver.LiteralValue(b)

		require.True(t, a == b)
	}
}

func TestElement(t *testing.T) {
	model := NewModel()
	var array []IntVar
	index := model.NewIntVar(0, 10, "index")
	target := model.NewIntVar(10, 100, "target")

	for i := 0; i <= 10; i += 1 {
		array = append(array, model.NewConstant(int64(i*10)))
	}

	model.AddElement(target, index, array...)
	solver := NewSolver(model)
	require.True(t, solver.Solve().Optimal(), "expected solver to find solution")
	require.True(t, solver.Value(target) == 10*solver.Value(index))
}

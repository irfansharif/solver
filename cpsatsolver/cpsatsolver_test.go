package cpsatsolver

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntVarAllDifferent(t *testing.T) {
	model := NewModel()

	var numVals int64 = 3
	x := model.NewIntVar(0, numVals-1, "x")
	y := model.NewIntVar(0, numVals-1, "y")
	z := model.NewIntVar(0, numVals-1, "z")

	ct := NewAllDifferentConstraint(x, y, z)
	model.AddConstraints(ct)

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		x := result.Value(x)
		y := result.Value(y)
		z := result.Value(z)

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

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")
	require.Equal(t, value, result.Value(c))
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
	ct := NewAllowedAssignmentsConstraint([]IntVar{x, y, z}, assignments)
	model.AddConstraints(ct)

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		assignment := []int64{result.Value(x), result.Value(y), result.Value(z)}
		require.True(t, reflect.DeepEqual(assignment, assignments[0]) ||
			reflect.DeepEqual(assignment, assignments[1]))
	}
}

func TestForbiddenAssignments(t *testing.T) {
	model := NewModel()

	x := model.NewIntVar(1, 2, "x")
	y := model.NewIntVar(1, 2, "y")

	forbiddenAssignments := [][]int64{
		{1, 2},
		{2, 1},
	}
	ct := NewForbiddenAssignmentsConstraint([]IntVar{x, y}, forbiddenAssignments)
	model.AddConstraints(ct)

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")
	require.True(t, result.Value(x) == result.Value(y))
}

func TestConflictingAssignments(t *testing.T) {
	model := NewModel()

	x := model.NewIntVar(1, 2, "x")
	y := model.NewIntVar(1, 2, "y")

	assignments := [][]int64{
		{1, 2},
		{2, 1},
	}

	ct1 := NewForbiddenAssignmentsConstraint([]IntVar{x, y}, assignments)
	ct2 := NewAllowedAssignmentsConstraint([]IntVar{x, y}, assignments)
	model.AddConstraints(ct1, ct2)

	result := model.Solve()
	require.True(t, result.Infeasible(), "didn't expect solver to find solution")
}

func TestBooleanConstraints(t *testing.T) {
	model := NewModel()

	a := model.NewLiteral("a")
	b := model.NewLiteral("b")
	c := model.NewLiteral("c")
	d := model.NewLiteral("d")
	e := model.NewLiteral("e")
	f := model.NewLiteral("f")

	and := NewBooleanAndConstraint(a, b) // a && b
	or := NewBooleanOrConstraint(c, d)   // c || d
	xor := NewBooleanXorConstraint(e, f) // e != f
	model.AddConstraints(and, or, xor)

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		a := result.BooleanValue(a)
		b := result.BooleanValue(b)
		c := result.BooleanValue(c)
		d := result.BooleanValue(d)
		e := result.BooleanValue(e)
		f := result.BooleanValue(f)

		require.True(t, a && b)
		require.True(t, c || d)
		require.True(t, e != f)
	}
}

func TestAllowedBooleanAssignments(t *testing.T) {
	model := NewModel()

	a := model.NewLiteral("a")
	b := model.NewLiteral("b")

	ct := NewAllowedLiteralAssignmentsConstraint([]Literal{a, b}, [][]bool{
		{true, false},
		{false, true},
	})
	model.AddConstraints(ct)

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		a := result.BooleanValue(a)
		b := result.BooleanValue(b)

		require.True(t, a != b)
	}
}

func TestForbiddenBooleanAssignments(t *testing.T) {
	model := NewModel()

	a := model.NewLiteral("a")
	b := model.NewLiteral("b")

	ct := NewForbiddenLiteralAssignmentsConstraint([]Literal{a, b}, [][]bool{
		{true, false},
		{false, true},
	})
	model.AddConstraints(ct)

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		a := result.BooleanValue(a)
		b := result.BooleanValue(b)

		require.True(t, a == b)
	}
}

// TestLinearExprMaximization is based on
// https://developers.google.com/optimization/lp/glop.
func TestLinearExprMaximization(t *testing.T) {
	model := NewModel()
	x := model.NewIntVar(0, 100, "x")
	y := model.NewIntVar(0, 100, "y")

	// Constraint 1: x + 2y <= 14.
	ct1 := NewLinearConstraint(
		NewLinearExpr([]IntVar{x, y}, []int64{1, 2}, 0),
		NewDomain(math.MinInt64, 14),
	)

	// Constraint 2: 3x - y >= 0.
	ct2 := NewLinearConstraint(
		NewLinearExpr([]IntVar{x, y}, []int64{3, -1}, 0),
		NewDomain(0, math.MaxInt64),
	)

	// Constraint 3: x - y <= 2.
	ct3 := NewLinearConstraint(
		NewLinearExpr([]IntVar{x, y}, []int64{1, -1}, 0),
		NewDomain(0, 2),
	)

	model.AddConstraints(ct1, ct2, ct3)

	// Objective function: 3x + 4y.
	model.Maximize(NewLinearExpr([]IntVar{x, y}, []int64{3, 4}, 0))

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		x := result.Value(x)
		y := result.Value(y)

		require.Equal(t, int64(6), x)
		require.Equal(t, int64(4), y)
		require.Equal(t, float64(34), result.ObjectiveValue())
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

	ct := NewElementConstraint(target, index, array...)
	model.AddConstraints(ct)

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")
	require.True(t, result.Value(target) == 10*result.Value(index))
}

func TestIterateThroughSolutions(t *testing.T) {
	model := NewModel()

	var numVals int64 = 3
	_ = model.NewIntVar(1, numVals, "x")

	results := model.SolveAll()
	require.Len(t, results, int(numVals))
}

func TestNegation(t *testing.T) {
	model := NewModel()

	A := model.NewLiteral("A")
	notA := model.NewNegation(A, "~A")

	model.AddConstraints(NewBooleanOrConstraint(A, notA))
	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		A, notA := result.BooleanValue(A), result.BooleanValue(notA)
		require.True(t, A || notA)
		require.True(t, A != notA)
	}
}

func TestNegationInfeasible(t *testing.T) {
	model := NewModel()

	A := model.NewLiteral("A")
	notA := model.NewNegation(A, "~A")

	model.AddConstraints(NewBooleanAndConstraint(A, notA))
	result := model.Solve()
	require.True(t, result.Infeasible(), "expected solver to not find solution")
}

func TestModelValidation(t *testing.T) {
	model := NewModel()

	_ = model.NewIntVar(0, math.MaxInt64, "a")

	ok, err := model.Validate()
	require.False(t, ok)
	require.True(t, strings.Contains(err.Error(),
		"domain do not fall in [kint64min + 2, kint64max - 1]"))
}

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
	"fmt"
	"math"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestingSetName is a testing-only helper to set the name of the model.
func (m *Model) TestingSetName(name string) {
	m.pb.Name = name
}

func TestIntVarAllDifferent(t *testing.T) {
	model := NewModel("")

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
	model := NewModel("")
	value := int64(42)
	c := model.NewConstant(value, "")

	t.Log(model.String())
	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")
	require.Equal(t, value, result.Value(c))
}

func TestAllowedAssignments(t *testing.T) {
	model := NewModel("")

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
	model := NewModel("")

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
	model := NewModel("")

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
	model := NewModel("")

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

	t.Log(model.String())
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
	model := NewModel("")

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
	model := NewModel("")

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
	model := NewModel("")
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

	t.Log(model.String())
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
	model := NewModel("")
	var array []IntVar
	index := model.NewIntVar(0, 10, "index")
	target := model.NewIntVar(10, 100, "target")

	for i := 0; i <= 10; i += 1 {
		array = append(array, model.NewConstant(int64(i*10), ""))
	}

	ct := NewElementConstraint(target, index, array...)
	model.AddConstraints(ct)

	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")
	require.True(t, result.Value(target) == 10*result.Value(index))
}

func TestEnumerateSolutions(t *testing.T) {
	model := NewModel("")

	var numVals int64 = 3
	_ = model.NewIntVar(1, numVals, "x")

	var results []Result
	_ = model.Solve(
		WithEnumeration(func(r Result) { results = append(results, r) }),
	)
	require.Len(t, results, int(numVals))
}

func TestNegation(t *testing.T) {
	model := NewModel("")

	A := model.NewLiteral("A")
	notA := A.Not()

	model.AddConstraints(NewBooleanOrConstraint(A, notA))

	t.Log(model.String())
	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		A, notA := result.BooleanValue(A), result.BooleanValue(notA)
		require.True(t, A || notA)
		require.True(t, A != notA)
	}
}

func TestNegationInfeasible(t *testing.T) {
	model := NewModel("")

	A := model.NewLiteral("A")
	notA := A.Not()

	model.AddConstraints(NewBooleanAndConstraint(A, notA))

	t.Log(model.String())
	result := model.Solve()
	require.True(t, result.Infeasible(), "expected solver to not find solution")
}

func TestModelValidation(t *testing.T) {
	model := NewModel("")

	_ = model.NewIntVar(0, math.MaxInt64, "a")

	ok, err := model.Validate()
	require.False(t, ok)
	require.True(t, strings.Contains(err.Error(),
		"domain do not fall in [kint64min + 2, kint64max - 1]"))
}

func TestAllSame(t *testing.T) {
	model := NewModel("")

	A := model.NewLiteral("A")
	B := model.NewLiteral("B")
	C := model.NewLiteral("C")

	model.AddConstraints(NewAllSameConstraint(A, B, C))

	t.Log(model.String())
	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		A, B, C := result.BooleanValue(A), result.BooleanValue(B), result.BooleanValue(C)
		require.True(t, A == B && B == C)
	}
}

func TestExactlyKLiterals(t *testing.T) {
	model := NewModel("")

	A := model.NewLiteral("A")
	B := model.NewLiteral("B")
	C := model.NewLiteral("C")
	D := model.NewLiteral("D")

	const k = 2
	model.AddConstraints(NewExactlyKConstraint(k, []Literal{A, B, C, D}...))

	t.Log(model.String())
	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		A := result.BooleanValue(A)
		B := result.BooleanValue(B)
		C := result.BooleanValue(C)
		D := result.BooleanValue(D)

		count := 0
		for _, b := range []bool{A, B, C, D} {
			if b {
				count += 1
			}
		}

		require.Equal(t, k, count)
	}
}

func TestNonOverlappingIntervalsWithEnforcement(t *testing.T) {
	model := NewModel("")

	lit := model.NewLiteral("a")
	model.AddConstraints(NewBooleanAndConstraint(lit))

	var intervals []Interval
	for i := 0; i < 3; i++ {
		start := model.NewIntVar(0, 10, fmt.Sprintf("start-%d", i))
		end := model.NewIntVar(0, 10, fmt.Sprintf("end-%d", i))
		size := model.NewIntVar(int64(i), 10, fmt.Sprintf("size-%d", i))

		intervals = append(intervals,
			model.NewInterval(start, end, size, "").OnlyEnforceIf(lit).(Interval))
	}

	model.AddConstraints(NewNonOverlappingConstraint(intervals...))
	valid, err := model.Validate()
	require.True(t, valid, err)

	t.Log(model.String())
	result := model.Solve()
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		type span struct {
			start, end, size int64
		}
		var spans []span

		for _, interval := range intervals {
			start, end, size := interval.Parameters()
			sp := span{
				start: result.Value(start),
				end:   result.Value(end),
				size:  result.Value(size),
			}
			spans = append(spans, sp)
		}

		sort.Slice(spans, func(i, j int) bool {
			return spans[i].start < spans[j].end
		})

		var last int64
		for _, sp := range spans {
			require.True(t, sp.start <= sp.end)
			require.True(t, sp.start+sp.size == sp.end)
			require.True(t, last <= sp.start)
			last = sp.end
		}
	}
}

func TestSolverOptions(t *testing.T) {
	model := NewModel("")

	A := model.NewLiteral("A")
	B := model.NewLiteral("B")
	C := model.NewLiteral("C")

	model.AddConstraints(NewAllSameConstraint(A, B, C))

	t.Log(model.String())
	result := model.Solve(
		WithLogger(os.Stdout, "[solver]  "),
		WithParallelism(4),
		WithTimeout(time.Second),
	)
	require.True(t, result.Optimal(), "expected solver to find solution")

	{
		A, B, C := result.BooleanValue(A), result.BooleanValue(B), result.BooleanValue(C)
		require.True(t, A == B && B == C)
	}
}

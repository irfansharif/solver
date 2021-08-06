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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinearExprString(t *testing.T) {
	model := NewModel("")
	a := model.NewIntVar(0, 10, "a")
	b := model.NewIntVar(0, 10, "b")
	c := model.NewIntVar(0, 10, "c")

	require.Equal(t, "a + b + c", Sum(a, b, c).String())
	require.Equal(t, "a + b", Sum(a, b).String())
	require.Equal(t, "a", Sum(a).String())
	require.Equal(t, "0a - b + 42c + 32", NewLinearExpr([]IntVar{a, b, c}, []int64{0, -1, 42}, 32).String())
	require.Equal(t, "-b + 42c", NewLinearExpr([]IntVar{b, c}, []int64{-1, 42}, 0).String())
	require.Equal(t, "-b + 42c + 10", NewLinearExpr([]IntVar{b, c}, []int64{-1, 42}, 10).String())
}

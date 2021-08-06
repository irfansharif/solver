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

func TestDomainValidation(t *testing.T) {
	require.NotPanics(t, func() { NewDomain(0, 0) })
	require.NotPanics(t, func() { NewDomain(0, 42) })
	require.NotPanics(t, func() { NewDomain(0, 42, 55, 70) })

	require.PanicsWithValue(t, "malformed domain: expected even number of interval boundaries",
		func() { NewDomain(0, 10, 11) })
	require.PanicsWithValue(t, "malformed domain: expected min <= max for 1st interval, found [42, 0]",
		func() { NewDomain(42, 0) })
	require.PanicsWithValue(t, "malformed domain: expected min <= max for 2nd interval, found [32, 12]",
		func() { NewDomain(0, 2, 32, 12) })
	require.PanicsWithValue(t, "malformed domain: expected 1st interval's max + 1 <  2nd interval's curMin, found [..., 2] [3, ...]",
		func() { NewDomain(0, 2, 3, 4) })
}

func TestDomainString(t *testing.T) {
	require.Equal(t, "[0,12]", NewDomain(0, 12).String())
	require.Equal(t, "[0,12] [24,32]", NewDomain(0, 12, 24, 32).String())
}

func TestDomainList(t *testing.T) {
	require.Equal(t, []int64{0, 12}, NewDomain(0, 12).list(0))
	require.Equal(t, []int64{0, 12, 24, 32}, NewDomain(0, 12, 24, 32).list(0))
	require.Equal(t, []int64{-2, 10, 22, 30}, NewDomain(0, 12, 24, 32).list(2))
}

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
	swigpb "github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

// Interval represents an interval parameterized by a start, end, and
// size. When added to a model, it automatically enforces the following
// relation:
//
//      start + size == end
//
// It can be used to define interval-based constraints. Constraints differ in how
// they interpret zero-sized intervals, and whether the end is exclusive.
type Interval interface {
	Constraint

	// Parameters returns the variables the interval is comprised of.
	Parameters() (start, end, size IntVar)

	index() int32
}

type interval struct {
	pb  *swigpb.ConstraintProto
	idx int32

	start, end, size IntVar
}

var _ Interval = &interval{}

func newInterval(start, end, size IntVar, idx int32) Interval {
	return &interval{
		start: start, end: end, size: size,
		idx: idx,
		pb: &swigpb.ConstraintProto{
			Constraint: &swigpb.ConstraintProto_Interval{
				Interval: &swigpb.IntervalConstraintProto{
					Start: start.index(),
					End:   end.index(),
					Size:  size.index(),
				},
			},
		},
	}
}

// Parameters is part of the Interval interface.
func (i *interval) Parameters() (start, end, size IntVar) {
	return i.start, i.end, i.size
}

// OnlyEnforceIf is part of the Interval interface.
func (i *interval) OnlyEnforceIf(literals ...Literal) Constraint {
	i.pb.EnforcementLiteral = asIntVars(literals).indexes()
	return i
}

func (i *interval) index() int32 {
	return i.idx
}

func (i *interval) protos() []*swigpb.ConstraintProto {
	return []*swigpb.ConstraintProto{i.pb}
}

type itrvals []Interval // named differently to avoid conflicts with variables using the plural form

func (is itrvals) indexes() []int32 {
	var indexes []int32
	for _, iv := range is {
		indexes = append(indexes, iv.index())
	}
	return indexes
}

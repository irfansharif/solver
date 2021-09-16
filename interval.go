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
	"strings"

	"github.com/irfansharif/solver/internal/pb"
)

// Interval represents an interval parameterized by a start, end, and
// size. When added to a model, it automatically enforces the following
// properties:
//
//      start + size == end
//      size >= 0
//
// It can be used to define interval-based constraints. Constraints differ in how
// they interpret zero-sized intervals and whether the end is considered
// exclusive.
type Interval interface {
	Constraint

	// Parameters returns the variables the interval is comprised of.
	Parameters() (start, end, size IntVar)

	// Stringer provides a printable format representation for the interval.
	fmt.Stringer

	index() int32
	name() string
}

type interval struct {
	pb  *pb.ConstraintProto
	idx int32

	start, end, size IntVar
	enforcement      Literal
}

var _ Interval = &interval{}

func newInterval(start, end, size IntVar, idx int32, name string) Interval {
	return &interval{
		start: start, end: end, size: size,
		idx: idx,
		pb: &pb.ConstraintProto{
			Name: name,
			Constraint: &pb.ConstraintProto_Interval{
				Interval: &pb.IntervalConstraintProto{
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
	if len(literals) > 1 {
		panic("intervals can only be enforced with a single literal")
	}
	i.pb.EnforcementLiteral = asIntVars(literals).indexes()
	if len(literals) == 1 {
		i.enforcement = literals[0]
	}
	return i
}

// String is part of the Interval interface.
func (i *interval) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%s, %s | %s]", i.start.name(), i.end.name(), i.size.name()))
	if i.enforcement != nil {
		b.WriteString(" if [")
		b.WriteString(i.enforcement.name())
		b.WriteString("]")
	}

	return b.String()
}

// index is part of the Interval interface.
func (i *interval) index() int32 {
	return i.idx
}

// name is part of the Interval interface.
func (i *interval) name() string {
	name := i.pb.GetName()
	if name == "" {
		name = "<unnamed>"
	}
	return name
}

// protos is part of the Constraint interface.
func (i *interval) protos() []*pb.ConstraintProto {
	return []*pb.ConstraintProto{i.pb}
}

type intervalList []Interval

func (is intervalList) indexes() []int32 {
	var indexes []int32
	for _, iv := range is {
		indexes = append(indexes, iv.index())
	}
	return indexes
}

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
	"fmt"

	"github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

// IntVar is an integer variable. It's typically constructed using a domain and
// is used as part of a model's constraints/objectives. When solving a model,
// the variable's integer value is decided on.
type IntVar interface {
	name() string
	index() int32
}

// Literal is a boolean variable. It's represented using an IntVar with a fixed
// domain [0, 1].
type Literal interface {
	IntVar

	// Not returns the negated form of literal. It uses a more efficient encoding
	// than two literals with a model constraint xor-ing them together.
	Not() Literal

	negated() bool
}

type intVar struct {
	pb  *pb.IntegerVariableProto
	idx int32
}

func newIntVar(d Domain, idx int32, name string) *intVar {
	return &intVar{
		pb: &pb.IntegerVariableProto{
			Name:   name,
			Domain: d.list(0),
		},
		idx: idx,
	}
}
func (i *intVar) name() string {
	return i.pb.Name
}

func (i *intVar) index() int32 {
	return i.idx
}

func (i *intVar) negated() bool {
	return i.idx < 0
}

// Not is part of the Literal interface.
func (i *intVar) Not() Literal {
	return &intVar{
		pb: &pb.IntegerVariableProto{
			Name:   fmt.Sprintf("~%s", i.name()),
			Domain: i.pb.Domain,
		},
		idx: -i.idx - 1,
	}
}

type intVars []IntVar

func (is intVars) indexes() []int32 {
	var indexes []int32
	for _, iv := range is {
		indexes = append(indexes, iv.index())
	}
	return indexes
}

type lits []Literal

func (ls lits) indexes() []int32 {
	var indexes []int32
	for _, l := range ls {
		indexes = append(indexes, l.index())
	}
	return indexes
}

func (ls lits) intVars() intVars {
	var res []IntVar
	for _, l := range ls {
		res = append(res, l.(IntVar))
	}
	return res
}

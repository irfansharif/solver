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
	"github.com/irfansharif/or-tools/internal/cpsatsolver/pb"
)

type IntVar = *intVar

type Literal = *intVar

type intVar struct {
	proto *pb.IntegerVariableProto
	idx   int32
}

func newIntVar(d *domain, idx int32, name string) *intVar {
	return &intVar{
		proto: &pb.IntegerVariableProto{
			Name:   name,
			Domain: d.list(0),
		},
		idx: idx,
	}
}

func (i *intVar) index() int32 {
	return i.idx
}

func (i *intVar) negated() bool {
	return i.idx < 0
}

func (i *intVar) name() string {
	return i.proto.Name
}

func (i *intVar) negation(name string) *intVar {
	return &intVar{
		proto: &pb.IntegerVariableProto{
			Name:   name,
			Domain: i.proto.Domain,
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

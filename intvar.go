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

	"github.com/irfansharif/solver/internal/pb"
)

// IntVar is an integer variable. It's typically constructed using a domain and
// is used as part of a model's constraints/objectives. When solving a model,
// the variable's integer value is decided on.
type IntVar interface {
	// Stringer provides a printable format representation for the int var.
	fmt.Stringer

	name() string
	index() int32
	domain() Domain
}

// Literal is a boolean variable. It's represented using an IntVar with a fixed
// domain [0, 1].
type Literal interface {
	IntVar

	// Not returns the negated form of literal. It uses a more efficient encoding
	// than two literals with a model constraint xor-ing them together.
	Not() Literal

	// isNegated returns true if the literal is negated.
	isNegated() bool
}

type intVar struct {
	pb  *pb.IntegerVariableProto
	idx int32
	d   Domain

	isLiteral, isConst bool
}

var _ IntVar = &intVar{}
var _ Literal = &intVar{}

func newIntVar(d Domain, idx int32, isLiteral, isConst bool, name string) *intVar {
	return &intVar{
		pb: &pb.IntegerVariableProto{
			Name:   name,
			Domain: d.list(0),
		},
		idx:       idx,
		d:         d,
		isLiteral: isLiteral,
		isConst:   isConst,
	}
}

// String is part of IntVar interface.
func (i *intVar) String() string {
	var domainStr string
	if i.isLiteral {
		domainStr = ""
	} else if i.isConst {
		domainStr = fmt.Sprintf(" == %d", i.d.list(0)[0])
	} else {
		domainStr = fmt.Sprintf(" in %s", i.d.String())
	}

	return fmt.Sprintf("%s%s", i.name(), domainStr)
}

// name is part of the IntVar interface.
func (i *intVar) name() string {
	name := i.pb.GetName()
	if name == "" {
		name = "<unnamed>"
	}
	return name
}

// index is part of the IntVar interface.
func (i *intVar) index() int32 {
	return i.idx
}

func (i *intVar) isNegated() bool {
	return i.idx < 0
}

func (i *intVar) domain() Domain {
	return i.d
}

// Not is part of the Literal interface.
func (i *intVar) Not() Literal {
	return &intVar{
		pb: &pb.IntegerVariableProto{
			Name:   fmt.Sprintf("~%s", i.name()),
			Domain: i.d.list(0),
		},
		idx: -i.idx - 1,
		d:   i.d,
	}
}

// AsIntVars is a convenience function to convert a slice of Literals to
// IntVars.
func AsIntVars(literals []Literal) []IntVar {
	var res []IntVar
	for _, l := range literals {
		res = append(res, l)
	}
	return res
}

type intVarList []IntVar

func (is intVarList) indexes() []int32 {
	var indexes []int32
	for _, iv := range is {
		indexes = append(indexes, iv.index())
	}
	return indexes
}

func asIntVars(literals []Literal) intVarList {
	return AsIntVars(literals)
}

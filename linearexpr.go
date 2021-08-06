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
	"strings"

	"github.com/irfansharif/solver/internal/pb"
)

// LinearExpr represents a linear expression of the form:
//
//   5x - 7y + 21z - 42
//
// In the expression above {x, y, z} are variables (IntVars) to be decided on by
// the model, {5, -7, 21} are coefficients for said variables, and -42 is the
// offset.
type LinearExpr interface {
	// Parameters returns the variables, coefficients, and offset the linear
	// expression is comprised of.
	Parameters() (vars []IntVar, coeffs []int64, offset int64)

	fmt.Stringer

	vars() []int32
	offset() int64
	coeffs() []int64
	proto() *pb.LinearExpressionProto
}

type linearExpr struct {
	pb *pb.LinearExpressionProto

	intVars []IntVar
}

// Parameters is part of the LinearExpr interface.
func (l *linearExpr) Parameters() (vars []IntVar, coeffs []int64, offset int64) {
	return l.intVars, l.coeffs(), l.offset()
}

var _ LinearExpr = &linearExpr{}

// Sum instantiates a new linear expression representing the sum of the given
// variables. It's a shorthand for NewLinearExpr with no offset and coefficients
// equal to one.
func Sum(vars ...IntVar) LinearExpr {
	var coeffs []int64
	for range vars {
		coeffs = append(coeffs, 1)
	}
	return NewLinearExpr(vars, coeffs, 0)
}

// TODO(irfansharif): We could instead construct a linear constraint bit-by-bit,
// setting coefficient per int var, setting offset, etc.
//
// 	expr := NewLinearExpr(
// 		WithVars(...),
// 		WithOffset(),
// 		WithCoeffs(),
// 	)
// 	expr.SetCoefficient(v, 2)
// 	expr.SetOffset(2)

// NewLinearExpr instantiates a new linear expression, representing:
//
//   sum(coefficients[i] * vars[i]) + offset
func NewLinearExpr(vars []IntVar, coeffs []int64, offset int64) LinearExpr {
	return &linearExpr{
		intVars: vars,
		pb: &pb.LinearExpressionProto{
			Vars:   intVarList(vars).indexes(),
			Coeffs: coeffs,
			Offset: offset,
		},
	}
}

// String is part of the LinearExpr interface.
func (l *linearExpr) String() string {
	var b strings.Builder
	for idx, v := range l.intVars {
		coeff := l.coeffs()[idx]
		var coeffStr, signStr string
		switch coeff {
		case 1:
			coeffStr = ""
		case -1:
			coeffStr = ""
		default:
			abs := int64(math.Abs(float64(coeff)))
			coeffStr = fmt.Sprintf("%d", abs)
		}

		signStr = "+"
		if coeff < 0 {
			signStr = "-"
		}

		if idx == 0 && coeff < 0 {
			coeffStr = fmt.Sprintf("-%s", coeffStr)
		} else if idx != 0 {
			b.WriteString(fmt.Sprintf(" %s ", signStr))
		}
		b.WriteString(fmt.Sprintf("%s%s", coeffStr, v.name()))
	}

	if offset := l.offset(); offset != 0 {
		abs := int64(math.Abs(float64(offset)))
		signStr := "+"
		if offset < 0 {
			signStr = "-"
		}
		b.WriteString(fmt.Sprintf(" %s ", signStr))
		b.WriteString(fmt.Sprintf("%d", abs))
	}
	return b.String()
}

// vars is part of the LinearExpr interface.
func (l *linearExpr) vars() []int32 {
	return l.pb.GetVars()
}

// offset is part of the LinearExpr interface.
func (l *linearExpr) offset() int64 {
	return l.pb.GetOffset()
}

// coeffs is part of the LinearExpr interface.
func (l *linearExpr) coeffs() []int64 {
	return l.pb.GetCoeffs()
}

// proto is part of the LinearExpr interface.
func (l *linearExpr) proto() *pb.LinearExpressionProto {
	return l.pb
}

type linearExprList []LinearExpr

func (le linearExprList) protos() []*pb.LinearExpressionProto {
	var ls []*pb.LinearExpressionProto
	for _, expr := range le {
		ls = append(ls, expr.proto())
	}
	return ls
}

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

// LinearExpr represents a linear expression of the form:
//
//   5x - 7y + 21z - 42
//
// In the expression above {x, y, z} are variables (IntVars) to be decided on by
// the model, {5, -7, 21} are coefficients for said variables, and -42 is the
// offset.
type LinearExpr interface {
	vars() []int32
	offset() int64
	coeffs() []int64
	proto() *swigpb.LinearExpressionProto
}

type linearExpr struct {
	pb *swigpb.LinearExpressionProto
}

func (l *linearExpr) proto() *swigpb.LinearExpressionProto {
	return l.pb
}

var _ LinearExpr = &linearExpr{}

// Sum instantiates a new linear expression representing the sum of the given
// variables.
func Sum(vars ...IntVar) LinearExpr {
	var coeffs []int64
	for range vars {
		coeffs = append(coeffs, 1)
	}
	return NewLinearExpr(vars, coeffs, 0)
}

// TODO(irfansharif): Could instead construct a linear constraint iteratively,
// setting coefficient per int var, setting offset, etc.
//
// 	expr := NewLinearExpr(
// 		WithVars(...),
// 		WithOffset(),
// 		WithCoeffs(),
// 	)
// 	expr := NewLinearExpr(vars...)
// 	expr.SetCoefficient(v, 2)
// 	expr.SetOffset(2)

// NewLinearExpr instantiates a new linear expression, representing:
//
//   sum(coefficients[i] * vars[i]) + offset
func NewLinearExpr(vars []IntVar, coeffs []int64, offset int64) LinearExpr {
	return &linearExpr{
		pb: &swigpb.LinearExpressionProto{
			Vars:   intVars(vars).indexes(),
			Coeffs: coeffs,
			Offset: offset,
		},
	}
}

func (l *linearExpr) vars() []int32 {
	return l.pb.GetVars()
}

func (l *linearExpr) offset() int64 {
	return l.pb.GetOffset()
}

func (l *linearExpr) coeffs() []int64 {
	return l.pb.GetCoeffs()
}

type linearExprs []LinearExpr

func (le linearExprs) protos() []*swigpb.LinearExpressionProto {
	var ls []*swigpb.LinearExpressionProto
	for _, expr := range le {
		ls = append(ls, expr.proto())
	}
	return ls
}

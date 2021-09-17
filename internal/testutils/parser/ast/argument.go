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

package ast

import (
	"fmt"
	"strings"

	"github.com/irfansharif/solver"
)

// Argument represents a statement argument (see Statement).
//
//   Argument = AssignmentsArgument
//            | BinaryOpArgument
//            | ConstantsArgument
//            | CumulativeArgument
//            | DomainArgument
//            | ElementArgument
//            | ImplicationArgument
//            | IntervalsArgument
//            | KArgument
//            | LinearEqualityArgument
//            | LinearExprsArgument
//            | NonOverlapping2DArgument
//            | VariableEqualityArgument
//            | VariablesArgument .
type Argument interface {
	fmt.Stringer
	argument()
}

// AssignmentsArgument represents an assignment argument: [a,b] ∉ [0,0] ∪ [1,1].
// It's used to test New{Allowed,Forbidden}{,Literal}AssignmentsConstraint.
//
//   AssignmentsArgument = "[" Variables "]" ( "∈" | "∉" ) ( NumbersList | BooleanList ) .
type AssignmentsArgument struct {
	Variables []string
	In        bool

	AllowedLiteralAssignments [][]bool // either-or
	AllowedIntVarAssignments  [][]int
}

func (a *AssignmentsArgument) String() string {
	var out strings.Builder
	out.WriteString("[")
	out.WriteString(strings.Join(a.Variables, ", "))
	out.WriteString("]")
	if a.In {
		out.WriteString(" ∈ ")
	} else {
		out.WriteString(" ∉ ")
	}

	if len(a.AllowedIntVarAssignments) != 0 {
		for i, arr := range a.AllowedIntVarAssignments {
			if i != 0 {
				out.WriteString(" ∪ ")
			}

			var inner []string
			for _, n := range arr {
				inner = append(inner, fmt.Sprintf("%d", n))
			}
			out.WriteString(fmt.Sprintf("[%s]", strings.Join(inner, ", ")))
		}
	} else {
		for i, arr := range a.AllowedLiteralAssignments {
			if i != 0 {
				out.WriteString(" ∪ ")
			}
			var inner []string
			for _, b := range arr {
				inner = append(inner, fmt.Sprintf("%t", b))
			}
			out.WriteString(fmt.Sprintf("[%s]", strings.Join(inner, ", ")))
		}
	}
	return out.String()
}

// ForLiterals returns true iff this argument refers to literals.
func (a *AssignmentsArgument) ForLiterals() bool {
	return len(a.AllowedLiteralAssignments) != 0
}

// ForIntVars returns true iff this argument refers to int vars.
func (a *AssignmentsArgument) ForIntVars() bool {
	return !a.ForLiterals()
}

// AsInt64s converts the underlying assignments int matrix to a matrix of
// int64s.
func (a *AssignmentsArgument) AsInt64s() [][]int64 {
	if !a.ForIntVars() {
		panic("not for int vars")
	}

	var assignments [][]int64
	for _, assignment := range a.AllowedIntVarAssignments {
		var conv []int64
		for _, v := range assignment {
			conv = append(conv, int64(v))
		}
		assignments = append(assignments, conv)
	}
	return assignments
}

// BinaryOpArgument represents a binary operation argument: a * b == c.
// It's used to test New{Product,Division,Modulo}Constraint.
//
//   BinaryOpArgument = Identifier ( "/" | "%" | "*" ) Identifier "==" Identifier .
type BinaryOpArgument struct {
	Left, Right, Op, Target string
}

func (b *BinaryOpArgument) String() string {
	return fmt.Sprintf("%s %s %s == %s", b.Left, b.Op, b.Right, b.Target)
}

// ConstantsArgument represents a constants argument: a, b to c == 42.
// It's used to test NewConstant.
//
//   ConstantsArgument = Variables "==" Number .
type ConstantsArgument struct {
	Variables []string
	Constant  int
}

func (c *ConstantsArgument) String() string {
	return fmt.Sprintf("%s == %d", strings.Join(c.Variables, ", "), c.Constant)
}

// CumulativeArgument represents a cumulative argument: i:2, j:4 | C.
// It's used to test NewCumulativeConstraint.
//
//   CumulativeArgument = IntervalDemands "|" Variable .
type CumulativeArgument struct {
	IntervalDemands []*IntervalDemand
	Capacity        string
}

// Intervals is a helper method that returns a slice of the underlying interval
// names.
func (c *CumulativeArgument) Intervals() []string {
	var intervals []string
	for _, id := range c.IntervalDemands {
		intervals = append(intervals, id.Name)
	}
	return intervals
}

// Demands is a helper method that returns a slice of the underlying demands.
func (c *CumulativeArgument) Demands() []int32 {
	var ds []int32
	for _, id := range c.IntervalDemands {
		ds = append(ds, int32(id.Demand))
	}
	return ds
}

func (c *CumulativeArgument) String() string {
	var strs []string
	for _, demand := range c.IntervalDemands {
		strs = append(strs, demand.String())
	}
	return fmt.Sprintf("%s | %s", strings.Join(strs, ", "), c.Capacity)
}

// DomainArgument represents a domain argument: a, b to d in [0, 2].
// It's used to test New{IntVar,LinearExpr}.
//
//   DomainArgument = ( Variables | LinearExprsMethod ) "in" Domains .
type DomainArgument struct {
	Variables   []string // either-or
	LinearExprs []*LinearExpr

	Domains []*Domain
}

func (d *DomainArgument) String() string {
	var terms, domains string
	if len(d.Variables) > 0 {
		terms = strings.Join(d.Variables, ", ")
	} else {
		var strs []string
		for _, expr := range d.LinearExprs {
			strs = append(strs, expr.String())
		}
		terms = strings.Join(strs, ", ")
	}

	var strs []string
	for _, domain := range d.Domains {
		strs = append(strs, domain.String())
	}
	domains = strings.Join(strs, " ∪ ")
	return fmt.Sprintf("%s in %s", terms, domains)
}

func (d *DomainArgument) AsSolverDomain() solver.Domain {
	var ls []int64
	for _, domain := range d.Domains {
		ls = append(ls, int64(domain.LowerBound), int64(domain.UpperBound))
	}
	return solver.NewDomain(ls[0], ls[1], ls[2:]...)
}

// ElementArgument represents an element argument: t == [a,b,c][i].
// It's used to test NewElementConstraint.
//
//   ElementArgument = Identifier "==" "[" Variables "]" "[" Identifier "]" .
type ElementArgument struct {
	Target, Index string
	Variables     []string
}

func (e *ElementArgument) String() string {
	return fmt.Sprintf("%s == [%s][%s]", e.Target, strings.Join(e.Variables, ", "), e.Index)
}

// ImplicationArgument represents an implication argument: a → b.
// It's used to test NewImplicationConstraint.
//
//   ImplicationArgument = Identifier "→" Identifier .
type ImplicationArgument struct {
	Left, Right string
}

func (i *ImplicationArgument) String() string {
	return fmt.Sprintf("%s → %s", i.Left, i.Right)
}

// IntervalsArgument represents an intervals argument: i as [s, e| sz]
// It's used to test NewInterval.
//
//   IntervalsArgument = Intervals .
type IntervalsArgument struct {
	Intervals []*Interval
}

func (i *IntervalsArgument) String() string {
	var strs []string
	for _, interval := range i.Intervals {
		strs = append(strs, interval.String())
	}
	return strings.Join(strs, ", ")
}

// KArgument represents a k-argument: a, b to f | 4.
// It's used to test New{AtLeast,AtMost,Exactly}KConstraint.
//
//   KArgument = Variables "|" Digits .
type KArgument struct {
	Literals []string
	K        int
}

func (k *KArgument) String() string {
	return fmt.Sprintf("%s | %d", strings.Join(k.Literals, ", "), k.K)
}

// LinearEqualityArgument represents a linear expression equality argument: 2j == max(k+i, i+2f).
// It's used to test NewLinear{Maximum,Minimum}Constraint
//
//   LinearEqualityArgument   = LinearExpr "==" ( "max" | "min" ) "(" LinearExprsMethod ")" .
type LinearEqualityArgument struct {
	Target *LinearExpr
	Exprs  []*LinearExpr
	Op     string
}

func (l *LinearEqualityArgument) String() string {
	var strs []string
	for _, expr := range l.Exprs {
		strs = append(strs, expr.String())
	}
	return fmt.Sprintf("%s == %s(%s)", l.Target, l.Op, strings.Join(strs, ", "))
}

// LinearExprsArgument represents an argument comprised of linear expressions.
//
//   LinearExprsArgument = LinearExprsMethod .
type LinearExprsArgument struct {
	Exprs []*LinearExpr
}

func (l *LinearExprsArgument) String() string {
	var strs []string
	for _, expr := range l.Exprs {
		strs = append(strs, expr.String())
	}
	return strings.Join(strs, ", ")
}

// NonOverlapping2DArgument represents an argument x and y interval variables,
// and a boolean indicating whether or not zero area boxes can overlap: [i, j], [k, l], false.
// It's used to test NewNonOverlapping2DConstraint.
//
//   NonOverlapping2DArgument = "[" Variables "]" "," "[" Variables "]" "," Boolean .
type NonOverlapping2DArgument struct {
	XVariables, YVariables    []string
	BoxesWithNoAreaCanOverlap bool
}

func (n *NonOverlapping2DArgument) String() string {
	return fmt.Sprintf("[%s], [%s], %t",
		strings.Join(n.XVariables, ", "), strings.Join(n.YVariables, ", "), n.BoxesWithNoAreaCanOverlap)
}

//   VariableEqualityArgument = Identifier "==" ( "max" | "min" ) "(" Variables ")" .

// VariableEqualityArgument represents a variable equality argument: j == min(k, i, f).
// It's used to test New{Minimum,Maximum}Constraint.
//
//   VariableEqualityArgument = Identifier "==" ( "max" | "min" ) "(" Variables ")" .
type VariableEqualityArgument struct {
	Target    string
	Variables []string
	Op        string
}

func (v *VariableEqualityArgument) String() string {
	return fmt.Sprintf("%s == %s(%s)", v.Target, v.Op, strings.Join(v.Variables, ", "))
}

// VariablesArgument represents an argument comprised of variables.
//
//   VariablesArgument        = Variables .
type VariablesArgument struct {
	Variables []string
}

func (v *VariablesArgument) String() string {
	return strings.Join(v.Variables, ", ")
}

// AsLinearExprsArgument returns a LinearExprsArgument representation of the
// VariablesArgument.
func (v *VariablesArgument) AsLinearExprsArgument() *LinearExprsArgument {
	argument := &LinearExprsArgument{}
	for _, variable := range v.Variables {
		argument.Exprs = append(argument.Exprs, &LinearExpr{
			LinearTerms: []*LinearTerm{
				{Coefficient: 1, Variable: variable},
			},
		})
	}
	return argument
}

var _ Argument = &AssignmentsArgument{}
var _ Argument = &BinaryOpArgument{}
var _ Argument = &ConstantsArgument{}
var _ Argument = &CumulativeArgument{}
var _ Argument = &DomainArgument{}
var _ Argument = &ElementArgument{}
var _ Argument = &ImplicationArgument{}
var _ Argument = &IntervalsArgument{}
var _ Argument = &KArgument{}
var _ Argument = &LinearEqualityArgument{}
var _ Argument = &LinearExprsArgument{}
var _ Argument = &NonOverlapping2DArgument{}
var _ Argument = &VariableEqualityArgument{}
var _ Argument = &VariablesArgument{}

func (*AssignmentsArgument) argument()      {}
func (*BinaryOpArgument) argument()         {}
func (*ConstantsArgument) argument()        {}
func (*CumulativeArgument) argument()       {}
func (*DomainArgument) argument()           {}
func (*ElementArgument) argument()          {}
func (*ImplicationArgument) argument()      {}
func (*IntervalsArgument) argument()        {}
func (*KArgument) argument()                {}
func (*LinearEqualityArgument) argument()   {}
func (*LinearExprsArgument) argument()      {}
func (*NonOverlapping2DArgument) argument() {}
func (*VariableEqualityArgument) argument() {}
func (*VariablesArgument) argument()        {}

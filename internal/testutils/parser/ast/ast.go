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
	"strconv"
	"strings"
)

// Statement represents a single statement.
//
//   Statement   = Receiver "." Method "(" [ Argument ] ")" [ Enforcement ] .
type Statement struct {
	Receiver    string
	Method      Method
	Argument    Argument
	Enforcement *Enforcement
}

func (s *Statement) String() string {
	argument, enforcement := "", ""
	if s.Argument != nil {
		argument = s.Argument.String()
	}
	if s.Enforcement != nil {
		enforcement = fmt.Sprintf(" %s", s.Enforcement.String())
	}
	return fmt.Sprintf("%s.%s(%s)%s", s.Receiver, s.Method, argument, enforcement)
}

// Enforcement represents the enforcement clause (see Statement).
//
//   Enforcement = "if" Variables .
type Enforcement struct {
	Literals []string
}

func (e *Enforcement) String() string {
	return fmt.Sprintf("if %s", strings.Join(e.Literals, ", "))
}

// Interval represents a single interval.
//
//   Interval       = Identifier "as" "[" Identifier "," Identifier "|" Identifier "]" .
type Interval struct {
	Name, Start, End, Size string // variables
}

func (i *Interval) String() string {
	return fmt.Sprintf("%s as [%s, %s | %s]", i.Name, i.Start, i.End, i.Size)
}

// Domain represents a unit domain.
//
//   Domain         = "[" Number "," Number "]" .
type Domain struct {
	LowerBound, UpperBound int
}

func (d *Domain) String() string {
	return fmt.Sprintf("[%d, %d]", d.LowerBound, d.UpperBound)
}

// LinearTerm represents an individual term in a linear expression (see
// LinearExpr). If the embedded variable is the empty string, the term is a just
// a constant.
//
//   LinearTerm     = { Digits } Identifier | Digits .
type LinearTerm struct {
	Coefficient int
	Variable    string
}

func (l *LinearTerm) String() string {
	if l.Coefficient == 1 {
		return fmt.Sprintf("%s", l.Variable)
	}
	return fmt.Sprintf("%d%s", l.Coefficient, l.Variable)
}

// LinearExpr represents a linear expression.
//
//   LinearExpr     = [ "-" ] LinearTerm { ( "+" | "-" ) LinearTerm } | "Σ" "(" Variables ")" .
type LinearExpr struct {
	LinearTerms []*LinearTerm
}

func (l *LinearExpr) String() string {
	var b strings.Builder
	for i, term := range l.LinearTerms {
		var sign, coeff string
		if i == 0 {
			if term.Coefficient < 0 {
				sign = "-"
			}
		} else {
			if term.Coefficient < 0 {
				sign = " - "
			} else {
				sign = " + "
			}
		}

		if term.Coefficient != 1 && term.Coefficient != -1 {
			abs := int64(term.Coefficient)
			if term.Coefficient < 0 {
				abs = -abs
			}
			coeff = strconv.FormatInt(abs, 10)
		}

		b.WriteString(fmt.Sprintf("%s%s%s", sign, coeff, term.Variable))
	}
	return b.String()
}

// IntervalDemand represents an interval identifier and it's corresponding
// demand.
//
//   IntervalDemand = Identifier ":" Identifier .
type IntervalDemand struct {
	Name   string
	Demand string
}

func (i *IntervalDemand) String() string {
	return fmt.Sprintf("%s: %s", i.Name, i.Demand)
}

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

type Method int

const (
	Unrecognized Method = 0

	AllDifferentMethod Method = iota + 128
	AllSameMethod
	AssignmentsMethod
	AtLeastKMethod
	AtMostKMethod
	BinaryOpMethod
	BooleanAndMethod
	BooleanOrMethod
	BooleanXorMethod
	BooleansMethod
	ConstantsMethod
	CumulativeMethod
	ElementMethod
	EqualityMethod
	ExactlyKMethod
	FeasibleMethod
	ImplicationMethod
	InfeasibleMethod
	IntervalsMethod
	InvalidMethod
	LinearExprsMethod
	LiteralsMethod
	MaximizeMethod
	MinimizeMethod
	NameMethod
	NonOverlappingMethod
	NonOverlapping2DMethod
	ObjectiveValueMethod
	OptimalMethod
	PrintMethod
	SolveMethod
	SolveAllMethod
	ValidateMethod
	ValuesMethod
	VarsMethod
)

var methods = map[Method]string{
	AllDifferentMethod:     "all-different",
	AllSameMethod:          "all-same",
	AssignmentsMethod:      "assignments",
	AtLeastKMethod:         "at-least-k",
	AtMostKMethod:          "at-most-k",
	BinaryOpMethod:         "binary-op",
	BooleanAndMethod:       "boolean-and",
	BooleanOrMethod:        "boolean-or",
	BooleanXorMethod:       "boolean-xor",
	BooleansMethod:         "booleans",
	ConstantsMethod:        "constants",
	CumulativeMethod:       "cumulative",
	ElementMethod:          "element",
	EqualityMethod:         "equality",
	ExactlyKMethod:         "exactly-k",
	FeasibleMethod:         "feasible",
	ImplicationMethod:      "implication",
	InfeasibleMethod:       "infeasible",
	IntervalsMethod:        "intervals",
	InvalidMethod:          "invalid",
	LinearExprsMethod:      "linear-exprs",
	LiteralsMethod:         "literals",
	MaximizeMethod:         "maximize",
	MinimizeMethod:         "minimize",
	NameMethod:             "name",
	NonOverlappingMethod:   "non-overlapping",
	NonOverlapping2DMethod: "non-overlapping-2D",
	ObjectiveValueMethod:   "objective-value",
	OptimalMethod:          "optimal",
	PrintMethod:            "print",
	SolveMethod:            "solve",
	SolveAllMethod:         "solve-all",
	ValidateMethod:         "validate",
	ValuesMethod:           "values",
	VarsMethod:             "vars",
}

var lookup = make(map[string]Method)

func init() {
	for m, s := range methods {
		lookup[s] = m
	}
}

func LookupMethod(s string) (Method, bool) {
	m, ok := lookup[s]
	return m, ok
}

func (m Method) String() string {
	s, ok := methods[m]
	if !ok {
		panic("unrecognized method")
	}
	return s
}

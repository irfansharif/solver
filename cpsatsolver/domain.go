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
	"math"
)

// Domain represents n disjoint intervals, each of the form [min, max]:
//
// 		[min_0, max_0,  ..., min_{n-1}, max_{n-1}].
//
// The following constraints hold:
// 		- For all i < n   :      min_i <= max_i
// 		- For all i < n-1 :  max_i + 1 < min_{i+1}.
//
// The most common example being just [min, max]. If min == max, then this is a
// constant variable.
//
// NB: We check at validation that a variable domain is small enough so
// that we don't run into integer overflow in our algorithms. Avoid having
// "unbounded" variables like [0, math.MaxInt64], opting instead for tighter
// domains.
type Domain interface {
	list(shift int64) []int64
}

type domain struct {
	intervals []int64
}

// NewDomain instantiates a new domain using the given intervals.
func NewDomain(lb, ub int64, ds ...int64) Domain {
	if len(ds)%2 != 0 {
		panic("malformed domain")
	}
	intervals := []int64{lb, ub}
	intervals = append(intervals, ds...)
	return &domain{intervals: intervals}
}

func (d *domain) list(shift int64) []int64 {
	var ls []int64
	for _, v := range d.intervals {
		if v == math.MaxInt64 {
			ls = append(ls, v)
		} else {
			ls = append(ls, v-shift)
		}
	}

	return ls
}

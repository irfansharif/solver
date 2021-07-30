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

type domain struct {
	lb, ub int64
}

func NewDomain(lb, ub int64) Domain {
	return &domain{lb, ub}
}

func (d *domain) list(offset int64) []int64 {
	var ls []int64
	for _, v := range []int64{d.lb, d.ub} {
		if v == math.MaxInt64 {
			ls = append(ls, v)
		} else {
			ls = append(ls, v-offset)
		}
	}

	return ls
}

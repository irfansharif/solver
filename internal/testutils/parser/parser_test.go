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

package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cockroachdb/datadriven"
	"github.com/irfansharif/solver/internal/testutils/bazel"
	"github.com/irfansharif/solver/internal/testutils/parser/ast"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/ebnf"
)

func TestDatadriven(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		// path, implant = bazel.WritableSandboxPathFor(t, filepath.Join("internal", "testutils", "parser", path))
		path, implant := bazel.WritableSandboxPathFor(t, "internal/testutils/parser", path)
		defer implant()

		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			p := New(d.Input)
			var out string
			var err error
			switch d.Cmd {
			case "receiver":
				out, err = p.Receiver()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
			case "identifier":
				out, err = p.Identifier()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
			case "method":
				var method ast.Method
				method, err = p.Method()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = method.String()
			case "variable":
				out, err = p.Variable()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
			case "variables":
				var variables []string
				variables, err = p.Variables()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = strings.Join(variables, ", ")
			case "enforcement":
				var e *ast.Enforcement
				e, err = p.Enforcement()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = e.String()
			case "interval":
				var i *ast.Interval
				i, err = p.Interval()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = i.String()
			case "boolean":
				var boolean bool
				boolean, err = p.Boolean()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = fmt.Sprintf("%t", boolean)
			case "booleans":
				var booleans []bool
				booleans, err = p.Booleans()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				var strs []string
				for _, boolean := range booleans {
					strs = append(strs, fmt.Sprintf("%t", boolean))
				}
				out = strings.Join(strs, ", ")
			case "number":
				var n int
				n, err = p.Number()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = fmt.Sprintf("%d", n)
			case "numbers":
				var numbers []int
				numbers, err = p.Numbers()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				var strs []string
				for _, number := range numbers {
					strs = append(strs, fmt.Sprintf("%d", number))
				}
				out = strings.Join(strs, ", ")
			case "intervals":
				var intervals []*ast.Interval
				intervals, err = p.Intervals()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				var strs []string
				for _, interval := range intervals {
					strs = append(strs, interval.String())
				}
				out = strings.Join(strs, ", ")
			case "interval-demand":
				var demand *ast.IntervalDemand
				demand, err = p.IntervalDemand()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = demand.String()
			case "domain":
				var domain *ast.Domain
				domain, err = p.Domain()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = domain.String()
			case "linear-term":
				var term *ast.LinearTerm
				term, err = p.LinearTerm()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = term.String()
			case "linear-expr":
				var expr *ast.LinearExpr
				expr, err = p.LinearExpr()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = expr.String()
			case "linear-exprs":
				var exprs []*ast.LinearExpr
				exprs, err = p.LinearExprs()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				var strs []string
				for _, expr := range exprs {
					strs = append(strs, expr.String())
				}
				out = strings.Join(strs, ", ")
			case "domains":
				var domains []*ast.Domain
				domains, err = p.Domains()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				var strs []string
				for _, domain := range domains {
					strs = append(strs, domain.String())
				}
				out = strings.Join(strs, " ∪ ")
			case "statement":
				var s *ast.Statement
				s, err = p.Statement()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = s.String()
			case "numbers-list":
				var list [][]int
				list, err = p.NumbersList()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}

				var strs []string
				for _, l := range list {
					var inner []string
					for _, n := range l {
						inner = append(inner, fmt.Sprintf("%d", n))
					}
					strs = append(strs, fmt.Sprintf("[%s]", strings.Join(inner, ", ")))
				}
				out = strings.Join(strs, " ∪ ")
			case "booleans-list":
				var list [][]bool
				list, err = p.BooleansList()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}

				var strs []string
				for _, l := range list {
					var inner []string
					for _, b := range l {
						inner = append(inner, fmt.Sprintf("%t", b))
					}
					strs = append(strs, fmt.Sprintf("[%s]", strings.Join(inner, ", ")))
				}
				out = strings.Join(strs, " ∪ ")
			case "assignments-argument":
				var arg ast.Argument
				arg, err = p.AssignmentsArgument()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = arg.String()
			case "binary-op-argument":
				var arg ast.Argument
				arg, err = p.BinaryOpArgument()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = arg.String()
			case "constants-argument":
				var arg ast.Argument
				arg, err = p.ConstantsArgument()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = arg.String()
			case "cumulative-argument":
				var arg ast.Argument
				arg, err = p.CumulativeArgument()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = arg.String()
			case "k-argument":
				var arg ast.Argument
				arg, err = p.KArgument()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = arg.String()
			case "domain-argument":
				var arg ast.Argument
				arg, err = p.DomainArgument()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = arg.String()
			case "element-argument":
				var arg ast.Argument
				arg, err = p.ElementArgument()
				if err != nil {
					return fmt.Sprintf("err: %s", err)
				}
				out = arg.String()
			default:
				t.Errorf("unrecognized command: %s", d.Cmd)
			}

			require.Nil(t, err)
			if !p.EOF() {
				return fmt.Sprintf("err: expected EOF; parsed %q", out)
			}
			return out
		})
	})
}

func TestGrammar(t *testing.T) {
	filename := `grammar.ebnf`
	contents, err := ioutil.ReadFile(filename)
	require.Nil(t, err)

	grammar, err := ebnf.Parse(filename, bytes.NewReader(contents))
	if err != nil {
		t.Fatal(err)
	}
	if err := ebnf.Verify(grammar, `Statement`); err != nil { // verify the top-level statement
		t.Fatal(err)
	}
}

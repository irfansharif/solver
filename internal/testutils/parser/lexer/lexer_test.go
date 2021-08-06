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

package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/datadriven"
	"github.com/irfansharif/solver/internal/testutils/bazel"
	"github.com/irfansharif/solver/internal/testutils/parser/lexer"
	"github.com/irfansharif/solver/internal/testutils/parser/token"
)

func TestDatadriven(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		if bazel.BuiltWithBazel() {
			var implant func()
			path, implant = bazel.WritableSandboxPathFor(t, "internal/testutils/parser/lexer", path)
			defer implant()
		}

		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			var out strings.Builder
			l := lexer.New(d.Input)
			for {
				tok := l.Next()
				if tok.Type == token.EOF {
					break
				}

				out.WriteString(fmt.Sprintf("%s %q", tok.Type, tok.Value))
				out.WriteString("\n")
			}
			return out.String()
		})
	})
}

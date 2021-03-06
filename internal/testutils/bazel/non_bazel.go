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

// +build !bazel

package bazel

import (
	"testing"
)

// BuiltWithBazel returns true iff this library was built with Bazel.
func BuiltWithBazel() bool {
	return false
}

// WritableSandboxPathFor returns a writable path for a given path in the source
// tree. It also returns a callback to copy over the contents of the writable
// path back into the source tree (`--sandbox_writable_path` is expected to have
// been declared over it).
func WritableSandboxPathFor(t *testing.T, pkg, path string) (writable string, implant func()) {
	panic("not built with Bazel")
}

// WorkspacePath returns the path of the bazel workspace.
func WorkspacePath(t *testing.T) string {
	panic("not built with Bazel")
}

// ScratchDirectory returns the path of the scratch directory in the bazel
// sandbox.
func ScratchDirectory(t *testing.T) string {
	panic("not built with Bazel")
}

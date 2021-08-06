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

// +build bazel

package bazel

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// BuiltWithBazel returns true iff this library was built with Bazel.
func BuiltWithBazel() bool {
	return true
}

// WritableSandboxPathFor returns a writable path for a given path in the source
// tree. It also returns a callback to copy over the contents of the writable
// path back into the source tree (`--sandbox_writable_path` is expected to have
// been declared over it).
func WritableSandboxPathFor(t *testing.T, pkg, path string) (writable string, implant func()) {
	cp := func(t *testing.T, src, dest string) { // shorthand to copy a file from src to dest
		input, err := ioutil.ReadFile(src)
		require.Nil(t, err)
		require.Nil(t, os.MkdirAll(filepath.Dir(dest), 0755))
		require.Nil(t, ioutil.WriteFile(dest, input, 0644))
	}

	workspace := os.Getenv("BAZEL_WORKSPACE")
	if workspace == "" {
		t.Fatal("BAZEL_WORKSPACE unset")
	}

	outdir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if outdir == "" {
		t.Fatal("expected to find TEST_UNDECLARED_OUTPUTS_DIR")
	}
	dest := filepath.Join(outdir, path)

	cp(t, path, dest)
	return dest, func() { cp(t, dest, filepath.Join(workspace, pkg, path)) }
}

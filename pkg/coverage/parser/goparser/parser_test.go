// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goparser

import (
	"go/token"
	"io/ioutil"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_NodeFromFilePath_bad_path(t *testing.T) {
	g := NewGomegaWithT(t)

	srcPath := "foo.go"

	fset := token.NewFileSet()
	_, err := NodeFromFilePath(srcPath, fset)
	g.Expect(err).ToNot(BeNil())
}

func Test_NodeFromFilePath_bad_file(t *testing.T) {
	g := NewGomegaWithT(t)

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Errorf("could not create temp dir")
		t.FailNow()
	}

	srcContent := `
func Meow(x, y int) bool {
  if x > y {
	  return true
  }
	return false
}
`
	srcPath := filepath.Join(dir, "src.go")
	err = ioutil.WriteFile(srcPath, []byte(srcContent), 0644)

	if err != nil {
		t.Errorf("could not write to temp file %v", err)
		t.FailNow()
	}

	fset := token.NewFileSet()
	_, err = NodeFromFilePath(srcPath, fset)
	g.Expect(err).ToNot(BeNil())
}

func Test_NodeFromFilePath(t *testing.T) {
	g := NewGomegaWithT(t)

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Errorf("could not create temp dir")
		t.FailNow()
	}

	srcContent := `
package foo

func Meow(x, y int) bool {
  if x > y {
	  return true
  }
	return false
}
`
	srcPath := filepath.Join(dir, "src.go")
	err = ioutil.WriteFile(srcPath, []byte(srcContent), 0644)

	if err != nil {
		t.Errorf("could not write to temp file %v", err)
		t.FailNow()
	}

	fset := token.NewFileSet()
	astFile, err := NodeFromFilePath(srcPath, fset)
	g.Expect(err).To(BeNil())
	g.Expect(astFile).ToNot(BeNil())
}

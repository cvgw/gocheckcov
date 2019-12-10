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

package functions

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_CollectFunctions(t *testing.T) {
	g := NewGomegaWithT(t)

	file, err := ioutil.TempFile("", "profile.test")
	if err != nil {
		t.Errorf("could not create temp file")
		t.FailNow()
	}

	profileFileContent := `
package foo

func Meow(x, y int) bool {
  if x > y {
	  return true
  }
	return false
}
`
	err = ioutil.WriteFile(file.Name(), []byte(profileFileContent), 0644)

	if err != nil {
		t.Errorf("could not write to temp file %v", err)
		t.FailNow()
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file.Name(), []byte(profileFileContent), 0)

	if err != nil {
		t.Errorf("could not create ast for file %v %v", file.Name(), err)
		t.FailNow()
	}

	funcs, err := CollectFunctions(f, fset, file.Name())
	g.Expect(err).To(BeNil())
	g.Expect(funcs).ToNot(BeNil())
	g.Expect(funcs).To(HaveLen(1))
}

//func MatchFunc(expected Function) types.GomegaMatcher {
//  return &funcMatcher{
//    expected: expected,
//  }
//}

//type funcMatcher struct {
//  expected Function
//}

//func (f *funcMatcher) Match(actual interface{}) (success bool, err error) {
//  switch x := actual.(type) {
//  case Function:

//    if x.Name != f.expected.Name {
//      return false, nil
//    }

//    if x.SrcPath != f.expected.SrcPath {
//      return false, nil
//    }

//    if x.StartOffset != f.expected.StartOffset {
//      return false, nil
//    }

//    if x.StartCol != f.expected.StartCol {
//      return false, nil
//    }

//    if x.StartLine != f.expected.StartLine {
//      return false, nil
//    }

//    if x.EndOffset != f.expected.EndOffset {
//      return false, nil
//    }

//    if x.EndCol != f.expected.EndCol {
//      return false, nil
//    }

//    if x.EndLine != f.expected.EndLine {
//      return false, nil
//    }

//    if len(x.Statements) != len(f.expected.Statements) {
//      return false, nil
//    }

//    return true, nil
//  default:
//    return false, fmt.Errorf("actual must be of type Function")
//  }
//}

//func (f *funcMatcher) FailureMessage(actual interface{}) (message string) {
//  return fmt.Sprintf("expected %v to equal %v", f.expected, actual)
//}

//func (f *funcMatcher) NegatedFailureMessage(actual interface{}) (message string) {
//  return fmt.Sprintf("expected %v to not equal %v", f.expected, actual)
//}

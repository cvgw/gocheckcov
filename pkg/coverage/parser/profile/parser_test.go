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

package profile

import (
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/cvgw/gocheckcov/pkg/coverage/parser/goparser"
	"github.com/cvgw/gocheckcov/pkg/coverage/parser/goparser/functions"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/cover"
)

//func Test_NodesFromProfiles(t *testing.T) {
//  type testcase struct {
//    description    string
//    srcFile        string
//    srcFileContent string
//    expectErr      bool
//    profiles       []*cover.Profile
//  }

//  testCases := []testcase{
//    func() testcase {
//      filename := "src.go"
//      file, err := ioutil.TempFile("", filename)
//      if err != nil {
//        t.Errorf("could not create temp file")
//        t.FailNow()
//      }

//      srcFileContent := `
//package foo

//func Meow(x, y int) bool {
//  if x > y {
//    return true
//  }
//  return false
//}
//`
//      err = ioutil.WriteFile(file.Name(), []byte(srcFileContent), 0644)

//      if err != nil {
//        t.Errorf("could not write to temp file %v", err)
//        t.FailNow()
//      }

//      return testcase{
//        description: "one profile with a valid src file",
//        srcFile:     "profile.test",
//        profiles: []*cover.Profile{
//          &cover.Profile{FileName: file.Name()},
//        },
//        srcFileContent: srcFileContent,
//      }
//    }(),
//    testcase{
//      description: "one profile with a blank file name",
//      srcFile:     "profile.test",
//      profiles: []*cover.Profile{
//        &cover.Profile{},
//      },
//      expectErr: true,
//    },
//    testcase{
//      description: "one profile with a bad file name",
//      srcFile:     "profile.test",
//      profiles: []*cover.Profile{
//        &cover.Profile{FileName: "meow"},
//      },
//      expectErr: true,
//    },
//  }
//  for i := range testCases {
//    tc := testCases[i]
//    t.Run(tc.description, func(t *testing.T) {
//      g := NewGomegaWithT(t)

//      profiles := tc.profiles

//      fset := token.NewFileSet()

//      res, err := NodesFromProfiles("", profiles, fset)
//      if tc.expectErr {
//        g.Expect(err).ToNot(BeNil())
//        g.Expect(res).To(BeNil())
//      } else {
//        g.Expect(err).To(BeNil())
//        g.Expect(res).To(HaveLen(1))
//      }
//    })
//  }
//}

func Test_NodeFromFilePath(t *testing.T) {
	type testcase struct {
		description    string
		srcFile        string
		srcFileContent string
		noSrcFile      bool
		expectErr      bool
	}

	testCases := []testcase{
		testcase{
			description: "src file with one function",
			srcFile:     "profile.test",
			srcFileContent: `
package foo

func Meow(x, y int) bool {
  if x > y {
	  return true
  }
	return false
}
`,
		},
		testcase{
			description: "bad src filepath",
			noSrcFile:   true,
			srcFile:     "profile.test",
			expectErr:   true,
		},
		testcase{
			description: "bad src file",
			srcFile:     "profile.test",
			expectErr:   true,
			srcFileContent: `
meow
`,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			srcFilePath := tc.srcFile
			if !tc.noSrcFile {
				file, err := ioutil.TempFile("", tc.srcFile)
				if err != nil {
					t.Errorf("could not create temp file")
					t.FailNow()
				}

				srcContent := tc.srcFileContent
				err = ioutil.WriteFile(file.Name(), []byte(srcContent), 0644)

				if err != nil {
					t.Errorf("could not write to temp file %v", err)
					t.FailNow()
				}

				srcFilePath = file.Name()
			}
			fset := token.NewFileSet()

			res, err := goparser.NodeFromFilePath(srcFilePath, fset)

			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
				g.Expect(res).To(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(res).ToNot(BeNil())
			}
		})
	}
}

func Test_Parser_RecordFunctionCoverage(t *testing.T) {
	type testcase struct {
		description     string
		functions       []functions.Function
		expectCoverages []FunctionCoverage
		profile         *cover.Profile
	}

	profile := &cover.Profile{
		Blocks: []cover.ProfileBlock{
			cover.ProfileBlock{
				StartLine: 3,
				StartCol:  1,
				EndLine:   5,
				EndCol:    1,
				Count:     4,
				NumStmt:   10,
			},
		},
	}

	testCases := []testcase{
		{
			description: "empty list of functions, blank file path, and nil profile",
		},
		{
			description: "list of functions, blank file path, and nil profile",
			functions: []functions.Function{
				{},
			},
			expectCoverages: []FunctionCoverage{
				FunctionCoverage{
					Function: functions.Function{},
				},
			},
		},
		{
			description: "list of functions, blank file path, and non-matching profile",
			profile:     &cover.Profile{},
			functions: []functions.Function{
				{},
			},
			expectCoverages: []FunctionCoverage{
				FunctionCoverage{
					Function: functions.Function{},
					Profile:  &cover.Profile{},
				},
			},
		},
		{
			description: "list of functions, blank file path, and matching profile",
			profile:     profile,
			functions: []functions.Function{
				{
					StartLine: 2,
					StartCol:  1,
					EndLine:   10,
					EndCol:    1,
				},
			},
			expectCoverages: []FunctionCoverage{
				FunctionCoverage{
					StatementCount: 10,
					CoveredCount:   10,
					Function: functions.Function{
						StartLine: 2,
						StartCol:  1,
						EndLine:   10,
						EndCol:    1,
					},
					Profile: profile,
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			fset := token.NewFileSet()
			p := Parser{Fset: fset, Profile: tc.profile}
			coverages := p.RecordFunctionCoverage(tc.functions)
			g.Expect(coverages).To(HaveLen(len(tc.functions)))
			g.Expect(coverages).To(ConsistOf(tc.expectCoverages))
		})
	}
}

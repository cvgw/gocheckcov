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

package analyzer

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cvgw/gocheckcov/pkg/coverage/parser/profile"
	. "github.com/onsi/gomega"
)

func Test_NewPackageCoverages(t *testing.T) {
	type testcase struct {
		pkgFuncMap map[string][]profile.FunctionCoverage
	}

	testCases := map[string]testcase{
		"one pkg no coverage data": testcase{
			pkgFuncMap: map[string][]profile.FunctionCoverage{
				"github.com/foo/bar/pkg/baz": []profile.FunctionCoverage{},
			},
		},
		"one pkg with coverage data": testcase{
			pkgFuncMap: map[string][]profile.FunctionCoverage{
				"github.com/foo/bar/pkg/baz": []profile.FunctionCoverage{
					profile.FunctionCoverage{
						StatementCount: 10,
						CoveredCount:   10,
					},
				},
			},
		},
	}

	for desc := range testCases {
		desc := desc
		t.Run(desc, func(t *testing.T) {
			g := NewGomegaWithT(t)
			tc := testCases[desc]

			p := NewPackageCoverages(tc.pkgFuncMap)
			g.Expect(p).ToNot(BeNil())
		})
	}
}

func Test_PackageCoverages_Coverage(t *testing.T) {
	g := NewGomegaWithT(t)

	pkgToFuncs := map[string][]profile.FunctionCoverage{
		"github.com/foo/bar/pkg/baz": []profile.FunctionCoverage{},
	}

	p := NewPackageCoverages(pkgToFuncs)
	cov, ok := p.Coverage("github.com/foo/bar/pkg/baz")
	g.Expect(ok).To(BeTrue())
	g.Expect(cov.CoveragePercent).To(Equal(float64(100)))
}

func Test_MapPackagesToFunctions(t *testing.T) {
	type testcase struct {
		srcPath   string
		covPath   string
		dir       string
		expectErr bool
	}

	testCases := map[string]testcase{
		"valid profile": func() testcase {
			dir, err := ioutil.TempDir("", "test")
			if err != nil {
				t.Errorf("could not create temp dir")
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
			goSrcPath := filepath.Join(dir, "src.go")

			err = ioutil.WriteFile(goSrcPath, []byte(profileFileContent), 0644)
			if err != nil {
				t.Errorf("could not write to temp file %v", err)
				t.FailNow()
			}

			profilePath := filepath.Join(dir, "profile.out")

			coverageContent := fmt.Sprintf(
				"mode: set\ngithub.com/cvgw/cov-analyzer/pkg/" +
					"coverage/config/config.go:21.66,22.31 1 1")

			err = ioutil.WriteFile(profilePath, []byte(coverageContent), 0644)

			if err != nil {
				t.Errorf("could not write to temp file %v", err)
				t.FailNow()
			}

			return testcase{
				srcPath: goSrcPath,
				covPath: profilePath,
				dir:     dir,
			}
		}(),
		"invalid profile": func() testcase {
			dir, err := ioutil.TempDir("", "test")
			if err != nil {
				t.Errorf("could not create temp dir")
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
			goSrcPath := filepath.Join(dir, "src.go")

			err = ioutil.WriteFile(goSrcPath, []byte(profileFileContent), 0644)
			if err != nil {
				t.Errorf("could not write to temp file %v", err)
				t.FailNow()
			}

			return testcase{
				srcPath:   goSrcPath,
				covPath:   "foo.out",
				dir:       dir,
				expectErr: true,
			}
		}(),
	}

	for desc := range testCases {
		desc := desc
		t.Run(desc, func(t *testing.T) {
			g := NewGomegaWithT(t)

			tc := testCases[desc]

			fset := token.NewFileSet()

			res, err := MapPackagesToFunctions(tc.covPath, []string{tc.srcPath}, fset, "")
			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(res).ToNot(BeNil())
				g.Expect(res).To(HaveLen(1))
				g.Expect(res).To(HaveKey(strings.TrimPrefix(tc.dir, "/")))
			}
		})
	}
}

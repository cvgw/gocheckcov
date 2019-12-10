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

package config

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_GetConfigFile(t *testing.T) {
	type testcase struct {
		path      string
		expectErr bool
		content   []byte
	}

	testCases := map[string]testcase{
		"blank config path": func() testcase {
			tc := testcase{}
			return tc
		}(),
		"bad config path": func() testcase {
			tc := testcase{
				path:      "foo.yml",
				expectErr: true,
			}

			return tc
		}(),
		"bad config file": func() testcase {
			fi, err := ioutil.TempFile("", "foo.yml")
			if err != nil {
				t.Errorf("could not create tempfile %v", err)
				t.FailNow()
			}

			content := []byte("meow")

			if err := ioutil.WriteFile(fi.Name(), content, 0777); err != nil {
				t.Errorf("could not write to tempfile %v", err)
				t.FailNow()
			}

			tc := testcase{
				path:    fi.Name(),
				content: content,
			}

			return tc
		}(),
		"good config file": func() testcase {
			fi, err := ioutil.TempFile("", "foo.yml")
			if err != nil {
				t.Errorf("could not create tempfile %v", err)
				t.FailNow()
			}

			content := []byte(`
min_coverage_percentage: 10
`)
			if err := ioutil.WriteFile(
				fi.Name(),
				content,
				0777,
			); err != nil {
				t.Errorf("could not write to tempfile %v", err)
				t.FailNow()
			}

			tc := testcase{
				path:    fi.Name(),
				content: content,
			}

			return tc
		}(),
	}

	for desc := range testCases {
		desc := desc
		t.Run(desc, func(t *testing.T) {
			g := NewGomegaWithT(t)
			tc := testCases[desc]
			content, err := GetConfigFile(tc.path)
			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(content).To(Equal(tc.content))
			}
		})
	}
}

func Test_ConfigFile_GetPackage(t *testing.T) {
	g := NewGomegaWithT(t)

	pkgs := []ConfigPackage{
		{Name: "github.com/foo/bar/pkg/baz",
			MinCoveragePercentage: 22,
		},
	}
	c := ConfigFile{
		Packages: pkgs,
	}

	pkg, ok := c.GetPackage("github.com/foo/bar/pkg/baz")
	g.Expect(ok).To(BeTrue())
	g.Expect(pkg).To(Equal(pkgs[0]))
}

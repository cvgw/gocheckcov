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

package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_SetSrcPath(t *testing.T) {
	type testcase struct {
		args        []string
		expected    string
		description string
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get working directory %v", err)
		t.FailNow()
	}

	testCases := []testcase{
		testcase{
			description: "relative path with ...",
			args:        []string{"./pkg/..."},
			expected:    filepath.Join(cwd, "pkg", "..."),
		},
		testcase{
			description: "relative path",
			args:        []string{"./pkg"},
			expected:    filepath.Join(cwd, "pkg"),
		},
		testcase{
			description: "relative path",
			args:        []string{},
			expected:    filepath.Join(cwd, "..."),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.description, func(t *testing.T) {
			g := NewGomegaWithT(t)

			actual := SetSrcPath(tc.args)
			g.Expect(actual).To(Equal(tc.expected))
		})
	}
}

func Test_FilesForPath(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Errorf("couldnt create temp dir %v", err)
		t.FailNow()
	}

	err = os.Mkdir(filepath.Join(dir, "meow"), 0777)
	if err != nil {
		t.Errorf("could not create temp dir %v", err)
		t.FailNow()
	}

	err = ioutil.WriteFile(filepath.Join(dir, "meow", "foo.go"), []byte(`meow`), 0644)
	if err != nil {
		t.Errorf("could not create temp file %v", err)
		t.FailNow()
	}

	err = ioutil.WriteFile(filepath.Join(dir, "bar.go"), []byte(`meow`), 0644)
	if err != nil {
		t.Errorf("could not create temp file %v", err)
		t.FailNow()
	}

	type testcase struct {
		description       string
		expectErr         bool
		dir               string
		ignoreDirs        []string
		expectedFileCount int
	}

	testCases := []testcase{
		{
			description: "bad directory path",
			dir:         "foobar",
			expectErr:   true,
		},
		{
			description: "bad directory path, is file",
			dir:         filepath.Join(dir, "bar.go"),
			expectErr:   true,
		},
		{
			description:       "valid directory path, no recursion",
			dir:               dir,
			expectedFileCount: 1,
		},
		{
			description:       "valid directory path, with recursion",
			dir:               filepath.Join(dir, "..."),
			expectedFileCount: 2,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			files, err := FilesForPath(tc.dir, tc.ignoreDirs)
			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(files).To(HaveLen(tc.expectedFileCount))
			}
		})
	}
}

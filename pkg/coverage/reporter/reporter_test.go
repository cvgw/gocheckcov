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

package reporter

import (
	"testing"

	"github.com/cvgw/gocheckcov/mocks/coverage/mock_reporter"
	"github.com/cvgw/gocheckcov/pkg/coverage/analyzer"
	"github.com/cvgw/gocheckcov/pkg/coverage/config"
	"github.com/cvgw/gocheckcov/pkg/coverage/parser/profile"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func Test_Verifier_VerifyCoverage(t *testing.T) {
	g := NewGomegaWithT(t)

	type testcase struct {
		verifier  *Verifier
		pkg       config.ConfigPackage
		coverages *analyzer.PackageCoverages
		result    bool
		expectErr bool
	}

	type tcFn func(*gomock.Controller) testcase

	testCases := map[string]tcFn{
		"empty package and coverages": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)

			return testcase{
				verifier:  &Verifier{Out: mockLogger},
				coverages: &analyzer.PackageCoverages{},
				pkg:       config.ConfigPackage{},
				expectErr: true,
			}
		},
		"cov is less than pkg min": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)
			mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).Times(1)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				coverages: analyzer.NewPackageCoverages(map[string][]profile.FunctionCoverage{
					"foo/bar": []profile.FunctionCoverage{
						{
							CoveredCount:   10,
							StatementCount: 100,
						},
					},
				}),
				pkg: config.ConfigPackage{
					Name:                  "foo/bar",
					MinCoveragePercentage: 100,
				},
			}
		},
	}

	for i := range testCases {
		i := i
		t.Run(i, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := testCases[i](ctrl)

			v := tc.verifier
			pkg := tc.pkg
			ok, err := v.VerifyCoverage(pkg, tc.coverages)
			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(ok).To(Equal(tc.result))
			}
		})
	}
}

func Test_Verifier_ReportCoverage(t *testing.T) {
	type testcase struct {
		verifier   *Verifier
		input      map[string][]profile.FunctionCoverage
		configData []byte
		printFuncs bool
		expectErr  bool
	}

	type tcFn func(*gomock.Controller) testcase

	testCases := map[string]tcFn{
		"empty function map and nil coverage data": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				input:    map[string][]profile.FunctionCoverage{},
			}
		},
		"one empty function": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)

			mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).MinTimes(1)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				input: map[string][]profile.FunctionCoverage{
					"baz": []profile.FunctionCoverage{
						profile.FunctionCoverage{},
					},
				},
			}
		},
		"one function with one statement not in coverage data": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)

			mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).MinTimes(1)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				input: map[string][]profile.FunctionCoverage{
					"foo/bar": []profile.FunctionCoverage{
						{CoveredCount: 1, StatementCount: 1},
					},
				},
			}
		},
		"one function with one statement in coverage data": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)

			mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).MinTimes(1)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				input: map[string][]profile.FunctionCoverage{
					"foo/bar": []profile.FunctionCoverage{
						{CoveredCount: 1, StatementCount: 1},
					},
				},
				configData: []byte(`
packages:
- name: foo/bar
  min_coverage_percentage: 10
`),
			}
		},
		"one function does not meet min coverage for package": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)

			mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).MinTimes(1)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				input: map[string][]profile.FunctionCoverage{
					"foo/bar": []profile.FunctionCoverage{
						{CoveredCount: 0, StatementCount: 1},
					},
				},
				expectErr: true,
				configData: []byte(`
min_coverage_percentage: 0
packages:
- name: foo/bar
  min_coverage_percentage: 10
`),
			}
		},
		"one function does not meet global min coverage": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)

			mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).MinTimes(1)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				input: map[string][]profile.FunctionCoverage{
					"foo/bar": []profile.FunctionCoverage{
						{CoveredCount: 0, StatementCount: 1},
					},
				},
				expectErr: true,
				configData: []byte(`
min_coverage_percentage: 20
packages:
- name: baz
  min_coverage_percentage: 0
`),
			}
		},
		"one function with one statement with a bad config file": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				input: map[string][]profile.FunctionCoverage{
					"foo/bar": []profile.FunctionCoverage{
						{CoveredCount: 1, StatementCount: 1},
					},
				},
				expectErr:  true,
				configData: []byte("meow"),
			}
		},
	}

	for description := range testCases {
		description := description

		t.Run(description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := testCases[description](ctrl)
			v := tc.verifier
			_, err := v.ReportCoverage(tc.input, tc.printFuncs, tc.configData)
			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func Test_Verifier_PrintReport(t *testing.T) {
	type testcase struct {
		verifier  *Verifier
		functions []profile.FunctionCoverage
	}

	type tcFn func(*gomock.Controller) testcase

	testCases := map[string]tcFn{
		"empty function list": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)
			mockLogger.EXPECT().Printf(gomock.Any()).Times(1)
			return testcase{
				verifier:  &Verifier{Out: mockLogger},
				functions: []profile.FunctionCoverage{},
			}
		},
		"one empty function": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)
			mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).MinTimes(1)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				functions: []profile.FunctionCoverage{
					profile.FunctionCoverage{},
				},
			}
		},
		"one function with one statement": func(ctrl *gomock.Controller) testcase {
			mockLogger := mock_reporter.NewMocklogger(ctrl)
			mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).MinTimes(1)

			return testcase{
				verifier: &Verifier{Out: mockLogger},
				functions: []profile.FunctionCoverage{
					{CoveredCount: 1, StatementCount: 1},
				},
			}
		},
	}

	for description := range testCases {
		description := description

		t.Run(description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := testCases[description](ctrl)
			v := tc.verifier
			v.PrintFunctionReport(tc.functions)
		})
	}
}

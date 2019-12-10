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
	"math"
	"path/filepath"
	"strings"

	"github.com/cvgw/gocheckcov/pkg/coverage/parser/goparser"
	"github.com/cvgw/gocheckcov/pkg/coverage/parser/goparser/functions"
	"github.com/cvgw/gocheckcov/pkg/coverage/parser/profile"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/cover"
)

type PackageCoverages struct {
	coverages map[string]coverage
}

type coverage struct {
	StatementCount  int64
	ExecutedCount   int64
	CoveragePercent float64
	Functions       []profile.FunctionCoverage
}

func (p *PackageCoverages) Coverage(pkg string) (coverage, bool) {
	cov, ok := p.coverages[pkg]
	return cov, ok
}

func NewPackageCoverages(packagesToFunctions map[string][]profile.FunctionCoverage) *PackageCoverages {
	pkgToCoverage := make(map[string]coverage)

	for pkg, functions := range packagesToFunctions {
		var statementCount int64

		var executedCount int64

		for _, function := range functions {
			statementCount += function.StatementCount
			executedCount += function.CoveredCount
		}

		var covPer float64

		if executedCount == 0 && statementCount == 0 {
			covPer = 100
		} else {
			covPer = math.Floor((float64(executedCount)/float64(statementCount))*10000) / 100
		}

		c := coverage{
			StatementCount:  statementCount,
			ExecutedCount:   executedCount,
			CoveragePercent: covPer,
			Functions:       functions,
		}
		pkgToCoverage[pkg] = c
	}

	return &PackageCoverages{
		coverages: pkgToCoverage,
	}
}

func MapPackagesToFunctions(
	filePath string,
	projectFiles []string,
	fset *token.FileSet,
	goSrc string,
) (map[string][]profile.FunctionCoverage, error) {
	profiles, err := cover.ParseProfiles(filePath)
	if err != nil {
		e := fmt.Errorf("could not parse profiles from %v %v", filePath, err)
		return nil, e
	}

	filePathToProfileMap := make(map[string]*cover.Profile)
	for _, prof := range profiles {
		filePathToProfileMap[prof.FileName] = prof
	}

	packageToFunctions := make(map[string][]profile.FunctionCoverage)

	for _, filePath := range projectFiles {
		node, err := goparser.NodeFromFilePath(filePath, goSrc, fset)
		if err != nil {
			e := fmt.Errorf("could not retrieve node from filepath %v", err)
			return nil, e
		}

		functions, err := functions.CollectFunctions(node, fset, filePath)
		if err != nil {
			e := fmt.Errorf("could not collect functions for filepath %v %v", filePath, err)
			return nil, e
		}

		log.Debugf("functions for file %v %v", filePath, functions)
		pkg := strings.TrimPrefix(filePath, fmt.Sprintf("%s/", goSrc))
		pkg = filepath.Dir(pkg)

		var funcCoverages []profile.FunctionCoverage

		p := profile.Parser{FilePath: filePath, Fset: fset}

		if prof, ok := filePathToProfileMap[filePath]; ok {
			p.Profile = prof
		}

		funcCoverages = p.RecordFunctionCoverage(functions)

		packageToFunctions[pkg] = append(packageToFunctions[pkg], funcCoverages...)
	}

	log.Debugf("map of packages to functions %v", packageToFunctions)

	return packageToFunctions, nil
}

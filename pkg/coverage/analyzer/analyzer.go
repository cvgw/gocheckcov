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
	"context"
	"fmt"
	"go/token"
	"math"
	"path/filepath"

	"github.com/cvgw/gocheckcov/pkg/coverage/parser/goparser"
	"github.com/cvgw/gocheckcov/pkg/coverage/parser/goparser/functions"
	"github.com/cvgw/gocheckcov/pkg/coverage/parser/profile"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/cover"
	"golang.org/x/tools/go/packages"
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

type packageList struct {
	cache map[string]*pkg
}

type pkg struct {
	*packages.Package
}

func (p *pkg) Path() string {
	return p.PkgPath
}

func (p *packageList) get(dirPath string) (*pkg, error) {
	if p.cache == nil {
		p.cache = make(map[string]*pkg)
	}

	bPkg, ok := p.cache[dirPath]
	if !ok {
		conf := &packages.Config{
			Mode:    packages.NeedName,
			Tests:   false,
			Context: context.Background(),
		}

		pkgs, err := packages.Load(conf, dirPath)
		if err != nil {
			return nil, err
		}

		if len(pkgs) != 1 {
			return nil, fmt.Errorf("expected 1 pkg")
		}

		pk := pkgs[0]
		bPkg = &pkg{pk}
		p.cache[dirPath] = bPkg
	}

	return bPkg, nil
}

func MapPackagesToFunctions(
	filePath string,
	projectFiles []string,
	fset *token.FileSet,
) (map[string][]profile.FunctionCoverage, error) {
	profiles, err := cover.ParseProfiles(filePath)
	if err != nil {
		e := fmt.Errorf("could not parse profiles from %v %v", filePath, err)
		return nil, e
	}

	filePathToProfileMap := make(map[string]*cover.Profile)

	for _, prof := range profiles {
		pPath := prof.FileName
		log.Debugf("adding profile with file name %v to map", pPath)

		filePathToProfileMap[pPath] = prof
	}

	packageToFunctions := make(map[string][]profile.FunctionCoverage)

	packageList := &packageList{}

	for _, filePath := range projectFiles {
		node, err := goparser.NodeFromFilePath(filePath, fset)
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

		pkg, err := packageList.get(filepath.Dir(filePath))
		if err != nil {
			return nil, err
		}

		log.Debugf("found pkg %v for filepath %v", pkg.Path(), filePath)

		var funcCoverages []profile.FunctionCoverage

		p := profile.Parser{FilePath: filePath, Fset: fset}

		profilePath := fmt.Sprintf("%v/%v", pkg.Path(), filepath.Base(filePath))
		if prof, ok := filePathToProfileMap[profilePath]; ok {
			p.Profile = prof
		} else {
			log.Debugf("no profile found for path %v", profilePath)
		}

		funcCoverages = p.RecordFunctionCoverage(functions)

		packageToFunctions[pkg.Path()] = append(packageToFunctions[pkg.Path()], funcCoverages...)
	}

	log.Debugf("map of packages to functions %v", packageToFunctions)

	return packageToFunctions, nil
}

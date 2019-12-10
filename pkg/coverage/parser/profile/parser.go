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

	"github.com/cvgw/gocheckcov/pkg/coverage/parser/goparser/functions"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/cover"
)

type FunctionCoverage struct {
	StatementCount int64
	CoveredCount   int64
	Name           string
	Function       functions.Function
	Profile        *cover.Profile
}

type Parser struct {
	Fset     *token.FileSet
	FilePath string
	Profile  *cover.Profile
}

func (p Parser) RecordFunctionCoverage(functions []functions.Function) []FunctionCoverage {
	out := make([]FunctionCoverage, 0, len(functions))

	for _, function := range functions {
		fc := FunctionCoverage{
			Name:     function.Name,
			Function: function,
		}

		if p.Profile != nil {
			fc = p.recordCoverageHits(fc, function)
			fc.Profile = p.Profile
		}

		if int(fc.StatementCount) != len(function.Statements) {
			log.Debugf(
				"function %v statement counts don't match Profile: %v AST: %v",
				function.Name,
				fc.StatementCount,
				len(function.Statements),
			)

			if int(fc.StatementCount) == 0 && len(function.Statements) > 0 {
				fc.StatementCount = int64(len(function.Statements))
			}
		}

		out = append(out, fc)
	}

	return out
}

func (p Parser) recordCoverageHits(fc FunctionCoverage, function functions.Function) FunctionCoverage {
	for _, block := range p.Profile.Blocks {
		startLine := function.StartLine
		startCol := function.StartCol
		endLine := function.EndLine
		endCol := function.EndCol

		if block.StartLine > endLine || (block.StartLine == endLine && block.StartCol >= endCol) {
			// Block starts after the function statement ends
			continue
		}

		if block.EndLine < startLine || (block.EndLine == startLine && block.EndCol <= startCol) {
			// Block ends before the function statement starts
			continue
		}

		fc.StatementCount += int64(block.NumStmt)
		if block.Count > 0 {
			fc.CoveredCount += int64(block.NumStmt)
		}
	}

	return fc
}

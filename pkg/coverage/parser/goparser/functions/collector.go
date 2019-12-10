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
	"go/ast"
	"go/token"

	"github.com/cvgw/gocheckcov/pkg/coverage/parser/goparser/statements"
	log "github.com/sirupsen/logrus"
)

func CollectFunctions(f *ast.File, fset *token.FileSet, filePath string) ([]Function, error) {
	functions := []Function{}

	for i := range f.Decls {
		switch x := f.Decls[i].(type) {
		case *ast.FuncDecl:
			name := x.Name.Name

			start := fset.Position(x.Pos())
			end := fset.Position(x.End())
			startLine := start.Line
			startCol := start.Column
			endLine := end.Line
			endCol := end.Column
			f := Function{
				Name:        name,
				StartLine:   startLine,
				StartCol:    startCol,
				EndLine:     endLine,
				EndCol:      endCol,
				SrcPath:     filePath,
				StartOffset: start.Offset,
				EndOffset:   end.Offset,
			}

			sc := &statements.StmtCollector{}
			if err := sc.Collect(x.Body, fset); err != nil {
				return nil, err
			}

			stmts := sc.Statements
			log.Debugf("statements for function %v %v", f.Name, stmts)
			convertedStmts := make([]statements.Statement, 0, len(stmts))

			for _, stmnt := range stmts {
				start := fset.Position(stmnt.Pos())
				end := fset.Position(stmnt.End())
				startLine := start.Line
				startCol := start.Column
				endLine := end.Line
				endCol := end.Column
				s := statements.Statement{
					StartLine: int64(startLine),
					StartCol:  int64(startCol),
					EndLine:   int64(endLine),
					EndCol:    int64(endCol),
				}
				convertedStmts = append(convertedStmts, s)
			}

			f.Statements = convertedStmts
			functions = append(functions, f)
		}
	}

	log.Debugf("found functions %v", functions)

	return functions, nil
}

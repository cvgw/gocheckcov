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

package statements

import (
	"go/ast"
	"go/token"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_CollectStatements(t *testing.T) {
	type testcase struct {
		description       string
		stmt              ast.Stmt
		expectErr         bool
		expectedStmtCount int
	}

	var nilBlock *ast.BlockStmt

	testCases := []testcase{
		{
			description: "nil block statement",
			stmt:        nilBlock,
			expectErr:   true,
		},
		{
			description: "empty block statement",
			stmt:        &ast.BlockStmt{},
		},
		{
			description: "block statement with a bad switch",
			stmt: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.SwitchStmt{},
				},
			},
			expectErr: true,
		},
		{
			description: "block statement with list containing empty stmts",
			stmt: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.BlockStmt{},
				},
			},
		},
		{
			description: "empty case clause",
			stmt:        &ast.CaseClause{},
		},
		{
			description: "empty comm clause",
			stmt:        &ast.CommClause{},
		},
		{
			description: "some other kind of empty stmt",
			stmt:        &ast.AssignStmt{},
		},
		{
			description: "switch statement with bad body",
			stmt:        &ast.SwitchStmt{},
			expectErr:   true,
		},
		{
			description: "empty for statement",
			stmt:        &ast.ForStmt{},
			expectErr:   true,
		},
		{
			description: "empty if statement",
			stmt:        &ast.IfStmt{},
			expectErr:   true,
		},
		{
			description: "if statement with non empty body and bad else statement",
			stmt: &ast.IfStmt{
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.BlockStmt{},
					},
				},
				Else: &ast.AssignStmt{},
			},
			expectErr: true,
		},
		{
			description: "if statement with non empty body and if stmt else statement not match",
			stmt: &ast.IfStmt{
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.BlockStmt{},
					},
				},
				Else: &ast.IfStmt{Body: &ast.BlockStmt{}},
			},
			expectedStmtCount: 1,
		},
		{
			description: "if statement with non empty body and block else statement not match",
			stmt: &ast.IfStmt{
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.BlockStmt{},
					},
				},
				Else: &ast.BlockStmt{},
			},
		},
		{
			description: "empty labeled statement",
			stmt:        &ast.LabeledStmt{},
		},
		{
			description: "empty range statement",
			stmt:        &ast.RangeStmt{},
			expectErr:   true,
		},
		{
			description: "empty select statement",
			stmt:        &ast.SelectStmt{},
			expectErr:   true,
		},
		{
			description: "empty switch statement",
			stmt:        &ast.SwitchStmt{},
			expectErr:   true,
		},
		{
			description: "switch statement with bad init",
			stmt: &ast.SwitchStmt{
				Init: &ast.IfStmt{},
			},
			expectErr: true,
		},
		{
			description: "switch statement with empty body",
			stmt: &ast.SwitchStmt{
				Body: &ast.BlockStmt{},
			},
		},
		{
			description: "bad type switch",
			stmt:        &ast.TypeSwitchStmt{},
			expectErr:   true,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			fset := token.NewFileSet()
			s := &StmtCollector{}
			err := s.Collect(tc.stmt, fset)
			stmts := s.Statements
			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(stmts).ToNot(BeNil())
				g.Expect(stmts).To(HaveLen(tc.expectedStmtCount))
			}
		})
	}
}

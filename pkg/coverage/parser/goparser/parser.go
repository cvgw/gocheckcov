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

package goparser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func NodeFromFilePath(filePath, goSrcPath string, fset *token.FileSet) (*ast.File, error) {
	pFilePath := filepath.Join(goSrcPath, filePath)

	src, err := ioutil.ReadFile(pFilePath)
	if err != nil {
		log.Debugf("could not read file from profile %v %v", pFilePath, err)
		return nil, err
	}

	f, err := parser.ParseFile(fset, pFilePath, src, 0)
	if err != nil {
		log.Debugf("could not parse file %v %v", pFilePath, err)
		return nil, err
	}

	return f, nil
}

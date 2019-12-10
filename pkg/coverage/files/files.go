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
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func SetSrcPath(args []string) string {
	var srcPath string
	if len(args) > 0 {
		srcPath = args[0]
		log.Debugf("srcPath %v", srcPath)
		absSrcPath, err := filepath.Abs(srcPath)

		if err != nil {
			log.Debugf("could not get absolute path from %v %v", srcPath, err)
		} else {
			log.Debugf("absSrcPath %v", absSrcPath)
			srcPath = absSrcPath
		}
	}

	if srcPath == "" {
		var err error
		srcPath, err = os.Getwd()

		if err != nil {
			log.Debugf("could not get working directory %v", err)
		}

		srcPath = filepath.Join(srcPath, "...")
	}

	return srcPath
}

type dirsToIgnore []string

func (d dirsToIgnore) Includes(dir string) bool {
	for _, ignore := range d {
		if ignore == dir {
			return true
		}
	}

	return false
}

func FilesForPath(dir string, ignoreDirs dirsToIgnore) ([]string, error) {
	base := filepath.Base(dir)
	if base == "..." {
		dir = filepath.Dir(dir)
	}

	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("path must be a directory")
	}

	if base == "..." {
		return recusiveFilesForPath(dir, ignoreDirs)
	}

	return filesForDir(dir)
}

func recusiveFilesForPath(dir string, ignoreDirs dirsToIgnore) ([]string, error) {
	goPath := build.Default.GOPATH
	files := make([]string, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("could not access path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() && ignoreDirs.Includes(info.Name()) {
			return filepath.SkipDir
		}

		if info.Mode().IsRegular() {
			if regexp.MustCompile(".go$").Match([]byte(path)) {
				if regexp.MustCompile("_test.go$").Match([]byte(path)) {
					return nil
				}
				path = strings.TrimPrefix(path, fmt.Sprintf("%v/", filepath.Join(goPath, "src")))
				files = append(files, path)
			}
		}
		return nil
	})

	log.Debugf("files for %v %v", dir, files)

	return files, err
}

func filesForDir(dir string) ([]string, error) {
	goPath := build.Default.GOPATH

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)

	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}

		if fi.Mode().IsRegular() {
			if regexp.MustCompile(".go$").Match([]byte(fi.Name())) {
				if regexp.MustCompile("_test.go$").Match([]byte(fi.Name())) {
					continue
				}

				path := filepath.Join(dir, fi.Name())
				path = strings.TrimPrefix(path, fmt.Sprintf("%v/", filepath.Join(goPath, "src")))
				files = append(files, path)
			}
		}
	}

	log.Debugf("files for %v %v", dir, files)

	return files, err
}

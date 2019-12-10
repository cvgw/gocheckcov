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

package cmd

import (
	"go/build"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cvgw/gocheckcov/pkg/coverage/analyzer"
	"github.com/cvgw/gocheckcov/pkg/coverage/config"
	"github.com/cvgw/gocheckcov/pkg/coverage/files"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// checkInitCmd represents the checkInit command
var checkInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new config file using current coverage",
	Long: `Create a new configuration file for gocheckcov which lists all of packages in the specified path and sets ` +
		`the minimum converage percentage for each to the current coverage percentage for that package`,
	Run: func(cmd *cobra.Command, args []string) {
		profilePath := ProfileFile
		fset := token.NewFileSet()

		srcPath := files.SetSrcPath(args)
		dir := srcPath
		ignoreDirs := strings.Split(skipDirs, ",")
		projectFiles, err := files.FilesForPath(dir, ignoreDirs)
		if err != nil {
			log.Printf("could not retrieve files for path %v %v", dir, err)
			os.Exit(1)
		}

		goSrc := filepath.Join(build.Default.GOPATH, "src")

		packageToFunctions, err := analyzer.MapPackagesToFunctions(profilePath, projectFiles, fset, goSrc)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		pc := analyzer.NewPackageCoverages(packageToFunctions)

		configFile := config.ConfigFile{}
		for pkg := range packageToFunctions {
			cov, ok := pc.Coverage(pkg)
			if !ok {
				log.Printf("could not get coverage for package %v", pkg)
				os.Exit(1)
			}
			cfgPkg := config.ConfigPackage{MinCoveragePercentage: cov.CoveragePercent, Name: pkg}
			configFile.Packages = append(configFile.Packages, cfgPkg)
		}

		configContent, err := yaml.Marshal(configFile)
		if err != nil {
			log.Printf("couldn't marshal config file %v", err)
			os.Exit(1)
		}
		if err := ioutil.WriteFile("config.yaml", configContent, 0644); err != nil {
			log.Printf("could not read config file %v %v", configFile, err)
			os.Exit(1)
		}
	},
}

func init() {
	checkCmd.AddCommand(checkInitCmd)

	checkInitCmd.Flags().StringVarP(&ProfileFile, "profile-file", "p", "", "path to coverage profile file")

	if err := checkInitCmd.MarkFlagRequired("profile-file"); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

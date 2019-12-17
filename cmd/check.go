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
	"bufio"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/cvgw/gocheckcov/pkg/coverage/analyzer"
	"github.com/cvgw/gocheckcov/pkg/coverage/config"
	"github.com/cvgw/gocheckcov/pkg/coverage/files"
	"github.com/cvgw/gocheckcov/pkg/coverage/reporter"
	"github.com/spf13/cobra"
)

var (
	noConfig       bool
	configFile     string
	ProfileFile    string
	printFunctions bool
	printSrc       bool
	minCov         float64
	skipDirs       string
	checkCmd       = &cobra.Command{
		Use:   "check",
		Short: "Check whether pkg coverage meets specified minimum",
		Run: func(cmd *cobra.Command, args []string) {
			err := runCheckCommand(args)
			if err != nil {
				os.Exit(1)
			}
		},
	}
)

func runCheckCommand(args []string) error {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	ignoreDirs := strings.Split(skipDirs, ",")
	srcPath := files.SetSrcPath(args)
	dir := srcPath

	projectFiles, err := files.FilesForPath(dir, ignoreDirs)
	if err != nil {
		log.Printf("could not retrieve project files from path %v %v", dir, err)
		return err
	}

	profilePath := ProfileFile
	if profilePath == "" {
		pf, e := runTestsAndGenerateProfile(srcPath)
		if e != nil {
			return fmt.Errorf("could not run tests %v", e)
		}

		profilePath = pf.Name()

		defer func() {
			if e := os.Remove(pf.Name()); e != nil {
				log.Print(e)
			}
		}()
	}

	fset := token.NewFileSet()

	packageToFunctions, err := analyzer.MapPackagesToFunctions(profilePath, projectFiles, fset)
	if err != nil {
		log.Print(err)
		return err
	}

	cfContent, err := getConfig()
	if err != nil {
		return err
	}

	cliL := reporter.NewCliTabLogger()
	defer cliL.Close()

	if printSrc {
		printFunctions = true
	}

	v := reporter.Verifier{
		Out:            cliL,
		PrintFunctions: printFunctions,
		PrintSrc:       printSrc,
		MinCov:         minCov,
	}

	if _, err := v.ReportCoverage(packageToFunctions, printFunctions, cfContent); err != nil {
		cliL.Printf("%v", err)
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolVar(&printFunctions, "print-functions", false, "print coverage for individual functions")

	checkCmd.Flags().BoolVar(
		&printSrc,
		"print-src",
		false,
		"print src coverage for each function (print-functions automatically set to true)",
	)

	checkCmd.Flags().BoolVar(&noConfig, "no-config", false, "do not read configuration from file")

	checkCmd.Flags().Float64VarP(
		&minCov,
		"minimum-coverage",
		"m",
		0,
		"minimum coverage percentage to enforce for all packages (defaults to 0)",
	)

	checkCmd.Flags().StringVarP(&ProfileFile, "profile-file", "p", "", "path to coverage profile file")

	checkCmd.Flags().StringVarP(
		&configFile,
		"config-file",
		"c",
		"",
		"path to configuration file",
	)

	checkCmd.PersistentFlags().StringVarP(
		&skipDirs,
		"skip-dirs",
		"s",
		"vendor",
		"command separted list of directories to skip when reporting coverage",
	)
}

func getConfig() ([]byte, error) {
	var cfContent []byte

	var err error

	if !noConfig {
		cfContent, err = config.GetConfigFile(configFile)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
	}

	return cfContent, nil
}

func runTestsAndGenerateProfile(srcPath string) (*os.File, error) {
	f, err := ioutil.TempFile("", "profile.out")
	if err != nil {
		return nil, err
	}

	args := []string{
		"test",
		"-coverprofile=" + f.Name(),
	}

	pkgPath := srcPath
	args = append(args, pkgPath)
	c := exec.Command("go", args...)

	stderr, err := c.StderrPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := c.Start(); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go scan(&wg, stderr)

	go scan(&wg, stdout)

	if err := c.Wait(); err != nil {
		return nil, err
	}

	wg.Wait()

	return f, nil
}

func scan(wg *sync.WaitGroup, r io.ReadCloser) {
	defer wg.Done()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}

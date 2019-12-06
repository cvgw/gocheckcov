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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/fatih/color"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/cover"

	"github.com/cvgw/gocheckcov/pkg/coverage/analyzer"
	"github.com/cvgw/gocheckcov/pkg/coverage/config"
	"github.com/cvgw/gocheckcov/pkg/coverage/parser/profile"
	"gopkg.in/yaml.v2"
)

func NewCliTabLogger() *CliLogger {
	out := bufio.NewWriter(os.Stdout)
	tabber := tabwriter.NewWriter(out, 1, 8, 1, '\t', 0)

	return &CliLogger{
		out:    out,
		tabber: tabber,
	}
}

type CliLogger struct {
	tabber *tabwriter.Writer
	out    *bufio.Writer
}

func (l *CliLogger) Printf(fmtString string, args ...interface{}) {
	if l.tabber == nil {
		fmt.Println(fmt.Sprintf(fmtString, args...))
		return
	}

	if _, err := fmt.Fprintf(l.tabber, fmtString, args...); err != nil {
		panic(fmt.Errorf("could not write to output %v", err))
	}
}

func (l *CliLogger) Close() {
	if err := l.tabber.Flush(); err != nil {
		log.Debug(err)
	}

	if err := l.out.Flush(); err != nil {
		log.Debug(err)
	}
}

type Logger interface {
	Printf(string, ...interface{})
}

type Verifier struct {
	Out            Logger
	MinCov         float64
	PrintSrc       bool
	PrintFunctions bool
}

func (v Verifier) ReportCoverage(
	packageToFunctions map[string][]profile.FunctionCoverage,
	printFunctions bool,
	configFile []byte,
) (map[string]float64, error) {
	pkgToCoverage := make(map[string]float64)
	pc := analyzer.NewPackageCoverages(packageToFunctions)

	fail := false
	keys := make([]string, 0, len(packageToFunctions))

	for pkg := range packageToFunctions {
		keys = append(keys, pkg)
	}

	sort.Strings(keys)

	for _, pkg := range keys {
		var cfgPkg config.ConfigPackage

		if len(configFile) != 0 {
			cfg := config.ConfigFile{}
			if err := yaml.Unmarshal(configFile, &cfg); err != nil {
				err = errors.Wrap(err, "could not unmarshal yaml for config file %v")
				log.Debug(err)

				return nil, err
			}

			var ok bool

			cfgPkg, ok = cfg.GetPackage(pkg)
			if !ok {
				log.Debugf("could not find package for name %v", pkg)

				cfgPkg = config.ConfigPackage{
					Name:                  pkg,
					MinCoveragePercentage: cfg.MinCoveragePercentage,
				}
			}
		} else {
			cfgPkg = config.ConfigPackage{
				Name:                  pkg,
				MinCoveragePercentage: v.MinCov,
			}
		}

		ok, err := v.VerifyCoverage(cfgPkg, pc)
		if err != nil {
			log.Debug(err)
			return nil, err
		}

		if !ok {
			fail = true
		}
	}

	if fail {
		return nil, fmt.Errorf("packages failed to meet minimum coverage")
	}

	return pkgToCoverage, nil
}

func (v Verifier) VerifyCoverage(pkg config.ConfigPackage, pc *analyzer.PackageCoverages) (bool, error) {
	if pc == nil {
		err := fmt.Errorf("can't report coverages because coverage data is nil")
		log.Debug(err)

		return false, err
	}

	cov, ok := pc.Coverage(pkg.Name)

	if !ok {
		err := fmt.Errorf("could not get coverage for package %v", pkg)
		log.Debug(err)

		return false, err
	}

	v.Out.Printf(
		"pkg  %v\tcoverage %v%% \tminimum %v%% \tstatements\t%v/%v\n",
		pkg.Name,
		cov.CoveragePercent,
		pkg.MinCoveragePercentage,
		cov.ExecutedCount,
		cov.StatementCount,
	)

	if v.PrintFunctions {
		if err := v.PrintFunctionReport(cov.Functions); err != nil {
			return false, err
		}
	}

	if pkg.MinCoveragePercentage > cov.CoveragePercent {
		return false, nil
	}

	return true, nil
}

func (v Verifier) PrintFunctionReport(functions []profile.FunctionCoverage) error {
	for _, function := range functions {
		if function.StatementCount == 0 {
			continue
		}

		var executedStatementsCount int64

		executedStatementsCount += function.CoveredCount

		val := (float64(executedStatementsCount) / float64(function.StatementCount)) * 10000
		percent := (math.Floor(val) / 100)
		v.Out.Printf(
			"func %v\tcoverage %v%% \t\tstatements\t%v/%v\n",
			function.Name,
			percent,
			executedStatementsCount,
			function.StatementCount,
		)

		if v.PrintSrc {
			filePath := function.Function.SrcPath
			src, err := ioutil.ReadFile(filePath)

			if err != nil {
				return err
			}

			if err := v.printSrcWithCoverage(function, src); err != nil {
				return err
			}
		}
	}

	v.Out.Printf("\n")

	return nil
}

func (v *Verifier) printSrcWithCoverage(fc profile.FunctionCoverage, src []byte) error {
	boundaries := []cover.Boundary{}
	if fc.Profile != nil {
		boundaries = fc.Profile.Boundaries(src)
	}

	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	wht := color.New(color.FgWhite)
	out := bytes.NewBuffer(make([]byte, 0))
	buf := bytes.NewBuffer(make([]byte, 0))

	clr := wht

	for i := fc.Function.StartOffset - 1; i < fc.Function.EndOffset+1; i++ {
		for _, b := range boundaries {
			if b.Offset != i {
				continue
			}

			if b.Start {
				if _, err := clr.Fprint(out, buf.String()); err != nil {
					return err
				}

				buf = bytes.NewBuffer(make([]byte, 0))

				if b.Norm == 0 {
					clr = red
				} else {
					clr = green
				}
			} else {
				if _, err := clr.Fprint(out, buf.String()); err != nil {
					return err
				}

				buf = bytes.NewBuffer(make([]byte, 0))
				clr = wht
			}
		}

		buf.Write([]byte{src[i]})
	}

	if _, err := clr.Fprint(out, buf.String()); err != nil {
		return err
	}

	v.Out.Printf("%s\n", out.String())

	return nil
}

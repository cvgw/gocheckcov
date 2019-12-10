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

package config

import (
	"io/ioutil"
	"os"
)

const (
	defaultConfigPath = ".gocheckcov-config.yml"
)

func GetConfigFile(configPath string) ([]byte, error) {
	var cfContent []byte

	if configPath == "" {
		configPath = defaultConfigPath
	}

	_, err := os.Stat(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		if configPath != defaultConfigPath {
			return nil, err
		}
	} else {
		cfContent, err = ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
	}

	return cfContent, nil
}

type ConfigFile struct {
	MinCoveragePercentage float64         `yaml:"min_coverage_percentage"`
	Packages              []ConfigPackage `yaml:"packages"`
}

func (c ConfigFile) GetPackage(pkg string) (ConfigPackage, bool) {
	for _, p := range c.Packages {
		if p.Name == pkg {
			return p, true
		}
	}

	return ConfigPackage{}, false
}

type ConfigPackage struct {
	Name                  string  `yaml:"name"`
	MinCoveragePercentage float64 `yaml:"min_coverage_percentage"`
}

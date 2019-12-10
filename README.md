# gocheckcov

![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/cvgw/gocheckcov/Go/master?style=plastic)
[![Coveralls github branch](https://img.shields.io/coveralls/github/cvgw/gocheckcov/master?style=plastic)](https://coveralls.io/github/cvgw/gocheckcov?branch=master)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/cvgw/gocheckcov?style=plastic)
[![license](https://img.shields.io/github/license/cvgw/gocheckcov?style=plastic)](./LICENSE)

[![asciicast](https://asciinema.org/a/GpOwjIQUFxTq2lr9m5HDsQ7E3.svg)](https://asciinema.org/a/GpOwjIQUFxTq2lr9m5HDsQ7E3)

## Description

gocheckcov allows users to assert that a set of golang packages meets a minimum
level of test coverage. Users can specify a minimum coverage percentage for all
packages or specify a minimum for each package via a configuration file.
gocheckcov executes the tests in the given path and generates a coverage
profile. If each package does not meet the specified minimum coverage gocheckcov
will exit with code 1.

```
$ gocheckcov check --minimum-coverage 66.6 $GOPATH/src/github.com/bar/foo/pkg/baz

pkg github.com/bar/foo/pkg/baz coverage is 10% (10/100 statements)
coverage 10% for package github.com/bar/foo/pkg/baz did not meet minimum 66.6%

$ echo $?
1
```

*   [Status](#status)
*   [Getting Started](#getting-started)
*   [Usage](#usage)
    *   [Supported Golang Versions](#supported-golang-versions)
    *   [Install](#install)
    *   [Configuration](#configuration)
*   [Development](#development)
*   [Contributing](#contributing)

## Status

Alpha

**This is not an officially supported Google product**

## Getting Started

Install gocheckcov `$ go get github.com/cvgw/gocheckcov`

Generate a coverage profile.

gocheckcov doesn't have any special requirements for the coverage profile. Just
add the `-coverprofile=${out-file-path}` flag to your normal `go test` command.

Initialize a new config for your project `/path/to/project$ gocheckcov check
init --profile-file ${coverprofile_path}`

Write some new code and generate a new coverage profile

Check that code meets minimum coverage percentage `/path/to/project$ gocheckcov
check`

## Usage

```
$ gocheckcov
analyzes coverage

Usage:
  gocheckcov [command]

Available Commands:
  check       Check whether pkg coverage meets specified minimum
  help        Help about any command
  version     print the gocheckcov version

Flags:
  -h, --help      help for gocheckcov
  -v, --verbose   turn on verbose logging

Use "gocheckcov [command] --help" for more information about a command.
```

### Enforce Minimum Coverage From A Configuration File

```
$ gocheckcov check
```

If all packages do not meet the specifed minimum coverage percentage gocheckcov
will return exit code 1.

gocheckcov will search for a configuration file `.gocheckcov-config.yml` at the
current working directory.

You can also specify the path to a configuration file using the `--config-file`
option.

#### Print out source and coverage for each function

gocheckcov can optionally print out the source and coverage for each function by
using the `--print-src` flag.

Statements covered are colored green, statements uncovered are colored red, and
source which is not applicable to coverage is color white. ``` $ gocheckcov
check --profile-file ${coverprofile_path} --print-src

func handleIfStmt coverage 77.77% statements 7/9

func (sc *StmtCollector) handleIfStmt(s *ast.IfStmt, fset *token.FileSet) error
{ if s.Init != nil { if err := sc.Collect(s.Init, fset); err != nil { return err
} }

          if err := sc.Collect(s.Body, fset); err != nil {
                  return err
          }

          if s.Else != nil {
                  if err := sc.handleIfStmtElse(s, fset); err != nil {
                          return err
                  }
          }

          return nil

} ```

### Initialize A New Configuration File Using Current Coverage Percentages

```
$ gocheckcov check init --profile-file ${coverprofile_path}
$ ls
.gocheckcov-config.yml
....
$ cat .gocheckcov-config.yml
packages:
- name: some/pkg/in/your/path
  min_coverage_percentage: 34.5 # whatever the current coverage is measured at
```

gocheckcov will write out a configuration file `.gocheckcov-config.yml` at the
current working directory using the current coverage measured for each package
in the specified path.

### Supported Golang Versions

*   1.11.x
*   1.12.x
*   1.13.x

### Install

`go get github.com/cvgw/gocheckcov`

### Configuration

Specify minimum coverage for each package via a configuration file ```

# .gocheckcov-config.yaml

min_coverage_percentage: 25 packages: - name: github.com/bar/foo/pkg/baz # this
overrides the global val of min_coverage_percentage for only this package
mininum_coverage_percentage: 66.6 ```

## Development

gocheckcov uses `dep` for dependency management and `golangci-lint` for linting.
See the [development guide](./DEVELOPMENT.md) for more info.

## Contributing

Contributors are welcome and appreciated. Please read the
[contributing guide](./CONTRIBUTING.md) for more info.

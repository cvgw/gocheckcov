# Development

This doc explains the development workflow so you can get started
[contributing](./CONTRIBUTING.md) to gocheckcov!

* [Supported Golang Versions](#supported-golang-versions)
* [Getting Started](#getting-started)
* [Checkout Your Fork](#checkout-your-fork)
* [Testing](#testing-gocheckcov)
* [Creating A PR](#creating-a-pr)

## Supported Golang versions
* 1.11.x
* 1.12.x
* 1.13.x

## Getting started

First you will need to setup your GitHub account and create a fork:

1. Create [a GitHub account](https://github.com/join)
2. Setup [GitHub access via
   SSH](https://help.github.com/articles/connecting-to-github-with-ssh/)
3. [Create and checkout a repo fork](#checkout-your-fork)

When you're ready, you can [create a PR](#creating-a-pr)!

## Checkout your fork

The Go tools require that you clone the repository to the `src/github.com/cvgw/gocheckcov` directory
in your [`GOPATH`](https://github.com/golang/go/wiki/SettingGOPATH).

To check out this repository:

1. Create your own [fork of this
  repo](https://help.github.com/articles/fork-a-repo/)
2. Clone it to your machine:

  ```shell
  mkdir -p ${GOPATH}/src/github.com/cvgw
  cd ${GOPATH}/src/github.com/cvgw
  git clone git@github.com:${YOUR_GITHUB_USERNAME}/gocheckcov.git
  cd gocheckcov
  git remote add upstream git@github.com:cvgw/gocheckcov.git
  git remote set-url --push upstream no_push
  ```

_Adding the `upstream` remote sets you up nicely for regularly [syncing your
fork](https://help.github.com/articles/syncing-a-fork/)._

## Testing gocheckcov

gocheckcov has [unit tests](#unit-tests)

### Unit Tests

The unit tests live with the code they test and can be run with:

```shell
make test
```

_These tests will not run correctly unless you have [checked out your fork into your `$GOPATH`](#checkout-your-fork)._

## Creating a PR

When you have changes you would like to propose to gocheckcov, you will need to:

1. Ensure the commit message(s) describe what issue you are fixing and how you are fixing it
   (include references to [issue numbers](https://help.github.com/articles/closing-issues-using-keywords/)
   if appropriate)
1. [Create a pull request](https://help.github.com/articles/creating-a-pull-request-from-a-fork/)

### Reviews

Each PR must be reviewed by a maintainer.

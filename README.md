# 🚅 Tokaido by Ironstar

[![CircleCI](https://circleci.com/gh/ironstar-io/tokaido.svg?style=shield)](https://circleci.com/gh/ironstar-io/tokaido)
[![GitHub stars](https://img.shields.io/github/stars/ironstar-io/tokaido.svg)](https://github.com/ironstar-io/tokaido/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/ironstar-io/tokaido.svg)](https://github.com/ironstar-io/tokaido/issues)
[![GitHub license](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/ironstar-io/tokaido)

A command line utility for quickly and painlessly spinning up your Drupal environment

## Why

## Dependencies

### Docker and Docker Compose

TODO: Docker desktop install instructions

## Installation

### The Easy Way

### Build From Source

TODO: Better instructions
Tokaido is build in Golang 1.10.2

- Install Go and dep with `brew install dep` (Latest go version is included)
- Clone this repository
- Install package dependencies with `dep ensure`
- From the root of the cloned repo run `make build`
- Your local executable is now avaialable with `./dist/tok [command]`

Getting this message or something similar?

```sh
import /Users/jimmycann/www/golang/pkg/darwin_amd64/bitbucket.org/ironstar/tokaido-cli/utils.a: object is [darwin amd64 go1.9.2 X:framepointer] expected [darwin amd64 go1.10.2 X:framepointer]
```

There are two versions of go installed on your system. Easiest is to remove the version in `/usr/local/go` with `sudo rm -rfv /usr/local/go`. Your brew installed version should now be working, check your version with `go version`

## Usage

## Roadmap

## About Tokaido

## About Ironstar

## Authors

The team at Ironstar

- Mike Richardson (@otakumike)
- Jimmy Cann (@yjimk)

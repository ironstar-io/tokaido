# 🚅 Tokaido by Ironstar

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
- Install package dependencies with `dep ensure`
- Clone this repository
- From the root of the cloned repo run `go build -o tok`
- Your executable is now avaialable with `./tok [command]`

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

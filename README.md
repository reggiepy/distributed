# golang Implemented Simple Distributed System

[comment]: <> "[![Build Status]&#40;https://travis-ci.org/reggiepy/distributed.svg?branch=master&#41;]&#40;https://travis-ci.org/reggiepy/distributed&#41;"

[comment]: <> "[![Coverage Status]&#40;https://coveralls.io/repos/github/reggiepy/distributed/badge.svg?branch=master&#41;]&#40;https://coveralls.io/github/reggiepy/distributed?branch=master&#41;"

[comment]: <> ([![go version]&#40;https://img.shields.io/badge/go-1.16-success.svg?style=flat&#41;]&#40;&#41;)

[![go version](https://img.shields.io/github/go-mod/go-version/reggiepy/distributed/master?color=success&filename=go.mod&style=flat)]()
[![release](https://img.shields.io/github/v/tag/reggiepy/distributed?color=success&label=release)]()
[![build status](https://img.shields.io/badge/build-pass-success.svg?style=flat)]()
[![Coverage Status](https://img.shields.io/badge/coverage-80%25-success.svg?style=flat)]()
[![License](https://img.shields.io/badge/license-GNU%203.0-success.svg?style=flat)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/reggiepy/distributed)](https://goreportcard.com/report/github.com/reggiepy/distributed)

## Installation

```bash
git clone https://github.com/reggiepy/distributed.git
go mod tidy
cd distributed
```

## Usage

```bash
go run cmd/registryservice/main.go
go run cmd/logservice/main.go
go run cmd/gradingservice/main.go
go run cmd/portal/main.go
```

## Architecture

![image-20220328155409148](README.assets/image-20220328155409148.png)
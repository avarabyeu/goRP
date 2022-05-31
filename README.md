![GitHub Workflow Status](https://img.shields.io/github/workflow/status/reportportal/goRP/Build)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/eBay/fabio/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/reportportal/goRP)](https://goreportcard.com/report/github.com/reportportal/goRP)

# goRP

Golang Client and CLI Utility for [ReportPortal](https://reportportal.io)

## Installation

## Usage

```
gorp [global options] command [command options] [arguments...]   

COMMANDS:
   launch   Operations over launches
   report   Reports input to report portal
   init     Initializes configuration cache
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --uuid value, -u value     Access Token [$GORP_UUID]
   --project value, -p value  ReportPortal Project Name [$GORP_PROJECT]
   --host value               ReportPortal Server Name
   --help, -h                 show help (default: false)
   --version, -v              print the version (default: false)
```

### Init command

    NAME:
        gorp init - Initializes configuration cache
    USAGE:
        gorp init [command options] [arguments...]
    OPTIONS:
        --help, -h  show help (default: false)

### Launch command

```
USAGE:
   goRP launch command [command options] [arguments...]

COMMANDS:
   list     List launches
   merge    Merge Launches
   help, h  Shows a list of commands or help for one command
```

#### List Launches

```
USAGE:
   goRP launch list [command options] [arguments...]

OPTIONS:
   --filter-name value, --fn value  Filter Name [$FILTER_NAME]
   --filter value, -f value         Filter [$Filter]
   --help, -h                       show help (default: false)
```

### Report command

    NAME:
        goRP report - Reports input to report portal
    USAGE:
        goRP report command [command options] [arguments...]
    COMMANDS:
        test2json  Input format: test2json
        help, h    Shows a list of commands or help for one command
    OPTIONS:
        --help, -h  show help (default: false)
   

## Using as Golang Test Results Agent
Run tests with JSON output
```
go test -json ./... > results.txt
```
Report The results
```
gorp report test2json -f results.txt
```
Report directly from go test output
```
go test -json ./... | bin/gorp report test2json
```
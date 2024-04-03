![Build Status](https://github.com/reportportal/goRP/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/reportportal/goRP)](https://goreportcard.com/report/github.com/reportportal/goRP)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/reportportal/goRP/master/LICENSE)
[![Release](https://img.shields.io/github/release/reportportal/goRP.svg)](https://github.com/reportportal/goRP/releases/latest)
[![GitHub Releases Stats of goRP](https://img.shields.io/github/downloads/reportportal/goRP/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=reportportal&repository=gorP)

# goRP

Golang Client and CLI Utility for [ReportPortal](https://reportportal.io)

## Installation

- Via Go Install
```sh
go install github.com/reportportal/goRP@latest
```
- Via cURL (passing version and arch)
```sh
curl -sL https://github.com/avarabyeu/goRP/releases/download/v5.0.2/goRP_5.0.2_darwin_amd64.tar.gz | tar zx -C .
```
- Via cURL (latest one)
```sh
curl -s https://api.github.com/repos/reportportal/goRP/releases/latest | \
  jq -r '.assets[] | select(.name | contains ("tar.gz")) | .browser_download_url' | \
  grep "$(uname)_$(arch)" | \
  xargs curl -sL |  tar zx -C .
```
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
go test -json ./... | gorp report test2json
```
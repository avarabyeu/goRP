![GitHub Workflow Status](https://img.shields.io/github/workflow/status/avarabyeu/goRP/Build)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/eBay/fabio/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/avarabyeu/goRP)](https://goreportcard.com/report/github.com/avarabyeu/goRP)

# goRP
Golang Client and CLI Utility for [ReportPortal](https://reportportal.io)

## Installation

## Usage
```
gorp [global options] command [command options] [arguments...]   

COMMANDS:
     launch   Operations over launches
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -u value, --uuid value     Access Token [$GORP_UUID]
   -p value, --project value  ReportPortal Project Name [$GORP_PROJECT]
   --host value               ReportPortal Server Name
   --help, -h                 show help
   --version, -v              print the version
```

### Launch command
```
USAGE:
   goRP launch command [command options] [arguments...]

COMMANDS:
     list  List launches
```

#### List Launches
```
USAGE:
   goRP launch list [command options] [arguments...]

OPTIONS:
   --fn value, --filter-name value  Filter Name [$FILTER_NAME]
   -f value, --filter value         Filter [$Filter]
```

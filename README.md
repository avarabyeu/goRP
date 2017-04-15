[![Build Status](https://travis-ci.org/avarabyeu/goRP.svg?branch=master)](https://travis-ci.org/avarabyeu/goRP)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/eBay/fabio/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/avarabyeu/goRP)](https://goreportcard.com/report/github.com/avarabyeu/goRP)

# goRP
Lightweight repack of ReportPortal

### Notable changes

##### Lightweight. Three java-based components were replaced with golang alternatives
1. Service Registry: [Eureka](https://github.com/Netflix/eureka) replaced with [Consul](https://www.consul.io/)
2. Gateway: [Zuul](https://github.com/Netflix/zuul) replaced with [fabio](https://github.com/eBay/fabio)   
3. Added golang-based root service that replaces logic of ReportPortal's service-gateway component
3. UI Service: Java-based replaced with golang-based alternative

##### Common modules:
1. commons. General-purpose utilities
2. conf. Utilities for service configuration
3. registry. Consul and Eureka clients to register/deregister services
4. server. Utilities for configuring and bootstrapping REST services

### Installation
##### All-in-one pack:
All ReportPortal/goRP components without mongodb.

##### All-in-one pack:
Consists of all components and mongodb inside one container
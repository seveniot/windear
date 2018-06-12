# Windear

[![Build Status](https://travis-ci.org/SevenIOT/windear.svg?branch=master)](https://travis-ci.org/SevenIOT/windear)
[![Code Coverage](https://codecov.io/gh/SevenIOT/windear/branch/master/graph/badge.svg)](https://codecov.io/gh/SevenIOT/windear)
[![GoDoc](https://godoc.org/github.com/SevenIOT/windear?status.svg)](https://godoc.org/github.com/SevenIOT/windear)

## Introduction

A simple MQTT broker write in Golang, which support cluster.

## Features

* Support MQTT v3.1.1
* Support Cluster

## Getting Started

#### dependent libraries
```
$ dep ensure
```

### build
```
$ go build
```

#### run
```
$ ./windear -c conf/config.yaml
```

## License

Released under the [MIT License](https://github.com/SevenIOT/windear/blob/master/License)
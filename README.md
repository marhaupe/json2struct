# json2struct


[![Build Status](https://travis-ci.com/marhaupe/json2struct.svg?branch=master)](https://travis-ci.com/marhaupe/json2struct)
[![codecov](https://codecov.io/gh/marhaupe/json2struct/branch/master/graph/badge.svg)](https://codecov.io/gh/marhaupe/json2struct)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com) 
<!-- [![GoDoc](https://godoc.org/github.com/marhaupe/json2struct?status.svg)](https://godoc.org/github.com/marhaupe/json2struct) -->

> CLI tool to convert JSON to Go type definitions

At some point when dealing with JSONs in Go, you will have to write types to `json.Unmarshal` your JSONs into. Doing this by hand is not only repetitive and time consuming, but also error prone. `json2struct` saves you this work by automatically parsing the JSON and generating you the matching type definitions ready to be used.

Different to other tools, `json2struct` tries to avoid generating `interface{}` and `map[string]interface{}` as much as possible. Nonetheless it's very fast 🚀.

# Installation

Simply grab the latest release [binaries](https://github.com/marhaupe/json2struct/releases). 

# Usage

> json2struct [flags]

Calling `json2struct` without flags opens a text editor. Simply input your JSON and save and exit. 

![Example](.github/demo.gif)

## Flags:

### File
>  -f, --file string:     path to JSON file 

```bash
json2struct -f apiResult.json
```

### String
>  -s, --string string:   JSON string


```bash
 json2struct -s "$(curl "https://reqres.in/api/users?page=2")"
```


### Benchmark
>  -b, --benchmark:       measure execution time
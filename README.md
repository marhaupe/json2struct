# json2struct

[![Build Status](https://travis-ci.com/marhaupe/json2struct.svg?branch=master)](https://travis-ci.com/marhaupe/json2struct)
[![codecov](https://codecov.io/gh/marhaupe/json2struct/branch/master/graph/badge.svg)](https://codecov.io/gh/marhaupe/json2struct)
[![GoDoc](https://godoc.org/github.com/marhaupe/json2struct?status.svg)](https://godoc.org/github.com/marhaupe/json2struct)

This project aims to make your life a lot easier by automatically generating structs for a given JSON. 

## Installation

```bash
go get github.com/marhaupe/json2struct
```

## Usage

### Base command:
Calling `json2struct` without arguments opens the superior text editor for unix systems, vim 🤖. Insert the JSON data you want to parse and save and exit - in case you always forget how to do that: `:wq!`.

![Showcase](.github/with_editor.gif)

### Flags:
Call `json2struct -s` or `json2struct --string` with the JSON data as argument. 

```bash
 json2struct -s "$(curl "https://reqres.in/api/users?page=2")"
```

![Showcase](.github/direct_input.gif)


The `string` option lets you pipe JSON data as input. The current implementation lacks some features, e.g. you need to escape quotes manually. PR's are more than welcome.


## Lastly

Please feel free to open a pull request for missing features or bugs.

Credits to Matt Holt (https://github.com/mholt/json-to-go), from whom I got the idea.

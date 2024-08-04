# json2struct

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)
[![codecov](https://codecov.io/gh/marhaupe/json2struct/branch/master/graph/badge.svg)](https://codecov.io/gh/marhaupe/json2struct)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

<!-- [![GoDoc](https://godoc.org/github.com/marhaupe/json2struct?status.svg)](https://godoc.org/github.com/marhaupe/json2struct) -->

> CLI tool to convert JSON to struct type definitions

At some point when dealing with JSONs in Go, you will have to write struct types to `json.Unmarshal` your JSONs into. Doing this by hand is not only repetitive and time consuming, but also error prone. `json2struct` saves you this work by automatically parsing the JSON and generating you the matching struct type definitions ready to be used.

Unlike other tools, `json2struct` tries to avoid generating `interface{}` and `map[string]interface{}` as much as possible. Nonetheless it's very fast 🚀.

## Installation

### Homebrew

```bash
brew tap marhaupe/json2struct https://github.com/marhaupe/json2struct

brew install marhaupe/json2struct/json2struct
```

### Manually

Grab the latest release [binaries](https://github.com/marhaupe/json2struct/releases).

## Usage

> json2struct [options]

Calling `json2struct` without flags opens a text editor. Simply input your JSON and save and exit.

![Example](.github/demo.gif)

### Options

You probably don't want to manually write that 1MB JSON you have to generate a struct for by hand. I mean, if you really want to, I'm not here to judge, but that's not the point. These options will make your life easier. If you miss some, feel free to open an issue.

#### Generating a struct from a string

> -s, --string string: JSON string

This is basically your bread and butter thanks to pipes. Usage:

```bash
 json2struct -s "$(curl "https://reqres.in/api/users?page=2")"
```

#### Generating a struct from an existing file

> -f, --file string: path to JSON file

This is useful if you have a JSON file stored in your filesystem and are too lazy to use pipes. Usage:

```bash
json2struct -f input.json
```

#### Generating a struct from the clipboard to the clipboard

> -c, --clipboard: read from and write types to clipboard

Reads JSON from clipboard, generates types and writes those types to the clipboard.

```bash
json2struct -c
```

#### Other options

> -b, --benchmark: measure execution time

> -h, --help: help for json2struct

> --version: version for json2struct

## Specifics

A lot of the lexing/parsing in this project is inspired by the [`text/template`](https://go.dev/src/text/template/) package. Rob Pike also gave a talk about how and why they wrote the package. The link can be found [here](https://www.youtube.com/watch?v=HxaD_trXwRE), and the slides to it [here](https://go.dev/talks/2011/lex.slide#1)

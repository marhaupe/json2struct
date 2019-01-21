# json2struct

[![Build Status](https://travis-ci.com/marhaupe/json2struct.svg?branch=master)](https://travis-ci.com/marhaupe/json2struct)
[![codecov](https://codecov.io/gh/marhaupe/json2struct/branch/master/graph/badge.svg)](https://codecov.io/gh/marhaupe/json2struct)

This projects aims to make your life a lot easier by automatically generating structs for a given JSON. 

## Installing

```bash
go get github.com/marhaupe/json2struct
```

## Usage

### Base command:
Calling `json2struct` without arguments opens the superior text editor for unix systems, vim 🤖. Write the JSON you want to parse in there and save and exit - in case you always forget how to do that: `:wq!`.

![Showcase](.github/with_editor.gif)

### Flags:
Call `json2struct -s` or `json2struct --string` with the JSON string as an argument. 

```bash
 json2struct -s "$(curl "https://reqres.in/api/users?page=2")"
```

![Showcase](.github/direct_input.gif)


This option lets you pipe JSONs as input, but does not handle inputting JSONs directly very well since you need to escape quotes and stuff like that.

## Lastly

Please feel free to open pull requests for features you miss, stuff that doesn't work, the usual.

Credits to Matt Holt (https://github.com/mholt/json-to-go), from whom I got the idea from.
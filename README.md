# json2struct

[![Build Status](https://travis-ci.com/marhaupe/json2struct.svg?branch=master)](https://travis-ci.com/marhaupe/json2struct)
[![Coverage Status](https://coveralls.io/repos/github/marhaupe/json2struct/badge.svg?branch=master)](https://coveralls.io/github/marhaupe/json2struct?branch=master)

This projects aims to make your life a lot easier by automatically generating structs for a given JSON. 

## Installing

```bash
go get github.com/marhaupe/json2struct
```

## Usage

### No additional commands:
Call `json2struct` with the JSON string as an argument. 

```bash
 json2struct "$(curl "https://reqres.in/api/users?page=2")"
```
![Showcase](.github/showcase.gif)


This option lets you pipe JSONs as input, but does not handle inputting JSONs directly very well since you need to escape quotes and stuff like that.

### Command `create`:
Calling `json2struct create` opens the superior text editor for unix systems, vim 🤖. Write the JSON you want to parse in there and save and exit - just in case you get stuck: `:wq!`. 

## Lastly

Please feel free to open pull requests for features you miss, stuff that doesn't work, the usual.

Credits to Matt Holt (https://github.com/mholt/json-to-go), from whom I got the idea from.
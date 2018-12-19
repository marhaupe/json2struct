# json2struct

[![Build Status](https://travis-ci.com/marhaupe/json2struct.svg?branch=master)](https://travis-ci.com/marhaupe/json2struct)
[![Coverage Status](https://coveralls.io/repos/github/marhaupe/json2struct/badge.svg?branch=master)](https://coveralls.io/github/marhaupe/json2struct?branch=master)

This projects aims to make your life a lot easier by automatically generating structs for a given json. 

### Usage

At the moment, there are two possible options to call `json2struct`

#### Option 1

Just call `json2struct`. This opens the superior text editor for unix systems, vim 🤖. Write the json you want to parse in there and save and exit - just in case you get stuck: `:wq!`. 

#### Option 2

Call `json2struct` with the json string as an argument. 

```bash
 json2struct "$(curl "https://reqres.in/api/users?page=2")" >> generated.go
```
![Showcase](.github/showcase.gif)


This option lets you pipe jsons as input, but does not handle inputting jsons directly very vell since you need to escape quotes and stuff like that.

### Lastly

Please feel free to open pull requests for features you miss, stuff that doesn't work, the usual.

Credits to Matt Holt (https://github.com/mholt/json-to-go), from whom I got the idea from.
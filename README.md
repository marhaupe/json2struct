# json2struct

[![Build Status](https://travis-ci.com/marhaupe/json2struct.svg?branch=master)](https://travis-ci.com/marhaupe/json2struct)


This projects aims to make your life a lot easier by automatically generating structs for a given json. Just call `json2struct` with your json as first argument: 

![Showcase](.github/showcase.gif)

The program will output a valid go struct to stdout. Please check whether it has been correctly generated, and if not, I'd greatly appreciate you opening an issue!

Generating a new file `generated.go` based on a received json: 
```bash
 json2struct "$(curl "https://reqres.in/api/users?page=2")" >> generated.go
```

Credits to Matt Holt (https://github.com/mholt/json-to-go), from whom I got the idea from.
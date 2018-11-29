# JSON-to-struct

[![Build Status](https://travis-ci.com/marhaupe/json-to-struct.svg?branch=master)](https://travis-ci.com/marhaupe/json-to-struct)

-- WORK IN PROGRESS -

This projects aims to make your life a lot easier by automatically generating structs for a given json. Just call `json-to-struct` with your json as first argument. 

The program will output a valid go struct to stdout. Please check whether it has been correctly generated, and if not, I'd greatly appreciate you opening an issue!

Generating a new file `generated.go` based on a received json: 
```bash
curl https://pokeapi.co/api/v2/pokemon/1/ |json-to-struct >> generated.go
```

Credits to Matt Holt (https://github.com/mholt/json-to-go), who inspired me to write a solution without relying on a web-application
# JSON-to-struct

-- WORK IN PROGRESS -

This projects aims to make your life a lot easier by automatically generating structs for a given json. Just call `json-to-struct` with your json as first argument. 

Example: 
```bash
curl https://pokeapi.co/api/v2/pokemon/1/ |json-to-struct
```

The program will generate a new .go file holding the generated struct. Please check whether the file has been correctly generated, and if not, I'd greatly appreciate you opening an issue!

Credits to Matt Holt (https://github.com/mholt/json-to-go), who inspired me to write a solution without relying on a web-application
## Different valid JSONs:

### Different objects:
```json
[
  {
    "thissucks": true
  },
  {
    "thisdoesntsuck": {
          "value": false
    }
  }
]
```

is going to be this:
```go
type JsonToStruct []struct{
  Thissucks bool `json:"thissucks,omitempty"`
  Thisdoesntsuck struct{
    Value bool `json:"value"`
  } `json:"thisdoesntsuck,omitempty"`
}
```

### Different datatypes:
```json
[
  "test",
  "auch ein test",
  true,
  5000
]
```

is going to be this:
```go
type JsonToStruct []interface{}
```
because there's no way to store different datatypes in a go slice


### Same, primitive datatypes:
```json
[
  5000,
  8000
]
```

is going to be this:
```go
type JsonToStruct []int
```
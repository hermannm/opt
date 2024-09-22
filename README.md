# opt

A simple Go package that provides a generic `Option` type: a container that either has a value or is
empty.

Run `go get hermannm.dev/opt` to add it to your project!

## Motivation

One of the Go proverbs is to
["make the zero value useful"](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=385s). But this is not
always practical - for example, one may want to distinguish between a blank string and an optional
string that has not been set. In these cases, one may choose to use a pointer instead, with `nil` as
the zero value. But working with pointers can be more cumbersome, risk panics from nil
dereferencing, and have negative performance implications if values are moved to the heap.

This package provides an alternative, through the generic `Option[T]` wrapper type. It aims to
provide a zero value for types (the empty option), and clearly signal to readers of the code that
the value is optional.

## Usage

`Option` has three constructors:

- `Value(T)` creates an `Option[T]` containing the given value:

  ```go
  name := opt.Value("hermannm")
  if name.HasValue() {
      fmt.Printf("Hello, %s\n", name.Value)
  }
  ```

- `Empty()` creates an empty option:

  ```go
  name := opt.Empty[string]()
  if name.IsEmpty() {
      fmt.Println("No name given")
  }
  ```

- `FromPointer(*T)` creates an `Option[T]` with the pointer's value, or an empty option if the
  pointer is `nil`:

  ```go
  var pointer *string
  option := opt.FromPointer(pointer)
  if value, ok := option.Get(); ok {
      fmt.Printf("Got value %s\n", value)
  } else {
      fmt.Println("No value received")
  }
  ```

Once an option is created, you can replace/remove its value with `Put`/`Clear`:

```go
option := opt.Empty[string]()

option.Put("value")
fmt.Println(option) // Prints "value"

option.Clear()
fmt.Println(option) // Prints "<empty>"
```

Finally, `Option` implements `json.Marshaler` and `json.Unmarshaler`. An empty option marshals to
`null`, and a `null` JSON value unmarshals to an empty option:

```go
type Person struct {
    Name opt.Option[string] `json:"name"`
    Age  opt.Option[int]    `json:"age"`
}

person1 := Person{
    Name: opt.Value("hermannm"),
    Age:  opt.Empty[int](),
}
jsonOutput, err := json.Marshal(object1)
// Output: {"name":"hermannm","age":null}

jsonInput := []byte(`{"name":null,"age":25}`)
var person2 Object
err := json.Unmarshal(jsonInput, &person2)
// Name is now empty, while Age has value 25
```

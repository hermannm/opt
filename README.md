# opt

A simple Go package that provides a generic `Option` type: a container that either has a value or is
empty.

Run `go get hermannm.dev/opt` to add it to your project!

**Docs:** [pkg.go.dev/hermannm.dev/opt](https://pkg.go.dev/hermannm.dev/opt)

**Contents:**

- [Motivation](#motivation)
- [Usage](#usage)
- [Maintainer's guide](#maintainers-guide)

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

  <!-- @formatter:off -->
  ```go
  name := opt.Value("hermannm")
  if name.HasValue() {
  	fmt.Printf("Hello, %s\n", name.Value)
  }
  ```
  <!-- @formatter:on -->

- `Empty()` creates an empty option:

  <!-- @formatter:off -->
  ```go
  name := opt.Empty[string]()
  if name.IsEmpty() {
  	fmt.Println("No name given")
  }
  ```
  <!-- @formatter:on -->

- `FromPointer(*T)` creates an `Option[T]` with the pointer's value, or an empty option if the
  pointer is `nil`:

  <!-- @formatter:off -->
  ```go
  var pointer *string
  option := opt.FromPointer(pointer)
  if value, ok := option.Get(); ok {
  	fmt.Printf("Got value %s\n", value)
  } else {
  	fmt.Println("No value received")
  }
  ```
  <!-- @formatter:on -->

Once an option is created, you can replace/remove its value with `Put`/`Clear`:

<!-- @formatter:off -->
```go
option := opt.Empty[string]()

option.Put("value")
fmt.Println(option) // Prints "value"

option.Clear()
fmt.Println(option) // Prints "<empty>"
```
<!-- @formatter:on -->

Finally, `Option` implements `json.Marshaler` and `json.Unmarshaler`. An empty option marshals to
`null`, and a `null` JSON value unmarshals to an empty option:

<!-- @formatter:off -->
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
<!-- @formatter:on -->

## Maintainer's guide

### Publishing a new release

- Run tests and linter ([`golangci-lint`](https://golangci-lint.run/)):
  ```
  go test ./... && golangci-lint run
  ```
- Add an entry to `CHANGELOG.md` (with the current date)
    - Remember to update the link section, and bump the version for the `[Unreleased]` link
- Create commit and tag for the release (update `TAG` variable in below command):
  ```
  TAG=vX.Y.Z && git commit -m "Release ${TAG}" && git tag -a "${TAG}" -m "Release ${TAG}" && git log --oneline -2
  ```
- Push the commit and tag:
  ```
  git push && git push --tags
  ```
    - Our release workflow will then create a GitHub release with the pushed tag's changelog entry

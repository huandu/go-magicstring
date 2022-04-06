# Magic String

[![Go](https://github.com/huandu/go-magicstring/workflows/Go/badge.svg)](https://github.com/huandu/go-magicstring/actions)
[![Go Doc](https://godoc.org/github.com/huandu/go-magicstring?status.svg)](https://pkg.go.dev/github.com/huandu/go-magicstring)
[![Go Report](https://goreportcard.com/badge/github.com/huandu/go-magicstring)](https://goreportcard.com/report/github.com/huandu/go-magicstring)
[![Coverage Status](https://coveralls.io/repos/github/huandu/go-magicstring/badge.svg?branch=main)](https://coveralls.io/github/huandu/go-magicstring?branch=main)

This `magicstring` package is designed to attach arbitrary data to a Go built-in `string` type. The string with arbitrary data is called "magic string" here. Such string can be used as an ordinary string. We can read the attached data from a magic string freely.

## Usage

### Attach data and then read it

Call `Attach` to attach data into a string and `Read` to read the attached data in the magic string.

```go
type T struct {
    Name string
}

s1 := "Hello, world!"
data := &T{Name: "Kanon"}
s2 := Attach(s1, data)

attached := Read(s2).(*T)
fmt.Println(s1 == s2)         // true
fmt.Println(attached == data) // true
```

### Check whether a string is a magic string

Call `Is` if we want to know whether a string is a magic string, .

```go
s1 := "ordinary string"
s2 := Attach("magic string", 123)
s3 := s2
s4 := fmt.Sprint(s2)

fmt.Println(Is(s1)) // false
fmt.Println(Is(s2)) // true
fmt.Println(Is(s3)) // true
fmt.Println(Is(s4)) // false
```

### Copy a magic string

In general, we can use a magic string like an ordinary string. The attached data will be kept during all kinds of assignments. However, if we copy the content of a string to a buffer and create a new string from the buffer, we will lose the attached data.

The simplest way to create an ordinary string from a magic string is to call `Detach`. This function is optimized for ordinary strings. If a string is not a magic string, the `Detach` simply returns the string to avoid an unnecessary memory allocation and memory copy.

## Performance

Memory allocation is highly optimized for small strings. The maximum size of a small string is 18,408 bytes right now. It's the maximum memory span size, which is 18,432 bytes provided by `runtime.ReadMemStats()`, minus the size of magic string payload struct, which is 24 bytes right now.

Here is the performance data running on my MacBook.

```text
goos: darwin
goarch: amd64
pkg: github.com/huandu/go-magicstring
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAttachSmallString-12        13208282         82.08 ns/op       32 B/op        1 allocs/op
BenchmarkAttachLarge1MBString-12         7812        149331 ns/op  1057068 B/op        3 allocs/op
```

## License

This package is licensed under MIT license. See LICENSE for details.

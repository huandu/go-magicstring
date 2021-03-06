# Magic String

[![Go](https://github.com/huandu/go-magicstring/workflows/Go/badge.svg)](https://github.com/huandu/go-magicstring/actions)
[![Go Doc](https://godoc.org/github.com/huandu/go-magicstring?status.svg)](https://pkg.go.dev/github.com/huandu/go-magicstring)
[![Go Report](https://goreportcard.com/badge/github.com/huandu/go-magicstring)](https://goreportcard.com/report/github.com/huandu/go-magicstring)
[![Coverage Status](https://coveralls.io/repos/github/huandu/go-magicstring/badge.svg?branch=main)](https://coveralls.io/github/huandu/go-magicstring?branch=main)

This `magicstring` package is designed to attach arbitrary data to a Go built-in `string` type and read the data later. The string with attached data is called "magic string" here.

## Usage

### Attach data and then read it

Call `Attach` to attach data to a string and `Read` to read the attached data in the magic string. `Read` is thread-safe and extremely fast.

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

### Replace the attached data

Call `Replace` if we want to replace the attached data to a new one in a magic string. As `Replace` modifies the payload in a magic string, this call will affect all copies of this magic string and is not thread-safe.

```go
s := Attach("magic string", 123)
fmt.Println(Read(s)) // 123

success := Replace(s, "replaced")
fmt.Println(success) // true
fmt.Println(Read(s)) // replaced
```

### Copy a magic string

In general, we can use a magic string like an ordinary string. The attached data will be kept during all kinds of assignments. However, if we copy the content of a string to a buffer and create a new string from the buffer, we will lose the attached data.

```go
s1 := Attach("magic string", 123)
buf := make([]byte, len(s1))
copy(buf, s1)
s2 := string(buf)

fmt.Println(Is(s1)) // true
fmt.Println(Is(s2)) // false
```

The simplest way to create an ordinary string from a magic string is to call `Detach`. This function is optimized for ordinary strings. If a string is not a magic string, the `Detach` simply returns the string to avoid an unnecessary memory allocation and memory copy.

### Slice a magic string

A magic string cannot be sliced by built-in slice expression. If we want to keep the attachment in a magic string, we must call `Slice` to slice it. If a string is not a magic string, `Slice` just works the same as slice expression.

## Performance

Memory allocation is highly optimized for small strings. The maximum size of a small string is 18,408 bytes right now. It's the maximum size of memory span classes, which is 18,432 bytes provided by `runtime.ReadMemStats()`, minus the size of magic string payload struct, which is 24 bytes right now.

Here is the performance data running on my MacBook.

```text
goos: darwin
goarch: amd64
pkg: github.com/huandu/go-magicstring
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAttachSmallString-12       12312567     83.78 ns/op       32 B/op   1 allocs/op
BenchmarkAttachLarge1MBString-12        8511    165978 ns/op  1057046 B/op   3 allocs/op
BenchmarkReplaceSmallString-12      225952905    5.329 ns/op        0 B/op   0 allocs/op
BenchmarkReplaceLarge1MBString-12   228580016    5.268 ns/op        0 B/op   0 allocs/op
BenchmarkReadSmallString-12         312613485    3.915 ns/op        0 B/op   0 allocs/op
BenchmarkReadLarge1MBString-12      301951902    3.964 ns/op        0 B/op   0 allocs/op
```

## License

This package is licensed under MIT license. See LICENSE for details.

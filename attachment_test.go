// Copyright 2022 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package magicstring

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/huandu/go-assert"
)

func TestAttachData(t *testing.T) {
	a := assert.New(t)
	cases := []interface{}{
		1, "abcd", true, 1.2, complex(1, 2),
		struct{}{}, []byte("123"), &testStruct{},
		map[string][]int{"foo": {1, 2, 3}},
	}

	for _, c := range cases {
		iterateStrings(func(s string) {
			copied := Attach(s, c)
			a.Equal(s, copied)
			a.Assert(Is(copied))

			data := Read(copied)
			a.Equal(data, c)
		})
	}
}

func TestAttachNilData(t *testing.T) {
	a := assert.New(t)
	s1 := "nil data"
	s2 := Attach(s1, nil)
	s3 := Attach(s1, 123)
	s4 := Attach(s3, nil)
	a.Equal(s1, s2)
	a.Equal(s1, s3)
	a.Equal(s1, s4)

	// There is no magic attached to s2.
	a.Assert(!Is(s2))

	// The buffer in s1 and s2 must be identical.
	data1 := (*reflect.StringHeader)(unsafe.Pointer(&s1)).Data
	data2 := (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data
	a.Equal(data1, data2)

	// The buffer in s3 and s4 must be different.
	data3 := (*reflect.StringHeader)(unsafe.Pointer(&s3)).Data
	data4 := (*reflect.StringHeader)(unsafe.Pointer(&s4)).Data
	a.NotEqual(data3, data4)
}

func TestReadInvalidString(t *testing.T) {
	a := assert.New(t)
	s := "dummy"
	a.Assert(!Is(s))
	a.Assert(Read(s) == nil)
}

func TestAttachReplaceData(t *testing.T) {
	a := assert.New(t)
	s := "sample string"
	data1 := &testStruct{
		Foo: 123,
	}
	data2 := &testStruct{
		Foo: 567,
	}

	s1 := Attach(s, data1)
	a.Equal(s, s1)
	payload := Read(s1)
	a.Equal(data1, payload)

	s2 := Attach(s1, data2)
	a.Equal(s, s2)
	a.Equal(s1, s2)
	payload = Read(s2)
	a.NotEqual(data1, payload)
	a.Equal(data2, payload)
}

func TestAttachMapKey(t *testing.T) {
	a := assert.New(t)
	key := "foo"
	m := map[string]int{
		key: 123,
	}

	// If a key exists in a map, map will replace old key with the magic string.
	// WARNING: It's not guaranteed by Go languange spec. Don't rely on this behavior.
	data := 567
	foo := Attach(key, data)
	m[foo] = data

	for k, v := range m {
		a.Assert(Is(k))
		a.Equal(Read(k), v)
	}

	// If a key absents in a map, map will use the magic string as key.
	data = 999
	bar := Attach("bar", data)
	m[bar] = data
	delete(m, key)

	for k, v := range m {
		a.Assert(Is(k))
		a.Equal(Read(k), v)
	}
}

func TestDetach(t *testing.T) {
	a := assert.New(t)

	iterateStrings(func(s string) {
		data := &testStruct{
			Foo: 398,
		}
		attached := Attach(s, data)
		a.Equal(s, attached)

		// Call detach will not affect attached string.
		detached := Detach(attached)
		a.Equal(detached, s)
		payload := Read(attached)
		a.Equal(data, payload)

		// Attached data is nil in detached string.
		payload = Read(detached)
		a.Assert(payload == nil)

		// It's OK to detach twice.
		detached = Detach(detached)
		a.Equal(detached, s)
		payload = Read(detached)
		a.Assert(payload == nil)
	})
}

func ExampleAttach() {
	type T struct {
		Name string
	}
	s1 := "Hello, world!"
	data := &T{Name: "Kanon"}
	s2 := Attach(s1, data)

	attached := Read(s2).(*T)
	fmt.Println(s1 == s2)
	fmt.Println(attached == data)

	// Output:
	// true
	// true
}

func ExampleIs() {
	s1 := "ordinary string"
	s2 := Attach("magic string", 123)
	s3 := s2
	s4 := fmt.Sprint(s2)

	fmt.Println(Is(s1))
	fmt.Println(Is(s2))
	fmt.Println(Is(s3))
	fmt.Println(Is(s4))

	// Output:
	// false
	// true
	// true
	// false
}

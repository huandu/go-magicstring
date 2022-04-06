// Copyright 2022 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package magicstring

import (
	"runtime"
	"testing"
	"time"

	"github.com/huandu/go-assert"
)

func TestGCSafety(t *testing.T) {
	a := assert.New(t)
	value := 123

	iterateStrings(func(s string) {
		data := &testStruct{
			Foo: value,
		}
		copiedData := *data
		attached := Attach(s, data)

		dataChan := make(chan int, 1)
		gcDone := make(chan bool, 1)
		runtime.SetFinalizer(data, func(data *testStruct) {
			dataChan <- data.Foo
			gcDone <- true
		})
		data = nil // Don't keep reference to the data anymore.

		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond)
			runtime.GC()
		}

		// Decode must work as expected after a few GC calls.
		payload := Read(attached)
		a.Equal(payload, &copiedData)
		payload = nil

		// Clear attached string to release internal data.
		attached = ""

		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond)
			runtime.GC()
		}

		// The original data must be finalized.
		select {
		case <-gcDone:
		default:
			a.Fatalf("The data must be finalized")
		}

		select {
		case v := <-dataChan:
			a.Equal(v, value)
		default:
			a.Fatalf("The value in the data must be set in finalizer")
		}
	})
}

func BenchmarkAttachSmallString(b *testing.B) {
	s := "small"
	data := map[string]int{
		"foo": 123,
		"bar": 345,
	}

	for i := 0; i < b.N; i++ {
		Attach(s, data)
	}
}

func BenchmarkAttachLarge1MBString(b *testing.B) {
	s := makeString(1 * 1024 * 1024)
	data := map[string]int{
		"foo": 123,
		"bar": 345,
	}

	for i := 0; i < b.N; i++ {
		Attach(s, data)
	}
}

// Copyright 2022 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package magicstring

import (
	"reflect"
	"runtime"
	"unsafe"
)

const pointerMask uintptr = unsafe.Sizeof(uintptr(0)) - 1

// Attach associates a newly allocated string with data.
// The value of the returned string is guaranteed to be identical to str.
func Attach(str string, data interface{}) string {
	sz := len(str)

	if sz == 0 {
		return attachEmptyString(data)
	}

	if sz <= smallStringMax {
		return attachSmallString(str, data)
	}

	return attachLargeString(str, data)
}

// attachEmptyString allocates a new holder for empty string.
func attachEmptyString(data interface{}) string {
	payload := &magicStringPayload{}
	payload.checksum = makeChecksum(uintptr(unsafe.Pointer(payload)))
	payload.data = data
	dst := (*[sizeofPayload]byte)(unsafe.Pointer(payload))[0:0:0]
	return *(*string)(unsafe.Pointer(&dst))
}

// attachSmallString allocates a new holder type to hold both string content and data.
func attachSmallString(str string, data interface{}) string {
	sz := len(str)
	payload, dst := newPayload(sz)
	copy(dst, str)
	payload.checksum = makeChecksum(uintptr(unsafe.Pointer(&dst[0])))
	payload.data = data
	return *(*string)(unsafe.Pointer(&dst))
}

// attachLargeString allocates new memory to hold both string content and
// magic string payload. The data is referenced by string's finalizer.
func attachLargeString(str string, data interface{}) string {
	holder := make([]byte, len(str)+int(sizeofPayload))
	payload := (*magicStringPayload)(unsafe.Pointer(&holder[0]))
	dst := holder[sizeofPayload:]

	payload.checksum = makeChecksum(uintptr(unsafe.Pointer(&dst[0])))
	payload.data = data
	copy(dst, str)

	if data != nil {
		runtime.SetFinalizer(payload, func(payload *magicStringPayload) {
			// Hold data in this finalizer and clear it after finalized.
			data = nil
		})
	}

	return *(*string)(unsafe.Pointer(&dst))
}

// Replace replaces the attached data in str if str is a magic string.
// If str is an ordinary string, Replace creates a magic string with data.
//
// Replace returns true if the data is replaced, false if not.
func Replace(str string, data interface{}) bool {
	payload := read(str)

	if payload == nil {
		return false
	}

	payload.data = data
	return true
}

// Read reads the attached data inside the str.
// If there is no such data, it returns nil.
func Read(str string) interface{} {
	payload := read(str)

	if payload == nil {
		return nil
	}

	return payload.data
}

// Is checks if there is any data attached to str.
func Is(str string) bool {
	payload := read(str)
	return payload != nil
}

func read(str string) (payload *magicStringPayload) {
	data := unsafe.Pointer(uintptr(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&str)).Data)) &^ pointerMask)

	if data == nil {
		return
	}

	checksum := makeChecksum(uintptr(data))

	if len(str) == 0 {
		payload = (*magicStringPayload)(data)
	} else {
		payload = (*magicStringPayload)(unsafe.Pointer(uintptr(data) - sizeofPayload))
	}

	if payload.checksum != checksum {
		payload = nil
		return
	}

	return
}

// Detach returns a new string without any attached data.
// If str is an ordinary string, Detach just simply returns str.
func Detach(str string) string {
	if len(str) == 0 {
		return ""
	}

	payload := read(str)

	if payload == nil {
		return str
	}

	dst := make([]byte, len(str))
	copy(dst, str)
	return *(*string)(unsafe.Pointer(&dst))
}

// Slice returns a result of str[begin:end].
// If str is a magic string with attachment, the returned string references the same attachment of str.
func Slice(str string, begin, end int) (sliced string) {
	if begin > end {
		panic("go-magicstring: begin must not be greater than end in Slice")
	}

	payload := read(str)
	sliced = str[begin:end]

	if payload == nil {
		return
	}

	sliced = Attach(sliced, payload.data)
	return
}

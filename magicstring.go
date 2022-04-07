// Copyright 2022 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

// This `magicstring` package is designed to attach arbitrary data to a Go built-in `string` type
// and read the data later.
// The string with attached data is called "magic string" here.
package magicstring

import (
	"math"
	"reflect"
	"runtime"
	"sort"
	"unsafe"
)

const (
	sizeofPayload     = unsafe.Sizeof(magicStringPayload{})
	sizeofSizeClasses = unsafe.Sizeof(runtime.MemStats{}.BySize) / unsafe.Sizeof(runtime.MemStats{}.BySize[0])
)

// magicStringPayload is the payload holding extra data.
type magicStringPayload struct {
	checksum uint64
	data     interface{}
}

var (
	typeofMagicStringPayload = reflect.TypeOf(magicStringPayload{})
	typeofByte               = reflect.TypeOf(byte(0))

	holderWithSizeClasses []int
	holderTypes           []reflect.Type
	smallStringMax        int
)

func init() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	holderWithSizeClasses = make([]int, 0, sizeofSizeClasses)
	holderTypes = make([]reflect.Type, 0, sizeofSizeClasses)

	// Create struct types which perfectly fit a size class.
	for i := 0; i < int(sizeofSizeClasses); i++ {
		size := int(stats.BySize[i].Size)

		if size <= int(sizeofPayload) {
			continue
		}

		diff := size - int(sizeofPayload)
		t := reflect.StructOf([]reflect.StructField{
			{
				Name: "Payload",
				Type: typeofMagicStringPayload,
			},
			{
				Name: "Data",
				Type: reflect.ArrayOf(diff, typeofByte),
			},
		})
		holderWithSizeClasses = append(holderWithSizeClasses, diff)
		holderTypes = append(holderTypes, t)
	}

	// As sizes in runtime.MemStats.BySize is sorted by Size,
	// max size of a small string is the last of the slice.
	smallStringMax = holderWithSizeClasses[len(holderWithSizeClasses)-1]
}

// newPayload allocates a payload struct with enough space.
func newPayload(size int) (payload *magicStringPayload, dst []byte) {
	idx := sort.SearchInts(holderWithSizeClasses, size)
	sz := holderWithSizeClasses[idx] + int(sizeofPayload)
	t := holderTypes[idx]
	v := reflect.New(t)
	data := (*[math.MaxInt32]byte)(unsafe.Pointer(v.Pointer()))[:sz:sz]
	payload = (*magicStringPayload)(unsafe.Pointer(&data[0]))
	dst = data[sizeofPayload : int(sizeofPayload)+size : int(sizeofPayload)+size]
	return
}

var (
	checksumMask uint64 = 0xcfb6b2da518bead
)

// makeChecksum calculates a checksum for ptr.
// This checksum will be checked when decoding.
func makeChecksum(ptr uintptr) uint64 {
	return uint64(ptr) ^ checksumMask
}

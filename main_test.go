// Copyright 2022 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package magicstring

var testStrings = []string{
	"", "1", "123", "1234567", "12345678",
	makeString(9),
	makeString(16),
	makeString(31),
	makeString(65),
	makeString(31293),
	makeString(32943),
	makeString(34543),
	makeString(64930),
	makeString(67938),
	makeString(129485),
	makeString(249485),
	makeString(1_000_000),
	makeString(4*1024*1024 + 1),
}

type testStruct struct {
	Foo int
}

func makeString(size int) string {
	data := make([]byte, size)
	return string(data)
}

func iterateStrings(fn func(s string)) {
	for _, s := range testStrings {
		fn(s)
	}
}

// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// +build gofuzz

package messages

// Fuzz is to fuzz-testing message parser with go-fuzz.
// DO NOT CALL THIS.
//
//  go-fuzz -bin <go-m3ua dir>/messages/fuzz/messages-fuzz.zip -workdir <go-m3ua dir>/messages/fuzz/
func Fuzz(data []byte) int {
	if _, err := Parse(data); err != nil {
		return 0
	}

	return 1
}

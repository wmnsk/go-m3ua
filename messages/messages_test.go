// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"
)

func TestParseMalformed(t *testing.T) {
	cases := []struct {
		data []byte
		err  error
	}{
		{[]byte{0x00}, ErrTooShortToParse},
		{[]byte{0x00, 0x00}, ErrTooShortToParse},
		{[]byte{0x00, 0x00, 0x00}, ErrTooShortToParse},
		{[]byte{0x00, 0x00, 0x00, 0x00}, ErrTooShortToParse},
	}

	for _, c := range cases {
		if _, err := Parse(c.data); err != c.err {
			t.Errorf("Parse/unexpected error: got: %v, want: %v", err, c.err)
		}
	}
}

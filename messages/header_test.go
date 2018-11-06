// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"
)

func TestHeader(t *testing.T) {
	cases := []testCase{
		{
			"has-all",
			NewHeader(
				1,  // Version
				16, // Class
				16, // Type
				[]byte{
					0xde, 0xad, 0xbe, 0xef,
				},
			),
			[]byte{
				// Header
				0x01, 0x00, 0x10, 0x10, 0x00, 0x00, 0x00, 0x0c,
				// dummy Payload
				0xde, 0xad, 0xbe, 0xef,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := DecodeHeader(b)
		if err != nil {
			return nil, err
		}
		return v, nil
	})
}

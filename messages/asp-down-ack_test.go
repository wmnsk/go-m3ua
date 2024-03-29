// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"

	"github.com/wmnsk/go-m3ua/messages/params"
)

func TestAspDownAck(t *testing.T) {
	cases := []testCase{
		{
			"has-info",
			NewAspDownAck(params.NewInfoString("deadbeef")),
			[]byte{
				// Header
				0x01, 0x00, 0x03, 0x05, 0x00, 0x00, 0x00, 0x14,
				// InfoString
				0x00, 0x04, 0x00, 0x0c, 0x64, 0x65, 0x61, 0x64,
				0x62, 0x65, 0x65, 0x66,
			},
		},
		{
			"has-none",
			NewAspDownAck(nil),
			[]byte{
				// Header
				0x01, 0x00, 0x03, 0x05, 0x00, 0x00, 0x00, 0x08,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := ParseAspDownAck(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}

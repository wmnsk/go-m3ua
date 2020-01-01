// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"

	"github.com/wmnsk/go-m3ua/messages/params"
)

func TestAspInactiveAck(t *testing.T) {
	cases := []testCase{
		{
			"has-all",
			NewAspInactiveAck(
				params.NewRoutingContext(1),
				params.NewInfoString("deadbeef"),
			),
			[]byte{
				// Header
				0x01, 0x00, 0x04, 0x04, 0x00, 0x00, 0x00, 0x24,
				// RoutingContext
				0x00, 0x06, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// InfoString
				0x00, 0x04, 0x00, 0x0c, 0x64, 0x65, 0x61, 0x64,
				0x62, 0x65, 0x65, 0x66,
			},
		},
		{
			"has-rc",
			NewAspInactiveAck(params.NewRoutingContext(1), nil),
			[]byte{
				// Header
				0x01, 0x00, 0x04, 0x04, 0x00, 0x00, 0x00, 0x18,
				// RoutingContext
				0x00, 0x06, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
			},
		},
		{
			"has-info",
			NewAspInactiveAck(nil, params.NewInfoString("deadbeef")),
			[]byte{
				// Header
				0x01, 0x00, 0x04, 0x04, 0x00, 0x00, 0x00, 0x1c,
				// InfoString
				0x00, 0x04, 0x00, 0x0c, 0x64, 0x65, 0x61, 0x64,
				0x62, 0x65, 0x65, 0x66,
			},
		},
		{
			"has-none",
			NewAspInactiveAck(nil, nil),
			[]byte{
				// Header
				0x01, 0x00, 0x04, 0x04, 0x00, 0x00, 0x00, 0x10,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := ParseAspInactiveAck(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}

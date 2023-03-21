// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"

	"github.com/wmnsk/go-m3ua/messages/params"
)

func TestHeartbeatAck(t *testing.T) {
	cases := []testCase{
		{
			"has-all",
			NewHeartbeatAck(
				params.NewHeartbeatData([]byte{0xde, 0xad, 0xbe, 0xef}),
			),
			[]byte{
				// Header
				0x01, 0x00, 0x03, 0x06, 0x00, 0x00, 0x00, 0x10,
				// HeartbeatData
				0x00, 0x09, 0x00, 0x08, 0xde, 0xad, 0xbe, 0xef,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := ParseHeartbeatAck(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}

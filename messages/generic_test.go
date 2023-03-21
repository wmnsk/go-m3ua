// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"

	"github.com/wmnsk/go-m3ua/messages/params"
)

func TestGeneric(t *testing.T) {
	cases := []testCase{
		{
			// Class: 127, Type: 127, 1 Network Appearance and 2 Routing Context
			"has-all",
			New(
				1,   // Version
				127, // Message Class (Reserved)
				127, // Message Type (Reserved)
				params.NewNetworkAppearance(1),
				params.NewRoutingContext(1, 255),
			),
			[]byte{
				// Header
				0x01, 0x00, 0x7f, 0x7f, 0x00, 0x00, 0x00, 0x1c,
				// NetworkAppearance
				0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// RoutingContext
				0x00, 0x06, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0xff,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := ParseGeneric(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}

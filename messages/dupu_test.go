// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"

	"github.com/wmnsk/go-m3ua/messages/params"
)

func TestDestinationUserPartUnavailable(t *testing.T) {
	cases := []testCase{
		{
			"has-all",
			NewDestinationUserPartUnavailable(
				params.NewNetworkAppearance(1),
				params.NewRoutingContext(2),
				params.NewAffectedPointCode(3, 4),
				params.NewUserCause(params.UserIdentityUnknown, params.SCCP),
				params.NewInfoString("deadbeef"),
			),
			[]byte{
				// Header
				0x01, 0x00, 0x02, 0x05, 0x00, 0x00, 0x00, 0x40,
				// NetworkAppearance
				0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// RoutingContext
				0x00, 0x06, 0x00, 0x08, 0x00, 0x00, 0x00, 0x02,
				// AffectedPointCode
				0x00, 0x12, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x04,
				// User/Cause
				0x02, 0x04, 0x00, 0x08, 0x00, 0x01, 0x00, 0x00,
				// InfoString
				0x00, 0x04, 0x00, 0x0c, 0x64, 0x65, 0x61, 0x64, 0x62, 0x65, 0x65, 0x66,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := ParseDestinationUserPartUnavailable(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}

// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"

	"github.com/wmnsk/go-m3ua/messages/params"
)

func TestError(t *testing.T) {
	cases := []testCase{
		{
			"has-all",
			NewError(
				params.NewErrorCode(params.ErrInvalidVersion),
				params.NewRoutingContext(1),
				params.NewNetworkAppearance(1),
				params.NewAffectedPointCode(0x11111111, 0x22222222),
				params.NewDiagnosticInformation([]byte{0xde, 0xad, 0xbe, 0xef}),
			),
			[]byte{
				// Header
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3c,
				// ErrorCode
				0x00, 0x0c, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// RoutingContext
				0x00, 0x06, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// NetworkAppearance
				0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// AffectedPointCode
				0x00, 0x12, 0x00, 0x0c, 0x11, 0x11, 0x11, 0x11,
				0x22, 0x22, 0x22, 0x22,
				// DiagnosticInformation (+padding)
				0x00, 0x07, 0x00, 0x08, 0xde, 0xad, 0xbe, 0xef,
			},
		},
		{
			"has-none",
			NewError(nil, nil, nil, nil, nil),
			[]byte{
				// Header
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := DecodeError(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}

// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"

	"github.com/wmnsk/go-m3ua/messages/params"
)

func TestData(t *testing.T) {
	cases := []testCase{
		{
			"has-all",
			NewData(
				params.NewNetworkAppearance(1),
				params.NewRoutingContext(1),
				params.NewProtocolData(
					1, // OriginatingPointCode
					2, // DestinationPointCode
					3, // ServiceIndicator
					1, // NetworkIndicator
					0, // MessagePriority
					1, // SignalingLinkSelection
					[]byte{ // Data
						0xde, 0xad, 0xbe, 0xef,
					},
				),
				nil,
			),
			[]byte{
				// Header
				0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x34,
				// NetworkAppearance
				0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// RoutingContext
				0x00, 0x06, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// ProtocolData
				// Param Header
				0x02, 0x10, 0x00, 0x14,
				// OPC, DPC
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02,
				// SI, NI, MP, SLS
				0x03, 0x01, 0x00, 0x01,
				// Data
				0xde, 0xad, 0xbe, 0xef,
			},
		},
		{
			"has-rc",
			NewData(
				params.NewNetworkAppearance(1),
				nil,
				params.NewProtocolData(
					1, // OriginatingPointCode
					2, // DestinationPointCode
					3, // ServiceIndicator
					1, // NetworkIndicator
					0, // MessagePriority
					1, // SignalingLinkSelection
					[]byte{ // Data
						0xde, 0xad, 0xbe, 0xef,
					},
				),
				nil,
			),
			[]byte{
				// Header
				0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x2c,
				// NetworkAppearance
				0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// ProtocolData
				// Param Header
				0x02, 0x10, 0x00, 0x14,
				// OPC, DPC
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02,
				// SI, NI, MP, SLS
				0x03, 0x01, 0x00, 0x01,
				// Data
				0xde, 0xad, 0xbe, 0xef,
			},
		},
		{
			"has-info",
			NewData(
				nil,
				params.NewRoutingContext(1, 255),
				params.NewProtocolData(
					1, // OriginatingPointCode
					2, // DestinationPointCode
					3, // ServiceIndicator
					1, // NetworkIndicator
					0, // MessagePriority
					1, // SignalingLinkSelection
					[]byte{ // Data
						0xde, 0xad, 0xbe, 0xef,
					},
				),
				nil,
			),
			[]byte{
				// Header
				0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x30,
				// RoutingContexts
				0x00, 0x06, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0xff,
				// ProtocolData
				// Param Header
				0x02, 0x10, 0x00, 0x14,
				// OPC, DPC
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02,
				// SI, NI, MP, SLS
				0x03, 0x01, 0x00, 0x01,
				// Data
				0xde, 0xad, 0xbe, 0xef,
			},
		},
		{
			"has-none",
			NewData(
				nil, nil,
				params.NewProtocolData(
					1, // OriginatingPointCode
					2, // DestinationPointCode
					3, // ServiceIndicator
					1, // NetworkIndicator
					0, // MessagePriority
					1, // SignalingLinkSelection
					[]byte{ // Data
						0xde, 0xad, 0xbe, 0xef,
					},
				),
				nil,
			),
			[]byte{
				// Header
				0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x24,
				// ProtocolData
				// Param Header
				0x02, 0x10, 0x00, 0x14,
				// OPC, DPC
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02,
				// SI, NI, MP, SLS
				0x03, 0x01, 0x00, 0x01,
				// Data
				0xde, 0xad, 0xbe, 0xef,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := ParseData(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}

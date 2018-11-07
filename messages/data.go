// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// Data is a PayloadData type of M3UA message.
type Data struct {
	*Header
	NetworkAppearance *params.Param
	RoutingContext    *params.Param
	ProtocolData      *params.Param
	CorrelationID     *params.Param
}

// NewData creates a new Data.
func NewData(nwApr, rtCtx, pd, corrID *params.Param) *Data {
	d := &Data{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassTransfer,
			Type:     MsgTypePayloadData,
		},
		NetworkAppearance: nwApr,
		RoutingContext:    rtCtx,
		ProtocolData:      pd,
		CorrelationID:     corrID,
	}
	d.SetLength()

	return d
}

// Serialize returns the byte sequence generated from a Data.
func (d *Data) Serialize() ([]byte, error) {
	b := make([]byte, d.Len())
	if err := d.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize Data")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (d *Data) SerializeTo(b []byte) error {
	if len(b) < d.Len() {
		return ErrTooShortToSerialize
	}

	d.Header.Payload = make([]byte, d.Len()-8)

	var offset = 0
	if d.NetworkAppearance != nil {
		if err := d.NetworkAppearance.SerializeTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += d.NetworkAppearance.Len()
	}

	if d.RoutingContext != nil {
		if err := d.RoutingContext.SerializeTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += d.RoutingContext.Len()
	}

	if d.ProtocolData != nil {
		if err := d.ProtocolData.SerializeTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += d.ProtocolData.Len()
	}

	if d.CorrelationID != nil {
		if err := d.CorrelationID.SerializeTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
	}

	return d.Header.SerializeTo(b)
}

// DecodeData decodes given byte sequence as a Data.
func DecodeData(b []byte) (*Data, error) {
	d := &Data{}
	if err := d.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return d, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (d *Data) DecodeFromBytes(b []byte) error {
	var err error
	d.Header, err = DecodeHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.DecodeMultiParams(d.Header.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode Params")
	}
	for _, pr := range prs {
		switch pr.Tag {
		case params.NetworkAppearance:
			d.NetworkAppearance = pr
		case params.RoutingContext:
			d.RoutingContext = pr
		case params.ProtocolData:
			d.ProtocolData = pr
		case params.CorrelationID:
			d.CorrelationID = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (d *Data) SetLength() {
	if d.NetworkAppearance != nil {
		d.NetworkAppearance.SetLength()
	}
	if d.RoutingContext != nil {
		d.RoutingContext.SetLength()
	}
	if d.ProtocolData != nil {
		d.ProtocolData.SetLength()
	}
	if d.CorrelationID != nil {
		d.CorrelationID.SetLength()
	}

	d.Header.SetLength()
	d.Header.Length += uint32(d.Len())
}

// Len returns the actual length of Data.
func (d *Data) Len() int {
	l := 8 + d.ProtocolData.Len()
	if d.NetworkAppearance != nil {
		l += d.NetworkAppearance.Len()
	}
	if d.RoutingContext != nil {
		l += d.RoutingContext.Len()
	}
	if d.CorrelationID != nil {
		l += d.CorrelationID.Len()
	}
	return l
}

// String returns the Data values in human readable format.
func (d *Data) String() string {
	return fmt.Sprintf("{Header: %s, NetworkAppearance: %s, RoutingContext: %s, ProtocolData %s, CorrelationID: %s}",
		d.Header.String(),
		d.NetworkAppearance.String(),
		d.RoutingContext.String(),
		d.ProtocolData.String(),
		d.CorrelationID.String(),
	)
}

// Version returns the version of M3UA in int.
func (d *Data) Version() uint8 {
	return d.Header.Version
}

// MessageType returns the message type in int.
func (d *Data) MessageType() uint8 {
	return MsgTypePayloadData
}

// MessageClass returns the message class in int.
func (d *Data) MessageClass() uint8 {
	return MsgClassTransfer
}

// MessageClassName returns the name of message class.
func (d *Data) MessageClassName() string {
	return "Transfer"
}

// MessageTypeName returns the name of message type.
func (d *Data) MessageTypeName() string {
	return "Payload Data"
}

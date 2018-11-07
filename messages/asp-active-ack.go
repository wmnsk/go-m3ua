// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// AspActiveAck is a AspActiveAck type of M3UA message.
type AspActiveAck struct {
	*Header
	TrafficModeType *params.Param
	RoutingContext  *params.Param
	InfoString      *params.Param
}

// NewAspActiveAck creates a new AspActiveAck.
func NewAspActiveAck(tmt, rtCtx, info *params.Param) *AspActiveAck {
	a := &AspActiveAck{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPTM,
			Type:     MsgTypeAspActiveAck,
		},
		TrafficModeType: tmt,
		RoutingContext:  rtCtx,
		InfoString:      info,
	}
	a.SetLength()

	return a
}

// Serialize returns the byte sequence generated from a AspActiveAck.
func (a *AspActiveAck) Serialize() ([]byte, error) {
	b := make([]byte, a.Len())
	if err := a.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize AspActiveAck")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (a *AspActiveAck) SerializeTo(b []byte) error {
	if len(b) < a.Len() {
		return ErrTooShortToSerialize
	}

	a.Header.Payload = make([]byte, a.Len()-8)

	var offset = 0
	if a.TrafficModeType != nil {
		if err := a.TrafficModeType.SerializeTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += a.TrafficModeType.Len()
	}

	if a.RoutingContext != nil {
		if err := a.RoutingContext.SerializeTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += a.RoutingContext.Len()
	}

	if a.InfoString != nil {
		if err := a.InfoString.SerializeTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
	}

	return a.Header.SerializeTo(b)
}

// DecodeAspActiveAck decodes given byte sequence as a AspActiveAck.
func DecodeAspActiveAck(b []byte) (*AspActiveAck, error) {
	a := &AspActiveAck{}
	if err := a.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return a, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspActiveAck) DecodeFromBytes(b []byte) error {
	var err error
	a.Header, err = DecodeHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.DecodeMultiParams(a.Header.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode Params")
	}
	for _, pr := range prs {
		switch pr.Tag {
		case params.TrafficModeType:
			a.TrafficModeType = pr
		case params.RoutingContext:
			a.RoutingContext = pr
		case params.InfoString:
			a.InfoString = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (a *AspActiveAck) SetLength() {
	if a.TrafficModeType != nil {
		a.TrafficModeType.SetLength()
	}
	if a.RoutingContext != nil {
		a.RoutingContext.SetLength()
	}
	if a.InfoString != nil {
		a.InfoString.SetLength()
	}

	a.Header.SetLength()
	a.Header.Length += uint32(a.Len())
}

// Len returns the actual length of AspActiveAck.
func (a *AspActiveAck) Len() int {
	l := 8
	if a.TrafficModeType != nil {
		l += a.TrafficModeType.Len()
	}
	if a.RoutingContext != nil {
		l += a.RoutingContext.Len()
	}
	if a.InfoString != nil {
		l += a.InfoString.Len()
	}
	return l
}

// String returns the AspActiveAck values in human readable format.
func (a *AspActiveAck) String() string {
	return fmt.Sprintf("{Header: %s, TrafficModeType: %s, RoutingContext: %s, InfoString: %s}",
		a.Header.String(),
		a.TrafficModeType.String(),
		a.RoutingContext.String(),
		a.InfoString.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *AspActiveAck) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *AspActiveAck) MessageType() uint8 {
	return MsgTypeAspActiveAck
}

// MessageClass returns the message class in int.
func (a *AspActiveAck) MessageClass() uint8 {
	return MsgClassASPTM
}

// MessageClassName returns the name of message class.
func (a *AspActiveAck) MessageClassName() string {
	return "ASPTM"
}

// MessageTypeName returns the name of message type.
func (a *AspActiveAck) MessageTypeName() string {
	return "ASP Active Ack"
}

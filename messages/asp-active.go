// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// AspActive is a AspActive type of M3UA message.
type AspActive struct {
	*Header
	TrafficModeType *params.Param
	RoutingContext  *params.Param
	InfoString      *params.Param
}

// NewAspActive creates a new AspActive.
func NewAspActive(tmt, rtCtx, info *params.Param) *AspActive {
	a := &AspActive{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPTM,
			Type:     MsgTypeAspActive,
		},
		TrafficModeType: tmt,
		RoutingContext:  rtCtx,
		InfoString:      info,
	}
	a.SetLength()

	return a
}

// Serialize returns the byte sequence generated from a AspActive.
func (a *AspActive) Serialize() ([]byte, error) {
	b := make([]byte, a.Len())
	if err := a.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize AspActive")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (a *AspActive) SerializeTo(b []byte) error {
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
		offset += a.InfoString.Len()
	}

	return a.Header.SerializeTo(b)
}

// DecodeAspActive decodes given byte sequence as a AspActive.
func DecodeAspActive(b []byte) (*AspActive, error) {
	a := &AspActive{}
	if err := a.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return a, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspActive) DecodeFromBytes(b []byte) error {
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
func (a *AspActive) SetLength() {
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

// Len returns the actual length of AspActive.
func (a *AspActive) Len() int {
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

// String returns the AspActive values in human readable format.
func (a *AspActive) String() string {
	return fmt.Sprintf("{Header: %s, TrafficModeType: %s, RoutingContext: %s, InfoString: %s}",
		a.Header.String(),
		a.TrafficModeType.String(),
		a.RoutingContext.String(),
		a.InfoString.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *AspActive) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *AspActive) MessageType() uint8 {
	return MsgTypeAspActive
}

// MessageClass returns the message class in int.
func (a *AspActive) MessageClass() uint8 {
	return MsgClassASPTM
}

// MessageClassName returns the name of message class.
func (a *AspActive) MessageClassName() string {
	return "ASPTM"
}

// MessageTypeName returns the name of message type.
func (a *AspActive) MessageTypeName() string {
	return "ASP Active"
}

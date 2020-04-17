// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// AspActiveAck is a AspActiveAck type of M3UA message.
//
// Spec: 3.7.2, RFC4666.
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

// MarshalBinary returns the byte sequence generated from a AspActiveAck.
func (a *AspActiveAck) MarshalBinary() ([]byte, error) {
	b := make([]byte, a.MarshalLen())
	if err := a.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize AspActiveAck")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (a *AspActiveAck) MarshalTo(b []byte) error {
	if len(b) < a.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	a.Header.Payload = make([]byte, a.MarshalLen()-8)

	var offset = 0
	if param := a.TrafficModeType; param != nil {
		if err := param.MarshalTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += param.MarshalLen()
	}

	if param := a.RoutingContext; param != nil {
		if err := param.MarshalTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += param.MarshalLen()
	}

	if param := a.InfoString; param != nil {
		if err := param.MarshalTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
	}

	return a.Header.MarshalTo(b)
}

// ParseAspActiveAck decodes given byte sequence as a AspActiveAck.
func ParseAspActiveAck(b []byte) (*AspActiveAck, error) {
	a := &AspActiveAck{}
	if err := a.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return a, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspActiveAck) UnmarshalBinary(b []byte) error {
	var err error
	a.Header, err = ParseHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.ParseMultiParams(a.Header.Payload)
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
	if param := a.TrafficModeType; param != nil {
		param.SetLength()
	}
	if param := a.RoutingContext; param != nil {
		param.SetLength()
	}
	if param := a.InfoString; param != nil {
		param.SetLength()
	}

	a.Header.Length = uint32(a.MarshalLen())
}

// MarshalLen returns the serial length of AspActiveAck.
func (a *AspActiveAck) MarshalLen() int {
	l := 8
	if param := a.TrafficModeType; param != nil {
		l += param.MarshalLen()
	}
	if param := a.RoutingContext; param != nil {
		l += param.MarshalLen()
	}
	if param := a.InfoString; param != nil {
		l += param.MarshalLen()
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

// Serialize returns the byte sequence generated from a AspActiveAck.
//
// DEPRECATED: use MarshalBinary instead.
func (a *AspActiveAck) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return a.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (a *AspActiveAck) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return a.MarshalTo(b)
}

// DecodeAspActiveAck decodes given byte sequence as a AspActiveAck.
//
// DEPRECATED: use ParseAspActiveAck instead.
func DecodeAspActiveAck(b []byte) (*AspActiveAck, error) {
	log.Println("DEPRECATED: use ParseAspActiveAck instead")
	return ParseAspActiveAck(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (a *AspActiveAck) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return a.UnmarshalBinary(b)
}

// Len returns the serial length of AspActiveAck.
//
// DEPRECATED: use MarshalLen instead.
func (a *AspActiveAck) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return a.MarshalLen()
}

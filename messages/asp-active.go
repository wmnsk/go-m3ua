// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"
	"log"

	"github.com/wmnsk/go-m3ua/messages/params"
)

// AspActive is a AspActive type of M3UA message.
//
// Spec: 3.7.1, RFC4666.
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

// MarshalBinary returns the byte sequence generated from a AspActive.
func (a *AspActive) MarshalBinary() ([]byte, error) {
	b := make([]byte, a.MarshalLen())
	if err := a.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (a *AspActive) MarshalTo(b []byte) error {
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

// ParseAspActive decodes given byte sequence as a AspActive.
func ParseAspActive(b []byte) (*AspActive, error) {
	a := &AspActive{}
	if err := a.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return a, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspActive) UnmarshalBinary(b []byte) error {
	var err error
	a.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}

	prs, err := params.ParseMultiParams(a.Header.Payload)
	if err != nil {
		return err
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

// MarshalLen returns the serial length of AspActive.
func (a *AspActive) MarshalLen() int {
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
	return MsgClassNameASPTM
}

// MessageTypeName returns the name of message type.
func (a *AspActive) MessageTypeName() string {
	return "ASP Active"
}

// Serialize returns the byte sequence generated from a AspActive.
//
// DEPRECATED: use MarshalBinary instead.
func (a *AspActive) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return a.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (a *AspActive) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return a.MarshalTo(b)
}

// DecodeAspActive decodes given byte sequence as a AspActive.
//
// DEPRECATED: use ParseAspActive instead.
func DecodeAspActive(b []byte) (*AspActive, error) {
	log.Println("DEPRECATED: use ParseAspActive instead")
	return ParseAspActive(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (a *AspActive) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return a.UnmarshalBinary(b)
}

// Len returns the serial length of AspActive.
//
// DEPRECATED: use MarshalLen instead.
func (a *AspActive) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return a.MarshalLen()
}

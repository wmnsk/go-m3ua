// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"
	"log"

	"github.com/wmnsk/go-m3ua/messages/params"
)

// AspInactive is a AspInactive type of M3UA message.
//
// Spec: 3.7.3, RFC4666.
type AspInactive struct {
	*Header
	RoutingContext *params.Param
	InfoString     *params.Param
}

// NewAspInactive creates a new AspInactive.
func NewAspInactive(rtCtx, info *params.Param) *AspInactive {
	a := &AspInactive{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPTM,
			Type:     MsgTypeAspInactive,
		},
		RoutingContext: rtCtx,
		InfoString:     info,
	}
	a.SetLength()

	return a
}

// MarshalBinary returns the byte sequence generated from a AspInactive.
func (a *AspInactive) MarshalBinary() ([]byte, error) {
	b := make([]byte, a.MarshalLen())
	if err := a.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (a *AspInactive) MarshalTo(b []byte) error {
	if len(b) < a.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	a.Header.Payload = make([]byte, a.MarshalLen()-8)

	var offset = 0
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

// ParseAspInactive decodes given byte sequence as a AspInactive.
func ParseAspInactive(b []byte) (*AspInactive, error) {
	a := &AspInactive{}
	if err := a.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return a, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspInactive) UnmarshalBinary(b []byte) error {
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
func (a *AspInactive) SetLength() {
	if param := a.RoutingContext; param != nil {
		param.SetLength()
	}
	if param := a.InfoString; param != nil {
		param.SetLength()
	}

	a.Header.Length = uint32(a.MarshalLen())
}

// MarshalLen returns the serial length of AspInactive.
func (a *AspInactive) MarshalLen() int {
	l := 8
	if param := a.RoutingContext; param != nil {
		l += param.MarshalLen()
	}
	if param := a.InfoString; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// String returns the AspInactive values in human readable format.
func (a *AspInactive) String() string {
	return fmt.Sprintf("{Header: %s, RoutingContext: %s, InfoString: %s}",
		a.Header.String(),
		a.RoutingContext.String(),
		a.InfoString.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *AspInactive) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *AspInactive) MessageType() uint8 {
	return MsgTypeAspInactive
}

// MessageClass returns the message class in int.
func (a *AspInactive) MessageClass() uint8 {
	return MsgClassASPTM
}

// MessageClassName returns the name of message class.
func (a *AspInactive) MessageClassName() string {
	return MsgClassNameASPTM
}

// MessageTypeName returns the name of message type.
func (a *AspInactive) MessageTypeName() string {
	return "ASP Inactive"
}

// Serialize returns the byte sequence generated from a AspInactive.
//
// DEPRECATED: use MarshalBinary instead.
func (a *AspInactive) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return a.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (a *AspInactive) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return a.MarshalTo(b)
}

// DecodeAspInactive decodes given byte sequence as a AspInactive.
//
// DEPRECATED: use ParseAspInactive instead.
func DecodeAspInactive(b []byte) (*AspInactive, error) {
	log.Println("DEPRECATED: use ParseAspInactive instead")
	return ParseAspInactive(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (a *AspInactive) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return a.UnmarshalBinary(b)
}

// Len returns the serial length of AspInactive.
//
// DEPRECATED: use MarshalLen instead.
func (a *AspInactive) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return a.MarshalLen()
}

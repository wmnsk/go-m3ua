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

// AspDownAck is a AspDownAck type of M3UA message.
//
// Spec: 3.5.4, RFC4666.
type AspDownAck struct {
	*Header
	AspIdentifier *params.Param
	InfoString    *params.Param
}

// NewAspDownAck creates a new AspDownAck.
func NewAspDownAck(info *params.Param) *AspDownAck {
	a := &AspDownAck{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPSM,
			Type:     MsgTypeAspDownAck,
		},
		InfoString: info,
	}
	a.SetLength()

	return a
}

// MarshalBinary returns the byte sequence generated from a AspDownAck.
func (a *AspDownAck) MarshalBinary() ([]byte, error) {
	b := make([]byte, a.MarshalLen())
	if err := a.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize AspDownAck")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (a *AspDownAck) MarshalTo(b []byte) error {
	if len(b) < a.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	a.Header.Payload = make([]byte, a.MarshalLen()-8)

	if param := a.InfoString; param != nil {
		if err := param.MarshalTo(a.Header.Payload); err != nil {
			return err
		}
	}

	return a.Header.MarshalTo(b)
}

// ParseAspDownAck decodes given byte sequence as a AspDownAck.
func ParseAspDownAck(b []byte) (*AspDownAck, error) {
	a := &AspDownAck{}
	if err := a.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return a, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspDownAck) UnmarshalBinary(b []byte) error {
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
		case params.InfoString:
			a.InfoString = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (a *AspDownAck) SetLength() {
	if param := a.InfoString; param != nil {
		param.SetLength()
	}

	a.Header.Length = uint32(a.MarshalLen())
}

// MarshalLen returns the serial length of AspDownAck.
func (a *AspDownAck) MarshalLen() int {
	l := 8
	if param := a.InfoString; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// String returns the AspDownAck values in human readable format.
func (a *AspDownAck) String() string {
	return fmt.Sprintf("{Header: %s, InfoString: %s}",
		a.Header.String(),
		a.InfoString.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *AspDownAck) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *AspDownAck) MessageType() uint8 {
	return MsgTypeAspDownAck
}

// MessageClass returns the message class in int.
func (a *AspDownAck) MessageClass() uint8 {
	return MsgClassASPSM
}

// MessageClassName returns the name of message class.
func (a *AspDownAck) MessageClassName() string {
	return "ASPSM"
}

// MessageTypeName returns the name of message type.
func (a *AspDownAck) MessageTypeName() string {
	return "ASP Down Ack"
}

// Serialize returns the byte sequence generated from a AspDownAck.
//
// DEPRECATED: use MarshalBinary instead.
func (a *AspDownAck) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return a.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (a *AspDownAck) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return a.MarshalTo(b)
}

// DecodeAspDownAck decodes given byte sequence as a AspDownAck.
//
// DEPRECATED: use ParseAspDownAck instead.
func DecodeAspDownAck(b []byte) (*AspDownAck, error) {
	log.Println("DEPRECATED: use ParseAspDownAck instead")
	return ParseAspDownAck(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (a *AspDownAck) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return a.UnmarshalBinary(b)
}

// Len returns the serial length of AspDownAck.
//
// DEPRECATED: use MarshalLen instead.
func (a *AspDownAck) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return a.MarshalLen()
}

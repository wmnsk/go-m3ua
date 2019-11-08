// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// AspUpAck is a AspUpAck type of M3UA message.
//
// Spec: 3.5.2, RFC4666.
type AspUpAck struct {
	*Header
	AspIdentifier *params.Param
	InfoString    *params.Param
}

// NewAspUpAck creates a new AspUpAck.
func NewAspUpAck(aspID, info *params.Param) *AspUpAck {
	a := &AspUpAck{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPSM,
			Type:     MsgTypeAspUpAck,
		},
		AspIdentifier: aspID,
		InfoString:    info,
	}
	a.SetLength()

	return a
}

// MarshalBinary returns the byte sequence generated from a AspUpAck.
func (a *AspUpAck) MarshalBinary() ([]byte, error) {
	b := make([]byte, a.MarshalLen())
	if err := a.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize AspUpAck")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (a *AspUpAck) MarshalTo(b []byte) error {
	if len(b) < a.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	a.Header.Payload = make([]byte, a.MarshalLen()-8)

	var offset = 0
	if param := a.AspIdentifier; param != nil {
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

// ParseAspUpAck decodes given byte sequence as a AspUpAck.
func ParseAspUpAck(b []byte) (*AspUpAck, error) {
	a := &AspUpAck{}
	if err := a.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return a, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspUpAck) UnmarshalBinary(b []byte) error {
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
		case params.AspIdentifier:
			a.AspIdentifier = pr
		case params.InfoString:
			a.InfoString = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (a *AspUpAck) SetLength() {
	if param := a.AspIdentifier; param != nil {
		param.SetLength()
	}
	if param := a.InfoString; param != nil {
		param.SetLength()
	}

	a.Header.SetLength()
	a.Header.Length += uint32(a.MarshalLen())
}

// MarshalLen returns the serial length of AspUpAck.
func (a *AspUpAck) MarshalLen() int {
	l := 8
	if param := a.AspIdentifier; param != nil {
		l += param.MarshalLen()
	}
	if param := a.InfoString; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// String returns the AspUpAck values in human readable format.
func (a *AspUpAck) String() string {
	return fmt.Sprintf("{Header: %s, AspIdentifier: %s, InfoString: %s}",
		a.Header.String(),
		a.AspIdentifier.String(),
		a.InfoString.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *AspUpAck) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *AspUpAck) MessageType() uint8 {
	return MsgTypeAspUpAck
}

// MessageClass returns the message class in int.
func (a *AspUpAck) MessageClass() uint8 {
	return MsgClassASPSM
}

// MessageClassName returns the name of message class.
func (a *AspUpAck) MessageClassName() string {
	return "ASPSM"
}

// MessageTypeName returns the name of message type.
func (a *AspUpAck) MessageTypeName() string {
	return "ASP Up Ack"
}

// Serialize returns the byte sequence generated from a AspUpAck.
//
// DEPRECATED: use MarshalBinary instead.
func (a *AspUpAck) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return a.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (a *AspUpAck) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return a.MarshalTo(b)
}

// DecodeAspUpAck decodes given byte sequence as a AspUpAck.
//
// DEPRECATED: use ParseAspUpAck instead.
func DecodeAspUpAck(b []byte) (*AspUpAck, error) {
	log.Println("DEPRECATED: use ParseAspUpAck instead")
	return ParseAspUpAck(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (a *AspUpAck) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return a.UnmarshalBinary(b)
}

// Len returns the serial length of AspUpAck.
//
// DEPRECATED: use MarshalLen instead.
func (a *AspUpAck) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return a.MarshalLen()
}

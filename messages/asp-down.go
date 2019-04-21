// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// AspDown is a AspDown type of M3UA message.
type AspDown struct {
	*Header
	AspIdentifier *params.Param
	InfoString    *params.Param
}

// NewAspDown creates a new AspDown.
func NewAspDown(info *params.Param) *AspDown {
	a := &AspDown{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPSM,
			Type:     MsgTypeAspDown,
		},
		InfoString: info,
	}
	a.SetLength()

	return a
}

// Serialize returns the byte sequence generated from a AspDown.
func (a *AspDown) Serialize() ([]byte, error) {
	b := make([]byte, a.Len())
	if err := a.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize AspDown")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (a *AspDown) SerializeTo(b []byte) error {
	if len(b) < a.Len() {
		return ErrTooShortToSerialize
	}

	a.Header.Payload = make([]byte, a.Len()-8)

	if a.InfoString != nil {
		if err := a.InfoString.SerializeTo(a.Header.Payload); err != nil {
			return err
		}
	}

	return a.Header.SerializeTo(b)
}

// DecodeAspDown decodes given byte sequence as a AspDown.
func DecodeAspDown(b []byte) (*AspDown, error) {
	a := &AspDown{}
	if err := a.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return a, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspDown) DecodeFromBytes(b []byte) error {
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
		case params.InfoString:
			a.InfoString = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (a *AspDown) SetLength() {
	if a.InfoString != nil {
		a.InfoString.SetLength()
	}

	a.Header.SetLength()
	a.Header.Length += uint32(a.Len())
}

// Len returns the actual length of AspDown.
func (a *AspDown) Len() int {
	l := 8
	if a.InfoString != nil {
		l += a.InfoString.Len()
	}
	return l
}

// String returns the AspDown values in human readable format.
func (a *AspDown) String() string {
	return fmt.Sprintf("{Header: %s, InfoString: %s}",
		a.Header.String(),
		a.InfoString.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *AspDown) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *AspDown) MessageType() uint8 {
	return MsgTypeAspDown
}

// MessageClass returns the message class in int.
func (a *AspDown) MessageClass() uint8 {
	return MsgClassASPSM
}

// MessageClassName returns the name of message class.
func (a *AspDown) MessageClassName() string {
	return "ASPSM"
}

// MessageTypeName returns the name of message type.
func (a *AspDown) MessageTypeName() string {
	return "ASP Down"
}

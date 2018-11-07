// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// AspUp is a AspUp type of M3UA message.
type AspUp struct {
	*Header
	AspIdentifier *params.Param
	InfoString    *params.Param
}

// NewAspUp creates a new AspUp.
func NewAspUp(aspID, info *params.Param) *AspUp {
	a := &AspUp{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPSM,
			Type:     MsgTypeAspUp,
		},
		AspIdentifier: aspID,
		InfoString:    info,
	}
	a.SetLength()

	return a
}

// Serialize returns the byte sequence generated from a AspUp.
func (a *AspUp) Serialize() ([]byte, error) {
	b := make([]byte, a.Len())
	if err := a.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize AspUp")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (a *AspUp) SerializeTo(b []byte) error {
	if len(b) < a.Len() {
		return ErrTooShortToSerialize
	}

	a.Header.Payload = make([]byte, a.Len()-8)

	var offset = 0
	if a.AspIdentifier != nil {
		if err := a.AspIdentifier.SerializeTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += a.AspIdentifier.Len()
	}

	if a.InfoString != nil {
		if err := a.InfoString.SerializeTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
	}

	return a.Header.SerializeTo(b)
}

// DecodeAspUp decodes given byte sequence as a AspUp.
func DecodeAspUp(b []byte) (*AspUp, error) {
	a := &AspUp{}
	if err := a.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return a, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (a *AspUp) DecodeFromBytes(b []byte) error {
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
func (a *AspUp) SetLength() {
	if a.AspIdentifier != nil {
		a.AspIdentifier.SetLength()
	}
	if a.InfoString != nil {
		a.InfoString.SetLength()
	}

	a.Header.SetLength()
	a.Header.Length += uint32(a.Len())
}

// Len returns the actual length of AspUp.
func (a *AspUp) Len() int {
	l := 8
	if a.AspIdentifier != nil {
		l += a.AspIdentifier.Len()
	}
	if a.InfoString != nil {
		l += a.InfoString.Len()
	}
	return l
}

// String returns the AspUp values in human readable format.
func (a *AspUp) String() string {
	return fmt.Sprintf("{Header: %s, AspIdentifier: %s, InfoString: %s}",
		a.Header.String(),
		a.AspIdentifier.String(),
		a.InfoString.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *AspUp) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *AspUp) MessageType() uint8 {
	return MsgTypeAspUp
}

// MessageClass returns the message class in int.
func (a *AspUp) MessageClass() uint8 {
	return MsgClassASPSM
}

// MessageClassName returns the name of message class.
func (a *AspUp) MessageClassName() string {
	return "ASPSM"
}

// MessageTypeName returns the name of message type.
func (a *AspUp) MessageTypeName() string {
	return "ASP Up"
}

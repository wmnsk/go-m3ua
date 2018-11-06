// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// Error is a Error type of M3UA message.
type Error struct {
	*Header
	ErrorCode             *params.Param
	RoutingContext        *params.Param
	NetworkAppearance     *params.Param
	AffectedPointCode     *params.Param
	DiagnosticInformation *params.Param
}

// NewError creates a new Error.
func NewError(code, rtCtx, nwApr, apc, info *params.Param) *Error {
	e := &Error{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassManagement,
			Type:     MsgTypeError,
		},
		ErrorCode:             code,
		RoutingContext:        rtCtx,
		NetworkAppearance:     nwApr,
		AffectedPointCode:     apc,
		DiagnosticInformation: info,
	}
	e.SetLength()

	return e
}

// Serialize returns the byte sequence generated from a Error.
func (e *Error) Serialize() ([]byte, error) {
	b := make([]byte, e.Len())
	if err := e.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize Error")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (e *Error) SerializeTo(b []byte) error {
	if len(b) < e.Len() {
		return ErrTooShortToSerialize
	}

	e.Header.Payload = make([]byte, e.Len()-8)

	var offset = 0

	if e.ErrorCode != nil {
		if err := e.ErrorCode.SerializeTo(e.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += e.ErrorCode.Len()
	}

	if e.RoutingContext != nil {
		if err := e.RoutingContext.SerializeTo(e.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += e.RoutingContext.Len()
	}

	if e.NetworkAppearance != nil {
		if err := e.NetworkAppearance.SerializeTo(e.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += e.NetworkAppearance.Len()
	}

	if e.AffectedPointCode != nil {
		if err := e.AffectedPointCode.SerializeTo(e.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += e.AffectedPointCode.Len()
	}

	if e.DiagnosticInformation != nil {
		if err := e.DiagnosticInformation.SerializeTo(e.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += e.DiagnosticInformation.Len()
	}

	return e.Header.SerializeTo(b)
}

// DecodeError decodes given byte sequence as a Error.
func DecodeError(b []byte) (*Error, error) {
	e := &Error{}
	if err := e.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return e, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (e *Error) DecodeFromBytes(b []byte) error {
	var err error
	e.Header, err = DecodeHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.DecodeMultiParams(e.Header.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode Params")
	}
	for _, pr := range prs {
		switch pr.Tag {
		case params.ErrorCode:
			e.ErrorCode = pr
		case params.RoutingContext:
			e.RoutingContext = pr
		case params.NetworkAppearance:
			e.NetworkAppearance = pr
		case params.AffectedPointCode:
			e.AffectedPointCode = pr
		case params.DiagnosticInformation:
			e.DiagnosticInformation = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (e *Error) SetLength() {
	if e.ErrorCode != nil {
		e.ErrorCode.SetLength()
	}
	if e.RoutingContext != nil {
		e.RoutingContext.SetLength()
	}
	if e.NetworkAppearance != nil {
		e.NetworkAppearance.SetLength()
	}
	if e.AffectedPointCode != nil {
		e.AffectedPointCode.SetLength()
	}
	if e.DiagnosticInformation != nil {
		e.DiagnosticInformation.SetLength()
	}

	e.Header.SetLength()
	e.Header.Length += uint32(e.Len())
}

// Len returns the actual length of Error.
func (e *Error) Len() int {
	l := 8
	if e.ErrorCode != nil {
		l += e.ErrorCode.Len()
	}
	if e.RoutingContext != nil {
		l += e.RoutingContext.Len()
	}
	if e.NetworkAppearance != nil {
		l += e.NetworkAppearance.Len()
	}
	if e.AffectedPointCode != nil {
		l += e.AffectedPointCode.Len()
	}
	if e.DiagnosticInformation != nil {
		l += e.DiagnosticInformation.Len()
	}
	return l
}

// String returns the Error values in human readable format.
func (e *Error) String() string {
	return fmt.Sprintf("{Header: %s, ErrorCode: %s, RoutingContext: %s, NetworkAppearance: %s, AffectedPointCode: %s, DiagnosticInformation: %s}",
		e.Header.String(),
		e.ErrorCode.String(),
		e.RoutingContext.String(),
		e.NetworkAppearance.String(),
		e.AffectedPointCode.String(),
		e.DiagnosticInformation.String(),
	)
}

// Version returns the version of M3UA in int.
func (e *Error) Version() uint8 {
	return e.Header.Version
}

// MessageType returns the message type in int.
func (e *Error) MessageType() uint8 {
	return MsgTypeError
}

// MessageClass returns the message class in int.
func (e *Error) MessageClass() uint8 {
	return MsgClassManagement
}

// MessageClassName returns the name of message class.
func (e *Error) MessageClassName() string {
	return "Management"
}

// MessageTypeName returns the name of message type.
func (e *Error) MessageTypeName() string {
	return "Error"
}

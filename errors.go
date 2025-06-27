// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"errors"
	"fmt"

	"github.com/wmnsk/go-m3ua/messages"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// Error definitions.
var (
	ErrSCTPNotAlive        = errors.New("SCTP is no longer alive")
	ErrInvalidState        = errors.New("invalid state")
	ErrNotEstablished      = errors.New("M3UA Conn not established")
	ErrFailedToEstablish   = errors.New("failed to establish M3UA Conn")
	ErrTimeout             = errors.New("timed out")
	ErrHeartbeatExpired    = errors.New("heartbeat timer expired")
	ErrFailedToPeelOff     = errors.New("failed to peel off Protocol Data")
	ErrFailedToWriteSignal = errors.New("failed to write signal")

	// ErrAspIDRequired is used by an SGP in response to an ASP Up message that
	// does not contain an ASP Identifier parameter when the SGP requires one.
	ErrAspIDRequired = errors.New("ASP Identifier required")
)

// InvalidVersionError is used if a message with an unsupported version is received.
type InvalidVersionError struct {
	Ver uint8
}

// NewInvalidVersionError creates InvalidVersionError.
func NewInvalidVersionError(ver uint8) *InvalidVersionError {
	return &InvalidVersionError{Ver: ver}
}

// Error returns error string with violating version.
func (e *InvalidVersionError) Error() string {
	return fmt.Sprintf("invalid version: %d", e.Ver)
}

// UnsupportedClassError is used if a message with an unexpected or
// unsupported Message Class is received.
type UnsupportedClassError struct {
	Msg messages.M3UA
}

// NewUnsupportedClassError creates UnsupportedClassError
func NewUnsupportedClassError(msg messages.M3UA) *UnsupportedClassError {
	return &UnsupportedClassError{Msg: msg}
}

// Error returns error string with message class.
func (e *UnsupportedClassError) Error() string {
	return fmt.Sprintf("message class unsupported. class: %s", e.Msg.MessageClassName())
}

func (e *UnsupportedClassError) first40Octets() []byte {
	b, err := e.Msg.MarshalBinary()
	if err != nil {
		return nil
	}
	if len(b) < 40 {
		return b
	}

	return b[:40]
}

// UnsupportedMessageError is used if a message with an
// unexpected or unsupported Message Type is received.
type UnsupportedMessageError struct {
	Msg messages.M3UA
}

// NewUnsupportedMessageError creates UnsupportedMessageError
func NewUnsupportedMessageError(msg messages.M3UA) *UnsupportedMessageError {
	return &UnsupportedMessageError{Msg: msg}
}

// Error returns error string with message class and type.
func (e *UnsupportedMessageError) Error() string {
	return fmt.Sprintf("message unsupported. class: %s, type: %s", e.Msg.MessageClassName(), e.Msg.MessageTypeName())
}

func (e *UnsupportedMessageError) first40Octets() []byte {
	b, err := e.Msg.MarshalBinary()
	if err != nil {
		return nil
	}
	if len(b) < 40 {
		return b
	}

	return b[:40]
}

// UnexpectedMessageError is used if a defined and recognized message is received
// that is not expected in the current state (in some cases, the ASP may optionally
// silently discard the message and not send an Error message).
type UnexpectedMessageError struct {
	Msg messages.M3UA
}

// NewUnexpectedMessageError creates UnexpectedMessageError
func NewUnexpectedMessageError(msg messages.M3UA) *UnexpectedMessageError {
	return &UnexpectedMessageError{Msg: msg}
}

// Error returns error string with message class and type.
func (e *UnexpectedMessageError) Error() string {
	return fmt.Sprintf("unexpected message. class: %s, type: %s", e.Msg.MessageClassName(), e.Msg.MessageTypeName())
}

// InvalidSCTPStreamIDError is used if a message is received on an unexpected SCTP stream.
type InvalidSCTPStreamIDError struct {
	ID uint16
}

// NewInvalidSCTPStreamIDError creates InvalidSCTPStreamIDError
func NewInvalidSCTPStreamIDError(id uint16) *InvalidSCTPStreamIDError {
	return &InvalidSCTPStreamIDError{ID: id}
}

// Error returns error string with violating stream ID.
func (e *InvalidSCTPStreamIDError) Error() string {
	return fmt.Sprintf("invalid SCTP Stream ID: %d", e.ID)
}

func (c *Conn) handleErrors(e error) error {
	var res messages.M3UA
	var InvalidVersionError *InvalidVersionError
	if errors.As(e, &InvalidVersionError) {
		res = messages.NewError(
			params.NewErrorCode(params.InvalidVersionError),
			nil, nil, nil, nil,
		)
	}
	//nolint:errorlint
	if err, ok := e.(*UnsupportedClassError); ok {
		res = messages.NewError(
			params.NewErrorCode(params.UnsupportedMessageErrorClass),
			nil, nil, nil,
			params.NewDiagnosticInformation(err.first40Octets()),
		)
	}
	//nolint:errorlint
	if err, ok := e.(*UnsupportedMessageError); ok {
		res = messages.NewError(
			params.NewErrorCode(params.UnsupportedMessageErrorType),
			nil, nil, nil,
			params.NewDiagnosticInformation(err.first40Octets()),
		)
	}
	var UnexpectedMessageError *UnexpectedMessageError
	if errors.As(e, &UnexpectedMessageError) {
		res = messages.NewError(
			params.NewErrorCode(params.UnexpectedMessageError),
			c.cfg.RoutingContexts,
			c.cfg.NetworkAppearance,
			params.NewAffectedPointCode(
				c.cfg.OriginatingPointCode,
			),
			nil,
		)
	}
	var InvalidSCTPStreamIDError *InvalidSCTPStreamIDError
	if errors.As(e, &InvalidSCTPStreamIDError) {
		res = messages.NewError(
			params.NewErrorCode(params.ErrInvalidStreamIdentifier),
			nil, nil, nil, nil,
		)
	}
	if errors.Is(e, ErrAspIDRequired) {
		res = messages.NewError(
			params.NewErrorCode(params.ErrAspIdentifierRequired),
			nil, nil, nil, nil,
		)
	}

	if res == nil {
		return e
	}

	if _, err := c.WriteSignal(res); err != nil {
		return err
	}

	return nil
}

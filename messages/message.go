// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"log"

	"github.com/pkg/errors"
)

// Message Class definitions.
const (
	MsgClassManagement uint8 = iota
	MsgClassTransfer
	MsgClassSSNM
	MsgClassASPSM
	MsgClassASPTM
	_
	_
	_
	_
	MsgClassRKM
)

// Message Type definitions (Management).
const (
	MsgTypeError = iota
	MsgTypeNotify
)

// Message Type definitions (SSNM).
const (
	_ = iota
	MsgTypeDestinationUnavailable
	MsgTypeDestinationAvailable
	MsgTypeDestinationStateAudit
	MsgTypeSignallingCongestion
	MsgTypeDestinationUserPartUnavailable
	MsgTypeDestinationRestricted
)

// Message Type definitions (Transfer).
const (
	_ uint8 = iota
	MsgTypePayloadData
)

// Message Type definitions (ASPSM).
const (
	_ uint8 = iota
	MsgTypeAspUp
	MsgTypeAspDown
	MsgTypeHeartbeat
	MsgTypeAspUpAck
	MsgTypeAspDownAck
	MsgTypeHeartbeatAck
)

// Message Type definitions (ASPTM).
const (
	_ uint8 = iota
	MsgTypeAspActive
	MsgTypeAspInactive
	MsgTypeAspActiveAck
	MsgTypeAspInactiveAck
)

// Message Type definitions (RKM).
const (
	_ uint8 = iota
	MsgTypeRegistrationRequest
	MsgTypeRegistrationResponse
	MsgTypeDeregistrationRequest
	MsgTypeDeregistrationResponse
)

// M3UA is an interface that defines M3UA messages.
type M3UA interface {
	MarshalBinary() ([]byte, error)
	MarshalTo([]byte) error
	UnmarshalBinary([]byte) error
	MarshalLen() int
	Version() uint8
	MessageClass() uint8
	MessageType() uint8
	MessageClassName() string
	MessageTypeName() string
}

// MarshalBinary returns the byte sequence generated from a M3UA instance.
// Better to use MarshalBinaryXxx instead if you know the type of data to be serialized.
func MarshalBinary(m M3UA) ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// Parse decodes the given bytes.
// This function checks the Message Class and Message Type and chooses the appropriate type.
func Parse(b []byte) (M3UA, error) {
	if len(b) < 4 {
		return nil, ErrTooShortToParse
	}
	var m M3UA
	combine := func(c, t uint8) uint16 {
		return uint16(c<<4 | t)
	}
	t := combine(b[2], b[3])

	switch t {
	// Transfer Messages
	case combine(MsgClassTransfer, MsgTypePayloadData):
		m = &Data{}
		// SSNM Messages
	case combine(MsgClassSSNM, MsgTypeDestinationUnavailable):
		m = &DestinationUnavailable{}
	case combine(MsgClassSSNM, MsgTypeDestinationAvailable):
		m = &DestinationAvailable{}
	case combine(MsgClassSSNM, MsgTypeDestinationStateAudit):
		m = &DestinationStateAudit{}
	case combine(MsgClassSSNM, MsgTypeSignallingCongestion):
		m = &SignallingCongestion{}
	case combine(MsgClassSSNM, MsgTypeDestinationUserPartUnavailable):
		m = &DestinationUserPartUnavailable{}
	case combine(MsgClassSSNM, MsgTypeDestinationRestricted):
		m = &DestinationRestricted{}
		// ASPSM Messages
	case combine(MsgClassASPSM, MsgTypeAspUp):
		m = &AspUp{}
	case combine(MsgClassASPSM, MsgTypeAspDown):
		m = &AspDown{}
	case combine(MsgClassASPSM, MsgTypeHeartbeat):
		m = &Heartbeat{}
	case combine(MsgClassASPSM, MsgTypeAspUpAck):
		m = &AspUpAck{}
	case combine(MsgClassASPSM, MsgTypeAspDownAck):
		m = &AspDownAck{}
	case combine(MsgClassASPSM, MsgTypeHeartbeatAck):
		m = &HeartbeatAck{}
	// ASPTM Messages
	case combine(MsgClassASPTM, MsgTypeAspActive):
		m = &AspActive{}
	case combine(MsgClassASPTM, MsgTypeAspActiveAck):
		m = &AspActiveAck{}
	case combine(MsgClassASPTM, MsgTypeAspInactive):
		m = &AspInactive{}
	case combine(MsgClassASPTM, MsgTypeAspInactiveAck):
		m = &AspInactiveAck{}
	// Management Messages
	case combine(MsgClassManagement, MsgTypeError):
		m = &Error{}
	case combine(MsgClassManagement, MsgTypeNotify):
		m = &Notify{}
	default:
		// If the combination of class and type is unknown or not supported, *Generic is used.
		m = &Generic{}
	}

	if err := m.UnmarshalBinary(b); err != nil {
		return nil, errors.Wrap(err, "failed to decode M3UA")
	}
	return m, nil
}

// Decode decodes the given bytes.
// This function checks the Message Class and Message Type and chooses the appropriate type.
//
// DEPRECATED: use Parse instead.
func Decode(b []byte) (M3UA, error) {
	log.Println("DEPRECATED: use Parse instead")
	return Parse(b)
}

// Error definitions.
var (
	ErrTooShortToMarshalBinary = errors.New("insufficient buffer to serialize M3UA to")
	ErrTooShortToParse         = errors.New("too short to decode as M3UA")
	ErrInvalidParameter        = errors.New("got invalid parameter inside a message")
)

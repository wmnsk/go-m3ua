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

// Notify is a Notify type of M3UA message.
//
// Spec: 3.8.2, RFC4666.
type Notify struct {
	*Header
	Status         *params.Param
	AspIdentifier  *params.Param
	RoutingContext *params.Param
	InfoString     *params.Param
}

// NewNotify creates a new Notify.
func NewNotify(status, aspID, rtCtx, info *params.Param) *Notify {
	n := &Notify{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassManagement,
			Type:     MsgTypeNotify,
		},
		Status:         status,
		AspIdentifier:  aspID,
		RoutingContext: rtCtx,
		InfoString:     info,
	}
	n.SetLength()

	return n
}

// MarshalBinary returns the byte sequence generated from a Notify.
func (n *Notify) MarshalBinary() ([]byte, error) {
	b := make([]byte, n.MarshalLen())
	if err := n.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize Notify")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (n *Notify) MarshalTo(b []byte) error {
	if len(b) < n.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	n.Header.Payload = make([]byte, n.MarshalLen()-8)

	var offset = 0

	// NOTE:
	// Precisely, it should validate whether the `Status` parameter exists or not because
	// the `Status` parameter has to be contained in the `Notify` message.
	// (ref: https://tools.ietf.org/html/rfc4666#section-3.8.2)
	// However, this library aims to be flexible using and/or verifying,
	// so it doesn't check the existence of the parameter for now.
	// Discussion: https://github.com/wmnsk/go-m3ua/pull/10#discussion_r304225571
	if param := n.Status; param != nil {
		if err := param.MarshalTo(n.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += param.MarshalLen()
	}

	if param := n.AspIdentifier; param != nil {
		if err := param.MarshalTo(n.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += param.MarshalLen()
	}

	if param := n.RoutingContext; param != nil {
		if err := param.MarshalTo(n.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += param.MarshalLen()
	}

	if param := n.InfoString; param != nil {
		if err := param.MarshalTo(n.Header.Payload[offset:]); err != nil {
			return err
		}
	}

	return n.Header.MarshalTo(b)
}

// ParseNotify decodes given byte sequence as a Notify.
func ParseNotify(b []byte) (*Notify, error) {
	n := &Notify{}
	if err := n.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return n, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (n *Notify) UnmarshalBinary(b []byte) error {
	var err error
	n.Header, err = ParseHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.ParseMultiParams(n.Header.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode Params")
	}
	for _, pr := range prs {
		switch pr.Tag {
		case params.Status:
			n.Status = pr
		case params.AspIdentifier:
			n.AspIdentifier = pr
		case params.RoutingContext:
			n.RoutingContext = pr
		case params.InfoString:
			n.InfoString = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (n *Notify) SetLength() {
	if param := n.Status; param != nil {
		param.SetLength()
	}
	if param := n.AspIdentifier; param != nil {
		param.SetLength()
	}
	if param := n.RoutingContext; param != nil {
		param.SetLength()
	}
	if param := n.InfoString; param != nil {
		param.SetLength()
	}

	n.Header.SetLength()
	n.Header.Length += uint32(n.MarshalLen())
}

// MarshalLen returns the serial length of Notify.
func (n *Notify) MarshalLen() int {
	l := 8

	if param := n.Status; param != nil {
		l += param.MarshalLen()
	}
	if param := n.AspIdentifier; param != nil {
		l += param.MarshalLen()
	}
	if param := n.RoutingContext; param != nil {
		l += param.MarshalLen()
	}
	if param := n.InfoString; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// String returns the Notify values in human readable format.
func (n *Notify) String() string {
	return fmt.Sprintf("{Header: %s, Status: %s, AspIdentifier: %s, RoutingContext: %s, InfoString: %s}",
		n.Header.String(),
		n.Status.String(),
		n.AspIdentifier.String(),
		n.RoutingContext.String(),
		n.InfoString.String(),
	)
}

// Version returns the version of M3UA in int.
func (n *Notify) Version() uint8 {
	return n.Header.Version
}

// MessageType returns the message type in int.
func (n *Notify) MessageType() uint8 {
	return MsgTypeNotify
}

// MessageClass returns the message class in int.
func (n *Notify) MessageClass() uint8 {
	return MsgClassManagement
}

// MessageClassName returns the name of message class.
func (n *Notify) MessageClassName() string {
	return "Management"
}

// MessageTypeName returns the name of message type.
func (n *Notify) MessageTypeName() string {
	return "Notify"
}

// Serialize returns the byte sequence generated from a Notify.
//
// DEPRECATED: use MarshalBinary instead.
func (n *Notify) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return n.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (n *Notify) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return n.MarshalTo(b)
}

// DecodeNotify decodes given byte sequence as a Notify.
//
// DEPRECATED: use ParseNotify instead.
func DecodeNotify(b []byte) (*Notify, error) {
	log.Println("DEPRECATED: use ParseNotify instead")
	return ParseNotify(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (n *Notify) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return n.UnmarshalBinary(b)
}

// Len returns the serial length of Notify.
//
// DEPRECATED: use MarshalLen instead.
func (n *Notify) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return n.MarshalLen()
}

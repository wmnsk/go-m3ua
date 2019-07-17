// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// Notify is a Notify type of M3UA message.
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

// Serialize returns the byte sequence generated from a Notify.
func (n *Notify) Serialize() ([]byte, error) {
	b := make([]byte, n.Len())
	if err := n.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize Notify")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (n *Notify) SerializeTo(b []byte) error {
	if len(b) < n.Len() {
		return ErrTooShortToSerialize
	}

	n.Header.Payload = make([]byte, n.Len()-8)

	var offset = 0

	// NOTE:
	// Precisely, it should validate whether the `Status` parameter exists or not because
	// the `Status` parameter has to be contained in the `Notify` message.
	// (ref: https://tools.ietf.org/html/rfc3332#section-3.8.2)
	// However, this library aims to be flexible using and/or verifying,
	// so it doesn't check the existence of the parameter for now.
	// Discussion: https://github.com/wmnsk/go-m3ua/pull/10#discussion_r304225571
	if n.Status != nil {
		if err := n.Status.SerializeTo(n.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += n.Status.Len()
	}

	if n.AspIdentifier != nil {
		if err := n.AspIdentifier.SerializeTo(n.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += n.AspIdentifier.Len()
	}

	if n.RoutingContext != nil {
		if err := n.RoutingContext.SerializeTo(n.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += n.RoutingContext.Len()
	}

	if n.InfoString != nil {
		if err := n.InfoString.SerializeTo(n.Header.Payload[offset:]); err != nil {
			return err
		}
	}

	return n.Header.SerializeTo(b)
}

// DecodeNotify decodes given byte sequence as a Notify.
func DecodeNotify(b []byte) (*Notify, error) {
	n := &Notify{}
	if err := n.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return n, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (n *Notify) DecodeFromBytes(b []byte) error {
	var err error
	n.Header, err = DecodeHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.DecodeMultiParams(n.Header.Payload)
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
	if n.Status != nil {
		n.Status.SetLength()
	}
	if n.AspIdentifier != nil {
		n.AspIdentifier.SetLength()
	}
	if n.RoutingContext != nil {
		n.RoutingContext.SetLength()
	}
	if n.InfoString != nil {
		n.InfoString.SetLength()
	}

	n.Header.SetLength()
	n.Header.Length += uint32(n.Len())
}

// Len returns the actual length of Notify.
func (n *Notify) Len() int {
	l := 8

	if n.Status != nil {
		l += n.Status.Len()
	}
	if n.AspIdentifier != nil {
		l += n.AspIdentifier.Len()
	}
	if n.RoutingContext != nil {
		l += n.RoutingContext.Len()
	}
	if n.InfoString != nil {
		l += n.InfoString.Len()
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

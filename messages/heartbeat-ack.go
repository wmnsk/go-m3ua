// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// HeartbeatAck is a HeartbeatAck type of M3UA message.
type HeartbeatAck struct {
	*Header
	AspIdentifier *params.Param
	HeartbeatData *params.Param
}

// NewHeartbeatAck creates a new HeartbeatAck.
func NewHeartbeatAck(hbData *params.Param) *HeartbeatAck {
	h := &HeartbeatAck{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPSM,
			Type:     MsgTypeHeartbeatAck,
		},
		HeartbeatData: hbData,
	}
	h.SetLength()

	return h
}

// Serialize returns the byte sequence generated from a HeartbeatAck.
func (h *HeartbeatAck) Serialize() ([]byte, error) {
	b := make([]byte, h.Len())
	if err := h.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize HeartbeatAck")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (h *HeartbeatAck) SerializeTo(b []byte) error {
	if len(b) < h.Len() {
		return ErrTooShortToSerialize
	}

	h.Header.Payload = make([]byte, h.Len()-8)

	if h.HeartbeatData != nil {
		if err := h.HeartbeatData.SerializeTo(h.Header.Payload); err != nil {
			return err
		}
	}

	return h.Header.SerializeTo(b)
}

// DecodeHeartbeatAck decodes given byte sequence as a HeartbeatAck.
func DecodeHeartbeatAck(b []byte) (*HeartbeatAck, error) {
	h := &HeartbeatAck{}
	if err := h.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return h, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (h *HeartbeatAck) DecodeFromBytes(b []byte) error {
	var err error
	h.Header, err = DecodeHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.DecodeMultiParams(h.Header.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode Params")
	}
	for _, pr := range prs {
		switch pr.Tag {
		case params.HeartbeatData:
			h.HeartbeatData = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (h *HeartbeatAck) SetLength() {
	if h.HeartbeatData != nil {
		h.HeartbeatData.SetLength()
	}

	h.Header.SetLength()
	h.Header.Length += uint32(h.Len())
}

// Len returns the actual length of HeartbeatAck.
func (h *HeartbeatAck) Len() int {
	l := 8
	if h.HeartbeatData != nil {
		l += h.HeartbeatData.Len()
	}
	return l
}

// String returns the HeartbeatAck values in human readable format.
func (h *HeartbeatAck) String() string {
	return fmt.Sprintf("{Header: %s, HeartbeatData: %s}",
		h.Header.String(),
		h.HeartbeatData.String(),
	)
}

// Version returns the version of M3UA in int.
func (h *HeartbeatAck) Version() uint8 {
	return h.Header.Version
}

// MessageType returns the message type in int.
func (h *HeartbeatAck) MessageType() uint8 {
	return MsgTypeHeartbeatAck
}

// MessageClass returns the message class in int.
func (h *HeartbeatAck) MessageClass() uint8 {
	return MsgClassASPSM
}

// MessageClassName returns the name of message class.
func (h *HeartbeatAck) MessageClassName() string {
	return "ASPSM"
}

// MessageTypeName returns the name of message type.
func (h *HeartbeatAck) MessageTypeName() string {
	return "Heartbeat Ack"
}

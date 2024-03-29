// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"
	"log"

	"github.com/wmnsk/go-m3ua/messages/params"
)

// HeartbeatAck is a HeartbeatAck type of M3UA message.
//
// Spec: 3.5.6, RFC4666.
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

// MarshalBinary returns the byte sequence generated from a HeartbeatAck.
func (h *HeartbeatAck) MarshalBinary() ([]byte, error) {
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (h *HeartbeatAck) MarshalTo(b []byte) error {
	if len(b) < h.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	h.Header.Payload = make([]byte, h.MarshalLen()-8)

	if param := h.HeartbeatData; param != nil {
		if err := param.MarshalTo(h.Header.Payload); err != nil {
			return err
		}
	}

	return h.Header.MarshalTo(b)
}

// ParseHeartbeatAck decodes given byte sequence as a HeartbeatAck.
func ParseHeartbeatAck(b []byte) (*HeartbeatAck, error) {
	h := &HeartbeatAck{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (h *HeartbeatAck) UnmarshalBinary(b []byte) error {
	var err error
	h.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}

	prs, err := params.ParseMultiParams(h.Header.Payload)
	if err != nil {
		return err
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
	if param := h.HeartbeatData; param != nil {
		param.SetLength()
	}

	h.Header.Length = uint32(h.MarshalLen())
}

// MarshalLen returns the serial length of HeartbeatAck.
func (h *HeartbeatAck) MarshalLen() int {
	l := 8
	if param := h.HeartbeatData; param != nil {
		l += param.MarshalLen()
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
	return MsgClassNameASPSM
}

// MessageTypeName returns the name of message type.
func (h *HeartbeatAck) MessageTypeName() string {
	return "Heartbeat Ack"
}

// Serialize returns the byte sequence generated from a HeartbeatAck.
//
// DEPRECATED: use MarshalBinary instead.
func (h *HeartbeatAck) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return h.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (h *HeartbeatAck) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return h.MarshalTo(b)
}

// DecodeHeartbeatAck decodes given byte sequence as a HeartbeatAck.
//
// DEPRECATED: use ParseHeartbeatAck instead.
func DecodeHeartbeatAck(b []byte) (*HeartbeatAck, error) {
	log.Println("DEPRECATED: use ParseHeartbeatAck instead")
	return ParseHeartbeatAck(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (h *HeartbeatAck) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return h.UnmarshalBinary(b)
}

// Len returns the serial length of HeartbeatAck.
//
// DEPRECATED: use MarshalLen instead.
func (h *HeartbeatAck) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return h.MarshalLen()
}

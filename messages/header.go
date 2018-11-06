// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"encoding/binary"
	"fmt"
)

// Header is a M3UA common header.
type Header struct {
	Version  uint8
	Reserved uint8
	Class    uint8
	Type     uint8
	Length   uint32
	Payload  []byte
}

// NewHeader creates a new Header.
func NewHeader(version, class, mtype uint8, payload []byte) *Header {
	h := &Header{
		Version:  version,
		Reserved: 0,
		Class:    class,
		Type:     mtype,
		Payload:  payload,
	}
	h.SetLength()

	return h
}

// Serialize returns the byte sequence generated from a Header instance.
func (h *Header) Serialize() ([]byte, error) {
	b := make([]byte, h.Len())
	if err := h.SerializeTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (h *Header) SerializeTo(b []byte) error {
	if len(b) < h.Len() {
		return ErrTooShortToSerialize
	}

	b[0] = h.Version
	b[1] = h.Reserved
	b[2] = h.Class
	b[3] = h.Type
	binary.BigEndian.PutUint32(b[4:8], h.Length)
	copy(b[8:h.Len()], h.Payload)

	return nil
}

// DecodeHeader decodes given byte sequence as a M3UA common header.
func DecodeHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.DecodeFromBytes(b); err != nil {
		return nil, err
	}

	return h, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (h *Header) DecodeFromBytes(b []byte) error {
	l := len(b)
	if l < 8 {
		return ErrTooShortToDecode
	}
	h.Version = b[0]
	h.Reserved = b[1]
	h.Class = b[2]
	h.Type = b[3]
	h.Length = binary.BigEndian.Uint32(b[4:8])
	h.Payload = b[8:l]

	return nil
}

// Len returns the actual length of Header.
func (h *Header) Len() int {
	return 8 + len(h.Payload)
}

// SetLength sets the length in Length field.
func (h *Header) SetLength() {
	h.Length = uint32(8 + len(h.Payload))
}

// String returns the M3UA common header values in human readable format.
func (h *Header) String() string {
	return fmt.Sprintf("{Version: %d, Reserved: %#x, Class: %d, Type: %d, Length: %d, Payload: %x}",
		h.Version,
		h.Reserved,
		h.Class,
		h.Type,
		h.Length,
		h.Payload,
	)
}

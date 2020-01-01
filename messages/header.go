// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"encoding/binary"
	"fmt"
	"log"
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

// MarshalBinary returns the byte sequence generated from a Header instance.
func (h *Header) MarshalBinary() ([]byte, error) {
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (h *Header) MarshalTo(b []byte) error {
	if len(b) < h.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	b[0] = h.Version
	b[1] = h.Reserved
	b[2] = h.Class
	b[3] = h.Type
	binary.BigEndian.PutUint32(b[4:8], h.Length)
	copy(b[8:h.MarshalLen()], h.Payload)

	return nil
}

// ParseHeader decodes given byte sequence as a M3UA common header.
func ParseHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return h, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (h *Header) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < 8 {
		return ErrTooShortToParse
	}
	h.Version = b[0]
	h.Reserved = b[1]
	h.Class = b[2]
	h.Type = b[3]
	h.Length = binary.BigEndian.Uint32(b[4:8])
	h.Payload = b[8:l]

	return nil
}

// MarshalLen returns the serial length of Header.
func (h *Header) MarshalLen() int {
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

// Serialize returns the byte sequence generated from a Header.
//
// DEPRECATED: use MarshalBinary instead.
func (h *Header) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return h.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (h *Header) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return h.MarshalTo(b)
}

// DecodeHeader decodes given byte sequence as a Header.
//
// DEPRECATED: use ParseHeader instead.
func DecodeHeader(b []byte) (*Header, error) {
	log.Println("DEPRECATED: use ParseHeader instead")
	return ParseHeader(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (h *Header) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return h.UnmarshalBinary(b)
}

// Len returns the serial length of Header.
//
// DEPRECATED: use MarshalLen instead.
func (h *Header) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return h.MarshalLen()
}

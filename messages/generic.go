// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// Generic is generic structure of M3UA.
// This is used when
//   - Parse() method does not understand the class & type of M3UA message.
//   - users manually create this for the specific purpose.
type Generic struct {
	*Header
	Params []*params.Param
}

// New creates a new generic M3UA message.
func New(ver, mcls, mtype uint8, params ...*params.Param) *Generic {
	g := &Generic{
		Header: &Header{
			Version:  ver,
			Reserved: 0,
			Class:    mcls,
			Type:     mtype,
		},
		Params: params,
	}
	g.SetLength()

	return g
}

// MarshalBinary returns the byte sequence generated from a Generic.
func (g *Generic) MarshalBinary() ([]byte, error) {
	b := make([]byte, g.MarshalLen())
	if err := g.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize Generic")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (g *Generic) MarshalTo(b []byte) error {
	if len(b) < g.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	g.Header.Payload = make([]byte, g.MarshalLen()-8)

	var offset = 0
	for _, pr := range g.Params {
		if err := pr.MarshalTo(g.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += pr.MarshalLen()
	}

	return g.Header.MarshalTo(b)
}

// ParseGeneric decodes given byte sequence as a M3UA Generic message.
func ParseGeneric(b []byte) (*Generic, error) {
	g := &Generic{}
	if err := g.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return g, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (g *Generic) UnmarshalBinary(b []byte) error {
	var err error
	g.Header, err = ParseHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.ParseMultiParams(g.Header.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode Params")
	}
	g.Params = append(g.Params, prs...)
	return nil
}

// MarshalLen returns the serial length of Data.
func (g *Generic) MarshalLen() int {
	l := 8
	for _, pr := range g.Params {
		l += pr.MarshalLen()
	}
	return l
}

// SetLength sets the length in Length field.
func (g *Generic) SetLength() {
	for _, pr := range g.Params {
		pr.SetLength()
	}

	g.Header.Length = uint32(g.MarshalLen())
}

// String returns the Generic values in human readable format.
func (g *Generic) String() string {
	var paramStr []string
	for _, pr := range g.Params {
		paramStr = append(paramStr, pr.String())
	}

	return fmt.Sprintf("{Header: %s, Params: %s}",
		g.Header.String(),
		paramStr,
	)
}

// Version returns the version of M3UA in int.
func (g *Generic) Version() uint8 {
	return g.Header.Version
}

// MessageType returns the message type in int.
func (g *Generic) MessageType() uint8 {
	return g.Header.Type
}

// MessageClass returns the message class in int.
func (g *Generic) MessageClass() uint8 {
	return g.Header.Class
}

// MessageClassName returns the name of message class.
func (g *Generic) MessageClassName() string {
	return "Unknown"
}

// MessageTypeName returns the name of message type.
func (g *Generic) MessageTypeName() string {
	return "Unknown"
}

// Serialize returns the byte sequence generated from a Generic.
//
// DEPRECATED: use MarshalBinary instead.
func (g *Generic) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return g.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (g *Generic) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return g.MarshalTo(b)
}

// DecodeGeneric decodes given byte sequence as a Generic.
//
// DEPRECATED: use ParseGeneric instead.
func DecodeGeneric(b []byte) (*Generic, error) {
	log.Println("DEPRECATED: use ParseGeneric instead")
	return ParseGeneric(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (g *Generic) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return g.UnmarshalBinary(b)
}

// Len returns the serial length of Generic.
//
// DEPRECATED: use MarshalLen instead.
func (g *Generic) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return g.MarshalLen()
}

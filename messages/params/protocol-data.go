// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

import (
	"encoding/binary"
	"fmt"
	"log"
)

// ServiceIndicator definitions.
const (
	ServiceIndUnused uint8 = iota
	_
	_
	ServiceIndSCCP
	ServiceIndTUP
	ServiceIndISUP
	_
	ServiceIndBroadbandISUP
	ServiceIndSatelliteISUP
	_
	ServiceIndAALType2Signalling
	ServiceIndBICC
	ServiceIndGatewayControlProtocol
	_
)

// ProtocolDataPayload is a M3UA ProtocolDataPayload.
type ProtocolDataPayload struct {
	OriginatingPointCode   uint32
	DestinationPointCode   uint32
	ServiceIndicator       uint8
	NetworkIndicator       uint8
	MessagePriority        uint8
	SignalingLinkSelection uint8
	Data                   []byte
}

// NewProtocolDataPayload creates a new ProtocolDataPayload payload.
// Note that this does not contain the Tag and Length inside.
// You need to create new Param using serialized ProtocolDataPayload.
func NewProtocolDataPayload(opc, dpc uint32, si, ni, mp, sls uint8, data []byte) *ProtocolDataPayload {
	return &ProtocolDataPayload{
		OriginatingPointCode:   opc,
		DestinationPointCode:   dpc,
		ServiceIndicator:       si,
		NetworkIndicator:       ni,
		MessagePriority:        mp,
		SignalingLinkSelection: sls,
		Data:                   data,
	}
}

// NewProtocolData creates a new ProtocolData.
// Note that this returns *Param, as no specific structure in this parameter.
// Also, Payload will be serialized and not accessible until calling ProtocolData() func.
func NewProtocolData(opc, dpc uint32, si, ni, mp, sls uint8, data []byte) *Param {
	pd, _ := NewProtocolDataPayload(opc, dpc, si, ni, mp, sls, data).MarshalBinary()
	p := &Param{
		Tag:  ProtocolData,
		Data: pd,
	}
	p.SetLength()

	return p
}

// ProtocolData returns ProtocolDataPayload
func (p *Param) ProtocolData() (*ProtocolDataPayload, error) {
	if p.Tag != ProtocolData {
		return nil, ErrInvalidType
	}

	return ParseProtocolDataPayload(p.Data)
}

// MarshalBinary returns the byte sequence generated from a M3UA ProtocolDataPayload instance.
func (p *ProtocolDataPayload) MarshalBinary() ([]byte, error) {
	b := make([]byte, p.MarshalLen())
	if err := p.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (p *ProtocolDataPayload) MarshalTo(b []byte) error {
	if len(b) < p.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	binary.BigEndian.PutUint32(b[0:4], p.OriginatingPointCode)
	binary.BigEndian.PutUint32(b[4:8], p.DestinationPointCode)
	b[8] = p.ServiceIndicator
	b[9] = p.NetworkIndicator
	b[10] = p.MessagePriority
	b[11] = p.SignalingLinkSelection
	copy(b[12:p.MarshalLen()], p.Data)
	return nil
}

// ParseProtocolDataPayload decodes given byte sequence as a M3UA ProtocolDataPayload.
func ParseProtocolDataPayload(b []byte) (*ProtocolDataPayload, error) {
	p := &ProtocolDataPayload{}
	if err := p.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return p, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA ProtocolDataPayload.
func (p *ProtocolDataPayload) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < 12 {
		return ErrTooShortToParse
	}

	p.OriginatingPointCode = binary.BigEndian.Uint32(b[0:4])
	p.DestinationPointCode = binary.BigEndian.Uint32(b[4:8])
	p.ServiceIndicator = b[8]
	p.NetworkIndicator = b[9]
	p.MessagePriority = b[10]
	p.SignalingLinkSelection = b[11]
	p.Data = b[12:]
	return nil
}

// MarshalLen returns serial length in integer.
func (p *ProtocolDataPayload) MarshalLen() int {
	return 12 + len(p.Data)
}

// String returns the M3UA header values in human readable format.
func (p *ProtocolDataPayload) String() string {
	return fmt.Sprintf("{OriginatingPointCode: %d, DestinationPointCode: %d, ServiceIndicator: %d, NetworkIndicator: %d, MessagePriority: %d, SignalingLinkSelection: %d, Data: %x}",
		p.OriginatingPointCode,
		p.DestinationPointCode,
		p.ServiceIndicator,
		p.NetworkIndicator,
		p.MessagePriority,
		p.SignalingLinkSelection,
		p.Data,
	)
}

// Serialize returns the byte sequence generated from a ProtocolDataPayload.
//
// DEPRECATED: use MarshalBinary instead.
func (p *ProtocolDataPayload) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return p.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (p *ProtocolDataPayload) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return p.MarshalTo(b)
}

// DecodeProtocolDataPayload decodes given byte sequence as a ProtocolDataPayload.
//
// DEPRECATED: use ParseProtocolDataPayload instead.
func DecodeProtocolDataPayload(b []byte) (*ProtocolDataPayload, error) {
	log.Println("DEPRECATED: use ParseProtocolDataPayload instead")
	return ParseProtocolDataPayload(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (p *ProtocolDataPayload) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return p.UnmarshalBinary(b)
}

// Len returns the serial length of ProtocolDataPayload.
//
// DEPRECATED: use MarshalLen instead.
func (p *ProtocolDataPayload) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return p.MarshalLen()
}

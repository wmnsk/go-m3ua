// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"
	"log"

	"github.com/wmnsk/go-m3ua/messages/params"
)

// DestinationStateAudit is a DestinationStateAudit type of M3UA message.
//
// Spec: 3.4.3, RFC4666.
type DestinationStateAudit struct {
	*Header
	NetworkAppearance *params.Param
	RoutingContext    *params.Param
	AffectedPointCode *params.Param
	InfoString        *params.Param
}

// NewDestinationStateAudit creates a new DestinationStateAudit.
func NewDestinationStateAudit(nwApr, rtCtx, apcs, info *params.Param) *DestinationStateAudit {
	d := &DestinationStateAudit{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassSSNM,
			Type:     MsgTypeDestinationStateAudit,
		},
		NetworkAppearance: nwApr,
		RoutingContext:    rtCtx,
		AffectedPointCode: apcs,
		InfoString:        info,
	}
	d.SetLength()

	return d
}

// MarshalBinary returns the byte sequence generated from a DestinationStateAudit.
func (d *DestinationStateAudit) MarshalBinary() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (d *DestinationStateAudit) MarshalTo(b []byte) error {
	if len(b) < d.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	d.Header.Payload = make([]byte, d.MarshalLen()-8)

	var offset = 0
	if param := d.NetworkAppearance; param != nil {
		if err := param.MarshalTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += param.MarshalLen()
	}
	if param := d.RoutingContext; param != nil {
		if err := param.MarshalTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += param.MarshalLen()
	}
	if param := d.AffectedPointCode; param != nil {
		if err := param.MarshalTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += param.MarshalLen()
	}
	if param := d.InfoString; param != nil {
		if err := param.MarshalTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
	}
	return d.Header.MarshalTo(b)
}

// ParseDestinationStateAudit decodes given byte sequence as a DestinationStateAudit.
func ParseDestinationStateAudit(b []byte) (*DestinationStateAudit, error) {
	d := &DestinationStateAudit{}
	if err := d.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (d *DestinationStateAudit) UnmarshalBinary(b []byte) error {
	var err error
	d.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}

	prs, err := params.ParseMultiParams(d.Header.Payload)
	if err != nil {
		return err
	}
	for _, pr := range prs {
		switch pr.Tag {
		case params.NetworkAppearance:
			d.NetworkAppearance = pr
		case params.RoutingContext:
			d.RoutingContext = pr
		case params.AffectedPointCode:
			d.AffectedPointCode = pr
		case params.InfoString:
			d.InfoString = pr
		default:
			return fmt.Errorf("failed to decode DAUD: %w", ErrInvalidParameter)
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (d *DestinationStateAudit) SetLength() {
	if param := d.NetworkAppearance; param != nil {
		param.SetLength()
	}
	if param := d.RoutingContext; param != nil {
		param.SetLength()
	}
	if param := d.AffectedPointCode; param != nil {
		param.SetLength()
	}
	if param := d.InfoString; param != nil {
		param.SetLength()
	}

	d.Header.Length = uint32(d.MarshalLen())
}

// MarshalLen returns the serial length of DestinationStateAudit.
func (d *DestinationStateAudit) MarshalLen() int {
	l := 8
	if param := d.NetworkAppearance; param != nil {
		l += param.MarshalLen()
	}
	if param := d.RoutingContext; param != nil {
		l += param.MarshalLen()
	}
	if param := d.AffectedPointCode; param != nil {
		l += param.MarshalLen()
	}
	if param := d.InfoString; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// Version returns the version of M3UA in int.
func (d *DestinationStateAudit) Version() uint8 {
	return d.Header.Version
}

// MessageType returns the message type in int.
func (d *DestinationStateAudit) MessageType() uint8 {
	return MsgTypeDestinationStateAudit
}

// MessageClass returns the message class in int.
func (d *DestinationStateAudit) MessageClass() uint8 {
	return MsgClassSSNM
}

// MessageClassName returns the name of message class.
func (d *DestinationStateAudit) MessageClassName() string {
	return MsgClassNameSSNM
}

// MessageTypeName returns the name of message type.
func (d *DestinationStateAudit) MessageTypeName() string {
	return "Destination State Audit"
}

// Serialize returns the byte sequence generated from a DestinationStateAudit.
//
// DEPRECATED: use MarshalBinary instead.
func (d *DestinationStateAudit) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return d.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (d *DestinationStateAudit) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return d.MarshalTo(b)
}

// DecodeDestinationStateAudit decodes given byte sequence as a DestinationStateAudit.
//
// DEPRECATED: use ParseDestinationStateAudit instead.
func DecodeDestinationStateAudit(b []byte) (*DestinationStateAudit, error) {
	log.Println("DEPRECATED: use ParseDestinationStateAudit instead")
	return ParseDestinationStateAudit(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (d *DestinationStateAudit) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return d.UnmarshalBinary(b)
}

// Len returns the serial length of DestinationStateAudit.
//
// DEPRECATED: use MarshalLen instead.
func (d *DestinationStateAudit) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return d.MarshalLen()
}

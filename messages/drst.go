// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"log"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// DestinationRestricted is a DestinationRestricted type of M3UA message.
//
// Spec: 3.4.6, RFC4666.
type DestinationRestricted struct {
	*Header
	NetworkAppearance *params.Param
	RoutingContext    *params.Param
	AffectedPointCode *params.Param
	InfoString        *params.Param
}

// NewDestinationRestricted creates a new DestinationRestricted.
func NewDestinationRestricted(nwApr, rtCtx, apcs, info *params.Param) *DestinationRestricted {
	d := &DestinationRestricted{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassSSNM,
			Type:     MsgTypeDestinationRestricted,
		},
		NetworkAppearance: nwApr,
		RoutingContext:    rtCtx,
		AffectedPointCode: apcs,
		InfoString:        info,
	}
	d.SetLength()

	return d
}

// MarshalBinary returns the byte sequence generated from a DestinationRestricted.
func (d *DestinationRestricted) MarshalBinary() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize DestinationRestricted")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (d *DestinationRestricted) MarshalTo(b []byte) error {
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

// ParseDestinationRestricted decodes given byte sequence as a DestinationRestricted.
func ParseDestinationRestricted(b []byte) (*DestinationRestricted, error) {
	d := &DestinationRestricted{}
	if err := d.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (d *DestinationRestricted) UnmarshalBinary(b []byte) error {
	var err error
	d.Header, err = ParseHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode DUNA")
	}

	prs, err := params.ParseMultiParams(d.Header.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode DUNA")
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
			return errors.Wrap(ErrInvalidParameter, "failed to decode DUNA")
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (d *DestinationRestricted) SetLength() {
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

	d.Header.SetLength()
	d.Header.Length += uint32(d.MarshalLen())
}

// MarshalLen returns the serial length of DestinationRestricted.
func (d *DestinationRestricted) MarshalLen() int {
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
func (d *DestinationRestricted) Version() uint8 {
	return d.Header.Version
}

// MessageType returns the message type in int.
func (d *DestinationRestricted) MessageType() uint8 {
	return MsgTypeDestinationRestricted
}

// MessageClass returns the message class in int.
func (d *DestinationRestricted) MessageClass() uint8 {
	return MsgClassSSNM
}

// MessageClassName returns the name of message class.
func (d *DestinationRestricted) MessageClassName() string {
	return "SSNM"
}

// MessageTypeName returns the name of message type.
func (d *DestinationRestricted) MessageTypeName() string {
	return "Destination Unavailable"
}

// Serialize returns the byte sequence generated from a DestinationRestricted.
//
// DEPRECATED: use MarshalBinary instead.
func (d *DestinationRestricted) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return d.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (d *DestinationRestricted) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return d.MarshalTo(b)
}

// DecodeDestinationRestricted decodes given byte sequence as a DestinationRestricted.
//
// DEPRECATED: use ParseDestinationRestricted instead.
func DecodeDestinationRestricted(b []byte) (*DestinationRestricted, error) {
	log.Println("DEPRECATED: use ParseDestinationRestricted instead")
	return ParseDestinationRestricted(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (d *DestinationRestricted) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return d.UnmarshalBinary(b)
}

// Len returns the serial length of DestinationRestricted.
//
// DEPRECATED: use MarshalLen instead.
func (d *DestinationRestricted) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return d.MarshalLen()
}

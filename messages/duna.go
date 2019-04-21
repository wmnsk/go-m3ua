// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// DestinationUnavailable is a DestinationUnavailable type of M3UA message.
//
// Spec: 3.4.1, RFC4666.
type DestinationUnavailable struct {
	*Header
	NetworkAppearance *params.Param
	RoutingContext    *params.Param
	AffectedPointCode *params.Param
	InfoString        *params.Param
}

// NewDestinationUnavailable creates a new DestinationUnavailable.
func NewDestinationUnavailable(nwApr, rtCtx, apcs, info *params.Param) *DestinationUnavailable {
	d := &DestinationUnavailable{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassSSNM,
			Type:     MsgTypeDestinationUnavailable,
		},
		NetworkAppearance: nwApr,
		RoutingContext:    rtCtx,
		AffectedPointCode: apcs,
		InfoString:        info,
	}
	d.SetLength()

	return d
}

// Serialize returns the byte sequence generated from a DestinationUnavailable.
func (d *DestinationUnavailable) Serialize() ([]byte, error) {
	b := make([]byte, d.Len())
	if err := d.SerializeTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize DestinationUnavailable")
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (d *DestinationUnavailable) SerializeTo(b []byte) error {
	if len(b) < d.Len() {
		return ErrTooShortToSerialize
	}

	d.Header.Payload = make([]byte, d.Len()-8)

	var offset = 0
	if p := d.NetworkAppearance; p != nil {
		if err := p.SerializeTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += p.Len()
	}
	if p := d.RoutingContext; p != nil {
		if err := p.SerializeTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += p.Len()
	}
	if p := d.AffectedPointCode; p != nil {
		if err := p.SerializeTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += p.Len()
	}
	if p := d.InfoString; p != nil {
		if err := p.SerializeTo(d.Header.Payload[offset:]); err != nil {
			return err
		}
	}
	return d.Header.SerializeTo(b)
}

// DecodeDestinationUnavailable decodes given byte sequence as a DestinationUnavailable.
func DecodeDestinationUnavailable(b []byte) (*DestinationUnavailable, error) {
	d := &DestinationUnavailable{}
	if err := d.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return d, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
func (d *DestinationUnavailable) DecodeFromBytes(b []byte) error {
	var err error
	d.Header, err = DecodeHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode DUNA")
	}

	prs, err := params.DecodeMultiParams(d.Header.Payload)
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
func (d *DestinationUnavailable) SetLength() {
	if p := d.NetworkAppearance; p != nil {
		p.SetLength()
	}
	if p := d.RoutingContext; p != nil {
		p.SetLength()
	}
	if p := d.AffectedPointCode; p != nil {
		p.SetLength()
	}
	if p := d.InfoString; p != nil {
		p.SetLength()
	}

	d.Header.SetLength()
	d.Header.Length += uint32(d.Len())
}

// Len returns the actual length of DestinationUnavailable.
func (d *DestinationUnavailable) Len() int {
	l := 8
	if p := d.NetworkAppearance; p != nil {
		l += p.Len()
	}
	if p := d.RoutingContext; p != nil {
		l += p.Len()
	}
	if p := d.AffectedPointCode; p != nil {
		l += p.Len()
	}
	if p := d.InfoString; p != nil {
		l += p.Len()
	}
	return l
}

// Version returns the version of M3UA in int.
func (d *DestinationUnavailable) Version() uint8 {
	return d.Header.Version
}

// MessageType returns the message type in int.
func (d *DestinationUnavailable) MessageType() uint8 {
	return MsgTypeDestinationUnavailable
}

// MessageClass returns the message class in int.
func (d *DestinationUnavailable) MessageClass() uint8 {
	return MsgClassSSNM
}

// MessageClassName returns the name of message class.
func (d *DestinationUnavailable) MessageClassName() string {
	return "SSNM"
}

// MessageTypeName returns the name of message type.
func (d *DestinationUnavailable) MessageTypeName() string {
	return "Destination Unavailable"
}

// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// DestinationUserPartUnavailable is a DestinationUserPartUnavailable type of M3UA message.
//
// Spec: 3.4.5, RFC4666.
type DestinationUserPartUnavailable struct {
	*Header
	NetworkAppearance *params.Param
	RoutingContext    *params.Param
	AffectedPointCode *params.Param
	UserCause         *params.Param
	InfoString        *params.Param
}

// NewDestinationUserPartUnavailable creates a new DestinationUserPartUnavailable.
func NewDestinationUserPartUnavailable(nwApr, rtCtx, apcs, cause, info *params.Param) *DestinationUserPartUnavailable {
	d := &DestinationUserPartUnavailable{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassSSNM,
			Type:     MsgTypeDestinationUserPartUnavailable,
		},
		NetworkAppearance: nwApr,
		RoutingContext:    rtCtx,
		AffectedPointCode: apcs,
		UserCause:         cause,
		InfoString:        info,
	}
	d.SetLength()

	return d
}

// MarshalBinary returns the byte sequence generated from a DestinationUserPartUnavailable.
func (d *DestinationUserPartUnavailable) MarshalBinary() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize DestinationUserPartUnavailable")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (d *DestinationUserPartUnavailable) MarshalTo(b []byte) error {
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
	if param := d.UserCause; param != nil {
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

// ParseDestinationUserPartUnavailable decodes given byte sequence as a DestinationUserPartUnavailable.
func ParseDestinationUserPartUnavailable(b []byte) (*DestinationUserPartUnavailable, error) {
	d := &DestinationUserPartUnavailable{}
	if err := d.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (d *DestinationUserPartUnavailable) UnmarshalBinary(b []byte) error {
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
		case params.UserCause:
			d.UserCause = pr
		case params.InfoString:
			d.InfoString = pr
		default:
			return errors.Wrap(ErrInvalidParameter, "failed to decode DUNA")
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (d *DestinationUserPartUnavailable) SetLength() {
	if param := d.NetworkAppearance; param != nil {
		param.SetLength()
	}
	if param := d.RoutingContext; param != nil {
		param.SetLength()
	}
	if param := d.AffectedPointCode; param != nil {
		param.SetLength()
	}
	if param := d.UserCause; param != nil {
		param.SetLength()
	}
	if param := d.InfoString; param != nil {
		param.SetLength()
	}

	d.Header.Length = uint32(d.MarshalLen())
}

// MarshalLen returns the serial length of DestinationUserPartUnavailable.
func (d *DestinationUserPartUnavailable) MarshalLen() int {
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
	if param := d.UserCause; param != nil {
		l += param.MarshalLen()
	}
	if param := d.InfoString; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// Version returns the version of M3UA in int.
func (d *DestinationUserPartUnavailable) Version() uint8 {
	return d.Header.Version
}

// MessageType returns the message type in int.
func (d *DestinationUserPartUnavailable) MessageType() uint8 {
	return MsgTypeDestinationUserPartUnavailable
}

// MessageClass returns the message class in int.
func (d *DestinationUserPartUnavailable) MessageClass() uint8 {
	return MsgClassSSNM
}

// MessageClassName returns the name of message class.
func (d *DestinationUserPartUnavailable) MessageClassName() string {
	return "SSNM"
}

// MessageTypeName returns the name of message type.
func (d *DestinationUserPartUnavailable) MessageTypeName() string {
	return "Destination User Part Unavailable"
}

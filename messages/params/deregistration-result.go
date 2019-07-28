// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

import "log"

// DeregResultPayload is the payload of DeregistrationResult.
type DeregResultPayload struct {
	RoutingContext, DeregistrationStatus *Param
}

// NewDeregResultPayload creates a new DeregResultPayload.
func NewDeregResultPayload(rtCtx, deregStatus *Param) *DeregResultPayload {
	return &DeregResultPayload{
		RoutingContext:       rtCtx,
		DeregistrationStatus: deregStatus,
	}
}

// NewDeregistrationResult creates a new DeregistrationResult.
// Note that this returns *Param, as no specific structure in this parameter.
func NewDeregistrationResult(dr *DeregResultPayload) *Param {
	return newNestedParam(
		DeregistrationResult,
		dr.RoutingContext,
		dr.DeregistrationStatus,
	)
}

// DeregistrationResult returns DeregResultPayload.
func (p *Param) DeregistrationResult() (*DeregResultPayload, error) {
	if p.Tag != DeregistrationResult {
		return nil, ErrInvalidType
	}

	d, err := ParseDeregResultPayload(p.Data)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// ParseDeregResultPayload decodes given byte sequence as a DeregResultPayload.
func ParseDeregResultPayload(b []byte) (*DeregResultPayload, error) {
	d := &DeregResultPayload{}
	if err := d.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a Param.
func (d *DeregResultPayload) UnmarshalBinary(b []byte) error {
	ps, err := ParseMultiParams(b)
	if err != nil {
		return err
	}
	if len(ps) != 2 {
		return ErrInvalidLength
	}

	for _, p := range ps {
		switch p.Tag {
		case RoutingContext:
			d.RoutingContext = p
		case DeregistrationStatus:
			d.DeregistrationStatus = p
		}
	}
	return nil
}

// DecodeDeregResultPayload decodes given byte sequence as a DeregResultPayload.
//
// DEPRECATED: use ParseDeregResultPayload instead.
func DecodeDeregResultPayload(b []byte) (*DeregResultPayload, error) {
	log.Println("DEPRECATED: use ParseDeregResultPayload instead")
	return ParseDeregResultPayload(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a Param.
//
// DEPRECATED: use UnmarshalBinary instead.
func (d *DeregResultPayload) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return d.UnmarshalBinary(b)
}

// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

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

	d, err := DecodeDeregResultPayload(p.Data)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// DecodeDeregResultPayload decodes given byte sequence as a DeregResultPayload.
func DecodeDeregResultPayload(b []byte) (*DeregResultPayload, error) {
	d := &DeregResultPayload{}
	if err := d.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return d, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a Param.
func (d *DeregResultPayload) DecodeFromBytes(b []byte) error {
	ps, err := DecodeMultiParams(b)
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

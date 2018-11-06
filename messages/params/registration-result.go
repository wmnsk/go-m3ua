// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// RegistrationResultPayload is the payload of RegistrationResult.
type RegistrationResultPayload struct {
	LocalRoutingKeyIdentifier, RegistrationStatus, RoutingContext *Param
}

// NewRegistrationResultPayload creates a new RegistrationResultPayload.
func NewRegistrationResultPayload(rkID, deregStatus, rtCtx *Param) *RegistrationResultPayload {
	return &RegistrationResultPayload{
		LocalRoutingKeyIdentifier: rkID,
		RegistrationStatus:        deregStatus,
		RoutingContext:            rtCtx,
	}
}

// NewRegistrationResult creates a new RegistrationResult.
// Note that this returns *Param, as no specific structure in this parameter.
func NewRegistrationResult(rr *RegistrationResultPayload) *Param {
	return newNestedParam(
		RegistrationResult,
		rr.LocalRoutingKeyIdentifier,
		rr.RegistrationStatus,
		rr.RoutingContext,
	)
}

// RegistrationResult returns RegistrationResultPayload.
func (p *Param) RegistrationResult() (*RegistrationResultPayload, error) {
	if p.Tag != RegistrationResult {
		return nil, ErrInvalidType
	}

	d, err := DecodeRegistrationResultPayload(p.Data)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// DecodeRegistrationResultPayload decodes given byte sequence as a RegistrationResultPayload.
func DecodeRegistrationResultPayload(b []byte) (*RegistrationResultPayload, error) {
	d := &RegistrationResultPayload{}
	if err := d.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return d, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a Param.
func (d *RegistrationResultPayload) DecodeFromBytes(b []byte) error {
	ps, err := DecodeMultiParams(b)
	if err != nil {
		return err
	}
	if len(ps) != 3 {
		return ErrInvalidLength
	}

	d.LocalRoutingKeyIdentifier = ps[0]
	d.RegistrationStatus = ps[1]
	d.RoutingContext = ps[2]

	return nil
}

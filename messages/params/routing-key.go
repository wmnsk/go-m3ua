// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// RoutingKeyPayload is the payload of RoutingKey.
type RoutingKeyPayload struct {
	LocalRoutingKeyIdentifier, RoutingContext, TrafficModeType, DestinationPointCode, NetworkAppearance, ServiceIndicators, OriginatingPointCodeList *Param
}

// NewRoutingKeyPayload creates a new RoutingKeyPayload.
func NewRoutingKeyPayload(rkID, rtCtx, tmType, dpc, nwApr, si, opcs *Param) *RoutingKeyPayload {
	return &RoutingKeyPayload{
		LocalRoutingKeyIdentifier: rkID,
		RoutingContext:            rtCtx,
		TrafficModeType:           tmType,
		DestinationPointCode:      dpc,
		NetworkAppearance:         nwApr,
		ServiceIndicators:         si,
		OriginatingPointCodeList:  opcs,
	}
}

// Note that this parameter contains some optional parameters inside.

// NewRoutingKey creates a new RoutingKey.
// Note that this returns *Param, as no specific structure in this parameter.
func NewRoutingKey(rk *RoutingKeyPayload) *Param {
	return newNestedParam(
		RoutingKey,
		rk.LocalRoutingKeyIdentifier,
		rk.RoutingContext,
		rk.TrafficModeType,
		rk.DestinationPointCode,
		rk.NetworkAppearance,
		rk.ServiceIndicators,
		rk.OriginatingPointCodeList,
	)
}

// RoutingKey returns RoutingKeyPayload.
func (p *Param) RoutingKey() (*RoutingKeyPayload, error) {
	if p.Tag != RoutingKey {
		return nil, ErrInvalidType
	}

	r, err := DecodeRoutingKeyPayload(p.Data)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DecodeRoutingKeyPayload decodes given byte sequence as a RoutingKeyPayload.
func DecodeRoutingKeyPayload(b []byte) (*RoutingKeyPayload, error) {
	r := &RoutingKeyPayload{}
	if err := r.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return r, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a Param.
func (r *RoutingKeyPayload) DecodeFromBytes(b []byte) error {
	ps, err := DecodeMultiParams(b)
	if err != nil {
		return err
	}
	if len(ps) < 3 {
		return ErrInvalidLength
	}

	for _, p := range ps {
		switch p.Tag {
		case LocalRoutingKeyIdentifier:
			r.LocalRoutingKeyIdentifier = p
		case RoutingContext:
			r.RoutingContext = p
		case TrafficModeType:
			r.TrafficModeType = p
		case DestinationPointCode:
			r.DestinationPointCode = p
		case NetworkAppearance:
			r.NetworkAppearance = p
		case ServiceIndicators:
			r.ServiceIndicators = p
		case OriginatingPointCodeList:
			r.OriginatingPointCodeList = p
		default:
			return ErrInvalidType
		}
	}
	return nil
}

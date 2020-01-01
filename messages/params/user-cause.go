// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// User Identity definitions.
const (
	UserIdentityUnknown uint16 = iota
	Unequipped
	Inaccessible
)

// Unavailability Cause definitions.
const (
	_ uint16 = iota
	SCCP
	TUP
	ISUP
	_
	BroadbandISUP
	SatelliteISUP
	_
	AAL2Signalling
	BICC
	GatewayControlProtocol
	_
)

// NewUserCause creates the User/Cause Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewUserCause(user, cause uint16) *Param {
	comb := uint32(cause)<<16 | uint32(user)
	return newUint32ValParam(UserCause, comb)
}

// UserCause returns multiple UserCause from Param.
func (p *Param) UserCause() uint32 {
	if p.Tag != UserCause {
		return 0
	}

	return p.decodeUint32ValData()
}

// UserIdentity returns multiple UserIdentity from Param.
func (p *Param) UserIdentity() uint16 {
	if p.Tag != UserCause {
		return 0
	}

	return uint16(p.decodeUint32ValData() & 0xffff)
}

// UnavailabilityCause returns multiple UnavailabilityCause from Param.
func (p *Param) UnavailabilityCause() uint16 {
	if p.Tag != UserCause {
		return 0
	}

	return uint16(p.decodeUint32ValData() >> 16)
}

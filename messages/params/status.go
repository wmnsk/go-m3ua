// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// Status Type definitions.
const (
	AsStateChange uint16 = iota + 1
	Other
)

// Status Information definitions (StatusType==AsStateChange).
// Note that this contains the type.
const (
	_ uint32 = uint32((0x01 << 16) | iota + 1)
	AsStateInactive
	AsStateActive
	AsStatePending
)

// Status Information definitions (StatusType==Other).
// Note that this contains the type.
const (
	InsufficientAspResources uint32 = uint32((0x02 << 16) | iota + 1)
	AlternateAspActive
	AspFailure
)

// NewStatus creates the Status Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
// The argument typeInfo can be chosen from Status Information definitions above.
func NewStatus(typeInfo uint32) *Param {
	return newUint32ValParam(Status, typeInfo)
}

// Status returns multiple Status from Param.
func (p *Param) Status() uint32 {
	if p.Tag != Status {
		return 0
	}

	return p.decodeUint32ValData()
}

// StatusType returns multiple StatusType from Param.
func (p *Param) StatusType() uint16 {
	if p.Tag != Status {
		return 0
	}

	return uint16(p.decodeUint32ValData() >> 16)
}

// StatusInfo returns multiple StatusInfo from Param.
func (p *Param) StatusInfo() uint16 {
	if p.Tag != Status {
		return 0
	}

	return uint16(p.decodeUint32ValData() & 0xffff)
}

// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// Deregistration Status definitions.
const (
	SuccessfullyDeregistered uint32 = iota
	DeregStatusUnknown
	DeregInvalidRoutingContext
	DeregPermissionDenied
	DeregNotRegistered
	DeregASPActiveForRoutingContext
)

// NewDeregistrationStatus creates the DeregistrationStatus Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewDeregistrationStatus(deregStatus uint32) *Param {
	return newUint32ValParam(DeregistrationStatus, deregStatus)
}

// DeregistrationStatus returns multiple DeregistrationStatus from Param.
func (p *Param) DeregistrationStatus() uint32 {
	if p.Tag != DeregistrationStatus {
		return 0
	}
	return p.decodeUint32ValData()
}

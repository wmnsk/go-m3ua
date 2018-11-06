// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// Registration Status definitions.
const (
	SuccessfullyRegistered uint32 = iota
	RegistrationStatusUnknown
	InvalidDPC
	InvalidNetworkAppearance
	InvalidRoutingKey
	PermissionDenied
	CannotSupportUniqueRouting
	RoutingKeynotCurrentlyProvisioned
	InsufficientResources
	UnsupportedRKparameterField
	UnsupportedTrafficHandlingMode
	RoutingKeyChangeRefused
	RoutingKeyAlreadyRegistered
)

// NewRegistrationStatus creates the RegistrationStatus Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewRegistrationStatus(regStatus uint32) *Param {
	return newUint32ValParam(RegistrationStatus, regStatus)
}

// RegistrationStatus returns multiple RegistrationStatus from Param.
func (p *Param) RegistrationStatus() uint32 {
	if p.Tag != RegistrationStatus {
		return 0
	}

	return p.decodeUint32ValData()
}

// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// Error Code definitions.
const (
	_ uint32 = iota
	ErrInvalidVersion
	_
	ErrUnsupportedMessageClass
	ErrUnsupportedMessageType
	ErrUnsupportedTrafficModeType
	ErrUnexpectedMessage
	ErrProtocolError
	_
	ErrInvalidStreamIdentifier
	_
	_
	_
	ErrRefusedManagementBlocking
	ErrAspIdentifierRequired
	ErrInvalidAspIdentifier
	_
	ErrInvalidParameterValue
	ErrParameterFieldError
	ErrUnexpectedParameter
	ErrDestinationStatusUnknown
	ErrInvalidNetworkAppearance
	ErrMissingParameter
	_
	_
	ErrInvalidRoutingContext
	ErrNoConfiguredAsForAsp
)

// NewErrorCode creates the ErrorCode Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewErrorCode(ec uint32) *Param {
	return newUint32ValParam(ErrorCode, ec)
}

// ErrorCode returns multiple ErrorCode from Param.
func (p *Param) ErrorCode() uint32 {
	if p.Tag != ErrorCode {
		return 0
	}
	return p.decodeUint32ValData()
}

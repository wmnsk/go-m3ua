// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewAspIdentifier creates the AspIdentifier Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewAspIdentifier(aspID uint32) *Param {
	return newUint32ValParam(AspIdentifier, aspID)
}

// AspIdentifier returns multiple AspIdentifier from Param.
func (p *Param) AspIdentifier() uint32 {
	if p.Tag != AspIdentifier {
		return 0
	}
	return p.decodeUint32ValData()
}

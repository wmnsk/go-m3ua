// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewLocalRoutingKeyIdentifier creates the LocalRoutingKeyIdentifier Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewLocalRoutingKeyIdentifier(rkID uint32) *Param {
	return newUint32ValParam(LocalRoutingKeyIdentifier, rkID)
}

// LocalRoutingKeyIdentifier returns multiple LocalRoutingKeyIdentifier from Param.
func (p *Param) LocalRoutingKeyIdentifier() uint32 {
	if p.Tag != LocalRoutingKeyIdentifier {
		return 0
	}
	return p.decodeUint32ValData()
}

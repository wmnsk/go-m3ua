// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewDestinationPointCode creates the DestinationPointCode Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewDestinationPointCode(dpc uint32) *Param {
	return newUint24ValParam(DestinationPointCode, dpc)
}

// DestinationPointCode returns multiple DestinationPointCode from Param.
func (p *Param) DestinationPointCode() uint32 {
	if p.Tag != DestinationPointCode {
		return 0
	}
	return p.decodeUint32ValData() & 0xffffff
}

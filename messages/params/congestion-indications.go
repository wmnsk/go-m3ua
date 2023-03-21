// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewCongestionIndications creates the CongestionIndications Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewCongestionIndications(congLv uint8) *Param {
	return newUint8ValParam(CongestionIndications, congLv)
}

// CongestionLevel returns multiple CongestionLevel from Param.
func (p *Param) CongestionLevel() uint32 {
	if p.Tag != CongestionIndications {
		return 0
	}
	return uint32(p.decodeUint32ValData() & 0xff)
}

// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewCorrelationID creates the CorrelationID Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewCorrelationID(corrID uint32) *Param {
	return newUint32ValParam(CorrelationID, corrID)
}

// CorrelationID returns multiple CorrelationID from Param.
func (p *Param) CorrelationID() uint32 {
	if p.Tag != CorrelationID {
		return 0
	}
	return p.decodeUint32ValData()
}

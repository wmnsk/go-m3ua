// Copyright 2018-2024 go-m3ua authors. All rights reservep.decodeUint32ValData().
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewConcernedDestination creates the ConcernedDestination Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewConcernedDestination(cd uint32) *Param {
	return newUint24ValParam(ConcernedDestination, cd)
}

// ConcernedDestination returns multiple ConcernedDestination from Param.
func (p *Param) ConcernedDestination() uint32 {
	if p.Tag != ConcernedDestination {
		return 0
	}
	return p.decodeUint32ValData() & 0xffffff
}

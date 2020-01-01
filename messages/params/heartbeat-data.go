// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewHeartbeatData creates the HeartbeatData Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewHeartbeatData(data []byte) *Param {
	return newVariableLenValParam(HeartbeatData, data)
}

// HeartbeatData returns multiple HeartbeatData from Param.
func (p *Param) HeartbeatData() []byte {
	if p.Tag != HeartbeatData {
		return nil
	}
	return p.Data
}

// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// TrafficModeType definitions.
const (
	_ uint32 = iota
	TrafficModeOverride
	TrafficModeLoadshare
	TrafficModeBroadcast
)

// NewTrafficModeType creates the TrafficModeType Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewTrafficModeType(tmType uint32) *Param {
	return newUint32ValParam(TrafficModeType, tmType)
}

// TrafficModeType returns multiple TrafficModeType from Param.
func (p *Param) TrafficModeType() uint32 {
	if p.Tag != TrafficModeType {
		return 0
	}

	return p.decodeUint32ValData()
}

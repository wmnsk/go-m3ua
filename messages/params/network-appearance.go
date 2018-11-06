// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewNetworkAppearance creates the NetworkAppearance Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewNetworkAppearance(nwApr uint32) *Param {
	return newUint32ValParam(NetworkAppearance, nwApr)
}

// NetworkAppearance returns multiple NetworkAppearance from Param.
func (p *Param) NetworkAppearance() uint32 {
	if p.Tag != NetworkAppearance {
		return 0
	}
	return p.decodeUint32ValData()
}

// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewServiceIndicators creates the ServiceIndicators Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewServiceIndicators(si ...uint8) *Param {
	return newMultiUint8ValParam(ServiceIndicators, si...)
}

// ServiceIndicators returns multiple ServiceIndicators from Param.
func (p *Param) ServiceIndicators() []uint8 {
	if p.Tag != ServiceIndicators {
		return nil
	}

	return p.Data
}

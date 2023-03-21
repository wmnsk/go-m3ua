// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewInfoString creates the InfoString Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewInfoString(infoStr string) *Param {
	return newVariableLenValParam(InfoString, []byte(infoStr))
}

// InfoString returns multiple InfoString from Param.
func (p *Param) InfoString() string {
	if p.Tag != InfoString {
		return ""
	}
	return string(p.Data)
}

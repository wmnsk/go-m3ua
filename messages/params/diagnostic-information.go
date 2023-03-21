// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewDiagnosticInformation creates the DiagnosticInformation Parameter.
// Note that this returns *Param, as no specific structure in this parameter.
func NewDiagnosticInformation(di []byte) *Param {
	return newVariableLenValParam(DiagnosticInformation, di)
}

// DiagnosticInformation returns multiple DiagnosticInformation from Param.
func (p *Param) DiagnosticInformation() []byte {
	if p.Tag != DiagnosticInformation {
		return nil
	}
	return p.Data
}

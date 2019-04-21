// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewAffectedPointCode creates the AffectedPointCode Parameter.
// Multiple number of AffectedPointCode will be accepted, but
// the mask for each point code should be included inside arguments.
// TODO: Handle masks and point codes separately.
// Note that this returns *Param, as no specific structure in this parameter.
func NewAffectedPointCode(apcs ...uint32) *Param {
	return newMultiUint32ValParam(AffectedPointCode, apcs...)
}

// AffectedPointCode returns single AffectedPointCode from Param.
func (p *Param) AffectedPointCode() uint32 {
	if p.Tag != AffectedPointCode {
		return 0
	}
	return p.AffectedPointCodes()[0]
}

// AffectedPointCodes returns multiple AffectedPointCode from Param.
func (p *Param) AffectedPointCodes() []uint32 {
	if p.Tag != AffectedPointCode {
		return nil
	}
	return p.decodeMultiUint32ValData()
}

/* TODO: Might be implemented in the following way?
// PointCodeWithMask is a set of Mask and Point Code.
type PointCodeWithMask struct {
	Mask      uint8
	PointCode uint32
}

// Serialize creates the 32bit-sized []byte from PointCodeWithMask.
func (p *PointCodeWithMask) Serialize() ([]byte, error) {
	b := make([]byte, 4)
	// to be written?
}

func (p *PointCodeWithMask) SerializeTo(b []bytes) error {
	// to be written?
}

func (p *PointCodeWithMask) Decode(b []bytes) (*PointCodeWithMask, error) {
	// to be written?
}

func (p *PointCodeWithMask) DecodeFromBytes(b []bytes) error {
	// to be written?
}
*/

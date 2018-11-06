// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewOriginatingPointCodeList creates the OriginatingPointCodeList Parameter.
// Multiple number of OriginatingPointCodeList will be accepted, but
// the mask for each point code should be included inside arguments.
// TODO: Handle masks and point codes separately.
// Note that this returns *Param, as no specific structure in this parameter.
func NewOriginatingPointCodeList(opcs ...uint32) *Param {
	return newMultiUint32ValParam(OriginatingPointCodeList, opcs...)
}

// OriginatingPointCodeList returns multiple OriginatingPointCode from Param.
func (p *Param) OriginatingPointCodeList() []uint32 {
	if p.Tag != OriginatingPointCodeList {
		return nil
	}
	return p.decodeMultiUint32ValData()
}

/* XXX - Might be implemented in the following way?
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

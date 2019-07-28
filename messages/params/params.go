// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

// Common Parameter Tag definitions.
const (
	_ uint16 = iota
	_
	_
	_
	InfoString
	_
	RoutingContext
	DiagnosticInformation
	_
	HeartbeatData
	_
	TrafficModeType
	ErrorCode
	Status
	_
	_
	_
	AspIdentifier
	AffectedPointCode
	CorrelationID
)

// M3UA-specific Parameter Tag definitions.
const (
	NetworkAppearance uint16 = uint16(0x200 | iota)
	_
	_
	_
	UserCause
	CongestionIndications
	ConcernedDestination
	RoutingKey           // specific: later
	RegistrationResult   // specific: later
	DeregistrationResult // specific: later
	LocalRoutingKeyIdentifier
	DestinationPointCode
	ServiceIndicators
	_
	OriginatingPointCodeList
	_
	ProtocolData
	_
	RegistrationStatus
	DeregistrationStatus
)

// Error definitions.
var (
	ErrInvalidType             = errors.New("got invalid type in parameter")
	ErrInvalidLength           = errors.New("parameter has invalid length value")
	ErrTooShortToMarshalBinary = errors.New("insufficient buffer to serialize parameter to")
	ErrTooShortToParse         = errors.New("too short to decode as parameter")
)

// Param is a M3UA Param.
type Param struct {
	Tag    uint16
	Length uint16
	Data   []byte
}

func newUint32ValParam(t uint16, u uint32) *Param {
	p := &Param{
		Tag:    t,
		Length: 8,
		Data:   make([]byte, 4),
	}
	binary.BigEndian.PutUint32(p.Data, u)
	return p
}

func newUint24ValParam(t uint16, u uint32) *Param {
	p := &Param{
		Tag:    t,
		Length: 8,
		Data:   make([]byte, 1),
	}
	p.Data = append(p.Data, uint32To24(u)...)
	return p
}

func uint32To24(n uint32) []byte {
	return []byte{uint8(n >> 16), uint8(n >> 8), uint8(n)}
}

func newUint8ValParam(t uint16, u uint8) *Param {
	return &Param{
		Tag:    t,
		Length: 8,
		Data:   []byte{0, 0, 0, u},
	}
}

func newMultiUint32ValParam(t uint16, ux ...uint32) *Param {
	p := &Param{
		Tag: t,
	}

	p.Data = make([]byte, len(ux)*4)
	for i, u := range ux {
		binary.BigEndian.PutUint32(p.Data[i*4:(i+1)*4], u)
	}
	p.SetLength()
	return p
}

func newMultiUint8ValParam(t uint16, ux ...uint8) *Param {
	l := len(ux) + (4 - len(ux)%4)
	p := &Param{
		Tag:  t,
		Data: make([]byte, l),
	}

	copy(p.Data, ux)
	p.SetLength()
	return p
}

func newVariableLenValParam(t uint16, b []byte) *Param {
	p := &Param{
		Tag:  t,
		Data: b,
	}
	p.SetLength()
	return p
}

func newNestedParam(t uint16, ps ...*Param) *Param {
	p := &Param{
		Tag: t,
	}

	for _, pr := range ps {
		if pr != nil {
			x, _ := pr.MarshalBinary()
			p.Data = append(p.Data, x...)
		}
	}
	p.SetLength()
	return p
}

func (p *Param) decodeUint32ValData() uint32 {
	l := len(p.Data)
	if l != 4 {
		return 0
	}
	return binary.BigEndian.Uint32(p.Data)
}

func (p *Param) decodeMultiUint32ValData() []uint32 {
	l := len(p.Data)
	if l%4 != 0 {
		return nil
	}

	var us []uint32
	for i := 0; i < l/4; i++ {
		us = append(us, binary.BigEndian.Uint32(p.Data[i*4:(i+1)*4]))
	}
	return us
}

// NewParam creates a new Param.
// This is for generic use. NewXXX(ParamName) functions are available to create the parameters defined in RFC4666.
func NewParam(tag int, data []byte) *Param {
	p := &Param{
		Tag:  uint16(tag),
		Data: data,
	}
	p.SetLength()
	return p
}

// MarshalBinary creates the byte sequence generated from a M3UA Param instance.
func (p *Param) MarshalBinary() ([]byte, error) {
	b := make([]byte, p.MarshalLen())
	if err := p.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (p *Param) MarshalTo(b []byte) error {
	binary.BigEndian.PutUint16(b[0:2], p.Tag)
	binary.BigEndian.PutUint16(b[2:4], p.Length)
	copy(b[4:p.MarshalLen()], p.Data)
	return nil
}

// Parse decodes given byte sequence as a M3UA Param.
func Parse(b []byte) (*Param, error) {
	p := &Param{}
	if err := p.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return p, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA Param.
func (p *Param) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < 4 {
		return ErrTooShortToParse
	}

	p.Tag = binary.BigEndian.Uint16(b[0:2])
	p.Length = binary.BigEndian.Uint16(b[2:4])
	if int(p.Length) > l {
		return ErrInvalidLength
	}

	p.Data = b[4:p.Length]
	return nil
}

// Padding creates the padded length of a M3UA Param.
func (p *Param) Padding() int {
	x := len(p.Data) % 4
	if x == 0 {
		return 0
	}
	return 4 - x
}

// MarshalLen returns serial length in integer.
func (p *Param) MarshalLen() int {
	return 4 + len(p.Data) + p.Padding()
}

// SetLength sets the length in Length field.
func (p *Param) SetLength() {
	p.Length = uint16(4 + len(p.Data))
}

// String creates the M3UA header values in human readable format.
func (p *Param) String() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("{Tag: %d, Length: %d, Data: %x}",
		p.Tag,
		p.Length,
		p.Data,
	)
}

// MarshalMultiParams creates the byte sequence from multiple Param instances.
func MarshalMultiParams(params []*Param) ([]byte, error) {
	var b []byte
	for _, param := range params {
		c, err := param.MarshalBinary()
		if err != nil {
			return nil, err
		}
		b = append(b, c...)
	}
	return b, nil
}

// ParseMultiParams decodes multiple Params at a time.
//
// This is easy and useful but slower than decoding one by one.
// When you don't know the number of Params, this is the only way to decode them.
// See benchmarks in diameter_test.go for the detail.
func ParseMultiParams(b []byte) ([]*Param, error) {
	var prms []*Param
	for {
		if len(b) == 0 {
			break
		}

		p, err := Parse(b)
		if err != nil {
			return nil, err
		}
		prms = append(prms, p)
		if len(b) < int(p.Length)+p.Padding() {
			return prms, nil
		}

		b = b[int(p.Length)+p.Padding():]
	}
	return prms, nil
}

// Serialize returns the byte sequence generated from a Param.
//
// DEPRECATED: use MarshalBinary instead.
func (p *Param) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return p.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (p *Param) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return p.MarshalTo(b)
}

// Decode decodes given byte sequence as a Param.
//
// DEPRECATED: use Parse instead.
func Decode(b []byte) (*Param, error) {
	log.Println("DEPRECATED: use Parse instead")
	return Parse(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (p *Param) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return p.UnmarshalBinary(b)
}

// Len returns the serial length of Param.
//
// DEPRECATED: use MarshalLen instead.
func (p *Param) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return p.MarshalLen()
}

// SerializeMultiParams creates the byte sequence from multiple Param instances.
//
// DEPRECATED: use MarshalMultiParams instead.
func SerializeMultiParams(params []*Param) ([]byte, error) {
	log.Println("DEPRECATED: use MarshalMultiParams instead")
	return MarshalMultiParams(params)
}

// DecodeMultiParams decodes multiple Params at a time.
//
// This is easy and useful but slower than decoding one by one.
// When you don't know the number of Params, this is the only way to decode them.
// See benchmarks in diameter_test.go for the detail.
//
// DEPRECATED: use ParseMultiParams instead.
func DecodeMultiParams(b []byte) ([]*Param, error) {
	log.Println("DEPRECATED: use ParseMultiParams instead")
	return ParseMultiParams(b)
}

// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package pc provides Point Code converting from some variants and translation to IP.
package pc

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Variant is a variant of Signaling Point Code represented in string.
type Variant string

// PointCode variant definitions.
const (
	VariantNone Variant = ""
	Variant383  Variant = "3-8-3"   // ITU
	Variant437  Variant = "4-3-7"   // ITU
	Variant4343 Variant = "4-3-4-3" // ITU
	Variant446  Variant = "4-4-6"   // ??
	Variant545  Variant = "5-4-5"   // ??
	Variant662  Variant = "6-6-2"   // ??
	Variant68   Variant = "6-8"     // ??
	Variant745  Variant = "7-4-5"   // Japan
	Variant77   Variant = "7-7"     // ??
	Variant888  Variant = "8-8-8"   // ANSI & China
)

// BitLength returns the defined bit length of Variant in int.
func (v Variant) BitLength() int {
	switch v {
	case Variant383, Variant437, Variant4343, Variant545, Variant662, Variant68, Variant77:
		return 14
	case Variant745:
		return 16
	case Variant888:
		return 24
	default:
		return 0
	}
}

func (v Variant) slice() []uint32 {
	if v == VariantNone {
		return nil
	}

	var s []uint32
	for _, digit := range strings.Split(v.String(), "-") {
		d, err := strconv.Atoi(digit)
		if err != nil {
			return nil
		}
		s = append(s, uint32(d))
	}
	return s
}

// String returns Variant in string representation.
func (v Variant) String() string {
	return string(v)
}

// PointCode represents a Signaling Point Code with its variant.
type PointCode struct {
	raw       uint32
	formatted string
	form      Variant
}

// NewPointCode creates a new PointCode from raw(uint32) value.
func NewPointCode(raw uint32, variant Variant) *PointCode {
	p := &PointCode{
		raw: raw, form: variant,
	}
	// apply bitmask
	p.raw &= (1 << uint32(variant.BitLength())) - 1

	var err error
	p.formatted, err = p.ConvertTo(variant)
	if err != nil {
		return nil
	}
	return p
}

// NewPointCodeFrom creates a new PointCode from formatted Signaling Point Code.
func NewPointCodeFrom(pc string, variant Variant) *PointCode {
	if variant == VariantNone {
		return nil
	}

	raw, err := convStrToRaw(pc, variant)
	if err != nil {
		return nil
	}
	return &PointCode{
		raw: raw, formatted: pc, form: variant,
	}
}

// Uint32 returns PointCode values in uint32.
func (pc *PointCode) Uint32() uint32 {
	return pc.raw
}

// Variant returns the variant of PointCode in string.
func (pc *PointCode) Variant() string {
	return pc.form.String()
}

// ConvertTo converts raw Signaling Point Code into specified Variant
// and returns converted PC value in string.
// The converted value is stored in PointCode and can be retrieved with
// String() without re-calculation.
func (pc *PointCode) ConvertTo(variant Variant) (string, error) {
	str, err := convRawToStr(pc.raw, variant)
	if err != nil {
		return "", err
	}
	pc.formatted = str
	return str, nil
}

func convRawToStr(n uint32, v Variant) (string, error) {
	if v == VariantNone {
		return "", errors.New("invalid Variant given")
	}

	s := v.slice()
	r := uint32(v.BitLength())
	n &= (1 << r) - 1 // apply bitmask

	d := make([]string, len(s))
	for i, v := range s {
		x := n & ((1 << r) - (1 << (r - v)))
		r -= v
		d[i] = strconv.Itoa(int(x >> r))
	}

	return strings.Join(d, "-"), nil
}

func convStrToRaw(f string, v Variant) (uint32, error) {
	if v == VariantNone {
		return 0, errors.New("invalid Variant given")
	}

	ds := strings.Split(f, "-")
	if len(ds) == 0 {
		return 0, errors.Errorf("PC: %s is invalid; digits should be splitted with \"-\"", f)
	}
	s := v.slice()
	if len(ds) != len(s) {
		return 0, errors.Errorf("PC: %s and Variant: %s doesn't match", f, v)
	}

	r := uint32(v.BitLength())
	var n uint32
	for i, d := range ds {
		x, err := strconv.Atoi(d)
		if err != nil {
			return 0, errors.Wrap(err, "failed to convert PC")
		}
		r -= s[i]
		n |= (uint32(x) << r)
	}

	return n, nil
}

// String returns PointCode in formatted string.
func (pc *PointCode) String() string {
	if pc.form == VariantNone {
		return ""
	}

	return pc.formatted
}

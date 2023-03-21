// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
This utility is to simplify the test cases used throughout the package,
much inspired from github.com/wmnsk/gopcua/blob/master/utils/codectest/testcase.go,
which is written by @magiconair (https://github.com/magiconair).
*/

package messages

import (
	"testing"

	"github.com/pascaldekloe/goe/verify"
)

type serializeable interface {
	MarshalBinary() ([]byte, error)
	MarshalLen() int
}

type testCase struct {
	name       string
	structured serializeable
	serialized []byte
}

type decoderFunc func([]byte) (serializeable, error)

func runTests(t *testing.T, cases []testCase, decode decoderFunc) {
	t.Helper()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Run("decode", func(t *testing.T) {
				v, err := decode(c.serialized)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := v, c.structured; !verify.Values(t, "", got, want) {
					t.Fail()
				}
			})

			t.Run("encode", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}

				if got, want := b, c.serialized; !verify.Values(t, "", got, want) {
					t.Fail()
				}
			})

			t.Run("len", func(t *testing.T) {
				if got, want := c.structured.MarshalLen(), len(c.serialized); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
			})

			t.Run("interface", func(t *testing.T) {
				// Ignore *Header in this tests.
				if _, ok := c.structured.(*Header); ok {
					return
				}
				decoded, err := Parse(c.serialized)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := decoded.MessageClass(), c.structured.(M3UA).MessageClass(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
				if got, want := decoded.MessageClassName(), c.structured.(M3UA).MessageClassName(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
				if got, want := decoded.MessageType(), c.structured.(M3UA).MessageType(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
				if got, want := decoded.MessageTypeName(), c.structured.(M3UA).MessageTypeName(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
			})
		})
	}
}

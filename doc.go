// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Package m3ua provides easy and painless handling of M3UA protocol in pure Golang.

The API design is kept as similar as possible to other protocols in standard net package.
To establish M3UA connection as client/server, you can use Dial() and Listen() / Accept()
without caring about the underlying SCTP association, as go-m3ua handles it together
with M3UA ASPSM & ASPTM procedures.

This package relies much on github.com/ishidawataru/sctp, as M3UA requires underlying SCTP connection,

Specification: https://tools.ietf.org/html/rfc4666
*/
package m3ua

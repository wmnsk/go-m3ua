// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// XXX - implement!

package m3ua

// PCFormat represents the format of PointCode.
//
// NOT IMPLEMENTED YET!
type PCFormat string

// PointCode Format definitions.
//
// NOT IMPLEMENTED YET!
const (
	PCFormat3_2_3 = "3-2-3"
)

// ParsePC parses formatted PointCode as uint32.
//
// NOT IMPLEMENTED YET!
func ParsePC(pc string, format PCFormat) uint32 {
	return 0
}

// FormatPC formats decimal(uint32) PointCode as given format.
//
// NOT IMPLEMENTED YET!
func FormatPC(pc uint32, format PCFormat) string {
	return ""
}

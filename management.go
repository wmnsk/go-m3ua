// Copyright 2018 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import "github.com/wmnsk/go-m3ua/messages"

// XXX - implement!
func (c *Conn) handleError(e *messages.Error) error {
	switch c.state {
	case StateSCTPCDI, StateSCTPRI:
		return NewErrUnexpectedMessage(e)
	}

	return nil
}

// XXX - implement!
func (c *Conn) handleNotify(e *messages.Notify) error {
	switch c.state {
	case StateSCTPCDI, StateSCTPRI:
		return NewErrUnexpectedMessage(e)
	}

	return nil
}

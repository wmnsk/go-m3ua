// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import "github.com/wmnsk/go-m3ua/messages"

func (c *Conn) initiateASPSM() error {
	if _, err := c.WriteSignal(
		messages.NewAspUp(c.cfg.AspIdentifier, nil),
	); err != nil {
		return err
	}

	return nil
}
func (c *Conn) handleAspUp(aspUp *messages.AspUp) error {
	if c.State() != StateAspDown {
		defer c.Close() // Provided to handle bugs from peer STPs
		return NewErrUnexpectedMessage(aspUp)

	}
	if c.StreamID() != 0 {
		return NewErrInvalidSCTPStreamID(c.StreamID())
	}

	if _, err := c.WriteSignal(
		messages.NewAspUpAck(
			c.cfg.AspIdentifier,
			nil,
		),
	); err != nil {
		return err
	}

	return nil
}

func (c *Conn) handleAspUpAck(aspUpAck *messages.AspUpAck) error {
	if c.State() != StateAspDown {
		return NewErrUnexpectedMessage(aspUpAck)
	}
	if c.StreamID() != 0 {
		return NewErrInvalidSCTPStreamID(c.StreamID())
	}

	return nil
}

func (c *Conn) handleAspDown(aspDown *messages.AspDown) error {
	c.Close() // Closing the connection to close the dataChan, to avoid the Read function keeping blocked on the channel.
	switch c.State() {
	case StateAspInactive, StateAspActive:
		return NewErrUnexpectedMessage(aspDown)
	}
	if c.StreamID() != 0 {
		return NewErrInvalidSCTPStreamID(c.StreamID())
	}

	// XXX - Validate the params.

	if _, err := c.WriteSignal(messages.NewAspDownAck(nil)); err != nil {
		return err
	}

	return nil
}

func (c *Conn) handleAspDownAck(aspDownAck *messages.AspDownAck) error {
	switch c.State() {
	case StateAspInactive, StateAspActive:
		return NewErrUnexpectedMessage(aspDownAck)
	}
	if c.StreamID() != 0 {
		return NewErrInvalidSCTPStreamID(c.StreamID())
	}

	return nil
}

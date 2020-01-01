// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/wmnsk/go-m3ua/messages"
	"github.com/wmnsk/go-m3ua/messages/params"
)

func (c *Conn) initiateASPTM() error {
	if _, err := c.WriteSignal(messages.NewAspActive(
		c.cfg.TrafficModeType, c.cfg.RoutingContexts, nil,
	)); err != nil {
		return err
	}

	return nil
}

func (c *Conn) heartbeat(ctx context.Context) {
	data := make([]byte, 128)
	beat := messages.NewHeartbeat(params.NewHeartbeatData(data))
	for {
		if _, err := rand.Read(data); err != nil {
			c.errChan <- err
			return
		}
		beat.HeartbeatData = params.NewHeartbeatData(data)
		if _, err := c.WriteSignal(beat); err != nil {
			c.errChan <- ErrFailedToWriteSignal
			return
		}
		c.cfg.HeartbeatInfo.Data = data

		// wait for response
		select {
		case <-ctx.Done():
			return
		case _, ok := <-c.beatAckChan: // got valid BEAT response from peer
			if !ok {
				return
			}
			break
		case <-time.After(c.cfg.HeartbeatInfo.Timer): // timer expired
			c.errChan <- ErrHeartbeatExpired
			return
		}

		// wait while next time
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.cfg.HeartbeatInfo.Interval):
			continue
		}
	}
}

func (c *Conn) handleAspActive(aspActive *messages.AspActive) error {
	if c.state != StateAspInactive {
		return NewErrUnexpectedMessage(aspActive)
	}

	if _, err := c.WriteSignal(
		messages.NewAspActiveAck(c.cfg.TrafficModeType, c.cfg.RoutingContexts, nil),
	); err != nil {
		return err
	}

	return nil
}

func (c *Conn) handleAspActiveAck(aspAcAck *messages.AspActiveAck) error {
	if c.state != StateAspInactive {
		return NewErrUnexpectedMessage(aspAcAck)
	}

	// XXX - Add some additional validation for aspAcAck here.

	return nil
}

func (c *Conn) handleAspInactive(aspInactive *messages.AspInactive) error {
	if c.state != StateAspActive {
		return NewErrUnexpectedMessage(aspInactive)
	}

	if _, err := c.WriteSignal(
		messages.NewAspInactiveAck(c.cfg.RoutingContexts, nil),
	); err != nil {
		return err
	}

	return nil
}

func (c *Conn) handleAspInactiveAck(aspAcAck *messages.AspInactiveAck) error {
	if c.state != StateAspActive {
		return NewErrUnexpectedMessage(aspAcAck)
	}

	// XXX - Add some additional validation for aspAcAck here.

	return nil
}

func (c *Conn) handleHeartbeat(beat *messages.Heartbeat) error {
	if c.state != StateAspActive {
		return NewErrUnexpectedMessage(beat)
	}

	// No need to create new HeartbeatAck, as it's identical to Heartbeat except the MessageType.
	beat.Type = messages.MsgTypeHeartbeatAck
	if _, err := c.WriteSignal(beat); err != nil {
		return err
	}
	return nil
}

func (c *Conn) handleHeartbeatAck(beatAck *messages.HeartbeatAck) error {
	if c.state != StateAspActive {
		return NewErrUnexpectedMessage(beatAck)
	}

	myData := c.cfg.HeartbeatInfo.Data
	dataFromPeer := beatAck.HeartbeatData.HeartbeatData()
	if len(dataFromPeer) != len(myData) {
		return NewErrUnexpectedMessage(beatAck)
	}
	for i, p := range dataFromPeer {
		if p != myData[i] {
			return NewErrUnexpectedMessage(beatAck)
		}
	}

	return nil
}

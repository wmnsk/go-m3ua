// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// State represents ASP State.
type State uint8

// M3UA status definitions.
const (
	StateAspDown State = iota
	StateAspInactive
	StateAspActive
	StateSCTPCDI
	StateSCTPRI
)

func (c *Conn) handleStateUpdate(current State) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	previous := c.state
	c.state = current

	switch c.mode {
	case modeClient:
		if err := c.handleStateUpdateAsClient(current, previous); err != nil {
			return err
		}
		return nil
	case modeServer:
		if err := c.handleStateUpdateAsServer(current, previous); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("not implemented yet")
	}
}

func (c *Conn) handleStateUpdateAsClient(current, previous State) error {
	switch current {
	case StateAspDown:
		c.sctpInfo.Stream = 0
		return c.initiateASPSM()
	case StateAspInactive:
		return c.initiateASPTM()
	case StateAspActive:
		if current != previous {
			c.established <- struct{}{}
			c.beatAllow.Broadcast()
		}
		return nil
	case StateSCTPCDI, StateSCTPRI:
		return ErrSCTPNotAlive
	default:
		return ErrInvalidState
	}
}

func (c *Conn) handleStateUpdateAsServer(current, previous State) error {
	switch current {
	case StateAspDown:
		// do nothing. just wait for the message from peer and state is updated
		return nil
	case StateAspInactive:
		// do nothing. just wait for the message from peer and state is updated
		// XXX - send DAVA to notify peer?
		return nil
	case StateAspActive:
		if current != previous {
			c.established <- struct{}{}
			c.beatAllow.Broadcast()
		}
		return nil
	case StateSCTPCDI, StateSCTPRI:
		return ErrSCTPNotAlive
	default:
		return ErrInvalidState
	}
}

func (c *Conn) handleSignals(ctx context.Context, m3 messages.M3UA) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	// Signal validations
	if m3.Version() != 1 {
		c.errChan <- NewErrInvalidVersion(m3.Version())
		return
	}

	switch msg := m3.(type) {
	// Transfer message
	case *messages.Data:
		go c.handleData(ctx, msg)
		c.stateChan <- c.state
	// ASPSM
	case *messages.AspUp:
		if err := c.handleAspUp(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- StateAspInactive
	case *messages.AspUpAck:
		if err := c.handleAspUpAck(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- StateAspInactive
	case *messages.AspDown:
		if err := c.handleAspDown(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- StateAspDown
	case *messages.AspDownAck:
		if err := c.handleAspDownAck(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- StateAspDown
	// ASPTM
	case *messages.AspActive:
		if err := c.handleAspActive(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- StateAspActive
	case *messages.AspActiveAck:
		if err := c.handleAspActiveAck(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- StateAspActive
	case *messages.AspInactive:
		if err := c.handleAspInactive(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- StateAspInactive
	case *messages.AspInactiveAck:
		if err := c.handleAspInactiveAck(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- StateAspInactive
	case *messages.Heartbeat:
		if err := c.handleHeartbeat(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- c.state
	case *messages.HeartbeatAck:
		if err := c.handleHeartbeatAck(msg); err != nil {
			c.errChan <- err
		}
		c.beatAckChan <- struct{}{}
		c.stateChan <- c.state
		// Management
	case *messages.Error:
		if err := c.handleError(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- c.state
	case *messages.Notify:
		if err := c.handleNotify(msg); err != nil {
			c.errChan <- err
		}
		c.stateChan <- c.state
	// Others: SSNM and RKM is not implemented.
	default:
		c.errChan <- NewErrUnsupportedMessage(m3)
	}
}

func (c *Conn) monitor(ctx context.Context) {
	c.errChan = make(chan error)
	c.dataChan = make(chan *params.ProtocolDataPayload, 0xffff)
	c.beatAckChan = make(chan struct{})
	c.beatAllow = sync.NewCond(&sync.Mutex{})
	c.beatAllow.L.Lock()
	go c.heartbeat(ctx)
	defer c.beatAllow.Broadcast()
	buf := make([]byte, 0xffff)
	for {
		select {
		case <-ctx.Done():
			c.Close()
			return
		case err := <-c.errChan:
			if e := c.handleErrors(err); e != nil {
				c.Close()
				return
			}
			continue
		case state := <-c.stateChan:
			// Act properly based on current state.
			if err := c.handleStateUpdate(state); err != nil {
				if errors.Is(err, ErrSCTPNotAlive) {
					c.Close()
					return
				}
			}

			// Read from conn to see something coming from the peer.
			n, info, err := c.sctpConn.SCTPRead(buf)
			if err != nil {
				c.Close()
				return
			}
			if info != nil {
				if info.Stream != c.sctpInfo.Stream {
					c.Close()
					return
				}
			}

			// Parse the received packet as M3UA. Undecodable packets are ignored.
			msg, err := messages.Parse(buf[:n])
			if err != nil {
				continue
			}

			go c.handleSignals(ctx, msg)
		}
	}
}

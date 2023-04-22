// Copyright 2018-2023 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ishidawataru/sctp"
	"github.com/wmnsk/go-m3ua/messages"
	"github.com/wmnsk/go-m3ua/messages/params"
)

type mode uint8

const (
	modeClient mode = iota
	modeServer
)

// Conn represents a M3UA connection, which satisfies standard net.Conn interface.
type Conn struct {
	// mu is to Lock when updating state
	mu *sync.Mutex
	// mode represents the endpoint works as client or server
	mode mode
	// state is to see the current state
	state State
	// stateChan is to update the state and handle it
	stateChan chan State
	// established notifies client/server the conn is established
	established chan struct{}
	// beatAckChan notifies that heartbeat gets the ack as expected
	beatAckChan chan struct{}
	// dataChan is to pass the ProtocolDataPayload(=payload on M3UA DATA) to user
	dataChan chan *params.ProtocolDataPayload
	// errChan is to pass errors to goroutine that monitors status
	errChan chan error
	// sctpConn is the underlying SCTP association
	sctpConn *sctp.SCTPConn
	// sctpInfo is SndRcvInfo in SCTP association
	sctpInfo *sctp.SndRcvInfo
	// cfg is a configuration that is required to communicate between M3UA endpoints
	cfg *Config
	// Condition to allow heartbeat, only after the state is AspUp
	beatAllow *sync.Cond
}

var netMap = map[string]string{
	"m3ua":  "sctp",
	"m3ua4": "sctp4",
	"m3ua6": "sctp6",
}

// Read reads data from the connection.
func (c *Conn) Read(b []byte) (n int, err error) {
	err = func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.state != StateAspActive {
			return ErrNotEstablished
		}
		return nil
	}()
	if err != nil {
		return 0, err
	}

	pd, ok := <-c.dataChan
	if !ok {
		return 0, ErrNotEstablished
	}

	copy(b, pd.Data)
	return len(pd.Data), nil

}

// Write writes data to the connection.
func (c *Conn) Write(b []byte) (n int, err error) {
	if c.state != StateAspActive {
		return 0, ErrNotEstablished
	}
	d, err := messages.NewData(
		c.cfg.NetworkAppearance, c.cfg.RoutingContexts, params.NewProtocolData(
			c.cfg.OriginatingPointCode, c.cfg.DestinationPointCode,
			c.cfg.ServiceIndicator, c.cfg.NetworkIndicator,
			c.cfg.MessagePriority, c.cfg.SignalingLinkSelection, b,
		), c.cfg.CorrelationID,
	).MarshalBinary()
	if err != nil {
		return 0, err
	}

	n, err = c.sctpConn.SCTPWrite(d, c.sctpInfo)
	if err != nil {
		return 0, err
	}

	n += len(d)
	return n, nil
}

// Write writes data to the connection and specific stream
func (c *Conn) WriteToStream(b []byte, streamId uint16) (n int, err error) {
	if c.state != StateAspActive {
		return 0, ErrNotEstablished
	}
	d, err := messages.NewData(
		c.cfg.NetworkAppearance, c.cfg.RoutingContexts, params.NewProtocolData(
			c.cfg.OriginatingPointCode, c.cfg.DestinationPointCode,
			c.cfg.ServiceIndicator, c.cfg.NetworkIndicator,
			c.cfg.MessagePriority, c.cfg.SignalingLinkSelection, b,
		), c.cfg.CorrelationID,
	).MarshalBinary()
	if err != nil {
		return 0, err
	}

	c.sctpInfo.Stream = streamId
	n, err = c.sctpConn.SCTPWrite(d, c.sctpInfo)
	if err != nil {
		return 0, err
	}

	n += len(d)
	return n, nil
}

// WriteSignal writes any type of M3UA signals on top of SCTP Connection.
func (c *Conn) WriteSignal(m3 messages.M3UA) (n int, err error) {
	n = m3.MarshalLen()
	buf := make([]byte, n)
	if err := m3.MarshalTo(buf); err != nil {
		return 0, fmt.Errorf("failed to create %T: %w", m3, err)
	}

	nn, err := c.sctpConn.SCTPWrite(buf, c.sctpInfo)
	if err != nil {
		return 0, errors.Wrap(err, "failed to write M3UA")
	}

	n += nn
	return
}

// Close closes the connection.
func (c *Conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == StateAspDown {
		return c.sctpConn.Close()
	}

	close(c.established)
	close(c.beatAckChan)
	close(c.dataChan)
	c.state = StateAspDown
	return c.sctpConn.Close()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.sctpConn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.sctpConn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated.
func (c *Conn) SetDeadline(t time.Time) error {
	return c.sctpConn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.sctpConn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.sctpConn.SetWriteDeadline(t)
}

// State returns current state of Conn.
func (c *Conn) State() State {
	return c.state
}

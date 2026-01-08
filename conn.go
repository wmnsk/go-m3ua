// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/ishidawataru/sctp"
	"github.com/wmnsk/go-m3ua/messages"
	"github.com/wmnsk/go-m3ua/messages/params"
)

type mode uint8

const (
	modeClient mode = iota
	modeServer
)

// Conn represents a M3UA connection, which satisfies the standard net.Conn interface.
type Conn struct {
	// maxMessageStreamID is the maximum negotiated sctp stream ID used,
	// must not be zero, must vary from 1 to maxMessageStreamID
	maxMessageStreamID uint16
	// muState is to Lock when updating state
	muState *sync.RWMutex
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
	// dataChan is to pass the ProtocolDataPayload(=payload on M3UA DATA) to the user
	dataChan chan *params.ProtocolDataPayload
	// errChan is to pass errors to a goroutine that monitors status
	errChan chan error
	// cfg is a configuration required to communicate between M3UA endpoints
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
		if c.State() != StateAspActive {
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

// ReadPD reads the next ProtocolDataPayload from the connection.
func (c *Conn) ReadPD() (pd *params.ProtocolDataPayload, err error) {
	err = func() error {
		if c.State() != StateAspActive {
			return ErrNotEstablished
		}
		return nil
	}()
	if err != nil {
		return nil, err
	}

	pd, ok := <-c.dataChan
	if !ok {
		return nil, ErrNotEstablished
	}

	return pd, nil
}

// Write writes data to the connection.
func (c *Conn) Write(b []byte) (n int, err error) {
	stream := c.chooseStreamID()

	return c.WriteToStream(b, stream)
}

// WriteToStream writes data to the connection and specific stream
func (c *Conn) WriteToStream(b []byte, streamID uint16) (n int, err error) {
	if c.State() != StateAspActive {
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

	// taken by value to avoid race condition on the stream id
	info := *c.cfg.SCTPConfig.sctpInfo
	info.Stream = streamID
	n, err = c.cfg.SCTPConfig.sctpConn.SCTPWrite(d, &info)
	if err != nil {
		return 0, err
	}

	n += len(d)
	return n, nil
}

// WritePD writes data with a specific mtp3 protocol data to the connection.
func (c *Conn) WritePD(protocolData *params.Param) (n int, err error) {
	stream := c.chooseStreamID()

	return c.WritePDToStream(protocolData, stream)
}

// WritePDToStream writes data with a specific mtp3 protocol data to the connection and specific stream
func (c *Conn) WritePDToStream(protocolData *params.Param, streamID uint16) (n int, err error) {
	if c.State() != StateAspActive {
		return 0, ErrNotEstablished
	}
	d, err := messages.NewData(
		c.cfg.NetworkAppearance, // cannot be changed on an active connection
		c.cfg.RoutingContexts,   // cannot be changed on an active connection
		protocolData,            // custom mtp3 protocol data OPC, DPC, SI, NI, MP, and SLS, flexible on active connections
		c.cfg.CorrelationID,
	).MarshalBinary()
	if err != nil {
		return 0, err
	}

	// taken by value to avoid race condition on the stream id
	info := *c.cfg.SCTPConfig.sctpInfo
	info.Stream = streamID
	n, err = c.cfg.SCTPConfig.sctpConn.SCTPWrite(d, &info)
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

	// taken by value to avoid race condition on the stream id
	sctpInfo := *c.cfg.SCTPConfig.sctpInfo
	if m3.MessageClass() != messages.MsgClassTransfer {
		sctpInfo.Stream = 0
	}

	nn, err := c.cfg.SCTPConfig.sctpConn.SCTPWrite(buf, &sctpInfo)
	if err != nil {
		return 0, fmt.Errorf("failed to write M3UA: %w", err)
	}

	n += nn
	return
}

// Close closes the connection.
func (c *Conn) Close() error {
	c.muState.Lock()
	defer c.muState.Unlock()

	if c.state == StateAspDown {
		return c.cfg.SCTPConfig.sctpConn.Close()
	}

	close(c.established)
	close(c.beatAckChan)
	close(c.dataChan)
	c.state = StateAspDown
	return c.cfg.SCTPConfig.sctpConn.Close()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.cfg.SCTPConfig.sctpConn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.cfg.SCTPConfig.sctpConn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated.
func (c *Conn) SetDeadline(t time.Time) error {
	return c.cfg.SCTPConfig.sctpConn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.cfg.SCTPConfig.sctpConn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.cfg.SCTPConfig.sctpConn.SetWriteDeadline(t)
}

// State returns current state of Conn.
func (c *Conn) State() State {
	c.muState.RLock()
	defer c.muState.RUnlock()
	return c.state
}

// StreamID returns sctpInfo.Stream of Conn.
func (c *Conn) StreamID() uint16 {
	return c.cfg.SCTPConfig.sctpInfo.Stream
}

// MaxMessageStreamID returns the maximum negotiated sctp stream ID
// The streamID for sending a message must start from 1 up to maxMessageStreamID, 0 is reserved for management messages
func (c *Conn) MaxMessageStreamID() uint16 {
	return c.maxMessageStreamID
}

// SetSctpSackConfig sets the SCTP SACK timer configuration on an active connection.
//
// sackDelay is the number of milliseconds for the delayed SACK timer
// (per RFC4960, should be between 200 and 500 ms).
//
// sackFrequency is the number of packets to receive before sending a SACK
// without waiting for the delay timer. Setting to 1 disables the delayed
// SACK algorithm.
//
// Note: sackDelay=0, sackFrequency=1 (disables delayed SACK)
func (c *Conn) SetSctpSackConfig(sackDelay, sackFrequency uint32) error {
	return c.cfg.SCTPConfig.sctpConn.SetSackTimer(&sctp.SackTimer{
		SackDelay:     sackDelay,
		SackFrequency: sackFrequency,
	})
}

// chooseStreamID generates a random uint16 from 1 to max (inclusive)
func (c *Conn) chooseStreamID() uint16 {
	if c.maxMessageStreamID == 1 {
		return 1
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := uint16(r.Intn(int(c.maxMessageStreamID)))
	return randomNum + 1
}

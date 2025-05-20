// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"fmt"
	"log"
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

// Conn represents a M3UA connection, which satisfies standard net.Conn interface.
type Conn struct {
	// maxMessageStreamID is the maximum negotiated sctp stream ID used, must not be zero
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
	stream := RandomUint16(c.maxMessageStreamID) // choose a random stream number from 1 to a certain maximum

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
	info := *c.sctpInfo
	info.Stream = streamID
	n, err = c.sctpConn.SCTPWrite(d, &info)
	if err != nil {
		log.Printf("go-m3ua: error writing on sctp connection, stream id: %v, max negotiated stream id: %v, error : %v", streamID, c.maxMessageStreamID, err)
		return 0, err
	}

	n += len(d)
	return n, nil
}

// WritePD writes data with a specific mtp3 protocol data to the connection.
func (c *Conn) WritePD(protocolData *params.Param) (n int, err error) {
	stream := RandomUint16(c.maxMessageStreamID) // choose a random stream number from 1 to a certain maximum

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
	info := *c.sctpInfo
	info.Stream = streamID
	n, err = c.sctpConn.SCTPWrite(d, &info)
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
	sctpInfo := *c.sctpInfo
	if m3.MessageClass() != messages.MsgClassTransfer {
		sctpInfo.Stream = 0
	}

	nn, err := c.sctpConn.SCTPWrite(buf, &sctpInfo)
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
	c.muState.RLock()
	defer c.muState.RUnlock()
	return c.state
}

// StreamID returns sctpInfo.Stream of Conn.
func (c *Conn) StreamID() uint16 {
	return c.sctpInfo.Stream
}

func (c *Conn) MaxMessageStreamID() uint16 {
	return c.maxMessageStreamID
}

// RandomUint16 generates a random uint16 from 1 to max (inclusive)
// If max is 1, it always returns 1
func RandomUint16(max uint16) uint16 {
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// If max is 1, just return 1
	if max == 1 {
		return 1
	}

	// Generate a random number from 0 to (max-1)
	randomNum := uint16(r.Intn(int(max)))

	// Add 1 to get a number from 1 to max
	return randomNum + 1
}

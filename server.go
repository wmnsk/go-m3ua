// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ishidawataru/sctp"
)

// Listener is a M3UA listener.
type Listener struct {
	sctpListener *sctp.SCTPListener
	*Config
}

// Listen returns a M3UA listener.
func Listen(net string, laddr *sctp.SCTPAddr, cfg *Config) (*Listener, error) {
	var err error
	l := &Listener{Config: cfg}

	n, ok := netMap[net]
	if !ok {
		return nil, fmt.Errorf("invalid network: %s", net)
	}

	l.sctpListener, err = sctp.ListenSCTP(n, laddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen SCTP: %w", err)
	}
	return l, nil
}

// Accept waits for and returns the next connection to the listener.
// After successfully establishing the association with peer, Payload can be read with Read() func.
// Other signals are automatically handled background in another goroutine.
func (l *Listener) Accept(ctx context.Context) (*Conn, error) {
	conn := &Conn{
		muState:     new(sync.RWMutex),
		mode:        modeServer,
		stateChan:   make(chan State),
		established: make(chan struct{}),
		sctpInfo:    &sctp.SndRcvInfo{PPID: 0x03000000, Stream: 0},
		cfg:         l.Config,
	}

	if conn.cfg.HeartbeatInfo.Interval == 0 {
		conn.cfg.HeartbeatInfo.Enabled = false
	}

	c, err := l.sctpListener.Accept()
	if err != nil {
		return nil, err
	}

	var ok bool
	conn.sctpConn, ok = c.(*sctp.SCTPConn)
	if !ok {
		c.Close()
		return nil, fmt.Errorf("failed to assert server connection")
	}

	if conn.cfg.SctpSackInfo != nil && conn.cfg.SctpSackInfo.Enabled {
		err = conn.sctpConn.SetSackTimer(&sctp.SackTimer{
			SackDelay:     conn.cfg.SctpSackInfo.SackDelay,
			SackFrequency: conn.cfg.SctpSackInfo.SackFrequency,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to set sack timer: %w", err)
		}
	}

	r, err := conn.sctpConn.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get sctpConn status: %w", err)
	}
	conn.maxMessageStreamID = r.Ostreams - 1 // removing 1 for management messages of stream ID 0

	go func() {
		conn.stateChan <- StateAspDown
	}()

	go conn.monitor(ctx)
	select {
	case _, ok := <-conn.established:
		if !ok {
			conn.sctpConn.Close()
			return nil, ErrFailedToEstablish
		}
		return conn, nil
	case <-time.After(10 * time.Second):
		conn.sctpConn.Close()
		return nil, ErrTimeout
	}
}

// Close closes the listener.
func (l *Listener) Close() error {
	// XXX - should close on M3UA layer.
	return l.sctpListener.Close()
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.sctpListener.Addr()
}

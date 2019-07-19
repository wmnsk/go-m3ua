// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"

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
		return nil, errors.Wrap(err, "failed to listen SCTP")
	}
	return l, nil
}

// Accept waits for and returns the next connection to the listener.
// After successfully established the association with peer, Payload can be read with Read() func.
// Other signals are automatically handled background in another goroutine.
func (l *Listener) Accept(ctx context.Context) (*Conn, error) {
	conn := &Conn{
		mu:          new(sync.Mutex),
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
		return nil, errors.New("failed to assert conn")
	}

	go func() {
		conn.stateChan <- StateAspDown
	}()

	go conn.monitor(ctx)
	select {
	case _, ok := <-conn.established:
		if !ok {
			return nil, ErrFailedToEstablish
		}
		return conn, nil
	case <-time.After(10 * time.Second):
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

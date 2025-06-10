// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ishidawataru/sctp"
)

// Dial establishes a M3UA connection as a client.
// After successfully establishing the connection with peer, state-changing
// signals and heartbeats are automatically handled background in another goroutine.
func Dial(ctx context.Context, net string, laddr, raddr *sctp.SCTPAddr, cfg *Config) (*Conn, error) {
	var err error
	conn := &Conn{
		muState:     new(sync.RWMutex),
		mode:        modeClient,
		stateChan:   make(chan State),
		established: make(chan struct{}),
		sctpInfo:    &sctp.SndRcvInfo{PPID: 0x03000000, Stream: 0},
		cfg:         cfg,
	}

	if conn.cfg.HeartbeatInfo.Interval == 0 {
		conn.cfg.HeartbeatInfo.Enabled = false
	}

	n, ok := netMap[net]
	if !ok {
		return nil, fmt.Errorf("invalid network: %s", net)
	}

	conn.sctpConn, err = sctp.DialSCTP(n, laddr, raddr)
	if err != nil {
		if conn.sctpConn != nil {
			return nil, fmt.Errorf("go-m3ua: issue dialing connection. closing error: %w", conn.sctpConn.Close())
		}
		return nil, err
	}

	r, err := conn.sctpConn.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("go-m3ua: failed to retrive sctpConnection status for Dial: %w", err)
	}
	conn.maxMessageStreamID = r.Ostreams - 1 // removing 1 for management messages of stream ID 0

	go func() {
		conn.stateChan <- StateAspDown
	}()

	go conn.monitor(ctx)
	select {
	case _, ok := <-conn.established:
		if !ok {
			return nil, fmt.Errorf("go-m3ua: issue having established client connection. error: %w, closing error: %w", ErrFailedToEstablish, conn.sctpConn.Close())
		}
		return conn, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("go-m3ua: issue client connection timeout. error: %w, closing error: %w", ErrTimeout, conn.sctpConn.Close())
	}
}

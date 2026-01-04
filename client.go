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

	if cfg.SCTPConfig == nil {
		cfg.SCTPConfig = &SCTPConfig{}
	}
	cfg.SCTPConfig.sctpInfo = &sctp.SndRcvInfo{PPID: 0x03000000, Stream: 0}

	conn := &Conn{
		muState:     new(sync.RWMutex),
		mode:        modeClient,
		stateChan:   make(chan State),
		established: make(chan struct{}),
		cfg:         cfg,
	}

	if conn.cfg.HeartbeatInfo.Interval == 0 {
		conn.cfg.HeartbeatInfo.Enabled = false
	}

	n, ok := netMap[net]
	if !ok {
		return nil, fmt.Errorf("invalid network: %s", net)
	}

	conn.cfg.SCTPConfig.sctpConn, err = sctp.DialSCTP(n, laddr, raddr)
	if err != nil {
		return nil, err
	}

	if conn.cfg.SCTPConfig.SctpSackInfo != nil && conn.cfg.SCTPConfig.SctpSackInfo.Enabled {
		err = conn.cfg.SCTPConfig.sctpConn.SetSackTimer(&sctp.SackTimer{
			SackDelay:     conn.cfg.SCTPConfig.SctpSackInfo.SackDelay,
			SackFrequency: conn.cfg.SCTPConfig.SctpSackInfo.SackFrequency,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to set sack timer: %w", err)
		}
	}

	r, err := conn.cfg.SCTPConfig.sctpConn.GetStatus()
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
			conn.cfg.SCTPConfig.sctpConn.Close()
			return nil, ErrFailedToEstablish
		}
		return conn, nil
	case <-time.After(10 * time.Second):
		conn.cfg.SCTPConfig.sctpConn.Close()
		return nil, ErrTimeout
	}
}

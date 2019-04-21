// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ishidawataru/sctp"
	"github.com/wmnsk/go-m3ua/messages/params"
)

// NewClientConfig creates a new Config for Client.
//
// The optional parameters that is not required (like CorrelationID),
// omit it by setting it to nil after created *Config.
func NewClientConfig(hbInfo *HeartbeatInfo, opc, dpc, aspID, tmt, nwApr, corrID uint32, rtCtxs []uint32, si, ni, mp, sls uint8) *Config {
	return &Config{
		HeartbeatInfo:          hbInfo,
		AspIdentifier:          params.NewAspIdentifier(aspID),
		TrafficModeType:        params.NewTrafficModeType(tmt),
		NetworkAppearance:      params.NewNetworkAppearance(nwApr),
		RoutingContexts:        params.NewRoutingContext(rtCtxs...),
		CorrelationID:          params.NewCorrelationID(corrID),
		OriginatingPointCode:   opc,
		DestinationPointCode:   dpc,
		ServiceIndicator:       si,
		NetworkIndicator:       ni,
		MessagePriority:        mp,
		SignalingLinkSelection: sls,
	}
}

// Dial establishes a M3UA connection as a client.
//
// After successfully established the connection with peer, state-changing
// signals and heartbeats are automatically handled background in another goroutine.
func Dial(ctx context.Context, net string, laddr, raddr *sctp.SCTPAddr, cfg *Config) (*Conn, error) {
	var err error
	conn := &Conn{
		mu:          new(sync.Mutex),
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
		return nil, err
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

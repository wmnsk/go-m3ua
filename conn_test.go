// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/wmnsk/go-m3ua/messages/params"

	"github.com/ishidawataru/sctp"
)

func setupConn(ctx context.Context) (*Conn, *Conn, error) {
	var (
		srvConnChan = make(chan *Conn)
		errChan     = make(chan error)
	)

	srvCfg := NewServerConfig(
		&HeartbeatInfo{Enabled: false},
		0x22222222, // OriginatingPointCode
		0x11111111, // DestinationPointCode
		1,          // AspIdentifier
		params.TrafficModeLoadshare, // TrafficModeType
		0,                     // NetworkAppearance
		0,                     // CorrelationID
		[]uint32{1, 2},        // RoutingContexts
		params.ServiceIndSCCP, // ServiceIndicator
		0, // NetworkIndicator
		0, // MessagePriority
		1, // SignalingLinkSelection
	)
	// set nil on unnecessary parameters.
	srvCfg.AspIdentifier = nil
	srvCfg.CorrelationID = nil

	// setup SCTP peer on the specified IPs and Port.
	raddr, err := sctp.ResolveSCTPAddr("sctp", "127.0.0.2:2905")
	if err != nil {
		return nil, nil, err
	}

	listener, err := Listen("m3ua", raddr, srvCfg)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		srvConn, err := listener.Accept(ctx)
		if err != nil {
			errChan <- err
		}

		srvConnChan <- srvConn
	}()

	cliCfg := NewClientConfig(
		&HeartbeatInfo{Enabled: false},
		0x11111111, // OriginatingPointCode
		0x22222222, // DestinationPointCode
		1,          // AspIdentifier
		params.TrafficModeLoadshare, // TrafficModeType
		0,                     // NetworkAppearance
		0,                     // CorrelationID
		[]uint32{1, 2},        // RoutingContexts
		params.ServiceIndSCCP, // ServiceIndicator
		0, // NetworkIndicator
		0, // MessagePriority
		1, // SignalingLinkSelection
	)
	// set nil on unnecessary parameters.
	cliCfg.CorrelationID = nil

	laddr, err := sctp.ResolveSCTPAddr("sctp", "127.0.0.1:2905")
	if err != nil {
		return nil, nil, err
	}

	cliConn, err := Dial(ctx, "m3ua", laddr, raddr, cliCfg)
	if err != nil {
		return nil, nil, err
	}

	select {
	case srvConn := <-srvConnChan:
		return cliConn, srvConn, nil
	case err := <-errChan:
		return nil, nil, err
	case <-time.After(10 * time.Second):
		return nil, nil, errors.New("timed out")
	}
}

func TestReadWrite(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cliConn, srvConn, err := setupConn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cliConn.Close()
		srvConn.Close()
	}()

	msg := []byte{0xde, 0xad, 0xbe, 0xef}
	buf := make([]byte, 1024)

	t.Run("client-write", func(t *testing.T) {
		if _, err := cliConn.Write(msg); err != nil {
			t.Fatal(err)
		}

		n, err := srvConn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(buf[:n], msg); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("server-write", func(t *testing.T) {
		if _, err := srvConn.Write(msg); err != nil {
			t.Fatal(err)
		}

		n, err := cliConn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(buf[:n], msg); diff != "" {
			t.Error(diff)
		}
	})
}

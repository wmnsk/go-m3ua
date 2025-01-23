// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Command m3ua-client works as M3UA client.
*/
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/dmisol/go-m3ua/messages/params"

	"github.com/dmisol/go-m3ua"
	"github.com/ishidawataru/sctp"
)

func main() {
	var (
		addr  = flag.String("addr", "127.0.0.1:2905", "Remote IP and Port to connect to.")
		data  = []byte{0x9, 0x0, 0x3, 0x7, 0xb, 0x4, 0x43, 0xb9, 0x0, 0xfe, 0x4, 0x43, 0x1, 0x0, 0xfe, 0x6, 0x0, 0x4, 0x30, 0x4, 0x1, 0x20}
		hbInt = flag.Duration("hb-interval", 0, "Interval for M3UA BEAT. Put 0 to disable")
	)
	flag.Parse()

	// create *Config to be used in M3UA connection
	config := m3ua.NewConfig(
		0x123,                 // OriginatingPointCode
		0x22222222,            // DestinationPointCode
		params.ServiceIndSCCP, // ServiceIndicator
		0,                     // NetworkIndicator
		0,                     // MessagePriority
		1,                     // SignalingLinkSelection
	)
	config. // set parameters to use
		EnableHeartbeat(*hbInt, 10*time.Second).
		SetAspIdentifier(1).
		//SetTrafficModeType(params.TrafficModeLoadshare).
		SetNetworkAppearance(0).
		//SetRoutingContexts(5951). // todo: fetch from RegResp
		SetRoutingKey(0x123)

	/* or, you can define config in the following way.
	config := m3ua.NewClientConfig(
		&m3ua.HeartbeatInfo{
			Enabled:  true,
			Interval: *hbInt,
			Timer:    time.Duration(10 * time.Second),
		},
		0x11111111,                  // OriginatingPointCode
		0x22222222,                  // DestinationPointCode
		1,                           // AspIdentifier
		params.TrafficModeLoadshare, // TrafficModeType
		0,                           // NetworkAppearance
		0,                           // CorrelationID
		[]uint32{1, 2},              // RoutingContexts
		params.ServiceIndSCCP,       // ServiceIndicator
		0,                           // NetworkIndicator
		0,                           // MessagePriority
		1,                           // SignalingLinkSelection
	)
	// set nil on unnecessary parameters.
	config.CorrelationID = nil
	*/

	// setup SCTP peer on the specified IPs and Port.
	raddr, err := sctp.ResolveSCTPAddr("sctp", *addr)
	if err != nil {
		log.Fatalf("Failed to resolve SCTP address: %s", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := m3ua.Dial(ctx, "m3ua", nil, raddr, config, make(chan *m3ua.ServeEvent, 5))
	if err != nil {
		log.Fatalf("Failed to dial M3UA: %s", err)
	}
	defer conn.Close()

	// send data once in 3 seconds.
	for {
		if _, err := conn.Write(data); err != nil {
			log.Fatalf("Failed to write M3UA data: %s", err)
		}

		time.Sleep(10 * time.Second)
	}
}

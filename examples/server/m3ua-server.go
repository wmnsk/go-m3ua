// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Command m3ua-server works as M3UA server.
*/
package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/dmisol/go-m3ua/messages/params"

	"github.com/dmisol/go-m3ua"
	"github.com/ishidawataru/sctp"
)

var (
	mu        sync.Mutex
	conns     map[int]*m3ua.Conn
	serveChan chan *m3ua.ServeEvent
)

func serve(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down, ctx")
			return
		case ev := <-serveChan:
			if ev.Err != nil {
				log.Println("conn is closed:", ev.Err)

				go func(id int) {
					mu.Lock()
					defer mu.Unlock()
					delete(conns, id)
				}(ev.Id)
			} else {
				log.Printf("%d(%X) -> Read: %x", ev.Id, ev.PD.OriginatingPointCode, ev.PD.Data)
			}
		}
	}
}

func main() {
	var (
		addr    = flag.String("addr", "127.0.0.1:2905", "Source IP and Port listen.")
		hbInt   = flag.Duration("hb-interval", 0, "Interval for M3UA BEAT. Put 0 to disable")
		hbTimer = flag.Duration("hb-timer", time.Duration(5*time.Second), "Expiration timer for M3UA BEAT. Ignored when hb-interval is 0")
	)
	flag.Parse()

	// create *Config to be used in M3UA connection
	config := m3ua.NewServerConfig(
		&m3ua.HeartbeatInfo{
			Enabled:  true,
			Interval: *hbInt,
			Timer:    *hbTimer,
		},
		0x22222222,                  // OriginatingPointCode
		0x11111111,                  // DestinationPointCode
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
	config.AspIdentifier = nil
	config.CorrelationID = nil

	// setup SCTP listener on the specified IPs and Port.
	laddr, err := sctp.ResolveSCTPAddr("sctp", *addr)
	if err != nil {
		log.Fatalf("Failed to resolve SCTP address: %s", err)
	}

	listener, err := m3ua.Listen("m3ua", laddr, config)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	log.Printf("Waiting for connection on: %s", listener.Addr())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	id := 1
	serveChan = make(chan *m3ua.ServeEvent, 10)
	conns = make(map[int]*m3ua.Conn)
	go serve(ctx)

	for {
		conn, err := listener.Accept(ctx, serveChan, id)

		if err != nil {
			log.Fatalf("Failed to accept M3UA: %s", err)
			continue
		}

		mu.Lock()
		conns[id] = conn
		mu.Unlock()
		id++

		log.Printf("Connected with: %s", conn.RemoteAddr())
	}
}

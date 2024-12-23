// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"time"

	"github.com/dmisol/go-m3ua/messages/params"
)

// HeartbeatInfo is a set of information for M3UA BEAT.
type HeartbeatInfo struct {
	Enabled  bool
	Interval time.Duration
	Timer    time.Duration
	Data     []byte
}

// NewHeartbeatInfo creates a new HeartbeatInfo.
func NewHeartbeatInfo(interval, timer time.Duration, data []byte) *HeartbeatInfo {
	return &HeartbeatInfo{
		Enabled: true, Interval: interval, Timer: timer, Data: data,
	}
}

// Config is a configration that defines a M3UA server.
type Config struct {
	*HeartbeatInfo
	AspIdentifier          *params.Param
	TrafficModeType        *params.Param
	NetworkAppearance      *params.Param
	RoutingContexts        *params.Param
	CorrelationID          *params.Param
	SelfSPC                uint32
	DefaultDPC             uint32
	ServiceIndicator       uint8
	NetworkIndicator       uint8
	MessagePriority        uint8
	SignalingLinkSelection uint8
}

// NewConfig creates a new Config.
//
// To set additional parameters, use constructors in param package or
// setters defined in this package. Note that the params left nil won't
// appear in the packets but the initialized params will, with zero
// values.
func NewConfig(own, remote uint32, si, ni, mp, sls uint8) *Config {
	return &Config{
		SelfSPC:                own,
		DefaultDPC:             remote,
		ServiceIndicator:       si,
		NetworkIndicator:       ni,
		MessagePriority:        mp,
		SignalingLinkSelection: sls,
	}
}

// EnableHeartbeat enables M3UA BEAT with interval and expiration timer
// given.
//
// The data is hard-coded by default. Manipulate the exported field
// Config.HeartbeatInfo.Data to customize it such as including current
// time to identify the BEAT and BEAT ACK pair.
func (c *Config) EnableHeartbeat(interval, timer time.Duration) *Config {
	c.HeartbeatInfo = NewHeartbeatInfo(
		interval, timer,
		[]byte("Hi, this is a BEAT from go-m3ua. Are you alive?"),
	)
	return c
}

// SetAspIdentifier sets AspIdentifier in Config.
func (c *Config) SetAspIdentifier(id uint32) *Config {
	c.AspIdentifier = params.NewAspIdentifier(id)
	return c
}

// SetTrafficModeType sets TrafficModeType in Config.
func (c *Config) SetTrafficModeType(tmType uint32) *Config {
	c.TrafficModeType = params.NewTrafficModeType(tmType)
	return c
}

// SetNetworkAppearance sets NetworkAppearance in Config.
func (c *Config) SetNetworkAppearance(nwApr uint32) *Config {
	c.NetworkAppearance = params.NewNetworkAppearance(nwApr)
	return c
}

// SetRoutingContexts sets RoutingContexts in Config.
func (c *Config) SetRoutingContexts(rtCtxs ...uint32) *Config {
	c.RoutingContexts = params.NewRoutingContext(rtCtxs...)
	return c
}

// SetCorrelationID sets CorrelationID in Config.
func (c *Config) SetCorrelationID(id uint32) *Config {
	c.CorrelationID = params.NewCorrelationID(id)
	return c
}

// NewClientConfig creates a new Config for Client.
//
// The optional parameters that is not required (like CorrelationID)
// can be omitted by setting it to nil after created *Config.
func NewClientConfig(hbInfo *HeartbeatInfo, own, remote, aspID, tmt, nwApr, corrID uint32, rtCtxs []uint32, si, ni, mp, sls uint8) *Config {
	return &Config{
		HeartbeatInfo:          hbInfo,
		AspIdentifier:          params.NewAspIdentifier(aspID),
		TrafficModeType:        params.NewTrafficModeType(tmt),
		NetworkAppearance:      params.NewNetworkAppearance(nwApr),
		RoutingContexts:        params.NewRoutingContext(rtCtxs...),
		CorrelationID:          params.NewCorrelationID(corrID),
		SelfSPC:                own,
		DefaultDPC:             remote,
		ServiceIndicator:       si,
		NetworkIndicator:       ni,
		MessagePriority:        mp,
		SignalingLinkSelection: sls,
	}
}

// NewServerConfig creates a new Config for Server.
//
// The optional parameters that is not required (like CorrelationID)
// can be omitted by setting it to nil after created *Config.
func NewServerConfig(hbInfo *HeartbeatInfo, own, remote, aspID, tmt, nwApr, corrID uint32, rtCtxs []uint32, si, ni, mp, sls uint8) *Config {
	return &Config{
		HeartbeatInfo:          hbInfo,
		AspIdentifier:          params.NewAspIdentifier(aspID),
		TrafficModeType:        params.NewTrafficModeType(tmt),
		NetworkAppearance:      params.NewNetworkAppearance(nwApr),
		RoutingContexts:        params.NewRoutingContext(rtCtxs...),
		CorrelationID:          params.NewCorrelationID(corrID),
		SelfSPC:                own,
		DefaultDPC:             remote,
		ServiceIndicator:       si,
		NetworkIndicator:       ni,
		MessagePriority:        mp,
		SignalingLinkSelection: sls,
	}
}

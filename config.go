// Copyright 2018-2024 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"time"

	"github.com/wmnsk/go-m3ua/messages/params"
)

// HeartbeatInfo is a set of information for M3UA BEAT.
type HeartbeatInfo struct {
	Enabled  bool
	Interval time.Duration
	Timer    time.Duration
	Data     []byte
}

// SctpSackInfo is a set of information for SCTP SACK timer configuration.
//
// SackDelay sack_delay: This parameter contains the number of milliseconds the
// user is requesting that the delayed SACK timer be set to.  Note
// that this value is defined in [RFC4960] to be between 200 and 500
// milliseconds.
//
// SackFrequency sack_freq: This parameter contains the number of packets that must
// be received before a SACK is sent without waiting for the delay
// timer to expire.  The default value is 2; setting this value to 1
// will disable the delayed SACK algorithm.
type SctpSackInfo struct {
	Enabled       bool
	SackDelay     uint32
	SackFrequency uint32
}

// NewHeartbeatInfo creates a new HeartbeatInfo.
func NewHeartbeatInfo(interval, timer time.Duration, data []byte) *HeartbeatInfo {
	return &HeartbeatInfo{
		Enabled: true, Interval: interval, Timer: timer, Data: data,
	}
}

// Config is a configuration that defines a M3UA server.
type Config struct {
	*HeartbeatInfo
	*SctpSackInfo
	AspIdentifier          *params.Param
	TrafficModeType        *params.Param
	NetworkAppearance      *params.Param
	RoutingContexts        *params.Param
	CorrelationID          *params.Param
	OriginatingPointCode   uint32
	DestinationPointCode   uint32
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
func NewConfig(opc, dpc uint32, si, ni, mp, sls uint8) *Config {
	return &Config{
		OriginatingPointCode:   opc,
		DestinationPointCode:   dpc,
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

// SetSackConfig sets the SCTP SACK timer configuration.
//
// sackDelay is the number of milliseconds for the delayed SACK timer
// (per RFC4960, should be between 200 and 500 ms).
//
// sackFrequency is the number of packets to receive before sending a SACK
// without waiting for the delay timer. Setting to 1 disables the delayed
// SACK algorithm.
//
// Note: sackDelay=0, sackFrequency=1 (disables delayed SACK)
func (c *Config) SetSackConfig(sackDelay, sackFrequency uint32) *Config {
	c.SctpSackInfo = &SctpSackInfo{
		Enabled:       true,
		SackDelay:     sackDelay,
		SackFrequency: sackFrequency,
	}
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

// NewServerConfig creates a new Config for Server.
//
// The optional parameters that is not required (like CorrelationID)
// can be omitted by setting it to nil after created *Config.
func NewServerConfig(hbInfo *HeartbeatInfo, opc, dpc, aspID, tmt, nwApr, corrID uint32, rtCtxs []uint32, si, ni, mp, sls uint8) *Config {
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

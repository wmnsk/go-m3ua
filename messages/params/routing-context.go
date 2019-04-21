// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// NewRoutingContext creates the RoutingContext Parameter.
// Multiple number of RoutingContext will be accepted.
// Note that this returns *Param, as no specific structure in this parameter.
func NewRoutingContext(rtCxts ...uint32) *Param {
	return newMultiUint32ValParam(RoutingContext, rtCxts...)
}

// RoutingContext returns single RoutingContext from Param.
func (p *Param) RoutingContext() uint32 {
	if p.Tag != RoutingContext {
		return 0
	}

	return p.RoutingContexts()[0]
}

// RoutingContexts returns multiple RoutingContexts from Param.
func (p *Param) RoutingContexts() []uint32 {
	if p.Tag != RoutingContext {
		return nil
	}

	return p.decodeMultiUint32ValData()
}

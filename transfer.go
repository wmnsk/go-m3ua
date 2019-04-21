// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"

	"github.com/wmnsk/go-m3ua/messages"
)

func (c *Conn) handleData(ctx context.Context, data *messages.Data) {
	c.mu.Lock()
	if c.state != StateAspActive {
		c.errChan <- NewErrUnexpectedMessage(data)
		return
	}
	c.mu.Unlock()

	pd, err := data.ProtocolData.ProtocolData()
	if err != nil {
		c.errChan <- ErrFailedToPeelOff
		return
	}

	if c.cfg.OriginatingPointCode != pd.DestinationPointCode {
		c.errChan <- NewErrUnexpectedMessage(data)
		return
	}

	select {
	case c.dataChan <- pd:
		return
	case <-ctx.Done():
		return
	}
}

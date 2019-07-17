// Copyright 2018-2019 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"errors"

	"github.com/wmnsk/go-m3ua/messages"
)

func (c *Conn) handleData(ctx context.Context, data *messages.Data) {
	err := func() error {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.state != StateAspActive {
			c.errChan <- NewErrUnexpectedMessage(data)
			return errors.New(data.String())
		}
		return nil
	}()
	if err != nil {
		// it has already emitted the error into the errChan, early exit
		return
	}

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

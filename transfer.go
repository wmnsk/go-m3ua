// Copyright 2018-2024 go-m3ua authors. All rights reserved.
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
		if c.State() != StateAspActive {
			c.errChan <- NewUnexpectedMessageError(data)
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

	select {
	case c.dataChan <- pd:
		return
	case <-ctx.Done():
		return
	}
}

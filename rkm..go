package m3ua

import "github.com/dmisol/go-m3ua/messages"

func (c *Conn) initiateREQREQ() error {
	if _, err := c.WriteSignal(
		messages.NewRegReq(c.cfg.RoutingKey),
	); err != nil {
		return err
	}

	return nil
}

package m3ua

import (
	"errors"

	"github.com/dmisol/go-m3ua/messages/params"
)

var ErrConnectionClosed = errors.New("connection closed")

type ServeEvent struct {
	PD  *params.ProtocolDataPayload
	Id  int
	Err error
}

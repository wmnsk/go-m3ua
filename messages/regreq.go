package messages

import (
	"fmt"
	"log"

	"github.com/dmisol/go-m3ua/messages/params"
)

type RegReq struct {
	*Header
	RoutingKey *params.Param
}

// todo: build RoutingKey from SPC
// NewRegReq creates a new RegReq.
func NewRegReq(rk *params.Param) *RegReq {
	a := &RegReq{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassRKM,
			Type:     MsgTypeRegistrationRequest,
		},
		RoutingKey: rk,
	}
	a.SetLength()

	return a
}

// MarshalBinary returns the byte sequence generated from a RegReq.
func (a *RegReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, a.MarshalLen())
	if err := a.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (a *RegReq) MarshalTo(b []byte) error {
	if len(b) < a.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	a.Header.Payload = make([]byte, a.MarshalLen()-8)

	var offset = 0

	if param := a.RoutingKey; param != nil {
		if err := param.MarshalTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
	}

	return a.Header.MarshalTo(b)
}

// ParseRegReq decodes given byte sequence as a RegReq.
func ParseRegReq(b []byte) (*RegReq, error) {
	a := &RegReq{}
	if err := a.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return a, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (a *RegReq) UnmarshalBinary(b []byte) error {
	var err error
	a.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}

	prs, err := params.ParseMultiParams(a.Header.Payload)
	if err != nil {
		return err
	}
	for _, pr := range prs {
		switch pr.Tag {
		case params.RoutingKey:
			a.RoutingKey = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (a *RegReq) SetLength() {
	if param := a.RoutingKey; param != nil {
		param.SetLength()
	}

	a.Header.Length = uint32(a.MarshalLen())
}

// MarshalLen returns the serial length of RegReq.
func (a *RegReq) MarshalLen() int {
	l := 8
	if param := a.RoutingKey; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// String returns the RegReq values in human readable format.
func (a *RegReq) String() string {
	return fmt.Sprintf("{Header: %s, RoutingKey: %s}",
		a.Header.String(),
		a.RoutingKey.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *RegReq) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *RegReq) MessageType() uint8 {
	return MsgTypeRegistrationRequest
}

// MessageClass returns the message class in int.
func (a *RegReq) MessageClass() uint8 {
	return MsgClassRKM
}

// MessageClassName returns the name of message class.
func (a *RegReq) MessageClassName() string {
	return MsgClassNameRKM
}

// MessageTypeName returns the name of message type.
func (a *RegReq) MessageTypeName() string {
	return "REG REQ"
}

// Serialize returns the byte sequence generated from a RegReq.
//
// DEPRECATED: use MarshalBinary instead.
func (a *RegReq) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return a.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (a *RegReq) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return a.MarshalTo(b)
}

// DecodeRegReq decodes given byte sequence as a RegReq.
//
// DEPRECATED: use ParseRegReq instead.
func DecodeRegReq(b []byte) (*RegReq, error) {
	log.Println("DEPRECATED: use ParseRegReq instead")
	return ParseRegReq(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (a *RegReq) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return a.UnmarshalBinary(b)
}

// Len returns the serial length of RegReq.
//
// DEPRECATED: use MarshalLen instead.
func (a *RegReq) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return a.MarshalLen()
}

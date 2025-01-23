package messages

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/dmisol/go-m3ua/messages/params"
)

const tagRoutingCtx = uint16(6)

type RegRsp struct {
	*Header
	RegistrationResult *params.Param
}

// todo: build RegistrationResult from SPC
// NewRegRsp creates a new RegRsp.
func NewRegRsp(rr *params.Param) *RegRsp {
	a := &RegRsp{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassRKM,
			Type:     MsgTypeRegistrationResponse,
		},
		RegistrationResult: rr,
	}
	a.SetLength()

	return a
}

func (a *RegRsp) FetchRC() (uint32, error) {
	if a.RegistrationResult == nil {
		return 0, fmt.Errorf("no RegistrationResult")
	}
	rr := a.RegistrationResult
	if rr.Length < 8 && len(rr.Data) != int(rr.Length)-4 {
		return 0, fmt.Errorf("wrong RR length")
	}
	l := rr.Length - 4
	offs := 0
	fmt.Println(rr.Length, rr.Data)
	for l > 0 {
		fmt.Println("offs:", offs)
		pt := binary.BigEndian.Uint16(rr.Data[offs : offs+2])
		pl := binary.BigEndian.Uint16(rr.Data[offs+2 : offs+4])
		if pt == tagRoutingCtx {
			if pl != 8 {
				return 0, fmt.Errorf("wrong RC length")
			}
			return binary.BigEndian.Uint32(rr.Data[offs+4 : offs+8]), nil
		}
		offs += int(pl)
		l -= pl
	}
	return 0, fmt.Errorf("RC not found")
}

// MarshalBinary returns the byte sequence generated from a RegRsp.
func (a *RegRsp) MarshalBinary() ([]byte, error) {
	b := make([]byte, a.MarshalLen())
	if err := a.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (a *RegRsp) MarshalTo(b []byte) error {
	if len(b) < a.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	a.Header.Payload = make([]byte, a.MarshalLen()-8)

	var offset = 0

	if param := a.RegistrationResult; param != nil {
		if err := param.MarshalTo(a.Header.Payload[offset:]); err != nil {
			return err
		}
	}

	return a.Header.MarshalTo(b)
}

// ParseRegRsp decodes given byte sequence as a RegRsp.
func ParseRegRsp(b []byte) (*RegRsp, error) {
	a := &RegRsp{}
	if err := a.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return a, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (a *RegRsp) UnmarshalBinary(b []byte) error {
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
		case params.RegistrationResult:
			a.RegistrationResult = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (a *RegRsp) SetLength() {
	if param := a.RegistrationResult; param != nil {
		param.SetLength()
	}

	a.Header.Length = uint32(a.MarshalLen())
}

// MarshalLen returns the serial length of RegRsp.
func (a *RegRsp) MarshalLen() int {
	l := 8
	if param := a.RegistrationResult; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// String returns the RegRsp values in human readable format.
func (a *RegRsp) String() string {
	return fmt.Sprintf("{Header: %s, RegistrationResult: %s}",
		a.Header.String(),
		a.RegistrationResult.String(),
	)
}

// Version returns the version of M3UA in int.
func (a *RegRsp) Version() uint8 {
	return a.Header.Version
}

// MessageType returns the message type in int.
func (a *RegRsp) MessageType() uint8 {
	return MsgTypeRegistrationResponse
}

// MessageClass returns the message class in int.
func (a *RegRsp) MessageClass() uint8 {
	return MsgClassRKM
}

// MessageClassName returns the name of message class.
func (a *RegRsp) MessageClassName() string {
	return MsgClassNameRKM
}

// MessageTypeName returns the name of message type.
func (a *RegRsp) MessageTypeName() string {
	return "REG REQ"
}

// Serialize returns the byte sequence generated from a RegRsp.
//
// DEPRECATED: use MarshalBinary instead.
func (a *RegRsp) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return a.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (a *RegRsp) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return a.MarshalTo(b)
}

// DecodeRegRsp decodes given byte sequence as a RegRsp.
//
// DEPRECATED: use ParseRegRsp instead.
func DecodeRegRsp(b []byte) (*RegRsp, error) {
	log.Println("DEPRECATED: use ParseRegRsp instead")
	return ParseRegRsp(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (a *RegRsp) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return a.UnmarshalBinary(b)
}

// Len returns the serial length of RegRsp.
//
// DEPRECATED: use MarshalLen instead.
func (a *RegRsp) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return a.MarshalLen()
}

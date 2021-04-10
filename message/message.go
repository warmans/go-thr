package message

import (
	"github.com/warmans/go-thr/encoding"
	"github.com/warmans/go-thr/util"
)

// msg type one seems to be used predominantly for "getter" type commands while
// msg type two seems to be used for "setter" commands. However this doesn't seem to be the actual
// system because sometimes things I would expect to use 00 use 01.
type MessageType byte

const TypeOne MessageType = 0x00
const TypeTwo MessageType = 0x01

type MessageSet []Encodable

type Encodable interface {
	Bytes(seqNum uint32) []byte
}

type RawMessage struct {
	Data []byte
}

func (c *RawMessage) Bytes(seqNum uint32) []byte {
	return c.Data
}

type THRMessage struct {
	Type        MessageType
	PayloadType [3]byte
	Payload     []byte
}

func (c *THRMessage) Bytes(seqNum uint32) []byte {
	msg := encoding.Message{
		ManufacturerCode: yamahaManufacturerCode(),
		Preamble:         preamble(),
		MessageType:      byte(c.Type),
		SequenceNum:      sequenceNumber(seqNum),
		PayloadType:      c.PayloadType,
		Payload:          c.Payload,
	}
	return msg.Encode()
}

func yamahaManufacturerCode() [3]byte {
	// manufacturer ID 0x00 indicates a 3 byte ID
	return [3]byte{0x00, 0x01, 0x0C}
}

// mystery stuff that gets send with every command
func preamble() [3]byte {
	return [3]byte{0x24, 0x00, 0x4d}
}

func sequenceNumber(seqNum uint32) byte {
	// sequence number should only be 1 byte so it needs to roll over at 127
	if seqNum > 127 {
		seqNum = seqNum - (127 * (seqNum / 127))
	}
	return util.SingleByteInt(int8(seqNum))
}

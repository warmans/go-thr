package encoding

import (
	"bytes"
	"encoding/hex"
	"html/template"
)

// standard sysex start and end bytes
const msgStart = 0xf0
const msgEnd = 0xf7

// next will take a buffer of bytes and return the next complete sysex message
// as well as the remaining buffer data. Incomplete messages are not really handled
// because without knowing the protocol, it's difficult not to get stuck in a loop.
func Next(buff []byte) (*Message, []byte) {

	msg := &Message{}

	// scan through the buffer looking for the start of the message
	var i int
	for i = 0; i < len(buff); i++ {
		if buff[i] != msgStart {
			continue
		} else {
			break
		}
	}
	// discard all read bytes from the buffer before start byte
	if len(buff) > i {
		buff = buff[i+1:]
	} else {
		return nil, []byte{}
	}

	// there should be 11 bytes of headers
	if len(buff) < 11 {
		return nil, buff
	}
	copy(msg.ManufacturerCode[:], buff[:3])
	copy(msg.Preamble[:], buff[3:6])
	msg.MessageType = buff[6]
	msg.SequenceNum = buff[7]
	copy(msg.PayloadType[:], buff[8:11])

	// advance to payload
	buff = buff[11:]

	payloadSize := payloadSize(msg.PayloadType)
	if len(buff) < payloadSize {
		// not enough data in buffer
		return nil, buff
	}
	msg.Payload = buff[:payloadSize]

	buff = buff[payloadSize:]
	if len(buff) == 0 || buff[0] != msgEnd {
		// end byte was not at expected position
		return nil, buff
	}
	buff = buff[1:]

	return msg, buff
}

func payloadSize(payloadType [3]byte) int {
	var payloadSize int
	switch payloadType {
	case [3]byte{0x00, 0x00, 0x03}:
		// 0x03 denotes a short message (12 bytes) + msgEnd
		payloadSize = 12
	case [3]byte{0x00, 0x00, 0x07}, [3]byte{0x00, 0x00, 0x08}, [3]byte{0x00, 0x00, 0x0b}:
		// consider all others (e.g. 0x07, 0x08, 0x0b) 16 bytes
		payloadSize = 16
	case [3]byte{0x00, 0x01, 0x07}:
		payloadSize = 32
	default:
		// todo what are the possible values?
		payloadSize = 16
	}
	return payloadSize
}

type Message struct {
	ManufacturerCode [3]byte
	Preamble         [3]byte
	MessageType      byte
	SequenceNum      byte
	PayloadType      [3]byte
	Payload          []byte
}

func (m *Message) Encode() []byte {
	// standard sysex start
	buff := []byte{msgStart}
	// The manufacturer's code uses the extended 3 byte format.
	buff = append(buff, m.ManufacturerCode[:]...)
	// All messages have this same preamble - is it a device code or something?
	buff = append(buff, m.Preamble[:]...)
	// Commands seems to be prefixed with 00 or 01. see msgType above
	buff = append(buff, m.MessageType)
	// There is a 1 byte sequence that rolls over when at the maximum value
	buff = append(buff, m.SequenceNum)
	// the payload seems to have a type that dictates its length and probably other stuff.
	buff = append(buff, m.PayloadType[:]...)
	// The payload is arbitrary bytes
	buff = append(buff, m.Payload...)
	// standard sysex end
	buff = append(buff, msgEnd)
	return buff
}

func (m *Message) Hex() string {
	return hex.EncodeToString(m.Encode())
}

func (m *Message) Print(tmpl string) string {
	parsed, err := template.New("printer").Parse(tmpl)
	if err != nil {
		return "invalid template: " + err.Error()
	}
	buff := &bytes.Buffer{}
	if err := parsed.Execute(buff, m); err != nil {
		return "invalid template: " + err.Error()
	}
	return buff.String()
}

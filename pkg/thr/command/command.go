package command

import (
	"encoding/hex"

	"github.com/rakyll/portmidi"
	"github.com/warmans/go-thr/pkg/thr/encoding"
	"github.com/warmans/go-thr/pkg/thr/util"
	"go.uber.org/zap"
)

// msg type one seems to be used predominantly for "getter" type commands while
// msg type two seems to be used for "setter" commands. However this doesn't seem to be the actual
// system because sometimes things I would expect to use 00 use 01.
type commandType byte

const TypeOne commandType = 0x00
const TypeTwo commandType = 0x01

type CommandSet []Command

type Command interface {
	Bytes(seqNum uint32) []byte
}

type RawCommmand struct {
	Data []byte
}

func (c *RawCommmand) Bytes(seqNum uint32) []byte {
	return c.Data
}

type THRCommand struct {
	Type        commandType
	PayloadType byte
	Payload     []byte
}

func (c *THRCommand) Bytes(seqNum uint32) []byte {
	msg := encoding.Message{
		ManufacturerCode: yamahaManufacturerCode(),
		Preamble:         preamble(),
		MessageType:      byte(c.Type),
		SequenceNum:      sequenceNumber(seqNum),
		Reserved1:        [2]byte{0x00, 0x00},
		PayloadType:      c.PayloadType,
		Payload:          c.Payload,
	}
	return msg.Encode()
}

func NewSession(out *portmidi.Stream, logger *zap.Logger) *Session {
	return &Session{out: out, logger: logger}
}

// Session manages sequence numbers
type Session struct {
	sequenceNum uint32
	out         *portmidi.Stream
	logger      *zap.Logger
}

func (s *Session) Send(cmds CommandSet) error {
	for _, cmd := range cmds {
		data := cmd.Bytes(s.sequenceNum)
		if err := s.out.WriteSysExBytes(portmidi.Time(), data); err != nil {
			return err
		}
		if s.logger != nil {
			s.logger.Debug("sent", zap.String("data", hex.EncodeToString(data)))
		}
		s.sequenceNum++

	}
	return nil
}

func yamahaManufacturerCode() [3]byte {
	// manufacturer ID 0x00 indicates a 3 byte ID
	return [3]byte{0x00, 0x01, 0x0C}
}

// mystery stuff that gets send with every command
func preamble() [3]byte {
	return [3]byte{0x22, 0x02, 0x4d}
}

func sequenceNumber(seqNum uint32) byte {
	// sequence number should only be 1 byte so it needs to roll over at 127
	if seqNum > 127 {
		seqNum = seqNum - (127 * (seqNum / 127))
	}
	return util.SingleByteInt(int8(seqNum))
}

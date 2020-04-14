package encoding

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rakyll/portmidi"
)

// standard sysex start and end bytes
const msgStart = 0xf0
const msgEnd = 0xf7

func NewListener(in *portmidi.Stream) *Listener {
	return &Listener{
		in:     in,
		c:      make(chan *Message, 100),
		closed: make(chan struct{}, 0),
	}
}

type Listener struct {
	in *portmidi.Stream
	c  chan *Message

	close  bool
	closed chan struct{}
}

func (l *Listener) Close() {
	l.close = true
	<-l.closed
}

func (l *Listener) shutdown() {
	close(l.c)
	l.closed <- struct{}{}
}

func (l *Listener) Data() <-chan *Message {
	return l.c
}

func (l *Listener) Listen() error {
	buff := []byte{}
	for {
		time.Sleep(10 * time.Millisecond)
		if l.close == true {
			l.shutdown()
			return nil
		}
		b, err := l.in.ReadSysExBytes(1024)
		if err != nil {
			l.shutdown()
			return err
		}
		if len(b) > 0 {
			buff = append(buff, b...)
		}

		// read as many messages from the buffer as possible
		for {
			var msg *Message
			msg, buff = Next(buff)
			if msg != nil {
				l.c <- msg
				continue
			}
			break
		}
	}
}

// next will take a buffer of bytes and return the next complete sysex message, as well as the
// remaining buffer that was not read.
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

	// the manufacturer code should be 3 bytes
	if len(buff) > 3 {
		// check for incorrect manufacturers code
		if buff[0] != 0x00 && buff[1] != 0x01 && buff[2] != 0x0c {
			return nil, buff
		}
		copy(msg.ManufacturerCode[:], buff[:3])
		buff = buff[3:]
	}

	// there should be 8 bytes of data leading up to the payload. The last byte
	// defines how big the payload will be.
	fmt.Println(hex.EncodeToString(buff))
	if len(buff) > 8 {
		copy(msg.Preamble[:], buff[:3])
		msg.MessageType = buff[3]
		msg.SequenceNum = buff[4]
		copy(msg.Reserved1[:], buff[5:6])
		msg.PayloadType = buff[7]

		buff = buff[8:]
	}

	switch msg.PayloadType {
	// 0x03 denotes a short message (12 bytes) + msgEnd
	case 0x03:
		if len(buff) < 12 {
			// not enough data in buffer
			return nil, buff
		}
		msg.Payload = buff[:12]
		buff = buff[12:]

	// consider all others (e.g. 0x07, 0x08, 0x0b) 16 bytes
	default:
		if len(buff) < 16 {
			// not enough data in buffer
			return nil, buff
		}
		msg.Payload = buff[:16]
		buff = buff[16:]
	}

	if len(buff) == 0 || buff[0] != msgEnd {
		// end byte was not at expected position
		return nil, buff
	}
	buff = buff[1:]

	return msg, buff
}

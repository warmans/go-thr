package amp

import (
	"fmt"
	"github.com/rakyll/portmidi"
	"go.uber.org/zap"
	"regexp"
	"time"
)

var ampName = regexp.MustCompile("THR[1-3]0II.+")

func GetThrInput(log *zap.Logger) (*portmidi.Stream, error) {
	return GetStream(log, ampName, StreamTypeInput)
}

func GetThrOutput(log *zap.Logger) (*portmidi.Stream, error) {
	return GetStream(log, ampName, StreamTypeOutput)
}

type StreamType string

const StreamTypeInput StreamType = "in"
const StreamTypeOutput StreamType = "out"

func GetStream(log *zap.Logger, nameMatcher *regexp.Regexp, streamType StreamType) (*portmidi.Stream, error) {
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		log.Debug(
			"found device",
			zap.Int("id", i),
			zap.String("name", info.Name),
			zap.String("interface", info.Interface),
			zap.Bool("input", info.IsInputAvailable),
			zap.Bool("output", info.IsOutputAvailable),
		)
		if nameMatcher.MatchString(info.Name) {
			if streamType == StreamTypeInput && info.IsInputAvailable {
				return portmidi.NewInputStream(portmidi.DeviceID(i), 1024)
			}
			if streamType == StreamTypeOutput && info.IsOutputAvailable {
				return portmidi.NewOutputStream(portmidi.DeviceID(i), 1024, 0)
			}
		}
	}

	return nil, fmt.Errorf("device not found")
}

func NewListener(in *portmidi.Stream) *Listener {
	return &Listener{
		in:     in,
		c:      make(chan []byte, 100),
		closed: make(chan struct{}, 0),
	}
}

type Listener struct {
	in *portmidi.Stream
	c  chan []byte

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

func (l *Listener) Data() <-chan []byte {
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
	}
}

// next will take a buffer of bytes and return the next complete sysex message, as well as the
// remaining buffer that was not read.
func Next(buff []byte) ([]byte, []byte) {

	msg := []byte{}

	// scan through the buffer looking for the start of the message
	for i := 0; i < len(buff); i++ {
		if buff[i] != msgStart {
			continue
		} else {
			msg = append(msg, buff[i])
			// discard all read bytes from the buffer
			if len(buff) > i {
				buff = buff[i+1:]
			} else {
				buff = []byte{}
			}
			break
		}
	}
	if len(msg) == 0 || len(buff) == 0 {
		return nil, []byte{}
	}
	// the manufacturer code should be 3 bytes
	if len(buff) > 3 {
		// check for incorrect manufacturers code
		if buff[0] != 0x00 && buff[1] != 0x01 && buff[2] != 0x0c {
			return nil, buff
		}
		msg = append(msg, buff[:3]...)
		buff = buff[3:]
	}
	// there should be 12 bytes of preamble
	if len(buff) > 9 {
		msg = append(msg, buff[:11]...)
		buff = buff[11:]
	}
	switch msg[len(msg)-1] {
	// 0x03 denotes a short message (12 bytes)  + msgEnd
	case 0x03:
		if len(buff) < 13 {
			// not enough data in buffer
			return nil, buff
		}
	msg = append(msg, buff[:12]...)
	buff = buff[:12]

	// consider all others (e.g. 0x07) 32 bytes + msgEnd
	default:
		if len(buff) < 33 {
			// not enough data in buffer
			return nil, buff
		}
		msg = append(msg, buff[:33]...)
		buff = buff[:33]
	}

	return msg, buff
}

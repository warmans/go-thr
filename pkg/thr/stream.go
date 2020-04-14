package thr

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/rakyll/portmidi"
	"github.com/warmans/go-thr/pkg/thr/encoding"
	"github.com/warmans/go-thr/pkg/thr/message"
	"go.uber.org/zap"
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
		c:      make(chan *encoding.Message, 100),
		closed: make(chan struct{}, 0),
	}
}

type Listener struct {
	in *portmidi.Stream
	c  chan *encoding.Message

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

func (l *Listener) Data() <-chan *encoding.Message {
	return l.c
}

func (l *Listener) Listen() error {
	defer l.shutdown()

	buff := []byte{}
	for {
		time.Sleep(10 * time.Millisecond)
		if l.close == true {
			return nil
		}
		b, err := l.in.ReadSysExBytes(1024)
		if err != nil {
			return err
		}
		if len(b) > 0 {
			//fmt.Printf("%s\n\n", hex.EncodeToString(b))
			buff = append(buff, b...)
		}

		// read as many messages from the buffer as possible
		if len(buff) > 0 {
			for {
				var msg *encoding.Message
				msg, buff = encoding.Next(buff)
				if msg != nil {
					l.c <- msg
					continue
				}
				break
			}
		}
	}
}

func NewSession(out *portmidi.Stream, logger *zap.Logger) *Session {
	return &Session{out: out, logger: logger}
}

// Session just tracks sequence numbers for messages.
type Session struct {
	sequenceNum uint32
	out         *portmidi.Stream
	logger      *zap.Logger
}

func (s *Session) Send(cmds message.MessageSet) error {
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

package thr

import (
	"fmt"
	"regexp"

	"github.com/rakyll/portmidi"
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
package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/warmans/go-thr/pkg/amp"
	"go.uber.org/zap"
	"log"
	"os"

	"github.com/rakyll/portmidi"
)

//todo: this only seems to work correctly if the app has been started at least once.
// there must be some kind of initialization missing.
func main() {
	os.Setenv("DEBUG", "true")

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	if err := portmidi.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer portmidi.Terminate()

	in, err := amp.GetThrInput(logger)
	if err != nil {
		logger.Fatal("failed to input find device", zap.Error(err))
	}
	defer in.Close()

	out, err := amp.GetThrOutput(logger)
	if err != nil {
		logger.Fatal("failed to output find device", zap.Error(err))
	}
	defer in.Close()

	session := amp.NewSession(out)

	fmt.Println("Switch channel 2...")
	if err := session.Send(amp.SelectPreset(2)); err != nil {
		logger.Fatal("failed to send command to enable events", zap.Error(err))
	}

	fmt.Println("Listening for responses...")
	for e := range in.Listen() {
		spew.Dump(e)
	}
}

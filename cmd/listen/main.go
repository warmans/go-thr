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

func main() {

	os.Setenv("DEBUG", "true")

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	if err := portmidi.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer portmidi.Terminate()

	fmt.Printf("there are %d devices\n", portmidi.CountDevices())
	fmt.Printf("default input deviceID is %v\n", portmidi.DefaultInputDeviceID())

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

	if err := session.Send(amp.EnableEvents); err != nil {
		logger.Fatal("failed to send command to enable events", zap.Error(err))
	}

	fmt.Println("Listening...")
	for e := range in.Listen() {
		spew.Dump(e)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/warmans/go-thr/pkg/thr"
	"github.com/warmans/go-thr/pkg/thr/command"
	"github.com/warmans/go-thr/pkg/thr/message"
	"go.uber.org/zap"

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

	in, err := thr.GetThrInput(logger)
	if err != nil {
		logger.Fatal("failed to input find device", zap.Error(err))
	}
	defer in.Close()

	out, err := thr.GetThrOutput(logger)
	if err != nil {
		logger.Fatal("failed to output find device", zap.Error(err))
	}
	defer in.Close()

	session := thr.NewSession(out, logger)

	if err := session.Send(thr.Init); err != nil {
		logger.Fatal("failed to init communication with device", zap.Error(err))
	}

	fmt.Println("Listening...")
	for e := range in.Listen() {
		spew.Dump(e)
	}
}

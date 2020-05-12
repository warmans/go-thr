package main

import (
	"fmt"
	"log"
	"os"

	"github.com/warmans/go-thr"
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

	listener := thr.NewListener(in)
	defer listener.Close()
	go func() {
		if err := listener.Listen(); err != nil {
			logger.Error("listen failed", zap.Error(err))
		}
	}()

	fmt.Println("Listening...")
	for msg := range listener.Data() {

		fmt.Println(msg.Print("PAYLOAD > {{.Payload}}"))
	}
}

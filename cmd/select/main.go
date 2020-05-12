package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/warmans/go-thr"
	"go.uber.org/zap"

	"github.com/rakyll/portmidi"
)

// This command will cycle through all the presets once then go into an endless listen loop for responses.
func main() {
	os.Setenv("DEBUG", "true")

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	if err := portmidi.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer portmidi.Terminate()

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

	listener := thr.NewListener(in)
	defer listener.Close()

	session := thr.NewSession(out, logger)

	if err := session.Send(thr.Init); err != nil {
		logger.Fatal("failed to init communication with device", zap.Error(err))
	}

	for i := int8(0); i < 5; i++ {
		fmt.Printf("Switch channel %d...\n", i)
		if err := session.Send(thr.SelectPreset(i)); err != nil {
			logger.Fatal("failed to send command to enable events", zap.Error(err))
		}
		time.Sleep(time.Second)
	}

	go func() {
		fmt.Println("Listening for responses...")
		for e := range listener.Data() {
			fmt.Println("DATA >", e.Hex())
		}
	}()

	if err := listener.Listen(); err != nil {
		logger.Fatal("listener failed", zap.Error(err))
	}
}

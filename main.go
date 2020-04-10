package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/rakyll/portmidi"
)

func main() {
	if err := portmidi.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer portmidi.Terminate()

	fmt.Printf("there are %d devices\n", portmidi.CountDevices())
	fmt.Printf("default input deviceID is %v\n", portmidi.DefaultInputDeviceID())

	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		fmt.Println(i, info.Name, info.Interface, info.IsInputAvailable, info.IsOutputAvailable)
	}

	in, err := portmidi.NewInputStream(portmidi.DeviceID(3), 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	fmt.Println("Listening...")
	for e := range in.Listen() {
		spew.Dump(e)
	}
}
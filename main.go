package main

import (
	"github.com/floreks/led-max7219-client/device"
	"github.com/spf13/pflag"

	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TempResponse struct {
	Temperature float32 `json:"temperature"`
}

const POLLING_TIME = 5 * time.Second

func pollTemperature(pollServer *string, tempChan chan float32) {
	tempResponse := new(TempResponse)

	for {
		r, err := http.Get(*pollServer)
		if err != nil {
			log.Printf("Error durring temperature polling: %s", err.Error())
			os.Exit(1)
		}

		err = json.NewDecoder(r.Body).Decode(tempResponse)
		if err != nil {
			log.Printf("Error during server response decoding: %s", err.Error())
			os.Exit(1)
		}

		tempChan <- tempResponse.Temperature
		time.Sleep(POLLING_TIME)
		r.Body.Close()
	}
}

func registerSigtermHandler(dev *device.Max7219) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		dev.ClearDisplay()
		os.Exit(1)
	}()
}

var (
	pollServer = pflag.String("poll-server", "", "Server address that should be used for polling data")
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if *pollServer == "" {
		log.Printf("Poll server not defined. Please specify server to poll data from.")
		os.Exit(1)
	}

	dev, err := device.NewMax7219()
	if err != nil {
		log.Printf("Error during device creation: %s", err.Error())
		return
	}

	defer dev.Close()

	tempChan := make(chan float32)
	go pollTemperature(pollServer, tempChan)

	registerSigtermHandler(dev)

	for temp := range tempChan {
		err := dev.DisplayTemperature(temp)
		if err != nil {
			log.Fatalf("Error during attempt to display temperature: %s", err.Error())
		}
	}
}

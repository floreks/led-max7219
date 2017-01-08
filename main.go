package main

import (
	"github.com/floreks/led-max7219-client/device"

	"encoding/json"
	"fmt"
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

const tempServerUrl = "http://192.168.0.100:30024/ds18b20"
const POLLING_TIME = 5 * time.Second

func pollTemperature(tempChan chan float32) {
	for {
		tempResponse := new(TempResponse)
		r, err := http.Get(tempServerUrl)
		if err != nil {
			log.Printf("Error durring temperature polling: %s", err.Error())
		}

		json.NewDecoder(r.Body).Decode(tempResponse)
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

func main() {
	dev, err := device.NewMax7219()
	if err != nil {
		fmt.Printf("Error during device creation: %s", err.Error())
		return
	}

	defer dev.Close()

	tempChan := make(chan float32)
	go pollTemperature(tempChan)

	registerSigtermHandler(dev)

	for temp := range tempChan {
		dev.DisplayTemperature(temp)
	}
}

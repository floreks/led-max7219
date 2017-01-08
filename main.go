package main

import (
	"fmt"

	"github.com/floreks/led-max7219-client/device"
	"github.com/fulr/spidev"
)

type Max7219Reg byte

var MAX7219_CLEAR_DISPLAY = []byte{0, 0, 0, 0, 0, 0, 0, 0}

const (
	MAX7219_REG_NOOP   Max7219Reg = 0
	MAX7219_REG_DIGIT0            = iota
	MAX7219_REG_DIGIT1
	MAX7219_REG_DIGIT2
	MAX7219_REG_DIGIT3
	MAX7219_REG_DIGIT4
	MAX7219_REG_DIGIT5
	MAX7219_REG_DIGIT6
	MAX7219_REG_DIGIT7
	MAX7219_REG_DECODEMODE
	MAX7219_REG_INTENSITY
	MAX7219_REG_SCANLIMIT
	MAX7219_REG_SHUTDOWN
	MAX7219_REG_DISPLAYTEST = 0x0F
	MAX7219_REG_LASTDIGIT   = MAX7219_REG_DIGIT7
)

const MAX7219_DIGIT_COUNT = MAX7219_REG_LASTDIGIT - MAX7219_REG_DIGIT0 + 1

func send(reg Max7219Reg, value byte, device *spidev.SPIDevice) error {
	buf := []byte{byte(reg), value}

	_, err := device.Xfer(buf)
	if err != nil {
		return err
	}

	return nil
}

func sendMulti(reg Max7219Reg, value, value2 byte, device *spidev.SPIDevice) error {
	buf := []byte{byte(reg), value, byte(reg), value2}

	_, err := device.Xfer(buf)
	if err != nil {
		return err
	}

	return nil
}

func sendBufferLine(position int, buffer []byte, device *spidev.SPIDevice) error {
	reg := MAX7219_REG_DIGIT0 + position
	fmt.Printf("Register: %#x\n", reg)
	buf := make([]byte, 2)
	for i := 0; i < 1; i++ {
		b := buffer[i*MAX7219_DIGIT_COUNT+position]
		fmt.Printf("Buffer value: %#x\n", b)
		buf[i*2] = byte(reg)
		buf[i*2+1] = b
	}
	fmt.Printf("Send to bus: %v\n", buf)
	_, err := device.Xfer(buf)
	if err != nil {
		return err
	}
	return nil
}

func flush(device *spidev.SPIDevice, buffer []byte) error {
	for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
		err := sendBufferLine(i, buffer, device)
		if err != nil {
			return err
		}
	}
	return nil
}

func clearDisplay(device *spidev.SPIDevice) error {
	err := flush(device, MAX7219_CLEAR_DISPLAY)
	if err != nil {
		return err
	}

	return nil
}

func displayTemp(first, second int, spi *spidev.SPIDevice) {
	for i := 0; i < 8; i++ {
		sendMulti(Max7219Reg(MAX7219_REG_DIGIT0+i), DEFAULT_FONT[first][i], DEFAULT_FONT[second][i], spi)
	}
}

func main() {



	// Setup
	devstr := fmt.Sprintf("/dev/spidev%d.%d", 0, 0)
	spi, err := spidev.NewSPIDevice(devstr)
	if err != nil {
		fmt.Printf("Error during device init: %s", err.Error())
		return
	}

	defer spi.Close()

	// Init driver
	send(MAX7219_REG_SCANLIMIT, 7, spi)   // show all 8 digits
	send(MAX7219_REG_DECODEMODE, 0, spi)  // use matrix (not digits)
	send(MAX7219_REG_DISPLAYTEST, 0, spi) // no display test
	send(MAX7219_REG_SHUTDOWN, 1, spi)    // not shutdown mode
	send(MAX7219_REG_INTENSITY, 7, spi)   // Set brightness (0-15)

	clearDisplay(spi)

	displayTemp(5, 2, spi)
}

//const POLLING_TIME = 2 * time.Second
//
//func pollTemperature(tempChan chan float32) {
//	for {
//		// Switch to poll from server
//		randTemp := rand.Float32() * 30.0
//
//		tempChan <- randTemp
//		time.Sleep(POLLING_TIME)
//	}
//}
//
//func main() {
//	tempChan := make(chan float32)
//	go pollTemperature(tempChan)
//
//	for temp := range tempChan {
//		fmt.Printf("Current temperature: %d\n", int32(temp))
//	}
//}

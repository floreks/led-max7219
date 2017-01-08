package device

import (
	"github.com/fulr/spidev"
	"fmt"
)

// Max7219 - represents 2 cascaded 8x8 devices. Code is written specially for this device
// https://botland.com.pl/led-do-raspberry-pi-32/4509-matryca-128-led-16x8-max7219-do-raspberry-pi.html
// may not work correctly with other configurations.
type Max7219 struct {
	buffer []byte
	device *spidev.SPIDevice
}

func (this Max7219) Close() {
	this.device.Close()
}

func NewMax7219() (*spidev.SPIDevice, error) {
	devstr := fmt.Sprintf("/dev/spidev%d.%d", 0, 0)

	spi, err := spidev.NewSPIDevice(devstr)
	if err != nil {
		return nil, err
	}

	return spi, nil
}
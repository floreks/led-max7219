package device

import (
	"fmt"

	"github.com/fulr/spidev"
)

type Max7219Reg byte

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

// Max7219 - represents 2 cascaded 8x8 devices. Code is written specially for this device
// https://botland.com.pl/led-do-raspberry-pi-32/4509-matryca-128-led-16x8-max7219-do-raspberry-pi.html
// may not work correctly with other configurations.
type Max7219 struct {
	device *spidev.SPIDevice
}

// Init device and Max7219 driver
func (this Max7219) Init() error {
	err := this.Command(MAX7219_REG_SCANLIMIT, 7)   // show all 8 digits
	if err != nil {
		return err
	}

	this.Command(MAX7219_REG_DECODEMODE, 0)  // use matrix (not digits)
	if err != nil {
		return err
	}

	this.Command(MAX7219_REG_DISPLAYTEST, 0) // no display test
	if err != nil {
		return err
	}

	this.Command(MAX7219_REG_SHUTDOWN, 1)    // not shutdown mode
	if err != nil {
		return err
	}

	this.Command(MAX7219_REG_INTENSITY, 0)   // Set brightness (0-15)
	if err != nil {
		return err
	}

	return nil
}

func (this Max7219) SetBrightness(value int) {
	brightness := value
	if brightness < 0 || brightness > 15 {
		brightness = 7
	}

	this.Command(MAX7219_REG_INTENSITY, byte(brightness))
}

func (this Max7219) Command(reg Max7219Reg, value byte) error {
	buf := []byte{byte(reg), value, byte(reg), value}

	_, err := this.device.Xfer(buf)
	if err != nil {
		return err
	}

	return nil
}

func (this Max7219) send(reg Max7219Reg, values ...byte) error {
	buf := make([]byte, 0)
	for _, value := range values {
		buf = append(buf, byte(reg)+1)
		buf = append(buf, value)
	}

	_, err := this.device.Xfer(buf)
	if err != nil {
		return err
	}

	return nil
}

func (this Max7219) ClearDisplay() error {
	for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
		err := this.send(Max7219Reg(i), 0, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO add doc
func (this Max7219) DisplayTemperature(temperature float32) error {
	// Round to integer number as we don't have space for floating point numbers
	temp := int(temperature)

	buf, err := DEFAULT_FONT.FromInt(temp)
	if err != nil {
		return err
	}

	for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
		err := this.send(Max7219Reg(i), buf[1][i], buf[0][i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (this Max7219) DisplayAlcoholConcentration(perMilles float32) error {
	// Round to integer number as we don't have space for floating point numbers
	if perMilles < 1.0 {
		perMilles *= 10
	}

	display := int(perMilles)

	buf, err := DEFAULT_FONT.FromInt(display)
	if err != nil {
		return err
	}

	for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
		// Add dot to indicate floating point number on first device (row/col) 7/8
		if i == 6 {
			buf[0][i] |= 1
		}

		err := this.send(Max7219Reg(i), buf[1][i], buf[0][i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (this Max7219) Close() {
	this.device.Close()
}

func NewMax7219() (*Max7219, error) {
	devstr := fmt.Sprintf("/dev/spidev%d.%d", 0, 0)

	spi, err := spidev.NewSPIDevice(devstr)
	if err != nil {
		return nil, err
	}

	return &Max7219{device: spi}, nil
}

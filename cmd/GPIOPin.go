package cmd

import (
	"fmt"
	"strconv"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type GPIOPin struct {
	Pin gpio.PinIO
}

// NewGPIOPin creates a new PTT instance for the specified GPIO pin number.
func NewGPIOPin(pinNum int) (*GPIOPin, error) {
	_, err := host.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize host: %w", err)
	}

	pin := gpioreg.ByName(strconv.Itoa(pinNum))
	if pin == nil {
		return nil, fmt.Errorf("failed to get GPIO %d", pinNum)
	}
	pin.Out(gpio.Low)

	return &GPIOPin{Pin: pin}, nil
}

func (p *GPIOPin) Set(level gpio.Level) {
	p.Pin.Out(level)
}

func (p *GPIOPin) On() {
	p.Pin.Out(gpio.High)
	// sleep for 300ms to ensure the PTT is fully engaged
	time.Sleep(270 * time.Millisecond)
}

func (p *GPIOPin) Off() {
	p.Pin.Out(gpio.Low)
}

func (p *GPIOPin) Read() gpio.Level {
	return p.Pin.Read()
}

func (p *GPIOPin) SetInput() {
	p.Pin.In(gpio.PullDown, gpio.NoEdge)
}

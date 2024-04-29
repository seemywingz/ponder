package cmd

import (
	"fmt"
	"strconv"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type PTT struct {
	Pin gpio.PinIO
}

// NewPTT creates a new PTT instance for the specified GPIO pin number.
func NewPTT(pinNum int) (*PTT, error) {
	_, err := host.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize host: %w", err)
	}

	pin := gpioreg.ByName(strconv.Itoa(pinNum))
	if pin == nil {
		return nil, fmt.Errorf("failed to get GPIO %d", pinNum)
	}
	pin.Out(gpio.Low)

	return &PTT{Pin: pin}, nil
}

func (p *PTT) Set(level gpio.Level) {
	p.Pin.Out(level)
}

func (p *PTT) On() {
	p.Pin.Out(gpio.High)
}

func (p *PTT) Off() {
	p.Pin.Out(gpio.Low)
}

/*
Copyright © 2024 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

var ptt int = -1
var pttPin gpio.PinIO = nil

// radioCmd represents the radio command
var radioCmd = &cobra.Command{
	Use:   "radio",
	Short: "An AI HAM Radio Operator",
	Long: `An AI HAM Radio Operator that can communicate with other HAM Radio Operators.
	When connected to a Raspberry Pi, or another computer with GPIO pins, the AI can transmit and receive messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		radio()
	},
}

func init() {
	rootCmd.AddCommand(radioCmd)

	radioCmd.Flags().IntVarP(&ptt, "ptt", "p", -1, "GPIO pin for Push To Talk (PTT) control")
}

func radio() {
	// Load all the drivers:
	_, err := host.Init()
	catchErr(err, "fatal")

	if ptt >= 0 {
		pttPin = gpioreg.ByName(strconv.Itoa(ptt))
		if pttPin == nil {
			catchErr(errors.New("Failed to get GPIO"+strconv.Itoa(ptt)), "fatal")
		}
		pttPin.Out(gpio.Low)
	}
	ttsRes, err := ai.TTS("Say hello and introduce yourself.")
	catchErr(err, "warn")
	if pttPin != nil {
		pttPin.Out(gpio.High)
	}
	playAudio(ttsRes)
}

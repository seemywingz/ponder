/*
Copyright © 2024 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var ptt int = -1
var pttPin gpio.PinIO

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

	// Load all the drivers:
	_, err := host.Init()
	catchErr(err, "fatal")

	if ptt >= 0 {
		pttPin := gpioreg.ByName(string(ptt))
		if pttPin == nil {
			catchErr(errors.New("Failed to get GPIO"+string(ptt)), "fatal")
		}
		pttPin.Out(gpio.OUT_LOW)
	}

}

func radio() {
	ai.TTS("Say hello and introduce yourself.")
}

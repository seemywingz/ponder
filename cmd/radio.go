/*
Copyright © 2024 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/seemywingz/goai"
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

	radioCmd.Flags().IntVar(&ptt, "ptt", -1, "GPIO pin for Push To Talk (PTT)")
}

func togglePTT() {
	if pttPin != nil {
		if pttPin.Read() == gpio.Low {
			pttPin.Out(gpio.High)
		} else {
			pttPin.Out(gpio.Low)
		}
	}
}

func radio() {

	// Load all the drivers:
	_, err := host.Init()
	catchErr(err, "fatal")

	if ptt >= 0 {
		if verbose {
			fmt.Print("Setting up GPIO" + strconv.Itoa(ptt) + " for Push To Talk (PTT)")
		}
		pttPin = gpioreg.ByName(strconv.Itoa(ptt))
		if pttPin == nil {
			catchErr(errors.New("Failed to get GPIO"+strconv.Itoa(ptt)), "fatal")
		}
		pttPin.Out(gpio.Low)
	}

	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "user",
		Content: "Say Hello and introduce yourself.",
	})

	ttsText, err := ai.ChatCompletion(ponderMessages)
	catchErr(err, "warn")
	ttsAudio, err := ai.TTS(ttsText.Choices[0].Message.Content)
	catchErr(err, "warn")

	pttPin.Out(gpio.High)
	playAudio(ttsAudio)
	pttPin.Out(gpio.Low)

	tick := time.Tick(1 * time.Second)
	quit := make(chan bool)

	go func() {
		time.Sleep(10 * time.Second)
		quit <- true
	}()

	for {
		select {
		case <-tick:
			fmt.Println("Tick")
		case <-quit:
			fmt.Println("Quit")
			return
		}
	}
}

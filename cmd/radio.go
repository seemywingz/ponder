/*
Copyright © 2024 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seemywingz/goai"
	"github.com/spf13/cobra"
	"periph.io/x/conn/v3/gpio"
)

var pttPinNum int = -1
var spkPinNum int = -1

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

	radioCmd.Flags().IntVarP(&pttPinNum, "ptt", "p", -1, "GPIO pin for PTT")
	radioCmd.Flags().IntVarP(&spkPinNum, "speaker", "s", -1, "GPIO pin for Speaker")
}

func radio() {

	var ptt *GPIOPin
	var spk *GPIOPin
	var err error

	if pttPinNum >= 0 {
		ptt, err = NewGPIOPin(pttPinNum)
		catchErr(err, "warn")
	}

	if spkPinNum >= 0 {
		spk, err = NewGPIOPin(spkPinNum)
		catchErr(err, "warn")
		spk.SetInput()
	}

	cleanup := make(chan bool)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGSTOP, syscall.SIGKILL)

	go func() {
		<-sigs
		cleanup <- true
	}()

	go func() {
		for {
			select {
			case <-cleanup:
				ptt.Off()
				fmt.Println("Ponder Radio Exiting now.")
				os.Exit(0)
			default:
				if spk != nil {
					if spk.Read() == gpio.High {
						fmt.Println("Data received on speaker pin")
					}
					time.Sleep(time.Millisecond * 100) // Adjust polling rate as necessary
				}
			}
		}
	}()

	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "user",
		Content: "Say Hello and introduce yourself. Share some radio knowledge.",
	})

	ttsText, err := ai.ChatCompletion(ponderMessages)
	catchErr(err, "warn")
	ttsAudio, err := ai.TTS(ttsText.Choices[0].Message.Content)
	catchErr(err, "warn")

	ptt.On()
	playAudio(ttsAudio)
	ptt.Off()

	<-cleanup // Wait for cleanup signal before exiting normally
}

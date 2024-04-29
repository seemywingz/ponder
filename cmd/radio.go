/*
Copyright © 2024 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/seemywingz/goai"
	"github.com/spf13/cobra"
)

var pttPinNum int = -1
var ptt *PTT

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

	radioCmd.Flags().IntVar(&pttPinNum, "ptt", -1, "GPIO pin for Push To Talk (PTT)")
}

func radio() {

	if pttPinNum >= 0 {
		ptt = new(PTT)
	}

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

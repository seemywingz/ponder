/*
Copyright © 2024 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/seemywingz/goai"
	"github.com/spf13/cobra"
)

var pttPinNum int = -1

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

	var ptt *PTT
	var err error

	if pttPinNum >= 0 {
		ptt, err = NewPTT(pttPinNum)
		catchErr(err, "fatal")
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

	// Create a channel to receive signals
	sigs := make(chan os.Signal, 1)
	// Notify sigs channel on SIGINT or SIGTERM
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// Create a channel to signal to finish
	done := make(chan bool, 1)

	// Goroutine to handle received signals
	go func() {
		sig := <-sigs
		fmt.Println("\nReceived signal:", sig)
		ptt.Off()
		done <- true
	}()

	<-done
}

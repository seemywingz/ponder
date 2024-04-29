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
		catchErr(err, "warn")
	}

	cleanupChan := make(chan bool)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGSTOP, syscall.SIGKILL)

	go func() {
		<-sigs
		cleanupChan <- true
	}()

	go func() {
		<-cleanupChan
		ptt.Off()
		fmt.Println("Cleanup complete. Exiting now.")
		os.Exit(0)
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

	// If the program reaches here without interruption, we should wait for an interrupt
	<-cleanupChan // Wait for cleanup signal before exiting normally
}

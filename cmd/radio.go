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

var ptt, spk *GPIOPin

var err error

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

	radioCmd.Flags().IntVar(&pttPinNum, "ptt", -1, "GPIO pin for PTT")
	radioCmd.Flags().IntVar(&spkPinNum, "spk", -1, "GPIO pin for Speaker")
}

func tx(audio []byte) {

	ptt.On()
	// wait for the PTT to engage
	time.Sleep(300 * time.Millisecond)
	playAudio(audio)
	ptt.Off()
}

func radio() {

	if pttPinNum >= 0 {
		ptt, err = NewGPIOPin(pttPinNum)
		catchErr(err, "warn")
	}

	if spkPinNum >= 0 {
		spk, err = NewGPIOPin(spkPinNum)
		catchErr(err, "warn")
		spk.SetInput()
	}

	// ponderMessages = append(ponderMessages, goai.Message{
	// 	Role:    "user",
	// 	Content: "Say Hello and introduce yourself.",
	// })

	ttsText := chatCompletion("Say Hello and introduce yourself.")
	ttsAudio, err := ai.TTS(ttsText)
	catchErr(err, "warn")
	tx(ttsAudio)

	cleanup := make(chan bool)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGSTOP, syscall.SIGKILL)

	go func() {
		<-sigs
		cleanup <- true
	}()

	lastSpeakerState := gpio.Low
	debounceDuration := time.Millisecond * 30

	go func() {
		for {
			select {
			case <-cleanup:
				ptt.Off()
				fmt.Println("Ponder Radio Exiting now.")
				os.Exit(0)
			default:
				currentSpeakerState, stable := spk.debouncePin(debounceDuration)
				if stable && currentSpeakerState != lastSpeakerState {
					if currentSpeakerState == gpio.High {
						fmt.Println("Data receiving started")
					} else if currentSpeakerState == gpio.Low {
						fmt.Println("Data receiving ended")
						ttsText := chatCompletion("Provide a question and answer from the HAM radio technician's manual.")
						ttsAudio, err := ai.TTS(ttsText)
						catchErr(err, "warn")
						tx(ttsAudio)
					}
					lastSpeakerState = currentSpeakerState
				}
				time.Sleep(time.Millisecond * 10) // polling interval
			}
		}
	}()

	<-cleanup // Wait for cleanup signal before exiting normally
}

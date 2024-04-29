/*
Copyright © 2024 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/spf13/cobra"
)

var audioFile,
	voice string

// ttsCmd represents the tts command
var ttsCmd = &cobra.Command{
	Use:   "tts",
	Short: "OpenAI Text to Speech API - TTS",
	Long: `OpenAI Text to Speech API - TTS
	You can use the TTS API to generate audio from text.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		audio := tts(prompt)
		if audio != nil {
			playAudio(audio)
		}
	},
}

func init() {
	rootCmd.AddCommand(ttsCmd)
	ttsCmd.Flags().StringVarP(&audioFile, "file", "f", "", "File to save audio to")
}

func tts(text string) []byte {
	ai.Voice = voice
	audioData, err := ai.TTS(text)
	catchErr(err, "fatal")
	if audioFile != "" {
		file, err := os.Create(audioFile)
		catchErr(err)
		defer file.Close()
		_, err = io.Copy(file, bytes.NewReader(audioData))
		catchErr(err)
		return nil
	}
	return audioData
}

// playAudio plays audio from a byte slice.
func playAudio(audioContent []byte) {
	if verbose {
		fmt.Println("🔊 Playing audio...")
	}

	// Create an io.Reader from the byte slice
	reader := bytes.NewReader(audioContent)

	// Wrap the reader in a NopCloser to make it an io.ReadCloser.
	readCloser := io.NopCloser(reader)

	// Decode the MP3 stream.
	streamer, format, err := mp3.Decode(readCloser)
	catchErr(err)
	defer streamer.Close()

	// Initialize the speaker with the sample rate of the audio and a buffer size.
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	catchErr(err)

	// Play the decoded audio.
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for the audio to finish playing.
	<-done
}

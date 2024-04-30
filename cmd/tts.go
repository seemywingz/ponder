/*
Copyright © 2024 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"bytes"
	"io"
	"os"

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

/*
Copyright © 2023 NAME HERE Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/seemywingz/goai"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(chatCmd)
}

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Open ended chat with OpenAI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if convo {
			for {
				response, audio := chatResponse(prompt)
				fmt.Println("\nPonder:\n  " + response + "\n")
				if narrate {
					playAudio(audio)
				}
				fmt.Print("You:\n  ")
				prompt, err = getUserInput()
				catchErr(err, "warn")
			}
		} else {
			response, audio := chatResponse(prompt)
			fmt.Println(response)
			if narrate {
				playAudio(audio)
			}
		}
	},
}

func chatResponse(prompt string) (string, []byte) {
	var audio []byte
	var response string
	spinner, _ := pterm.DefaultSpinner.Start("Pondering...")
	spinner.RemoveWhenDone = true
	response = chatCompletion(prompt)
	if narrate {
		audio = tts(response)
	}
	spinner.Stop()
	return response, audio
}

func chatCompletion(prompt string) string {
	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "user",
		Content: prompt,
	})

	// Send the messages to OpenAI
	res, err := ai.ChatCompletion(ponderMessages)
	catchErr(err)
	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "assistant",
		Content: res.Choices[0].Message.Content,
	})
	return res.Choices[0].Message.Content
}

func getUserInput() (string, error) {
	// ReadString will block until the delimiter is entered
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		trace()
		return "", err
	}
	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	if verbose {
		trace()
		fmt.Println(input)
	}
	return input, nil
}

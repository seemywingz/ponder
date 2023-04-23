/*
Copyright © 2023 NAME HERE Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/seemywingz/goai"
	"github.com/spf13/cobra"
)

var convo bool
var sayText bool
var ponderMessages = []goai.Message{{
	Role:    "system",
	Content: ponder_SystemMessage,
}}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().BoolVarP(&convo, "convo", "c", false, "Conversational Style chat")
	chatCmd.Flags().BoolVarP(&sayText, "say", "s", false, "Say text out loud (MacOS only)")
}

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Open ended chat with OpenAI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if convo {
			for {
				fmt.Println("\nYou: ")
				prompt, err := getUserInput()
				catchErr(err)
				fmt.Println("\nPonder:\n", chatCompletion(prompt))
			}
		} else {
			textCompletion(prompt)
		}

	},
}

func chatCompletion(prompt string) string {
	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "user",
		Content: prompt,
	})

	// Send the messages to OpenAI
	oaiResponse, err := openai.ChatCompletion(ponderMessages)
	catchErr(err)
	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "assistant",
		Content: oaiResponse,
	})
	return oaiResponse
}

func textCompletion(prompt string) {

	oaiResponse, err := openai.TextCompletion(prompt)
	catchErr(err)

	for _, v := range oaiResponse.Choices {
		text := v.Text
		if runtime.GOOS == "darwin" && sayText {
			say(text)
		}
		fmt.Println(text[2:])
	}

}

func getUserInput() (string, error) {
	// ReadString will block until the delimiter is entered
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		trace()
		return "", err
	}
	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\n")
	if verbose {
		trace()
		fmt.Println(input)
	}
	return input, nil
}

func say(phrase string) {
	say := exec.Command(`say`, phrase)
	err := say.Start()
	if err != nil {
		fmt.Println(err)
	}
}

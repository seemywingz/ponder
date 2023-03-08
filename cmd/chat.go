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

	"github.com/spf13/cobra"
)

var loop bool
var sayText bool

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Openended chat with ChatGPT",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if loop {
			for {
				fmt.Println("\nYou: ")
				prompt, err := getUserInput()
				catchErr(err)
				fmt.Println("\nPonder: ")
				chat(prompt)
			}
		} else {
			chat(prompt)
		}

	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().BoolVarP(&loop, "loop", "l", false, "Loop chat")
	chatCmd.Flags().BoolVarP(&sayText, "say", "s", false, "Say text out loud (MacOS only)")
}

func say(phrase string) {
	say := exec.Command(`say`, phrase)
	err := say.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func chat(prompt string) {

	oaiResponse := openAI_Chat(prompt)

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

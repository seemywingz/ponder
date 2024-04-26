/*
Copyright © 2023 NAME HERE Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
			fmt.Println(chatCompletion(prompt))
		}
	},
}

func chatCompletion(prompt string) string {
	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "user",
		Content: prompt,
	})

	fmt.Println("Pondering...")

	// Send the messages to OpenAI
	res, err := ai.ChatCompletion(ponderMessages)
	catchErr(err)
	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "assistant",
		Content: res.Choices[0].Message.Content,
	})
	return res.Choices[0].Message.Content
}

// func cliCommand(command string, args ...string) {
// 	cli := exec.Command(command, args...)
// 	output, err := cli.Output()
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println(string(output))
// 	}
// }

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

func say(phrase string) {
	say := exec.Command(`say`, phrase)
	err := say.Start()
	if err != nil {
		fmt.Println(err)
	}
}

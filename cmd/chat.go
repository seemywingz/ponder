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
	oaiResponse, err := ai.ChatCompletion(ponderMessages)
	catchErr(err)
	ponderMessages = append(ponderMessages, goai.Message{
		Role:    "assistant",
		Content: oaiResponse.Choices[0].Message.Content,
	})
	return oaiResponse.Choices[0].Message.Content
}

func textCompletion(prompt string) {

	if perform {
		prompt = command_SystemMessage + "\n here is the prompt:\n" + prompt
	}

	oaiResponse, err := ai.TextCompletion(prompt)
	catchErr(err)

	for _, v := range oaiResponse.Choices {
		text := v.Text
		if runtime.GOOS == "darwin" && sayText {
			say(text)
		}
		fmt.Println(text[2:])
	}

	if perform {
		command := strings.Split(oaiResponse.Choices[0].Text, " ")
		// fmt.Println("Running command: ", strings.ReplaceAll(command[0], "\n", ""), command[1:])
		cliCommand(strings.ReplaceAll(command[0], "\n", ""), command[1:]...)
	}

}

func cliCommand(command string, args ...string) {
	cli := exec.Command(command, args...)
	output, err := cli.Output()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
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

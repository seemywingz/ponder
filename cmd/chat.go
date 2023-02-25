/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Openended chat with ChatGPT",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		chat(prompt)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}

func chat(prompt string) {
	// fmt.Println("Welcome to Ponder! Ask me anything!")

	// for {
	// q, err := getUserInput()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	ans, err := getChatResponse(prompt)
	catchErr(err)

	fmt.Println()
	// fmt.Print("Ponder: ")
	for _, v := range ans.Choices {
		say(v.Text)
		fmt.Println(v.Text)
		fmt.Println()
	}

	// }
}

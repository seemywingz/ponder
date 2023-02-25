/*
Copyright © 2023 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ponder",
	Short: "A ChatGPT Chat Bot",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Ponder! Ask me anything!")

		for {
			q, err := getUserInput()
			if err != nil {
				fmt.Println(err)
			}

			ans, err := getChatResponse(q)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println()
				fmt.Print("Ponder: ")
				for _, v := range ans.Choices {
					say(v.Text)
					fmt.Println(v.Text)
					fmt.Println()
				}
			}
		}
	},
}

func say(phrase string) {
	say := exec.Command(`say`, phrase)
	err := say.Start()
	if err != nil {
		fmt.Println(err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ponder.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

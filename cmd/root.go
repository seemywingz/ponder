/*
Copyright © 2023 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"hash/fnv"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var verbose bool
var prompt string
var openAIUser string
var OPENAI_API_KEY string
var ETSY_API_KEY string
var PRINTIFY_API_KEY string
var DISCORD_API_KEY string
var DISCORD_PUB_KEY string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ponder",
	Short: "A OpenAI Chat Bot",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {

	// },
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&prompt, "prompt", "p", "", "Prompt AI generation")
	rootCmd.MarkFlagRequired("prompt")

	// Check for Required Environment Variables
	OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		catchErr(errors.New("OPENAI_API_KEY environment variable is not set"))
	}
	PRINTIFY_API_KEY = os.Getenv("PRINTIFY_API_KEY")
	if PRINTIFY_API_KEY == "" {
		catchErr(errors.New("PRINTIFY_API_KEY environment variable is not set"))
	}
	DISCORD_API_KEY = os.Getenv("DISCORD_API_KEY")
	if DISCORD_API_KEY == "" {
		catchErr(errors.New("DISCORD_API_KEY environment variable is not set"))
	}
	DISCORD_PUB_KEY = os.Getenv("DISCORD_PUB_KEY")
	if DISCORD_PUB_KEY == "" {
		catchErr(errors.New("DISCORD_PUB_KEY environment variable is not set"))
	}

	// Create a unique user for OpenAI
	h := fnv.New32a()
	h.Write([]byte(OPENAI_API_KEY))
	openAIUser = "ponder" + strconv.Itoa(int(h.Sum32()))
}

func trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d\n%s\n", file, line, f.Name())
}

func catchErr(err error) {
	if err != nil {
		// trace()
		fmt.Println("💔", err)
		os.Exit(1)
	}
}

func formatPrompt(prompt string) string {
	prompt = strings.ReplaceAll(prompt, " ", "_")
	prompt = strings.ReplaceAll(prompt, "/", "-")
	prompt = strings.ReplaceAll(prompt, ",", "")
	return prompt
}

func fileNameFromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	catchErr(err)

	// Get the last path component of the URL
	filename := filepath.Base(u.Path)

	// Replace any characters that are not letters, numbers, or underscores with dashes
	filename = regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(filename, "-")

	return filename
}

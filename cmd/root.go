/*
Copyright © 2023 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
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
var APP_VERSION = "v0.0.0"
var prompt,
	openAIUser,
	ponderID,
	OPENAI_API_KEY,
	PRINTIFY_API_KEY,
	DISCORD_API_KEY,
	DISCORD_PUB_KEY string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ponder",
	Short: "Ponder OpenAI Chat Bot " + APP_VERSION,
	Long: `
	Ponder
	GitHub: https://github.com/seemywingz/ponder
	App Version: ` + APP_VERSION + `

  Ponder uses OpenAI's GPT-3.5-Turbo API to generate text responses to user input.
  You can use Ponder as a Discord chat bot or to generate images using the DALL-E API.
  Or whatever else you can think of...
	`,
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
		fmt.Println("OPENAI_API_KEY environment variable is not set, continuing without OpenAI API Key")
	}
	PRINTIFY_API_KEY = os.Getenv("PRINTIFY_API_KEY")
	if PRINTIFY_API_KEY == "" {
		fmt.Println("PRINTIFY_API_KEY environment variable is not set, continuing without Printify API Key")
	}
	DISCORD_API_KEY = os.Getenv("DISCORD_API_KEY")
	if DISCORD_API_KEY == "" {
		fmt.Println("DISCORD_API_KEY environment variable is not set, continuing without Discord API Key")
	}
	DISCORD_PUB_KEY = os.Getenv("DISCORD_PUB_KEY")
	if DISCORD_PUB_KEY == "" {
		fmt.Println("DISCORD_PUB_KEY environment variable is not set, continuing without Discord Public Key")
	}

	// Create a unique user for OpenAI
	h := fnv.New32a()
	h.Write([]byte(OPENAI_API_KEY))
	ponderID = "ponder-" + strconv.Itoa(int(h.Sum32())) + "-"
	openAIUser = ponderID + "user"
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

/*
Copyright © 2023 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/seemywingz/goai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var APP_VERSION = "v0.1.0"
var verbose bool
var ai *goai.Client
var prompt,
	configFile,
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

	cobra.OnInitialize(viperConfig)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&prompt, "prompt", "p", "", "Prompt AI generation")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "$HOME/.ponder/config", "config file")
	rootCmd.MarkFlagRequired("prompt")

	// Check for Required Environment Variables
	OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		fmt.Println("⚠️ OPENAI_API_KEY environment variable is not set, continuing without OpenAI API Key")
	}
	PRINTIFY_API_KEY = os.Getenv("PRINTIFY_API_KEY")
	if PRINTIFY_API_KEY == "" {
		fmt.Println("⚠️ PRINTIFY_API_KEY environment variable is not set, continuing without Printify API Key")
	}
	DISCORD_API_KEY = os.Getenv("DISCORD_API_KEY")
	if DISCORD_API_KEY == "" {
		fmt.Println("⚠️ DISCORD_API_KEY environment variable is not set, continuing without Discord API Key")
	}
	DISCORD_PUB_KEY = os.Getenv("DISCORD_PUB_KEY")
	if DISCORD_PUB_KEY == "" {
		fmt.Println("⚠️ DISCORD_PUB_KEY environment variable is not set, continuing without Discord Public Key")
	}
}

func viperConfig() {
	// use spf13/viper to read config file
	// viper.AddConfigPath("$HOME/.ponder") // call multiple times to add many search paths
	// viper.SetConfigName("config")        // name of config file (without extension)
	// viper.AddConfigPath(".")             // optionally look for config in the working directory

	viper.SetDefault("openAI_endpoint", "https://api.openai.com/v1/")

	viper.SetDefault("openAI_image_size", "1024x1024")
	viper.SetDefault("openAI_image_downloadPath", "~/Ponder/Images/")

	viper.SetDefault("openAI_topP", "0.9")
	viper.SetDefault("openAI_frequencyPenalty", "0.0")
	viper.SetDefault("openAI_presencePenalty", "0.6")
	viper.SetDefault("openAI_temperature", "0")
	viper.SetDefault("openAI_maxTokens", "999")
	viper.SetDefault("openAI_chat_model", "gpt-3.5-turbo")
	viper.SetDefault("openAI_text_model", "text-davinci-003")

	viper.SetDefault("discord_message_context_count", "15")

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	if verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.ReadInConfig(); err != nil {
		if err != nil {
			// Config file not found; ignore error if desired
			fmt.Println("⚠️  Error Opening Config File:", err.Error(), "- Using Defaults")
		}
	}

	ai = &goai.Client{
		Endpoint:         viper.GetString("openAI_endpoint"),
		API_KEY:          OPENAI_API_KEY,
		Verbose:          verbose,
		ImageSize:        viper.GetString("openAI_image_size"),
		User:             goai.HashAPIKey(OPENAI_API_KEY),
		TopP:             viper.GetFloat64("openAI_topP"),
		ChatModel:        viper.GetString("openAI_chat_model"),
		TextModel:        viper.GetString("openAI_text_model"),
		MaxTokens:        viper.GetInt("openAI_maxTokens"),
		Temperature:      viper.GetFloat64("openAI_temperature"),
		FrequencyPenalty: viper.GetFloat64("openAI_frequencyPenalty"),
		PresencePenalty:  viper.GetFloat64("openAI_presencePenalty"),
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
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
		fmt.Println("💔", err)
		// os.Exit(1)
	}
}

func formatPrompt(prompt string) string {
	// Replace any characters that are not letters, numbers, or underscores with dashes
	return regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(prompt, "-")
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

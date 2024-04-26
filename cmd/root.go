/*
Copyright © 2023 Kevin Jayne <kevin.jayne@icloud.com>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/seemywingz/goai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verbose bool
var perform bool
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

  Ponder uses OpenAI's API to generate text responses to user input.
  You can use Ponder as a Discord chat bot or to generate images using the DALL-E API.
  Or whatever else you can think of...
	`,
	Run: func(cmd *cobra.Command, args []string) {
		chatCmd.Run(cmd, args)
	},
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

	rootCmd.MarkFlagRequired("prompt")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	rootCmd.Flags().BoolVarP(&convo, "convo", "c", false, "Conversational Style chat")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&sayText, "say", "s", false, "Say text out loud (MacOS only)")
	rootCmd.PersistentFlags().StringVarP(&prompt, "prompt", "p", "", "Prompt AI generation")
	rootCmd.Flags().BoolVarP(&perform, "perform", "x", false, "Attempt to perform the response as cli command")

	// Check for Required Environment Variables
	OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" && verbose {
		fmt.Println("⚠️ OPENAI_API_KEY environment variable is not set, continuing without OpenAI API Key")
	}
	PRINTIFY_API_KEY = os.Getenv("PRINTIFY_API_KEY")
	if PRINTIFY_API_KEY == "" && verbose {
		fmt.Println("⚠️ PRINTIFY_API_KEY environment variable is not set, continuing without Printify API Key")
	}
	DISCORD_API_KEY = os.Getenv("DISCORD_API_KEY")
	if DISCORD_API_KEY == "" && verbose {
		fmt.Println("⚠️ DISCORD_API_KEY environment variable is not set, continuing without Discord API Key")
	}
	DISCORD_PUB_KEY = os.Getenv("DISCORD_PUB_KEY")
	if DISCORD_PUB_KEY == "" && verbose {
		fmt.Println("⚠️ DISCORD_PUB_KEY environment variable is not set, continuing without Discord Public Key")
	}
}

func viperConfig() {
	// use spf13/viper to read config file

	viper.SetDefault("openAI_endpoint", "https://api.openai.com/v1/")

	viper.SetDefault("openAI_image_size", "1024x1024")
	viper.SetDefault("openAI_image_downloadPath", "~/Ponder/Images/")

	viper.SetDefault("openAI_topP", "0.9")
	viper.SetDefault("openAI_frequencyPenalty", "0.0")
	viper.SetDefault("openAI_presencePenalty", "0.6")
	viper.SetDefault("openAI_temperature", "0")
	viper.SetDefault("openAI_maxTokens", "999")
	viper.SetDefault("openAI_chat_model", "gpt-4")
	viper.SetDefault("openAI_image_model", "dall-e-3")
	viper.SetDefault("openAI_text_model", "text-davinci-003")

	viper.SetDefault("discord_message_context_count", "15")

	viper.SetConfigName("config")        // name of config file (without extension)
	viper.SetConfigType("yaml")          // REQUIRED the config file does not have an extension
	viper.AddConfigPath("$HOME/.ponder") // call multiple times to add many search paths
	viper.AddConfigPath("./files")       // look for config in the working directory /files
	viper.AddConfigPath(".")             // look for config in the working directory

	if configFile != "" {
		viper.SetConfigFile(configFile)
		if verbose {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if err != nil {
			// Config file not found; ignore error if desired
			fmt.Println("⚠️  Error Opening Config File:", err.Error(), "- Using Defaults")
		}
	} else {
		if verbose {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
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
		ImageModel:       viper.GetString("openAI_image_model"),
		TTSModel:         viper.GetString("openAI_tts_model"),
		MaxTokens:        viper.GetInt("openAI_maxTokens"),
		Temperature:      viper.GetFloat64("openAI_temperature"),
		FrequencyPenalty: viper.GetFloat64("openAI_frequencyPenalty"),
		PresencePenalty:  viper.GetFloat64("openAI_presencePenalty"),
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

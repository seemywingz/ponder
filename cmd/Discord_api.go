package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var discord *discordgo.Session

func initDiscord() {
	var err error
	discord, err = discordgo.New("Bot " + DISCORD_API_KEY)
	catchErr(err)

	discord.Client = httpClient // Set the HTTP client for the Discord session.

	// Open a websocket connection to Discord
	err = discord.Open()
	catchErr(err)
	defer discord.Close()

	setStatusOnline()
	registerHandlers()
	registerSlashCommand()
	deregisterSlashCommands()

	log.Println("Ponder Discord Bot is now running...")
	select {}
}

func setStatusOnline() {
	fmt.Println("Setting Status to Online...")
	// Set status to online with active activity
	err := discord.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "online",
		Activities: []*discordgo.Activity{
			{ // Activity Type 0 is "Playing"
				Name: "With the API",
				Type: discordgo.ActivityType(0),
			},
		},
		AFK: false,
	})
	catchErr(err)
}

func registerHandlers() {
	fmt.Println("Registering Handlers...")
	// Register a new slash command handler
	discord.AddHandler(handleCommands)
}

func deregisterSlashCommands() {
	fmt.Println("Deregistering Slash Commands...")
	// Set the ID of the slash command to deregister
	slashCommandIDs := []string{}

	for _, slashCommandID := range slashCommandIDs {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, "", slashCommandID)
		if err != nil {
			fmt.Println("Error deleting slash command: ", err)
			return
		}
	}
}

func registerSlashCommand() {
	fmt.Println("Registering Slash Commands...")
	// /chat command
	command := &discordgo.ApplicationCommand{
		Name:        "chat",
		Description: "Chat with Ponder Discord Bot, powered by OpenAI GPT-3!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "The prompt message",
				Required:    true,
			},
		},
	}
	_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
	catchErr(err)

	// /chat command
	command = &discordgo.ApplicationCommand{
		Name:        "scrape",
		Description: "Scrape Discord channel for Upscaled Midjourney Images!",
	}
	_, err = discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
	catchErr(err)

	// /hello command
	command = &discordgo.ApplicationCommand{
		Name:        "hello",
		Description: "Say hello to Ponder Discord Bot!",
	}
	_, err = discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
	catchErr(err)

}

func handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discordInitialResponse("Thinking...", s, i)
	switch i.ApplicationCommandData().Name {
	case "chat":
		discordChat(s, i)
	case "scrape":
		scrapeImages(s, i)
	case "hello":
		oaiResponse := openAI_Chat("Hello World, from Ponder Discord Bot that is kinda cheeky!")
		responseMessage := ""
		for _, v := range oaiResponse.Choices {
			responseMessage += v.Text[2:]
		}
		discordFollowUp(responseMessage, s, i)
	default: // Handle unknown slash commands
		log.Printf("Unknown Ponder Command: %s", i.ApplicationCommandData().Name)
	}
}

func scrapeImages(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discordFollowUp("Scraping Discord Channel for Upscaled Midjourney Images...", s, i)
	// Get the interaction channel ID
	channelID := i.ChannelID
	messages, err := discord.ChannelMessages(channelID, 100, "", "", "")
	catchErr(err)

	for _, v := range messages {
		if strings.Contains(v.Content, "Upscaled") {
			fmt.Println()
			fmt.Println("Message:")
			fmt.Println("ID: " + v.ID)
			fmt.Println("Author: " + v.Author.Username)
			fmt.Println("Content: " + v.Content)
			if len(v.Attachments) > 0 {
				fmt.Println("Attachments: " + v.Attachments[0].URL)
			}
		}
	}

}

func discordChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the value of the prompt parameter
	prompt := i.ApplicationCommandData().Options[0].StringValue()
	oaiResponse := openAI_Chat(prompt)
	responseMessage := ""

	for _, v := range oaiResponse.Choices {
		responseMessage += v.Text[2:]
	}

	responseMessage = "prompt: " + prompt + "\n" + responseMessage
	discordFollowUp(responseMessage, s, i)

	if verbose {
		fmt.Println("Ponder Discord Bot...")
		fmt.Println("Response: " + responseMessage)
	}
}

func discordInitialResponse(content string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Send initial defer response.
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	catchErr(err)
}

func discordFollowUp(message string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	followup := &discordgo.WebhookParams{
		Content: message,
	}
	_, err := s.FollowupMessageCreate(i.Interaction, false, followup)
	catchErr(err)
}

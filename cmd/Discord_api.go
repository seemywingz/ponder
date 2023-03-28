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
	registerCommands()
	deregisterCommands()

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
	// A handler function for the "interactionCreate" event
	discord.AddHandler(handleCommands)
	// Add a handler function for the "messageCreate" event
	discord.AddHandler(handleMessages)
}

func deregisterCommands() {
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

func registerCommands() {
	fmt.Println("Registering Slash Commands...")
	// /chat command
	command := &discordgo.ApplicationCommand{
		Name:        "scrape",
		Description: "Scrape Discord channel for Upscaled Midjourney Images!",
	}
	_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
	catchErr(err)
}

func handleMessages(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == discord.State.User.ID {
		return
	}

	// get the channel object using the ID
	channel, err := discord.Channel(m.ChannelID)
	if err != nil {
		log.Fatal(err)
	}

	// get the name of the channel
	channelName := channel.Name
	fmt.Println("Channel Name: " + channelName)

	// Respond to messages in the #ponder channel
	if channelName == "ponder" {
		discordOpenAIResponse(s, m, false)
		return
	}

	// Check if the message contains an @mention of the bot.
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			// Send a reply to the user who mentioned the bot.
			discordOpenAIResponse(s, m, true)
			return
		}
	}

}

func handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discordInitialResponse("Thinking...", s, i)
	switch i.ApplicationCommandData().Name {
	case "scrape":
		discordScrapeImages(s, i)
	default: // Handle unknown slash commands
		log.Printf("Unknown Ponder Command: %s", i.ApplicationCommandData().Name)
	}
}

func discordOpenAIResponse(s *discordgo.Session, m *discordgo.MessageCreate, mention bool) {
	discord.ChannelTyping(m.ChannelID)
	response := openai_ChatTXTonly(m.Content)
	if mention {
		response = m.Author.Mention() + " " + response
	}
	s.ChannelMessageSend(m.ChannelID, response)
}

func discordScrapeImages(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				url := v.Attachments[0].URL
				fmt.Println("Attachments: " + url)
				printify_UploadImage(fileNameFromURL(url), url)
			}
		}
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

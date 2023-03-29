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

func handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discordInitialResponse("Thinking...", s, i)
	switch i.ApplicationCommandData().Name {
	case "scrape":
		discordScrapeImages(s, i)
	default: // Handle unknown slash commands
		log.Printf("Unknown Ponder Command: %s", i.ApplicationCommandData().Name)
	}
}

func handleMessages(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == discord.State.User.ID {
		return
	}

	channelName := discordGetChannelName(m.ChannelID)

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

func discordOpenAIResponse(s *discordgo.Session, m *discordgo.MessageCreate, mention bool) {
	discord.ChannelTyping(m.ChannelID)
	openaiMessages := []OPENAI_Message{{
		Role:    "system",
		Content: discord_BotSystemMessage,
	}}

	discordMessages, err := discord.ChannelMessages(m.ChannelID, 30, "", "", "")
	catchErr(err)
	discordMessages = discordReverseMessageOrder(discordMessages)

	for _, message := range discordMessages {
		role := "user"
		if message.Author.ID == discord.State.User.ID {
			role = "assistant"
		}
		newMessage := OPENAI_Message{
			Role:    role,
			Content: message.Content,
		}
		openaiMessages = append(openaiMessages, newMessage)
	}

	oaiResponse := openai_ChatComplete(openaiMessages)
	responseMessage := oaiResponse.Choices[0].Message.Content
	if mention {
		responseMessage = m.Author.Mention() + " " + responseMessage
	}
	s.ChannelMessageSend(m.ChannelID, responseMessage)
}

func discordGetChannelName(channelID string) string {
	channel, err := discord.Channel(channelID)
	catchErr(err)
	return channel.Name
}

func discordGetChannelID(s *discordgo.Session, guildID string, channelName string) string {
	channels, err := s.GuildChannels(guildID)
	catchErr(err)

	for _, channel := range channels {
		if channel.Name == channelName {
			return channel.ID
		}
	}

	return ""
}

func discordScrapeImages(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discord.ChannelTyping(i.ChannelID)
	// Get the interaction channel ID
	channelID := i.ChannelID
	messages, err := discord.ChannelMessages(channelID, 100, "", "", "")
	catchErr(err)
	discordFollowUp("Scraping Discord Channel for ALL Image URLs and sending them to #saved-images.\nAll `Upscaled` Midjourney Images will be sent to Printify as well...", s, i)

	savedImagesChannelID := discordGetChannelID(s, i.GuildID, "saved-images")

	for _, v := range messages {
		if len(v.Attachments) > 0 {
			url := v.Attachments[0].URL
			s.ChannelMessageSend(savedImagesChannelID, url)
			if strings.Contains(v.Content, "Upscaled") {
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

// function to reverse the order of a slice
func discordReverseMessageOrder(s []*discordgo.Message) []*discordgo.Message {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

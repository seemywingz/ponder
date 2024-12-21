package cmd

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/seemywingz/goai"
	"github.com/spf13/viper"
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

	log.Println("🤖 Ponder Discord Bot is Running...")
	select {} // Block forever to prevent the program from terminating.
}

func setStatusOnline() {
	log.Println("🛜  Setting Status to Online...")
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
	log.Println("💾 Registering Handlers...")
	discord.AddHandler(handleCommands)
	discord.AddHandler(handleMessages)
}

func deregisterCommands() {

	if removeCMDIds != "" {
		slashCommandIDs := strings.Split(removeCMDIds, ",")
		for _, slashCommandID := range slashCommandIDs {
			log.Println("➖ Removing Command: ", slashCommandID)
			err := discord.ApplicationCommandDelete(discord.State.User.ID, "", slashCommandID)
			if err != nil {
				log.Println("Error deleting slash command: ", err)
				return
			}
		}
	}
}

func registerCommands() {

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "scrape",
			Description: "Scrape Discord channel for Upscaled Midjourney Images",
		},
		{
			Name:        "ponder-image",
			Description: "Use DALL-E 3 to generate an Image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "Prompt for Image Generation",
					Required:    true,
				},
			},
		},
	}

	for _, command := range commands {
		log.Println("➕ Adding Command: /"+command.Name, "-", command.Description)
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
		catchErr(err)
	}
}

func handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discordInitialResponse("Pondering...", s, i)
	switch i.ApplicationCommandData().Name {
	case "scrape":
		discordScrapeImages(s, i)
	case "ponder-image":
		discordPonderImage(s, i)
	default: // Handle unknown slash commands
		log.Printf("Unknown Ponder Command: %s", i.ApplicationCommandData().Name)
	}
}

func handleMessages(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == discord.State.User.ID {
		return
	}

	// channelName := discordGetChannelName(m.ChannelID)

	// Respond to messages in the #ponder channel
	if m.GuildID == "" {
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
	openaiMessages := []goai.Message{{
		Role:    "system",
		Content: viper.GetString("discord_bot_systemMessage"),
	}}

	discordMessages, err := discord.ChannelMessages(m.ChannelID, viper.GetInt("discord_message_context_count"), "", "", "")
	catchErr(err)
	discordMessages = discordReverseMessageOrder(discordMessages)

	for _, message := range discordMessages {
		role := "user"
		if message.Author.ID == discord.State.User.ID {
			role = "assistant"
		}
		newMessage := goai.Message{
			Role:    role,
			Content: message.Content,
		}
		openaiMessages = append(openaiMessages, newMessage)
	}

	// Send the messages to OpenAI
	ai.User = ai.User + m.Author.Username
	oaiResponse, err := ai.ChatCompletion(openaiMessages)
	catchErr(err)
	s.ChannelMessageSend(m.ChannelID, oaiResponse.Choices[0].Message.Content)
}

func discordPonderImage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelID := i.ChannelID
	discord.ChannelTyping(channelID)
	commandData := i.ApplicationCommandData()

	// Check if there are options and retrieve the prompt
	if len(commandData.Options) > 0 {
		promptOption := commandData.Options[0]
		prompt := promptOption.StringValue()
		discordFollowUp("Using DALL-E 3 to generate an image: "+prompt, s, i)
		res := ai.ImageGen(prompt, "", 1)
		s.ChannelMessageSend(channelID, res.Data[0].URL)
	} else {
		discordFollowUp("Please Provide a Prompt for Image Generation", s, i)
	}

}

func discordScrapeImages(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the interaction channel ID
	channelID := i.ChannelID
	discord.ChannelTyping(channelID)
	messages, err := discord.ChannelMessages(channelID, 100, "", "", "")
	catchErr(err)
	discordFollowUp("Scraping Discord Channel for ALL Image URLs and sending them to #saved-images.\nAll `Upscaled` Midjourney Images will be sent to Printify as well...", s, i)

	savedImagesChannelID := discordGetChannelID(s, i.GuildID, "saved-images")

	for _, v := range messages {
		if len(v.Attachments) > 0 {
			url := v.Attachments[0].URL
			s.ChannelMessageSend(savedImagesChannelID, url)
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

package cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var discord *discordgo.Session

func initDiscord() {
	var err error
	discord, err = discordgo.New("Bot " + DISCORD_API_KEY)
	catchErr(err)

	// Open a websocket connection to Discord
	err = discord.Open()
	catchErr(err)
	defer discord.Close()

	setStatusOnline()
	registerHandlers()
	registerSlashCommand()

	log.Println("Ponder Discord Bot is now running...")
	select {}
}

// func discord_GetImage() {

// 	initDiscord()

// 	messages, err := discord.ChannelMessages(discord_channel_midjourney, 10, "", "", "")
// 	catchErr(err)

// 	for _, v := range messages {
// 		fmt.Println()
// 		fmt.Println("Message:")
// 		fmt.Println("ID: " + v.ID)
// 		fmt.Println("Author: " + v.Author.Username)
// 		fmt.Println("Content: " + v.Content)
// 		if len(v.Attachments) > 0 {
// 			fmt.Println("Attachments: " + v.Attachments[0].URL)
// 		}
// 	}

// }

func setStatusOnline() {
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
	// Register a new slash command handler
	discord.AddHandler(slashCommandHandler)
}

func registerSlashCommand() {

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

	// /hello command
	command = &discordgo.ApplicationCommand{
		Name:        "hello",
		Description: "Say hello to Ponder Discord Bot!",
	}
	_, err = discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
	catchErr(err)

}

func slashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "hello":

		oaiResponse := openAI_Chat("Hello World, from Ponder Discord Bot that is kinda cheeky!")
		responseMessage := ""

		for _, v := range oaiResponse.Choices {
			responseMessage += v.Text[2:]
		}

		discordSendMessagesResponse(responseMessage, s, i)
	case "chat":
		discordChat(s, i)
	default:
		// Handle unknown slash commands
		log.Printf("Unknown Ponder Command: %s", i.ApplicationCommandData().Name)
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

	if verbose {
		fmt.Println("Chatting with Ponder Discord Bot...")
		fmt.Println("Prompt: " + prompt)
		fmt.Println("Response: " + responseMessage)
	}

	discordSendMessagesResponse(responseMessage, s, i)

}

func discordSendMessagesResponse(responseMessage string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Send a message to the channel
	// Handle the "test" slash command
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
		},
	}

	// Send the response back to the user
	err := s.InteractionRespond(i.Interaction, response)
	catchErr(err)
}

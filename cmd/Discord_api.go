package cmd

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var discord *discordgo.Session

func initDiscord() {
	var err error
	discord, err = discordgo.New("Bot " + DISCORD_API_KEY)
	catchErr(err)
}

func discord_GetImage() {

	initDiscord()

	messages, err := discord.ChannelMessages(discord_channel_midjourney, 10, "", "", "")
	catchErr(err)

	for _, v := range messages {
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

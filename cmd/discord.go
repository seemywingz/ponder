/*
Copyright © 2023 Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// discordCmd represents the discord command
var discordCmd = &cobra.Command{
	Use:   "discord-bot",
	Short: "Discord Chat Bot Integration",
	Long:  `Discord Chat Bot Integration utilizing Secure Gateway Websocket`,
	Run: func(cmd *cobra.Command, args []string) {
		initDiscord()
	},
}

func init() {
	rootCmd.AddCommand(discordCmd)
}

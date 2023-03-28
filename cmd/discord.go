/*
Copyright © 2023 Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// discordCmd represents the discord command
var discordCmd = &cobra.Command{
	Use:   "discord",
	Short: "Discord Chat Bot Integration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		discord_GetImage()
	},
}

func init() {
	rootCmd.AddCommand(discordCmd)
}

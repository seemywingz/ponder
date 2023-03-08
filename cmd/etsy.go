/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// etsyCmd represents the etsy command
var etsyCmd = &cobra.Command{
	Use:   "etsy",
	Short: "Make Etsy API calls",
	Long:  `Make calls to the Etsy API.`,
	Run: func(cmd *cobra.Command, args []string) {
		etsy_CreateRequest()
	},
}

func init() {
	rootCmd.AddCommand(etsyCmd)
}

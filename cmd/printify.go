/*
Copyright © 2023 Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// printifyCmd represents the printify command
var printifyCmd = &cobra.Command{
	Use:   "printify",
	Short: "Interact with the Printify API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		generatImageAndPost()
	},
}

func init() {
	rootCmd.AddCommand(printifyCmd)
}

func generatImageAndPost() {

	prompt := "test image under construction"

	// Generate Image
	fmt.Println("🖼  Creating Image...")
	res := openAI_ImageGen(prompt, "", 1)
	fmt.Println("🌐  Image URL", res.Data[0].URL)

	prompt = strings.ReplaceAll(prompt, " ", "_")
	prompt = strings.ReplaceAll(prompt, "/", "-")

	// Create Printify Product
	fmt.Println("📦  Creating Printify Product...")
	printify_UploadImage(prompt+".jpg", res.Data[0].URL)
}

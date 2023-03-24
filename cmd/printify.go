/*
Copyright © 2023 Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"fmt"
	"strconv"

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
	printifyCmd.Flags().IntVarP(&n, "num-generated", "n", 1, "Number of images to generate")
}

func generatImageAndPost() {

	// Generate Image
	fmt.Println("🖼  Generating Image(s)...")
	res := openAI_ImageGen(prompt, "", n)

	for imgNum, data := range res.Data {
		url := data.URL
		// Format Prompt for use as Product Name
		fileName := formatPrompt(prompt)
		// Create Printify Product
		fmt.Println()
		fmt.Println("📦  Creating Printify Product...")
		fmt.Println("🌐 Image URL: " + url)
		printify_UploadImage(fileName+"_"+strconv.Itoa(imgNum)+".jpg", url)
	}

}

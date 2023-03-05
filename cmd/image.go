/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var open bool
var imageFile string
var n int

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate an image from a prompt",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		createImage(prompt, imageFile)

	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.Flags().BoolVarP(&open, "open", "o", false, "Open image in browser")
	imageCmd.Flags().StringVarP(&imageFile, "image", "i", "", "Image file to be used as prompt")
	imageCmd.Flags().IntVarP(&n, "n", "n", 1, "Number of images to generate")
}

func createImage(prompt, imageFile string) {
	fmt.Println("🖼  Creating Image...")
	res := openAI_ImageGen(prompt, imageFile, n)

	for _, data := range res.Data {
		url := data.URL
		fmt.Println("🌐 Image URL: " + url)

		err := error(nil)
		if open { // Open image in browser if open flag is set
			fmt.Println("💻 Opening Image URL...")
			switch runtime.GOOS {
			case "linux":
				err = exec.Command("xdg-open", url).Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
			case "darwin":
				err = exec.Command("open", url).Start()
			default:
				err = fmt.Errorf("unsupported platform")
			}
			if err != nil {
				trace()
				fmt.Println(err)
			}
		}
	}
}

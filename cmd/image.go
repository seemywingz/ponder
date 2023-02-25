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

var prompt string

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate an image from a prompt",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		createImage(prompt)
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "Prompt for image generation")
	imageCmd.MarkFlagRequired("prompt")
}

func createImage(prompt string) {
	fmt.Println("Creating Image...")
	res, err := getImageResponse(prompt)
	catchErr(err)
	url := res.Data[0].URL
	fmt.Println("Opening Image URL: " + url)
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

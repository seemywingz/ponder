/*
Copyright © 2023 NAME HERE Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
)

var filePath = "HOME"
var open, download bool
var file string
var n int

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate an image from a prompt",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		createImage(prompt, file)

	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.Flags().BoolVarP(&download, "download", "d", false, "Download image(s) to local directory")
	imageCmd.Flags().BoolVarP(&open, "open", "o", false, "Open image in browser")
	imageCmd.Flags().IntVarP(&n, "n", "n", 1, "Number of images to generate")
	imageCmd.Flags().StringVarP(&file, "file", "f", "", "Image file to edit")
}

func createImage(prompt, imageFile string) {
	fmt.Println("🖼  Creating Image...")
	res := openAI_ImageGen(prompt, imageFile, n)

	for imgNum, data := range res.Data {
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
		if download { // Download image to local directory if download flag is set
			if filePath == "HOME" { // If no path is specified, use the user's home directory
				currentUser, err := user.Current()
				catchErr(err)
				filePath = currentUser.HomeDir + "/Ponder/Images"
			}
			fileName := prompt + "_" + strconv.Itoa(imgNum) + ".jpg"
			fullFilePath := filepath.Join(filePath, fileName)
			// Create the directory (if it doesn't exist)
			err := os.MkdirAll(filePath, os.ModePerm)
			catchErr(err)
			fmt.Println("📥 Downloading Image...", fullFilePath)
			httpDownloadFile(url, fullFilePath)
		}
	}
}

package cmd

import (
	"log"
	"strconv"

	"github.com/disintegration/imaging"
)

// UpscaleImage upscales the image at the given path by the given scale factor.
// It saves the upscaled image to a new file and returns the path of the new file.
func UpscaleImage(imagePath string, scaleFactor int) string {
	// Open the original image.
	src, err := imaging.Open(imagePath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// Upscale the image.
	dstImage := imaging.Resize(src, src.Bounds().Dx()*scaleFactor, src.Bounds().Dy()*scaleFactor, imaging.Lanczos)

	// Create the path for the new file.
	newImagePath := "upscaled_" + strconv.Itoa(scaleFactor) + "x_" + imagePath

	// Save the resulting image as JPEG.
	err = imaging.Save(dstImage, newImagePath)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	// Return the path of the new file.
	return newImagePath
}

func main() {
	// Test the UpscaleImage function.
	newImagePath := UpscaleImage("test.jpg", 2)
	log.Println("Saved upscaled image to:", newImagePath)
}

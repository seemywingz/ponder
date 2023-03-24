package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func printify_UploadImage(fileName, imageURL string) {
	// Create the JSON Request and Response Objects
	responseJson := PRINTIFY_ImageUploadResponse{}
	requestJson := PRINTIFY_ImageUploadRequest{
		FileName: fileName,
		URL:      imageURL,
	}

	// Marshal the JSON Request Body
	requestBody, err := json.Marshal(requestJson)
	catchErr(err)

	// Create the HTTP Request
	req, err := http.NewRequest("POST", printify_endpoint+"uploads/images.json", bytes.NewBuffer(requestBody))
	catchErr(err)

	// Set the HTTP Request Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+PRINTIFY_API_KEY)

	// Verbose Output
	if verbose {
		trace()
		fmt.Println("🗂️  Request JSON:", requestJson)
		fmt.Println("🌐 HTTP Request", req)
	}

	// Make the HTTP Request
	httpMakeRequest(req, &responseJson)

	fmt.Println("Image Uploaded to Printify")
	fmt.Println("📁 ID:", responseJson.ID)
	fmt.Println("📁 Name:", responseJson.FileName)
	fmt.Println("📁 Height:", responseJson.Height)
	fmt.Println("📁 Width:", responseJson.Width)
	fmt.Println("📁 Size:", responseJson.Size)
	fmt.Println("📁 MimeType:", responseJson.MimeType)
	fmt.Println("📁 PreviewURL:", responseJson.PreviewURL)
	fmt.Println("📁 UploadTime:", responseJson.UploadTime)

}

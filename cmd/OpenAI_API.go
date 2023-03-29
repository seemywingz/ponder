package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func openAI_UploadImage(requestJson, responseJson interface{}, endpoint, filePath string) {

	// Get the absolute path of the file
	fullPath, err := filepath.Abs(filePath)
	catchErr(err)

	// https://platform.openai.com/docs/api-reference/images/create-edit#images/create-edit-image
	// The image to edit. Must be a valid PNG file, less than 4MB, and square.
	// If mask is not provided, image must have transparency, which will be used as the mask.
	//
	// Open the PNG image file
	file, err := os.Open(fullPath)
	catchErr(err)
	defer file.Close()

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the PNG file to the request
	part, err := writer.CreateFormFile("image", filePath)
	catchErr(err)
	_, err = io.Copy(part, file)
	catchErr(err)

	oaiImageEditJson := requestJson.(*OPENAI_ImageEditRequest)

	// Add the JSON payload to the request
	part, err = writer.CreateFormField("prompt")
	catchErr(err)
	part.Write([]byte(oaiImageEditJson.Prompt))

	part, err = writer.CreateFormField("n")
	catchErr(err)
	part.Write([]byte(strconv.Itoa(oaiImageEditJson.N)))

	part, err = writer.CreateFormField("size")
	catchErr(err)
	part.Write([]byte(oaiImageEditJson.Size))

	part, err = writer.CreateFormField("user")
	catchErr(err)
	part.Write([]byte(oaiImageEditJson.User))

	// Close the multipart writer
	err = writer.Close()
	catchErr(err)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", endpoint, body)
	catchErr(err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	if verbose {
		trace()
		fmt.Println("Request Body: ", req.Body)
		fmt.Println("Request JSON: ", oaiImageEditJson)
	}

	// Send the request
	fmt.Println("⏳ Uploading File: " + fullPath)
	resp, err := httpClient.Do(req)
	catchErr(err)

	// Read the JSON Response Body
	jsonString, err := io.ReadAll(resp.Body)
	catchErr(err)

	// Check for API Errors
	httpCatchErr(resp, jsonString)

	// Unmarshal the JSON Response Body
	err = json.Unmarshal([]byte(jsonString), &responseJson)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(jsonString))
	}

	// Close the HTTP Response Body
	defer resp.Body.Close()
}

func openAI_ImageGen(prompt, imageFile string, n int) OPENAI_ImageResponse {
	var oaiRequest interface{}
	oaiResponse := OPENAI_ImageResponse{}

	// Check if we are editing an image or generating a new one
	if imageFile != "" {

		// Create the JSON Request Body
		oaiRequest = &OPENAI_ImageEditRequest{
			Prompt:         prompt,
			N:              n,
			Size:           "1024x1024",
			ResponseFormat: "url",
			User:           openAIUser,
		}
		openAI_UploadImage(oaiRequest, &oaiResponse, openai_endpoint+"images/edits", imageFile)

	} else { // Generate a new image

		oaiRequest = &OPENAI_ImageRequest{
			Prompt:         prompt,
			N:              n,
			Size:           "1024x1024",
			ResponseFormat: "url",
			User:           openAIUser,
		}
		openAI_PostJson(oaiRequest, &oaiResponse, openai_endpoint+"images/generations")
	}
	if verbose {
		trace()
		fmt.Println(oaiRequest)
	}

	return oaiResponse
}

func openai_ChatComplete(messages []OPENAI_Message) OPENAI_ChatCompletionResponse {
	oaiResponse := OPENAI_ChatCompletionResponse{}
	oaiRequest := OPENAI_ChatCompletionRequest{
		Model:            "gpt-3.5-turbo",
		Messages:         messages,
		N:                1,
		Temperature:      0,
		TopP:             0.1,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.6,
		MaxTokens:        1000,
		User:             openAIUser,
	}
	openAI_PostJson(oaiRequest, &oaiResponse, openai_endpoint+"chat/completions")
	return oaiResponse
}

func openAI_Chat(prompt string) OPENAI_ChatResponse {
	oaiResponse := OPENAI_ChatResponse{}
	oaiRequest := &OPENAI_ChatRequest{
		Prompt:           prompt,
		MaxTokens:        1000,
		Model:            "text-davinci-003",
		Temperature:      0,
		TopP:             0.1,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.6,
		User:             openAIUser,
	}
	if verbose {
		trace()
		fmt.Println(oaiRequest)
	}
	openAI_PostJson(oaiRequest, &oaiResponse, openai_endpoint+"completions")
	return oaiResponse
}

func openAI_PostJson(requestJson, responseJson interface{}, endpoint string) {

	// Marshal the JSON Request Body
	requestBodyJson, err := json.Marshal(requestJson)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(requestBodyJson))
	}

	// Format HTTP Response and Set Headers
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBodyJson))
	catchErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	httpMakeRequest(req, responseJson)
}

func openai_ChatTXTonly(prompt string) string {
	oaiResponse := openAI_Chat(prompt)
	responseMessage := ""
	for _, v := range oaiResponse.Choices {
		responseMessage += v.Text[2:]
	}
	return responseMessage
}

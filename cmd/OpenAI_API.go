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

	"github.com/spf13/viper"
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
			N:              n,
			ResponseFormat: "url",
			Prompt:         prompt,
			User:           openAIUser,
			Size:           viper.GetString("openai.image.size"),
		}
		openAI_UploadImage(oaiRequest, &oaiResponse, viper.GetString("openai.endpoint")+"images/edits", imageFile)

	} else { // Generate a new image

		oaiRequest = &OPENAI_ImageRequest{
			N:              n,
			ResponseFormat: "url",
			Prompt:         prompt,
			User:           openAIUser,
			Size:           viper.GetString("openai.image.size"),
		}
		openAI_PostJson(oaiRequest, &oaiResponse, viper.GetString("openai.endpoint")+"images/generations")
	}
	if verbose {
		trace()
		fmt.Println(oaiRequest)
	}
	return oaiResponse
}

func openai_ChatCompletion(messages []OPENAI_Message) string {
	oaiResponse := OPENAI_ChatCompletionResponse{}
	oaiRequest := OPENAI_ChatCompletionRequest{
		N:                1,
		Messages:         messages,
		User:             openAIUser,
		TopP:             viper.GetFloat64("openai.completion.chat.topP"),
		Model:            viper.GetString("openai.completion.chat.model"),
		MaxTokens:        viper.GetInt("openai.completion.chat.maxTokens"),
		Temperature:      viper.GetFloat64("openai.completion.chat.temperature"),
		FrequencyPenalty: viper.GetFloat64("openai.completion.chat.frequencyPenalty"),
		PresencePenalty:  viper.GetFloat64("openai.completion.chat.presencePenalty"),
	}
	openAI_PostJson(oaiRequest, &oaiResponse, viper.GetString("openai.endpoint")+"chat/completions")
	return oaiResponse.Choices[0].Message.Content
}

func openAI_TextCompletion(prompt string) OPENAI_ChatResponse {
	oaiResponse := OPENAI_ChatResponse{}
	oaiRequest := &OPENAI_ChatRequest{
		Prompt:           prompt,
		User:             openAIUser,
		Model:            viper.GetString("openai.completion.text.model"),
		MaxTokens:        viper.GetInt("openai.completion.text.maxTokens"),
		Temperature:      viper.GetFloat64("openai.completion.text.temperature"),
		TopP:             viper.GetFloat64("openai.completion.text.topP"),
		FrequencyPenalty: viper.GetFloat64("openai.completion.text.frequencyPenalty"),
		PresencePenalty:  viper.GetFloat64("openai.completion.text.presencePenalty"),
	}
	if verbose {
		trace()
		fmt.Println(oaiRequest)
	}
	openAI_PostJson(oaiRequest, &oaiResponse, viper.GetString("openai.endpoint")+"completions")
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

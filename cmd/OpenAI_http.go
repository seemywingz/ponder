package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func openAI_API(requestJson, responseJson interface{}, endpoint string) {

	// Marshal the JSON Request Body
	requestBodyJson, err := json.Marshal(requestJson)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(requestBodyJson))
	}

	// Create HTTP Client
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}

	// Format HTTP Response and Set Headers
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBodyJson))
	catchErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	// Make the HTTP Request
	resp, err := httpClient.Do(req)
	catchErr(err)

	// Read the JSON Response Body
	jsonString, err := io.ReadAll(resp.Body)
	catchErr(err)

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

func openAI_ImageGen(prompt, imageFile string) OPENAI_ImageResponse {
	var oaiRequest interface{}
	oaiResponse := OPENAI_ImageResponse{}
	endpoint := enpoint_openai

	if imageFile != "" {
		endpoint += "images/edits"
		oaiRequest = &OPENAI_ImageEditRequest{
			Prompt:         prompt,
			N:              1,
			Size:           "1024x1024",
			ResponseFormat: "url",
			Mask:           "",
			Image:          imageFile,
		}

	} else {
		endpoint += "images/generations"
		oaiRequest = &OPENAI_ImageRequest{
			Prompt:         prompt,
			N:              1,
			Size:           "1024x1024",
			ResponseFormat: "url",
		}
	}
	if verbose {
		trace()
		fmt.Println(oaiRequest)
	}

	openAI_API(oaiRequest, &oaiResponse, endpoint)
	return oaiResponse
}

func openAI_Chat(prompt string) OPENAI_ChatResponse {
	oaiResponse := OPENAI_ChatResponse{}
	endpoint := enpoint_openai + "completions"

	oaiRequest := &OPENAI_ChatRequest{
		Prompt:           prompt,
		MaxTokens:        1000,
		Model:            "text-davinci-003",
		Temperature:      0,
		TopP:             0.1,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.6,
		User:             "ponder",
	}

	openAI_API(oaiRequest, &oaiResponse, endpoint)
	return oaiResponse
}

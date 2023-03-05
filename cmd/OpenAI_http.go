package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func openAI_API(requestJson, responseJson interface{}, endpoint string) {
	requestBodyJson, err := json.Marshal(requestJson)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(requestBodyJson))
	}

	apiKeyHeader := "Bearer " + OPENAI_API_KEY
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBodyJson))
	catchErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKeyHeader)

	// Make the HTTP Request
	resp, err := httpClient.Do(req)
	catchErr(err)

	// read the JSON Body
	jsonString, err := io.ReadAll(resp.Body)
	catchErr(err)

	err = json.Unmarshal([]byte(jsonString), &responseJson)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(jsonString))
	}
	defer resp.Body.Close()
}

func openAI_ImageGen(prompt, imageFile string) OPENAI_ImageResponse {
	var oaiRequest interface{}
	oaiResponse := OPENAI_ImageResponse{}
	endpoint := enpoint_openai

	if imageFile != "" {
		endpoint += "images/edits"
		oaiRequest = &ImageEditRequest{
			Prompt:         prompt,
			N:              1,
			Size:           "1024x1024",
			ResponseFormat: "url",
			Mask:           "",
			Image:          imageFile,
		}

	} else {
		endpoint += "images/generations"
		oaiRequest = &ImageRequest{
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

func getChatResponse(prompt string) (ChatResponse, error) {
	chatResponse := ChatResponse{}
	chatRequest := &ChatRequest{
		Prompt:           prompt,
		MaxTokens:        1000,
		Model:            "text-davinci-003",
		Temperature:      0,
		TopP:             0.1,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.6,
		User:             "ponder" + os.Getenv("OPENAI_API_KEY"),
	}
	requestBodyJson, err := json.Marshal(chatRequest)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(requestBodyJson))
	}

	resp, err := http_POST(requestBodyJson, enpoint_openai+"completions", OPENAI_API_KEY)
	catchErr(err)

	jsonString, err := io.ReadAll(resp.Body)
	catchErr(err)

	err = json.Unmarshal([]byte(jsonString), &chatResponse)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(jsonString))
	}
	defer resp.Body.Close()
	return chatResponse, nil
}

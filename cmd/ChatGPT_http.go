package cmd

import (
	"encoding/json"
	"fmt"
	"io"
)

func getImageResponse(prompt string) (ImageResponse, error) {
	imageResponse := ImageResponse{}
	imageRequest := &ImageRequest{
		Prompt:         prompt,
		N:              1,
		Size:           "1024x1024",
		ResponseFormat: "url",
	}

	requestBodyJson, err := json.Marshal(imageRequest)
	catchErr(err)

	if verbose {
		trace()
		fmt.Println(string(requestBodyJson))
	}
	resp, err := getResponse(requestBodyJson, enpoint_openai+"images/generations")
	catchErr(err)

	jsonString, err := io.ReadAll(resp.Body)
	catchErr(err)

	err = json.Unmarshal([]byte(jsonString), &imageResponse)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(jsonString))
	}
	defer resp.Body.Close()
	return imageResponse, nil
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
		User:             "ponder",
	}
	requestBodyJson, err := json.Marshal(chatRequest)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(requestBodyJson))
	}

	resp, err := getResponse(requestBodyJson, enpoint_openai+"completions")
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

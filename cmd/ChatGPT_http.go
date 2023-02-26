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
	resp, err := getResponse(requestBodyJson, "images/generations")
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
	}
	requestBodyJson, err := json.Marshal(chatRequest)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(requestBodyJson))
	}

	resp, err := getResponse(requestBodyJson, "completions")
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

func getResponse(requestBodyJson []byte, endpoint string) (*http.Response, error) {

	apiKey := "Bearer " + os.Getenv("OPENAI_API_KEY")
	requestUrl := "https://api.openai.com/v1/" + endpoint

	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(requestBodyJson))
	catchErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKey)
	resp, err := httpClient.Do(req)
	catchErr(err)
	return resp, nil
}

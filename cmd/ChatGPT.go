package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func getUserInput() (string, error) {
	fmt.Println("You:")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\n")
	return input, nil
}

func getChatResponse(prompt string) (ChatResponse, error) {
	var chatResponse ChatResponse
	chatRequest := &ChatRequest{
		Prompt:           prompt,
		MaxTokens:        1000,
		Model:            "text-davinci-003",
		Temperature:      0.5,
		TopP:             1.0,
		FrequencyPenalty: 0.5,
		PresencePenalty:  0.0,
	}
	requestBodyJson, err := json.Marshal(chatRequest)
	if err != nil {
		return chatResponse, err
	}
	apiKey := "Bearer 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"
	requestUrl := "https://api.openai.com/v1/completions"

	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}
	request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return chatResponse, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", apiKey)
	resp, err := httpClient.Do(request)
	if err != nil {
		return chatResponse, err
	}

	jsonString, err := io.ReadAll(resp.Body)
	if err != nil {
		return chatResponse, err
	}

	err = json.Unmarshal([]byte(jsonString), &chatResponse)
	if err != nil {
		return chatResponse, err
	}

	defer resp.Body.Close()
	return chatResponse, nil
}

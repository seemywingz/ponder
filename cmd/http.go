package cmd

import (
	"bytes"
	"net/http"
	"os"
	"time"
)

func getResponse(requestBodyJson []byte, endpoint string) (*http.Response, error) {

	apiKey := "Bearer " + os.Getenv("OPENAI_API_KEY")

	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBodyJson))
	catchErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKey)
	resp, err := httpClient.Do(req)
	catchErr(err)
	return resp, nil
}

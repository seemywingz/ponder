package cmd

import (
	"bytes"
	"net/http"
	"time"
)

func http_POST(requestBodyJson []byte, endpoint, apiKey string) (*http.Response, error) {
	apiKeyHeader := "Bearer " + apiKey
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBodyJson))
	catchErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKeyHeader)

	resp, err := httpClient.Do(req)
	catchErr(err)
	return resp, nil
}

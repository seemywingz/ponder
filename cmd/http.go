package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Create HTTP Client
var httpClient = &http.Client{
	Timeout: time.Second * 60,
}

func httpMakeRequest(request *http.Request, responseJson interface{}) {

	// Make the HTTP Request
	resp, err := httpClient.Do(request)
	catchErr(err)

	// Read the JSON Response Body
	jsonString, err := io.ReadAll(resp.Body)
	catchErr(err)

	// Check for HTTP Errors
	httpCatchErr(resp, jsonString)

	// Unmarshal the JSON Response Body into provided responseJson
	err = json.Unmarshal([]byte(jsonString), &responseJson)
	catchErr(err)
	if verbose {
		trace()
		fmt.Println(string(jsonString))
	}
	// Close the HTTP Response Body
	defer resp.Body.Close()
}

func httpCatchErr(resp *http.Response, jsonString []byte) {
	// Check for HTTP Response Errors
	if resp.StatusCode != 200 {
		catchErr(errors.New("API Error: " + strconv.Itoa(resp.StatusCode) + "\n" + string(jsonString)))
	}
}

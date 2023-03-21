package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

// download file from url and save to local directory
func httpDownloadFile(url string, filePath string) {

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File already exists, so rename the new file
		dir := filepath.Dir(filePath)
		ext := filepath.Ext(filePath)
		name := filepath.Base(filePath[:len(filePath)-len(ext)])
		i := 1
		for {
			newName := fmt.Sprintf("%s (%d)%s", name, i, ext)
			newFilepath := filepath.Join(dir, newName)
			_, err := os.Stat(newFilepath)
			if os.IsNotExist(err) {
				// New filename is available, use it
				filePath = newFilepath
				break
			}
			i++
		}
	}

	// Create the file
	out, err := os.Create(filePath)
	catchErr(err)
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	catchErr(err)
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	catchErr(err)
}

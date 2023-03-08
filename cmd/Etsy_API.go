package cmd

import (
	"fmt"
	"net/http"
)

func etsy_CreateRequest() {

	etsyResponseJson := &Etsy_Response{}

	// Create the Etsy API Request
	endpoint := etsy_endpoint + "application/listings/active" // Create HTTP Request
	request, err := http.NewRequest("GET", endpoint, nil)
	catchErr(err)

	// Add Headers
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-api-key", ETSY_API_KEY)

  if verbose {
		trace()
		fmt.Println(request)
	}

	// Make the HTTP Request
	httpMakeRequest(request, &etsyResponseJson)
	fmt.Println(etsyResponseJson)

}

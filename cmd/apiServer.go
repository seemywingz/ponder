/*
Copyright © 2023 Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var ponder_api_version = "v1"
var ponder_api_port = "8080"

// apiServerCmd represents the apiServer command
var apiServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Start the API Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		apiServer()
	},
}

func init() {
	rootCmd.AddCommand(apiServerCmd)
}

// Start the Gorilla MUX API Server
func apiServer() {
	// Create a new router
	r := mux.NewRouter()

	// Add routes
	r.HandleFunc("/api/"+ponder_api_version+"/discord", discordHandler).Methods("POST")

	// Start the server
	fmt.Println("Starting API Server on port", ponder_api_port, "...")
	err := http.ListenAndServe(":"+ponder_api_port, r)
	catchErr(err)

}

// Discord Handler
func discordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Discord Handler")
	// Get the Request Body
	fmt.Println("Request Body:", r)
}

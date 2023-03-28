/*
Copyright © 2023 Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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

	// port flag
	apiServerCmd.Flags().StringVarP(&ponder_api_port, "port", "P", "8080", "Port to run the API Server on")
}

// Start the Gorilla MUX API Server
func apiServer() {
	// Create a new router
	r := mux.NewRouter()

	// add liveness and readiness probes
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add routes
	r.HandleFunc("/api/"+ponder_api_version+"/discord", discordHandler).Methods("POST")

	// Start the server
	fmt.Println("Starting API Server on port", ponder_api_port, "...")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	err := http.ListenAndServe(":"+ponder_api_port, loggedRouter)
	catchErr(err)

}

// Discord Handler
func discordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Discord Handler")
	fmt.Println("Request Body:", r)
	log.Println("Processing request!")
}

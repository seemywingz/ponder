/*
Copyright © 2023 Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
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
		initDiscord()
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

func discordValidateRequest(w http.ResponseWriter, r *http.Request) bool {
	// Decode the hex string into bytes
	bytes, err := hex.DecodeString(DISCORD_PUB_KEY)
	catchErr(err)
	// Convert bytes to ed25519.PublicKey
	var publicKey ed25519.PublicKey = bytes
	// Validate the request using the Discord Go library
	return discordgo.VerifyInteraction(r, publicKey)
}

// Discord Handler
func discordHandler(w http.ResponseWriter, r *http.Request) {
	request := discordgo.Webhook{}

	if verbose {
		trace()
		httpDumpRequest(r)
	}

	if !discordValidateRequest(w, r) {
		http.Error(w, "Invalid Signature", http.StatusUnauthorized)
		fmt.Println("Discord Handler: Invalid Signature")
		return
	}

	// Read the JSON
	reqJson, err := io.ReadAll(r.Body)
	catchErr(err)

	if len(reqJson) <= 0 {
		fmt.Fprintf(w, "Discord Handler: No JSON received")
		return
	}

	// Unmarshal the JSON
	err = json.Unmarshal([]byte(reqJson), &request)
	catchErr(err)

	if request.Type == 1 {
		fmt.Println("Discord Ping")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(reqJson)
	}

}

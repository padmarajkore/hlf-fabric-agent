package main

import (
	"log"
	"net/http"

	"hlf-controller/internal/config"
	"hlf-controller/internal/handlers"
)

func main() {
	// Initialize config on startup. This will check for fabric-samples and download them if needed.
	// The server will not start if this fails.
	log.Println("[INFO] Initializing HLF Controller...")
	config.LoadConfig()
	log.Println("[SUCCESS] Initialization complete. Configuration loaded and prerequisites checked.")

	http.HandleFunc("/network/up", handlers.UpNetworkHandler)
	http.HandleFunc("/network/down", handlers.DownNetworkHandler)
	http.HandleFunc("/chaincode/deploy", handlers.DeployChaincodeHandler)
	http.HandleFunc("/chaincode/invoke", handlers.InvokeChaincodeHandler)
	http.HandleFunc("/chaincode/query", handlers.QueryChaincodeHandler)
	http.HandleFunc("/channel/create", handlers.CreateChannelHandler)
	log.Println("[INFO] HLF Controller server running on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

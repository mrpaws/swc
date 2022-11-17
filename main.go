package main

// strava activity webhook subscription callback endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("starting swc athlete event receiver server...")
	http.HandleFunc("/swc/receiver", SubscriptionEventHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func SubscriptionEventHandler(w http.ResponseWriter, r *http.Request) {
	/*
		handles events from strava webhook subscription
		https://developers.strava.com/docs/webhooks/
	*/
	query := r.URL.Query()

	if hubMode := query.Get("hub.mode"); hubMode == "subscribe" {
		hubChallenge := query.Get("hub.challenge")
		hubVerifyToken := query.Get("hub.verify_token")
		log.Printf("Handling subscription creation query (%s,%s)", hubChallenge, hubVerifyToken)

		header := w.Header()
		header.Set("Content-Type", "application/json")
		responseBody := make(map[string]string)
		responseBody["hub.challenge"] = hubChallenge
		b, err := json.Marshal(responseBody)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s", b)
	}

}

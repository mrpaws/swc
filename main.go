package main

/* strava activity webhook subscription callback endpoint manager
   ref: https://developers.strava.com/docs/webhooks/
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("starting swc athlete event receiver server...")
	http.HandleFunc("/swc/receiver", SubscriptionEventRouter)

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

func SubscriptionEventRouter(w http.ResponseWriter, r *http.Request) {
	/*
		handles events from strava webhook subscription
		https://developers.strava.com/docs/webhooks/
	*/
	method := r.Method

	// strava subscription service callback events
	if method == "GET" {
		SubscriptionEventHandler(w, r)
	} else if method == "POST" {
		AthleteActivityEventHandler(w, r)
	}
}

func SubscriptionEventHandler(w http.ResponseWriter, r *http.Request) {
	/* handler for receiving subscription event
	i.e. creation and callback challenge / token verification
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
func AthleteActivityEventHandler(w http.ResponseWriter, r *http.Request) {
	/* handler for receiving event about athlete or activity

		   https://developers.strava.com/docs/webhooks/
	       - 200 OK + async processing recommended;
		   - multiple fires per update possible;
		   - multiple values in updates kv

		   POST payload structure i.e.:
		     {
				"aspect_type": "update",
				"event_time": 1516126040,
				"object_id": 1360128428,
				"object_type": "activity",
				"owner_id": 134815,
				"subscription_id": 120475,
				"updates": {
					"title": "Messy"
			}
	*/
	r.ParseForm()
	log.Println("handling athlete or activity update:", r.PostForm)
	fmt.Fprintf(w, "")
}

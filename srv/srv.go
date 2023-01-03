package srv

/* strava activity webhook subscription callback endpoint manager
   ref: https://developers.strava.com/docs/webhooks/
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mrpaws/swc/lib"
)

func Server() {
	swcVerifyToken := lib.GetSwcVerifyToken()

	if swcSubscriptionId := lib.GetSwcSubscriptionId(); swcSubscriptionId == "" {
		log.Print("no subscription id set for server at startup")
	}

	log.Print("starting swc athlete event receiver server...")
	http.HandleFunc("/swc/receiver", subscriptionEventRouter(swcVerifyToken.String()))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("swc listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func subscriptionEventRouter(swcVerifyToken string) http.HandlerFunc {
	/*
		handles events from strava webhook subscription
		https://developers.strava.com/docs/webhooks/
	*/
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method

		// strava subscription service callback events
		if method == "GET" {
			subscriptionEventHandler(w, r, swcVerifyToken)
		} else if method == "POST" {
			athleteActivityEventHandler(w, r)
		}
	}
}

func subscriptionEventHandler(w http.ResponseWriter, r *http.Request, swcVerifyToken string) {
	/* handler for receiving subscription event
	i.e. creation and callback challenge / token verification
	*/
	query := r.URL.Query()
	
	if hubMode := query.Get("hub.mode"); hubMode == "subscribe" {
		hubChallenge := query.Get("hub.challenge")
		hubVerifyToken := query.Get("hub.verify_token")
		log.Printf("handling subscription creation query from (%s) with challengeL (%s)",
			r.Host, hubChallenge)
		if hubVerifyToken == swcVerifyToken {
			log.Print("provided verify token matches our verify token")
			header := w.Header()
			header.Set("Content-Type", "application/json")
			responseBody := make(map[string]string)
			responseBody["hub.challenge"] = hubChallenge
			b, err := json.Marshal(responseBody)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s", b)
		} else {
			log.Printf("unauthorized subscription attempt: %s ", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func athleteActivityEventHandler(w http.ResponseWriter, r *http.Request) {
	/* handlerfunc closure for receiving event about athlete or activity

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

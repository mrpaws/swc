package clt

/* strava activity webhook subscription client operation jobs
	create
	view
	delete
   ref: https://developers.strava.com/docs/webhooks/
*/

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mrpaws/swc/lib"
)

const (
	swcSubscriptionEndpoint = "https://www.strava.com/api/v3/push_subscriptions"
)

var (
	swcClientId string
	swcCallbackUrl string
	swcVerifyToken string
	swcClientSecret string
	swcClientOp string
	swcSubscriptionId string
	swcEventData map[string]string
)

func main() {

	log.Print("starting strava webhook callback client job")

	// switch better
	if swcClientOp = lib.LoadEnvVar("swcClientOp", true); swcClientOp == "" {
		log.Fatal("unspecified operation not yet supported (env var: swcClientOp")
	} else if swcClientOp == "create" {
		log.Print("attempting to create subscription...")
		swcClientId = lib.LoadEnvVar("swcClientId", true)
		swcCallbackUrl = lib.LoadEnvVar("swcCallbackUrl", true)
		swcVerifyToken = lib.LoadEnvVar("token", true)
		swcClientSecret = lib.LoadEnvVar("swcClientSecret", true)
		
		swcSubscriptionId = createSubscription(
			swcClientId, 
			swcClientSecret, 
			swcCallbackUrl, 
			swcVerifyToken,
		)

	} else if swcClientOp == "view" {
		log.Fatal("view is not yet implemented")
	}

	fmt.Printf("swcSubscriptionId: %v", swcSubscriptionId)
	log.Printf("ending strava webhook callback client job type: (%s)", swcClientOp)

}

func createSubscription(swcClientId string, swcClientSecret string, swcCallbackUrl string, swcVerifyToken string) string {
	//var respSize int
	var respBody []byte
	postData := url.Values{}
	postData.Set("client_id", swcClientId)
	postData.Set("client_secret", swcClientSecret)
	postData.Set("callback_url", swcCallbackUrl)
	postData.Set("verify_token", swcVerifyToken)
	resp, err := http.PostForm(swcSubscriptionEndpoint, postData)
	if err != nil {
		log.Fatal("error with http request to create subscription")
	}
	respSize, err := resp.Body.Read(respBody)
	if err != nil {
		log.Fatal("error reading http response body")
	}
	log.Printf("read %d bytes response for subscription creation: (%s)", respSize, respBody)
	return string(respBody)
}
